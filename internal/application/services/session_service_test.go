package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"arandu/internal/domain/session"
)

// mockSessionRepository is a mock implementation of session.Repository for testing
type mockSessionRepository struct {
	createFunc        func(ctx context.Context, s *session.Session) error
	getByIDFunc       func(ctx context.Context, id string) (*session.Session, error)
	listByPatientFunc func(ctx context.Context, patientID string) ([]*session.Session, error)
	updateFunc        func(ctx context.Context, s *session.Session) error
	deleteFunc        func(ctx context.Context, id string) error
}

func (m *mockSessionRepository) Create(ctx context.Context, s *session.Session) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, s)
	}
	return nil
}

func (m *mockSessionRepository) GetByID(ctx context.Context, id string) (*session.Session, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockSessionRepository) ListByPatient(ctx context.Context, patientID string) ([]*session.Session, error) {
	if m.listByPatientFunc != nil {
		return m.listByPatientFunc(ctx, patientID)
	}
	return nil, nil
}

func (m *mockSessionRepository) Update(ctx context.Context, s *session.Session) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, s)
	}
	return nil
}

func (m *mockSessionRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestNewSessionService(t *testing.T) {
	repo := &mockSessionRepository{}
	service := NewSessionService(repo)

	if service == nil {
		t.Error("Expected SessionService to be created")
	}
}

func TestSessionService_CreateSession_Success(t *testing.T) {
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

	service := NewSessionService(repo)
	session, err := service.CreateSession(ctx, patientID, date, summary)

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

func TestSessionService_CreateSession_RepositoryError(t *testing.T) {
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

	service := NewSessionService(repo)
	session, err := service.CreateSession(ctx, patientID, date, summary)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	if session != nil {
		t.Error("Expected no session to be returned on error")
	}
}

func TestSessionService_GetSession_Success(t *testing.T) {
	ctx := context.Background()
	sessionID := "session-123"
	expectedSession := &session.Session{
		ID:        sessionID,
		PatientID: "patient-123",
		Date:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Summary:   "Test session",
	}

	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			if id == sessionID {
				return expectedSession, nil
			}
			return nil, nil
		},
	}

	service := NewSessionService(repo)
	session, err := service.GetSession(ctx, sessionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if session != expectedSession {
		t.Error("Expected session to be returned from repository")
	}
}

func TestSessionService_GetSession_RepositoryError(t *testing.T) {
	ctx := context.Background()
	sessionID := "session-123"
	expectedErr := errors.New("repository error")

	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			return nil, expectedErr
		},
	}

	service := NewSessionService(repo)
	session, err := service.GetSession(ctx, sessionID)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}

	if session != nil {
		t.Error("Expected no session to be returned on error")
	}
}

func TestSessionService_ListSessionsByPatient_Success(t *testing.T) {
	ctx := context.Background()
	patientID := "patient-123"
	expectedSessions := []*session.Session{
		{
			ID:        "session-1",
			PatientID: patientID,
			Date:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Summary:   "Session 1",
		},
		{
			ID:        "session-2",
			PatientID: patientID,
			Date:      time.Date(2024, 1, 22, 14, 0, 0, 0, time.UTC),
			Summary:   "Session 2",
		},
	}

	repo := &mockSessionRepository{
		listByPatientFunc: func(ctx context.Context, pid string) ([]*session.Session, error) {
			if pid == patientID {
				return expectedSessions, nil
			}
			return nil, nil
		},
	}

	service := NewSessionService(repo)
	sessions, err := service.ListSessionsByPatient(ctx, patientID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(sessions) != len(expectedSessions) {
		t.Errorf("Expected %d sessions, got %d", len(expectedSessions), len(sessions))
	}

	for i, sess := range sessions {
		if sess != expectedSessions[i] {
			t.Errorf("Session %d doesn't match expected", i)
		}
	}
}

func TestSessionService_UpdateSession_Success(t *testing.T) {
	ctx := context.Background()
	sessionID := "session-123"
	originalDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	newDate := time.Date(2024, 1, 22, 14, 0, 0, 0, time.UTC)
	newSummary := "Updated session summary"

	existingSession := &session.Session{
		ID:        sessionID,
		PatientID: "patient-123",
		Date:      originalDate,
		Summary:   "Original summary",
		CreatedAt: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
	}

	var updatedSession *session.Session
	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			if id == sessionID {
				return existingSession, nil
			}
			return nil, nil
		},
		updateFunc: func(ctx context.Context, s *session.Session) error {
			updatedSession = s
			return nil
		},
	}

	service := NewSessionService(repo)
	input := UpdateSessionInput{
		ID:      sessionID,
		Date:    newDate,
		Summary: newSummary,
	}

	err := service.UpdateSession(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if updatedSession != existingSession {
		t.Error("Expected existing session to be updated")
	}

	if !updatedSession.Date.Equal(newDate) {
		t.Errorf("Expected Date %v, got %v", newDate, updatedSession.Date)
	}

	if updatedSession.Summary != newSummary {
		t.Errorf("Expected Summary %s, got %s", newSummary, updatedSession.Summary)
	}

	if !updatedSession.UpdatedAt.After(existingSession.CreatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

func TestSessionService_UpdateSession_SessionNotFound(t *testing.T) {
	ctx := context.Background()
	sessionID := "non-existent-session"

	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			return nil, nil
		},
	}

	service := NewSessionService(repo)
	input := UpdateSessionInput{
		ID:      sessionID,
		Date:    time.Date(2024, 1, 22, 14, 0, 0, 0, time.UTC),
		Summary: "Updated summary",
	}

	err := service.UpdateSession(ctx, input)

	if err == nil {
		t.Error("Expected error when session not found")
	}

	if err.Error() != "session not found" {
		t.Errorf("Expected 'session not found' error, got %v", err)
	}
}

func TestSessionService_UpdateSession_GetSessionError(t *testing.T) {
	ctx := context.Background()
	sessionID := "session-123"
	expectedErr := errors.New("repository error")

	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			return nil, expectedErr
		},
	}

	service := NewSessionService(repo)
	input := UpdateSessionInput{
		ID:      sessionID,
		Date:    time.Date(2024, 1, 22, 14, 0, 0, 0, time.UTC),
		Summary: "Updated summary",
	}

	err := service.UpdateSession(ctx, input)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestSessionService_UpdateSession_UpdateError(t *testing.T) {
	ctx := context.Background()
	sessionID := "session-123"
	expectedErr := errors.New("update error")

	existingSession := &session.Session{
		ID:        sessionID,
		PatientID: "patient-123",
		Date:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Summary:   "Original summary",
	}

	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			if id == sessionID {
				return existingSession, nil
			}
			return nil, nil
		},
		updateFunc: func(ctx context.Context, s *session.Session) error {
			return expectedErr
		},
	}

	service := NewSessionService(repo)
	input := UpdateSessionInput{
		ID:      sessionID,
		Date:    time.Date(2024, 1, 22, 14, 0, 0, 0, time.UTC),
		Summary: "Updated summary",
	}

	err := service.UpdateSession(ctx, input)

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestSessionService_UpdateSession_InvalidUpdate(t *testing.T) {
	ctx := context.Background()
	sessionID := "session-123"

	existingSession := &session.Session{
		ID:        sessionID,
		PatientID: "patient-123",
		Date:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Summary:   "Original summary",
	}

	repo := &mockSessionRepository{
		getByIDFunc: func(ctx context.Context, id string) (*session.Session, error) {
			if id == sessionID {
				return existingSession, nil
			}
			return nil, nil
		},
	}

	service := NewSessionService(repo)
	// Invalid date (future date)
	input := UpdateSessionInput{
		ID:      sessionID,
		Date:    time.Now().Add(24 * time.Hour),
		Summary: "Updated summary",
	}

	err := service.UpdateSession(ctx, input)

	if err == nil {
		t.Error("Expected error for invalid update")
	}
}
