package services

import (
	"context"

	"arandu/internal/domain/observation"
	"fmt"
)

type ObservationService struct {
	repo observation.Repository
}

func NewObservationService(repo observation.Repository) *ObservationService {
	return &ObservationService{repo: repo}
}

func (s *ObservationService) CreateObservation(ctx context.Context, sessionID, content string) (*observation.Observation, error) {
	if content == "" {
		return nil, fmt.Errorf("observation content cannot be empty")
	}

	if len(content) > 5000 {
		return nil, fmt.Errorf("observation content cannot exceed 5000 characters")
	}

	obs := &observation.Observation{
		SessionID: sessionID,
		Content:   content,
	}
	if err := s.repo.Save(ctx, obs); err != nil {
		return nil, err
	}
	return obs, nil
}

func (s *ObservationService) GetObservation(ctx context.Context, id string) (*observation.Observation, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ObservationService) ListObservations(ctx context.Context) ([]*observation.Observation, error) {
	return s.repo.FindAll(ctx)
}

func (s *ObservationService) ListObservationsBySession(ctx context.Context, sessionID string) ([]*observation.Observation, error) {
	return s.repo.FindBySessionID(ctx, sessionID)
}

func (s *ObservationService) UpdateObservation(ctx context.Context, id, content string) error {
	if content == "" {
		return fmt.Errorf("observation content cannot be empty")
	}

	if len(content) > 5000 {
		return fmt.Errorf("observation content cannot exceed 5000 characters")
	}

	obs, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if obs == nil {
		return fmt.Errorf("observation not found")
	}

	obs.Content = content
	return s.repo.Update(ctx, obs)
}

func (s *ObservationService) DeleteObservation(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
