package patient

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type GoalStatus string

const (
	GoalStatusInProgress GoalStatus = "in_progress"
	GoalStatusAchieved   GoalStatus = "achieved"
	GoalStatusArchived   GoalStatus = "archived"
)

func (s GoalStatus) IsValid() bool {
	switch s {
	case GoalStatusInProgress, GoalStatusAchieved, GoalStatusArchived:
		return true
	}
	return false
}

type TherapeuticGoal struct {
	ID          string     `json:"id"`
	PatientID   string     `json:"patient_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      GoalStatus `json:"status"`
	ClosureNote string     `json:"closure_note"`
	ClosedAt    *time.Time `json:"closed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTherapeuticGoal(patientID, title, description string) (*TherapeuticGoal, error) {
	if patientID == "" {
		return nil, errors.New("patient ID cannot be empty")
	}
	if title == "" {
		return nil, errors.New("goal title cannot be empty")
	}

	return &TherapeuticGoal{
		ID:          uuid.New().String(),
		PatientID:   patientID,
		Title:       title,
		Description: description,
		Status:      GoalStatusInProgress,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (g *TherapeuticGoal) SetStatus(status GoalStatus) error {
	if !status.IsValid() {
		return errors.New("invalid goal status")
	}
	g.Status = status
	g.UpdatedAt = time.Now()
	return nil
}

func (g *TherapeuticGoal) MarkAchieved() error {
	return g.SetStatus(GoalStatusAchieved)
}

func (g *TherapeuticGoal) CloseWithNote(closureNote string) error {
	now := time.Now()
	g.ClosureNote = closureNote
	g.ClosedAt = &now
	return g.SetStatus(GoalStatusAchieved)
}

func (g *TherapeuticGoal) Archive() error {
	return g.SetStatus(GoalStatusArchived)
}

func (g *TherapeuticGoal) Reopen() error {
	return g.SetStatus(GoalStatusInProgress)
}

func (g *TherapeuticGoal) Update(title, description string) error {
	if title == "" {
		return errors.New("goal title cannot be empty")
	}
	g.Title = title
	g.Description = description
	g.UpdatedAt = time.Now()
	return nil
}

func (g *TherapeuticGoal) IsAchieved() bool {
	return g.Status == GoalStatusAchieved
}

func (g *TherapeuticGoal) IsArchived() bool {
	return g.Status == GoalStatusArchived
}

func (g *TherapeuticGoal) IsInProgress() bool {
	return g.Status == GoalStatusInProgress
}

type TherapeuticGoalRepository interface {
	Save(ctx context.Context, goal *TherapeuticGoal) error
	FindByID(ctx context.Context, id string) (*TherapeuticGoal, error)
	FindByPatientID(ctx context.Context, patientID string) ([]*TherapeuticGoal, error)
	GetActiveGoals(ctx context.Context, patientID string) ([]*TherapeuticGoal, error)
	GetGoalsByStatus(ctx context.Context, patientID string, status GoalStatus) ([]*TherapeuticGoal, error)
	UpdateStatus(ctx context.Context, id string, status GoalStatus) error
	Update(ctx context.Context, goal *TherapeuticGoal) error
	Delete(ctx context.Context, id string) error
}
