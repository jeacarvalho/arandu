package sqlite

import (
	"context"
	"os"
	"testing"

	"arandu/internal/domain/patient"
)

func TestAnamnesisRepositoryIntegration(t *testing.T) {
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

	// Create repository
	repo := NewPatientRepository(db)

	ctx := context.Background()

	t.Run("Create and retrieve anamnesis", func(t *testing.T) {
		patientID := "test-patient-123"

		// Create anamnesis
		anamnesis, err := patient.NewAnamnesis(patientID)
		if err != nil {
			t.Fatalf("Failed to create anamnesis: %v", err)
		}

		// Update sections
		if err := anamnesis.UpdateSection("chief_complaint", "Ansiedade generalizada há 2 anos"); err != nil {
			t.Fatalf("Failed to update chief complaint: %v", err)
		}
		if err := anamnesis.UpdateSection("personal_history", "Solteiro, trabalha como desenvolvedor"); err != nil {
			t.Fatalf("Failed to update personal history: %v", err)
		}

		// Save anamnesis
		if err := repo.SaveAnamnesis(ctx, anamnesis); err != nil {
			t.Fatalf("Failed to save anamnesis: %v", err)
		}

		// Retrieve anamnesis
		retrieved, err := repo.GetAnamnesis(ctx, patientID)
		if err != nil {
			t.Fatalf("Failed to get anamnesis: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Retrieved anamnesis is nil")
		}

		// Verify fields
		if retrieved.PatientID != patientID {
			t.Errorf("PatientID mismatch: got %q, want %q", retrieved.PatientID, patientID)
		}
		if retrieved.ChiefComplaint != "Ansiedade generalizada há 2 anos" {
			t.Errorf("ChiefComplaint mismatch: got %q, want %q", retrieved.ChiefComplaint, "Ansiedade generalizada há 2 anos")
		}
		if retrieved.PersonalHistory != "Solteiro, trabalha como desenvolvedor" {
			t.Errorf("PersonalHistory mismatch: got %q, want %q", retrieved.PersonalHistory, "Solteiro, trabalha como desenvolvedor")
		}
	})

	t.Run("Update existing anamnesis", func(t *testing.T) {
		patientID := "test-patient-456"

		// Create initial anamnesis
		anamnesis, err := patient.NewAnamnesis(patientID)
		if err != nil {
			t.Fatalf("Failed to create anamnesis: %v", err)
		}

		if err := anamnesis.UpdateSection("chief_complaint", "Initial complaint"); err != nil {
			t.Fatalf("Failed to update chief complaint: %v", err)
		}

		if err := repo.SaveAnamnesis(ctx, anamnesis); err != nil {
			t.Fatalf("Failed to save initial anamnesis: %v", err)
		}

		// Update anamnesis
		if err := anamnesis.UpdateSection("chief_complaint", "Updated complaint"); err != nil {
			t.Fatalf("Failed to update chief complaint: %v", err)
		}

		if err := repo.SaveAnamnesis(ctx, anamnesis); err != nil {
			t.Fatalf("Failed to save updated anamnesis: %v", err)
		}

		// Retrieve and verify update
		retrieved, err := repo.GetAnamnesis(ctx, patientID)
		if err != nil {
			t.Fatalf("Failed to get updated anamnesis: %v", err)
		}

		if retrieved.ChiefComplaint != "Updated complaint" {
			t.Errorf("ChiefComplaint not updated: got %q, want %q", retrieved.ChiefComplaint, "Updated complaint")
		}
	})

	t.Run("Get non-existent anamnesis returns empty", func(t *testing.T) {
		nonExistentID := "non-existent-patient"
		anamnesis, err := repo.GetAnamnesis(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("Failed to get non-existent anamnesis: %v", err)
		}

		if anamnesis == nil {
			t.Fatal("Anamnesis should not be nil for non-existent patient")
		}

		if anamnesis.PatientID != nonExistentID {
			t.Errorf("PatientID mismatch: got %q, want %q", anamnesis.PatientID, nonExistentID)
		}

		if !anamnesis.IsEmpty() {
			t.Error("Anamnesis should be empty for non-existent patient")
		}
	})

	t.Run("Validate anamnesis sections", func(t *testing.T) {
		patientID := "test-patient-validation"
		anamnesis, err := patient.NewAnamnesis(patientID)
		if err != nil {
			t.Fatalf("Failed to create anamnesis: %v", err)
		}

		// Test valid sections
		validSections := []string{"chief_complaint", "personal_history", "family_history", "mental_state_exam"}
		for _, section := range validSections {
			if err := anamnesis.UpdateSection(section, "Test content"); err != nil {
				t.Errorf("Failed to update valid section %s: %v", section, err)
			}
		}

		// Test invalid section
		if err := anamnesis.UpdateSection("invalid_section", "Test content"); err == nil {
			t.Error("Expected error for invalid section")
		}
	})
}
