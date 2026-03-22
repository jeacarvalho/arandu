package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
	"arandu/internal/domain/timeline"

	layoutComponents "arandu/web/components/layout"
	patientComponents "arandu/web/components/patient"
	sessionComponents "arandu/web/components/session"
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
	ListPatientsPaginated(ctx context.Context, page, pageSize int) ([]*patient.Patient, int, error)
	CreatePatient(ctx context.Context, input services.CreatePatientInput) (*patient.Patient, error)
	SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error)
	GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error)
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

// BiopsychosocialService defines the interface for biopsychosocial context operations
type BiopsychosocialService interface {
	GetMedications(ctx context.Context, patientID string) ([]interface{}, error)
	GetLatestVitals(ctx context.Context, patientID string) (interface{}, error)
	GetAverageVitals(ctx context.Context, patientID string, days int) (interface{}, error)
}

// BiopsychosocialServiceFuncs is a helper type that implements BiopsychosocialService using functions
type BiopsychosocialServiceFuncs struct {
	GetMedicationsFunc   func(ctx context.Context, patientID string) ([]interface{}, error)
	GetLatestVitalsFunc  func(ctx context.Context, patientID string) (interface{}, error)
	GetAverageVitalsFunc func(ctx context.Context, patientID string, days int) (interface{}, error)
}

func (f BiopsychosocialServiceFuncs) GetMedications(ctx context.Context, patientID string) ([]interface{}, error) {
	return f.GetMedicationsFunc(ctx, patientID)
}

func (f BiopsychosocialServiceFuncs) GetLatestVitals(ctx context.Context, patientID string) (interface{}, error) {
	return f.GetLatestVitalsFunc(ctx, patientID)
}

func (f BiopsychosocialServiceFuncs) GetAverageVitals(ctx context.Context, patientID string, days int) (interface{}, error) {
	return f.GetAverageVitalsFunc(ctx, patientID, days)
}

// PatientHandler handles HTTP requests related to patients
type PatientHandler struct {
	patientService         PatientService
	sessionService         SessionService
	insightService         InsightService
	biopsychosocialService BiopsychosocialService
	timelineService        TimelineServicePort
}

// TimelineServicePort defines the interface for timeline operations
type TimelineServicePort interface {
	GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error)
}

// NewPatientHandler creates a new PatientHandler with dependency injection
func NewPatientHandler(
	patientService PatientService,
	sessionService SessionService,
	insightService InsightService,
	biopsychosocialService BiopsychosocialService,
	timelineService TimelineServicePort,
) *PatientHandler {
	return &PatientHandler{
		patientService:         patientService,
		sessionService:         sessionService,
		insightService:         insightService,
		biopsychosocialService: biopsychosocialService,
		timelineService:        timelineService,
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

	// Parse pagination parameters
	page := 1
	pageSize := 20 // Default batch size

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// For HTMX requests, we might want to use offset instead of page
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
			// Calculate page from offset
			page = (offset / pageSize) + 1
		}
	}

	patients, _, err := h.patientService.ListPatientsPaginated(ctx, page, pageSize)
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

	// Check if this is an HTMX request for infinite scroll
	isHTMXRequest := r.Header.Get("HX-Request") == "true"

	listData := patientComponents.PatientListData{
		Patients: patientItems,
		ErrorMsg: "",
		PageSize: pageSize,
		Offset:   offset,
	}

	if isHTMXRequest {
		// For HTMX requests, render just the patient list fragment
		patientComponents.PatientList(listData).Render(r.Context(), w)
	} else {
		// For full page requests, render with layout
		patientList := patientComponents.PatientList(listData)
		layoutComponents.BaseWithContent("Pacientes", patientList).Render(r.Context(), w)
	}
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
	id := extractIDFromPath(r.URL.Path, "/patients/")
	if id == "" {
		h.renderError(w, r, "ID do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	// 2. Chamada ao Serviço
	patientData, err := h.patientService.GetPatientByID(r.Context(), id)
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

	// Get theme frequency for the patient
	themes, err := h.patientService.GetThemeFrequency(r.Context(), id, 10)
	var totalCount int
	themeVM := patientComponents.ThemeCloudViewModel{
		PatientID: id,
		Themes:    []patientComponents.ThemeItem{},
	}
	if err == nil && len(themes) > 0 {
		// Calculate max count for weight normalization
		maxCount := 0
		for _, t := range themes {
			if c, ok := t["count"].(int); ok && c > maxCount {
				maxCount = c
			}
		}
		for _, t := range themes {
			term, _ := t["term"].(string)
			count, _ := t["count"].(int)
			weightClass := patientComponents.CalculateWeightClass(count, maxCount)
			themeVM.Themes = append(themeVM.Themes, patientComponents.ThemeItem{
				Name:        term,
				Count:       count,
				WeightClass: weightClass,
			})
			totalCount += count
		}
		themeVM.TotalCount = totalCount
	}

	// Get biopsychosocial context
	var bioContext *patientComponents.BiopsychosocialPanelViewModel
	if h.biopsychosocialService != nil {
		log.Printf("DEBUG: Getting biopsychosocial context for patient %s", id)

		ctx := r.Context()

		// Get medications
		meds, err := h.biopsychosocialService.GetMedications(ctx, id)
		log.Printf("DEBUG: GetMedications returned - err: %v, meds count: %d", err, len(meds))

		var medicationItems []patientComponents.MedicationListItemViewModel
		if err == nil && meds != nil {
			log.Printf("DEBUG: Processing %d medications", len(meds))
			for i, m := range meds {
				log.Printf("DEBUG: Medication %d - type: %T, value: %+v", i, m, m)

				// Handle different types
				switch med := m.(type) {
				case *patient.Medication:
					if med == nil {
						log.Printf("DEBUG: Medication is nil pointer")
						continue
					}
					// Convert patient.Medication to viewmodel
					statusLabel := "Desconhecido"
					switch med.Status {
					case patient.MedicationStatusActive:
						statusLabel = "Ativo"
					case patient.MedicationStatusSuspended:
						statusLabel = "Suspenso"
					case patient.MedicationStatusFinished:
						statusLabel = "Finalizado"
					}

					medicationItems = append(medicationItems, patientComponents.MedicationListItemViewModel{
						ID:          med.ID,
						Name:        med.Name,
						Dosage:      med.Dosage,
						Frequency:   med.Frequency,
						Prescriber:  med.Prescriber,
						Status:      string(med.Status),
						StatusKey:   string(med.Status),
						StatusLabel: statusLabel,
						StartedAt:   med.StartedAt.Format("02/01/2006"),
						IsActive:    med.Status == patient.MedicationStatusActive,
						IsSuspended: med.Status == patient.MedicationStatusSuspended,
						IsFinished:  med.Status == patient.MedicationStatusFinished,
					})
					log.Printf("DEBUG: Added medication from struct: %s", med.Name)

				case map[string]interface{}:
					// Keep existing conversion for backward compatibility
					status := "active"
					if s, ok := med["status"].(string); ok {
						status = s
					}

					statusLabel := "Desconhecido"
					switch status {
					case "active":
						statusLabel = "Ativo"
					case "suspended":
						statusLabel = "Suspenso"
					case "finished":
						statusLabel = "Finalizado"
					}

					medicationItems = append(medicationItems, patientComponents.MedicationListItemViewModel{
						ID:          getString(med, "id"),
						Name:        getString(med, "name"),
						Dosage:      getString(med, "dosage"),
						Frequency:   getString(med, "frequency"),
						Prescriber:  getString(med, "prescriber"),
						Status:      status,
						StatusKey:   status,
						StatusLabel: statusLabel,
						StartedAt:   getString(med, "started_at"),
						IsActive:    status == "active",
						IsSuspended: status == "suspended",
						IsFinished:  status == "finished",
					})
					log.Printf("DEBUG: Added medication from map: %s", getString(med, "name"))
				default:
					log.Printf("DEBUG: Unknown medication type: %T", m)
				}
			}
		}
		log.Printf("DEBUG: Total medication items: %d", len(medicationItems))

		// Get latest vitals
		latestVitals, errVitals := h.biopsychosocialService.GetLatestVitals(ctx, id)
		log.Printf("DEBUG: GetLatestVitals returned - err: %v, vitals: %+v", errVitals, latestVitals)

		var vitalsItem *patientComponents.VitalsItemViewModel
		if latestVitals != nil {
			log.Printf("DEBUG: Latest vitals type: %T", latestVitals)

			// Handle different types
			switch v := latestVitals.(type) {
			case *patient.Vitals:
				if v == nil {
					log.Printf("DEBUG: Latest vitals is nil pointer")
					break
				}
				// Convert patient.Vitals to viewmodel
				sleepHours := ""
				if v.SleepHours != nil {
					sleepHours = formatFloat(*v.SleepHours, 1)
				}

				appetiteLevel := ""
				if v.AppetiteLevel != nil {
					appetiteLevel = strconv.Itoa(*v.AppetiteLevel)
				}

				weight := ""
				if v.Weight != nil {
					weight = formatFloat(*v.Weight, 1)
				}

				vitalsItem = &patientComponents.VitalsItemViewModel{
					ID:               v.ID,
					Date:             v.Date.Format("02/01/2006"),
					SleepHours:       sleepHours,
					AppetiteLevel:    appetiteLevel,
					Weight:           weight,
					PhysicalActivity: strconv.Itoa(v.PhysicalActivity),
					Notes:            v.Notes,
					HasData:          true,
				}
				log.Printf("DEBUG: Created vitals item from struct: %+v", vitalsItem)

			case map[string]interface{}:
				// Keep existing conversion for backward compatibility
				log.Printf("DEBUG: Latest vitals map: %+v", v)
				vitalsItem = &patientComponents.VitalsItemViewModel{
					ID:               getString(v, "id"),
					Date:             getString(v, "date"),
					SleepHours:       getString(v, "sleep_hours"),
					AppetiteLevel:    getString(v, "appetite_level"),
					Weight:           getString(v, "weight"),
					PhysicalActivity: getString(v, "physical_activity"),
					Notes:            getString(v, "notes"),
					HasData:          true,
				}
				log.Printf("DEBUG: Created vitals item from map: %+v", vitalsItem)
			default:
				log.Printf("DEBUG: Unknown vitals type: %T", latestVitals)
			}
		}

		// Get average vitals
		avgVitals, errAvg := h.biopsychosocialService.GetAverageVitals(ctx, id, 30)
		log.Printf("DEBUG: GetAverageVitals returned - err: %v, avgVitals: %+v", errAvg, avgVitals)

		var avgItem *patientComponents.VitalsAverageItemViewModel
		if avgVitals != nil {
			log.Printf("DEBUG: Average vitals type: %T", avgVitals)

			// Handle different types
			switch a := avgVitals.(type) {
			case *patient.VitalsAverage:
				if a == nil {
					log.Printf("DEBUG: Average vitals is nil pointer")
					break
				}
				// Convert patient.VitalsAverage to viewmodel
				avgSleepHours := ""
				if a.AverageSleepHours != nil {
					avgSleepHours = formatFloat(*a.AverageSleepHours, 1)
				}

				avgAppetiteLevel := ""
				if a.AverageAppetiteLevel != nil {
					avgAppetiteLevel = formatFloat(*a.AverageAppetiteLevel, 1)
				}

				avgWeight := ""
				if a.AverageWeight != nil {
					avgWeight = formatFloat(*a.AverageWeight, 1)
				}

				avgPhysicalActivity := ""
				if a.AveragePhysicalActivity != nil {
					avgPhysicalActivity = formatFloat(*a.AveragePhysicalActivity, 0)
				}

				avgItem = &patientComponents.VitalsAverageItemViewModel{
					AvgSleepHours:       avgSleepHours,
					AvgAppetiteLevel:    avgAppetiteLevel,
					AvgWeight:           avgWeight,
					AvgPhysicalActivity: avgPhysicalActivity,
					RecordCount:         a.Count,
					HasData:             a.Count > 0,
				}
				log.Printf("DEBUG: Created average vitals item from struct: %+v", avgItem)

			case map[string]interface{}:
				// Keep existing conversion for backward compatibility
				log.Printf("DEBUG: Average vitals map: %+v", a)
				recordCount := 0
				if rc, ok := a["record_count"].(float64); ok {
					recordCount = int(rc)
				}

				avgItem = &patientComponents.VitalsAverageItemViewModel{
					AvgSleepHours:       getString(a, "avg_sleep_hours"),
					AvgAppetiteLevel:    getString(a, "avg_appetite_level"),
					AvgWeight:           getString(a, "avg_weight"),
					AvgPhysicalActivity: getString(a, "avg_physical_activity"),
					RecordCount:         recordCount,
					HasData:             recordCount > 0,
				}
				log.Printf("DEBUG: Created average vitals item from map: %+v", avgItem)
			default:
				log.Printf("DEBUG: Unknown average vitals type: %T", avgVitals)
			}
		}

		bioContext = &patientComponents.BiopsychosocialPanelViewModel{
			PatientID:     id,
			Medications:   medicationItems,
			LatestVitals:  vitalsItem,
			VitalsAverage: avgItem,
		}
		log.Printf("DEBUG: Created bioContext - Medications: %d, LatestVitals: %v, VitalsAverage: %v",
			len(medicationItems), vitalsItem != nil, avgItem != nil)
	}

	// Buscar eventos da timeline (5 mais recentes)
	var timelineEvents []patientComponents.TimelineEventItem
	if h.timelineService != nil {
		events, err := h.timelineService.GetPatientTimeline(r.Context(), id, nil, 5, 0)
		if err != nil {
			log.Printf("Erro ao buscar timeline: %v", err)
		} else {
			timelineEvents = make([]patientComponents.TimelineEventItem, 0, len(events))
			for _, event := range events {
				// Determinar ícone e cor baseado no tipo
				icon := "fa-circle"
				color := "var(--clinical-teal)"
				switch event.Type {
				case timeline.EventTypeSession:
					icon = "fa-calendar-check"
					color = "var(--primary-600)"
				case timeline.EventTypeObservation:
					icon = "fa-eye"
					color = "var(--accent-600)"
				case timeline.EventTypeIntervention:
					icon = "fa-hand-holding-medical"
					color = "var(--secondary-600)"
				}

				timelineEvents = append(timelineEvents, patientComponents.TimelineEventItem{
					ID:      event.ID,
					Type:    string(event.Type),
					Date:    event.Date.Format("02/01/2006"),
					Content: truncateText(event.Content, 50),
					Icon:    icon,
					Color:   color,
				})
			}
		}
	}

	// Map to templ components
	patientDetail := patientComponents.PatientDetailItem{
		ID:         patientData.ID,
		Name:       patientData.Name,
		Notes:      patientData.Notes,
		CreatedAt:  patientData.CreatedAt.Format("02/01/2006 às 15:04"),
		Themes:     themeVM,
		BioContext: bioContext,
		Timeline:   timelineEvents,
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

	detail := patientComponents.PatientDetail(patientDetail, sessionItems)
	layoutComponents.BaseWithContent(patientData.Name, detail).Render(r.Context(), w)
}

// Helper functions for biopsychosocial context
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func formatFloat(f float64, decimals int) string {
	return strconv.FormatFloat(f, 'f', decimals, 64)
}

// truncateText truncates text to maxLength characters, adding "..." if truncated
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

// NewPatient handles GET /patients/new - shows new patient form
func (h *PatientHandler) NewPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	formData := patientComponents.NewPatientFormData{
		FormData: &patientComponents.PatientFormValues{},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTMX-aware rendering
	if r.Header.Get("HX-Request") == "true" {
		// Render just the form fragment
		patientComponents.NewPatientForm(formData).Render(r.Context(), w)
		return
	}

	// Render with layout using templ
	form := patientComponents.NewPatientForm(formData)
	layoutComponents.BaseWithContent("Novo Paciente", form).Render(r.Context(), w)
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

// CreatePatient handles POST /patient/create - creates a new patient
func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderError(w, r, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	notes := r.FormValue("notes")

	if name == "" {
		h.renderError(w, r, "Nome do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Create patient using service
	input := services.CreatePatientInput{
		Name:  name,
		Notes: notes,
	}

	patient, err := h.patientService.CreatePatient(ctx, input)
	if err != nil {
		h.renderError(w, r, "Erro ao criar paciente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to patient detail page
	http.Redirect(w, r, "/patients/"+patient.ID, http.StatusSeeOther)
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
	id := strings.TrimPrefix(path, prefix)
	if id == "" {
		return ""
	}
	// Handle trailing slashes and get first segment
	id = strings.TrimSuffix(id, "/")
	// Get first segment if there are more path parts
	for i, c := range id {
		if c == '/' {
			return id[:i]
		}
	}
	return id
}

// ListSessions handles GET /patients/{id}/sessions - returns session list fragment via HTMX
func (h *PatientHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract patient ID from path
	id := extractIDFromPath(r.URL.Path, "/patients/")
	if id == "" {
		h.renderError(w, r, "ID do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Get sessions for the patient
	sessions, err := h.sessionService.ListSessionsByPatient(ctx, id)
	if err != nil {
		// Return empty list on error
		sessions = []*session.Session{}
	}

	// Map to view models
	sessionItems := make([]sessionComponents.SessionItem, len(sessions))
	for i, s := range sessions {
		sessionItems[i] = sessionComponents.SessionItem{
			ID:      s.ID,
			Date:    s.Date.Format("02/01/2006"),
			Summary: s.Summary,
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render just the session list fragment for HTMX
	sessionComponents.SessionList(id, sessionItems).Render(ctx, w)
}

// Search handles GET /patients/search - returns search results fragment via HTMX
func (h *PatientHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		patientComponents.SearchResults([]patientComponents.SearchResultItem{}).Render(r.Context(), w)
		return
	}

	ctx := r.Context()

	// Search patients with default limit of 15
	patients, err := h.patientService.SearchPatients(ctx, query, 15, 0)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		patientComponents.SearchResults([]patientComponents.SearchResultItem{}).Render(r.Context(), w)
		return
	}

	// Map to view models
	results := make([]patientComponents.SearchResultItem, len(patients))
	for i, p := range patients {
		results[i] = patientComponents.SearchResultItem{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render search results fragment for HTMX
	patientComponents.SearchResults(results).Render(ctx, w)
}
