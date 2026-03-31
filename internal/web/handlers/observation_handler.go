package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"arandu/internal/domain/observation"

	sessionComponents "arandu/web/components/session"
)

type ObservationHandlerServiceInterface interface {
	GetObservation(ctx context.Context, id string) (*observation.Observation, error)
	UpdateObservation(ctx context.Context, id, content string) error
}

type ObservationHandler struct {
	observationService ObservationHandlerServiceInterface
}

func NewObservationHandler(observationService ObservationHandlerServiceInterface) *ObservationHandler {
	return &ObservationHandler{
		observationService: observationService,
	}
}

// GetObservation handles GET /observations/{id} - returns observation item component
func (h *ObservationHandler) GetObservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract observation ID from URL
	id := extractObservationID(r.URL.Path)
	if id == "" {
		http.Error(w, "ID da observação é obrigatório", http.StatusBadRequest)
		return
	}

	// Get observation
	obs, err := h.observationService.GetObservation(r.Context(), id)
	if err != nil {
		log.Printf("Error getting observation: %v", err)
		http.Error(w, "Erro ao buscar observação", http.StatusInternalServerError)
		return
	}

	if obs == nil {
		http.Error(w, "Observação não encontrada", http.StatusNotFound)
		return
	}

	// Render observation item component
	component := sessionComponents.ObservationItem(sessionComponents.ObservationItemData{
		ID:        obs.ID,
		Content:   obs.Content,
		CreatedAt: obs.CreatedAt.Format("02/01/2006 15:04"),
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering observation item: %v", err)
		http.Error(w, "Erro ao renderizar componente", http.StatusInternalServerError)
	}
}

// GetObservationEditForm handles GET /observations/{id}/edit - returns edit form component
func (h *ObservationHandler) GetObservationEditForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract observation ID from URL
	id := extractObservationID(r.URL.Path)
	if id == "" {
		http.Error(w, "ID da observação é obrigatório", http.StatusBadRequest)
		return
	}

	// Get observation
	obs, err := h.observationService.GetObservation(r.Context(), id)
	if err != nil {
		log.Printf("Error getting observation: %v", err)
		http.Error(w, "Erro ao buscar observação", http.StatusInternalServerError)
		return
	}

	if obs == nil {
		http.Error(w, "Observação não encontrada", http.StatusNotFound)
		return
	}

	// Render edit form component
	component := sessionComponents.ObservationEditForm(sessionComponents.ObservationEditFormData{
		ID:      obs.ID,
		Content: obs.Content,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering observation edit form: %v", err)
		http.Error(w, "Erro ao renderizar formulário", http.StatusInternalServerError)
	}
}

// UpdateObservation handles PUT /observations/{id} - updates observation
func (h *ObservationHandler) UpdateObservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract observation ID from URL
	id := extractObservationID(r.URL.Path)
	if id == "" {
		http.Error(w, "ID da observação é obrigatório", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Conteúdo da observação não pode ser vazio", http.StatusBadRequest)
		return
	}

	// Update observation
	err := h.observationService.UpdateObservation(r.Context(), id, content)
	if err != nil {
		log.Printf("Error updating observation: %v", err)
		http.Error(w, "Erro ao atualizar observação: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated observation
	obs, err := h.observationService.GetObservation(r.Context(), id)
	if err != nil {
		log.Printf("Error getting updated observation: %v", err)
		http.Error(w, "Erro ao buscar observação atualizada", http.StatusInternalServerError)
		return
	}

	// Render updated observation item component
	component := sessionComponents.ObservationItem(sessionComponents.ObservationItemData{
		ID:        obs.ID,
		Content:   obs.Content,
		CreatedAt: obs.CreatedAt.Format("02/01/2006 15:04"),
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering updated observation item: %v", err)
		http.Error(w, "Erro ao renderizar componente", http.StatusInternalServerError)
	}
}

// Helper function to extract observation ID from URL path
func extractObservationID(path string) string {
	// Remove /observations/ prefix
	path = strings.TrimPrefix(path, "/observations/")

	// Remove /edit suffix if present
	path = strings.TrimSuffix(path, "/edit")

	// Trim any trailing slashes
	path = strings.Trim(path, "/")

	return path
}
