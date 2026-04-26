package search

// PatientSearchResultItem represents a patient search result
type PatientSearchResultItem struct {
	ID   string
	Name string
}

// SearchResultItem represents a single search result
type SearchResultItem struct {
	ID          string
	PatientID   string
	PatientName string
	SessionID   string
	Type        string
	Date        string
	Snippet    string
}

// SearchResultsViewModel is the view model for search results
type SearchResultsViewModel struct {
	Query     string
	Patients  []PatientSearchResultItem
	Results   []SearchResultItem
	Total     int
}