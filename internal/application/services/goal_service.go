package services

import (
	"context"
	"errors"

	"arandu/internal/domain/patient"
)

type GoalService interface {
	GetActiveGoals(ctx context.Context, patientID string) ([]*patient.TherapeuticGoal, error)
	CreateGoal(ctx context.Context, patientID, title, description string) (*patient.TherapeuticGoal, error)
	UpdateGoalStatus(ctx context.Context, goalID string, status patient.GoalStatus) error
}

type CreateGoalInput struct {
	PatientID   string `json:"patient_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

var (
	ErrGoalTitleRequired = errors.New("goal title is required")
)

func (input *CreateGoalInput) Validate() error {
	if input.PatientID == "" {
		return errors.New("patient ID is required")
	}
	if input.Title == "" {
		return ErrGoalTitleRequired
	}
	return nil
}

func (input *CreateGoalInput) Sanitize() {
	input.Title = sanitizeString(input.Title)
	input.Description = sanitizeString(input.Description)
}

func sanitizeString(s string) string {
	result := ""
	for _, r := range s {
		if r != ' ' || len(result) > 0 && result[len(result)-1] != ' ' {
			result += string(r)
		}
	}
	if len(result) > 0 && result[len(result)-1] == ' ' {
		result = result[:len(result)-1]
	}
	return result
}
