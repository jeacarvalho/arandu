package e2e

import (
	"context"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/infrastructure/repository/sqlite"
)

func TestBiopsychosocialMedicationFlow(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize database with migrations
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	patientRepo := sqlite.NewPatientRepository(db)
	medicationRepo := sqlite.NewMedicationRepository(db)

	// Create services
	patientService := services.NewPatientService(patientRepo)
	biopsychosocialService := services.NewBiopsychosocialService(medicationRepo, nil)

	// Create a test patient
	ctx := context.Background()
	input := services.CreatePatientInput{
		Name:  "Paciente Teste Medicação",
		Notes: "Paciente para testes de medicação",
	}
	testPatient, err := patientService.CreatePatient(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test patient: %v", err)
	}

	// Test 1: Add medication
	t.Run("AddMedication", func(t *testing.T) {
		med, err := biopsychosocialService.AddMedication(
			ctx,
			testPatient.ID,
			"Sertralina",
			"50mg",
			"Manhã",
			"Dr. Silva",
			time.Now(),
		)
		if err != nil {
			t.Errorf("AddMedication failed: %v", err)
		}

		if med == nil {
			t.Error("AddMedication returned nil")
		}

		if med.Name != "Sertralina" {
			t.Errorf("AddMedication name = %v, want %v", med.Name, "Sertralina")
		}

		if med.Dosage != "50mg" {
			t.Errorf("AddMedication dosage = %v, want %v", med.Dosage, "50mg")
		}

		if med.Status != patient.MedicationStatusActive {
			t.Errorf("AddMedication status = %v, want %v", med.Status, patient.MedicationStatusActive)
		}
	})

	// Test 2: Get medications for patient
	t.Run("GetMedications", func(t *testing.T) {
		meds, err := biopsychosocialService.GetMedications(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetMedications failed: %v", err)
		}

		if len(meds) != 1 {
			t.Errorf("GetMedications count = %v, want %v", len(meds), 1)
		}

		if meds[0].Name != "Sertralina" {
			t.Errorf("GetMedications[0].name = %v, want %v", meds[0].Name, "Sertralina")
		}
	})

	// Test 3: Suspend medication
	t.Run("SuspendMedication", func(t *testing.T) {
		meds, _ := biopsychosocialService.GetMedications(ctx, testPatient.ID)
		if len(meds) == 0 {
			t.Fatal("No medications to suspend")
		}

		suspendedMed, err := biopsychosocialService.SuspendMedication(ctx, meds[0].ID)
		if err != nil {
			t.Errorf("SuspendMedication failed: %v", err)
		}

		if suspendedMed.Status != patient.MedicationStatusSuspended {
			t.Errorf("SuspendMedication status = %v, want %v", suspendedMed.Status, patient.MedicationStatusSuspended)
		}

		if suspendedMed.EndedAt == nil {
			t.Error("SuspendMedication should set EndedAt")
		}
	})

	// Test 4: Activate suspended medication
	t.Run("ActivateMedication", func(t *testing.T) {
		meds, _ := biopsychosocialService.GetMedications(ctx, testPatient.ID)
		if len(meds) == 0 {
			t.Fatal("No medications to activate")
		}

		activatedMed, err := biopsychosocialService.ActivateMedication(ctx, meds[0].ID)
		if err != nil {
			t.Errorf("ActivateMedication failed: %v", err)
		}

		if activatedMed.Status != patient.MedicationStatusActive {
			t.Errorf("ActivateMedication status = %v, want %v", activatedMed.Status, patient.MedicationStatusActive)
		}

		if activatedMed.EndedAt != nil {
			t.Error("ActivateMedication should clear EndedAt")
		}
	})

	// Test 5: Get active medications
	t.Run("GetActiveMedications", func(t *testing.T) {
		activeMeds, err := medicationRepo.GetActiveMedications(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetActiveMedications failed: %v", err)
		}

		if len(activeMeds) != 1 {
			t.Errorf("GetActiveMedications count = %v, want %v", len(activeMeds), 1)
		}

		if activeMeds[0].Status != patient.MedicationStatusActive {
			t.Errorf("GetActiveMedications[0].status = %v, want %v", activeMeds[0].Status, patient.MedicationStatusActive)
		}
	})

	// Test 6: Add medication with empty name (should fail)
	t.Run("AddMedicationEmptyName", func(t *testing.T) {
		_, err := biopsychosocialService.AddMedication(
			ctx,
			testPatient.ID,
			"",
			"50mg",
			"Manhã",
			"Dr. Silva",
			time.Now(),
		)
		if err == nil {
			t.Error("AddMedication with empty name should fail")
		}

		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("AddMedication with empty name error = %v, want error containing 'cannot be empty'", err)
		}
	})

	// Test 7: Add medication with future date (should fail)
	t.Run("AddMedicationFutureDate", func(t *testing.T) {
		futureDate := time.Now().Add(24 * time.Hour)
		_, err := biopsychosocialService.AddMedication(
			ctx,
			testPatient.ID,
			"Paroxetina",
			"20mg",
			"Noite",
			"Dr. Santos",
			futureDate,
		)
		if err == nil {
			t.Error("AddMedication with future date should fail")
		}

		if !strings.Contains(err.Error(), "cannot be in the future") {
			t.Errorf("AddMedication with future date error = %v, want error containing 'cannot be in the future'", err)
		}
	})
}

func TestBiopsychosocialVitalsFlow(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize database with migrations
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	patientRepo := sqlite.NewPatientRepository(db)
	vitalsRepo := sqlite.NewVitalsRepository(db)

	// Create services
	patientService := services.NewPatientService(patientRepo)
	biopsychosocialService := services.NewBiopsychosocialService(nil, vitalsRepo)

	// Create a test patient
	ctx := context.Background()
	input := services.CreatePatientInput{
		Name:  "Paciente Teste Vitais",
		Notes: "Paciente para testes de sinais vitais",
	}
	testPatient, err := patientService.CreatePatient(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test patient: %v", err)
	}

	// Test 1: Record vitals
	t.Run("RecordVitals", func(t *testing.T) {
		sleepHours := 8.0
		appetiteLevel := 7
		weight := 70.5

		vitals, err := biopsychosocialService.RecordVitals(
			ctx,
			testPatient.ID,
			time.Now(),
			&sleepHours,
			&appetiteLevel,
			&weight,
			3,
			"Paciente dormiu bem",
		)
		if err != nil {
			t.Errorf("RecordVitals failed: %v", err)
		}

		if vitals == nil {
			t.Error("RecordVitals returned nil")
		}

		if *vitals.SleepHours != 8.0 {
			t.Errorf("RecordVitals sleepHours = %v, want %v", *vitals.SleepHours, 8.0)
		}

		if *vitals.AppetiteLevel != 7 {
			t.Errorf("RecordVitals appetiteLevel = %v, want %v", *vitals.AppetiteLevel, 7)
		}

		if *vitals.Weight != 70.5 {
			t.Errorf("RecordVitals weight = %v, want %v", *vitals.Weight, 70.5)
		}

		if vitals.PhysicalActivity != 3 {
			t.Errorf("RecordVitals physicalActivity = %v, want %v", vitals.PhysicalActivity, 3)
		}
	})

	// Test 2: Get latest vitals
	t.Run("GetLatestVitals", func(t *testing.T) {
		latest, err := biopsychosocialService.GetLatestVitals(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetLatestVitals failed: %v", err)
		}

		if latest == nil {
			t.Error("GetLatestVitals returned nil")
		}

		if *latest.SleepHours != 8.0 {
			t.Errorf("GetLatestVitals sleepHours = %v, want %v", *latest.SleepHours, 8.0)
		}
	})

	// Test 3: Record vitals with partial data
	t.Run("RecordVitalsPartialData", func(t *testing.T) {
		sleepHours := 7.5
		appetiteLevel := 6

		vitals, err := biopsychosocialService.RecordVitals(
			ctx,
			testPatient.ID,
			time.Now(),
			&sleepHours,
			&appetiteLevel,
			nil,
			0,
			"",
		)
		if err != nil {
			t.Errorf("RecordVitals with partial data failed: %v", err)
		}

		if vitals.Weight != nil {
			t.Error("RecordVitals with nil weight should keep weight as nil")
		}
	})

	// Test 4: Record vitals with invalid appetite level (should fail)
	t.Run("RecordVitalsInvalidAppetite", func(t *testing.T) {
		sleepHours := 8.0
		appetiteLevel := 0 // Invalid: less than 1
		weight := 70.5

		_, err := biopsychosocialService.RecordVitals(
			ctx,
			testPatient.ID,
			time.Now(),
			&sleepHours,
			&appetiteLevel,
			&weight,
			3,
			"",
		)
		if err == nil {
			t.Error("RecordVitals with appetite level 0 should fail")
		}

		if !strings.Contains(err.Error(), "must be between 1 and 10") {
			t.Errorf("RecordVitals with invalid appetite error = %v, want error containing 'must be between 1 and 10'", err)
		}
	})

	// Test 5: Record vitals with invalid sleep hours (should fail)
	t.Run("RecordVitalsInvalidSleep", func(t *testing.T) {
		sleepHours := 25.0 // Invalid: greater than 24
		appetiteLevel := 5
		weight := 70.5

		_, err := biopsychosocialService.RecordVitals(
			ctx,
			testPatient.ID,
			time.Now(),
			&sleepHours,
			&appetiteLevel,
			&weight,
			3,
			"",
		)
		if err == nil {
			t.Error("RecordVitals with sleep hours > 24 should fail")
		}

		if !strings.Contains(err.Error(), "must be between 0 and 24") {
			t.Errorf("RecordVitals with invalid sleep error = %v, want error containing 'must be between 0 and 24'", err)
		}
	})

	// Test 6: Get average vitals
	t.Run("GetAverageVitals", func(t *testing.T) {
		// Record multiple vitals for averaging
		for i := 0; i < 3; i++ {
			sleepHours := 7.0 + float64(i)
			appetiteLevel := 6 + i
			weight := 70.0 + float64(i)
			date := time.Now().AddDate(0, 0, -i) // Different dates

			_, err := biopsychosocialService.RecordVitals(
				ctx,
				testPatient.ID,
				date,
				&sleepHours,
				&appetiteLevel,
				&weight,
				i,
				"",
			)
			if err != nil {
				t.Errorf("Failed to record vitals for averaging: %v", err)
			}
		}

		avg, err := vitalsRepo.GetAverageVitals(ctx, testPatient.ID, 30)
		if err != nil {
			t.Errorf("GetAverageVitals failed: %v", err)
		}

		if avg == nil {
			t.Error("GetAverageVitals returned nil")
		}

		if avg.Count != 5 { // 1 from RecordVitals + 1 from RecordVitalsPartialData + 3 new ones
			t.Errorf("GetAverageVitals count = %v, want %v", avg.Count, 5)
		}
	})

	// Test 7: Find vitals by patient ID
	t.Run("FindVitalsByPatientID", func(t *testing.T) {
		vitalsList, err := vitalsRepo.FindByPatientID(ctx, testPatient.ID, 30)
		if err != nil {
			t.Errorf("FindVitalsByPatientID failed: %v", err)
		}

		if len(vitalsList) != 5 { // 1 from RecordVitals + 1 from RecordVitalsPartialData + 3 from averaging test
			t.Errorf("FindVitalsByPatientID count = %v, want %v", len(vitalsList), 5)
		}
	})
}

func TestBiopsychosocialContextFlow(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize database with migrations
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	patientRepo := sqlite.NewPatientRepository(db)
	medicationRepo := sqlite.NewMedicationRepository(db)
	vitalsRepo := sqlite.NewVitalsRepository(db)

	// Create services
	patientService := services.NewPatientService(patientRepo)
	biopsychosocialService := services.NewBiopsychosocialService(medicationRepo, vitalsRepo)

	// Create a test patient
	ctx := context.Background()
	input := services.CreatePatientInput{
		Name:  "Paciente Teste Contexto",
		Notes: "Paciente para testes de contexto completo",
	}
	testPatient, err := patientService.CreatePatient(ctx, input)
	if err != nil {
		t.Fatalf("Failed to create test patient: %v", err)
	}

	// Test 1: Get empty context
	t.Run("GetEmptyContext", func(t *testing.T) {
		context, err := biopsychosocialService.GetContext(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetContext failed: %v", err)
		}

		if context == nil {
			t.Error("GetContext returned nil")
		}

		if context.PatientID != testPatient.ID {
			t.Errorf("GetContext patientID = %v, want %v", context.PatientID, testPatient.ID)
		}

		if len(context.AllMedications) != 0 {
			t.Errorf("GetContext allMedications count = %v, want %v", len(context.AllMedications), 0)
		}

		if context.LatestVitals != nil {
			t.Error("GetContext latestVitals should be nil for empty context")
		}
	})

	// Test 2: Get context with data
	t.Run("GetContextWithData", func(t *testing.T) {
		// Add medication
		_, err := biopsychosocialService.AddMedication(
			ctx,
			testPatient.ID,
			"Fluoxetina",
			"20mg",
			"Manhã",
			"Dr. Costa",
			time.Now(),
		)
		if err != nil {
			t.Fatalf("Failed to add medication for context test: %v", err)
		}

		// Record vitals
		sleepHours := 7.5
		appetiteLevel := 8
		weight := 68.0
		_, err = biopsychosocialService.RecordVitals(
			ctx,
			testPatient.ID,
			time.Now(),
			&sleepHours,
			&appetiteLevel,
			&weight,
			4,
			"Paciente estável",
		)
		if err != nil {
			t.Fatalf("Failed to record vitals for context test: %v", err)
		}

		// Get context
		context, err := biopsychosocialService.GetContext(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetContext with data failed: %v", err)
		}

		if len(context.AllMedications) != 1 {
			t.Errorf("GetContext allMedications count = %v, want %v", len(context.AllMedications), 1)
		}

		if context.LatestVitals == nil {
			t.Error("GetContext latestVitals should not be nil")
		}

		if *context.LatestVitals.SleepHours != 7.5 {
			t.Errorf("GetContext latestVitals sleepHours = %v, want %v", *context.LatestVitals.SleepHours, 7.5)
		}

		// Check average vitals
		avg, err := vitalsRepo.GetAverageVitals(ctx, testPatient.ID, 30)
		if err != nil {
			t.Errorf("GetAverageVitals for context failed: %v", err)
		}

		if avg == nil {
			t.Error("GetAverageVitals returned nil")
		}

		if avg.Count != 1 {
			t.Errorf("GetAverageVitals count = %v, want %v", avg.Count, 1)
		}
	})

	// Test 3: Get context with multiple medications
	t.Run("GetContextMultipleMedications", func(t *testing.T) {
		// Add second medication
		_, err := biopsychosocialService.AddMedication(
			ctx,
			testPatient.ID,
			"Clonazepam",
			"0.5mg",
			"Noite",
			"Dr. Silva",
			time.Now(),
		)
		if err != nil {
			t.Fatalf("Failed to add second medication: %v", err)
		}

		context, err := biopsychosocialService.GetContext(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetContext with multiple medications failed: %v", err)
		}

		if len(context.AllMedications) != 2 {
			t.Errorf("GetContext allMedications count = %v, want %v", len(context.AllMedications), 2)
		}

		// Check medications are sorted by started_at (newest first)
		if context.AllMedications[0].Name != "Clonazepam" {
			t.Errorf("GetContext medications[0].name = %v, want %v", context.AllMedications[0].Name, "Clonazepam")
		}
	})

	// Test 4: Get context with suspended medication
	t.Run("GetContextWithSuspendedMedication", func(t *testing.T) {
		meds, _ := biopsychosocialService.GetMedications(ctx, testPatient.ID)
		if len(meds) == 0 {
			t.Fatal("No medications to suspend")
		}

		// Suspend first medication
		_, err := biopsychosocialService.SuspendMedication(ctx, meds[0].ID)
		if err != nil {
			t.Fatalf("Failed to suspend medication: %v", err)
		}

		context, err := biopsychosocialService.GetContext(ctx, testPatient.ID)
		if err != nil {
			t.Errorf("GetContext with suspended medication failed: %v", err)
		}

		// Count active medications
		activeCount := 0
		for _, med := range context.AllMedications {
			if med.Status == patient.MedicationStatusActive {
				activeCount++
			}
		}

		if activeCount != 1 { // One active (Fluoxetina), one suspended (Clonazepam)
			t.Errorf("GetContext active medications count = %v, want %v", activeCount, 1)
		}
	})
}
