package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"arandu/internal/domain/session"
)

func TestNewCreateSessionService(t *testing.T) {
	repo := &mockSessionRepository{}
	service := NewCreateSessionService(repo)

	if service == nil {
		t.Error("Expected CreateSessionService to be created")
	}
}

func TestCreateSessionService_Execute_Success(t *testing.T) {
	ctx := context.Background()
	patientID := "patient-123"
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment session"

	var createdSession *session.Session
	repo := &mockSessionRepository{
		createFunc: func(ctx context.Context, s *session.Session) error {
			createdSession = s
			return nil
		},
	}

	service := NewCreateSessionService(repo)
	input := CreateSessionInput{
		PatientID: patientID,
		Date:      date,
		Summary:   summary,
	}

	session, err := service.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if session == nil {
		t.Error("Expected session to be created")
	}

	if session.PatientID != patientID {
		t.Errorf("Expected PatientID %s, got %s", patientID, session.PatientID)
	}

	if !session.Date.Equal(date) {
		t.Errorf("Expected Date %v, got %v", date, session.Date)
	}

	if session.Summary != summary {
		t.Errorf("Expected Summary %s, got %s", summary, session.Summary)
	}

	if createdSession != session {
		t.Error("Expected session to be passed to repository")
	}
}

func TestCreateSessionService_Execute_EmptyPatientID(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment session"

	repo := &mockSessionRepository{}
	service := NewCreateSessionService(repo)
	input := CreateSessionInput{
		PatientID: "", // Empty patient ID
		Date:      date,
		Summary:   summary,
	}

	session, err := service.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error for empty patient_id")
	}

	if err.Error() != "patient_id is required" {
		t.Errorf("Expected 'patient_id is required' error, got %v", err)
	}

	if session != nil {
		t.Error("Expected no session to be returned on error")
	}
}

func TestCreateSessionService_Execute_EmptyDate(t *testing.T) {
	ctx := context.Background()
	patientID := "patient-123"
	summary := "Initial assessment session"

	repo := &mockSessionRepository{}
	service := NewCreateSessionService(repo)
	input := CreateSessionInput{
		PatientID: patientID,
		Date:      time.Time{}, // Zero time
		Summary:   summary,
	}

	session, err := service.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error for empty date")
	}

	if err.Error() != "date is required" {
		t.Errorf("Expected 'date is required' error, got %v", err)
	}

	if session != nil {
		t.Error("Expected no session to be returned on error")
	}
}

func TestCreateSessionService_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	patientID := "patient-123"
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment session"

	expectedErr := errors.New("repository error")
	repo := &mockSessionRepository{
		createFunc: func(ctx context.Context, s *session.Session) error {
			return expectedErr
		},
	}

	service := NewCreateSessionService(repo)
	input := CreateSessionInput{
		PatientID: patientID,
		Date:      date,
		Summary:   summary,
	}

	session, err := service.Execute(ctx, input)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	if session != nil {
		t.Error("Expected no session to be returned on error")
	}
}

func TestCreateSessionService_Execute_WithEmptySummary(t *testing.T) {
	ctx := context.Background()
	patientID := "patient-123"
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	var createdSession *session.Session
	repo := &mockSessionRepository{
		createFunc: func(ctx context.Context, s *session.Session) error {
			createdSession = s
			return nil
		},
	}

	service := NewCreateSessionService(repo)
	input := CreateSessionInput{
		PatientID: patientID,
		Date:      date,
		Summary:   "", // Empty summary
	}

	session, err := service.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error for empty summary, got %v", err)
	}

	if session == nil {
		t.Error("Expected session to be created even with empty summary")
	}

	if session.Summary != "" {
		t.Errorf("Expected empty summary, got %s", session.Summary)
	}

	if createdSession != session {
		t.Error("Expected session to be passed to repository")
	}
}
