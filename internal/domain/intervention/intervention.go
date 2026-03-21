package intervention

import (
	"context"
	"time"
)

type Intervention struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Save(ctx context.Context, intervention *Intervention) error
	FindByID(ctx context.Context, id string) (*Intervention, error)
	FindBySessionID(ctx context.Context, sessionID string) ([]*Intervention, error)
	FindAll(ctx context.Context) ([]*Intervention, error)
	Update(ctx context.Context, intervention *Intervention) error
	Delete(ctx context.Context, id string) error
}
