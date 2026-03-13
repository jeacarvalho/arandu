package patient

import (
	"time"
)

type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Repository interface {
	Save(patient *Patient) error
	FindByID(id string) (*Patient, error)
	FindAll() ([]*Patient, error)
	Update(patient *Patient) error
	Delete(id string) error
}
