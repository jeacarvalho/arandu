package agenda

import "time"

// AgendaViewModel represents the agenda view data
type AgendaViewModel struct {
	WeekStart        time.Time
	WeekEnd        time.Time
	WeekLabel      string // "6 – 12 de abril de 2026"
	MonthLabel    string // "abril 2026"
	Days          []DayViewModel
	TotalCount    int
	Confirmed    int
	Pending      int
	CurrentView  string // "dia" | "semana" | "mes"
	PrevDate     string // "2026-03-30"
	NextDate     string // "2026-04-13"
	MonthStartOffset int    // 0=Seg, 1=Ter, ..., 6=Dom — células vazias antes do dia 1
	Metrics      []AgendaMetric
}

// DayViewModel represents a single day in the agenda
type DayViewModel struct {
	Date             time.Time
	DayName          string // "Seg"
	DayNumber        string // "6"
	IsToday          bool
	IsCurrentMonth   bool // false for prev/next month filler days in month view
	Appointments     []AppointmentViewModel
	AppointmentCount int // used in month view (DayInfo only has count, not individual appts)
}

// AppointmentViewModel represents a single appointment
type AppointmentViewModel struct {
	ID          string
	PatientName string
	StartTime   string // "10:00"
	EndTime     string // "10:50"
	SessionType string // "Sessão individual"
	Status      string // "confirmed" | "pending" | "first_session" | "cancelled"
	Tone        string // "accent" | "info" | "warn" | "danger" | "moss" | "ghost" | "neutral"
	TopPx       int    // (startHour - 8) * 60 + startMinute
	HeightPx    int    // durationMinutes - 4
}

// AgendaMetric represents a metric displayed in the hero
type AgendaMetric struct {
	Value string
	Label string
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
	case "scheduled":
		return "Agendada"
	case "confirmed":
		return "Confirmada"
	case "pending":
		return "Aguardando confirmação"
	case "first_session":
		return "1ª consulta"
	case "cancelled":
		return "Cancelado"
	case "no_show":
		return "Não Compareceu"
	case "completed":
		return "Realizada"
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

// GetAppointmentCardClasses returns Tailwind CSS classes for appointment card styling
func GetAppointmentCardClasses(status string) string {
	switch status {
	case "confirmed":
		return "bg-emerald-100 text-emerald-900 border-emerald-500"
	case "pending":
		return "bg-amber-100 text-amber-900 border-amber-500"
	case "first_session":
		return "bg-blue-100 text-blue-900 border-blue-500"
	case "cancelled":
		return "bg-neutral-100 text-neutral-500 border-neutral-300 line-through opacity-65"
	default:
		return "bg-emerald-100 text-emerald-900 border-emerald-500"
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

// TranslateDayName translates English day abbreviations to Portuguese
func TranslateDayName(englishDay string) string {
	switch englishDay {
	case "Mon":
		return "Seg"
	case "Tue":
		return "Ter"
	case "Wed":
		return "Qua"
	case "Thu":
		return "Qui"
	case "Fri":
		return "Sex"
	case "Sat":
		return "Sáb"
	case "Sun":
		return "Dom"
	default:
		return englishDay
	}
}

// ViewTabClasses returns Tailwind CSS classes for view tabs (dia/semana/mes)
// IMPORTANT: Tailwind v4 scans source files as plain text and cannot detect
// class names constructed dynamically inside templ.KV().
// Always use helper functions that return complete class strings.
func ViewTabClasses(viewName string, currentView string) string {
	baseClasses := "h-9 px-4 text-sm font-medium rounded-lg transition-all focus:outline-none focus:ring-2 focus:ring-arandu-primary/50"

	if viewName == currentView {
		return baseClasses + " bg-arandu-primary text-white shadow-sm"
	}
	return baseClasses + " text-neutral-500 hover:text-neutral-800 hover:bg-neutral-100"
}

// TodayBadgeClasses returns Tailwind CSS classes for today badge in day headers
func TodayBadgeClasses(isToday bool) string {
	base := "w-6.5 h-6.5 rounded-full flex items-center justify-center text-sm font-semibold mt-0.5"
	if isToday {
		return base + " bg-arandu-primary text-white"
	}
	return base + " text-neutral-600"
}

// DayGridClass returns the correct grid-cols class based on number of days shown
func DayGridClass(numDays int) string {
	if numDays == 1 {
		return "grid-cols-1"
	}
	return "grid-cols-7"
}

// MonthDayNumberClass returns CSS for the day number bubble in month grid
func MonthDayNumberClass(isToday bool) string {
	base := "w-6 h-6 rounded-full flex items-center justify-center text-xs font-semibold"
	if isToday {
		return base + " bg-arandu-primary text-white"
	}
	return base + " text-neutral-600"
}

// MonthOffsetClass returns the col-start class for the first day of the month
func MonthOffsetClass(offset int) string {
	switch offset {
	case 1:
		return "col-start-2"
	case 2:
		return "col-start-3"
	case 3:
		return "col-start-4"
	case 4:
		return "col-start-5"
	case 5:
		return "col-start-6"
	case 6:
		return "col-start-7"
	default:
		return "col-start-1"
	}
}

// MonthCellClasses returns Tailwind CSS classes for a month grid cell
// Fixed height so all cells are the same size regardless of appointment count
func MonthCellClasses(isToday, isCurrentMonth bool) string {
	base := "p-2 h-[110px] overflow-hidden"
	if !isCurrentMonth {
		return base + " bg-neutral-50"
	}
	if isToday {
		return base + " bg-arandu-primary/5"
	}
	return base + " bg-white"
}

// MonthAppointmentPillClasses returns classes for an appointment pill in month view
func MonthAppointmentPillClasses(status string) string {
	base := "flex items-center gap-1 mt-0.5 px-1.5 py-0.5 rounded text-xs border-l-2 overflow-hidden cursor-pointer hover:brightness-95 transition-all"
	switch status {
	case "confirmed":
		return base + " border-emerald-500 bg-emerald-50 text-emerald-800"
	case "pending":
		return base + " border-amber-500 bg-amber-50 text-amber-800"
	case "first_session":
		return base + " border-blue-500 bg-blue-50 text-blue-800"
	case "cancelled":
		return base + " border-neutral-300 bg-neutral-50 text-neutral-400 line-through"
	default:
		return base + " border-emerald-500 bg-emerald-50 text-emerald-800"
	}
}

// MonthDayNumberOutsideClass returns CSS for a day number outside current month
func MonthDayNumberOutsideClass() string {
	return "w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium text-neutral-300"
}

// TodayColumnClasses returns Tailwind CSS classes for today column background
func TodayColumnClasses(isToday bool) string {
	base := "border-r border-neutral-200 relative last:border-r-0 min-h-[780px]"
	if isToday {
		return base + " bg-arandu-primary/5"
	}
	return base
}

// ConflictAlertClasses returns Tailwind CSS classes for conflict error alert
func ConflictAlertClasses(visible bool) string {
	if visible {
		return "mt-2 p-3 rounded-lg text-sm bg-red-50 border border-red-200 text-red-700"
	}
	return "hidden"
}

// RescheduleFormModel represents the reschedule form
type RescheduleFormModel struct {
	AppointmentID string
	CurrentDate   time.Time
	CurrentStart  string
	Duration      int
}

// StatusToTone maps appointment status to Sábio tone
func StatusToTone(status string) string {
	switch status {
	case "confirmed", "completed", "scheduled":
		return "accent"
	case "first_session":
		return "info"
	case "pending":
		return "warn"
	case "no_show", "risk":
		return "danger"
	case "cancelled":
		return "ghost"
	default:
		return "neutral"
	}
}

// ToneStyles returns background, border, and text colors for week view appointments
func ToneStyles(tone string) (bg, bd, fg string) {
	switch tone {
	case "accent":
		return "#EADFCB", "#6B4E3D", "#3E2A1E"
	case "info":
		return "#DDE6EE", "#3C5C7A", "#1F3A55"
	case "warn":
		return "#F3E3C4", "#B8842A", "#6B4A10"
	case "danger":
		return "#EDCFC7", "#A0463A", "#6F241A"
	case "moss":
		return "#D9E2D3", "#4A5D4F", "#263326"
	case "ghost":
		return "repeating-linear-gradient(135deg,var(--paper-3),var(--paper-3) 4px,var(--paper-2) 4px,var(--paper-2) 8px)", "var(--ink-4)", "var(--ink-3)"
	default:
		return "var(--paper)", "var(--ink-4)", "var(--ink-2)"
	}
}
