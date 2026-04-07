package agenda

import "time"

// AgendaViewModel represents the agenda view data
type AgendaViewModel struct {
	WeekStart   time.Time
	WeekEnd     time.Time
	WeekLabel   string // "6 – 12 de abril de 2026"
	Days        []DayViewModel
	TotalCount  int
	CurrentView string // "dia" | "semana" | "mes"
	PrevDate    string // "2026-03-30"
	NextDate    string // "2026-04-13"
}

// DayViewModel represents a single day in the agenda
type DayViewModel struct {
	Date         time.Time
	DayName      string // "Seg"
	DayNumber    string // "6"
	IsToday      bool
	Appointments []AppointmentViewModel
}

// AppointmentViewModel represents a single appointment
type AppointmentViewModel struct {
	ID          string
	PatientName string
	StartTime   string // "10:00"
	EndTime     string // "10:50"
	SessionType string // "Sessão individual"
	Status      string // "confirmed" | "pending" | "first_session" | "cancelled"
	TopPx       int    // (startHour - 8) * 60 + startMinute
	HeightPx    int    // durationMinutes - 4
}

// PatientOption represents a patient for autocomplete
type PatientOption struct {
	ID   string
	Name string
}

// TimeSlotViewModel represents a time slot for the UI
type TimeSlotViewModel struct {
	StartTime     string
	EndTime       string
	Available     bool
	AppointmentID string
	PatientName   string
}

// NewAppointmentFormData represents data for the new appointment form
type NewAppointmentFormData struct {
	Date        time.Time
	DefaultTime string
	Patients    []PatientOption
	Slots       []TimeSlotViewModel
}

// SlotsViewModel represents available slots view
type SlotsViewModel struct {
	Date  time.Time
	Slots []TimeSlotViewModel
}

// AppointmentDetailModel represents appointment details
type AppointmentDetailModel struct {
	ID          string
	PatientID   string
	PatientName string
	Date        time.Time
	StartTime   string
	EndTime     string
	Duration    int
	SessionType string
	Status      string
	Notes       string
	SessionID   *string
}

// StatusLabel returns the human-readable label for a status
func StatusLabel(status string) string {
	switch status {
	case "pending":
		return "Aguardando confirmação"
	case "first_session":
		return "1ª consulta"
	case "cancelled":
		return "Cancelado"
	default:
		return "Sessão individual"
	}
}

// GetAppointmentStatusClass returns the CSS class for appointment status
func GetAppointmentStatusClass(status string) string {
	switch status {
	case "confirmed":
		return "appt confirmed"
	case "pending":
		return "appt pending"
	case "first_session":
		return "appt first_session"
	case "cancelled":
		return "appt cancelled"
	default:
		return "appt confirmed"
	}
}

// GetLegendDotClass returns the CSS class for legend dot
func GetLegendDotClass(status string) string {
	switch status {
	case "confirmed":
		return "leg-dot confirmed"
	case "pending":
		return "leg-dot pending"
	case "first_session":
		return "leg-dot first_session"
	case "cancelled":
		return "leg-dot cancelled"
	default:
		return "leg-dot confirmed"
	}
}
