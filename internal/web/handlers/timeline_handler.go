package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func NewTimelineHandler(timelineService TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
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

	events, err := h.timelineService.GetPatientTimeline(ctx, patientID, filterType, limit, offset)
	if err != nil {
		log.Printf("Error getting patient timeline: %v", err)
		http.Error(w, "Failed to load patient history", http.StatusInternalServerError)
		return
	}

	data := patientComponents.TimelinePageData{
		PatientID: patientID,
		Events:    events,
		Filter:    filterType,
		Limit:     limit,
		Offset:    offset,
	}

	isHTMXRequest := r.Header.Get("HX-Request") == "true"

	if isHTMXRequest {
		// Para requisições HTMX de filtro, renderizar filtros + conteúdo
		// Para busca, apenas o conteúdo (já que a busca não muda os filtros)
		patientComponents.FiltersAndContent(data).Render(ctx, w)
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

// renderSearchResults renderiza os resultados da busca como HTML
func renderSearchResults(w http.ResponseWriter, ctx context.Context, patientID string, events TimelineViewModel) {
	// Escrever HTML diretamente
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Timeline container
	io.WriteString(w, `<!-- Timeline --><div style="position: relative;"><div style="position: absolute; left: 32px; top: 0; bottom: 0; width: 2px; background: var(--neutral-100);"></div><div id="timeline-content" style="display: flex; flex-direction: column; gap: var(--space-xl);">`)

	for _, event := range events {
		// Determinar estilos baseados no tipo
		badgeStyle := `background: var(--observation-bg); color: var(--observation-text);`
		dotStyle := `background: var(--observation-color);`
		borderStyle := `border-left: 4px solid var(--observation-color);`
		icon := `fas fa-sticky-note`

		if event.Type == timeline.EventTypeIntervention {
			badgeStyle = `background: var(--intervention-bg); color: var(--intervention-text);`
			dotStyle = `background: var(--intervention-color);`
			borderStyle = `border-left: 4px solid var(--intervention-color);`
			icon = `fas fa-hand-holding-heart`
		}

		// Escrever evento
		io.WriteString(w, `<div style="position: relative;"><div style="position: absolute; left: -32px; top: 24px; width: 16px; height: 16px; border-radius: var(--radius-full); border: 2px solid white; `+dotStyle+`"></div><div style="border-radius: var(--radius-lg); padding: var(--space-xl); box-shadow: var(--shadow-sm); background: white; `+borderStyle+`"><div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: var(--space-md);"><div style="display: flex; align-items: center;"><span style="display: inline-flex; align-items: center; padding: var(--space-xs) var(--space-md); border-radius: var(--radius-full); font-family: var(--font-sans); font-size: 0.75rem; font-weight: 500; `+badgeStyle+`"><i class="`+icon+`" style="margin-right: var(--space-xs);"></i> `)

		if event.Type == timeline.EventTypeObservation {
			io.WriteString(w, `Observação`)
		} else {
			io.WriteString(w, `Intervenção`)
		}

		io.WriteString(w, `</span> <span style="margin-left: var(--space-lg); font-family: var(--font-sans); font-size: 0.75rem; color: var(--neutral-500);">`+event.Date+`</span></div>`)

		// Link para sessão se houver
		if sessionID, ok := event.Metadata["session_id"]; ok && sessionID != "" {
			io.WriteString(w, `<a href="/session/`+sessionID+`" style="font-family: var(--font-sans); font-size: 0.75rem; font-weight: 500; color: var(--primary-600); text-decoration: none;" onmouseover="this.style.color='var(--primary-800)';" onmouseout="this.style.color='var(--primary-600)';"><i class="fas fa-external-link-alt" style="margin-right: 2px;"></i> Ver sessão</a>`)
		}

		io.WriteString(w, `</div><div style="font-family: var(--font-clinical); font-size: 1.125rem; color: var(--neutral-800); line-height: 1.75;">`)

		// Renderizar conteúdo (componente Templ)
		event.Content.Render(ctx, w)

		io.WriteString(w, `</div><div style="margin-top: var(--space-lg); padding-top: var(--space-md); border-top: 1px solid var(--neutral-100);"><div style="font-family: var(--font-sans); font-size: 0.75rem; color: var(--neutral-500);">Registrado em `+event.CreatedAt+`</div></div></div></div>`)
	}

	// Fechar container
	io.WriteString(w, `</div></div>`)
}
