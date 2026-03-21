package services

import (
	"context"

	"arandu/internal/domain/intervention"
)

type InterventionService struct {
	repo intervention.Repository
}

func NewInterventionService(repo intervention.Repository) *InterventionService {
	return &InterventionService{repo: repo}
}

func (s *InterventionService) CreateIntervention(ctx context.Context, sessionID, content string) (*intervention.Intervention, error) {
	interv := &intervention.Intervention{
		SessionID: sessionID,
		Content:   content,
	}
	if err := s.repo.Save(ctx, interv); err != nil {
		return nil, err
	}
	return interv, nil
}

func (s *InterventionService) GetIntervention(ctx context.Context, id string) (*intervention.Intervention, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *InterventionService) ListInterventions(ctx context.Context) ([]*intervention.Intervention, error) {
	return s.repo.FindAll(ctx)
}

func (s *InterventionService) ListInterventionsBySession(ctx context.Context, sessionID string) ([]*intervention.Intervention, error) {
	return s.repo.FindBySessionID(ctx, sessionID)
}

func (s *InterventionService) UpdateIntervention(ctx context.Context, id, content string) error {
	interv, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if interv == nil {
		return nil
	}

	interv.Content = content
	return s.repo.Update(ctx, interv)
}

func (s *InterventionService) DeleteIntervention(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
