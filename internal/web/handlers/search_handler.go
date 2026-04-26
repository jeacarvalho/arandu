package handlers

import (
	"context"
	"net/http"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/timeline"
	searchComponents "arandu/web/components/search"
)

// SearchHandler handles global search requests
type SearchHandler struct {
	timelineService TimelineSearchServiceInterface
	patientService PatientSearchServiceInterface
}

// TimelineSearchServiceInterface defines the interface for timeline search operations
type TimelineSearchServiceInterface interface {
	SearchGlobal(ctx context.Context, query string) ([]*services.SearchGlobalResult, error)
}

// PatientSearchServiceInterface defines the interface for patient search operations
type PatientSearchServiceInterface interface {
	SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error)
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(timelineService TimelineSearchServiceInterface, patientService PatientSearchServiceInterface) *SearchHandler {
	return &SearchHandler{
		timelineService: timelineService,
		patientService: patientService,
	}
}

// Search handles GET /search?q=termo
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")

	vm := searchComponents.SearchResultsViewModel{
		Query:    query,
		Patients: []searchComponents.PatientSearchResultItem{},
		Results:  []searchComponents.SearchResultItem{},
		Total:    0,
	}

	if len(query) < 2 {
		h.renderSearchResults(w, r, vm)
		return
	}

	if h.patientService != nil {
		patients, err := h.patientService.SearchPatients(r.Context(), query, 10, 0)
		if err == nil && patients != nil {
			patientItems := make([]searchComponents.PatientSearchResultItem, 0, len(patients))
			for _, p := range patients {
				patientItems = append(patientItems, searchComponents.PatientSearchResultItem{
					ID:   p.ID,
					Name: p.Name,
				})
			}
			vm.Patients = patientItems
		}
	}

	results, err := h.timelineService.SearchGlobal(r.Context(), query)
	if err != nil || results == nil {
		h.renderSearchResults(w, r, vm)
		return
	}

	items := make([]searchComponents.SearchResultItem, 0, len(results))
	for _, res := range results {
		var itemType string
		if res.Type == timeline.EventType("observation") {
			itemType = "Observação"
		} else if res.Type == timeline.EventType("intervention") {
			itemType = "Intervenção"
		}

		items = append(items, searchComponents.SearchResultItem{
			ID:          res.ID,
			PatientID:   res.PatientID,
			PatientName: res.PatientName,
			SessionID:   res.SessionID,
			Type:        itemType,
			Date:        res.Date.Format("02 de January de 2006"),
			Snippet:     res.Snippet,
		})
	}

	vm.Results = items
	vm.Total = len(vm.Patients) + len(items)

	h.renderSearchResults(w, r, vm)
}

func (h *SearchHandler) renderSearchResults(w http.ResponseWriter, r *http.Request, vm searchComponents.SearchResultsViewModel) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		searchComponents.SearchResults(vm).Render(r.Context(), w)
		return
	}

	searchComponents.SearchResultsPage(vm).Render(r.Context(), w)
}