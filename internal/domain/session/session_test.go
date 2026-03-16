package session

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewSession(t *testing.T) {
	patientID := uuid.New().String()
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment session"

	session := NewSession(patientID, date, summary)

	if session.ID == "" {
		t.Error("Expected session ID to be generated")
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

	if session.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if session.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	if !session.CreatedAt.Equal(session.UpdatedAt) {
		t.Error("Expected CreatedAt and UpdatedAt to be equal for new session")
	}
}

func TestUpdate_ValidInput(t *testing.T) {
	patientID := uuid.New().String()
	originalDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	originalSummary := "Initial assessment"
	session := NewSession(patientID, originalDate, originalSummary)

	newDate := time.Date(2024, 1, 22, 14, 0, 0, 0, time.UTC)
	newSummary := "Follow-up session"

	err := session.Update(newDate, newSummary)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !session.Date.Equal(newDate) {
		t.Errorf("Expected Date %v, got %v", newDate, session.Date)
	}

	if session.Summary != newSummary {
		t.Errorf("Expected Summary %s, got %s", newSummary, session.Summary)
	}

	if !session.UpdatedAt.After(session.CreatedAt) {
		t.Error("Expected UpdatedAt to be after CreatedAt after update")
	}
}

func TestUpdate_InvalidDate_Zero(t *testing.T) {
	patientID := uuid.New().String()
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment"
	session := NewSession(patientID, date, summary)

	err := session.Update(time.Time{}, "Updated summary")
	if err != ErrInvalidDate {
		t.Errorf("Expected ErrInvalidDate, got %v", err)
	}
}

func TestUpdate_InvalidDate_Future(t *testing.T) {
	patientID := uuid.New().String()
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment"
	session := NewSession(patientID, date, summary)

	futureDate := time.Now().Add(24 * time.Hour)
	err := session.Update(futureDate, "Updated summary")
	if err != ErrInvalidDate {
		t.Errorf("Expected ErrInvalidDate, got %v", err)
	}
}

func TestUpdate_SummaryTooLong(t *testing.T) {
	patientID := uuid.New().String()
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment"
	session := NewSession(patientID, date, summary)

	// Create a summary longer than 10000 characters
	longSummary := ""
	for i := 0; i < 10001; i++ {
		longSummary += "a"
	}

	err := session.Update(date, longSummary)
	if err != ErrSummaryTooLong {
		t.Errorf("Expected ErrSummaryTooLong, got %v", err)
	}
}

func TestUpdate_SummaryAtLimit(t *testing.T) {
	patientID := uuid.New().String()
	date := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	summary := "Initial assessment"
	session := NewSession(patientID, date, summary)

	// Create a summary exactly at the limit (10000 characters)
	limitSummary := ""
	for i := 0; i < 10000; i++ {
		limitSummary += "a"
	}

	err := session.Update(date, limitSummary)
	if err != nil {
		t.Errorf("Expected no error for summary at limit, got %v", err)
	}

	if session.Summary != limitSummary {
		t.Error("Expected summary to be updated")
	}
}
