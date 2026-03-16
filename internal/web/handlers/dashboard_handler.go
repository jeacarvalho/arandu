package handlers

import (
	"net/http"

	dashboardComponents "arandu/web/components/dashboard"
	layoutComponents "arandu/web/components/layout"
)

type DashboardHandler struct {
	patientService PatientService
	sessionService SessionService
}

func NewDashboardHandler(
	patientService PatientService,
	sessionService SessionService,
) *DashboardHandler {
	return &DashboardHandler{
		patientService: patientService,
		sessionService: sessionService,
	}
}

func (h *DashboardHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	patients, err := h.patientService.ListPatients(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var allSessions []interface{}
	var totalSessions int
	for _, patient := range patients {
		sessions, err := h.sessionService.ListSessionsByPatient(ctx, patient.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		totalSessions += len(sessions)

		for i, session := range sessions {
			if i >= 5 {
				break
			}
			allSessions = append(allSessions, struct {
				ID          string
				PatientName string
				Date        string
				Summary     string
			}{
				ID:          session.ID,
				PatientName: patient.Name,
				Date:        session.Date.Format("02/01/2006"),
				Summary:     session.Summary,
			})
		}
	}

	stats := dashboardComponents.Stats{
		TotalPatients:   len(patients),
		NewThisWeek:     0,
		ActiveThisMonth: len(patients),
		TotalSessions:   totalSessions,
	}

	patientItems := make([]dashboardComponents.PatientItem, len(patients))
	for i, p := range patients {
		patientItems[i] = dashboardComponents.PatientItem{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt.Format("02/01/2006"),
		}
	}

	sessionItems := make([]dashboardComponents.SessionItem, len(allSessions))
	for i, s := range allSessions {
		if s, ok := s.(struct {
			ID          string
			PatientName string
			Date        string
			Summary     string
		}); ok {
			sessionItems[i] = dashboardComponents.SessionItem{
				ID:          s.ID,
				PatientName: s.PatientName,
				Date:        s.Date,
				Summary:     s.Summary,
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	dashboard := dashboardComponents.Dashboard(stats, patientItems, sessionItems)
	layoutComponents.BaseWithContent("Dashboard", dashboard).Render(r.Context(), w)
}
