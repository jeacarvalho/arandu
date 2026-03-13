package services

import (
	"time"

	"arandu/internal/domain/session"
)

type SessionService struct {
	repo session.Repository
}

func NewSessionService(repo session.Repository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) CreateSession(patientID string, date time.Time, notes string) (*session.Session, error) {
	sess := &session.Session{
		PatientID: patientID,
		Date:      date,
		Notes:     notes,
	}
	if err := s.repo.Save(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *SessionService) GetSession(id string) (*session.Session, error) {
	return s.repo.FindByID(id)
}

func (s *SessionService) ListSessions() ([]*session.Session, error) {
	return s.repo.FindAll()
}

func (s *SessionService) ListSessionsByPatient(patientID string) ([]*session.Session, error) {
	return s.repo.FindByPatientID(patientID)
}

func (s *SessionService) UpdateSession(id, patientID string, date time.Time, notes string) error {
	sess, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if sess == nil {
		return nil
	}

	sess.PatientID = patientID
	sess.Date = date
	sess.Notes = notes
	return s.repo.Update(sess)
}

func (s *SessionService) DeleteSession(id string) error {
	return s.repo.Delete(id)
}
