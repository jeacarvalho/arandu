package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	PatientID string
	Date      time.Time
	Summary   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSession(patientID string, date time.Time, summary string) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.New().String(),
		PatientID: patientID,
		Date:      date,
		Summary:   summary,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
