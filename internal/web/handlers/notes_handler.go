package handlers

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"arandu/internal/domain/timeline"
	"arandu/web/components/layout"
	notesComponents "arandu/web/components/notes"
)

type NotesHandler struct {
	patientService  PatientService
	timelineService TimelineService
}

func NewNotesHandler(ps PatientService, ts TimelineService) *NotesHandler {
	return &NotesHandler{
		patientService:   ps,
		timelineService: ts,
	}
}

func (h *NotesHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	focusedID := r.URL.Query().Get("patient")
	filter := r.URL.Query().Get("filter")
	searchQuery := r.URL.Query().Get("q")
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	pageSize := 50

	patients, total, err := h.patientService.ListPatientsPaginated(ctx, page, pageSize)
	if err != nil {
		log.Printf("Error listing patients: %v", err)
		http.Error(w, "Failed to load patients", http.StatusInternalServerError)
		return
	}

	records := make([]notesComponents.NoteRecord, 0, len(patients))
	tagSet := make(map[string]bool)
	var totalEvents int

	for _, p := range patients {
		if filter != "" && filter != "all" && p.Tag != filter {
			continue
		}

		if searchQuery != "" {
			lowerName := strings.ToLower(p.Name)
			lowerQuery := strings.ToLower(searchQuery)
			if !strings.Contains(lowerName, lowerQuery) {
				continue
			}
		}

		if p.Tag != "" {
			tagSet[p.Tag] = true
		}

		tl, _ := h.timelineService.GetPatientTimeline(ctx, p.ID, nil, 1000, 0)
		eventCount := len(tl)
		totalEvents += eventCount

		lastUpdate := ""
		if !p.UpdatedAt.IsZero() {
			lastUpdate = p.UpdatedAt.Format("02/01/2006")
		}

		tags := []string{}
		if p.Tag != "" {
			tags = []string{p.Tag}
		}

		isFocused := p.ID == focusedID

		records = append(records, notesComponents.NoteRecord{
			PatientID:   p.ID,
			PatientName: p.Name,
			Initials:   notesComponents.BuildNotesInitials(p.Name),
			RecordID:  notesComponents.BuildRecordID(p.ID),
			EventCount: eventCount,
			LastUpdate: lastUpdate,
			Tags:     tags,
			Status:   "EM ACOMPANHAMENTO",
			IsFocused: isFocused,
		})
	}

	for i := range records {
		if records[i].PatientID == focusedID {
			records[i].IsFocused = true
		} else if focusedID == "" && i == 0 {
			records[i].IsFocused = true
		}
	}

	filterTags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		filterTags = append(filterTags, tag)
	}
	sort.Strings(filterTags)

	var detail notesComponents.NoteRecordDetail
	for _, rec := range records {
		if rec.IsFocused {
			detail = h.buildNoteRecordDetail(ctx, rec.PatientID, "evolucao")
			break
		}
	}

	currentFilter := filter
	if currentFilter == "" {
		currentFilter = "all"
	}

	vm := notesComponents.NotesLibraryViewModel{
		Records:        records,
		FocusedRecord: detail,
		FilterTags:    filterTags,
		TotalRecords: total,
		TotalEvents:  totalEvents,
		CurrentFilter: currentFilter,
		CurrentPage: page,
		PageSize:   pageSize,
		SearchQuery: searchQuery,
	}

	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		notesComponents.NotesLibraryContent(vm).Render(ctx, w)
		return
	}

	layout.Shell(layout.ShellConfig{
		PageTitle:  "Prontuários",
		ActivePage: "notes",
	}, notesComponents.NotesLibraryPage(vm)).Render(ctx, w)
}

func (h *NotesHandler) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patientID := extractNotesPatientID(r.URL.Path, "/notes/detail/")
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	focusedPatientID := r.URL.Query().Get("patient")
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		tab = "evolucao"
	}

	detail := h.buildNoteRecordDetail(ctx, patientID, tab)
	detail.FocusedPatientID = focusedPatientID

	notesComponents.NotesDetailPanel(detail).Render(ctx, w)
}

func (h *NotesHandler) buildNoteRecordDetail(ctx context.Context, patientID, tab string) notesComponents.NoteRecordDetail {
	p, err := h.patientService.GetPatientByID(ctx, patientID)
	if err != nil || p == nil {
		return notesComponents.NoteRecordDetail{}
	}

	tl, _ := h.timelineService.GetPatientTimeline(ctx, patientID, nil, 20, 0)

	var sessions, observations, interventions []*timeline.TimelineEvent
	for _, e := range tl {
		switch e.Type {
		case timeline.EventTypeSession:
			sessions = append(sessions, e)
		case timeline.EventTypeObservation:
			observations = append(observations, e)
		case timeline.EventTypeIntervention:
			interventions = append(interventions, e)
		}
	}

	detail := notesComponents.NoteRecordDetail{
		PatientID:   p.ID,
		PatientName: p.Name,
		RecordID:  notesComponents.BuildRecordID(p.ID),
		ActiveTab: tab,
		Sections: []notesComponents.NoteSection{
			{Key: "evolucao", Label: "Evolução", Count: len(sessions), Active: tab == "evolucao"},
			{Key: "anamnese", Label: "Anamnese", Count: 0, Active: tab == "anamnese"},
			{Key: "observacoes", Label: "Observações", Count: len(observations), Active: tab == "observacoes"},
			{Key: "intervencoes", Label: "Intervenções", Count: len(interventions), Active: tab == "intervencoes"},
		},
		SessionCount:      len(sessions),
		ObservationCount: len(observations),
		InterventionCount: len(interventions),
		TotalEvents:      len(tl),
	}

	if tab == "evolucao" && len(sessions) > 0 {
		last := sessions[0]
		detail.LastEntryDate = last.Date.Format("02/01/2006")
		detail.LastEntryTitle = notesComponents.BuildEventTitle(string(last.Type), last.Metadata["title"])
		detail.LastEntryContent = notesComponents.TruncateStr(last.Content, 400)

		for _, e := range tl {
			if e.Type == timeline.EventTypeObservation {
				detail.LastEntryQuote = notesComponents.TruncateStr(e.Content, 200)
				break
			}
		}
	}

	detail.ObservationItems = make([]notesComponents.NoteEventItem, 0, len(observations))
	for _, e := range observations {
		detail.ObservationItems = append(detail.ObservationItems, notesComponents.NoteEventItem{
			DateStr: e.Date.Format("02 Jan · 2006"),
			Content: notesComponents.TruncateStr(e.Content, 300),
		})
	}

	detail.InterventionItems = make([]notesComponents.NoteEventItem, 0, len(interventions))
	for _, e := range interventions {
		detail.InterventionItems = append(detail.InterventionItems, notesComponents.NoteEventItem{
			DateStr: e.Date.Format("02 Jan · 2006"),
			Content: notesComponents.TruncateStr(e.Content, 300),
		})
	}

	detail.AnamneseText = p.Notes

	return detail
}

// Search handles GET /notes/search?q=... — retorna apenas os itens da lista (HTMX)
func (h *NotesHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

	patients, err := h.patientService.ListPatients(ctx)
	if err != nil {
		http.Error(w, "Failed to search patients", http.StatusInternalServerError)
		return
	}

	records := make([]notesComponents.NoteRecord, 0)
	for _, p := range patients {
		if q != "" && !strings.Contains(strings.ToLower(p.Name), q) {
			continue
		}
		tags := []string{}
		if p.Tag != "" {
			tags = []string{p.Tag}
		}
		lastUpdate := ""
		if !p.UpdatedAt.IsZero() {
			lastUpdate = p.UpdatedAt.Format("02/01/2006")
		}
		records = append(records, notesComponents.NoteRecord{
			PatientID:   p.ID,
			PatientName: p.Name,
			Initials:    notesComponents.BuildNotesInitials(p.Name),
			RecordID:    notesComponents.BuildRecordID(p.ID),
			LastUpdate:  lastUpdate,
			Tags:        tags,
			Status:      "EM ACOMPANHAMENTO",
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	notesComponents.NoteRecordItems(records).Render(ctx, w)
}

func extractNotesPatientID(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	id := strings.TrimPrefix(path, prefix)
	if id == "" {
		return ""
	}
	return id
}