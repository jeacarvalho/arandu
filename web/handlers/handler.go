package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
) *Handler {
	h := &Handler{
		patientService:      patientService,
		sessionService:      sessionService,
		observationService:  observationService,
		interventionService: interventionService,
		insightService:      insightService,
	}

	h.loadTemplates()
	return h
}

func (h *Handler) loadTemplates() {
	tmpl := template.New("")

	// Parse all templates with explicit names to avoid caching issues
	templateFiles := []string{
		"web/templates/layout.html",
		"web/templates/dashboard.html",
		"web/templates/patients.html",
		"web/templates/patient.html",
		"web/templates/patient_detail.html",
		"web/templates/session.html",
		"web/templates/new_patient.html",
		"web/templates/session_new.html",
	}

	for _, file := range templateFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading template %s: %v", file, err)
			continue
		}

		_, err = tmpl.New(filepath.Base(file)).Parse(string(content))
		if err != nil {
			log.Printf("Error parsing template %s: %v", file, err)
		}
	}

	h.templates = tmpl
}

func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Usar dados mock para o dashboard clínico
	// Futuramente integrar com serviços reais
	dashboardData := MockDashboardData()

	// Estrutura de dados para o novo dashboard
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

	data := struct {
		Stats    DashboardStats
		Patients []interface{}
		Sessions []interface{}
		Insights []interface{}
	}{
		Stats: DashboardStats{
			TotalPatients:      len(dashboardData.ActivePatients),
			NewThisWeek:        2, // Placeholder - em produção viria do banco
			ActiveThisMonth:    len(dashboardData.ActivePatients),
			TotalSessions:      dashboardData.TotalSessions,
			SessionsThisWeek:   8, // Placeholder
			SessionsToday:      1, // Placeholder
			TotalInsights:      len(dashboardData.AIInsights),
			NewInsights:        3,  // Placeholder
			HighConfidence:     2,  // Placeholder
			AvgSessionDuration: 50, // Placeholder
		},
		Patients: make([]interface{}, len(dashboardData.ActivePatients)),
		Sessions: make([]interface{}, len(dashboardData.RecentSessions)),
		Insights: make([]interface{}, len(dashboardData.AIInsights)),
	}

	// Converter ActivePatients
	for i, p := range dashboardData.ActivePatients {
		data.Patients[i] = struct {
			ID        string
			Name      string
			CreatedAt time.Time
			Notes     string
		}{
			ID:        p.ID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			Notes:     p.Notes,
		}
	}

	// Converter RecentSessions
	for i, s := range dashboardData.RecentSessions {
		data.Sessions[i] = struct {
			ID            string
			PatientName   string
			Date          string
			Summary       string
			SessionNumber string
		}{
			ID:            s.ID,
			PatientName:   s.PatientName,
			Date:          s.Date.Format("02/01/2006"),
			Summary:       s.Summary,
			SessionNumber: formatSessionNumber(s.SessionNumber),
		}
	}

	// Converter AIInsights
	for i, ins := range dashboardData.AIInsights {
		data.Insights[i] = struct {
			ID         string
			Title      string
			Content    string
			Confidence string
			CreatedAt  string
		}{
			ID:         ins.ID,
			Title:      ins.Title,
			Content:    ins.Content,
			Confidence: formatConfidence(ins.Confidence),
			CreatedAt:  ins.CreatedAt.Format("02/01/2006"),
		}
	}

	h.renderTemplate(w, "dashboard", data)
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

	insights, err := h.insightService.ListInsights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Patients []interface{}
		Insights []interface{}
	}{
		Patients: make([]interface{}, len(patients)),
		Insights: make([]interface{}, len(insights)),
	}

	for i, p := range patients {
		data.Patients[i] = p
	}
	for i, ins := range insights {
		data.Insights[i] = ins
	}

	h.renderTemplate(w, "patients", data)
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

	insights, err := h.insightService.ListInsights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Insights []interface{}
	}{
		Insights: make([]interface{}, len(insights)),
	}

	for i, ins := range insights {
		data.Insights[i] = ins
	}

	h.renderTemplate(w, "new-patient", data)
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

	insights, err := h.insightService.ListInsights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Patient  interface{}
		Sessions []interface{}
		Insights []interface{}
	}{
		Patient:  patient,
		Sessions: make([]interface{}, len(sessions)),
		Insights: make([]interface{}, len(insights)),
	}

	for i, s := range sessions {
		data.Sessions[i] = s
	}
	for i, ins := range insights {
		data.Insights[i] = ins
	}

	h.renderTemplate(w, "patient-detail", data)
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

	observations, err := h.observationService.ListObservationsBySession(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	interventions, err := h.interventionService.ListInterventionsBySession(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insights, err := h.insightService.ListInsights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Session       interface{}
		Observations  []interface{}
		Interventions []interface{}
		Insights      []interface{}
	}{
		Session:       session,
		Observations:  make([]interface{}, len(observations)),
		Interventions: make([]interface{}, len(interventions)),
		Insights:      make([]interface{}, len(insights)),
	}

	for i, obs := range observations {
		data.Observations[i] = obs
	}
	for i, interv := range interventions {
		data.Interventions[i] = interv
	}
	for i, ins := range insights {
		data.Insights[i] = ins
	}

	h.renderTemplate(w, "session", data)
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

	data := struct {
		Patient interface{}
	}{
		Patient: patient,
	}

	h.renderTemplate(w, "session_new.html", data)
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

// Helper functions para o dashboard
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Hour {
		return "há menos de 1 hora"
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "há 1 hora"
		}
		return "há " + strconv.Itoa(hours) + " horas"
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "ontem"
		}
		return "há " + strconv.Itoa(days) + " dias"
	}

	return t.Format("02/01/2006")
}

func formatConfidence(c float64) string {
	return strconv.Itoa(int(c*100)) + "%"
}

func formatSessionNumber(n int) string {
	return "Sessão " + strconv.Itoa(n)
}
