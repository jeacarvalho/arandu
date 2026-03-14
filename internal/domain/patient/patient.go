package patient

import (
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
	Save(patient *Patient) error
	FindByID(id string) (*Patient, error)
	FindAll() ([]*Patient, error)
	Update(patient *Patient) error
	Delete(id string) error

	// Additional queries for enhanced functionality
	FindByName(name string) ([]*Patient, error)
	CountAll() (int, error)
	FindPaginated(limit, offset int) ([]*Patient, error)
}
