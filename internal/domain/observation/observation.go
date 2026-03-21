package observation

import (
	"context"
	"time"
)

type Observation struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Save(ctx context.Context, observation *Observation) error
	FindByID(ctx context.Context, id string) (*Observation, error)
	FindBySessionID(ctx context.Context, sessionID string) ([]*Observation, error)
	FindAll(ctx context.Context) ([]*Observation, error)
	Update(ctx context.Context, observation *Observation) error
	Delete(ctx context.Context, id string) error
}
