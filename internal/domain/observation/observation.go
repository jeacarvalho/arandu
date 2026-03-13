package observation

import (
	"time"
)

type Observation struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Save(observation *Observation) error
	FindByID(id string) (*Observation, error)
	FindBySessionID(sessionID string) ([]*Observation, error)
	FindAll() ([]*Observation, error)
	Update(observation *Observation) error
	Delete(id string) error
}
