package session

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patient_id"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Save(session *Session) error
	FindByID(id string) (*Session, error)
	FindByPatientID(patientID string) ([]*Session, error)
	FindAll() ([]*Session, error)
	Update(session *Session) error
	Delete(id string) error
}
