package patient

import (
	"context"
	"errors"
	"time"
)

type Anamnesis struct {
	PatientID       string    `json:"patient_id"`
	ChiefComplaint  string    `json:"chief_complaint"`
	PersonalHistory string    `json:"personal_history"`
	FamilyHistory   string    `json:"family_history"`
	MentalStateExam string    `json:"mental_state_exam"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewAnamnesis(patientID string) (*Anamnesis, error) {
	if patientID == "" {
		return nil, errors.New("patient ID cannot be empty")
	}

	if len(patientID) > 36 {
		return nil, errors.New("patient ID too long")
	}

	return &Anamnesis{
		PatientID: patientID,
		UpdatedAt: time.Now(),
	}, nil
}

func (a *Anamnesis) UpdateSection(section, content string) error {
	a.UpdatedAt = time.Now()

	switch section {
	case "chief_complaint":
		a.ChiefComplaint = content
	case "personal_history":
		a.PersonalHistory = content
	case "family_history":
		a.FamilyHistory = content
	case "mental_state_exam":
		a.MentalStateExam = content
	default:
		return errors.New("invalid section")
	}

	return nil
}

func (a *Anamnesis) Validate() error {
	if a == nil {
		return errors.New("anamnesis cannot be nil")
	}
	if a.PatientID == "" {
		return errors.New("patient ID cannot be empty")
	}

	if len(a.PatientID) > 36 {
		return errors.New("patient ID too long")
	}

	return nil
}

func (a *Anamnesis) IsEmpty() bool {
	return a.ChiefComplaint == "" &&
		a.PersonalHistory == "" &&
		a.FamilyHistory == "" &&
		a.MentalStateExam == ""
}

type AnamnesisRepository interface {
	GetAnamnesis(ctx context.Context, patientID string) (*Anamnesis, error)
	SaveAnamnesis(ctx context.Context, anamnesis *Anamnesis) error
}
