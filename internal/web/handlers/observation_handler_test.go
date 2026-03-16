package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"arandu/internal/domain/observation"
)

type mockObservationHandlerService struct {
	observations map[string]*observation.Observation
}

func (m *mockObservationHandlerService) GetObservation(ctx context.Context, id string) (*observation.Observation, error) {
	return m.observations[id], nil
}

func (m *mockObservationHandlerService) UpdateObservation(ctx context.Context, id, content string) error {
	obs, exists := m.observations[id]
	if !exists {
		return fmt.Errorf("observation not found")
	}

	if content == "" {
		return fmt.Errorf("observation content cannot be empty")
	}

	if len(content) > 5000 {
		return fmt.Errorf("observation content cannot exceed 5000 characters")
	}

	obs.Content = content
	obs.UpdatedAt = time.Now()
	return nil
}

func TestObservationHandler_GetObservation(t *testing.T) {
	service := &mockObservationHandlerService{
		observations: map[string]*observation.Observation{
			"obs-123": {
				ID:        "obs-123",
				SessionID: "session-456",
				Content:   "Conteúdo original da observação",
				CreatedAt: time.Now(),
			},
		},
	}

	handler := NewObservationHandler(service)

	req := httptest.NewRequest("GET", "/observations/obs-123", nil)
	w := httptest.NewRecorder()

	handler.GetObservation(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetObservation() status = %v, want %v", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Conteúdo original da observação") {
		t.Errorf("GetObservation() body doesn't contain observation content")
	}

	if !strings.Contains(body, "obs-123") {
		t.Errorf("GetObservation() body doesn't contain observation ID")
	}
}

func TestObservationHandler_GetObservationEditForm(t *testing.T) {
	service := &mockObservationHandlerService{
		observations: map[string]*observation.Observation{
			"obs-123": {
				ID:        "obs-123",
				SessionID: "session-456",
				Content:   "Conteúdo original da observação",
				CreatedAt: time.Now(),
			},
		},
	}

	handler := NewObservationHandler(service)

	req := httptest.NewRequest("GET", "/observations/obs-123/edit", nil)
	w := httptest.NewRecorder()

	handler.GetObservationEditForm(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetObservationEditForm() status = %v, want %v", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Conteúdo original da observação") {
		t.Errorf("GetObservationEditForm() body doesn't contain observation content")
	}

	if !strings.Contains(body, "hx-put") {
		t.Errorf("GetObservationEditForm() body doesn't contain hx-put attribute")
	}

	if !strings.Contains(body, "/observations/obs-123") {
		t.Errorf("GetObservationEditForm() body doesn't contain correct action URL")
	}
}

func TestObservationHandler_UpdateObservation(t *testing.T) {
	service := &mockObservationHandlerService{
		observations: map[string]*observation.Observation{
			"obs-123": {
				ID:        "obs-123",
				SessionID: "session-456",
				Content:   "Conteúdo original da observação",
				CreatedAt: time.Now(),
			},
		},
	}

	handler := NewObservationHandler(service)

	// Test successful update
	formData := "content=Conteúdo%20atualizado%20da%20observação"
	req := httptest.NewRequest("PUT", "/observations/obs-123", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateObservation(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("UpdateObservation() status = %v, want %v", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Conteúdo atualizado da observação") {
		t.Errorf("UpdateObservation() body doesn't contain updated content")
	}

	// Verify the observation was updated in the service
	obs, _ := service.GetObservation(context.Background(), "obs-123")
	if obs.Content != "Conteúdo atualizado da observação" {
		t.Errorf("UpdateObservation() didn't update content in service, got %v", obs.Content)
	}
}

func TestObservationHandler_UpdateObservation_EmptyContent(t *testing.T) {
	service := &mockObservationHandlerService{
		observations: map[string]*observation.Observation{
			"obs-123": {
				ID:        "obs-123",
				SessionID: "session-456",
				Content:   "Conteúdo original",
				CreatedAt: time.Now(),
			},
		},
	}

	handler := NewObservationHandler(service)

	formData := "content="
	req := httptest.NewRequest("PUT", "/observations/obs-123", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateObservation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdateObservation() with empty content status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestObservationHandler_UpdateObservation_NotFound(t *testing.T) {
	service := &mockObservationHandlerService{
		observations: map[string]*observation.Observation{},
	}

	handler := NewObservationHandler(service)

	formData := "content=Novo%20conteúdo"
	req := httptest.NewRequest("PUT", "/observations/non-existent", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateObservation(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("UpdateObservation() with non-existent observation status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
