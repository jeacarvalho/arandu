package patient

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	domainPatient "arandu/internal/domain/patient"
)

func newSabioTestPatient() *domainPatient.Patient {
	return &domainPatient.Patient{
		ID:         "p0446",
		Name:       "André Barbosa",
		Gender:     "m",
		Ethnicity:  "b",
		Occupation: "Engenheiro",
		Education:  "s",
		Notes:      "Burnout profissional com sintomas somáticos. Busca reencontro com sentido profissional.",
		CreatedAt:  time.Now().Add(-760 * 24 * time.Hour), // ~2 anos
		UpdatedAt:  time.Now(),
	}
}

func newSabioTimeline() []TimelineEventItem {
	return []TimelineEventItem{
		{ID: "s1", Type: "session", Date: "02/04/2026", Content: "Questionamento existencial.", Title: "Sessão 5 · Sentido e existência", Href: "/session/s1"},
		{ID: "n1", Type: "note", Date: "19/03/2026", Content: "Melhora no sono relatada.", Title: "Nota entre sessões"},
		{ID: "o1", Type: "observation", Date: "12/03/2026", Content: "Questiona escolhas profissionais."},
	}
}

func newSabioSessions() []SessionItem {
	return []SessionItem{
		{ID: "s1", SessionNumber: 5, Date: "02/04/2026", Summary: "Exploração de questionamentos existenciais."},
		{ID: "s2", SessionNumber: 4, Date: "12/03/2026", Summary: "Histórico profissional."},
	}
}

// TestPatientProfile_HasSabioHero verifica o hero editorial Sábio.
func TestPatientProfile_HasSabioHero(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, newSabioTimeline(), newSabioSessions()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-patient-hero") {
		t.Error("Perfil deve ter div.sabio-patient-hero")
	}
	if !strings.Contains(html, "André Barbosa") {
		t.Error("Hero deve conter o nome do paciente")
	}
}

// TestPatientProfile_HeroHasAvatarWithInitials verifica avatar com iniciais no hero.
func TestPatientProfile_HeroHasAvatarWithInitials(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, nil, nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-avatar") {
		t.Error("Hero deve ter div.sabio-avatar")
	}
	if !strings.Contains(html, "AB") {
		t.Error("Avatar deve conter as iniciais 'AB'")
	}
}

// TestPatientProfile_HeroHasStatBlocks verifica os stat blocks (sessões, tempo).
func TestPatientProfile_HeroHasStatBlocks(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, nil, newSabioSessions()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-stat-block") {
		t.Error("Hero deve ter sabio-stat-block para métricas")
	}
	if !strings.Contains(html, "em terapia") {
		t.Error("Deve exibir 'em terapia' como label de stat")
	}
}

// TestPatientProfile_HasTriageBlockquote verifica notas de triagem como blockquote editorial.
func TestPatientProfile_HasTriageBlockquote(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, nil, nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-triage-quote") {
		t.Error("Perfil deve ter sabio-triage-quote com notas de triagem")
	}
	if !strings.Contains(html, "Burnout profissional") {
		t.Error("Blockquote deve conter o conteúdo das notas de triagem")
	}
}

// TestPatientProfile_HasTimeline verifica a linha do tempo clínica.
func TestPatientProfile_HasTimeline(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, newSabioTimeline(), nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-timeline") {
		t.Error("Perfil deve ter div.sabio-timeline")
	}
	if !strings.Contains(html, "Sessão 5") {
		t.Error("Timeline deve exibir o título do evento")
	}
}

// TestPatientProfile_HasQuickActions verifica coluna de ações rápidas.
func TestPatientProfile_HasQuickActions(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, nil, nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-action-btn") {
		t.Error("Perfil deve ter sabio-action-btn para ações rápidas")
	}
	if !strings.Contains(html, "Nova sessão") {
		t.Error("Ações rápidas devem ter 'Nova sessão'")
	}
}

// TestPatientProfile_NewSessionIsPrimary verifica que "Nova sessão" usa variante primária.
func TestPatientProfile_NewSessionIsPrimary(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, nil, nil).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-action-btn--primary") {
		t.Error("Botão 'Nova sessão' deve ter classe sabio-action-btn--primary")
	}
}

// TestPatientProfile_RecentSessionsList verifica lista de sessões recentes.
func TestPatientProfile_RecentSessionsList(t *testing.T) {
	p := newSabioTestPatient()
	var buf bytes.Buffer
	if err := PatientProfileView(p, nil, nil, newSabioSessions()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-session-timeline") {
		t.Error("Perfil deve ter div.sabio-session-timeline para sessões")
	}
	if !strings.Contains(html, "Sessão 5") {
		t.Error("Lista de sessões deve conter 'Sessão 5'")
	}
}
