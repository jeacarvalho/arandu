package services

import (
	"arandu/internal/domain/observation"
)

type ObservationService struct {
	repo observation.Repository
}

func NewObservationService(repo observation.Repository) *ObservationService {
	return &ObservationService{repo: repo}
}

func (s *ObservationService) CreateObservation(sessionID, content string) (*observation.Observation, error) {
	obs := &observation.Observation{
		SessionID: sessionID,
		Content:   content,
	}
	if err := s.repo.Save(obs); err != nil {
		return nil, err
	}
	return obs, nil
}

func (s *ObservationService) GetObservation(id string) (*observation.Observation, error) {
	return s.repo.FindByID(id)
}

func (s *ObservationService) ListObservations() ([]*observation.Observation, error) {
	return s.repo.FindAll()
}

func (s *ObservationService) ListObservationsBySession(sessionID string) ([]*observation.Observation, error) {
	return s.repo.FindBySessionID(sessionID)
}

func (s *ObservationService) UpdateObservation(id, content string) error {
	obs, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if obs == nil {
		return nil
	}

	obs.Content = content
	return s.repo.Update(obs)
}

func (s *ObservationService) DeleteObservation(id string) error {
	return s.repo.Delete(id)
}
