package handlers

import (
	"html/template"
	"log"
	"net/http"
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
	err := h.templates.ExecuteTemplate(w, contentName, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", contentName, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (h *Handler) renderSimple(w http.ResponseWriter, name string, data interface{}) {
	log.Printf("Attempting to render template: %s", name)
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	} else {
		log.Printf("Successfully rendered template: %s", name)
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
	log.Printf("handleGetPatients called")

	patients, err := h.patientService.ListPatients(r.Context())
	if err != nil {
		log.Printf("Error listing patients: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d patients", len(patients))

	// Simple test response
	html := `<html><body><h1>Pacientes</h1><p>Total: ` + strconv.Itoa(len(patients)) + `</p></body></html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Printf("Sending test response, length: %d", len(html))
	w.Write([]byte(html))
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

	// Serve a complete HTML page for new patient
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Arandu — Novo Paciente</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Source+Serif+4:ital,opsz,wght@0,8..60,400;0,8..60,600;1,8..60,400&display=swap" rel="stylesheet">
    <link href="/static/css/tailwind.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://unpkg.com/alpinejs@3.13.5/dist/cdn.min.js" defer></script>
</head>
<body class="bg-arandu-background text-arandu-text font-sans antialiased">
    <div class="flex min-h-screen">
        
        <aside class="w-20 lg:w-64 bg-white border-r border-gray-100 flex flex-col transition-all">
            <div class="p-6">
                <h1 class="text-xl font-semibold text-arandu-primary tracking-tight">Arandu</h1>
            </div>
            
            <nav class="flex-1 px-4 space-y-2">
                <a href="/dashboard" class="flex items-center p-2 rounded-md hover:bg-gray-50 text-gray-600 hover:text-arandu-primary">
                    <span class="hidden lg:inline">Dashboard</span>
                </a>
                <a href="/patients" class="flex items-center p-2 rounded-md hover:bg-gray-50 text-gray-600 hover:text-arandu-primary border-l-2 border-transparent hover:border-arandu-primary">
                    <span class="hidden lg:inline">Pacientes</span>
                </a>
            </nav>
        </aside>

        <main class="flex-1 flex flex-col lg:flex-row overflow-hidden">
            
            <section class="flex-1 overflow-y-auto p-8 lg:p-12">
                <div class="max-w-4xl mx-auto">
                    <header class="mb-12">
                        <h1 class="text-3xl font-bold text-arandu-text leading-tight">Novo Paciente</h1>
                        <p class="text-arandu-text-secondary mt-2">Cadastre um novo paciente para iniciar o acompanhamento clínico</p>
                    </header>

                    <div class="bg-white rounded-lg shadow-sm border border-gray-100">
                        <form action="/patients" method="POST" class="p-8">
                            <div class="mb-8">
                                <label for="name" class="block text-sm font-medium text-arandu-text-secondary mb-2">Nome Completo</label>
                                <input type="text" id="name" name="name" required class="block w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-md text-arandu-text placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-arandu-primary-focus focus:border-transparent transition" placeholder="Nome do paciente">
                            </div>

                            <div class="mb-8">
                                <label for="notes" class="block text-sm font-medium text-arandu-text-secondary mb-2">Observações Iniciais</label>
                                <textarea id="notes" name="notes" rows="6" placeholder="Observações relevantes sobre o paciente (histórico, queixa principal, contexto clínico, etc.)" class="block w-full px-4 py-3 bg-gray-50 border-gray-200 rounded-md text-arandu-text placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-arandu-primary-focus focus:border-transparent transition font-serif text-base"></textarea>
                            </div>

                            <div class="flex items-center justify-end space-x-4 pt-6 border-t border-gray-100">
                                <a href="/patients" class="px-6 py-2.5 rounded-md text-sm font-semibold text-arandu-text-secondary hover:bg-gray-100 transition">
                                    Cancelar
                                </a>
                                <button type="submit" class="px-6 py-2.5 rounded-md text-sm font-semibold text-white bg-arandu-primary hover:bg-arandu-primary-dark transition">
                                    Salvar Paciente
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </section>

            <aside class="w-full lg:w-72 bg-gray-50/50 border-l border-gray-100 p-6 overflow-y-auto">
                <header class="mb-6 flex items-center justify-between">
                    <h2 class="text-xs font-bold uppercase tracking-widest text-gray-400">Reflexões</h2>
                    <span class="h-2 w-2 rounded-full bg-arandu-insight animate-pulse"></span>
                </header>
                
                <div id="insights-panel" class="space-y-6 font-serif italic text-gray-600">
                    <p class="text-sm text-gray-400 font-sans not-italic">O sistema aguarda novos registros para gerar reflexões.</p>
                </div>
            </aside>

        </main>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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

	// For now, serve a simple response
	html := `<!DOCTYPE html>
<html>
<head><title>Paciente</title></head>
<body>
	<h1>Paciente: ` + patient.Name + `</h1>
	<p>ID: ` + patient.ID + `</p>
	<p>Notas: ` + patient.Notes + `</p>
	<p>Sessões: ` + string(len(sessions)) + `</p>
	<a href="/patients">Voltar</a>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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

	h.render(w, "session.html", data)
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

	h.render(w, "session_new.html", data)
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

	h.render(w, "session_edit.html", data)
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
