package patient

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Vitals struct {
	ID               string    `json:"id"`
	PatientID        string    `json:"patient_id"`
	Date             time.Time `json:"date"`
	SleepHours       *float64  `json:"sleep_hours,omitempty"`
	AppetiteLevel    *int      `json:"appetite_level,omitempty"`
	Weight           *float64  `json:"weight,omitempty"`
	PhysicalActivity int       `json:"physical_activity"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewVitals(patientID string, date time.Time, sleepHours *float64, appetiteLevel *int, weight *float64, physicalActivity int, notes string) (*Vitals, error) {
	if patientID == "" {
		return nil, errors.New("patient ID cannot be empty")
	}
	if date.After(time.Now()) {
		return nil, errors.New("vitals date cannot be in the future")
	}
	if appetiteLevel != nil && (*appetiteLevel < 1 || *appetiteLevel > 10) {
		return nil, errors.New("appetite level must be between 1 and 10")
	}
	if sleepHours != nil && (*sleepHours < 0 || *sleepHours > 24) {
		return nil, errors.New("sleep hours must be between 0 and 24")
	}

	return &Vitals{
		ID:               uuid.New().String(),
		PatientID:        patientID,
		Date:             date,
		SleepHours:       sleepHours,
		AppetiteLevel:    appetiteLevel,
		Weight:           weight,
		PhysicalActivity: physicalActivity,
		Notes:            notes,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}, nil
}

func (v *Vitals) Update(sleepHours *float64, appetiteLevel *int, weight *float64, physicalActivity int, notes string) error {
	if appetiteLevel != nil && (*appetiteLevel < 1 || *appetiteLevel > 10) {
		return errors.New("appetite level must be between 1 and 10")
	}
	if sleepHours != nil && (*sleepHours < 0 || *sleepHours > 24) {
		return errors.New("sleep hours must be between 0 and 24")
	}

	v.SleepHours = sleepHours
	v.AppetiteLevel = appetiteLevel
	v.Weight = weight
	v.PhysicalActivity = physicalActivity
	v.Notes = notes
	v.UpdatedAt = time.Now()
	return nil
}

type VitalsAverage struct {
	AverageSleepHours       *float64
	AverageAppetiteLevel    *float64
	AverageWeight           *float64
	AveragePhysicalActivity *float64
	Count                   int
}

type VitalsRepository interface {
	Save(ctx context.Context, v *Vitals) error
	FindByID(ctx context.Context, id string) (*Vitals, error)
	FindByPatientID(ctx context.Context, patientID string, limit int) ([]*Vitals, error)
	GetLatestVitals(ctx context.Context, patientID string) (*Vitals, error)
	GetAverageVitals(ctx context.Context, patientID string, days int) (*VitalsAverage, error)
	Update(ctx context.Context, v *Vitals) error
	Delete(ctx context.Context, id string) error
	FindByPatientIDAndTimeframe(ctx context.Context, patientID string, startTime time.Time) ([]*Vitals, error)
}
