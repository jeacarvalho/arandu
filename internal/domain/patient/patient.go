package patient

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPatient(name string, notes string) (*Patient, error) {
	if name == "" {
		return nil, errors.New("patient name cannot be empty")
	}

	patient := &Patient{
		ID:        uuid.New().String(),
		Name:      name,
		Notes:     notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return patient, nil
}

func (p *Patient) Update(name, notes string) error {
	if name == "" {
		return errors.New("patient name cannot be empty")
	}

	p.Name = name
	p.Notes = notes
	p.UpdatedAt = time.Now()

	return nil
}

type Repository interface {
	Save(ctx context.Context, patient *Patient) error
	FindByID(ctx context.Context, id string) (*Patient, error)
	FindAll(ctx context.Context) ([]*Patient, error)
	Update(ctx context.Context, patient *Patient) error
	Delete(ctx context.Context, id string) error

	// Additional queries for enhanced functionality
	FindByName(ctx context.Context, name string) ([]*Patient, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*Patient, error)
	CountAll(ctx context.Context) (int, error)
	FindPaginated(ctx context.Context, limit, offset int) ([]*Patient, error)

	// Theme analysis
	GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error)
}

// Theme represents a theme with its frequency
type Theme struct {
	Term  string
	Count int
}
