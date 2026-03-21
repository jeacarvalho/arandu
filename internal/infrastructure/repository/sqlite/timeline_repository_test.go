package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"arandu/internal/domain/intervention"
	"arandu/internal/domain/observation"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
	"arandu/internal/domain/timeline"
)

func TestTimelineRepositoryIntegration(t *testing.T) {
	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Initialize database
	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations to create all tables
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	timelineRepo := NewTimelineRepository(db)
	patientRepo := NewPatientRepository(db)
	sessionRepo := NewSessionRepository(db)
	observationRepo := NewObservationRepository(db)
	interventionRepo := NewInterventionRepository(db)

	ctx := context.Background()

	t.Run("Get timeline for patient with all event types", func(t *testing.T) {
		// Create a patient
		p, err := patient.NewPatient("Test Patient", "Test Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}
		if err := patientRepo.Save(ctx, p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Create a session
		sess := session.NewSession(p.ID, time.Now(), "Test session summary")
		if err := sessionRepo.Create(ctx, sess); err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		// Create an observation
		obs := &observation.Observation{
			ID:        "test-observation-1",
			SessionID: sess.ID,
			Content:   "Test observation content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := observationRepo.Save(ctx, obs); err != nil {
			t.Fatalf("Failed to save observation: %v", err)
		}

		// Create an intervention
		intv := &intervention.Intervention{
			ID:        "test-intervention-1",
			SessionID: sess.ID,
			Content:   "Test intervention content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := interventionRepo.Save(ctx, intv); err != nil {
			t.Fatalf("Failed to save intervention: %v", err)
		}

		// Test getting timeline without filter
		events, err := timelineRepo.GetTimelineByPatientID(ctx, p.ID, nil, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get timeline: %v", err)
		}

		if len(events) != 3 {
			t.Errorf("Expected 3 events, got %d", len(events))
		}

		// Verify event types
		eventTypes := make(map[timeline.EventType]int)
		for _, event := range events {
			eventTypes[event.Type]++
		}

		if eventTypes[timeline.EventTypeSession] != 1 {
			t.Errorf("Expected 1 session event, got %d", eventTypes[timeline.EventTypeSession])
		}
		if eventTypes[timeline.EventTypeObservation] != 1 {
			t.Errorf("Expected 1 observation event, got %d", eventTypes[timeline.EventTypeObservation])
		}
		if eventTypes[timeline.EventTypeIntervention] != 1 {
			t.Errorf("Expected 1 intervention event, got %d", eventTypes[timeline.EventTypeIntervention])
		}
	})

	t.Run("Get timeline with filter", func(t *testing.T) {
		// Create another patient
		p, err := patient.NewPatient("Test Patient 2", "Test Notes 2")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}
		if err := patientRepo.Save(ctx, p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Create a session
		sess := session.NewSession(p.ID, time.Now(), "Test session summary 2")
		if err := sessionRepo.Create(ctx, sess); err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		// Create an observation
		obs := &observation.Observation{
			ID:        "test-observation-2",
			SessionID: sess.ID,
			Content:   "Test observation content 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := observationRepo.Save(ctx, obs); err != nil {
			t.Fatalf("Failed to save observation: %v", err)
		}

		// Test filtering by observation type
		obsFilter := timeline.EventTypeObservation
		obsEvents, err := timelineRepo.GetTimelineByPatientID(ctx, p.ID, &obsFilter, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get filtered timeline: %v", err)
		}

		if len(obsEvents) != 1 {
			t.Errorf("Expected 1 observation event after filter, got %d", len(obsEvents))
		}

		if obsEvents[0].Type != timeline.EventTypeObservation {
			t.Errorf("Expected observation event type, got %s", obsEvents[0].Type)
		}
	})

	t.Run("Get timeline for non-existent patient", func(t *testing.T) {
		emptyEvents, err := timelineRepo.GetTimelineByPatientID(ctx, "non-existent-patient", nil, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get timeline for non-existent patient: %v", err)
		}

		if len(emptyEvents) != 0 {
			t.Errorf("Expected 0 events for non-existent patient, got %d", len(emptyEvents))
		}
	})

	t.Run("Search in patient history", func(t *testing.T) {
		// Create a patient
		p, err := patient.NewPatient("Search Test Patient", "Test Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}
		if err := patientRepo.Save(ctx, p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Create a session
		sess := session.NewSession(p.ID, time.Now(), "Test session with transferência theme")
		if err := sessionRepo.Create(ctx, sess); err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		// Create an observation with searchable content
		obs := &observation.Observation{
			ID:        "test-observation-search",
			SessionID: sess.ID,
			Content:   "Paciente demonstrou sinais de transferência durante a sessão. Relatou sentimentos de luto pela perda recente.",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := observationRepo.Save(ctx, obs); err != nil {
			t.Fatalf("Failed to save observation: %v", err)
		}

		// Create an intervention with searchable content
		intv := &intervention.Intervention{
			ID:        "test-intervention-search",
			SessionID: sess.ID,
			Content:   "Intervenção focada em processar o luto e trabalhar a transferência de forma terapêutica.",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := interventionRepo.Save(ctx, intv); err != nil {
			t.Fatalf("Failed to save intervention: %v", err)
		}

		// Test search for "transferência"
		results, err := timelineRepo.SearchInHistory(ctx, p.ID, "transferência")
		if err != nil {
			t.Fatalf("Failed to search in history: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 search results for 'transferência', got %d", len(results))
		}

		// Test search for "luto"
		results, err = timelineRepo.SearchInHistory(ctx, p.ID, "luto")
		if err != nil {
			t.Fatalf("Failed to search in history: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 search results for 'luto', got %d", len(results))
		}

		// Test search for non-existent term
		results, err = timelineRepo.SearchInHistory(ctx, p.ID, "inexistente")
		if err != nil {
			t.Fatalf("Failed to search in history: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 search results for 'inexistente', got %d", len(results))
		}

		// Test empty query
		results, err = timelineRepo.SearchInHistory(ctx, p.ID, "")
		if err != nil {
			t.Fatalf("Failed to search in history with empty query: %v", err)
		}

		if results != nil {
			t.Errorf("Expected nil results for empty query, got %v", results)
		}
	})
}
