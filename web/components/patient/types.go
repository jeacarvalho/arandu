package patient

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
