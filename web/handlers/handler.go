package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"arandu/internal/application/services"
)

type Handler struct {
	patientService      *services.PatientService
	sessionService      *services.SessionService
	observationService  *services.ObservationService
	interventionService *services.InterventionService
	insightService      *services.InsightService
	templates           *template.Template
}

func NewHandler(
	patientService *services.PatientService,
	sessionService *services.SessionService,
	observationService *services.ObservationService,
	interventionService *services.InterventionService,
	insightService *services.InsightService,
	templatePath string,
) *Handler {
	h := &Handler{
		patientService:      patientService,
		sessionService:      sessionService,
		observationService:  observationService,
		interventionService: interventionService,
		insightService:      insightService,
	}

	h.LoadTemplates(templatePath)
	return h
}

func (h *Handler) LoadTemplates(basePath string) {
	// Using a function to add custom functions to the templates
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	// Using template.ParseGlob to find and parse all .html files in the directory
	templates, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(basePath, "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	h.templates = templates
}

func (h *Handler) render(w http.ResponseWriter, contentName string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	// Execute the layout template which includes the content template
	err := h.templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", contentName, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *Handler) renderSimple(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		// More detailed error logging
		log.Printf("Error rendering template %s: %v", name, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
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

	type DashboardStats struct {
		TotalPatients      int
		NewThisWeek        int
		ActiveThisMonth    int
		TotalSessions      int
		SessionsThisWeek   int
		SessionsToday      int
		TotalInsights      int
		NewInsights        int
		HighConfidence     int
		AvgSessionDuration int
	}

	stats := DashboardStats{
		TotalPatients:      len(patients),
		NewThisWeek:        0,
		ActiveThisMonth:    len(patients),
		TotalSessions:      totalSessions,
		SessionsThisWeek:   0,
		SessionsToday:      0,
		TotalInsights:      0,
		NewInsights:        0,
		HighConfidence:     0,
		AvgSessionDuration: 0,
	}

	data := map[string]interface{}{
		"Stats":    stats,
		"Patients": patients,
		"Sessions": allSessions,
		"Insights": []interface{}{},
	}

	h.renderSimple(w, "dashboard.html", data)
}

func (h *Handler) Patients(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGetPatients(w, r)
	} else if r.Method == http.MethodPost {
		h.handleCreatePatient(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.patientService.ListPatients(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Patients": patients,
		"Insights": []interface{}{},
	}

	h.renderSimple(w, "patients.html", data)
}

func (h *Handler) handleCreatePatient(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	input := services.CreatePatientInput{
		Name:  r.FormValue("name"),
		Notes: r.FormValue("notes"),
	}

	patient, err := h.patientService.CreatePatient(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patient/"+patient.ID, http.StatusSeeOther)
}

func (h *Handler) NewPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.renderSimple(w, "new_patient.html", nil)
}

func (h *Handler) Patient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/patient/")
	if id == "" {
		http.Error(w, "Patient ID required", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetPatientByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	sessions, err := h.sessionService.ListSessionsByPatient(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Patient":  patient,
		"Sessions": make([]interface{}, len(sessions)),
	}

	for i, s := range sessions {
		data["Sessions"].([]interface{})[i] = s
	}

	h.renderSimple(w, "patient.html", data)
}

func (h *Handler) Session(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/session/")
	if id == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.GetSession(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Session": session,
	}

	h.renderSimple(w, "session.html", data)
}

func (h *Handler) NewSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	patient, err := h.patientService.GetPatientByID(r.Context(), patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Patient": patient,
	}

	h.renderSimple(w, "session_new.html", data)
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	session, err := h.sessionService.CreateSession(r.Context(), patientID, date, summary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patient/"+session.PatientID, http.StatusSeeOther)
}

func (h *Handler) EditSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/sessions/edit/")
	if id == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.GetSession(r.Context(), id)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	patient, err := h.patientService.GetPatientByID(r.Context(), session.PatientID)
	if err != nil {
		log.Printf("Error getting patient: %v", err)
		http.Error(w, "Failed to get patient", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Session": session,
		"Patient": patient,
	}

	h.renderSimple(w, "session_edit.html", data)
}

func (h *Handler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	sessionID := r.FormValue("session_id")
	if sessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
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

	input := services.UpdateSessionInput{
		ID:      sessionID,
		Date:    date,
		Summary: summary,
	}

	err = h.sessionService.UpdateSession(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.sessionService.GetSession(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patient/"+session.PatientID, http.StatusSeeOther)
}
