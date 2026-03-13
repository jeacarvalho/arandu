package intervention

import (
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
	Save(intervention *Intervention) error
	FindByID(id string) (*Intervention, error)
	FindBySessionID(sessionID string) ([]*Intervention, error)
	FindAll() ([]*Intervention, error)
	Update(intervention *Intervention) error
	Delete(id string) error
}
