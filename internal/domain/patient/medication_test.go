package patient

import (
	"testing"
	"time"
)

func TestNewMedication(t *testing.T) {
	patientID := "patient-123"
	name := "Sertralina"
	dosage := "50mg"
	frequency := "Manhã"
	prescriber := "Dr. João"
	startedAt := time.Now().AddDate(0, 0, -1)

	med, err := NewMedication(patientID, name, dosage, frequency, prescriber, startedAt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if med.ID == "" {
		t.Error("Expected medication ID to be set")
	}
	if med.PatientID != patientID {
		t.Errorf("Expected patient ID %s, got %s", patientID, med.PatientID)
	}
	if med.Name != name {
		t.Errorf("Expected name %s, got %s", name, med.Name)
	}
	if med.Dosage != dosage {
		t.Errorf("Expected dosage %s, got %s", dosage, med.Dosage)
	}
	if med.Frequency != frequency {
		t.Errorf("Expected frequency %s, got %s", frequency, med.Frequency)
	}
	if med.Prescriber != prescriber {
		t.Errorf("Expected prescriber %s, got %s", prescriber, med.Prescriber)
	}
	if med.Status != MedicationStatusActive {
		t.Errorf("Expected status %s, got %s", MedicationStatusActive, med.Status)
	}
	if med.EndedAt != nil {
		t.Error("Expected ended_at to be nil for active medication")
	}
}

func TestNewMedicationValidation(t *testing.T) {
	tests := []struct {
		name      string
		patientID string
		medName   string
		startedAt time.Time
		expectErr bool
		errMsg    string
	}{
		{
			name:      "empty patient ID",
			patientID: "",
			medName:   "Sertralina",
			startedAt: time.Now(),
			expectErr: true,
			errMsg:    "patient ID cannot be empty",
		},
		{
			name:      "empty medication name",
			patientID: "patient-123",
			medName:   "",
			startedAt: time.Now(),
			expectErr: true,
			errMsg:    "medication name cannot be empty",
		},
		{
			name:      "future start date",
			patientID: "patient-123",
			medName:   "Sertralina",
			startedAt: time.Now().AddDate(0, 0, 1),
			expectErr: true,
			errMsg:    "medication start date cannot be in the future",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMedication(tt.patientID, tt.medName, "", "", "", tt.startedAt)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error %s, got nil", tt.errMsg)
				} else if err.Error() != tt.errMsg {
					t.Errorf("Expected error %s, got %s", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestMedicationSuspend(t *testing.T) {
	med, _ := NewMedication("patient-123", "Sertralina", "50mg", "Manhã", "", time.Now())

	if err := med.Suspend(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if med.Status != MedicationStatusSuspended {
		t.Errorf("Expected status %s, got %s", MedicationStatusSuspended, med.Status)
	}
	if med.EndedAt == nil {
		t.Error("Expected ended_at to be set after suspension")
	}
}

func TestMedicationSuspendFinished(t *testing.T) {
	med, _ := NewMedication("patient-123", "Sertralina", "50mg", "Manhã", "", time.Now())
	med.Finish()

	err := med.Suspend()
	if err == nil {
		t.Error("Expected error when suspending finished medication")
	}
}

func TestMedicationActivate(t *testing.T) {
	med, _ := NewMedication("patient-123", "Sertralina", "50mg", "Manhã", "", time.Now())
	med.Suspend()

	if err := med.Activate(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if med.Status != MedicationStatusActive {
		t.Errorf("Expected status %s, got %s", MedicationStatusActive, med.Status)
	}
	if med.EndedAt != nil {
		t.Error("Expected ended_at to be nil after activation")
	}
}

func TestMedicationFinish(t *testing.T) {
	med, _ := NewMedication("patient-123", "Sertralina", "50mg", "Manhã", "", time.Now())

	if err := med.Finish(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if med.Status != MedicationStatusFinished {
		t.Errorf("Expected status %s, got %s", MedicationStatusFinished, med.Status)
	}
	if med.EndedAt == nil {
		t.Error("Expected ended_at to be set after finish")
	}
}

func TestMedicationIsActive(t *testing.T) {
	med, _ := NewMedication("patient-123", "Sertralina", "50mg", "Manhã", "", time.Now())

	if !med.IsActive() {
		t.Error("Expected medication to be active")
	}

	med.Suspend()
	if med.IsActive() {
		t.Error("Expected medication to not be active after suspension")
	}
}
