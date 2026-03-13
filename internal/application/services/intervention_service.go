package services

import (
	"arandu/internal/domain/intervention"
)

type InterventionService struct {
	repo intervention.Repository
}

func NewInterventionService(repo intervention.Repository) *InterventionService {
	return &InterventionService{repo: repo}
}

func (s *InterventionService) CreateIntervention(sessionID, content string) (*intervention.Intervention, error) {
	interv := &intervention.Intervention{
		SessionID: sessionID,
		Content:   content,
	}
	if err := s.repo.Save(interv); err != nil {
		return nil, err
	}
	return interv, nil
}

func (s *InterventionService) GetIntervention(id string) (*intervention.Intervention, error) {
	return s.repo.FindByID(id)
}

func (s *InterventionService) ListInterventions() ([]*intervention.Intervention, error) {
	return s.repo.FindAll()
}

func (s *InterventionService) ListInterventionsBySession(sessionID string) ([]*intervention.Intervention, error) {
	return s.repo.FindBySessionID(sessionID)
}

func (s *InterventionService) UpdateIntervention(id, content string) error {
	interv, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if interv == nil {
		return nil
	}

	interv.Content = content
	return s.repo.Update(interv)
}

func (s *InterventionService) DeleteIntervention(id string) error {
	return s.repo.Delete(id)
}
