package patient

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Tag        string    `json:"tag"`
	Gender     string    `json:"gender"`
	Ethnicity  string    `json:"ethnicity"`
	Occupation string    `json:"occupation"`
	Education  string    `json:"education"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DashboardSummary carries per-patient data for the dashboard patient list.
// Populated by a single CTE query — no N+1.
type DashboardSummary struct {
	ID              string
	Name            string
	Tag             string
	SessionCount    int
	LastSessionDate *time.Time
	NextApptDate    string // "2006-01-02"
	NextApptTime    string // "15:04"
}

func NewPatient(name, gender, ethnicity, occupation, education, notes string) (*Patient, error) {
	if name == "" {
		return nil, errors.New("patient name cannot be empty")
	}

	patient := &Patient{
		ID:         uuid.New().String(),
		Name:       name,
		Gender:     gender,
		Ethnicity:  ethnicity,
		Occupation: occupation,
		Education:  education,
		Notes:      notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return patient, nil
}

func (p *Patient) Update(name, gender, ethnicity, occupation, education, notes string) error {
	if name == "" {
		return errors.New("patient name cannot be empty")
	}

	p.Name = name
	p.Gender = gender
	p.Ethnicity = ethnicity
	p.Occupation = occupation
	p.Education = education
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

	// Dashboard enriched list (single CTE query)
	ListForDashboard(ctx context.Context, limit int) ([]*DashboardSummary, error)
}

// Theme represents a theme with its frequency
type Theme struct {
	Term  string
	Count int
}
