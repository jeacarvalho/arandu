package insight

import (
	"context"
	"time"
)

type Insight struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Source    string    `json:"source"` // "ai" or "therapist"
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Save(ctx context.Context, insight *Insight) error
	FindByID(ctx context.Context, id string) (*Insight, error)
	FindAll(ctx context.Context) ([]*Insight, error)
	Delete(ctx context.Context, id string) error
}
