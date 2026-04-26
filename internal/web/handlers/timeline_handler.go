package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"arandu/internal/domain/timeline"
	layoutComponents "arandu/web/components/layout"
	patientComponents "arandu/web/components/patient"
	"github.com/a-h/templ"
)

type TimelineService interface {
	GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error)
	SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error)
}

type TimelineHandler struct {
	timelineService TimelineService
	patientService PatientService
}

// TimelineEventViewModel é uma versão do TimelineEvent para templates
// que usa templ.Component para conteúdo HTML não escapado
type TimelineEventViewModel struct {
	ID        string
	Type      timeline.EventType
	Date      string
	Content   templ.Component
	Metadata  map[string]string
	CreatedAt string
}

type TimelineViewModel []*TimelineEventViewModel

func NewTimelineHandler(timelineService TimelineService, patientService PatientService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
		patientService: patientService,
	}
}

// rawHTML converte uma string HTML em um componente Templ sem escape
func rawHTML(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	})
}

func (h *TimelineHandler) ShowPatientHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patientID := extractPatientIDFromTimelinePath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	filterType := parseFilterType(r.URL.Query().Get("filter"))

	limit := 20
	offset := 0

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get timeline events
	events, err := h.timelineService.GetPatientTimeline(ctx, patientID, filterType, limit, offset)
	if err != nil {
		log.Printf("Error getting patient timeline: %v", err)
		http.Error(w, "Failed to load patient history", http.StatusInternalServerError)
		return
	}

	// Get patient info
	patient, err := h.patientService.GetPatientByID(ctx, patientID)
	if err != nil {
		log.Printf("Error getting patient: %v", err)
		http.Error(w, "Failed to load patient", http.StatusInternalServerError)
		return
	}

	// Build the new ViewModel
	currentFilter := r.URL.Query().Get("filter")
	if currentFilter == "" {
		currentFilter = "all"
	}

	// Age - not available in Patient struct, show therapy duration instead
	ageStr := ""
	_ = ageStr

	// Build therapy since
	since := ""
	if !patient.CreatedAt.IsZero() {
		since = "Desde " + strings.ToLower(patient.CreatedAt.Format("Jan/2006"))
	}

	// Count sessions
	sessionCount := 0
	timelineEvents := make([]patientComponents.PatientTimelineEvent, 0, len(events))
	for _, e := range events {
		isSession := string(e.Type) == "session"
		if isSession {
			sessionCount++
		}

		title := patientComponents.BuildTimelineTitle(string(e.Type), e.Metadata["title"])
		dotAccent := isSession

		dateStr := e.Date.Format("02 Jan · 2006")
		if e.Date.Year() == time.Now().Year() {
			dateStr = e.Date.Format("02 Jan")
		}

		summary := patientComponents.TruncateStr(e.Content, 120)

		href := ""
		if isSession {
			href = "/session/" + e.Metadata["session_id"]
		}

		kind := "Observação"
		kindTone := "neutral"
		switch string(e.Type) {
		case "session":
			kind = "Sessão"
			kindTone = "accent"
		case "intervention":
			kind = "Intervenção"
		}

		timelineEvents = append(timelineEvents, patientComponents.PatientTimelineEvent{
			ID:        e.ID,
			Kind:      kind,
			KindTone:  kindTone,
			DateStr:   dateStr,
			Title:     title,
			Summary:   summary,
			IsSession: isSession,
			Href:      href,
			DotAccent: dotAccent,
		})
	}

	// Recent observations (last 3)
	recentObs := make([]patientComponents.PatientRecentObs, 0)
	allEvents, _ := h.timelineService.GetPatientTimeline(ctx, patientID, nil, 1000, 0)
	obsCount := 0
	for _, e := range allEvents {
		if string(e.Type) == "observation" && obsCount < 3 {
			recentObs = append(recentObs, patientComponents.PatientRecentObs{
				Tag:     "obs",
				DateStr: e.Date.Format("02 Jan"),
				Text:   patientComponents.TruncateStr(e.Content, 160),
			})
			obsCount++
		}
	}

	// Therapy duration
	therapyDuration := "< 1 mês"
	if !patient.CreatedAt.IsZero() {
		therapyDuration = patientComponents.FormatTherapyDuration(patient.CreatedAt)
	}

	vm := patientComponents.PatientHistoryViewModel{
		PatientID:         patientID,
		PatientName:      patient.Name,
		Initials:         patientComponents.BuildInitials(patient.Name),
		AgeStr:           ageStr,
		Since:            since,
		SessionCount:      sessionCount,
		TherapyDuration:   therapyDuration,
		Frequency:        "Semanal",
		TriageContent:    patient.Notes,
		TriageDate:       patient.CreatedAt.Format("02/01/2006"),
		Events:           timelineEvents,
		CurrentFilter:   currentFilter,
		PatientIDForURL: patientID,
		RecentObservations: recentObs,
	}

	isHTMXRequest := r.Header.Get("HX-Request") == "true"

	if isHTMXRequest {
		patientComponents.PatientHistoryContent(vm).Render(ctx, w)
	} else {
		layoutComponents.Shell(layoutComponents.ShellConfig{
			PageTitle:      "Prontuário",
			ActivePage:     "patient-history",
			ShowSidebar:    true,
			SidebarVariant: "patient",
			PatientID:      patientID,
		}, patientComponents.PatientHistoryPage(vm)).Render(ctx, w)
	}
}

func extractPatientIDFromTimelinePath(path string) string {
	parts := strings.Split(path, "/")

	for i, part := range parts {
		if part == "patients" && i+1 < len(parts) {
			nextPart := parts[i+1]
			if nextPart != "" && nextPart != "history" {
				return nextPart
			}
		}
	}

	return ""
}

func parseFilterType(filterStr string) *timeline.EventType {
	switch filterStr {
	case "observation":
		filter := timeline.EventTypeObservation
		return &filter
	case "intervention":
		filter := timeline.EventTypeIntervention
		return &filter
	case "session":
		filter := timeline.EventTypeSession
		return &filter
	case "all":
		return nil
	default:
		return nil
	}
}

func (h *TimelineHandler) LoadMoreEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patientID := extractPatientIDFromTimelinePath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	filterType := parseFilterType(r.URL.Query().Get("filter"))

	limit := 20
	offset := 0

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	events, err := h.timelineService.GetPatientTimeline(ctx, patientID, filterType, limit, offset)
	if err != nil {
		log.Printf("Error loading more events: %v", err)
		http.Error(w, "Failed to load more events", http.StatusInternalServerError)
		return
	}

	data := patientComponents.TimelinePageData{
		PatientID: patientID,
		Events:    events,
		Filter:    filterType,
		Limit:     limit,
		Offset:    offset,
	}

	// Render apenas o conteúdo da timeline (sem filtros)
	patientComponents.TimelineContent(data).Render(ctx, w)
}

func (h *TimelineHandler) SearchPatientHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patientID := extractPatientIDFromTimelinePath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	query := r.URL.Query().Get("q")

	if query == "" {
		// Se a query estiver vazia, retorna a timeline normal
		h.ShowPatientHistory(w, r)
		return
	}

	results, err := h.timelineService.SearchInHistory(ctx, patientID, query)
	if err != nil {
		log.Printf("Error searching patient history: %v", err)
		http.Error(w, "Failed to search patient history", http.StatusInternalServerError)
		return
	}

	// Converter resultados para TimelineEventViewModel com conteúdo não escapado
	var viewModels TimelineViewModel
	for _, result := range results {
		viewModel := &TimelineEventViewModel{
			ID:      result.ID,
			Type:    result.Type,
			Date:    result.Date.Format("15:04"),
			Content: rawHTML(result.Snippet), // Usar helper para HTML não escapado
			Metadata: map[string]string{
				"session_id": result.SessionID,
			},
			CreatedAt: result.Date.Format("02/01/2006 às 15:04"),
		}
		viewModels = append(viewModels, viewModel)
	}

	// Renderizar resultados de busca diretamente
	renderSearchResults(w, ctx, patientID, viewModels)
}

// renderSearchResults renderiza os resultados da busca como HTML usando componentes templ
func renderSearchResults(w http.ResponseWriter, ctx context.Context, patientID string, events TimelineViewModel) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	io.WriteString(w, `<div class="timeline-container"><div class="timeline-line"></div><div id="timeline-content" class="timeline-content">`)

	for _, event := range events {
		eventTypeClass := getEventTypeCSSClass(event.Type)
		eventIcon := getEventIconClass(event.Type)
		eventLabel := getEventLabel(event.Type)

		io.WriteString(w, `<div class="timeline-event"><div class="timeline-dot timeline-dot-`+eventTypeClass+`"></div><div class="timeline-event-card timeline-event-card-`+eventTypeClass+`"><div class="timeline-event-header"><div class="flex items-center"><span class="timeline-event-type timeline-event-type-`+eventTypeClass+`"><i class="`+eventIcon+`"></i> `+eventLabel+`</span><span class="timeline-event-time">`+event.Date+`</span></div>`)
		if sessionID, ok := event.Metadata["session_id"]; ok && sessionID != "" {
			io.WriteString(w, `<a href="/session/`+sessionID+`" class="timeline-event-link"><i class="fas fa-external-link-alt"></i> Ver sessão</a>`)
		}
		io.WriteString(w, `</div><div class="timeline-event-content">`)
		event.Content.Render(ctx, w)
		io.WriteString(w, `</div><div class="timeline-event-footer"><div class="timeline-event-meta">Registrado em `+event.CreatedAt+`</div></div></div></div>`)
	}

	io.WriteString(w, `</div></div>`)
}

func getEventTypeCSSClass(eventType timeline.EventType) string {
	switch eventType {
	case timeline.EventTypeSession:
		return "session"
	case timeline.EventTypeObservation:
		return "observation"
	case timeline.EventTypeIntervention:
		return "intervention"
	default:
		return ""
	}
}

func getEventIconClass(eventType timeline.EventType) string {
	switch eventType {
	case timeline.EventTypeSession:
		return "fas fa-calendar-check"
	case timeline.EventTypeObservation:
		return "fas fa-sticky-note"
	case timeline.EventTypeIntervention:
		return "fas fa-hand-holding-heart"
	default:
		return "fas fa-calendar-day"
	}
}

func getEventLabel(eventType timeline.EventType) string {
	switch eventType {
	case timeline.EventTypeSession:
		return "Sessão"
	case timeline.EventTypeObservation:
		return "Observação"
	case timeline.EventTypeIntervention:
		return "Intervenção"
	default:
		return "Evento"
	}
}
