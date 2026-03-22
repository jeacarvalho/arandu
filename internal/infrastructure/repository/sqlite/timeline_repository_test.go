package sqlite

import (
	"context"
	"fmt"
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

		// Test filtering by session type
		sessionFilter := timeline.EventTypeSession
		sessionEvents, err := timelineRepo.GetTimelineByPatientID(ctx, p.ID, &sessionFilter, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get filtered timeline for sessions: %v", err)
		}

		if len(sessionEvents) != 1 {
			t.Errorf("Expected 1 session event after filter, got %d", len(sessionEvents))
		}

		if sessionEvents[0].Type != timeline.EventTypeSession {
			t.Errorf("Expected session event type, got %s", sessionEvents[0].Type)
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

	t.Run("Chronological parity - events ordered correctly", func(t *testing.T) {
		// Create a patient
		p, err := patient.NewPatient("Chronology Test Patient", "Test Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}
		if err := patientRepo.Save(ctx, p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Create sessions with specific dates (yesterday and day before)
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		dayBeforeYesterday := now.Add(-48 * time.Hour)

		sessionYesterday := session.NewSession(p.ID, yesterday, "Session from yesterday")
		if err := sessionRepo.Create(ctx, sessionYesterday); err != nil {
			t.Fatalf("Failed to save session from yesterday: %v", err)
		}

		sessionDayBefore := session.NewSession(p.ID, dayBeforeYesterday, "Session from day before yesterday")
		if err := sessionRepo.Create(ctx, sessionDayBefore); err != nil {
			t.Fatalf("Failed to save session from day before: %v", err)
		}

		// Create observation for yesterday's session
		obsYesterday := &observation.Observation{
			ID:        "obs-yesterday",
			SessionID: sessionYesterday.ID,
			Content:   "Observation from yesterday",
			CreatedAt: yesterday.Add(time.Hour),
			UpdatedAt: yesterday.Add(time.Hour),
		}
		if err := observationRepo.Save(ctx, obsYesterday); err != nil {
			t.Fatalf("Failed to save observation: %v", err)
		}

		// Create observation for day before's session (but with today's timestamp)
		obsToday := &observation.Observation{
			ID:        "obs-today",
			SessionID: sessionDayBefore.ID,
			Content:   "Observation created today for older session",
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := observationRepo.Save(ctx, obsToday); err != nil {
			t.Fatalf("Failed to save today's observation: %v", err)
		}

		// Get timeline
		events, err := timelineRepo.GetTimelineByPatientID(ctx, p.ID, nil, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get timeline: %v", err)
		}

		// Verify chronological order: today's observation should be first, then yesterday's session, then yesterday's observation, then day before's session
		if len(events) < 4 {
			t.Fatalf("Expected at least 4 events, got %d", len(events))
		}

		// First event should be the observation created today
		if events[0].ID != obsToday.ID {
			t.Errorf("First event should be observation from today (ID: %s), got %s", obsToday.ID, events[0].ID)
		}

		// Verify events are in descending order
		for i := 0; i < len(events)-1; i++ {
			if events[i].Date.Before(events[i+1].Date) {
				t.Errorf("Events not in descending order: event[%d] date %v is before event[%d] date %v",
					i, events[i].Date, i+1, events[i+1].Date)
			}
		}
	})

	t.Run("Isolation - timeline only shows events from selected patient", func(t *testing.T) {
		// Create patient 1
		p1, err := patient.NewPatient("Isolation Patient 1", "Notes 1")
		if err != nil {
			t.Fatalf("Failed to create patient 1: %v", err)
		}
		if err := patientRepo.Save(ctx, p1); err != nil {
			t.Fatalf("Failed to save patient 1: %v", err)
		}

		// Create patient 2
		p2, err := patient.NewPatient("Isolation Patient 2", "Notes 2")
		if err != nil {
			t.Fatalf("Failed to create patient 2: %v", err)
		}
		if err := patientRepo.Save(ctx, p2); err != nil {
			t.Fatalf("Failed to save patient 2: %v", err)
		}

		// Create session for patient 1
		session1 := session.NewSession(p1.ID, time.Now(), "Session for patient 1")
		if err := sessionRepo.Create(ctx, session1); err != nil {
			t.Fatalf("Failed to save session for patient 1: %v", err)
		}

		// Create observation for patient 1
		obs1 := &observation.Observation{
			ID:        "obs-patient-1",
			SessionID: session1.ID,
			Content:   "Observation for patient 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := observationRepo.Save(ctx, obs1); err != nil {
			t.Fatalf("Failed to save observation for patient 1: %v", err)
		}

		// Create session for patient 2
		session2 := session.NewSession(p2.ID, time.Now(), "Session for patient 2")
		if err := sessionRepo.Create(ctx, session2); err != nil {
			t.Fatalf("Failed to save session for patient 2: %v", err)
		}

		// Create observation for patient 2
		obs2 := &observation.Observation{
			ID:        "obs-patient-2",
			SessionID: session2.ID,
			Content:   "Observation for patient 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := observationRepo.Save(ctx, obs2); err != nil {
			t.Fatalf("Failed to save observation for patient 2: %v", err)
		}

		// Get timeline for patient 1
		events1, err := timelineRepo.GetTimelineByPatientID(ctx, p1.ID, nil, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get timeline for patient 1: %v", err)
		}

		// Verify patient 1 timeline only has their events
		for _, event := range events1 {
			// All events should be from patient 1's session
			sessionID := event.Metadata["session_id"]
			if sessionID == session2.ID {
				t.Errorf("Timeline for patient 1 contains event from patient 2's session: %v", event)
			}
		}

		// Should have exactly 2 events (session + observation)
		if len(events1) != 2 {
			t.Errorf("Expected 2 events for patient 1, got %d", len(events1))
		}

		// Get timeline for patient 2
		events2, err := timelineRepo.GetTimelineByPatientID(ctx, p2.ID, nil, 100, 0)
		if err != nil {
			t.Fatalf("Failed to get timeline for patient 2: %v", err)
		}

		// Verify patient 2 timeline only has their events
		for _, event := range events2 {
			sessionID := event.Metadata["session_id"]
			if sessionID == session1.ID {
				t.Errorf("Timeline for patient 2 contains event from patient 1's session: %v", event)
			}
		}

		// Should have exactly 2 events (session + observation)
		if len(events2) != 2 {
			t.Errorf("Expected 2 events for patient 2, got %d", len(events2))
		}
	})

	t.Run("Performance - timeline with 50 events loads efficiently", func(t *testing.T) {
		// Create patient
		p, err := patient.NewPatient("Performance Test Patient", "Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}
		if err := patientRepo.Save(ctx, p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Create session
		sess := session.NewSession(p.ID, time.Now().Add(-24*time.Hour), "Performance test session")
		if err := sessionRepo.Create(ctx, sess); err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		// Create 49 observations (plus 1 session = 50 events)
		for i := 0; i < 49; i++ {
			obs := &observation.Observation{
				ID:        fmt.Sprintf("perf-obs-%d", i),
				SessionID: sess.ID,
				Content:   fmt.Sprintf("Performance test observation number %d with some meaningful clinical content", i),
				CreatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
				UpdatedAt: time.Now().Add(-time.Duration(i) * time.Hour),
			}
			if err := observationRepo.Save(ctx, obs); err != nil {
				t.Fatalf("Failed to save observation %d: %v", i, err)
			}
		}

		// Measure query time
		start := time.Now()
		events, err := timelineRepo.GetTimelineByPatientID(ctx, p.ID, nil, 50, 0)
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Failed to get timeline: %v", err)
		}

		// Verify we got all 50 events
		if len(events) != 50 {
			t.Errorf("Expected 50 events, got %d", len(events))
		}

		// Performance requirement: should load in less than 500ms
		if duration > 500*time.Millisecond {
			t.Errorf("Timeline query took too long: %v (expected < 500ms)", duration)
		}

		t.Logf("Timeline with 50 events loaded in %v", duration)
	})
}

// BenchmarkTimelinePerformance runs a benchmark for timeline queries
func BenchmarkTimelineRepository_GetTimelineByPatientID(b *testing.B) {
	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "benchmark-*.db")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Initialize database
	db, err := NewDB(tmpfile.Name())
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		b.Fatalf("Failed to run migrations: %v", err)
	}

	// Create repositories
	timelineRepo := NewTimelineRepository(db)
	patientRepo := NewPatientRepository(db)
	sessionRepo := NewSessionRepository(db)
	observationRepo := NewObservationRepository(db)

	ctx := context.Background()

	// Create patient
	p, _ := patient.NewPatient("Benchmark Patient", "Notes")
	patientRepo.Save(ctx, p)

	// Create session
	sess := session.NewSession(p.ID, time.Now(), "Benchmark session")
	sessionRepo.Create(ctx, sess)

	// Create 100 observations
	for i := 0; i < 100; i++ {
		obs := &observation.Observation{
			ID:        fmt.Sprintf("bench-obs-%d", i),
			SessionID: sess.ID,
			Content:   fmt.Sprintf("Benchmark observation %d", i),
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute),
			UpdatedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
		observationRepo.Save(ctx, obs)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := timelineRepo.GetTimelineByPatientID(ctx, p.ID, nil, 50, 0)
		if err != nil {
			b.Fatalf("Failed to get timeline: %v", err)
		}
	}
}
