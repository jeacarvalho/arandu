package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/infrastructure/ai"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/platform/logger"
	"arandu/internal/platform/middleware"
	"arandu/internal/web"
	"arandu/internal/web/handlers"
	"arandu/web/components/dashboard"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found or error loading")
	}

	// Log startup information (version, commit e build_time são adicionados automaticamente pelo logger)
	logger.Info("Starting Arandu")

	// Initialize Control Plane (Central DB)
	logger.Info("Initializing Control Plane (Central DB)")
	centralDB, err := sqlite.NewCentralDB("storage")
	if err != nil {
		logger.Error("Failed to initialize central database", logger.String("error", err.Error()))
		os.Exit(1)
	}
	defer centralDB.Close()

	// Apply central migrations
	logger.Info("Applying central database migrations")
	if err := centralDB.Migrate(nil); err != nil {
		logger.Warn("Failed to apply central migrations", logger.String("error", err.Error()))
	}

	// Ensure tenants directory exists
	tenantsDir := "storage/tenants"
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		logger.Warn("Failed to create tenants directory", logger.String("error", err.Error()))
	} else {
		logger.Info("Tenants directory ready", logger.String("path", tenantsDir))
	}

	// Use central database for repository factory (multi-tenant data plane)
	db := &sqlite.DB{centralDB.DB}
	logger.Info("Using central database for data plane")

	// Initialize Tenant Pool for multi-tenant connections (must be before repositories)
	tenantPool := sqlite.NewTenantPool("storage", nil)
	logger.Info("Tenant pool initialized")

	// Create context-aware repository factory for multi-tenant support (clinical services)
	repoFactory := sqlite.NewContextAwareRepositoryFactory(db, tenantPool)
	patientRepo := sqlite.NewContextAwarePatientRepository(repoFactory)
	sessionRepo := sqlite.NewContextAwareSessionRepository(repoFactory)
	observationRepo := sqlite.NewContextAwareObservationRepository(repoFactory)
	interventionRepo := sqlite.NewContextAwareInterventionRepository(repoFactory)
	insightRepo := sqlite.NewContextAwareInsightRepository(repoFactory)
	medicationRepo := sqlite.NewContextAwareMedicationRepository(repoFactory)
	vitalsRepo := sqlite.NewContextAwareVitalsRepository(repoFactory)
	goalRepo := sqlite.NewContextAwareGoalRepository(repoFactory)
	appointmentRepo := sqlite.NewContextAwareAppointmentRepository(repoFactory)

	// Use context-aware repo for timeline
	timelineRepo := sqlite.NewContextAwareTimelineRepository(repoFactory)

	// Use context-aware repos for biopsychosocial (requires interface-based with context)
	biopsychosocialService := services.NewBiopsychosocialService(medicationRepo, vitalsRepo)

	// Create audit service (uses central DB)
	auditService := services.NewAuditService(centralDB.DB)
	defer auditService.Close()

	patientService := services.NewPatientServiceWithAudit(patientRepo, auditService)
	sessionService := services.NewSessionService(sessionRepo)
	observationService := services.NewObservationService(observationRepo)
	interventionService := services.NewInterventionService(interventionRepo)
	insightService := services.NewInsightService(insightRepo)
	timelineService := services.NewTimelineServiceContext(timelineRepo)

	// Create service adapters for the new handler interfaces
	sessionServiceAdapter := web.NewSessionServiceAdapter(sessionService)
	insightServiceAdapter := web.NewInsightServiceAdapter(insightService)
	patientServiceAdapter := web.NewPatientServiceAdapter(patientService)
	observationServiceAdapter := web.NewObservationServiceAdapter(observationService)
	interventionServiceAdapter := web.NewInterventionServiceAdapter(interventionService)
	timelineServiceAdapter := web.NewTimelineServiceAdapter(timelineService)
	goalServiceAdapter := web.NewGoalServiceAdapter(goalRepo)

	// Create biopsychosocial service adapter
	biopsychosocialServiceAdapterImpl := handlers.BiopsychosocialServiceFuncs{
		GetMedicationsFunc: func(ctx context.Context, patientID string) ([]interface{}, error) {
			meds, err := biopsychosocialService.GetMedications(ctx, patientID)
			if err != nil {
				return nil, err
			}
			// Convert to []interface{}
			result := make([]interface{}, len(meds))
			for i, m := range meds {
				result[i] = m
			}
			return result, nil
		},
		GetLatestVitalsFunc: func(ctx context.Context, patientID string) (interface{}, error) {
			return biopsychosocialService.GetLatestVitals(ctx, patientID)
		},
		GetAverageVitalsFunc: func(ctx context.Context, patientID string, days int) (interface{}, error) {
			return biopsychosocialService.GetAverageVitals(ctx, patientID, days)
		},
	}

	// Create anamnesis service adapter
	anamnesisServiceAdapter := web.NewAnamnesisServiceAdapter(patientRepo)

	// Create intervention classification service adapter
	interventionClassificationServiceAdapter := web.NewInterventionClassificationServiceAdapter(interventionRepo)

	// Create intervention classification handler
	interventionClassificationHandler := handlers.NewInterventionClassificationHandler(interventionClassificationServiceAdapter)

	// Create new handlers with dependency injection
	patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, insightServiceAdapter, biopsychosocialServiceAdapterImpl, timelineServiceAdapter, anamnesisServiceAdapter)
	sessionHandler := handlers.NewSessionHandler(sessionServiceAdapter, patientServiceAdapter, observationServiceAdapter, interventionServiceAdapter, goalServiceAdapter, observationServiceAdapter)
	observationHandler := handlers.NewObservationHandler(observationServiceAdapter)
	interventionHandler := handlers.NewInterventionHandler(interventionServiceAdapter)
	var dashboardHandler *handlers.DashboardHandler // initialized after agendaService below
	timelineHandler := handlers.NewTimelineHandler(timelineServiceAdapter)
	biopsychosocialHandler := handlers.NewBiopsychosocialHandler(biopsychosocialService)

	// Create Analysis Handler for theme cloud and pattern detection
	analysisHandler := handlers.NewAnalysisHandler(patientServiceAdapter, sessionServiceAdapter, observationServiceAdapter, interventionServiceAdapter, timelineServiceAdapter)

	// Initialize AI service
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		logger.Warn("GEMINI_API_KEY not set. AI features will be disabled.")
		geminiAPIKey = "dummy-key-for-initialization" // Use dummy key to allow initialization
	}

	geminiClient, err := ai.NewGeminiClient(geminiAPIKey)
	if err != nil {
		logger.Warn("Failed to initialize Gemini client", logger.String("error", err.Error()))
	}

	// Create cache for AI responses (24 hour TTL)
	cache := ai.NewCache(24 * time.Hour)

	// Create AI service with context-aware repositories (multi-tenant)
	aiService := services.NewAIService(
		geminiClient,
		cache,
		observationRepo,
		interventionRepo,
		vitalsRepo,
		medicationRepo,
	)

	aiHandler := handlers.NewAIHandler(aiService)

	// Auth handler with central DB for OAuth
	authHandler := handlers.NewAuthHandler(centralDB)
	logger.Info("Auth handler initialized")

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(centralDB, tenantPool)
	logger.Info("Auth middleware initialized")

	mux := http.NewServeMux()

	// Auth routes (public) - using ServeHTTP for all routes
	mux.HandleFunc("/login", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/login", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/google", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/google/callback", authHandler.ServeHTTP)
	mux.HandleFunc("/logout", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/signup", authHandler.ServeHTTP)

	// Agenda routes
	agendaService := services.NewAgendaService(appointmentRepo)
	agendaHandler := handlers.NewAgendaHandler(agendaService, patientServiceAdapter)
	dashboardHandler = handlers.NewDashboardHandler(patientServiceAdapter, sessionServiceAdapter, agendaService)

	mux.HandleFunc("/agenda", agendaHandler.View)
	mux.HandleFunc("/agenda/day", agendaHandler.DayView)
	mux.HandleFunc("/agenda/week", agendaHandler.WeekView)
	mux.HandleFunc("/agenda/month", agendaHandler.MonthView)
	mux.HandleFunc("/agenda/new", agendaHandler.NewForm)
	mux.HandleFunc("/agenda/slots", agendaHandler.GetSlots)

	// Agenda appointments CRUD
	mux.HandleFunc("/agenda/appointments/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.Method == "POST" && !strings.Contains(path, "/") {
			agendaHandler.Create(w, r)
		} else if r.Method == "GET" && !strings.HasSuffix(path, "/") {
			agendaHandler.Show(w, r)
		} else if r.Method == "PUT" && !strings.HasSuffix(path, "/") {
			agendaHandler.Update(w, r)
		} else if (r.Method == "DELETE" || (r.Method == "POST" && strings.Contains(path, "/cancel"))) && !strings.HasSuffix(path, "/") {
			agendaHandler.Cancel(w, r)
		} else if strings.Contains(path, "/reschedule") && r.Method == "POST" {
			agendaHandler.Reschedule(w, r)
		} else if strings.Contains(path, "/complete") && r.Method == "POST" {
			agendaHandler.Complete(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/agenda/appointments", agendaHandler.Create)

	// Dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/dashboard", dashboardHandler.Show)

	// Patient routes - using the actual method names from the new handlers
	mux.HandleFunc("/patients", patientHandler.ListPatients)
	// TODO: Migrate to templ
	mux.HandleFunc("/patients/new", patientHandler.NewPatient)
	mux.HandleFunc("/patients/search", patientHandler.Search)
	mux.HandleFunc("/patients/create", patientHandler.CreatePatient)

	// Session routes - using the actual method names from the new handlers
	mux.HandleFunc("/session/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") && r.Method == "GET" {
			sessionHandler.EditSession(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/update") && r.Method == "POST" {
			sessionHandler.UpdateSession(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/observations") && r.Method == "POST" {
			sessionHandler.CreateObservation(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/interventions") && r.Method == "POST" {
			sessionHandler.CreateIntervention(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/summary") && r.Method == "PATCH" {
			sessionHandler.PatchSummary(w, r)
		} else if r.Method == "GET" {
			sessionHandler.Show(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/session", sessionHandler.CreateSession)

	// Observation routes
	// Initialize classification handler
	classificationHandler := handlers.NewClassificationHandler(observationServiceAdapter)

	// Classification routes (must be registered before general observation routes)
	mux.HandleFunc("/observations/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Log for debugging
		log.Printf("[Router] /observations/ path=%s method=%s", path, r.Method)

		// Classification routes - check for /classify pattern
		// Pattern: /observations/{id}/classify or /observations/{id}/classify/{tag_id}
		if strings.Contains(path, "/classify") {
			log.Printf("[Router] Matched classify route")
			// Check if it's the edit form request
			if r.Method == "GET" && strings.HasSuffix(path, "/edit") {
				classificationHandler.GetClassificationEdit(w, r)
				return
			}
			// POST to /observations/{id}/classify - add tag
			if r.Method == "POST" && !strings.Contains(path, "/classify/") {
				classificationHandler.ClassifyObservation(w, r)
				return
			}
			// DELETE to /observations/{id}/classify/{tag_id} - remove tag
			if r.Method == "DELETE" || (r.Method == "POST" && strings.Contains(path, "/classify/")) {
				classificationHandler.RemoveClassification(w, r)
				return
			}
			// Fallback for GET without /edit
			if r.Method == "GET" {
				classificationHandler.GetClassificationEdit(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}

		// Standard observation routes
		if strings.HasSuffix(path, "/edit") && r.Method == "GET" {
			observationHandler.GetObservationEditForm(w, r)
		} else if r.Method == "GET" {
			observationHandler.GetObservation(w, r)
		} else if r.Method == "PUT" {
			observationHandler.UpdateObservation(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Tag routes for observations
	mux.HandleFunc("/tags", classificationHandler.GetTagsByType)

	// Tag routes for interventions
	mux.HandleFunc("/tags/interventions", interventionClassificationHandler.GetInterventionTagsByType)

	// Intervention routes
	mux.HandleFunc("/interventions/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Intervention classification routes
		if strings.Contains(path, "/classify") {
			// Check if it's the edit form request
			if r.Method == "GET" && strings.HasSuffix(path, "/edit") {
				interventionClassificationHandler.GetInterventionClassificationEdit(w, r)
				return
			}
			// POST to /interventions/{id}/classify - add tag
			if r.Method == "POST" && !strings.Contains(path, "/classify/") {
				interventionClassificationHandler.ClassifyIntervention(w, r)
				return
			}
			// DELETE to /interventions/{id}/classify/{tag_id} - remove tag
			if r.Method == "DELETE" || (r.Method == "POST" && strings.Contains(path, "/classify/")) {
				interventionClassificationHandler.RemoveInterventionClassification(w, r)
				return
			}
			// Fallback for GET without /edit
			if r.Method == "GET" {
				interventionClassificationHandler.GetInterventionClassificationEdit(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}

		// Standard intervention routes
		if strings.HasSuffix(r.URL.Path, "/edit") && r.Method == "GET" {
			interventionHandler.GetInterventionEditForm(w, r)
		} else if r.Method == "GET" {
			interventionHandler.GetIntervention(w, r)
		} else if r.Method == "PUT" {
			interventionHandler.UpdateIntervention(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Combined route for patient details and new sessions
	mux.HandleFunc("/patients/", func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "Route /patients/ called",
			logger.String("path", r.URL.Path),
			logger.String("method", r.Method),
		)
		if strings.HasSuffix(r.URL.Path, "/sessions/new") {
			logger.InfoContext(r.Context(), "Routing to NewSession")
			sessionHandler.NewSession(w, r)
		} else if strings.Contains(r.URL.Path, "/history/load-more") && r.Method == "GET" {
			timelineHandler.LoadMoreEvents(w, r)
		} else if strings.Contains(r.URL.Path, "/history/search") && r.Method == "GET" {
			timelineHandler.SearchPatientHistory(w, r)
		} else if strings.Contains(r.URL.Path, "/history") && r.Method == "GET" {
			timelineHandler.ShowPatientHistory(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/sessions") && r.Method == "GET" {
			patientHandler.ListSessions(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/context") && r.Method == "GET" {
			biopsychosocialHandler.GetContextPanel(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/medications") && r.Method == "POST" {
			biopsychosocialHandler.AddMedication(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/vitals") && r.Method == "POST" {
			biopsychosocialHandler.RecordVitals(w, r)
		} else if strings.Contains(r.URL.Path, "/medications/") && r.Method == "PUT" {
			biopsychosocialHandler.UpdateMedicationStatus(w, r)
		} else if strings.Contains(r.URL.Path, "/anamnesis/") && r.Method == "PATCH" {
			patientHandler.UpdateAnamnesisSection(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/anamnesis") && r.Method == "GET" {
			patientHandler.ShowAnamnesis(w, r)
		} else if strings.Contains(r.URL.Path, "/analysis/synthesis") && r.Method == "POST" {
			aiHandler.GeneratePatientSynthesis(w, r)
		} else if strings.Contains(r.URL.Path, "/analysis/themes") && r.Method == "GET" {
			analysisHandler.ShowThemes(w, r)
		} else if strings.Contains(r.URL.Path, "/plan/report") && r.Method == "GET" {
			sessionHandler.TherapeuticPlanReport(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/goals") && r.Method == "POST" {
			sessionHandler.CreateGoal(w, r)
		} else if strings.Contains(r.URL.Path, "/goals/") && strings.HasSuffix(r.URL.Path, "/close") && r.Method == "POST" {
			sessionHandler.CloseGoalWithNote(w, r)
		} else {
			logger.InfoContext(r.Context(), "Routing to patientHandler.Show")
			patientHandler.Show(w, r)
		}
	})

	// Test endpoint for network connectivity
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := `{"status": "ok", "message": "Server is running", "timestamp": "` + time.Now().Format(time.RFC3339) + `", "client_ip": "` + r.RemoteAddr + `"}`
		w.Write([]byte(jsonResponse))
	})

	// Favicon - return 204 No Content to avoid log pollution
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Screenshot test endpoint (bypasses auth for visual testing)
	mux.HandleFunc("/screenshot/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		vm := dashboard.DashboardVM{
			GreetingName: "Terapeuta",
			DateLabel:    "DOMINGO, 20 DE MARÇO",
			Stats: dashboard.Stats{TotalPatients: 12, TotalSessions: 47, SessionsThisWeek: 5, SessionsToday: 1},
			KpiItems: []dashboard.KpiItem{
				{Label: "Sessões registradas", Value: "47", Delta: "+5 esta semana", Tone: "neutral", Dark: true},
				{Label: "Pacientes ativos", Value: "12", Delta: "em acompanhamento", Tone: "up"},
				{Label: "Hoje", Value: "1", Delta: "0 agendamentos", Tone: "neutral"},
				{Label: "Esta semana", Value: "5", Delta: "sessões registradas", Tone: "neutral"},
			},
			Patients: []dashboard.PatientItem{
				{ID: "1", Name: "Maria Silva", CreatedAt: "10/01/2024"},
				{ID: "2", Name: "João Santos", CreatedAt: "15/02/2024"},
			},
			Sessions: []dashboard.SessionItem{
				{ID: "1", PatientName: "Maria Silva", Date: "20/03/2024", Summary: "Exploramos questões sobre autoeficácia.", Theme: "Cognição"},
				{ID: "2", PatientName: "João Santos", Date: "19/03/2024", Summary: "Gestão de estresse no trabalho.", Theme: "Ansiedade"},
			},
		}
		dashboard.Dashboard(vm).Render(r.Context(), w)
	})

	// File server with cache busting for development
	staticDir := http.Dir("web/static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable all caching during development - always check file modification time
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Vary", "Accept-Encoding")

		// For CSS/JS files, add version header based on mtime
		if strings.HasSuffix(r.URL.Path, ".css") || strings.HasSuffix(r.URL.Path, ".js") {
			if info, err := staticDir.Open(r.URL.Path); err == nil {
				if stat, ok := info.(os.FileInfo); ok {
					w.Header().Set("X-CSS-Version", fmt.Sprintf("%d", stat.ModTime().Unix()))
				}
				info.Close()
			}
		}

		http.FileServer(staticDir).ServeHTTP(w, r)
	})))

	// Create CORS middleware
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})

	// Apply auth middleware to all routes
	protectedHandler := authMiddleware.Middleware(corsHandler)

	// Apply RequestID middleware first (must be before auth)
	handlerWithRequestID := middleware.RequestIDMiddleware(protectedHandler)

	// Apply HTMX Cache middleware (prevents caching of HTMX fragments)
	handlerWithHTMXCache := middleware.HTMXCacheMiddleware(handlerWithRequestID)

	// Apply telemetry middleware after RequestID (so it has access to request_id)
	telemetryMiddleware := middleware.NewTelemetryMiddleware("/static/")
	handlerWithTelemetry := telemetryMiddleware.Middleware(handlerWithHTMXCache)

	// Create a recovery middleware
	recoveryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorContext(r.Context(), "PANIC recovered",
					logger.String("error", fmt.Sprintf("%v", err)),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		handlerWithTelemetry.ServeHTTP(w, r)
	})

	port := ":8080"
	logger.Info("Starting server",
		logger.String("address", "http://localhost"+port),
		logger.String("accessibility", "network"),
	)
	if err := http.ListenAndServe(port, recoveryHandler); err != nil {
		logger.Error("Server failed", logger.String("error", err.Error()))
		os.Exit(1)
	}
}
