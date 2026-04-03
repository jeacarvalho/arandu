package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"arandu/web/components/analysis"
	layoutComponents "arandu/web/components/layout"
)

type AnalysisHandler struct {
	patientService      PatientService
	sessionService      SessionService
	observationService  ObservationService
	interventionService InterventionService
	timelineService     TimelineServicePort
}

type ObservationService interface {
	GetObservationsBySession(ctx context.Context, sessionID string) ([]interface{}, error)
}

type InterventionService interface {
	GetInterventionsBySession(ctx context.Context, sessionID string) ([]interface{}, error)
}

func NewAnalysisHandler(
	patientService PatientService,
	sessionService SessionService,
	observationService ObservationService,
	interventionService InterventionService,
	timelineService TimelineServicePort,
) *AnalysisHandler {
	return &AnalysisHandler{
		patientService:      patientService,
		sessionService:      sessionService,
		observationService:  observationService,
		interventionService: interventionService,
		timelineService:     timelineService,
	}
}

func (h *AnalysisHandler) ShowThemes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	patientID := extractIDFromPath(r.URL.Path, "/patients/")
	if patientID == "" {
		h.renderAnalysisError(w, r, "ID do paciente é obrigatório", http.StatusBadRequest)
		return
	}

	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		timeframe = "all"
	}

	patient, err := h.patientService.GetPatientByID(r.Context(), patientID)
	if err != nil {
		h.renderAnalysisError(w, r, "Paciente não encontrado", http.StatusNotFound)
		return
	}

	themes, err := h.patientService.GetThemeFrequency(r.Context(), patientID, 50)
	if err != nil {
		themes = []map[string]interface{}{}
	}

	vm := h.buildThemeViewModel(patientID, patient.Name, timeframe, themes, "")

	isHTMX := r.Header.Get("HX-Request") == "true"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if isHTMX {
		analysis.ThemeCloudPanel(vm).Render(r.Context(), w)
		return
	}

	layoutComponents.Shell(layoutComponents.ShellConfig{
		PageTitle:   "Análise de Temas - " + patient.Name,
		ActivePage:  "patient-analysis",
		ShowSidebar: true,
		PatientID:   patientID,
	}, analysis.AnalysisPage(vm)).Render(r.Context(), w)
}

func (h *AnalysisHandler) FilterTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	patientID := parts[2]
	term := parts[5]

	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		timeframe = "all"
	}

	patient, err := h.patientService.GetPatientByID(r.Context(), patientID)
	if err != nil {
		h.renderAnalysisError(w, r, "Paciente não encontrado", http.StatusNotFound)
		return
	}

	themes, _ := h.patientService.GetThemeFrequency(r.Context(), patientID, 50)
	vm := h.buildThemeViewModel(patientID, patient.Name, timeframe, themes, term)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	analysis.AnalysisPage(vm).Render(r.Context(), w)
}

func (h *AnalysisHandler) GetTimelineByTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	patientID := parts[2]
	term := parts[5]

	events := h.getTimelineFilteredByTerm(r.Context(), patientID, term)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	analysis.TimelineFiltered(events).Render(r.Context(), w)
}

func (h *AnalysisHandler) buildThemeViewModel(patientID, patientName, timeframe string, themes []map[string]interface{}, selectedTerm string) analysis.ThemeAnalysisViewModel {
	var themeVMs []analysis.ThemeTermViewModel
	totalCount := 0
	maxFreq := 0

	for _, t := range themes {
		if freq, ok := t["count"].(int); ok {
			if freq > maxFreq {
				maxFreq = freq
			}
			totalCount += freq
		}
	}

	for _, t := range themes {
		term, _ := t["term"].(string)
		freq, _ := t["count"].(int)
		weight := calculateWeight(freq, maxFreq)

		themeVMs = append(themeVMs, analysis.ThemeTermViewModel{
			Term:      term,
			Frequency: freq,
			Weight:    weight,
		})
	}

	var filteredEvents []analysis.ThemeEventViewModel
	if selectedTerm != "" {
		filteredEvents = h.getTimelineFilteredByTerm(nil, patientID, selectedTerm)
	}

	return analysis.ThemeAnalysisViewModel{
		PatientID:      patientID,
		PatientName:    patientName,
		Timeframe:      timeframe,
		Themes:         themeVMs,
		SelectedTerm:   selectedTerm,
		TotalCount:     totalCount,
		FilteredEvents: filteredEvents,
		GeneratedAt:    time.Now().Format("02/01/2006 15:04"),
	}
}

func (h *AnalysisHandler) getTimelineFilteredByTerm(ctx context.Context, patientID, term string) []analysis.ThemeEventViewModel {
	sessions, err := h.sessionService.ListSessionsByPatient(ctx, patientID)
	if err != nil {
		return []analysis.ThemeEventViewModel{}
	}

	var events []analysis.ThemeEventViewModel

	for _, sess := range sessions {
		hasTerm := false

		observations, _ := h.observationService.GetObservationsBySession(ctx, sess.ID)
		for _, obs := range observations {
			if obsMap, ok := obs.(map[string]interface{}); ok {
				if content, ok := obsMap["content"].(string); ok {
					if strings.Contains(strings.ToLower(content), strings.ToLower(term)) {
						hasTerm = true
						events = append(events, analysis.ThemeEventViewModel{
							ID:        obsMap["id"].(string),
							Type:      "observation",
							Date:      sess.Date.Format("02/01/2006"),
							Content:   truncateForEvent(content),
							Icon:      "fa-eye",
							Color:     "#6366F1",
							SessionID: sess.ID,
						})
					}
				}
			}
		}

		if !hasTerm {
			interventions, _ := h.interventionService.GetInterventionsBySession(ctx, sess.ID)
			for _, inter := range interventions {
				if intMap, ok := inter.(map[string]interface{}); ok {
					if content, ok := intMap["content"].(string); ok {
						if strings.Contains(strings.ToLower(content), strings.ToLower(term)) {
							events = append(events, analysis.ThemeEventViewModel{
								ID:        intMap["id"].(string),
								Type:      "intervention",
								Date:      sess.Date.Format("02/01/2006"),
								Content:   truncateForEvent(content),
								Icon:      "fa-hand-holding-medical",
								Color:     "#8B5CF6",
								SessionID: sess.ID,
							})
						}
					}
				}
			}
		}
	}

	return events
}

func calculateWeight(freq, maxFreq int) int {
	if maxFreq == 0 || freq == 0 {
		return 1
	}
	ratio := float64(freq) / float64(maxFreq)
	if ratio >= 0.8 {
		return 5
	} else if ratio >= 0.6 {
		return 4
	} else if ratio >= 0.4 {
		return 3
	} else if ratio >= 0.2 {
		return 2
	}
	return 1
}

func (h *AnalysisHandler) renderAnalysisError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	isHTMX := r.Header.Get("HX-Request") == "true"
	errorData := layoutComponents.ErrorData{
		Error:  message,
		IsHTMX: isHTMX,
		Title:  "Erro na Análise",
	}

	if isHTMX {
		layoutComponents.ErrorFragment(errorData).Render(r.Context(), w)
		return
	}

	layoutComponents.Shell(layoutComponents.ShellConfig{
		PageTitle:   "Erro",
		ShowSidebar: true,
	}, layoutComponents.ErrorFragment(errorData)).Render(r.Context(), w)
}

func truncateForEvent(text string) string {
	if len(text) <= 100 {
		return text
	}
	return text[:100] + "..."
}
