package patient

// TimelineEventItem representa um evento da linha do tempo no ViewModel de perfil.
type TimelineEventItem struct {
	ID      string
	Type    string
	Date    string
	Content string
	Icon    string
	Color   string
}

// SessionItem representa uma sessão resumida no ViewModel de perfil.
type SessionItem struct {
	ID            string
	SessionNumber int
	Date          string
	Summary       string
}
