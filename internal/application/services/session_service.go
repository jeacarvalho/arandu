package services

import (
	"context"
	"errors"
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

type UpdateSessionInput struct {
	ID      string
	Date    time.Time
	Summary string
}

func (s *SessionService) UpdateSession(ctx context.Context, input UpdateSessionInput) error {
	sess, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if sess == nil {
		return errors.New("session not found")
	}

	if err := sess.Update(input.Date, input.Summary); err != nil {
		return err
	}

	return s.repo.Update(ctx, sess)
}
