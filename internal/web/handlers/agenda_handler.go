package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/appointment"

	agendaComponents "arandu/web/components/agenda"
)

// AgendaHandler handles agenda HTTP requests
type AgendaHandler struct {
	agendaService  AgendaServiceInterface
	patientService PatientServiceInterface
}

// AgendaServiceInterface defines the interface for agenda operations
type AgendaServiceInterface interface {
	CreateAppointment(ctx context.Context, patientID, patientName string, date time.Time, startTime, endTime string, duration int, apptType appointment.AppointmentType, notes string) (*appointment.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*appointment.Appointment, error)
	UpdateAppointment(ctx context.Context, id string, patientID, patientName string, date time.Time, startTime, endTime string, duration int, notes string) error
	CancelAppointment(ctx context.Context, id string) error
	CompleteAppointment(ctx context.Context, id string, sessionID string) error
	GetDayView(ctx context.Context, date time.Time) (*services.DayView, error)
	GetWeekView(ctx context.Context, date time.Time) (*services.WeekView, error)
	GetMonthView(ctx context.Context, year, month int) (*services.MonthView, error)
	GetAvailableSlots(ctx context.Context, date time.Time) ([]appointment.TimeSlot, error)
	CheckConflicts(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error)
	GetUpcomingAppointments(ctx context.Context, limit int) ([]*appointment.Appointment, error)
}

// NewAgendaHandler creates a new agenda handler
func NewAgendaHandler(agendaService AgendaServiceInterface, patientService PatientServiceInterface) *AgendaHandler {
	return &AgendaHandler{
		agendaService:  agendaService,
		patientService: patientService,
	}
}

// View handles GET /agenda - dispatches by ?view= param
func (h *AgendaHandler) View(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	view := r.URL.Query().Get("view")
	if view == "" {
		view = "semana"
	}

	dateStr := r.URL.Query().Get("date")
	var currentDate time.Time
	if dateStr != "" {
		var err error
		currentDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			currentDate = time.Now()
		}
	} else {
		currentDate = time.Now()
	}

	ctx := r.Context()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var (
		vm  agendaComponents.AgendaViewModel
		err error
	)

	switch view {
	case "dia":
		vm, err = h.viewModelForDay(ctx, currentDate)
	case "mes":
		vm, err = h.viewModelForMonth(ctx, currentDate)
	default:
		view = "semana"
		vm, err = h.viewModelForWeek(ctx, currentDate)
	}

	if err != nil {
		http.Error(w, "Failed to load agenda", http.StatusInternalServerError)
		return
	}
	vm.CurrentView = view

	if r.Header.Get("HX-Request") == "true" {
		agendaComponents.AgendaContent(vm).Render(ctx, w)
		return
	}
	agendaComponents.AgendaPage(vm).Render(ctx, w)
}

func (h *AgendaHandler) viewModelForWeek(ctx context.Context, date time.Time) (agendaComponents.AgendaViewModel, error) {
	weekView, err := h.agendaService.GetWeekView(ctx, date)
	if err != nil {
		return agendaComponents.AgendaViewModel{}, err
	}

	days := make([]agendaComponents.DayViewModel, 0, len(weekView.Days))
	total := 0
	for _, d := range weekView.Days {
		appts := convertAppointmentsForWeek(d.Appointments)
		total += len(appts)
		days = append(days, agendaComponents.DayViewModel{
			Date:         d.Date,
			DayName:      d.Date.Format("Mon"),
			DayNumber:    d.Date.Format("2"),
			IsToday:      isToday(d.Date),
			Appointments: appts,
		})
	}

	return agendaComponents.AgendaViewModel{
		WeekStart:  weekView.StartDate,
		WeekEnd:    weekView.EndDate,
		WeekLabel:  fmt.Sprintf("%d – %d de %s de %d", weekView.StartDate.Day(), weekView.EndDate.Day(), strings.ToLower(weekView.EndDate.Format("January")), weekView.EndDate.Year()),
		Days:       days,
		TotalCount: total,
		PrevDate:   weekView.StartDate.AddDate(0, 0, -7).Format("2006-01-02"),
		NextDate:   weekView.StartDate.AddDate(0, 0, 7).Format("2006-01-02"),
	}, nil
}

func (h *AgendaHandler) viewModelForDay(ctx context.Context, date time.Time) (agendaComponents.AgendaViewModel, error) {
	dayView, err := h.agendaService.GetDayView(ctx, date)
	if err != nil {
		return agendaComponents.AgendaViewModel{}, err
	}

	appts := convertAppointmentsForWeek(dayView.Appointments)
	days := []agendaComponents.DayViewModel{{
		Date:         date,
		DayName:      date.Format("Mon"),
		DayNumber:    date.Format("2"),
		IsToday:      isToday(date),
		Appointments: appts,
	}}

	return agendaComponents.AgendaViewModel{
		WeekStart:  date,
		WeekEnd:    date,
		WeekLabel:  fmt.Sprintf("%d de %s de %d", date.Day(), strings.ToLower(date.Format("January")), date.Year()),
		Days:       days,
		TotalCount: len(appts),
		PrevDate:   date.AddDate(0, 0, -1).Format("2006-01-02"),
		NextDate:   date.AddDate(0, 0, 1).Format("2006-01-02"),
	}, nil
}

func (h *AgendaHandler) viewModelForMonth(ctx context.Context, date time.Time) (agendaComponents.AgendaViewModel, error) {
	monthView, err := h.agendaService.GetMonthView(ctx, date.Year(), int(date.Month()))
	if err != nil {
		return agendaComponents.AgendaViewModel{}, err
	}

	// Index appointments by date string for O(1) lookup per day
	apptsByDate := make(map[string][]agendaComponents.AppointmentViewModel)
	for _, a := range monthView.Appointments {
		key := a.Date.Format("2006-01-02")
		apptsByDate[key] = append(apptsByDate[key], agendaComponents.AppointmentViewModel{
			ID:          a.ID,
			PatientName: a.PatientName,
			StartTime:   a.StartTime,
			EndTime:     a.EndTime,
			Status:      mapAppointmentStatus(string(a.Status)),
		})
	}

	days := make([]agendaComponents.DayViewModel, 0, len(monthView.Days))
	total := 0
	for _, d := range monthView.Days {
		key := d.Date.Format("2006-01-02")
		dayAppts := apptsByDate[key]
		total += len(dayAppts)
		days = append(days, agendaComponents.DayViewModel{
			Date:             d.Date,
			DayName:          d.Date.Format("Mon"),
			DayNumber:        d.Date.Format("2"),
			IsToday:          isToday(d.Date),
			IsCurrentMonth:   d.Date.Year() == date.Year() && d.Date.Month() == date.Month(),
			Appointments:     dayAppts,
			AppointmentCount: len(dayAppts),
		})
	}

	prev := date.AddDate(0, -1, 0)
	next := date.AddDate(0, 1, 0)
	return agendaComponents.AgendaViewModel{
		WeekStart:  date,
		WeekEnd:    date,
		WeekLabel:  fmt.Sprintf("%s de %d", strings.ToLower(date.Format("January")), date.Year()),
		Days:       days,
		TotalCount: total,
		PrevDate:   fmt.Sprintf("%d-%02d-01", prev.Year(), prev.Month()),
		NextDate:   fmt.Sprintf("%d-%02d-01", next.Year(), next.Month()),
	}, nil
}

func isToday(date time.Time) bool {
	now := time.Now()
	return date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day()
}

// convertAppointmentsForWeek converts appointments for the week view with pixel positioning
func convertAppointmentsForWeek(appointments []*appointment.Appointment) []agendaComponents.AppointmentViewModel {
	result := make([]agendaComponents.AppointmentViewModel, 0, len(appointments))
	for _, a := range appointments {
		startHour := a.GetStartDateTime().Hour()
		startMin := a.GetStartDateTime().Minute()
		topPx := ((startHour-8)*60 + startMin) * 80 / 60
		if topPx < 0 {
			topPx = 0
		}

		heightPx := (a.Duration * 80 / 60) - 4
		if heightPx < 20 {
			heightPx = 20
		}

		status := mapAppointmentStatus(string(a.Status))

		sessionType := "Sessão individual"
		if a.Type == appointment.AppointmentTypeFollowup {
			sessionType = "Acompanhamento"
		} else if a.Type == appointment.AppointmentTypeBlocked {
			sessionType = "Bloqueado"
		}

		result = append(result, agendaComponents.AppointmentViewModel{
			ID:          a.ID,
			PatientName: a.PatientName,
			StartTime:   a.StartTime,
			EndTime:     a.EndTime,
			SessionType: sessionType,
			Status:      status,
			TopPx:       topPx,
			HeightPx:    heightPx,
		})
	}
	return result
}

func mapAppointmentStatus(status string) string {
	switch status {
	case "confirmed":
		return "confirmed"
	case "scheduled":
		return "pending"
	case "completed":
		return "confirmed"
	case "cancelled":
		return "cancelled"
	case "no_show":
		return "cancelled"
	default:
		return "pending"
	}
}

// DayView handles GET /agenda/day
func (h *AgendaHandler) DayView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now()
		}
	} else {
		date = time.Now()
	}

	ctx := r.Context()
	dayView, err := h.agendaService.GetDayView(ctx, date)
	if err != nil {
		http.Error(w, "Failed to load day view", http.StatusInternalServerError)
		return
	}

	day := agendaComponents.DayViewModel{
		Date:         date,
		DayName:      date.Format("Mon"),
		DayNumber:    date.Format("2"),
		IsToday:      isToday(date),
		Appointments: convertAppointmentsForWeek(dayView.Appointments),
	}

	prevDate := date.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := date.AddDate(0, 0, 1).Format("2006-01-02")

	viewModel := agendaComponents.AgendaViewModel{
		WeekStart:   date,
		WeekEnd:     date,
		WeekLabel:   date.Format("2 de January de 2006"),
		Days:        []agendaComponents.DayViewModel{day},
		TotalCount:  len(dayView.Appointments),
		CurrentView: "dia",
		PrevDate:    prevDate,
		NextDate:    nextDate,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		agendaComponents.AgendaContent(viewModel).Render(ctx, w)
		return
	}

	agendaComponents.AgendaPage(viewModel).Render(ctx, w)
}

// WeekView handles GET /agenda/week
func (h *AgendaHandler) WeekView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now()
		}
	} else {
		date = time.Now()
	}

	ctx := r.Context()
	weekView, err := h.agendaService.GetWeekView(ctx, date)
	if err != nil {
		http.Error(w, "Failed to load week view", http.StatusInternalServerError)
		return
	}

	days := make([]agendaComponents.DayViewModel, 0, len(weekView.Days))
	for _, day := range weekView.Days {
		days = append(days, agendaComponents.DayViewModel{
			Date:         day.Date,
			DayName:      day.Date.Format("Mon"),
			DayNumber:    day.Date.Format("2"),
			IsToday:      isToday(day.Date),
			Appointments: convertAppointmentsForWeek(day.Appointments),
		})
	}

	viewModel := agendaComponents.AgendaViewModel{
		WeekStart:   weekView.StartDate,
		WeekEnd:     weekView.EndDate,
		WeekLabel:   fmt.Sprintf("%d – %d de %s", weekView.StartDate.Day(), weekView.EndDate.Day(), weekView.EndDate.Format("January")),
		Days:        days,
		CurrentView: "semana",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		agendaComponents.AgendaContent(viewModel).Render(ctx, w)
		return
	}

	agendaComponents.AgendaPage(viewModel).Render(ctx, w)
}

// MonthView handles GET /agenda/month
func (h *AgendaHandler) MonthView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	var year, month int
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}
	if monthStr != "" {
		month, _ = strconv.Atoi(monthStr)
	}

	if year == 0 || month == 0 {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	ctx := r.Context()
	monthView, err := h.agendaService.GetMonthView(ctx, year, month)
	if err != nil {
		http.Error(w, "Failed to load month view", http.StatusInternalServerError)
		return
	}

	days := make([]agendaComponents.DayViewModel, 0, len(monthView.Days))
	for _, day := range monthView.Days {
		days = append(days, agendaComponents.DayViewModel{
			Date:      day.Date,
			DayName:   day.Date.Format("Mon"),
			DayNumber: day.Date.Format("2"),
			IsToday:   isToday(day.Date),
		})
	}

	viewModel := agendaComponents.AgendaViewModel{
		WeekLabel:   fmt.Sprintf("%s %d", time.Month(month).String(), year),
		Days:        days,
		CurrentView: "mes",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		agendaComponents.AgendaContent(viewModel).Render(ctx, w)
		return
	}

	agendaComponents.AgendaPage(viewModel).Render(ctx, w)
}

// NewForm handles GET /agenda/new
func (h *AgendaHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slotStr := r.URL.Query().Get("slot")
	var slotDate time.Time
	if slotStr != "" {
		var err error
		slotDate, err = time.Parse(time.RFC3339, slotStr)
		if err != nil {
			slotDate = time.Now()
		}
	} else {
		slotDate = time.Now()
	}

	ctx := r.Context()

	patients, err := h.patientService.ListPatients(ctx)
	if err != nil {
		http.Error(w, "Failed to load patients", http.StatusInternalServerError)
		return
	}

	patientOptions := make([]agendaComponents.PatientOption, 0, len(patients))
	for _, p := range patients {
		patientOptions = append(patientOptions, agendaComponents.PatientOption{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	slots, err := h.agendaService.GetAvailableSlots(ctx, slotDate)
	if err != nil {
		slots = []appointment.TimeSlot{}
	}

	viewModel := agendaComponents.NewAppointmentFormData{
		Date:        slotDate,
		DefaultTime: fmt.Sprintf("%02d:00", slotDate.Hour()),
		Patients:    patientOptions,
		Slots:       convertTimeSlots(slots),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	agendaComponents.NewAppointmentForm(viewModel).Render(ctx, w)
}

// GetSlots handles GET /agenda/slots
func (h *AgendaHandler) GetSlots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	} else {
		date = time.Now()
	}

	ctx := r.Context()
	slots, err := h.agendaService.GetAvailableSlots(ctx, date)
	if err != nil {
		http.Error(w, "Failed to get slots", http.StatusInternalServerError)
		return
	}

	viewModel := agendaComponents.SlotsViewModel{
		Date:  date,
		Slots: convertTimeSlots(slots),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	agendaComponents.SlotsList(viewModel).Render(ctx, w)
}

// Create handles POST /agenda/appointments
func (h *AgendaHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	patientID := r.FormValue("patient_id")
	patientName := r.FormValue("patient_name")
	dateStr := r.FormValue("date")
	startTime := r.FormValue("start_time")
	durationStr := r.FormValue("duration")
	notes := r.FormValue("notes")
	apptType := r.FormValue("type")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	duration := 50
	if durationStr != "" {
		if d, err := strconv.Atoi(durationStr); err == nil {
			duration = d
		}
	}

	startHour, startMin := parseTime(startTime)
	startDateTime := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMin, 0, 0, time.Now().Location())
	endDateTime := startDateTime.Add(time.Duration(duration) * time.Minute)
	endTime := fmt.Sprintf("%02d:%02d", endDateTime.Hour(), endDateTime.Minute())

	var appointmentType appointment.AppointmentType
	switch apptType {
	case "followup":
		appointmentType = appointment.AppointmentTypeFollowup
	case "blocked":
		appointmentType = appointment.AppointmentTypeBlocked
	default:
		appointmentType = appointment.AppointmentTypeSession
	}

	ctx := r.Context()

	if patientName == "" && patientID != "" {
		patient, err := h.patientService.GetPatientByID(ctx, patientID)
		if err == nil && patient != nil {
			patientName = patient.Name
		}
	}

	_, err = h.agendaService.CreateAppointment(ctx, patientID, patientName, date, startTime, endTime, duration, appointmentType, notes)
	if err != nil {
		if strings.Contains(err.Error(), "conflicts") {
			http.Error(w, "Time slot conflicts with existing appointment", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create appointment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := "/agenda?view=day&date=" + dateStr

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// Show handles GET /agenda/appointments/{id}
func (h *AgendaHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
	if id == "" {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	appt, err := h.agendaService.GetAppointment(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get appointment", http.StatusInternalServerError)
		return
	}
	if appt == nil {
		http.NotFound(w, r)
		return
	}

	sessionType := "Sessão individual"
	if appt.Type == appointment.AppointmentTypeFollowup {
		sessionType = "Acompanhamento"
	} else if appt.Type == appointment.AppointmentTypeBlocked {
		sessionType = "Bloqueado"
	}

	viewModel := agendaComponents.AppointmentDetailModel{
		ID:          appt.ID,
		PatientID:   appt.PatientID,
		PatientName: appt.PatientName,
		Date:        appt.Date,
		StartTime:   appt.StartTime,
		EndTime:     appt.EndTime,
		Duration:    appt.Duration,
		SessionType: sessionType,
		Status:      string(appt.Status),
		Notes:       appt.Notes,
		SessionID:   appt.SessionID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	agendaComponents.AppointmentDetail(viewModel).Render(ctx, w)
}

// Update handles PUT /agenda/appointments/{id}
func (h *AgendaHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
	if id == "" {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	patientID := r.FormValue("patient_id")
	patientName := r.FormValue("patient_name")
	dateStr := r.FormValue("date")
	startTime := r.FormValue("start_time")
	durationStr := r.FormValue("duration")
	notes := r.FormValue("notes")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	duration := 50
	if durationStr != "" {
		if d, err := strconv.Atoi(durationStr); err == nil {
			duration = d
		}
	}

	startHour, startMin := parseTime(startTime)
	startDateTime := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMin, 0, 0, time.Now().Location())
	endDateTime := startDateTime.Add(time.Duration(duration) * time.Minute)
	endTime := fmt.Sprintf("%02d:%02d", endDateTime.Hour(), endDateTime.Minute())

	ctx := r.Context()

	err = h.agendaService.UpdateAppointment(ctx, id, patientID, patientName, date, startTime, endTime, duration, notes)
	if err != nil {
		if strings.Contains(err.Error(), "conflicts") {
			http.Error(w, "Time slot conflicts with existing appointment", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update appointment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Cancel handles DELETE /agenda/appointments/{id}
func (h *AgendaHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
	if id == "" {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := h.agendaService.CancelAppointment(ctx, id)
	if err != nil {
		http.Error(w, "Failed to cancel appointment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Reschedule handles POST /agenda/appointments/{id}/reschedule
func (h *AgendaHandler) Reschedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
	id = strings.TrimSuffix(id, "/reschedule")
	if id == "" {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("date")
	startTime := r.FormValue("start_time")
	durationStr := r.FormValue("duration")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	duration := 50
	if durationStr != "" {
		if d, err := strconv.Atoi(durationStr); err == nil {
			duration = d
		}
	}

	startHour, startMin := parseTime(startTime)
	startDateTime := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMin, 0, 0, time.Now().Location())
	endDateTime := startDateTime.Add(time.Duration(duration) * time.Minute)
	endTime := fmt.Sprintf("%02d:%02d", endDateTime.Hour(), endDateTime.Minute())

	ctx := r.Context()
	err = h.agendaService.UpdateAppointment(ctx, id, "", "", date, startTime, endTime, duration, "")
	if err != nil {
		if strings.Contains(err.Error(), "conflicts") {
			http.Error(w, "Time slot conflicts with existing appointment", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to reschedule appointment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/agenda", http.StatusSeeOther)
}

// Complete handles POST /agenda/appointments/{id}/complete
func (h *AgendaHandler) Complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
	id = strings.TrimSuffix(id, "/complete")
	if id == "" {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	sessionID := r.FormValue("session_id")

	ctx := r.Context()
	err := h.agendaService.CompleteAppointment(ctx, id, sessionID)
	if err != nil {
		http.Error(w, "Failed to complete appointment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Helper functions

func parseTime(timeStr string) (hour, min int) {
	parts := strings.Split(timeStr, ":")
	if len(parts) == 2 {
		hour, _ = strconv.Atoi(parts[0])
		min, _ = strconv.Atoi(parts[1])
	}
	return
}

func convertTimeSlots(slots []appointment.TimeSlot) []agendaComponents.TimeSlotViewModel {
	result := make([]agendaComponents.TimeSlotViewModel, 0, len(slots))
	for _, s := range slots {
		slot := agendaComponents.TimeSlotViewModel{
			StartTime: fmt.Sprintf("%02d:%02d", s.Start.Hour(), s.Start.Minute()),
			EndTime:   fmt.Sprintf("%02d:%02d", s.End.Hour(), s.End.Minute()),
			Available: s.Available,
		}
		if s.Appointment != nil {
			slot.AppointmentID = s.Appointment.ID
			slot.PatientName = s.Appointment.PatientName
		}
		result = append(result, slot)
	}
	return result
}
