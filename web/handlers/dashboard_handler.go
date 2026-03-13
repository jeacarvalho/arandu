package handlers

import (
	"net/http"
)

// DashboardHandler lida com requisições do dashboard clínico
type DashboardHandler struct {
	patientService interface{}
	sessionService interface{}
	insightService interface{}
}

// NewDashboardHandler cria um novo handler para dashboard
func NewDashboardHandler(patientService, sessionService, insightService interface{}) *DashboardHandler {
	return &DashboardHandler{
		patientService: patientService,
		sessionService: sessionService,
		insightService: insightService,
	}
}

// ServeHTTP implementa a interface http.Handler
func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Por enquanto, usamos dados mock
	// Futuramente, integrar com os serviços reais
	data := h.getDashboardData()

	// Renderizar template
	// Nota: Esta é uma implementação simplificada
	// Em produção, usaríamos o sistema de templates existente
	h.renderDashboard(w, data)
}

// getDashboardData retorna dados para o dashboard
// Por enquanto retorna dados mock, mas pode ser estendido para usar serviços reais
func (h *DashboardHandler) getDashboardData() DashboardData {
	// Usar dados mock por enquanto
	// Futuramente: integrar com patientService, sessionService, etc.
	return MockDashboardData()
}

// renderDashboard renderiza o template do dashboard
func (h *DashboardHandler) renderDashboard(w http.ResponseWriter, data DashboardData) {
	// Esta é uma implementação simplificada
	// Em produção, usaríamos o sistema de templates existente
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Por enquanto, apenas retornamos um placeholder
	// O template real será implementado separadamente
	html := `<html><body><h1>Dashboard Clínico (Em desenvolvimento)</h1>
	<p>Active Patients: ` + string(len(data.ActivePatients)) + `</p>
	<p>Recent Sessions: ` + string(len(data.RecentSessions)) + `</p>
	<p>AI Insights: ` + string(len(data.AIInsights)) + `</p>
	<p>Emerging Patterns: ` + string(len(data.EmergingPatterns)) + `</p>
	</body></html>`

	w.Write([]byte(html))
}

// Helper functions para formatar datas de forma amigável
// (Definidas no handler principal)
