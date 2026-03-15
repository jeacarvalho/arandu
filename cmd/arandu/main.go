package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"arandu/internal/application/services"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/web"
	"arandu/internal/web/handlers"
)

func main() {
	dbPath := filepath.Join(os.Getenv("HOME"), ".arandu", "database.db")
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	// First, run migrations for tables that have migration support
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	migrationsDir := filepath.Join(exeDir, "internal", "infrastructure", "repository", "sqlite", "migrations")

	// Fallback to relative path if absolute path doesn't work
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = filepath.Join(".", "internal", "infrastructure", "repository", "sqlite", "migrations")
	}

	log.Printf("Using migrations directory: %s", migrationsDir)
	if err := db.Migrate(migrationsDir); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Then, initialize other schemas that don't have migrations yet
	// TODO: Convert all tables to use migrations
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	patientRepo := sqlite.NewPatientRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)
	// observationRepo := sqlite.NewObservationRepository(db) // TODO: Migrate later
	// interventionRepo := sqlite.NewInterventionRepository(db) // TODO: Migrate later
	insightRepo := sqlite.NewInsightRepository(db)

	patientService := services.NewPatientService(patientRepo)
	sessionService := services.NewSessionService(sessionRepo)
	// observationService := services.NewObservationService(observationRepo) // TODO: Migrate later
	// interventionService := services.NewInterventionService(interventionRepo) // TODO: Migrate later
	insightService := services.NewInsightService(insightRepo)

	// Create template renderer
	templateRenderer := web.NewTemplateRendererAdapter("web/templates")

	// Create service adapters for the new handler interfaces
	sessionServiceAdapter := web.NewSessionServiceAdapter(sessionService)
	insightServiceAdapter := web.NewInsightServiceAdapter(insightService)
	patientServiceAdapter := web.NewPatientServiceAdapter(patientService)

	// Create new handlers with dependency injection
	patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, insightServiceAdapter, templateRenderer)
	sessionHandler := handlers.NewSessionHandler(sessionServiceAdapter, patientServiceAdapter, templateRenderer)
	dashboardHandler := handlers.NewDashboardHandler(patientServiceAdapter, sessionServiceAdapter, templateRenderer)

	mux := http.NewServeMux()

	// Dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/dashboard", dashboardHandler.Show)

	// Patient routes - using the actual method names from the new handlers
	mux.HandleFunc("/patients", patientHandler.ListPatients)
	mux.HandleFunc("/patients/new", patientHandler.NewPatient)
	mux.HandleFunc("/patient/create", patientHandler.CreatePatient)

	// Session routes - using the actual method names from the new handlers
	mux.HandleFunc("/session/", sessionHandler.Show)
	mux.HandleFunc("/sessions", sessionHandler.CreateSession)
	mux.HandleFunc("/sessions/edit/", sessionHandler.EditSession)
	mux.HandleFunc("/sessions/update/", sessionHandler.UpdateSession)

	// Combined route for patient details and new sessions
	mux.HandleFunc("/patient/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sessions/new") {
			sessionHandler.NewSession(w, r)
		} else {
			patientHandler.Show(w, r)
		}
	})

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":8080"
	log.Printf("Starting server on http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
