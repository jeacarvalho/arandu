package handlers

import (
	"context"
	"net/http"

	"arandu/internal/application/services"
	"arandu/internal/domain/timeline"
	searchComponents "arandu/web/components/search"
)

// SearchHandler handles global search requests
type SearchHandler struct {
	timelineService TimelineSearchServiceInterface
}

// TimelineSearchServiceInterface defines the interface for timeline search operations
type TimelineSearchServiceInterface interface {
	SearchGlobal(ctx context.Context, query string) ([]*services.SearchGlobalResult, error)
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(timelineService TimelineSearchServiceInterface) *SearchHandler {
	return &SearchHandler{timelineService: timelineService}
}

// Search handles GET /search?q=termo
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")

	vm := searchComponents.SearchResultsViewModel{
		Query:   query,
		Results: []searchComponents.SearchResultItem{},
		Total:   0,
	}

	if len(query) < 2 {
		h.renderSearchResults(w, r, vm)
		return
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
	vm.Total = len(items)

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