package handlers

import (
	"net/http"
	"strings"
	"time"

	"arandu/internal/application/services"
)

type SessionHandler struct {
	createSessionService *services.CreateSessionService
}

func NewSessionHandler(createSessionService *services.CreateSessionService) *SessionHandler {
	return &SessionHandler{
		createSessionService: createSessionService,
	}
}

func (h *SessionHandler) NewSession(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/patient/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "sessions" || parts[2] != "new" {
		http.NotFound(w, r)
		return
	}

	patientID := parts[0]
	if patientID == "" {
		http.Error(w, "patient ID is required", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"PatientID": patientID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "session_new.html", data)
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	patientID := r.FormValue("patient_id")
	if patientID == "" {
		http.Error(w, "patient_id is required", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("date")
	if dateStr == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	summary := r.FormValue("summary")

	input := services.CreateSessionInput{
		PatientID: patientID,
		Date:      date,
		Summary:   summary,
	}

	session, err := h.createSessionService.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patient/"+session.PatientID, http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	// This function should be implemented in a template rendering package
	// For now, we'll use a placeholder
	w.Write([]byte("Template rendering not implemented"))
}
