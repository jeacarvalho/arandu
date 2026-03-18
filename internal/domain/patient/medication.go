package patient

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type MedicationStatus string

const (
	MedicationStatusActive    MedicationStatus = "active"
	MedicationStatusSuspended MedicationStatus = "suspended"
	MedicationStatusFinished  MedicationStatus = "finished"
)

type Medication struct {
	ID         string           `json:"id"`
	PatientID  string           `json:"patient_id"`
	Name       string           `json:"name"`
	Dosage     string           `json:"dosage"`
	Frequency  string           `json:"frequency"`
	Prescriber string           `json:"prescriber"`
	Status     MedicationStatus `json:"status"`
	StartedAt  time.Time        `json:"started_at"`
	EndedAt    *time.Time       `json:"ended_at,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

func NewMedication(patientID, name, dosage, frequency, prescriber string, startedAt time.Time) (*Medication, error) {
	if patientID == "" {
		return nil, errors.New("patient ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("medication name cannot be empty")
	}
	if startedAt.After(time.Now()) {
		return nil, errors.New("medication start date cannot be in the future")
	}

	return &Medication{
		ID:         uuid.New().String(),
		PatientID:  patientID,
		Name:       name,
		Dosage:     dosage,
		Frequency:  frequency,
		Prescriber: prescriber,
		Status:     MedicationStatusActive,
		StartedAt:  startedAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m *Medication) Suspend() error {
	if m.Status == MedicationStatusFinished {
		return errors.New("cannot suspend a finished medication")
	}
	now := time.Now()
	m.Status = MedicationStatusSuspended
	m.EndedAt = &now
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Medication) Activate() error {
	if m.Status == MedicationStatusFinished {
		return errors.New("cannot activate a finished medication")
	}
	m.Status = MedicationStatusActive
	m.EndedAt = nil
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Medication) Finish() error {
	now := time.Now()
	m.Status = MedicationStatusFinished
	m.EndedAt = &now
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Medication) IsActive() bool {
	return m.Status == MedicationStatusActive
}
