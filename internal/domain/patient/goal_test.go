package patient

import (
	"testing"
	"time"
)

func TestNewTherapeuticGoal(t *testing.T) {
	patientID := "patient-123"
	title := "Reduzir esquiva social"
	description := "Trabalhar habilidades sociais em contextos grupais"

	goal, err := NewTherapeuticGoal(patientID, title, description)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.ID == "" {
		t.Error("Expected goal ID to be set")
	}
	if goal.PatientID != patientID {
		t.Errorf("Expected patient ID %s, got %s", patientID, goal.PatientID)
	}
	if goal.Title != title {
		t.Errorf("Expected title %s, got %s", title, goal.Title)
	}
	if goal.Description != description {
		t.Errorf("Expected description %s, got %s", description, goal.Description)
	}
	if goal.Status != GoalStatusInProgress {
		t.Errorf("Expected status %s, got %s", GoalStatusInProgress, goal.Status)
	}
	if goal.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}
	if goal.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}
}

func TestNewTherapeuticGoalValidation(t *testing.T) {
	tests := []struct {
		name        string
		patientID   string
		title       string
		description string
		expectErr   bool
		errMsg      string
	}{
		{
			name:        "empty patient ID",
			patientID:   "",
			title:       "Reduzir esquiva",
			description: "Descrição",
			expectErr:   true,
			errMsg:      "patient ID cannot be empty",
		},
		{
			name:        "empty title",
			patientID:   "patient-123",
			title:       "",
			description: "Descrição",
			expectErr:   true,
			errMsg:      "goal title cannot be empty",
		},
		{
			name:        "valid goal without description",
			patientID:   "patient-123",
			title:       "Reduzir esquiva",
			description: "",
			expectErr:   false,
		},
		{
			name:        "valid goal with all fields",
			patientID:   "patient-123",
			title:       "Melhorar sono",
			description: "Implementar higiene do sono",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTherapeuticGoal(tt.patientID, tt.title, tt.description)
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

func TestGoalStatusIsValid(t *testing.T) {
	tests := []struct {
		status  GoalStatus
		isValid bool
	}{
		{GoalStatusInProgress, true},
		{GoalStatusAchieved, true},
		{GoalStatusArchived, true},
		{GoalStatus("invalid"), false},
		{GoalStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if tt.status.IsValid() != tt.isValid {
				t.Errorf("Expected IsValid() = %v for status %s", tt.isValid, tt.status)
			}
		})
	}
}

func TestGoalSetStatus(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")

	tests := []struct {
		name      string
		newStatus GoalStatus
		expectErr bool
	}{
		{"set in_progress", GoalStatusInProgress, false},
		{"set achieved", GoalStatusAchieved, false},
		{"set archived", GoalStatusArchived, false},
		{"set invalid", GoalStatus("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := goal.SetStatus(tt.newStatus)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if goal.Status != tt.newStatus {
					t.Errorf("Expected status %s, got %s", tt.newStatus, goal.Status)
				}
			}
		})
	}
}

func TestGoalMarkAchieved(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")

	if err := goal.MarkAchieved(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.Status != GoalStatusAchieved {
		t.Errorf("Expected status %s, got %s", GoalStatusAchieved, goal.Status)
	}
}

func TestGoalArchive(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")

	if err := goal.Archive(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.Status != GoalStatusArchived {
		t.Errorf("Expected status %s, got %s", GoalStatusArchived, goal.Status)
	}
}

func TestGoalReopen(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")
	goal.MarkAchieved()

	if err := goal.Reopen(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.Status != GoalStatusInProgress {
		t.Errorf("Expected status %s, got %s", GoalStatusInProgress, goal.Status)
	}
}

func TestGoalUpdate(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Original title", "Original description")

	newTitle := "Updated title"
	newDescription := "Updated description"

	if err := goal.Update(newTitle, newDescription); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if goal.Title != newTitle {
		t.Errorf("Expected title %s, got %s", newTitle, goal.Title)
	}
	if goal.Description != newDescription {
		t.Errorf("Expected description %s, got %s", newDescription, goal.Description)
	}
}

func TestGoalUpdateValidation(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Original title", "")

	err := goal.Update("", "New description")
	if err == nil {
		t.Error("Expected error for empty title")
	}
	if err.Error() != "goal title cannot be empty" {
		t.Errorf("Expected error message 'goal title cannot be empty', got %s", err.Error())
	}
}

func TestGoalIsAchieved(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")

	if goal.IsAchieved() {
		t.Error("Expected goal to not be achieved initially")
	}

	goal.MarkAchieved()
	if !goal.IsAchieved() {
		t.Error("Expected goal to be achieved after MarkAchieved")
	}
}

func TestGoalIsArchived(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")

	if goal.IsArchived() {
		t.Error("Expected goal to not be archived initially")
	}

	goal.Archive()
	if !goal.IsArchived() {
		t.Error("Expected goal to be archived after Archive")
	}
}

func TestGoalIsInProgress(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")

	if !goal.IsInProgress() {
		t.Error("Expected goal to be in progress initially")
	}

	goal.MarkAchieved()
	if goal.IsInProgress() {
		t.Error("Expected goal to not be in progress after MarkAchieved")
	}
}

func TestGoalUpdatedAtChanges(t *testing.T) {
	goal, _ := NewTherapeuticGoal("patient-123", "Test goal", "")
	originalUpdatedAt := goal.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	goal.MarkAchieved()
	if !goal.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to change after status update")
	}

	originalUpdatedAt = goal.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	goal.Update("New title", "New description")
	if !goal.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to change after update")
	}
}
