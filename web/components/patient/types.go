package patient

import (
	"fmt"
	"strings"
	"time"
)

// TimelineEventItem representa um evento da linha do tempo no ViewModel de perfil.
type TimelineEventItem struct {
	ID      string
	Type    string // "session" | "observation" | "note"
	Date    string
	Content string
	Icon    string
	Color   string
	Title   string // used in Sábio timeline
	Href    string // link target for session events
}

// PatientHistoryViewModel — ViewModel completo da página de histórico
type PatientHistoryViewModel struct {
	PatientID         string
	PatientName      string
	Initials         string // primeiras letras de nome + sobrenome
	AgeStr           string // "34 anos"
	Since            string // "Desde jan/2026"

	// Stats (hero direita)
	SessionCount      int
	TherapyDuration   string // "3 meses" | "1 ano"
	Frequency        string // "Semanal" (fixo por ora)

	// Triagem
	TriageContent    string // patient.Notes (pode estar vazio)
	TriageDate     string // patient.CreatedAt formatada "02/01/2006"

	// Timeline
	Events          []PatientTimelineEvent
	CurrentFilter   string // "all" | "session" | "observation"
	PatientIDForURL string // usado em URLs HTMX

	// Sidebar: observações recentes (últimas 3 do tipo observation)
	RecentObservations []PatientRecentObs
}

// PatientTimelineEvent representa um evento na timeline
type PatientTimelineEvent struct {
	ID        string
	Kind      string // "Sessão" | "Observação" | "Intervenção"
	KindTone  string // "accent" (sessão) | "neutral" (outros)
	DateStr   string // "12 abr · 2026" — monospace
	Title     string // Metadata["title"] ou fallback por tipo
	Summary   string // Content truncado em 120 chars
	IsSession bool
	Href      string // "/session/{session_id}" ou ""
	DotAccent bool   // true = dot preenchido (sessão), false = neutro
}

// PatientRecentObs representa uma observação recente
type PatientRecentObs struct {
	Tag     string // tipo abreviado: "obs"
	DateStr string // "12 abr"
	Text    string // Content truncado em 160 chars
}

// SessionItem representa uma sessão resumida no ViewModel de perfil.
type SessionItem struct {
	ID            string
	SessionNumber int
	Date          string
	Summary       string
}

// AppointmentHistoryItem representa um agendamento no perfil do paciente.
type AppointmentHistoryItem struct {
	ID          string
	Date        string
	StartTime   string
	Duration   int
	StatusLabel string
	StatusClass string
	HasSession bool
	SessionID   string
}

func AppointmentStatusBadgeClass(status string) string {
	base := "inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
	switch status {
	case "scheduled":
		return base + " bg-amber-100 text-amber-800"
	case "confirmed":
		return base + " bg-emerald-100 text-emerald-800"
	case "completed":
		return base + " bg-arandu-primary/10 text-arandu-primary"
	case "cancelled":
		return base + " bg-neutral-100 text-neutral-500"
	case "no_show":
		return base + " bg-red-100 text-red-700"
	default:
		return base + " bg-neutral-100 text-neutral-600"
	}
}

// BuildInitials extrai as primeiras letras do nome
func BuildInitials(name string) string {
	if name == "" {
		return "?"
	}
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	if len(parts) == 1 {
		return strings.ToUpper(string([]rune(parts[0])[:1]))
	}
	return strings.ToUpper(string([]rune(parts[0])[:1]) + string([]rune(parts[len(parts)-1])[:1]))
}

// FormatTherapyDuration formata a duração em terapia
func FormatTherapyDuration(createdAt time.Time) string {
	months := int(time.Since(createdAt).Hours() / 730)
	if months < 1 {
		return "< 1 mês"
	}
	years := months / 12
	rem := months % 12
	if years == 0 {
		return fmt.Sprintf("%d meses", months)
	}
	if rem == 0 {
		return fmt.Sprintf("%d ano", years)
	}
	return fmt.Sprintf("%d ano e %d meses", years, rem)
}

// TruncateStr trunca uma string em n caracteres
func TruncateStr(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

// BuildTimelineTitle retorna o título do evento
func BuildTimelineTitle(eventType, title string) string {
	if title != "" {
		return title
	}
	switch eventType {
	case "session":
		return "Sessão clínica"
	case "observation":
		return "Observação"
	case "intervention":
		return "Intervenção"
	default:
		return "Evento"
	}
}
