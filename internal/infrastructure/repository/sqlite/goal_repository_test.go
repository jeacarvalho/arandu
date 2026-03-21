package sqlite

import (
	"context"
	"testing"
	"time"

	"arandu/internal/domain/patient"
)

func TestGoalRepository(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	repo := NewGoalRepository(db)
	patientRepo := NewPatientRepository(db)

	ctx := context.Background()

	patientObj := &patient.Patient{
		ID:        "test-patient-goal-1",
		Name:      "Test Patient for Goals",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := patientRepo.Save(ctx, patientObj); err != nil {
		t.Fatalf("Failed to save patient: %v", err)
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		goal := &patient.TherapeuticGoal{
			ID:          "goal-1",
			PatientID:   patientObj.ID,
			Title:       "Reduzir esquiva social",
			Description: "Trabalhar habilidades sociais",
			Status:      patient.GoalStatusInProgress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := repo.Save(ctx, goal); err != nil {
			t.Fatalf("Failed to save goal: %v", err)
		}

		found, err := repo.FindByID(ctx, "goal-1")
		if err != nil {
			t.Fatalf("Failed to find goal: %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find goal")
		}
		if found.Title != "Reduzir esquiva social" {
			t.Errorf("Expected title 'Reduzir esquiva social', got %s", found.Title)
		}
		if found.Description != "Trabalhar habilidades sociais" {
			t.Errorf("Expected description 'Trabalhar habilidades sociais', got %s", found.Description)
		}
		if found.Status != patient.GoalStatusInProgress {
			t.Errorf("Expected status %s, got %s", patient.GoalStatusInProgress, found.Status)
		}
	})

	t.Run("FindByPatientID", func(t *testing.T) {
		goals, err := repo.FindByPatientID(ctx, patientObj.ID)
		if err != nil {
			t.Fatalf("Failed to find goals by patient: %v", err)
		}
		if len(goals) != 1 {
			t.Errorf("Expected 1 goal, got %d", len(goals))
		}
	})

	t.Run("GetActiveGoals", func(t *testing.T) {
		goal2 := &patient.TherapeuticGoal{
			ID:          "goal-2",
			PatientID:   patientObj.ID,
			Title:       "Meta já alcançada",
			Description: "",
			Status:      patient.GoalStatusAchieved,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := repo.Save(ctx, goal2); err != nil {
			t.Fatalf("Failed to save goal 2: %v", err)
		}

		activeGoals, err := repo.GetActiveGoals(ctx, patientObj.ID)
		if err != nil {
			t.Fatalf("Failed to get active goals: %v", err)
		}
		if len(activeGoals) != 1 {
			t.Errorf("Expected 1 active goal, got %d", len(activeGoals))
		}
		if activeGoals[0].Status != patient.GoalStatusInProgress {
			t.Errorf("Expected active goal status %s, got %s", patient.GoalStatusInProgress, activeGoals[0].Status)
		}
	})

	t.Run("GetGoalsByStatus", func(t *testing.T) {
		achievedGoals, err := repo.GetGoalsByStatus(ctx, patientObj.ID, patient.GoalStatusAchieved)
		if err != nil {
			t.Fatalf("Failed to get achieved goals: %v", err)
		}
		if len(achievedGoals) != 1 {
			t.Errorf("Expected 1 achieved goal, got %d", len(achievedGoals))
		}
		if achievedGoals[0].Status != patient.GoalStatusAchieved {
			t.Errorf("Expected status %s, got %s", patient.GoalStatusAchieved, achievedGoals[0].Status)
		}
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		if err := repo.UpdateStatus(ctx, "goal-1", patient.GoalStatusAchieved); err != nil {
			t.Fatalf("Failed to update status: %v", err)
		}

		updated, err := repo.FindByID(ctx, "goal-1")
		if err != nil {
			t.Fatalf("Failed to find updated goal: %v", err)
		}
		if updated.Status != patient.GoalStatusAchieved {
			t.Errorf("Expected status %s, got %s", patient.GoalStatusAchieved, updated.Status)
		}
	})

	t.Run("Update", func(t *testing.T) {
		goal := &patient.TherapeuticGoal{
			ID:          "goal-3",
			PatientID:   patientObj.ID,
			Title:       "Original title",
			Description: "Original description",
			Status:      patient.GoalStatusInProgress,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := repo.Save(ctx, goal); err != nil {
			t.Fatalf("Failed to save goal 3: %v", err)
		}

		goal.Title = "Updated title"
		goal.Description = "Updated description"
		if err := repo.Update(ctx, goal); err != nil {
			t.Fatalf("Failed to update goal: %v", err)
		}

		updated, err := repo.FindByID(ctx, "goal-3")
		if err != nil {
			t.Fatalf("Failed to find updated goal: %v", err)
		}
		if updated.Title != "Updated title" {
			t.Errorf("Expected title 'Updated title', got %s", updated.Title)
		}
		if updated.Description != "Updated description" {
			t.Errorf("Expected description 'Updated description', got %s", updated.Description)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := repo.Delete(ctx, "goal-3"); err != nil {
			t.Fatalf("Failed to delete goal: %v", err)
		}

		deleted, err := repo.FindByID(ctx, "goal-3")
		if err == nil {
			t.Fatalf("Expected error when finding deleted goal, got nil. Found: %v", deleted)
		}
	})

	t.Run("FindByID not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, "non-existent-id")
		if err == nil {
			t.Error("Expected error for non-existent goal ID")
		}
	})
}

func TestGoalRepositoryMultiPatient(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	repo := NewGoalRepository(db)
	patientRepo := NewPatientRepository(db)

	ctx := context.Background()

	patient1 := &patient.Patient{
		ID:        "patient-1",
		Name:      "Patient One",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	patient2 := &patient.Patient{
		ID:        "patient-2",
		Name:      "Patient Two",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := patientRepo.Save(ctx, patient1); err != nil {
		t.Fatalf("Failed to save patient 1: %v", err)
	}
	if err := patientRepo.Save(ctx, patient2); err != nil {
		t.Fatalf("Failed to save patient 2: %v", err)
	}

	goal1 := &patient.TherapeuticGoal{
		ID:        "goal-p1-1",
		PatientID: patient1.ID,
		Title:     "Meta do Paciente 1",
		Status:    patient.GoalStatusInProgress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	goal2 := &patient.TherapeuticGoal{
		ID:        "goal-p2-1",
		PatientID: patient2.ID,
		Title:     "Meta do Paciente 2",
		Status:    patient.GoalStatusInProgress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.Save(ctx, goal1); err != nil {
		t.Fatalf("Failed to save goal 1: %v", err)
	}
	if err := repo.Save(ctx, goal2); err != nil {
		t.Fatalf("Failed to save goal 2: %v", err)
	}

	t.Run("Goals are isolated by patient", func(t *testing.T) {
		goals1, err := repo.FindByPatientID(ctx, patient1.ID)
		if err != nil {
			t.Fatalf("Failed to find goals for patient 1: %v", err)
		}
		if len(goals1) != 1 {
			t.Errorf("Expected 1 goal for patient 1, got %d", len(goals1))
		}
		if goals1[0].ID != goal1.ID {
			t.Errorf("Expected goal ID %s, got %s", goal1.ID, goals1[0].ID)
		}

		goals2, err := repo.FindByPatientID(ctx, patient2.ID)
		if err != nil {
			t.Fatalf("Failed to find goals for patient 2: %v", err)
		}
		if len(goals2) != 1 {
			t.Errorf("Expected 1 goal for patient 2, got %d", len(goals2))
		}
		if goals2[0].ID != goal2.ID {
			t.Errorf("Expected goal ID %s, got %s", goal2.ID, goals2[0].ID)
		}
	})

	t.Run("Active goals are isolated by patient", func(t *testing.T) {
		active1, err := repo.GetActiveGoals(ctx, patient1.ID)
		if err != nil {
			t.Fatalf("Failed to get active goals for patient 1: %v", err)
		}
		if len(active1) != 1 {
			t.Errorf("Expected 1 active goal for patient 1, got %d", len(active1))
		}

		active2, err := repo.GetActiveGoals(ctx, patient2.ID)
		if err != nil {
			t.Fatalf("Failed to get active goals for patient 2: %v", err)
		}
		if len(active2) != 1 {
			t.Errorf("Expected 1 active goal for patient 2, got %d", len(active2))
		}
	})
}

func TestGoalRepositoryArchivedGoals(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	repo := NewGoalRepository(db)
	patientRepo := NewPatientRepository(db)

	ctx := context.Background()

	patientObj := &patient.Patient{
		ID:        "patient-archived-test",
		Name:      "Test Patient Archived",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := patientRepo.Save(ctx, patientObj); err != nil {
		t.Fatalf("Failed to save patient: %v", err)
	}

	goals := []*patient.TherapeuticGoal{
		{
			ID:        "archived-1",
			PatientID: patientObj.ID,
			Title:     "Meta arquivada 1",
			Status:    patient.GoalStatusArchived,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "archived-2",
			PatientID: patientObj.ID,
			Title:     "Meta arquivada 2",
			Status:    patient.GoalStatusArchived,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "in-progress-1",
			PatientID: patientObj.ID,
			Title:     "Meta em progresso",
			Status:    patient.GoalStatusInProgress,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, g := range goals {
		if err := repo.Save(ctx, g); err != nil {
			t.Fatalf("Failed to save goal %s: %v", g.ID, err)
		}
	}

	t.Run("GetArchivedGoals", func(t *testing.T) {
		archived, err := repo.GetGoalsByStatus(ctx, patientObj.ID, patient.GoalStatusArchived)
		if err != nil {
			t.Fatalf("Failed to get archived goals: %v", err)
		}
		if len(archived) != 2 {
			t.Errorf("Expected 2 archived goals, got %d", len(archived))
		}
	})

	t.Run("GetInProgressGoals", func(t *testing.T) {
		inProgress, err := repo.GetGoalsByStatus(ctx, patientObj.ID, patient.GoalStatusInProgress)
		if err != nil {
			t.Fatalf("Failed to get in-progress goals: %v", err)
		}
		if len(inProgress) != 1 {
			t.Errorf("Expected 1 in-progress goal, got %d", len(inProgress))
		}
	})

	t.Run("GetAllGoals", func(t *testing.T) {
		all, err := repo.FindByPatientID(ctx, patientObj.ID)
		if err != nil {
			t.Fatalf("Failed to get all goals: %v", err)
		}
		if len(all) != 3 {
			t.Errorf("Expected 3 total goals, got %d", len(all))
		}
	})
}
