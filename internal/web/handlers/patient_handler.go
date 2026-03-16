package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"

	layoutComponents "arandu/web/components/layout"
	patientComponents "arandu/web/components/patient"
)

// PatientViewData is a ViewModel that protects the domain from template concerns
type PatientViewData struct {
	Patient  *PatientViewModel
	Sessions []*SessionViewModel
	Insights []InsightViewModel
	Error    string
}

// PatientViewModel is a view-specific representation of a patient
type PatientViewModel struct {
	ID        string
	Name      string
	Notes     string
	CreatedAt string
	UpdatedAt string
}

// SessionViewModel is a view-specific representation of a session
type SessionViewModel struct {
	ID        string
	PatientID string
	Date      string
	Summary   string
	CreatedAt string
	UpdatedAt string
}

// InsightViewModel represents an insight for the sidebar
type InsightViewModel struct {
	Content   string
	Source    string
	CreatedAt string
}

// PatientService defines the interface for patient operations (dependency inversion)
type PatientService interface {
	GetPatientByID(ctx context.Context, id string) (*patient.Patient, error)
	ListPatients(ctx context.Context) ([]*patient.Patient, error)
	CreatePatient(ctx context.Context, input services.CreatePatientInput) (*patient.Patient, error)
}

// SessionService defines the interface for session operations
type SessionService interface {
	ListSessionsByPatient(ctx context.Context, patientID string) ([]*session.Session, error)
	CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error)
	GetSession(ctx context.Context, id string) (*session.Session, error)
	UpdateSession(ctx context.Context, input services.UpdateSessionInput) error
}

// InsightService defines the interface for insight operations
type InsightService interface {
	GetInsightsByPatient(ctx context.Context, patientID string, limit int) ([]interface{}, error)
}

// PatientHandler handles HTTP requests related to patients
type PatientHandler struct {
	patientService PatientService
	sessionService SessionService
	insightService InsightService
	templates      TemplateRenderer
}

// TemplateRenderer defines the interface for template rendering (deprecated - use templ components)
type TemplateRenderer interface {
	ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error
}

// Render handles HTMX vs full page rendering automatically
func (h *PatientHandler) Render(w http.ResponseWriter, r *http.Request, pageTemplate string, data interface{}) {
	// This method is deprecated - use templ components directly
	// For now, do nothing
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Render method deprecated - use templ components"))
}

// NewPatientHandler creates a new PatientHandler with dependency injection
func NewPatientHandler(
	patientService PatientService,
	sessionService SessionService,
	insightService InsightService,
	templates TemplateRenderer,
) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
		sessionService: sessionService,
		insightService: insightService,
		templates:      templates,
	}
}

// mapPatientToViewModel maps domain patient to view model
func mapPatientToViewModel(p *patient.Patient) *PatientViewModel {
	if p == nil {
		return nil
	}
	return &PatientViewModel{
		ID:        p.ID,
		Name:      p.Name,
		Notes:     p.Notes,
		CreatedAt: p.CreatedAt.Format("Jan 2006"),
		UpdatedAt: p.UpdatedAt.Format("02/01/2006"),
	}
}

// mapSessionsToViewModel maps domain sessions to view models
func mapSessionsToViewModel(sessions []*session.Session) []*SessionViewModel {
	result := make([]*SessionViewModel, len(sessions))
	for i, s := range sessions {
		result[i] = &SessionViewModel{
			ID:        s.ID,
			PatientID: s.PatientID,
			Date:      s.Date.Format("02/01/2006"),
			Summary:   s.Summary,
			CreatedAt: s.CreatedAt.Format("02/01/2006 15:04"),
			UpdatedAt: s.UpdatedAt.Format("02/01/2006 15:04"),
		}
	}
	return result
}

// renderError handles error rendering with HTMX awareness
func (h *PatientHandler) renderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	data := PatientViewData{
		Error: message,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.templates.ExecuteTemplate(w, "error-fragment", data)
		return
	}

	h.templates.ExecuteTemplate(w, "layout", data)
}

// getInsights retrieves insights for the sidebar (mock implementation)
func (h *PatientHandler) getInsights(ctx context.Context, patientID string) []InsightViewModel {
	if h.insightService == nil {
		return []InsightViewModel{}
	}

	insights, err := h.insightService.GetInsightsByPatient(ctx, patientID, 5)
	if err != nil {
		return []InsightViewModel{}
	}

	result := make([]InsightViewModel, len(insights))
	for i, insight := range insights {
		// Type assertion based on expected insight structure
		if ins, ok := insight.(map[string]interface{}); ok {
			content, _ := ins["content"].(string)
			source, _ := ins["source"].(string)
			createdAt, _ := ins["created_at"].(string)
			result[i] = InsightViewModel{
				Content:   content,
				Source:    source,
				CreatedAt: createdAt,
			}
		}
	}
	return result
}

// ListPatients handles GET /patients - lists all patients
func (h *PatientHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	patients, err := h.patientService.ListPatients(ctx)
	if err != nil {
		h.renderError(w, r, "Erro ao listar pacientes", http.StatusInternalServerError)
		return
	}

	// Map to templ components
	patientItems := make([]patientComponents.PatientListItem, len(patients))
	for i, p := range patients {
		patientItems[i] = patientComponents.PatientListItem{
			ID:        p.ID,
			Name:      p.Name,
			Notes:     p.Notes,
			CreatedAt: p.CreatedAt.Format("02/01/2006"),
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTMX-aware rendering
	if r.Header.Get("HX-Request") == "true" {
		// Render just the patient list fragment
		patientComponents.PatientList(patientItems, "").Render(r.Context(), w)
		return
	}

	// Render with layout using templ
	patientList := patientComponents.PatientList(patientItems, "")
	layoutComponents.BaseWithContent("Pacientes", patientList).Render(r.Context(), w)
}

// PatientsViewData is a ViewModel for the patients list page
type PatientsViewData struct {
	Patients []*PatientViewModel
	Insights []InsightViewModel
	Error    string
}

// Show handles GET /patient/{id} - shows patient details
func (h *PatientHandler) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("PatientHandler.Show called for path: %s", r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Extração de Parâmetros
	id := extractIDFromPath(r.URL.Path, "/patient/")
	if id == "" {
		h.renderError(w, r, "ID do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	// 2. Chamada ao Serviço
	patient, err := h.patientService.GetPatientByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrPatientNotFound) {
			h.renderError(w, r, "Paciente não encontrado", http.StatusNotFound)
			return
		}
		h.renderError(w, r, "Erro ao buscar paciente", http.StatusInternalServerError)
		return
	}

	// Get sessions
	sessions, err := h.sessionService.ListSessionsByPatient(r.Context(), id)
	if err != nil {
		sessions = []*session.Session{}
	}

	// Map to templ components
	patientDetail := patientComponents.PatientDetailItem{
		ID:        patient.ID,
		Name:      patient.Name,
		Notes:     patient.Notes,
		CreatedAt: patient.CreatedAt.Format("02/01/2006 às 15:04"),
	}

	sessionItems := make([]patientComponents.SessionItem, len(sessions))
	for i, s := range sessions {
		sessionItems[i] = patientComponents.SessionItem{
			ID:            s.ID,
			SessionNumber: i + 1,
			Date:          s.Date.Format("02/01/2006"),
			Summary:       s.Summary,
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTMX-aware rendering
	if r.Header.Get("HX-Request") == "true" {
		patientComponents.PatientDetail(patientDetail, sessionItems).Render(r.Context(), w)
		return
	}

	detail := patientComponents.PatientDetail(patientDetail, sessionItems)
	layoutComponents.BaseWithContent(patient.Name, detail).Render(r.Context(), w)
}

// NewPatient handles GET /patients/new - shows new patient form
func (h *PatientHandler) NewPatient(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with templ components
	http.Error(w, "Not implemented - use templ components", http.StatusNotImplemented)
}

// NewPatientViewData is a ViewModel for the new patient form
type NewPatientViewData struct {
	Error       string
	FormData    *PatientFormValues
	ServerError string
}

// PatientFormValues holds form data for patient creation/update
type PatientFormValues struct {
	Name  string
	Notes string
}

// CreatePatient handles POST /patients - creates a new patient
func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with templ components
	http.Error(w, "Not implemented - use templ components", http.StatusNotImplemented)
}

// Helper function to convert interface{} insights to InsightViewModel
func convertInsights(rawInsights []interface{}) []InsightViewModel {
	result := make([]InsightViewModel, 0, len(rawInsights))
	for _, insight := range rawInsights {
		if ins, ok := insight.(map[string]interface{}); ok {
			content, _ := ins["content"].(string)
			source, _ := ins["source"].(string)
			createdAt, _ := ins["created_at"].(string)
			result = append(result, InsightViewModel{
				Content:   content,
				Source:    source,
				CreatedAt: createdAt,
			})
		}
	}
	return result
}

// extractIDFromPath extracts an ID from a URL path given a prefix
// e.g., extractIDFromPath("/patient/123", "/patient/") returns "123"
func extractIDFromPath(path, prefix string) string {
	id := trimPrefix(path, prefix)
	if id == "" {
		return ""
	}
	// Handle trailing slashes and get first segment
	id = trimSuffix(id, "/")
	// Get first segment if there are more path parts
	for i, c := range id {
		if c == '/' {
			return id[:i]
		}
	}
	return id
}
