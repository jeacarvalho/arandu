package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"arandu/internal/domain/intervention"

	layoutComponents "arandu/web/components/layout"
	sessionComponents "arandu/web/components/session"
)

type InterventionHandlerServiceInterface interface {
	GetIntervention(ctx context.Context, id string) (*intervention.Intervention, error)
	UpdateIntervention(ctx context.Context, id, content string) error
}

type InterventionHandler struct {
	interventionService InterventionHandlerServiceInterface
}

func NewInterventionHandler(interventionService InterventionHandlerServiceInterface) *InterventionHandler {
	return &InterventionHandler{
		interventionService: interventionService,
	}
}

// GetIntervention handles GET /interventions/{id} - returns intervention item component
func (h *InterventionHandler) GetIntervention(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract intervention ID from URL
	id := extractInterventionID(r.URL.Path)
	if id == "" {
		http.Error(w, "ID da intervenção é obrigatório", http.StatusBadRequest)
		return
	}

	// Get intervention
	intv, err := h.interventionService.GetIntervention(r.Context(), id)
	if err != nil {
		log.Printf("Error getting intervention: %v", err)
		http.Error(w, "Erro ao buscar intervenção", http.StatusInternalServerError)
		return
	}

	if intv == nil {
		http.Error(w, "Intervenção não encontrada", http.StatusNotFound)
		return
	}

	// Render intervention item component
	component := sessionComponents.InterventionItem(sessionComponents.InterventionItemData{
		ID:        intv.ID,
		Content:   intv.Content,
		CreatedAt: intv.CreatedAt.Format("02/01/2006 15:04"),
	})

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering intervention item: %v", err)
			http.Error(w, "Erro ao renderizar componente", http.StatusInternalServerError)
		}
		return
	}
	layoutComponents.Shell(layoutComponents.ShellConfig{
		PageTitle:   "Intervenção",
		ShowSidebar: true,
	}, component).Render(r.Context(), w)
}

// GetInterventionEditForm handles GET /interventions/{id}/edit - returns edit form component
func (h *InterventionHandler) GetInterventionEditForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract intervention ID from URL
	id := extractInterventionID(r.URL.Path)
	if id == "" {
		http.Error(w, "ID da intervenção é obrigatório", http.StatusBadRequest)
		return
	}

	// Get intervention
	intv, err := h.interventionService.GetIntervention(r.Context(), id)
	if err != nil {
		log.Printf("Error getting intervention: %v", err)
		http.Error(w, "Erro ao buscar intervenção", http.StatusInternalServerError)
		return
	}

	if intv == nil {
		http.Error(w, "Intervenção não encontrada", http.StatusNotFound)
		return
	}

	// Render edit form component
	component := sessionComponents.InterventionEditForm(sessionComponents.InterventionEditData{
		ID:      intv.ID,
		Content: intv.Content,
	})

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering intervention edit form: %v", err)
			http.Error(w, "Erro ao renderizar formulário", http.StatusInternalServerError)
		}
		return
	}
	layoutComponents.Shell(layoutComponents.ShellConfig{
		PageTitle:   "Editar Intervenção",
		ShowSidebar: true,
	}, component).Render(r.Context(), w)
}

// UpdateIntervention handles PUT /interventions/{id} - updates intervention
func (h *InterventionHandler) UpdateIntervention(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract intervention ID from URL
	id := extractInterventionID(r.URL.Path)
	if id == "" {
		http.Error(w, "ID da intervenção é obrigatório", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Conteúdo da intervenção não pode ser vazio", http.StatusBadRequest)
		return
	}

	// Update intervention
	err := h.interventionService.UpdateIntervention(r.Context(), id, content)
	if err != nil {
		log.Printf("Error updating intervention: %v", err)
		http.Error(w, "Erro ao atualizar intervenção: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated intervention
	intv, err := h.interventionService.GetIntervention(r.Context(), id)
	if err != nil {
		log.Printf("Error getting updated intervention: %v", err)
		http.Error(w, "Erro ao buscar intervenção atualizada", http.StatusInternalServerError)
		return
	}

	// Render updated intervention item component
	component := sessionComponents.InterventionItem(sessionComponents.InterventionItemData{
		ID:        intv.ID,
		Content:   intv.Content,
		CreatedAt: intv.CreatedAt.Format("02/01/2006 15:04"),
	})

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering updated intervention item: %v", err)
			http.Error(w, "Erro ao renderizar componente", http.StatusInternalServerError)
		}
		return
	}
	layoutComponents.Shell(layoutComponents.ShellConfig{
		PageTitle:   "Intervenção",
		ShowSidebar: true,
	}, component).Render(r.Context(), w)
}

// Helper function to extract intervention ID from URL path
func extractInterventionID(path string) string {
	// Remove /interventions/ prefix
	path = strings.TrimPrefix(path, "/interventions/")

	// Remove /edit suffix if present
	path = strings.TrimSuffix(path, "/edit")

	// Trim any trailing slashes
	path = strings.Trim(path, "/")

	return path
}
