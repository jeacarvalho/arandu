package notes

import (
	"fmt"
	"strings"
	"time"
)

type NotesLibraryViewModel struct {
	Records        []NoteRecord
	FocusedRecord  NoteRecordDetail
	FilterTags    []string
	TotalRecords int
	TotalEvents  int
	CurrentFilter string
	CurrentPage int
	PageSize   int
	SearchQuery string
}

type NoteRecord struct {
	PatientID   string
	PatientName string
	Initials    string
	RecordID    string
	EventCount  int
	LastUpdate  string
	Tags        []string
	Status     string
	IsFocused   bool
}

type NoteEventItem struct {
	DateStr  string
	Content string
}

type NoteRecordDetail struct {
	PatientID          string
	PatientName       string
	RecordID        string
	ActiveTab       string
	Sections        []NoteSection
	LastEntryDate    string
	LastEntryTitle   string
	LastEntryContent string
	LastEntryQuote   string
	SessionCount      int
	ObservationCount int
	InterventionCount int
	TotalEvents      int
	FocusedPatientID string
	ObservationItems  []NoteEventItem
	InterventionItems []NoteEventItem
	AnamneseText      string
}

type NoteSection struct {
	Key    string
	Label string
	Count int
	Active bool
}

func BuildNotesInitials(name string) string {
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

func BuildRecordID(patientID string) string {
	id := patientID
	if len(id) > 4 {
		id = id[len(id)-4:]
	}
	return "PR-" + strings.ToUpper(id)
}

func TruncateStr(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

func BuildEventTitle(eventType, title string) string {
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

func FormatTherapyMonths(createdAt time.Time) string {
	months := int(time.Since(createdAt).Hours() / 730)
	if months < 1 {
		return "< 1m"
	}
	years := months / 12
	rem := months % 12
	if years == 0 {
		return fmt.Sprintf("%dm", months)
	}
	if rem == 0 {
		return fmt.Sprintf("%da", years)
	}
	return fmt.Sprintf("%da %dm", years, rem)
}