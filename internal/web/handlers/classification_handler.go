package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"arandu/internal/domain/observation"
	"arandu/web/components/classification"

	"github.com/a-h/templ"
)

// ClassificationServiceInterface defines the service methods needed by the handler
type ClassificationServiceInterface interface {
	GetTags(ctx context.Context) ([]observation.Tag, error)
	GetTagsByType(ctx context.Context, tagType observation.TagType) ([]observation.Tag, error)
	AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error
	RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error
	GetObservationTags(ctx context.Context, observationID string) ([]observation.ObservationTag, error)
	GetObservation(ctx context.Context, id string) (*observation.Observation, error)
}

// ClassificationHandler handles classification/tagging of observations
type ClassificationHandler struct {
	service ClassificationServiceInterface
}

// NewClassificationHandler creates a new classification handler
func NewClassificationHandler(service ClassificationServiceInterface) *ClassificationHandler {
	return &ClassificationHandler{
		service: service,
	}
}

// ClassifyObservation handles POST /observations/{id}/classify
func (h *ClassificationHandler) ClassifyObservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	observationID := extractObservationIDFromPath(r.URL.Path)
	if observationID == "" {
		http.Error(w, "ID da observação é obrigatório", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	tagsJSON := r.FormValue("tags_json")

	if tagsJSON != "" && tagsJSON != "[]" {
		var tags []struct {
			TagID     string `json:"tagId"`
			TagName   string `json:"tagName"`
			TagColor  string `json:"tagColor"`
			Intensity int    `json:"intensity"`
		}
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			log.Printf("[ClassifyObservation] Error parsing tags_json: %v", err)
			http.Error(w, "Erro ao processar tags: "+err.Error(), http.StatusBadRequest)
			return
		}

		for _, tag := range tags {
			log.Printf("[ClassifyObservation] Adding tag %s with intensity %d", tag.TagID, tag.Intensity)
			if err := h.service.AddTagToObservation(r.Context(), observationID, tag.TagID, tag.Intensity); err != nil {
				log.Printf("[ClassifyObservation] Error adding tag %s: %v", tag.TagID, err)
				http.Error(w, "Erro ao adicionar tag: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		tagID := r.FormValue("tag_id")
		log.Printf("[ClassifyObservation] tag_id: %s", tagID)
		if tagID == "" {
			http.Error(w, "Tag é obrigatória", http.StatusBadRequest)
			return
		}

		intensityStr := r.FormValue("intensity")
		intensity := 3
		if intensityStr != "" {
			if i, err := strconv.Atoi(intensityStr); err == nil && i >= 1 && i <= 5 {
				intensity = i
			}
		}

		log.Printf("[ClassifyObservation] Adding single tag %s with intensity %d", tagID, intensity)
		if err := h.service.AddTagToObservation(r.Context(), observationID, tagID, intensity); err != nil {
			log.Printf("[ClassifyObservation] Error adding tag: %v", err)
			http.Error(w, "Erro ao adicionar tag: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tags, err := h.service.GetObservationTags(r.Context(), observationID)
	if err != nil {
		log.Printf("[ClassifyObservation] Error getting observation tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	log.Printf("[ClassifyObservation] Found %d tags for observation", len(tags))

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Render tags wrapped in the proper div for HTMX swap
		component := TagsWrapper(observationID, tags)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// TagsWrapper renders the tags list wrapped in the proper div for HTMX, with edit button
func TagsWrapper(observationID string, tags []observation.ObservationTag) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<div id="observation-%s-tags" class="observation-tags my-3 flex flex-wrap items-center gap-2">`, observationID)
		if err != nil {
			return err
		}
		component := classification.ObservationTagList(tags, true)
		if err := component.Render(ctx, w); err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, `<button hx-get="/observations/%s/classify/edit" hx-target="#observation-%s-tags" hx-swap="outerHTML" class="btn btn-ghost btn-xs"><i class="fas fa-tags"></i>Classificar</button></div>`, observationID, observationID)
		return err
	})
}

// RemoveClassification handles DELETE /observations/{id}/classify/{tag_id}
func (h *ClassificationHandler) RemoveClassification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	observationID := parts[2]
	tagID := parts[4]

	if err := h.service.RemoveTagFromObservation(r.Context(), observationID, tagID); err != nil {
		log.Printf("Error removing tag from observation: %v", err)
		http.Error(w, "Erro ao remover tag", http.StatusInternalServerError)
		return
	}

	tags, err := h.service.GetObservationTags(r.Context(), observationID)
	if err != nil {
		log.Printf("Error getting observation tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		component := TagsWrapper(observationID, tags)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetClassificationEdit handles GET /observations/{id}/classify/edit
func (h *ClassificationHandler) GetClassificationEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	observationID := extractObservationIDFromPath(r.URL.Path)
	if observationID == "" {
		http.Error(w, "ID da observação é obrigatório", http.StatusBadRequest)
		return
	}

	tags, err := h.service.GetTags(r.Context())
	if err != nil {
		log.Printf("Error getting tags: %v", err)
		http.Error(w, "Erro ao buscar tags disponíveis", http.StatusInternalServerError)
		return
	}

	selectedTags, err := h.service.GetObservationTags(r.Context(), observationID)
	if err != nil {
		log.Printf("Error getting observation tags: %v", err)
		http.Error(w, "Erro ao buscar tags da observação", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := classification.ObservationTagSelectorData{
			ObservationID: observationID,
			AvailableTags: tags,
			SelectedTags: selectedTags,
		}
		component := classification.ObservationTagSelectorInline(data)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// BulkClassifyObservation handles POST /observations/{id}/classify/bulk
func (h *ClassificationHandler) BulkClassifyObservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	observationID := extractObservationIDFromPath(r.URL.Path)
	if observationID == "" {
		http.Error(w, "ID da observação é obrigatório", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	tagIDs := r.Form["tag_ids"]
	log.Printf("[BulkClassifyObservation] Processing %d tag_ids for observation %s", len(tagIDs), observationID)

	currentTags, err := h.service.GetObservationTags(r.Context(), observationID)
	if err != nil {
		log.Printf("[BulkClassifyObservation] Error getting current tags: %v", err)
	}

	currentTagIDs := make(map[string]bool)
	for _, t := range currentTags {
		currentTagIDs[t.TagID] = true
	}

	availableTags, err := h.service.GetTags(r.Context())
	if err != nil {
		log.Printf("[BulkClassifyObservation] Error getting available tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	validTagIDs := make(map[string]bool)
	for _, tag := range availableTags {
		validTagIDs[tag.ID] = true
	}

	for _, tagID := range tagIDs {
		if !validTagIDs[tagID] {
			log.Printf("[BulkClassifyObservation] Invalid tag ID: %s", tagID)
			continue
		}
		if currentTagIDs[tagID] {
			log.Printf("[BulkClassifyObservation] Tag %s already exists, skipping", tagID)
			continue
		}
		if err := h.service.AddTagToObservation(r.Context(), observationID, tagID, 3); err != nil {
			log.Printf("[BulkClassifyObservation] Error adding tag %s: %v", tagID, err)
			http.Error(w, "Erro ao adicionar tag: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("[BulkClassifyObservation] Added tag %s", tagID)
	}

	tags, err := h.service.GetObservationTags(r.Context(), observationID)
	if err != nil {
		log.Printf("[BulkClassifyObservation] Error getting observation tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		component := classification.ObservationTagList(tags, true)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetTagsByType handles GET /tags?type={tag_type}
func (h *ClassificationHandler) GetTagsByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tagTypeStr := r.URL.Query().Get("type")
	if tagTypeStr == "" {
		http.Error(w, "Tipo de tag é obrigatório", http.StatusBadRequest)
		return
	}

	tagType := observation.TagType(tagTypeStr)
	if !observation.IsValidTagType(tagTypeStr) {
		http.Error(w, "Tipo de tag inválido", http.StatusBadRequest)
		return
	}

	tags, err := h.service.GetTagsByType(r.Context(), tagType)
	if err != nil {
		log.Printf("Error getting tags by type: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		var obsTags []observation.ObservationTag
		for _, tag := range tags {
			obsTags = append(obsTags, observation.ObservationTag{
				Tag: &tag,
			})
		}
		component := classification.TagList(obsTags, false)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Helper functions

func extractObservationIDFromPath(path string) string {
	path = strings.TrimPrefix(path, "/observations/")
	path = strings.TrimSuffix(path, "/classify/edit")
	path = strings.TrimSuffix(path, "/classify/bulk")
	path = strings.TrimSuffix(path, "/classify")
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
