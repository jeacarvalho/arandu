package main

import (
	"log"
	"net/http"
	"strings"

	"arandu/internal/application/services"
	"arandu/internal/infrastructure/repository/sqlite"
	"arandu/internal/web"
	"arandu/internal/web/handlers"
)

func main() {
	// Use production database
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

	// Apply database migrations
	if err := db.Migrate(); err != nil {
		log.Printf("Warning: Failed to apply database migrations: %v", err)
	}

	patientService := services.NewPatientService(patientRepo)
	sessionService := services.NewSessionService(sessionRepo)
	observationService := services.NewObservationService(observationRepo)
	interventionService := services.NewInterventionService(interventionRepo)
	insightService := services.NewInsightService(insightRepo)

	// Create service adapters for the new handler interfaces
	sessionServiceAdapter := web.NewSessionServiceAdapter(sessionService)
	insightServiceAdapter := web.NewInsightServiceAdapter(insightService)
	patientServiceAdapter := web.NewPatientServiceAdapter(patientService)
	observationServiceAdapter := web.NewObservationServiceAdapter(observationService)
	interventionServiceAdapter := web.NewInterventionServiceAdapter(interventionService)

	// Create new handlers with dependency injection
	patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, insightServiceAdapter)
	sessionHandler := handlers.NewSessionHandler(sessionServiceAdapter, patientServiceAdapter, observationServiceAdapter, interventionServiceAdapter)
	observationHandler := handlers.NewObservationHandler(observationServiceAdapter)
	interventionHandler := handlers.NewInterventionHandler(interventionServiceAdapter)
	dashboardHandler := handlers.NewDashboardHandler(patientServiceAdapter, sessionServiceAdapter)

	mux := http.NewServeMux()

	// Dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
	mux.HandleFunc("/dashboard", dashboardHandler.Show)

	// Patient routes - using the actual method names from the new handlers
	mux.HandleFunc("/patients", patientHandler.ListPatients)
	// TODO: Migrate to templ
	mux.HandleFunc("/patients/new", patientHandler.NewPatient)
	mux.HandleFunc("/patient/create", patientHandler.CreatePatient)

	// Session routes - using the actual method names from the new handlers
	mux.HandleFunc("/session/", sessionHandler.Show)
	mux.HandleFunc("/sessions", sessionHandler.CreateSession)
	mux.HandleFunc("/sessions/edit/", sessionHandler.EditSession)
	mux.HandleFunc("/sessions/update/", sessionHandler.UpdateSession)
	mux.HandleFunc("/sessions/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/observations") && r.Method == "POST" {
			sessionHandler.CreateObservation(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/interventions") && r.Method == "POST" {
			sessionHandler.CreateIntervention(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

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

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Create a recovery middleware
	recoveryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		mux.ServeHTTP(w, r)
	})

	port := ":8080"
	log.Printf("Starting server on http://localhost%s", port)
	if err := http.ListenAndServe(port, recoveryHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
