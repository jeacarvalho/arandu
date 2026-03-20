package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/infrastructure/ai"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/platform/middleware"
	"arandu/internal/web"
	"arandu/internal/web/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading: %v", err)
	}

	// Initialize Control Plane (Central DB)
	log.Printf("Initializing Control Plane (Central DB)...")
	centralDB, err := sqlite.NewCentralDB("storage")
	if err != nil {
		log.Fatalf("Failed to initialize central database: %v", err)
	}
	defer centralDB.Close()

	// Apply central migrations
	log.Printf("Applying central database migrations...")
	if err := centralDB.Migrate(nil); err != nil {
		log.Printf("Warning: Failed to apply central migrations: %v", err)
	}

	// Ensure tenants directory exists
	tenantsDir := "storage/tenants"
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		log.Printf("Warning: Failed to create tenants directory: %v", err)
	} else {
		log.Printf("Tenants directory ready: %s", tenantsDir)
	}

	// Use production database (Data Plane)
	dbPath := "arandu.db"
	log.Printf("Using database: %s", dbPath)
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	log.Printf("Initializing database schema...")

	patientRepo := sqlite.NewPatientRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)
	observationRepo := sqlite.NewObservationRepository(db)
	interventionRepo := sqlite.NewInterventionRepository(db)
	insightRepo := sqlite.NewInsightRepository(db)
	medicationRepo := sqlite.NewMedicationRepository(db)
	vitalsRepo := sqlite.NewVitalsRepository(db)

	// Apply database migrations
	if err := db.Migrate(); err != nil {
		log.Printf("Warning: Failed to apply database migrations: %v", err)
	}

	patientService := services.NewPatientService(patientRepo)
	sessionService := services.NewSessionService(sessionRepo)
	observationService := services.NewObservationService(observationRepo)
	interventionService := services.NewInterventionService(interventionRepo)
	insightService := services.NewInsightService(insightRepo)
	timelineRepo := sqlite.NewTimelineRepository(db)
	timelineService := services.NewTimelineService(timelineRepo)
	biopsychosocialService := services.NewBiopsychosocialService(medicationRepo, vitalsRepo)

	// Create service adapters for the new handler interfaces
	sessionServiceAdapter := web.NewSessionServiceAdapter(sessionService)
	insightServiceAdapter := web.NewInsightServiceAdapter(insightService)
	patientServiceAdapter := web.NewPatientServiceAdapter(patientService)
	observationServiceAdapter := web.NewObservationServiceAdapter(observationService)
	interventionServiceAdapter := web.NewInterventionServiceAdapter(interventionService)
	timelineServiceAdapter := web.NewTimelineServiceAdapter(timelineService)

	// Create biopsychosocial service adapter
	biopsychosocialServiceAdapterImpl := handlers.BiopsychosocialServiceFuncs{
		GetMedicationsFunc: func(patientID string) ([]interface{}, error) {
			meds, err := biopsychosocialService.GetMedications(patientID)
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
		GetLatestVitalsFunc: func(patientID string) (interface{}, error) {
			return biopsychosocialService.GetLatestVitals(patientID)
		},
		GetAverageVitalsFunc: func(patientID string, days int) (interface{}, error) {
			return biopsychosocialService.GetAverageVitals(patientID, days)
		},
	}

	// Create new handlers with dependency injection
	patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, insightServiceAdapter, biopsychosocialServiceAdapterImpl)
	sessionHandler := handlers.NewSessionHandler(sessionServiceAdapter, patientServiceAdapter, observationServiceAdapter, interventionServiceAdapter)
	observationHandler := handlers.NewObservationHandler(observationServiceAdapter)
	interventionHandler := handlers.NewInterventionHandler(interventionServiceAdapter)
	dashboardHandler := handlers.NewDashboardHandler(patientServiceAdapter, sessionServiceAdapter)
	timelineHandler := handlers.NewTimelineHandler(timelineServiceAdapter)
	biopsychosocialHandler := handlers.NewBiopsychosocialHandler(biopsychosocialService)

	// Initialize AI service
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Printf("Warning: GEMINI_API_KEY not set. AI features will be disabled.")
		geminiAPIKey = "dummy-key-for-initialization" // Use dummy key to allow initialization
	}

	geminiClient, err := ai.NewGeminiClient(geminiAPIKey)
	if err != nil {
		log.Printf("Warning: Failed to initialize Gemini client: %v", err)
	}

	// Create cache for AI responses (24 hour TTL)
	cache := ai.NewCache(24 * time.Hour)

	// Create AI service with repository adapters
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
	log.Printf("Auth handler initialized")

	// Initialize Tenant Pool for multi-tenant connections
	tenantPool := sqlite.NewTenantPool("storage", nil)
	log.Printf("Tenant pool initialized")

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(centralDB, tenantPool)
	log.Printf("Auth middleware initialized")

	mux := http.NewServeMux()

	// Auth routes (public) - using ServeHTTP for all routes
	mux.HandleFunc("/login", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/login", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/google", authHandler.ServeHTTP)
	mux.HandleFunc("/auth/google/callback", authHandler.ServeHTTP)
	mux.HandleFunc("/logout", authHandler.ServeHTTP)

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
	mux.HandleFunc("/patient/create", patientHandler.CreatePatient)

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
		} else if r.Method == "GET" {
			sessionHandler.Show(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/session", sessionHandler.CreateSession)

	// Observation routes
	mux.HandleFunc("/observations/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") && r.Method == "GET" {
			observationHandler.GetObservationEditForm(w, r)
		} else if r.Method == "GET" {
			observationHandler.GetObservation(w, r)
		} else if r.Method == "PUT" {
			observationHandler.UpdateObservation(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Intervention routes
	mux.HandleFunc("/interventions/", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("/patient/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Route /patient/ called: %s, Method: %s", r.URL.Path, r.Method)
		if strings.HasSuffix(r.URL.Path, "/sessions/new") {
			log.Printf("  -> Routing to NewSession")
			sessionHandler.NewSession(w, r)
		} else {
			log.Printf("  -> Routing to patientHandler.Show")
			patientHandler.Show(w, r)
		}
	})

	// Combined patient routes (plural)
	mux.HandleFunc("/patients/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/history/search") && r.Method == "GET" {
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
		} else if strings.Contains(r.URL.Path, "/analysis/synthesis") && r.Method == "POST" {
			aiHandler.GeneratePatientSynthesis(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Test endpoint for network connectivity
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonResponse := `{"status": "ok", "message": "Server is running", "timestamp": "` + time.Now().Format(time.RFC3339) + `", "client_ip": "` + r.RemoteAddr + `"}`
		w.Write([]byte(jsonResponse))
	})

	// File server with cache control for CSS files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable cache for CSS files during development
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		fs.ServeHTTP(w, r)
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

	// Create a recovery middleware
	recoveryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		protectedHandler.ServeHTTP(w, r)
	})

	port := ":8080"
	log.Printf("Starting server on http://localhost%s (accessible from network)", port)
	if err := http.ListenAndServe(port, recoveryHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
