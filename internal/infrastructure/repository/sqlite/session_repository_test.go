package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
)

func setupSessionTestDB(t *testing.T) (*DB, func()) {
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Run migrations to create all tables
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func TestSessionRepositoryIntegration(t *testing.T) {
	db, cleanup := setupSessionTestDB(t)
	defer cleanup()

	patientRepo := NewPatientRepository(db)
	sessionRepo := NewSessionRepository(db)

	ctx := context.Background()

	// Create a patient to associate sessions with
	p, err := patient.NewPatient("Test Patient", "Test Notes")
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}
	if err := patientRepo.Save(ctx, p); err != nil {
		t.Fatalf("Failed to save patient: %v", err)
	}

	t.Run("Create and retrieve session", func(t *testing.T) {
		sess := session.NewSession(p.ID, time.Now(), "Test Summary")

		if err := sessionRepo.Create(ctx, sess); err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		retrieved, err := sessionRepo.GetByID(ctx, sess.ID)
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}
		if retrieved == nil {
			t.Fatal("Retrieved session is nil")
		}

		if retrieved.ID != sess.ID {
			t.Errorf("ID mismatch: got %q, want %q", retrieved.ID, sess.ID)
		}
		if retrieved.Summary != sess.Summary {
			t.Errorf("Summary mismatch: got %q, want %q", retrieved.Summary, sess.Summary)
		}
	})

	t.Run("Update session", func(t *testing.T) {
		sess := session.NewSession(p.ID, time.Now(), "Original Summary")
		if err := sessionRepo.Create(ctx, sess); err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		originalUpdatedAt := sess.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure timestamp changes

		newDate := time.Now().Add(-24 * time.Hour)
		newSummary := "Updated Summary"

		if err := sess.Update(newDate, newSummary); err != nil {
			t.Fatalf("Failed to update session domain object: %v", err)
		}

		if err := sessionRepo.Update(ctx, sess); err != nil {
			t.Fatalf("Failed to save session update: %v", err)
		}

		retrieved, err := sessionRepo.GetByID(ctx, sess.ID)
		if err != nil {
			t.Fatalf("Failed to get updated session: %v", err)
		}

		if retrieved.Summary != newSummary {
			t.Errorf("Summary not updated: got %q, want %q", retrieved.Summary, newSummary)
		}
		if !retrieved.UpdatedAt.After(originalUpdatedAt) {
			t.Errorf("UpdatedAt not updated: got %v, original was %v", retrieved.UpdatedAt, originalUpdatedAt)
		}
		if !retrieved.Date.Equal(newDate) {
			// Truncate to millisecond precision for comparison to avoid SQLite precision issues
			if !retrieved.Date.Truncate(time.Millisecond).Equal(newDate.Truncate(time.Millisecond)) {
				t.Errorf("Date not updated: got %v, want %v", retrieved.Date, newDate)
			}
		}
	})

	t.Run("List sessions by patient", func(t *testing.T) {
		// Clear existing sessions for this patient by updating the patient ID
		// Since we can't delete, we'll create a new patient for this test
		newPatient, err := patient.NewPatient("Test Patient 2", "Notes 2")
		if err != nil {
			t.Fatalf("Failed to create new patient: %v", err)
		}
		if err := patientRepo.Save(ctx, newPatient); err != nil {
			t.Fatalf("Failed to save new patient: %v", err)
		}

		// Use the new patient for this test
		testPatientID := newPatient.ID

		// Create multiple sessions for the new patient
		sessions := []*session.Session{}
		for i := 1; i <= 3; i++ {
			sess := session.NewSession(testPatientID, time.Now().Add(time.Duration(i)*time.Hour),
				"Session "+string(rune('A'+i-1)))
			if err := sessionRepo.Create(ctx, sess); err != nil {
				t.Fatalf("Failed to create session %d: %v", i, err)
			}
			sessions = append(sessions, sess)
		}

		// List sessions for the new patient
		retrieved, err := sessionRepo.ListByPatient(ctx, testPatientID)
		if err != nil {
			t.Fatalf("Failed to list sessions: %v", err)
		}

		if len(retrieved) != len(sessions) {
			t.Errorf("Expected %d sessions, got %d", len(sessions), len(retrieved))
		}

		// Verify sessions are returned in reverse chronological order (ORDER BY date DESC)
		for i := 0; i < len(retrieved)-1; i++ {
			if retrieved[i].Date.Before(retrieved[i+1].Date) {
				t.Errorf("Sessions not in reverse chronological order: session %d is before session %d", i, i+1)
			}
		}
	})

	// Note: SessionRepository doesn't have a Delete method yet
	// This test is commented out until Delete is implemented
	/*
		t.Run("Delete session", func(t *testing.T) {
			sess := session.NewSession(p.ID, time.Now(), "Session to delete")
			if err := sessionRepo.Create(ctx, sess); err != nil {
				t.Fatalf("Failed to create session: %v", err)
			}

			// Verify session exists
			retrieved, err := sessionRepo.GetByID(ctx, sess.ID)
			if err != nil {
				t.Fatalf("Failed to get session: %v", err)
			}
			if retrieved == nil {
				t.Fatal("Session not found before delete")
			}

			// Delete session
			if err := sessionRepo.Delete(ctx, sess.ID); err != nil {
				t.Fatalf("Failed to delete session: %v", err)
			}

			// Verify session is deleted
			deleted, err := sessionRepo.GetByID(ctx, sess.ID)
			if err != nil {
				t.Fatalf("Failed to get deleted session: %v", err)
			}
			if deleted != nil {
				t.Error("Session still exists after delete")
			}
		})
	*/

	t.Run("Get non-existent session returns nil", func(t *testing.T) {
		nonExistentID := "non-existent-id"
		sess, err := sessionRepo.GetByID(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("Error getting non-existent session: %v", err)
		}
		if sess != nil {
			t.Error("Expected nil for non-existent session")
		}
	})

	t.Run("List sessions for non-existent patient returns empty", func(t *testing.T) {
		nonExistentPatientID := "non-existent-patient"
		sessions, err := sessionRepo.ListByPatient(ctx, nonExistentPatientID)
		if err != nil {
			t.Fatalf("Error listing sessions for non-existent patient: %v", err)
		}
		if len(sessions) != 0 {
			t.Errorf("Expected empty list for non-existent patient, got %d sessions", len(sessions))
		}
	})

	t.Run("Update non-existent session doesn't error (SQLite behavior)", func(t *testing.T) {
		nonExistentSession := &session.Session{
			ID:        "non-existent-id",
			PatientID: p.ID,
			Date:      time.Now(),
			Summary:   "Non-existent session",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := sessionRepo.Update(ctx, nonExistentSession)
		if err != nil {
			t.Errorf("Update on non-existent session returned error: %v", err)
		}
		// SQLite UPDATE on non-existent row doesn't error, just affects 0 rows
	})
}
