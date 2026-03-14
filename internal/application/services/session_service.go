package services

import (
	"context"
	"time"

	"arandu/internal/domain/session"
)

type SessionService struct {
	repo session.Repository
}

func NewSessionService(repo session.Repository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error) {
	sess := session.NewSession(patientID, date, summary)
	if err := s.repo.Create(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *SessionService) GetSession(ctx context.Context, id string) (*session.Session, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SessionService) ListSessionsByPatient(ctx context.Context, patientID string) ([]*session.Session, error) {
	return s.repo.ListByPatient(ctx, patientID)
}
