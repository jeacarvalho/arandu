package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"unicode"

	appcontext "arandu/internal/platform/context"
	"arandu/internal/application/services"
	dashboardComponents "arandu/web/components/dashboard"
	layoutComponents "arandu/web/components/layout"
)

// DashboardAgendaService is a minimal interface for today's appointments.
type DashboardAgendaService interface {
	GetDayView(ctx context.Context, date time.Time) (*services.DayView, error)
}

type DashboardHandler struct {
	patientService PatientService
	sessionService SessionService
	agendaService  DashboardAgendaService
}

func NewDashboardHandler(
	patientService PatientService,
	sessionService SessionService,
	agendaService DashboardAgendaService,
) *DashboardHandler {
	return &DashboardHandler{
		patientService: patientService,
		sessionService: sessionService,
		agendaService:  agendaService,
	}
}

func (h *DashboardHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	now := time.Now()
	todayStr := now.Format("2006-01-02")
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekStartStr := weekStart.Format("2006-01-02")

	// Enriched patient list — single query, no N+1
	dashSummaries, err := h.patientService.ListForDashboard(ctx, 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	patientItems := make([]dashboardComponents.PatientItem, len(dashSummaries))
	for i, s := range dashSummaries {
		patientItems[i] = dashboardComponents.PatientItem{
			ID:               s.ID,
			Name:             s.Name,
			Tag:              s.Tag,
			SessionCount:     s.SessionCount,
			LastSessionLabel: lastSessionLabel(s.LastSessionDate, now),
			NextApptLabel:    nextApptLabel(s.NextApptDate, s.NextApptTime, now),
		}
	}

	// Session stats and recent items (kept for KPIs — ListPatients is still used)
	patients, err := h.patientService.ListPatients(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var recentItems []dashboardComponents.SessionItem
	var totalSessions, sessionsThisWeek, sessionsToday int

	for _, patient := range patients {
		sessions, err := h.sessionService.ListSessionsByPatient(ctx, patient.ID)
		if err != nil {
			continue
		}
		totalSessions += len(sessions)
		for _, sess := range sessions {
			dateStr := sess.Date.Format("2006-01-02")
			if dateStr == todayStr {
				sessionsToday++
			}
			if dateStr >= weekStartStr {
				sessionsThisWeek++
			}
			recentItems = append(recentItems, dashboardComponents.SessionItem{
				ID:          sess.ID,
				PatientName: patient.Name,
				Date:        sess.Date.Format("02/01/2006"),
				RawDate:     sess.Date,
				Summary:     sess.Summary,
			})
		}
	}

	sort.Slice(recentItems, func(i, j int) bool {
		return recentItems[i].RawDate.After(recentItems[j].RawDate)
	})
	if len(recentItems) > 6 {
		recentItems = recentItems[:6]
	}

	// Today's schedule from agenda
	var todaySchedule []dashboardComponents.AppointmentItem
	if dayView, err := h.agendaService.GetDayView(ctx, now); err == nil {
		nextMarked := false
		for _, appt := range dayView.Appointments {
			status := "upcoming"
			if appt.Status == "completed" || appt.Status == "no_show" {
				status = "done"
			} else if !nextMarked {
				// first non-done is "next"
				if appt.StartTime >= now.Format("15:04") {
					status = "next"
					nextMarked = true
				}
			}
			todaySchedule = append(todaySchedule, dashboardComponents.AppointmentItem{
				Time:        appt.StartTime,
				PatientName: appt.PatientName,
				Type:        appointmentTypeLabel(string(appt.Type)),
				Status:      status,
			})
		}
	}

	// User greeting
	userEmail, _ := appcontext.GetUserEmail(ctx)
	greeting := greetingNameFromEmail(userEmail)

	vm := dashboardComponents.DashboardVM{
		GreetingName: greeting,
		DateLabel:    ptBRDateLabel(now),
		Stats: dashboardComponents.Stats{
			TotalPatients:    len(patients),
			TotalSessions:    totalSessions,
			SessionsThisWeek: sessionsThisWeek,
			SessionsToday:    sessionsToday,
		},
		KpiItems: []dashboardComponents.KpiItem{
			{Label: "Sessões registradas", Value: fmt.Sprintf("%d", totalSessions), Delta: fmt.Sprintf("+%d esta semana", sessionsThisWeek), Tone: "neutral", Dark: true},
			{Label: "Pacientes ativos", Value: fmt.Sprintf("%d", len(patients)), Delta: "em acompanhamento", Tone: "up"},
			{Label: "Hoje", Value: fmt.Sprintf("%d", sessionsToday), Delta: fmt.Sprintf("%d agendamentos", len(todaySchedule)), Tone: "neutral"},
			{Label: "Esta semana", Value: fmt.Sprintf("%d", sessionsThisWeek), Delta: "sessões registradas", Tone: "neutral"},
		},
		TodaySchedule: todaySchedule,
		Patients:      patientItems,
		Sessions:      recentItems,
	}

	dashboard := dashboardComponents.Dashboard(vm)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX {
		dashboard.Render(r.Context(), w)
		return
	}

	layoutComponents.Shell(layoutComponents.ShellConfig{
		PageTitle:   "Dashboard",
		ActivePage:  "dashboard",
		ShowSidebar: true,
		UserEmail:   userEmail,
	}, dashboard).Render(r.Context(), w)
}

func greetingNameFromEmail(email string) string {
	if email == "" {
		return "Terapeuta"
	}
	prefix := strings.SplitN(email, "@", 2)[0]
	parts := strings.FieldsFunc(prefix, func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	})
	if len(parts) == 0 || parts[0] == "" {
		return "Terapeuta"
	}
	runes := []rune(parts[0])
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ptBRDateLabel(t time.Time) string {
	weekdays := [...]string{"Domingo", "Segunda-feira", "Terça-feira", "Quarta-feira", "Quinta-feira", "Sexta-feira", "Sábado"}
	months := [...]string{"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"}
	return strings.ToUpper(fmt.Sprintf("%s, %d de %s", weekdays[t.Weekday()], t.Day(), months[t.Month()-1]))
}

// lastSessionLabel formats "última há N dias" / "última hoje" / "" from a nullable time.
func lastSessionLabel(t *time.Time, now time.Time) string {
	if t == nil {
		return ""
	}
	days := int(now.Truncate(24*time.Hour).Sub(t.Truncate(24*time.Hour)).Hours() / 24)
	switch {
	case days == 0:
		return "última hoje"
	case days == 1:
		return "última ontem"
	case days < 30:
		return fmt.Sprintf("última há %d dias", days)
	default:
		return "última há " + fmt.Sprintf("%d", days/30) + " meses"
	}
}

// nextApptLabel formats "próxima seg, 14:00" / "próxima hoje, 14:00" from date+time strings.
func nextApptLabel(date, startTime string, now time.Time) string {
	if date == "" || startTime == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	days := int(t.Truncate(24*time.Hour).Sub(now.Truncate(24*time.Hour)).Hours() / 24)
	weekdays := [...]string{"dom", "seg", "ter", "qua", "qui", "sex", "sáb"}
	var dayLabel string
	switch days {
	case 0:
		dayLabel = "hoje"
	case 1:
		dayLabel = "amanhã"
	default:
		dayLabel = weekdays[t.Weekday()]
	}
	return fmt.Sprintf("próxima %s, %s", dayLabel, startTime)
}

func appointmentTypeLabel(typ string) string {
	switch typ {
	case "first_session":
		return "Primeira consulta"
	case "follow_up":
		return "Retorno"
	case "evaluation":
		return "Avaliação"
	default:
		return "Consulta"
	}
}
