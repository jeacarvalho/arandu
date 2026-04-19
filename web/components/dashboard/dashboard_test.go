package dashboard

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func makeSampleVM() DashboardVM {
	return DashboardVM{
		GreetingName: "Helena",
		DateLabel:    "TERÇA-FEIRA, 15 DE ABRIL",
		Stats:        Stats{TotalSessions: 10, TotalPatients: 3, SessionsThisWeek: 2, SessionsToday: 2},
		KpiItems: []KpiItem{
			{Label: "Sessões registradas", Value: "618", Delta: "+24 esta semana", Tone: "neutral", Dark: true},
			{Label: "Pacientes ativos", Value: "42", Delta: "3 novos no mês", Tone: "up"},
			{Label: "Hoje", Value: "5", Delta: "próxima às 14h", Tone: "neutral"},
			{Label: "Anotações pendentes", Value: "2", Delta: "de ontem", Tone: "warn"},
		},
		TodaySchedule: []AppointmentItem{
			{Time: "09:00", PatientName: "Amanda Rocha", Type: "Retorno", Status: "done"},
			{Time: "14:00", PatientName: "Carolina Costa", Type: "Primeira consulta", Status: "next"},
		},
		Patients: []PatientItem{
			{ID: "p001", Name: "André Barbosa", CreatedAt: "29/01/2026"},
			{ID: "p002", Name: "Amanda Rocha", CreatedAt: "15/03/2026"},
		},
		Sessions: []SessionItem{
			{ID: "s001", PatientName: "André Barbosa", Date: "02/04/2026", Summary: "Sessão produtiva.", Theme: "Sentido"},
		},
	}
}

// TestDashboard_HasEditorialHeader verifica o cabeçalho editorial com saudação.
func TestDashboard_HasEditorialHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-dash-hero") {
		t.Error("Dashboard deve ter div.sabio-dash-hero")
	}
	if !strings.Contains(html, "Helena") {
		t.Error("Hero deve conter o nome de saudação")
	}
	if !strings.Contains(html, "TERÇA-FEIRA, 15 DE ABRIL") {
		t.Error("Hero deve conter o DateLabel")
	}
}

// TestDashboard_HasKpiGrid verifica que o grid de KPIs tem 4 cards.
func TestDashboard_HasKpiGrid(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-kpi-grid") {
		t.Error("Dashboard deve ter div.sabio-kpi-grid")
	}
	count := strings.Count(html, "sabio-kpi-card")
	if count < 4 {
		t.Errorf("Dashboard deve ter 4 sabio-kpi-card, tem %d", count)
	}
}

// TestDashboard_FirstKpiIsDark verifica que o primeiro KPI usa a variante dark.
func TestDashboard_FirstKpiIsDark(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-kpi-card--dark") {
		t.Error("Primeiro KPI card deve ter classe sabio-kpi-card--dark")
	}
}

// TestDashboard_HasTodaySchedule verifica a coluna de agenda do dia.
func TestDashboard_HasTodaySchedule(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-schedule") {
		t.Error("Dashboard deve ter seção sabio-schedule")
	}
	if !strings.Contains(html, "Amanda Rocha") {
		t.Error("Agenda deve listar Amanda Rocha")
	}
	if !strings.Contains(html, "sabio-schedule-item--done") {
		t.Error("Item com status done deve ter classe sabio-schedule-item--done")
	}
	if !strings.Contains(html, "sabio-schedule-item--next") {
		t.Error("Item com status next deve ter classe sabio-schedule-item--next")
	}
}

// TestDashboard_HasPatientsSection verifica a lista de pacientes.
func TestDashboard_HasPatientsSection(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-patients-list") {
		t.Error("Dashboard deve ter div.sabio-patients-list")
	}
	if !strings.Contains(html, "André Barbosa") {
		t.Error("Lista deve conter André Barbosa")
	}
}

// TestDashboard_PatientShowsInitials verifica que cada paciente tem iniciais no avatar.
func TestDashboard_PatientShowsInitials(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-avatar") {
		t.Error("Pacientes devem ter div.sabio-avatar com iniciais")
	}
	// "André Barbosa" → "AB"
	if !strings.Contains(html, "AB") {
		t.Error("Iniciais AB devem aparecer para André Barbosa")
	}
}

// TestDashboard_HasRecentSessions verifica o card de sessões recentes.
func TestDashboard_HasRecentSessions(t *testing.T) {
	var buf bytes.Buffer
	if err := Dashboard(makeSampleVM()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-sessions-list") {
		t.Error("Dashboard deve ter div.sabio-sessions-list")
	}
	if !strings.Contains(html, "Sentido") {
		t.Error("Sessões devem exibir o Theme")
	}
}

// TestDashboard_EmptyPatientsShowsEmptyState verifica estado vazio.
func TestDashboard_EmptyPatientsShowsEmptyState(t *testing.T) {
	vm := makeSampleVM()
	vm.Patients = nil
	var buf bytes.Buffer
	if err := Dashboard(vm).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Nenhum paciente") {
		t.Error("Lista vazia deve exibir texto 'Nenhum paciente'")
	}
}

// TestDashboard_KpiUpDeltaHasIcon verifica ícone de tendência no KPI com tone "up".
func TestDashboard_KpiUpDeltaHasIcon(t *testing.T) {
	vm := makeSampleVM()
	vm.KpiItems[1].Tone = "up"
	var buf bytes.Buffer
	if err := Dashboard(vm).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-kpi-delta-icon") {
		t.Error("KPI com tone 'up' deve ter ícone sabio-kpi-delta-icon")
	}
}

// TestDashboard_KpiWarnDeltaHasIcon verifica ícone de alerta no KPI com tone "warn".
func TestDashboard_KpiWarnDeltaHasIcon(t *testing.T) {
	vm := makeSampleVM()
	vm.KpiItems[3].Tone = "warn"
	var buf bytes.Buffer
	if err := Dashboard(vm).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-kpi-delta-icon") {
		t.Error("KPI com tone 'warn' deve ter ícone sabio-kpi-delta-icon")
	}
}

// TestDashboard_PatientShowsTag verifica que a tag clínica é exibida no item de paciente.
func TestDashboard_PatientShowsTag(t *testing.T) {
	vm := makeSampleVM()
	vm.Patients = []PatientItem{
		{ID: "p001", Name: "André Barbosa", Tag: "BURNOUT", SessionCount: 5, LastSessionLabel: "última há 2 dias", NextApptLabel: "próxima qua, 10:30"},
	}
	var buf bytes.Buffer
	if err := Dashboard(vm).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "BURNOUT") {
		t.Error("Item de paciente deve exibir a tag clínica")
	}
	if !strings.Contains(html, "sabio-patient-tag") {
		t.Error("Tag deve ter classe sabio-patient-tag")
	}
}

// TestDashboard_PatientShowsSessionMeta verifica contagem de sessões e próxima consulta.
func TestDashboard_PatientShowsSessionMeta(t *testing.T) {
	vm := makeSampleVM()
	vm.Patients = []PatientItem{
		{ID: "p001", Name: "André Barbosa", Tag: "BURNOUT", SessionCount: 5, LastSessionLabel: "última há 2 dias", NextApptLabel: "próxima qua, 10:30"},
	}
	var buf bytes.Buffer
	if err := Dashboard(vm).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "5 sessões") {
		t.Error("Item de paciente deve exibir contagem de sessões")
	}
	if !strings.Contains(html, "última há 2 dias") {
		t.Error("Item de paciente deve exibir label da última sessão")
	}
	if !strings.Contains(html, "próxima qua, 10:30") {
		t.Error("Item de paciente deve exibir label da próxima consulta")
	}
}
