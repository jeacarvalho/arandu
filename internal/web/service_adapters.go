package web

import (
	"context"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
)

// SessionServiceAdapter adapts services.SessionService to both handlers.SessionService and handlers.SessionServiceInterface
type SessionServiceAdapter struct {
	service *services.SessionService
}

// NewSessionServiceAdapter creates a new adapter
func NewSessionServiceAdapter(service *services.SessionService) *SessionServiceAdapter {
	return &SessionServiceAdapter{service: service}
}

// ListSessionsByPatient implements handlers.SessionService interface
func (a *SessionServiceAdapter) ListSessionsByPatient(ctx context.Context, patientID string) ([]*session.Session, error) {
	return a.service.ListSessionsByPatient(ctx, patientID)
}

// CreateSession implements handlers.SessionService interface (for patient_handler.go)
func (a *SessionServiceAdapter) CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error) {
	return a.service.CreateSession(ctx, patientID, date, summary)
}

// GetSession implements handlers.SessionService interface
func (a *SessionServiceAdapter) GetSession(ctx context.Context, id string) (*session.Session, error) {
	return a.service.GetSession(ctx, id)
}

// UpdateSession implements handlers.SessionService interface
func (a *SessionServiceAdapter) UpdateSession(ctx context.Context, input services.UpdateSessionInput) error {
	return a.service.UpdateSession(ctx, input)
}

// InsightServiceAdapter adapts services.InsightService to handlers.InsightService interface
type InsightServiceAdapter struct {
	service *services.InsightService
}

// NewInsightServiceAdapter creates a new adapter
func NewInsightServiceAdapter(service *services.InsightService) *InsightServiceAdapter {
	return &InsightServiceAdapter{service: service}
}

// GetInsightsByPatient implements handlers.InsightService interface
func (a *InsightServiceAdapter) GetInsightsByPatient(ctx context.Context, patientID string, limit int) ([]interface{}, error) {
	// For now, return empty list - TODO: implement proper filtering by patient
	return []interface{}{}, nil
}

// PatientServiceAdapter adapts services.PatientService to handlers.PatientServiceInterface
type PatientServiceAdapter struct {
	service *services.PatientService
}

// NewPatientServiceAdapter creates a new adapter
func NewPatientServiceAdapter(service *services.PatientService) *PatientServiceAdapter {
	return &PatientServiceAdapter{service: service}
}

// GetPatientByID implements handlers.PatientServiceInterface
func (a *PatientServiceAdapter) GetPatientByID(ctx context.Context, id string) (*patient.Patient, error) {
	return a.service.GetPatientByID(ctx, id)
}

// ListPatients implements handlers.PatientServiceInterface
func (a *PatientServiceAdapter) ListPatients(ctx context.Context) ([]*patient.Patient, error) {
	return a.service.ListPatients(ctx)
}

// CreatePatient implements handlers.PatientServiceInterface
func (a *PatientServiceAdapter) CreatePatient(ctx context.Context, input services.CreatePatientInput) (*patient.Patient, error) {
	return a.service.CreatePatient(ctx, input)
}

// ObservationServiceAdapter adapts services.ObservationService to handlers.ObservationServiceInterface
type ObservationServiceAdapter struct {
	service *services.ObservationService
}

// NewObservationServiceAdapter creates a new adapter
func NewObservationServiceAdapter(service *services.ObservationService) *ObservationServiceAdapter {
	return &ObservationServiceAdapter{service: service}
}

// CreateObservation implements handlers.ObservationServiceInterface
func (a *ObservationServiceAdapter) CreateObservation(ctx context.Context, sessionID, content string) (interface{}, error) {
	return a.service.CreateObservation(sessionID, content)
}

// GetObservationsBySession implements handlers.ObservationServiceInterface
func (a *ObservationServiceAdapter) GetObservationsBySession(ctx context.Context, sessionID string) ([]interface{}, error) {
	observations, err := a.service.ListObservationsBySession(sessionID)
	if err != nil {
		return nil, err
	}

	// Convert to []interface{}
	result := make([]interface{}, len(observations))
	for i, obs := range observations {
		result[i] = obs
	}
	return result, nil
}
