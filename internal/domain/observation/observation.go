package observation

import (
	"context"
	"time"
)

type Observation struct {
	ID        string           `json:"id"`
	SessionID string           `json:"session_id"`
	Content   string           `json:"content"`
	Tags      []ObservationTag `json:"tags,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type Repository interface {
	Save(ctx context.Context, observation *Observation) error
	FindByID(ctx context.Context, id string) (*Observation, error)
	FindBySessionID(ctx context.Context, sessionID string) ([]*Observation, error)
	FindAll(ctx context.Context) ([]*Observation, error)
	Update(ctx context.Context, observation *Observation) error
	Delete(ctx context.Context, id string) error

	// Tag-related methods
	GetTags(ctx context.Context) ([]Tag, error)
	GetTagsByType(ctx context.Context, tagType TagType) ([]Tag, error)
	AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error
	RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error
	GetObservationTags(ctx context.Context, observationID string) ([]ObservationTag, error)
	GetTagsSummary(ctx context.Context) ([]TagSummary, error)
	GetTagsSummaryByPatient(ctx context.Context, patientID string) ([]TagSummary, error)
	FindByTag(ctx context.Context, tagID string) ([]*Observation, error)
}
