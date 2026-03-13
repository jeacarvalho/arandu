package insight

import (
	"time"
)

type Insight struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Source    string    `json:"source"` // "ai" or "therapist"
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Save(insight *Insight) error
	FindByID(id string) (*Insight, error)
	FindAll() ([]*Insight, error)
	Delete(id string) error
}
