package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"arandu/internal/domain/timeline"
	layoutComponents "arandu/web/components/layout"
	patientComponents "arandu/web/components/patient"
)

type TimelineService interface {
	GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType) (timeline.Timeline, error)
}

type TimelineHandler struct {
	timelineService TimelineService
}

func NewTimelineHandler(timelineService TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
	}
}

func (h *TimelineHandler) ShowPatientHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	patientID := extractPatientIDFromTimelinePath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	filterType := parseFilterType(r.URL.Query().Get("filter"))

	events, err := h.timelineService.GetPatientTimeline(ctx, patientID, filterType)
	if err != nil {
		log.Printf("Error getting patient timeline: %v", err)
		http.Error(w, "Failed to load patient history", http.StatusInternalServerError)
		return
	}

	data := patientComponents.TimelinePageData{
		PatientID: patientID,
		Events:    events,
		Filter:    filterType,
	}

	isHTMXRequest := r.Header.Get("HX-Request") == "true"

	if isHTMXRequest {
		// Para requisições HTMX, renderizar apenas o conteúdo da timeline
		patientComponents.TimelineContainer(data).Render(ctx, w)
	} else {
		layoutComponents.BaseWithContent("Prontuário", patientComponents.TimelineContainer(data)).Render(ctx, w)
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
	case "session":
		filter := timeline.EventTypeSession
		return &filter
	case "observation":
		filter := timeline.EventTypeObservation
		return &filter
	case "intervention":
		filter := timeline.EventTypeIntervention
		return &filter
	case "all":
		return nil
	default:
		return nil
	}
}
