package agenda_test

import (
	"strings"
	"testing"

	"arandu/web/components/agenda"
)

func TestStatusLabel(t *testing.T) {
	cases := []struct {
		status string
		want   string
	}{
		{"scheduled", "Agendada"},
		{"confirmed", "Confirmada"},
		{"pending", "Aguardando confirmação"},
		{"first_session", "1ª consulta"},
		{"cancelled", "Cancelado"},
		{"no_show", "Não Compareceu"},
		{"completed", "Realizada"},
	}

	for _, tc := range cases {
		t.Run(tc.status, func(t *testing.T) {
			got := agenda.StatusLabel(tc.status)
			if got != tc.want {
				t.Errorf("StatusLabel(%q) = %q, want %q", tc.status, got, tc.want)
			}
		})
	}
}

func TestStatusLabel_UnknownFallback(t *testing.T) {
	got := agenda.StatusLabel("whatever")
	if got == "" {
		t.Error("StatusLabel should return a non-empty fallback for unknown status")
	}
}

func TestGetAppointmentCardClasses_AllStatusesReturnNonEmpty(t *testing.T) {
	statuses := []string{"confirmed", "pending", "first_session", "cancelled", "unknown"}
	for _, s := range statuses {
		t.Run(s, func(t *testing.T) {
			got := agenda.GetAppointmentCardClasses(s)
			if got == "" {
				t.Errorf("GetAppointmentCardClasses(%q) returned empty string", s)
			}
		})
	}
}

func TestViewTabClasses_ActiveTabHasPrimaryClass(t *testing.T) {
	got := agenda.ViewTabClasses("semana", "semana")
	if !strings.Contains(got, "arandu-primary") {
		t.Errorf("active tab should contain arandu-primary class, got: %q", got)
	}
}

func TestViewTabClasses_InactiveTabHasNoActiveStyle(t *testing.T) {
	got := agenda.ViewTabClasses("dia", "semana")
	if strings.Contains(got, "bg-arandu-primary") {
		t.Errorf("inactive tab should not contain bg-arandu-primary, got: %q", got)
	}
}

func TestMonthAppointmentPillClasses_AllStatusesReturnNonEmpty(t *testing.T) {
	statuses := []string{"confirmed", "pending", "first_session", "cancelled", "unknown"}
	for _, s := range statuses {
		t.Run(s, func(t *testing.T) {
			got := agenda.MonthAppointmentPillClasses(s)
			if got == "" {
				t.Errorf("MonthAppointmentPillClasses(%q) returned empty string", s)
			}
		})
	}
}

func TestTodayColumnClasses_TodayHasHighlight(t *testing.T) {
	today := agenda.TodayColumnClasses(true)
	notToday := agenda.TodayColumnClasses(false)
	if today == notToday {
		t.Error("TodayColumnClasses(true) and TodayColumnClasses(false) should differ")
	}
	if !strings.Contains(today, "arandu-primary") {
		t.Errorf("today column should highlight with arandu-primary, got: %q", today)
	}
}
