package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"arandu/internal/application/services"
	aiComponents "arandu/web/components/ai"

	"github.com/a-h/templ"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

func (h *AIHandler) GeneratePatientSynthesis(w http.ResponseWriter, r *http.Request) {
	// Extract patient ID from URL path
	path := r.URL.Path
	// Pattern: /patients/{id}/analysis/synthesis
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	patientID := parts[2]
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	timeframe := r.FormValue("timeframe")
	if timeframe == "" {
		timeframe = "3 months"
	}

	ctx := r.Context()

	req := services.PatientSynthesisRequest{
		PatientID:  patientID,
		Timeframe:  timeframe,
		MaxRetries: 5,
	}

	resp, err := h.aiService.GeneratePatientSynthesis(ctx, req)
	if err != nil {
		errorMsg := fmt.Sprintf("Erro ao gerar síntese: %v", err)
		renderError(w, errorMsg)
		return
	}

	generatedAt := resp.GeneratedAt.Format("02/01/2006 15:04")
	component := aiComponents.InsightCard(resp.Synthesis, generatedAt)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templ.Handler(component).ServeHTTP(w, r)
}

func renderError(w http.ResponseWriter, errorMsg string) {
	errorHTML := fmt.Sprintf(`
		<div style="background: var(--red-50); border: 1px solid var(--red-200); border-radius: var(--radius-lg); padding: var(--space-md);">
			<div style="display: flex; align-items: center;">
				<svg style="width: 20px; height: 20px; color: var(--red-400); margin-right: var(--space-sm);" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				<span style="color: var(--red-800); font-family: var(--font-sans);">%s</span>
			</div>
		</div>
	`, errorMsg)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errorHTML))
}
