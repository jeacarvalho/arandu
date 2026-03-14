package sqlite

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"arandu/internal/domain/patient"
)

func TestPatientRepositoryIntegration(t *testing.T) {
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

	// Create repository
	repo := NewPatientRepository(db)

	// Initialize schema
	// Note: In production, use db.Migrate() with the correct migrations directory
	// For unit tests, we create the table directly to avoid path issues
	query := `
	CREATE TABLE IF NOT EXISTS patients (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		notes TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
	`
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("Failed to create patients table: %v", err)
	}

	t.Run("Create and retrieve patient", func(t *testing.T) {
		// Create a patient using domain constructor
		p, err := patient.NewPatient("Test Patient", "Test Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}

		// Save patient
		if err := repo.Save(p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Retrieve patient
		retrieved, err := repo.FindByID(p.ID)
		if err != nil {
			t.Fatalf("Failed to find patient: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Retrieved patient is nil")
		}

		// Verify fields
		if retrieved.ID != p.ID {
			t.Errorf("ID mismatch: got %q, want %q", retrieved.ID, p.ID)
		}
		if retrieved.Name != p.Name {
			t.Errorf("Name mismatch: got %q, want %q", retrieved.Name, p.Name)
		}
		if retrieved.Notes != p.Notes {
			t.Errorf("Notes mismatch: got %q, want %q", retrieved.Notes, p.Notes)
		}
		if !retrieved.CreatedAt.Equal(p.CreatedAt) {
			t.Errorf("CreatedAt mismatch: got %v, want %v", retrieved.CreatedAt, p.CreatedAt)
		}
		if !retrieved.UpdatedAt.Equal(p.UpdatedAt) {
			t.Errorf("UpdatedAt mismatch: got %v, want %v", retrieved.UpdatedAt, p.UpdatedAt)
		}
	})

	t.Run("FindAll returns patients in correct order", func(t *testing.T) {
		// Clear existing patients
		allPatients, _ := repo.FindAll()
		for _, p := range allPatients {
			repo.Delete(p.ID)
		}

		// Create multiple patients with slight time differences
		patients := []*patient.Patient{}
		for i := 1; i <= 3; i++ {
			p, err := patient.NewPatient(
				"Patient "+string(rune('A'+i-1)),
				"Notes "+string(rune('A'+i-1)),
			)
			if err != nil {
				t.Fatalf("Failed to create patient %d: %v", i, err)
			}

			// Small delay to ensure different timestamps
			time.Sleep(1 * time.Millisecond)

			if err := repo.Save(p); err != nil {
				t.Fatalf("Failed to save patient %d: %v", i, err)
			}
			patients = append(patients, p)
		}

		// Retrieve all patients
		all, err := repo.FindAll()
		if err != nil {
			t.Fatalf("Failed to find all patients: %v", err)
		}

		if len(all) != 3 {
			t.Fatalf("Expected 3 patients, got %d", len(all))
		}

		// Verify order (should be reverse chronological - newest first)
		// Since we created patients in order A, B, C with delays,
		// they should be returned as C, B, A
		for i := 0; i < 3; i++ {
			expected := patients[2-i] // Reverse order
			actual := all[i]

			if actual.ID != expected.ID {
				t.Errorf("Patient %d ID mismatch: got %q, want %q", i, actual.ID, expected.ID)
			}
		}
	})

	t.Run("Update patient", func(t *testing.T) {
		// Create a patient
		p, err := patient.NewPatient("Original Name", "Original Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}

		if err := repo.Save(p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Update using domain method
		originalUpdatedAt := p.UpdatedAt
		time.Sleep(1 * time.Millisecond) // Ensure timestamp changes

		if err := p.Update("Updated Name", "Updated Notes"); err != nil {
			t.Fatalf("Failed to update patient: %v", err)
		}

		// Save update
		if err := repo.Update(p); err != nil {
			t.Fatalf("Failed to save update: %v", err)
		}

		// Retrieve and verify
		retrieved, err := repo.FindByID(p.ID)
		if err != nil {
			t.Fatalf("Failed to find updated patient: %v", err)
		}

		if retrieved.Name != "Updated Name" {
			t.Errorf("Name not updated: got %q, want %q", retrieved.Name, "Updated Name")
		}
		if retrieved.Notes != "Updated Notes" {
			t.Errorf("Notes not updated: got %q, want %q", retrieved.Notes, "Updated Notes")
		}
		if !retrieved.UpdatedAt.After(originalUpdatedAt) {
			t.Errorf("UpdatedAt not updated: got %v, original was %v", retrieved.UpdatedAt, originalUpdatedAt)
		}
	})

	t.Run("Delete patient", func(t *testing.T) {
		// Create a patient
		p, err := patient.NewPatient("To Delete", "Delete Notes")
		if err != nil {
			t.Fatalf("Failed to create patient: %v", err)
		}

		if err := repo.Save(p); err != nil {
			t.Fatalf("Failed to save patient: %v", err)
		}

		// Verify exists
		found, err := repo.FindByID(p.ID)
		if err != nil || found == nil {
			t.Fatal("Patient should exist before deletion")
		}

		// Delete
		if err := repo.Delete(p.ID); err != nil {
			t.Fatalf("Failed to delete patient: %v", err)
		}

		// Verify deleted
		found, err = repo.FindByID(p.ID)
		if err != nil {
			t.Fatalf("Error finding deleted patient: %v", err)
		}
		if found != nil {
			t.Fatal("Patient should not exist after deletion")
		}
	})

	t.Run("FindByID returns nil for non-existent patient", func(t *testing.T) {
		nonExistentID := "non-existent-id"
		p, err := repo.FindByID(nonExistentID)
		if err != nil {
			t.Fatalf("FindByID should not error for non-existent ID: %v", err)
		}
		if p != nil {
			t.Fatal("FindByID should return nil for non-existent ID")
		}
	})

	t.Run("FindByName searches patients", func(t *testing.T) {
		// Clear existing patients
		allPatients, _ := repo.FindAll()
		for _, p := range allPatients {
			repo.Delete(p.ID)
		}

		// Create test patients with different names
		patients := []struct {
			name  string
			notes string
		}{
			{"John Doe", "Patient 1"},
			{"Jane Smith", "Patient 2"},
			{"Robert Robertson", "Patient 3"}, // Changed from Johnson to avoid John match
			{"Johnny Appleseed", "Patient 4"},
		}

		for _, pt := range patients {
			p, err := patient.NewPatient(pt.name, pt.notes)
			if err != nil {
				t.Fatalf("Failed to create patient: %v", err)
			}
			if err := repo.Save(p); err != nil {
				t.Fatalf("Failed to save patient: %v", err)
			}
		}

		// Test search for "John" (should find John Doe and Johnny Appleseed)
		results, err := repo.FindByName("John")
		if err != nil {
			t.Fatalf("FindByName failed: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 patients named 'John', got %d", len(results))
		}

		// Verify names contain "John"
		for _, p := range results {
			if !containsIgnoreCase(p.Name, "John") {
				t.Errorf("Patient name %q does not contain 'John'", p.Name)
			}
		}

		// Test case-insensitive search
		results, err = repo.FindByName("jOhN")
		if err != nil {
			t.Fatalf("FindByName (case-insensitive) failed: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 patients (case-insensitive), got %d", len(results))
		}

		// Test search for non-existent name
		results, err = repo.FindByName("Nonexistent")
		if err != nil {
			t.Fatalf("FindByName (non-existent) failed: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 patients for 'Nonexistent', got %d", len(results))
		}
	})

	t.Run("CountAll returns correct count", func(t *testing.T) {
		// Clear existing patients
		allPatients, _ := repo.FindAll()
		for _, p := range allPatients {
			repo.Delete(p.ID)
		}

		// Create 5 patients
		for i := 1; i <= 5; i++ {
			p, err := patient.NewPatient(fmt.Sprintf("Patient %d", i), fmt.Sprintf("Notes %d", i))
			if err != nil {
				t.Fatalf("Failed to create patient: %v", err)
			}
			if err := repo.Save(p); err != nil {
				t.Fatalf("Failed to save patient: %v", err)
			}
		}

		count, err := repo.CountAll()
		if err != nil {
			t.Fatalf("CountAll failed: %v", err)
		}

		if count != 5 {
			t.Errorf("Expected 5 patients, got %d", count)
		}

		// Delete one and count again
		allPatients, _ = repo.FindAll()
		if len(allPatients) > 0 {
			repo.Delete(allPatients[0].ID)

			count, err = repo.CountAll()
			if err != nil {
				t.Fatalf("CountAll after delete failed: %v", err)
			}

			if count != 4 {
				t.Errorf("Expected 4 patients after delete, got %d", count)
			}
		}
	})

	t.Run("FindPaginated returns paginated results", func(t *testing.T) {
		// Clear existing patients
		allPatients, _ := repo.FindAll()
		for _, p := range allPatients {
			repo.Delete(p.ID)
		}

		// Create 10 patients
		for i := 1; i <= 10; i++ {
			p, err := patient.NewPatient(fmt.Sprintf("Patient %d", i), fmt.Sprintf("Notes %d", i))
			if err != nil {
				t.Fatalf("Failed to create patient: %v", err)
			}
			// Small delay to ensure different timestamps for ordering
			time.Sleep(1 * time.Millisecond)
			if err := repo.Save(p); err != nil {
				t.Fatalf("Failed to save patient: %v", err)
			}
		}

		// Test first page (limit 3, offset 0)
		page1, err := repo.FindPaginated(3, 0)
		if err != nil {
			t.Fatalf("FindPaginated page 1 failed: %v", err)
		}

		if len(page1) != 3 {
			t.Errorf("Expected 3 patients on page 1, got %d", len(page1))
		}

		// Test second page (limit 3, offset 3)
		page2, err := repo.FindPaginated(3, 3)
		if err != nil {
			t.Fatalf("FindPaginated page 2 failed: %v", err)
		}

		if len(page2) != 3 {
			t.Errorf("Expected 3 patients on page 2, got %d", len(page2))
		}

		// Verify pages don't overlap
		page1IDs := make(map[string]bool)
		for _, p := range page1 {
			page1IDs[p.ID] = true
		}
		for _, p := range page2 {
			if page1IDs[p.ID] {
				t.Errorf("Patient %s appears on both page 1 and page 2", p.ID)
			}
		}

		// Test last page with partial results (limit 3, offset 9)
		page4, err := repo.FindPaginated(3, 9)
		if err != nil {
			t.Fatalf("FindPaginated page 4 failed: %v", err)
		}

		if len(page4) != 1 {
			t.Errorf("Expected 1 patient on page 4, got %d", len(page4))
		}

		// Test invalid parameters
		_, err = repo.FindPaginated(0, 0)
		if err == nil {
			t.Error("Expected error for limit=0")
		}

		_, err = repo.FindPaginated(101, 0)
		if err == nil {
			t.Error("Expected error for limit>100")
		}

		_, err = repo.FindPaginated(10, -1)
		if err == nil {
			t.Error("Expected error for negative offset")
		}
	})

	t.Run("Validation prevents invalid operations", func(t *testing.T) {
		// Test Save with invalid patient
		invalidPatient := &patient.Patient{
			ID:   "", // Empty ID
			Name: "Test",
		}

		err := repo.Save(invalidPatient)
		if err == nil {
			t.Error("Expected error when saving patient with empty ID")
		}

		// Test FindByID with empty ID
		_, err = repo.FindByID("")
		if err == nil {
			t.Error("Expected error when finding patient with empty ID")
		}

		// Test Delete with empty ID
		err = repo.Delete("")
		if err == nil {
			t.Error("Expected error when deleting patient with empty ID")
		}

		// Test FindByName with empty name
		_, err = repo.FindByName("")
		if err == nil {
			t.Error("Expected error when searching with empty name")
		}

		// Test FindByName with very long name
		longName := strings.Repeat("a", 101)
		_, err = repo.FindByName(longName)
		if err == nil {
			t.Error("Expected error when searching with name > 100 characters")
		}
	})
}

// Helper function to check if string contains substring (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
