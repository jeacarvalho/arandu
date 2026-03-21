package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/intervention"
	"arandu/internal/domain/observation"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"

	layoutComponents "arandu/web/components/layout"
	sessionComponents "arandu/web/components/session"
)

// SessionViewData is a ViewModel that protects the domain from template concerns
type SessionViewData struct {
	Session     *SessionDetailViewModel
	Patient     *PatientViewModel
	Insights    []InsightViewModel
	Error       string
	FormData    *SessionFormValues
	ServerError string
}

// SessionDetailViewModel is a view-specific representation of a session with full details
type SessionDetailViewModel struct {
	ID            string
	PatientID     string
	Date          string
	Summary       string
	CreatedAt     string
	UpdatedAt     string
	Observations  []ObservationViewModel
	Interventions []InterventionViewModel
}

// ObservationViewModel represents a clinical observation
type ObservationViewModel struct {
	ID        string
	Content   string
	CreatedAt string
}

// InterventionViewModel represents a therapeutic intervention
type InterventionViewModel struct {
	ID        string
	Content   string
	CreatedAt string
}

// SessionFormValues holds form data for session creation/update
type SessionFormValues struct {
	PatientID string
	Date      string
	Summary   string
}

// SessionServiceInterface defines the interface for session operations (dependency inversion)
type SessionServiceInterface interface {
	GetSession(ctx context.Context, id string) (*session.Session, error)
	ListSessionsByPatient(ctx context.Context, patientID string) ([]*session.Session, error)
	CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error)
	UpdateSession(ctx context.Context, input services.UpdateSessionInput) error
}

// PatientServiceInterface defines the interface for patient operations
type PatientServiceInterface interface {
	GetPatientByID(ctx context.Context, id string) (*patient.Patient, error)
}

type ObservationServiceInterface interface {
	CreateObservation(ctx context.Context, sessionID, content string) (interface{}, error)
	GetObservationsBySession(ctx context.Context, sessionID string) ([]interface{}, error)
}

type InterventionServiceInterface interface {
	CreateIntervention(ctx context.Context, sessionID, content string) (interface{}, error)
	GetInterventionsBySession(ctx context.Context, sessionID string) ([]interface{}, error)
}

// SessionHandler handles HTTP requests related to sessions
type SessionHandler struct {
	sessionService      SessionServiceInterface
	patientService      PatientServiceInterface
	observationService  ObservationServiceInterface
	interventionService InterventionServiceInterface
}

// NewSessionHandler creates a new SessionHandler with dependency injection
func NewSessionHandler(
	sessionService SessionServiceInterface,
	patientService PatientServiceInterface,
	observationService ObservationServiceInterface,
	interventionService InterventionServiceInterface,
) *SessionHandler {
	return &SessionHandler{
		sessionService:      sessionService,
		patientService:      patientService,
		observationService:  observationService,
		interventionService: interventionService,
	}
}

// mapSessionToDetailViewModel maps domain session to detail view model
func mapSessionToDetailViewModel(s *session.Session) *SessionDetailViewModel {
	if s == nil {
		return nil
	}
	return &SessionDetailViewModel{
		ID:            s.ID,
		PatientID:     s.PatientID,
		Date:          s.Date.Format("02/01/2006"),
		Summary:       s.Summary,
		CreatedAt:     s.CreatedAt.Format("02/01/2006 15:04"),
		UpdatedAt:     s.UpdatedAt.Format("02/01/2006 15:04"),
		Observations:  []ObservationViewModel{},
		Interventions: []InterventionViewModel{},
	}
}

// renderError handles error rendering with HTMX awareness
func (h *SessionHandler) renderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	isHTMX := r.Header.Get("HX-Request") == "true"
	errorData := layoutComponents.ErrorData{
		Error:  message,
		IsHTMX: isHTMX,
		Title:  "Erro",
	}

	if isHTMX {
		layoutComponents.ErrorFragment(errorData).Render(r.Context(), w)
		return
	}

	// Render full page with layout
	errorComponent := layoutComponents.ErrorFragment(errorData)
	layoutComponents.BaseWithContent("Erro", errorComponent).Render(r.Context(), w)
}

// Show handles GET /session/{id} - shows session details
func (h *SessionHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Extração de Parâmetros
	id := extractIDFromPath(r.URL.Path, "/session/")
	if id == "" {
		h.renderError(w, r, "ID da sessão é obrigatório", http.StatusBadRequest)
		return
	}

	// 2. Chamada ao Serviço
	sess, err := h.sessionService.GetSession(r.Context(), id)
	if err != nil {
		h.renderError(w, r, "Sessão não encontrada", http.StatusNotFound)
		return
	}

	// Get patient info
	patient, err := h.patientService.GetPatientByID(r.Context(), sess.PatientID)
	if err != nil {
		h.renderError(w, r, "Erro ao buscar paciente", http.StatusInternalServerError)
		return
	}

	// Map to templ components
	sessionDetail := sessionComponents.SessionDetail{
		ID:          sess.ID,
		PatientID:   sess.PatientID,
		PatientName: patient.Name,
		Date:        sess.Date.Format("02/01/2006"),
		Summary:     sess.Summary,
		CreatedAt:   sess.CreatedAt.Format("02/01/2006 às 15:04"),
	}

	// Fetch observations from service
	obsList, err := h.observationService.GetObservationsBySession(r.Context(), sess.ID)
	if err != nil {
		log.Printf("Error fetching observations: %v", err)
	}

	// Convert to component observations
	observations := []sessionComponents.Observation{}
	for _, obs := range obsList {
		if o, ok := obs.(*observation.Observation); ok {
			observations = append(observations, sessionComponents.Observation{
				ID:        o.ID,
				Content:   o.Content,
				CreatedAt: o.CreatedAt.Format("02/01/2006 às 15:04"),
			})
		}
	}

	// Empty interventions for now (can be expanded later)
	interventions := []sessionComponents.Intervention{}

	// Fetch interventions from service
	intervList, err := h.interventionService.GetInterventionsBySession(r.Context(), sess.ID)
	if err != nil {
		log.Printf("Error fetching interventions: %v", err)
	}

	// Convert to component interventions
	for _, intv := range intervList {
		if i, ok := intv.(*intervention.Intervention); ok {
			interventions = append(interventions, sessionComponents.Intervention{
				ID:        i.ID,
				Content:   i.Content,
				CreatedAt: i.CreatedAt.Format("02/01/2006 às 15:04"),
			})
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	detail := sessionComponents.SessionDetailView(sessionDetail, observations, interventions, sess.PatientID, "")
	layoutComponents.BaseWithContent("Sessão "+sessionDetail.Date, detail).Render(r.Context(), w)
}

// NewSession handles GET /patient/{id}/sessions/new - shows new session form
func (h *SessionHandler) NewSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("NewSession called: %s", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract patient ID from URL path
	patientID := extractPatientIDFromPath(r.URL.Path)
	log.Printf("Extracted patientID: %s", patientID)
	if patientID == "" {
		h.renderError(w, r, "ID do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	// Get patient info
	patient, err := h.patientService.GetPatientByID(r.Context(), patientID)
	if err != nil {
		log.Printf("Error getting patient: %v", err)
		h.renderError(w, r, "Paciente não encontrado: "+patientID, http.StatusNotFound)
		return
	}
	log.Printf("Found patient: %s", patient.Name)

	formData := sessionComponents.NewSessionFormData{
		PatientName: patient.Name,
		FormData: &sessionComponents.SessionFormValues{
			PatientID: patientID,
			Date:      time.Now().Format("2006-01-02"),
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTMX-aware rendering
	if r.Header.Get("HX-Request") == "true" {
		// Render just the form fragment
		sessionComponents.NewSessionForm(formData).Render(r.Context(), w)
		return
	}

	// Render with layout using templ
	form := sessionComponents.NewSessionForm(formData)
	layoutComponents.BaseWithContent("Nova Sessão", form).Render(r.Context(), w)
}

// CreateSession handles POST /sessions - creates a new session
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, r, "Dados do formulário inválidos", http.StatusBadRequest)
		return
	}

	patientID := r.FormValue("patient_id")
	if patientID == "" {
		h.renderError(w, r, "ID do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("date")
	if dateStr == "" {
		h.renderError(w, r, "Data é obrigatória", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.renderError(w, r, "Formato de data inválido", http.StatusBadRequest)
		return
	}

	summary := r.FormValue("summary")

	sess, err := h.sessionService.CreateSession(r.Context(), patientID, date, summary)
	if err != nil {
		// For HTMX requests, return form with error
		if r.Header.Get("HX-Request") == "true" {
			patient, _ := h.patientService.GetPatientByID(r.Context(), patientID)
			formData := sessionComponents.NewSessionFormData{
				Error:       err.Error(),
				PatientName: patient.Name,
				FormData: &sessionComponents.SessionFormValues{
					PatientID: patientID,
					Date:      dateStr,
					Summary:   summary,
				},
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			sessionComponents.NewSessionForm(formData).Render(r.Context(), w)
			return
		}

		h.renderError(w, r, "Erro ao criar sessão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Instead of redirecting, render the edit page directly for the wizard flow
	// This avoids any client-side JavaScript interference with redirects

	// Get patient info
	patient, err := h.patientService.GetPatientByID(r.Context(), sess.PatientID)
	if err != nil {
		log.Printf("Error getting patient for wizard: %v", err)
		http.Redirect(w, r, "/patient/"+sess.PatientID, http.StatusSeeOther)
		return
	}

	// Create form data for the edit page
	formData := sessionComponents.EditSessionFormData{
		SessionID:   sess.ID,
		PatientName: patient.Name,
		FormData: &sessionComponents.SessionFormValues{
			PatientID: sess.PatientID,
			Date:      sess.Date.Format("2006-01-02"),
			Summary:   sess.Summary,
		},
		Observations:  []sessionComponents.Observation{},
		Interventions: []sessionComponents.Intervention{},
	}

	// Render the edit page directly
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	form := sessionComponents.EditSessionForm(formData)
	layoutComponents.BaseWithContent("Completar Sessão", form).Render(r.Context(), w)
}

// EditSession handles GET /sessions/edit/{id} - shows edit session form
func (h *SessionHandler) EditSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("EditSession called: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session ID from URL path /session/{id}/edit
	// Remove "/edit" suffix first
	pathWithoutEdit := strings.TrimSuffix(r.URL.Path, "/edit")
	log.Printf("Path without edit: %s", pathWithoutEdit)

	// Extract ID from /session/{id}
	sessionID := extractSessionIDFromPath(pathWithoutEdit, "/session/")
	log.Printf("Extracted sessionID: %s", sessionID)

	if sessionID == "" {
		h.renderError(w, r, "ID da sessão é obrigatório", http.StatusBadRequest)
		return
	}

	// Get session
	sess, err := h.sessionService.GetSession(r.Context(), sessionID)
	if err != nil {
		h.renderError(w, r, "Sessão não encontrada", http.StatusNotFound)
		return
	}

	// Get patient
	patient, err := h.patientService.GetPatientByID(r.Context(), sess.PatientID)
	if err != nil {
		h.renderError(w, r, "Erro ao buscar paciente", http.StatusInternalServerError)
		return
	}

	// Get observations and interventions for this session
	obsList, err := h.observationService.GetObservationsBySession(r.Context(), sessionID)
	if err != nil {
		log.Printf("Error getting observations: %v", err)
		obsList = []interface{}{}
	}

	intList, err := h.interventionService.GetInterventionsBySession(r.Context(), sessionID)
	if err != nil {
		log.Printf("Error getting interventions: %v", err)
		intList = []interface{}{}
	}

	// Convert to component observations
	observations := []sessionComponents.Observation{}
	for _, obs := range obsList {
		if o, ok := obs.(*observation.Observation); ok {
			observations = append(observations, sessionComponents.Observation{
				ID:        o.ID,
				Content:   o.Content,
				CreatedAt: o.CreatedAt.Format("02/01/2006 15:04"),
			})
		}
	}

	// Convert to component interventions
	interventions := []sessionComponents.Intervention{}
	for _, interv := range intList {
		if i, ok := interv.(*intervention.Intervention); ok {
			interventions = append(interventions, sessionComponents.Intervention{
				ID:        i.ID,
				Content:   i.Content,
				CreatedAt: i.CreatedAt.Format("02/01/2006 15:04"),
			})
		}
	}

	formData := sessionComponents.EditSessionFormData{
		SessionID:   sessionID,
		PatientName: patient.Name,
		FormData: &sessionComponents.SessionFormValues{
			PatientID: sess.PatientID,
			Date:      sess.Date.Format("2006-01-02"),
			Summary:   sess.Summary,
		},
		Observations:  observations,
		Interventions: interventions,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTMX-aware rendering
	if r.Header.Get("HX-Request") == "true" {
		// Render just the form fragment
		sessionComponents.EditSessionForm(formData).Render(r.Context(), w)
		return
	}

	// Render with layout using templ
	form := sessionComponents.EditSessionForm(formData)
	layoutComponents.BaseWithContent("Completar Sessão", form).Render(r.Context(), w)
}

// UpdateSession handles POST /sessions/update - updates an existing session
func (h *SessionHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateSession called: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, r, "Dados do formulário inválidos", http.StatusBadRequest)
		return
	}

	// Debug: log all form values
	log.Printf("Form values: %v", r.Form)

	sessionID := r.FormValue("session_id")
	log.Printf("Extracted session_id: %s", sessionID)
	if sessionID == "" {
		h.renderError(w, r, "ID da sessão é obrigatório", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("date")
	if dateStr == "" {
		h.renderError(w, r, "Data é obrigatória", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		h.renderError(w, r, "Formato de data inválido", http.StatusBadRequest)
		return
	}

	summary := r.FormValue("summary")

	input := services.UpdateSessionInput{
		ID:      sessionID,
		Date:    date,
		Summary: summary,
	}

	err = h.sessionService.UpdateSession(r.Context(), input)
	if err != nil {
		if errors.Is(err, errors.New("session not found")) {
			h.renderError(w, r, "Sessão não encontrada", http.StatusNotFound)
			return
		}

		// For HTMX requests, return form with error
		if r.Header.Get("HX-Request") == "true" {
			sess, _ := h.sessionService.GetSession(r.Context(), sessionID)
			patient, _ := h.patientService.GetPatientByID(r.Context(), sess.PatientID)
			formData := sessionComponents.EditSessionFormData{
				Error:       err.Error(),
				SessionID:   sessionID,
				PatientName: patient.Name,
				FormData: &sessionComponents.SessionFormValues{
					PatientID: sess.PatientID,
					Date:      dateStr,
					Summary:   summary,
				},
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			sessionComponents.EditSessionForm(formData).Render(r.Context(), w)
			return
		}

		h.renderError(w, r, "Erro ao atualizar sessão: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to session edit page to continue adding observations/interventions
	log.Printf("Redirecting to: /session/%s/edit", sessionID)
	http.Redirect(w, r, "/session/"+sessionID+"/edit", http.StatusSeeOther)
}

// extractPatientIDFromPath extracts patient ID from URL path like /patient/{id}/sessions/new or /patients/{id}/context
func extractPatientIDFromPath(path string) string {
	parts := splitPath(path)
	for i, part := range parts {
		if (part == "patient" || part == "patients") && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// extractSessionIDFromPath extracts session ID from URL path
func extractSessionIDFromPath(path, prefix string) string {
	id := trimPrefix(path, prefix)
	if id == "" {
		return ""
	}
	// Handle trailing slashes
	id = trimSuffix(id, "/")
	return id
}

// Helper functions to avoid importing strings package
func splitPath(path string) []string {
	result := make([]string, 0)
	current := ""
	for _, c := range path {
		if c == '/' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func trimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// CreateObservation handles POST /sessions/{id}/observations
func (h *SessionHandler) CreateObservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract session ID from URL
	sessionID := extractSessionID(r.URL.Path, "/session/", "/observations")
	if sessionID == "" {
		h.renderError(w, r, "ID da sessão não encontrado", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		h.renderError(w, r, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		h.renderError(w, r, "Conteúdo da observação não pode ser vazio", http.StatusBadRequest)
		return
	}

	// Create observation
	obs, err := h.observationService.CreateObservation(ctx, sessionID, content)
	if err != nil {
		h.renderError(w, r, "Erro ao criar observação: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to domain observation
	observation, ok := obs.(*observation.Observation)
	if !ok {
		h.renderError(w, r, "Erro ao converter observação", http.StatusInternalServerError)
		return
	}

	// Render observation item component
	obsVM := ObservationViewModel{
		ID:        observation.ID,
		Content:   observation.Content,
		CreatedAt: observation.CreatedAt.Format("02/01/2006 15:04"),
	}

	// Use templ component
	component := sessionComponents.ObservationItem(sessionComponents.ObservationItemData{
		ID:        obsVM.ID,
		Content:   obsVM.Content,
		CreatedAt: obsVM.CreatedAt,
	})

	// Render HTMX fragment
	w.Header().Set("Content-Type", "text/html")
	if err := component.Render(ctx, w); err != nil {
		log.Printf("Error rendering observation item: %v", err)
		http.Error(w, "Erro ao renderizar componente", http.StatusInternalServerError)
	}
}

// CreateIntervention handles POST /sessions/{id}/interventions
func (h *SessionHandler) CreateIntervention(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract session ID from URL
	sessionID := extractSessionID(r.URL.Path, "/session/", "/interventions")
	if sessionID == "" {
		h.renderError(w, r, "ID da sessão não encontrado", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		h.renderError(w, r, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		h.renderError(w, r, "Conteúdo da intervenção não pode ser vazio", http.StatusBadRequest)
		return
	}

	// Create intervention
	intv, err := h.interventionService.CreateIntervention(ctx, sessionID, content)
	if err != nil {
		h.renderError(w, r, "Erro ao criar intervenção: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to domain intervention
	intervention, ok := intv.(*intervention.Intervention)
	if !ok {
		h.renderError(w, r, "Erro ao converter intervenção", http.StatusInternalServerError)
		return
	}

	// Render intervention item component
	component := sessionComponents.InterventionItem(sessionComponents.InterventionItemData{
		ID:        intervention.ID,
		Content:   intervention.Content,
		CreatedAt: intervention.CreatedAt.Format("02/01/2006 15:04"),
	})

	// Render HTMX fragment
	w.Header().Set("Content-Type", "text/html")
	if err := component.Render(ctx, w); err != nil {
		log.Printf("Error rendering intervention item: %v", err)
		http.Error(w, "Erro ao renderizar componente", http.StatusInternalServerError)
	}
}

// Helper function to extract session ID from URL path
func extractSessionID(path, prefix, suffix string) string {
	path = strings.TrimPrefix(path, prefix)
	path = strings.TrimSuffix(path, suffix)
	path = strings.Trim(path, "/")

	// Split to get just the ID (remove any trailing parts)
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
