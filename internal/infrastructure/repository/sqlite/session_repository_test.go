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

	// Create tables
	patientQuery := `
	CREATE TABLE IF NOT EXISTS patients (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		notes TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
	`
	if _, err := db.Exec(patientQuery); err != nil {
		t.Fatalf("Failed to create patients table: %v", err)
	}
	sessionQuery := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		patient_id TEXT NOT NULL,
		date DATETIME NOT NULL,
		summary TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
	)
	`
	if _, err := db.Exec(sessionQuery); err != nil {
		t.Fatalf("Failed to create sessions table: %v", err)
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

	// Create a patient to associate sessions with
	p, err := patient.NewPatient("Test Patient", "Test Notes")
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}
	if err := patientRepo.Save(p); err != nil {
		t.Fatalf("Failed to save patient: %v", err)
	}

	ctx := context.Background()

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
}
