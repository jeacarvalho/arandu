package services

import (
	"context"
	"fmt"
	"time"

	"arandu/internal/domain/session"
)

type CreateSessionService struct {
	repo session.Repository
}

func NewCreateSessionService(repo session.Repository) *CreateSessionService {
	return &CreateSessionService{repo: repo}
}

type CreateSessionInput struct {
	PatientID string
	Date      time.Time
	Summary   string
}

func (s *CreateSessionService) Execute(ctx context.Context, input CreateSessionInput) (*session.Session, error) {
	if input.PatientID == "" {
		return nil, fmt.Errorf("patient_id is required")
	}

	if input.Date.IsZero() {
		return nil, fmt.Errorf("date is required")
	}

	sess := session.NewSession(input.PatientID, input.Date, input.Summary)

	if err := s.repo.Create(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}
