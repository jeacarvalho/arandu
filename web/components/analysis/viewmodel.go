package analysis

type ThemeAnalysisViewModel struct {
	PatientID      string
	PatientName    string
	Timeframe      string
	Themes         []ThemeTermViewModel
	SelectedTerm   string
	TotalCount     int
	FilteredEvents []ThemeEventViewModel
	GeneratedAt    string
}

type ThemeTermViewModel struct {
	Term      string
	Frequency int
	Weight    int
}

type ThemeEventViewModel struct {
	ID        string
	Type      string
	Date      string
	Content   string
	Icon      string
	Color     string
	SessionID string
}
