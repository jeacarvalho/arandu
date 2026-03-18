package sqlite

import (
	"testing"
	"time"

	"arandu/internal/domain/patient"
)

func TestMedicationRepository(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	repo := NewMedicationRepository(db)
	patientRepo := NewPatientRepository(db)

	patientObj := &patient.Patient{
		ID:        "test-patient-1",
		Name:      "Test Patient",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := patientRepo.Save(patientObj); err != nil {
		t.Fatalf("Failed to save patient: %v", err)
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		med := &patient.Medication{
			ID:         "med-1",
			PatientID:  patientObj.ID,
			Name:       "Sertralina",
			Dosage:     "50mg",
			Frequency:  "Manhã",
			Prescriber: "Dr. João",
			Status:     patient.MedicationStatusActive,
			StartedAt:  time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := repo.Save(med); err != nil {
			t.Fatalf("Failed to save medication: %v", err)
		}

		found, err := repo.FindByID("med-1")
		if err != nil {
			t.Fatalf("Failed to find medication: %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find medication")
		}
		if found.Name != "Sertralina" {
			t.Errorf("Expected name Sertralina, got %s", found.Name)
		}
	})

	t.Run("GetActiveMedications", func(t *testing.T) {
		activeMeds, err := repo.GetActiveMedications(patientObj.ID)
		if err != nil {
			t.Fatalf("Failed to get active medications: %v", err)
		}
		if len(activeMeds) == 0 {
			t.Error("Expected at least one active medication")
		}
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		if err := repo.UpdateStatus("med-1", patient.MedicationStatusSuspended); err != nil {
			t.Fatalf("Failed to update status: %v", err)
		}

		updated, err := repo.FindByID("med-1")
		if err != nil {
			t.Fatalf("Failed to find updated medication: %v", err)
		}
		if updated.Status != patient.MedicationStatusSuspended {
			t.Errorf("Expected status suspended, got %s", updated.Status)
		}
		if updated.EndedAt == nil {
			t.Error("Expected ended_at to be set after suspension")
		}
	})

	t.Run("FindByPatientID", func(t *testing.T) {
		meds, err := repo.FindByPatientID(patientObj.ID)
		if err != nil {
			t.Fatalf("Failed to find medications by patient: %v", err)
		}
		if len(meds) != 1 {
			t.Errorf("Expected 1 medication, got %d", len(meds))
		}
	})
}

func TestVitalsRepository(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	repo := NewVitalsRepository(db)
	patientRepo := NewPatientRepository(db)

	patientObj := &patient.Patient{
		ID:        "test-patient-2",
		Name:      "Test Patient 2",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := patientRepo.Save(patientObj); err != nil {
		t.Fatalf("Failed to save patient: %v", err)
	}

	sleep := 8.0
	appetite := 7
	weight := 70.5

	t.Run("Save and FindByID", func(t *testing.T) {
		vitals := &patient.Vitals{
			ID:               "vitals-1",
			PatientID:        patientObj.ID,
			Date:             time.Now(),
			SleepHours:       &sleep,
			AppetiteLevel:    &appetite,
			Weight:           &weight,
			PhysicalActivity: 3,
			Notes:            "Test notes",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := repo.Save(vitals); err != nil {
			t.Fatalf("Failed to save vitals: %v", err)
		}

		found, err := repo.FindByID("vitals-1")
		if err != nil {
			t.Fatalf("Failed to find vitals: %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find vitals")
		}
		if *found.SleepHours != 8.0 {
			t.Errorf("Expected sleep hours 8.0, got %f", *found.SleepHours)
		}
	})

	t.Run("GetLatestVitals", func(t *testing.T) {
		latest, err := repo.GetLatestVitals(patientObj.ID)
		if err != nil {
			t.Fatalf("Failed to get latest vitals: %v", err)
		}
		if latest == nil {
			t.Fatal("Expected to find latest vitals")
		}
		if *latest.SleepHours != 8.0 {
			t.Errorf("Expected sleep hours 8.0, got %f", *latest.SleepHours)
		}
	})

	t.Run("GetAverageVitals", func(t *testing.T) {
		avg, err := repo.GetAverageVitals(patientObj.ID, 30)
		if err != nil {
			t.Fatalf("Failed to get average vitals: %v", err)
		}
		if avg == nil {
			t.Fatal("Expected to find average vitals")
		}
		if avg.Count != 1 {
			t.Errorf("Expected 1 record, got %d", avg.Count)
		}
	})

	t.Run("FindByPatientID", func(t *testing.T) {
		vitals, err := repo.FindByPatientID(patientObj.ID, 10)
		if err != nil {
			t.Fatalf("Failed to find vitals by patient: %v", err)
		}
		if len(vitals) != 1 {
			t.Errorf("Expected 1 vitals record, got %d", len(vitals))
		}
	})
}
