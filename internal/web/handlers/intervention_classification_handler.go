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

	"arandu/internal/domain/intervention"
	interventionComponents "arandu/web/components/intervention"

	"github.com/a-h/templ"
)

// InterventionClassificationServiceInterface defines the service methods needed by the handler
type InterventionClassificationServiceInterface interface {
	GetAllInterventionTags(ctx context.Context) ([]*intervention.Tag, error)
	GetInterventionTagsByType(ctx context.Context, tagType intervention.TagType) ([]*intervention.Tag, error)
	AddTagToIntervention(ctx context.Context, interventionID, tagID string, intensity int) error
	RemoveTagFromIntervention(ctx context.Context, interventionID, tagID string) error
	GetInterventionTags(ctx context.Context, interventionID string) ([]*intervention.InterventionClassification, error)
	GetIntervention(ctx context.Context, id string) (*intervention.Intervention, error)
}

// InterventionClassificationHandler handles classification/tagging of interventions
type InterventionClassificationHandler struct {
	service InterventionClassificationServiceInterface
}

// NewInterventionClassificationHandler creates a new intervention classification handler
func NewInterventionClassificationHandler(service InterventionClassificationServiceInterface) *InterventionClassificationHandler {
	return &InterventionClassificationHandler{
		service: service,
	}
}

// ClassifyIntervention handles POST /interventions/{id}/classify
func (h *InterventionClassificationHandler) ClassifyIntervention(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	interventionID := extractInterventionIDFromPath(r.URL.Path)
	if interventionID == "" {
		http.Error(w, "ID da intervenção é obrigatório", http.StatusBadRequest)
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
			log.Printf("[ClassifyIntervention] Error parsing tags_json: %v", err)
			http.Error(w, "Erro ao processar tags: "+err.Error(), http.StatusBadRequest)
			return
		}

		for _, tag := range tags {
			log.Printf("[ClassifyIntervention] Adding tag %s with intensity %d", tag.TagID, tag.Intensity)
			if err := h.service.AddTagToIntervention(r.Context(), interventionID, tag.TagID, tag.Intensity); err != nil {
				log.Printf("[ClassifyIntervention] Error adding tag %s: %v", tag.TagID, err)
				http.Error(w, "Erro ao adicionar tag: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		tagID := r.FormValue("tag_id")
		log.Printf("[ClassifyIntervention] tag_id: %s", tagID)
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

		log.Printf("[ClassifyIntervention] Adding single tag %s with intensity %d", tagID, intensity)
		if err := h.service.AddTagToIntervention(r.Context(), interventionID, tagID, intensity); err != nil {
			log.Printf("[ClassifyIntervention] Error adding tag: %v", err)
			http.Error(w, "Erro ao adicionar tag: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tags, err := h.service.GetInterventionTags(r.Context(), interventionID)
	if err != nil {
		log.Printf("[ClassifyIntervention] Error getting intervention tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	log.Printf("[ClassifyIntervention] Found %d tags for intervention", len(tags))

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Render tags wrapped in the proper div for HTMX swap
		component := InterventionTagsWrapper(interventionID, tags)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// InterventionTagsWrapper renders the tags list wrapped in the proper div for HTMX
func InterventionTagsWrapper(interventionID string, tags []*intervention.InterventionClassification) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<div id="intervention-%s-tags" class="intervention-tags my-3 py-2 border-t border-b border-gray-100">`, interventionID)
		if err != nil {
			return err
		}
		component := interventionComponents.InterventionTagList(tags, true)
		if err := component.Render(ctx, w); err != nil {
			return err
		}
		_, err = fmt.Fprint(w, `</div>`)
		return err
	})
}

// BulkClassifyIntervention handles POST /interventions/{id}/classify/bulk
func (h *InterventionClassificationHandler) BulkClassifyIntervention(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	interventionID := extractInterventionIDFromPath(r.URL.Path)
	if interventionID == "" {
		http.Error(w, "ID da intervenção é obrigatório", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	tagIDs := r.Form["tag_ids"]

	currentTags, err := h.service.GetInterventionTags(r.Context(), interventionID)
	if err != nil {
		log.Printf("Error getting current tags: %v", err)
		http.Error(w, "Erro ao buscar tags atuais", http.StatusInternalServerError)
		return
	}

	currentTagIDs := make(map[string]bool)
	for _, t := range currentTags {
		currentTagIDs[t.TagID] = true
	}

	for _, tagID := range tagIDs {
		if !currentTagIDs[tagID] {
			if err := h.service.AddTagToIntervention(r.Context(), interventionID, tagID, 1); err != nil {
				log.Printf("Error adding tag %s: %v", tagID, err)
			}
		}
	}

	tags, err := h.service.GetInterventionTags(r.Context(), interventionID)
	if err != nil {
		log.Printf("Error getting intervention tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	component := InterventionTagsWrapper(interventionID, tags)
	component.Render(r.Context(), w)
}

// RemoveInterventionClassification handles DELETE /interventions/{id}/classify/{tag_id}
func (h *InterventionClassificationHandler) RemoveInterventionClassification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	interventionID := parts[2]
	tagID := parts[4]

	if err := h.service.RemoveTagFromIntervention(r.Context(), interventionID, tagID); err != nil {
		log.Printf("Error removing tag from intervention: %v", err)
		http.Error(w, "Erro ao remover tag", http.StatusInternalServerError)
		return
	}

	tags, err := h.service.GetInterventionTags(r.Context(), interventionID)
	if err != nil {
		log.Printf("Error getting intervention tags: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		component := InterventionTagsWrapper(interventionID, tags)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetInterventionClassificationEdit handles GET /interventions/{id}/classify/edit
func (h *InterventionClassificationHandler) GetInterventionClassificationEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	interventionID := extractInterventionIDFromPath(r.URL.Path)
	if interventionID == "" {
		http.Error(w, "ID da intervenção é obrigatório", http.StatusBadRequest)
		return
	}

	tags, err := h.service.GetAllInterventionTags(r.Context())
	if err != nil {
		log.Printf("Error getting intervention tags: %v", err)
		http.Error(w, "Erro ao buscar tags disponíveis", http.StatusInternalServerError)
		return
	}

	selectedTags, err := h.service.GetInterventionTags(r.Context(), interventionID)
	if err != nil {
		log.Printf("Error getting intervention classification tags: %v", err)
		http.Error(w, "Erro ao buscar tags da intervenção", http.StatusInternalServerError)
		return
	}

	data := interventionComponents.TagSelectorData{
		InterventionID: interventionID,
		AvailableTags:  tags,
		SelectedTags:   selectedTags,
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		component := interventionComponents.TagSelectorInline(data)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetInterventionTagsByType handles GET /tags/interventions?type={tag_type}
func (h *InterventionClassificationHandler) GetInterventionTagsByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tagTypeStr := r.URL.Query().Get("type")
	if tagTypeStr == "" {
		http.Error(w, "Tipo de tag é obrigatório", http.StatusBadRequest)
		return
	}

	tagType := intervention.TagType(tagTypeStr)
	validTypes := map[intervention.TagType]bool{
		intervention.TagTypeCognitive:       true,
		intervention.TagTypeBehavioral:      true,
		intervention.TagTypeEmotional:       true,
		intervention.TagTypePsychoeducation: true,
		intervention.TagTypeNarrative:       true,
		intervention.TagTypeBody:            true,
	}
	if !validTypes[tagType] {
		http.Error(w, "Tipo de tag inválido", http.StatusBadRequest)
		return
	}

	tags, err := h.service.GetInterventionTagsByType(r.Context(), tagType)
	if err != nil {
		log.Printf("Error getting intervention tags by type: %v", err)
		http.Error(w, "Erro ao buscar tags", http.StatusInternalServerError)
		return
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		var intTags []*intervention.InterventionClassification
		for _, tag := range tags {
			intTags = append(intTags, &intervention.InterventionClassification{
				Tag: tag,
			})
		}
		component := interventionComponents.InterventionTagList(intTags, false)
		component.Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Helper functions

func extractInterventionIDFromPath(path string) string {
	path = strings.TrimPrefix(path, "/interventions/")
	path = strings.TrimSuffix(path, "/classify/edit")
	path = strings.TrimSuffix(path, "/classify")
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
