package services

import (
	"arandu/internal/domain/patient"
)

type PatientService struct {
	repo patient.Repository
}

func NewPatientService(repo patient.Repository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(name, notes string) (*patient.Patient, error) {
	p := &patient.Patient{
		Name:  name,
		Notes: notes,
	}
	if err := s.repo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PatientService) GetPatient(id string) (*patient.Patient, error) {
	return s.repo.FindByID(id)
}

func (s *PatientService) ListPatients() ([]*patient.Patient, error) {
	return s.repo.FindAll()
}

func (s *PatientService) UpdatePatient(id, name, notes string) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if p == nil {
		return nil
	}

	p.Name = name
	p.Notes = notes
	return s.repo.Update(p)
}

func (s *PatientService) DeletePatient(id string) error {
	return s.repo.Delete(id)
}
