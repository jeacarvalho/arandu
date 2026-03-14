package patient

import (
	"testing"
	"time"
)

func TestNewPatient(t *testing.T) {
	tests := []struct {
		name      string
		notes     string
		wantError bool
	}{
		{"John Doe", "Initial notes", false},
		{"Jane Smith", "", false},
		{"", "Should fail", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patient, err := NewPatient(tt.name, tt.notes)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewPatient(%q, %q) should have returned an error", tt.name, tt.notes)
				}
				return
			}

			if err != nil {
				t.Errorf("NewPatient(%q, %q) returned unexpected error: %v", tt.name, tt.notes, err)
				return
			}

			if patient.Name != tt.name {
				t.Errorf("NewPatient(%q, %q) patient.Name = %q, want %q", tt.name, tt.notes, patient.Name, tt.name)
			}

			if patient.Notes != tt.notes {
				t.Errorf("NewPatient(%q, %q) patient.Notes = %q, want %q", tt.name, tt.notes, patient.Notes, tt.notes)
			}

			if patient.ID == "" {
				t.Errorf("NewPatient(%q, %q) patient.ID is empty", tt.name, tt.notes)
			}

			if patient.CreatedAt.IsZero() {
				t.Errorf("NewPatient(%q, %q) patient.CreatedAt is zero", tt.name, tt.notes)
			}

			if patient.UpdatedAt.IsZero() {
				t.Errorf("NewPatient(%q, %q) patient.UpdatedAt is zero", tt.name, tt.notes)
			}
		})
	}
}

func TestPatient_Update(t *testing.T) {
	// Create a patient first
	patient, err := NewPatient("Original Name", "Original Notes")
	if err != nil {
		t.Fatalf("Failed to create patient: %v", err)
	}

	originalUpdatedAt := patient.UpdatedAt

	// Test successful update
	t.Run("Successful update", func(t *testing.T) {
		newName := "Updated Name"
		newNotes := "Updated Notes"

		// Wait a bit to ensure timestamp changes
		time.Sleep(1 * time.Millisecond)

		err := patient.Update(newName, newNotes)
		if err != nil {
			t.Errorf("Update(%q, %q) returned unexpected error: %v", newName, newNotes, err)
		}

		if patient.Name != newName {
			t.Errorf("After Update, patient.Name = %q, want %q", patient.Name, newName)
		}

		if patient.Notes != newNotes {
			t.Errorf("After Update, patient.Notes = %q, want %q", patient.Notes, newNotes)
		}

		if !patient.UpdatedAt.After(originalUpdatedAt) {
			t.Errorf("After Update, patient.UpdatedAt = %v, should be after original %v", patient.UpdatedAt, originalUpdatedAt)
		}
	})

	// Test update with empty name (should fail)
	t.Run("Update with empty name", func(t *testing.T) {
		originalName := patient.Name
		originalNotes := patient.Notes
		originalUpdatedAt := patient.UpdatedAt

		err := patient.Update("", "New Notes")
		if err == nil {
			t.Error("Update with empty name should have returned an error")
		}

		// Patient should remain unchanged
		if patient.Name != originalName {
			t.Errorf("After failed Update, patient.Name changed to %q, should remain %q", patient.Name, originalName)
		}

		if patient.Notes != originalNotes {
			t.Errorf("After failed Update, patient.Notes changed to %q, should remain %q", patient.Notes, originalNotes)
		}

		if patient.UpdatedAt != originalUpdatedAt {
			t.Errorf("After failed Update, patient.UpdatedAt changed to %v, should remain %v", patient.UpdatedAt, originalUpdatedAt)
		}
	})

	// Test update with same name but different notes
	t.Run("Update notes only", func(t *testing.T) {
		currentName := patient.Name
		newNotes := "Different notes"
		originalUpdatedAt := patient.UpdatedAt

		// Wait a bit to ensure timestamp changes
		time.Sleep(1 * time.Millisecond)

		err := patient.Update(currentName, newNotes)
		if err != nil {
			t.Errorf("Update(%q, %q) returned unexpected error: %v", currentName, newNotes, err)
		}

		if patient.Name != currentName {
			t.Errorf("After Update, patient.Name changed to %q, should remain %q", patient.Name, currentName)
		}

		if patient.Notes != newNotes {
			t.Errorf("After Update, patient.Notes = %q, want %q", patient.Notes, newNotes)
		}

		if !patient.UpdatedAt.After(originalUpdatedAt) {
			t.Errorf("After Update, patient.UpdatedAt = %v, should be after original %v", patient.UpdatedAt, originalUpdatedAt)
		}
	})
}
