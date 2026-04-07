package appointment

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusNoShow    AppointmentStatus = "no_show"
)

// AppointmentType represents the type of appointment
type AppointmentType string

const (
	AppointmentTypeSession  AppointmentType = "session"
	AppointmentTypeFollowup AppointmentType = "followup"
	AppointmentTypeBlocked  AppointmentType = "blocked"
)

// Appointment represents a scheduled appointment in the clinical agenda
type Appointment struct {
	ID          string            `json:"id"`
	PatientID   string            `json:"patient_id,omitempty"`
	PatientName string            `json:"patient_name,omitempty"` // Denormalized for performance
	Date        time.Time         `json:"date"`                   // Date only (year, month, day)
	StartTime   string            `json:"start_time"`             // Format: "HH:MM"
	EndTime     string            `json:"end_time"`               // Format: "HH:MM"
	Duration    int               `json:"duration"`               // minutes
	Type        AppointmentType   `json:"type"`
	Status      AppointmentStatus `json:"status"`
	Notes       string            `json:"notes"`
	SessionID   *string           `json:"session_id,omitempty"` // Link to clinical session
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CancelledAt *time.Time        `json:"cancelled_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
}

// AppointmentFilter represents filters for querying appointments
type AppointmentFilter struct {
	View      string // day, week, month
	Date      time.Time
	PatientID string              // optional filter
	Status    []AppointmentStatus // optional filter
	Type      []AppointmentType   // optional filter
}

// TimeSlot represents an available or occupied time slot
type TimeSlot struct {
	Start       time.Time
	End         time.Time
	Available   bool
	Appointment *Appointment // if not available
}

// AgendaSettings represents the therapist's agenda configuration
type AgendaSettings struct {
	UserID            string    `json:"user_id"`
	SlotDuration      int       `json:"slot_duration"`       // minutes, default 50
	WorkStartTime     string    `json:"work_start_time"`     // "HH:MM", default "08:00"
	WorkEndTime       string    `json:"work_end_time"`       // "HH:MM", default "18:00"
	WorkDays          []int     `json:"work_days"`           // 0=Sunday, 6=Saturday, default [1,2,3,4,5]
	BreakBetweenSlots int       `json:"break_between_slots"` // minutes, default 10
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Validation errors
var (
	ErrInvalidDate      = errors.New("invalid appointment date")
	ErrInvalidTime      = errors.New("invalid appointment time")
	ErrInvalidDuration  = errors.New("invalid appointment duration")
	ErrPatientRequired  = errors.New("patient is required for session type")
	ErrInvalidStatus    = errors.New("invalid appointment status")
	ErrAlreadyCompleted = errors.New("cannot modify completed appointment")
	ErrAlreadyCancelled = errors.New("cannot modify cancelled appointment")
	ErrTimeInPast       = errors.New("cannot schedule appointment in the past")
)

// NewAppointment creates a new appointment
func NewAppointment(patientID, patientName string, date time.Time, startTime, endTime string, duration int, apptType AppointmentType, notes string) (*Appointment, error) {
	now := time.Now()

	// Validate required fields
	if date.IsZero() {
		return nil, ErrInvalidDate
	}

	if startTime == "" || endTime == "" {
		return nil, ErrInvalidTime
	}

	if duration < 30 || duration > 120 {
		return nil, ErrInvalidDuration
	}

	// Patient is required for session and followup types
	if (apptType == AppointmentTypeSession || apptType == AppointmentTypeFollowup) && patientID == "" {
		return nil, ErrPatientRequired
	}

	appointment := &Appointment{
		ID:          uuid.New().String(),
		PatientID:   patientID,
		PatientName: patientName,
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		Duration:    duration,
		Type:        apptType,
		Status:      AppointmentStatusScheduled,
		Notes:       notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return appointment, nil
}

// NewBlockedSlot creates a blocked time slot (e.g., lunch, vacation)
func NewBlockedSlot(date time.Time, startTime, endTime string, notes string) (*Appointment, error) {
	return NewAppointment("", "", date, startTime, endTime, 0, AppointmentTypeBlocked, notes)
}

// Update updates the appointment details
func (a *Appointment) Update(patientID, patientName string, date time.Time, startTime, endTime string, duration int, notes string) error {
	if a.Status == AppointmentStatusCompleted {
		return ErrAlreadyCompleted
	}

	if a.Status == AppointmentStatusCancelled {
		return ErrAlreadyCancelled
	}

	if !date.IsZero() {
		a.Date = date
	}

	if startTime != "" {
		a.StartTime = startTime
	}

	if endTime != "" {
		a.EndTime = endTime
	}

	if duration > 0 {
		a.Duration = duration
	}

	if patientID != "" {
		a.PatientID = patientID
	}

	if patientName != "" {
		a.PatientName = patientName
	}

	if notes != "" {
		a.Notes = notes
	}

	a.UpdatedAt = time.Now()

	return nil
}

// Confirm marks the appointment as confirmed
func (a *Appointment) Confirm() error {
	if a.Status == AppointmentStatusCompleted {
		return ErrAlreadyCompleted
	}

	if a.Status == AppointmentStatusCancelled {
		return ErrAlreadyCancelled
	}

	a.Status = AppointmentStatusConfirmed
	a.UpdatedAt = time.Now()

	return nil
}

// Complete marks the appointment as completed
func (a *Appointment) Complete() error {
	if a.Status == AppointmentStatusCancelled {
		return ErrAlreadyCancelled
	}

	now := time.Now()
	a.Status = AppointmentStatusCompleted
	a.CompletedAt = &now
	a.UpdatedAt = now

	return nil
}

// Cancel marks the appointment as cancelled
func (a *Appointment) Cancel() error {
	if a.Status == AppointmentStatusCompleted {
		return ErrAlreadyCompleted
	}

	now := time.Now()
	a.Status = AppointmentStatusCancelled
	a.CancelledAt = &now
	a.UpdatedAt = now

	return nil
}

// MarkNoShow marks the appointment as no-show
func (a *Appointment) MarkNoShow() error {
	if a.Status == AppointmentStatusCompleted {
		return ErrAlreadyCompleted
	}

	if a.Status == AppointmentStatusCancelled {
		return ErrAlreadyCancelled
	}

	a.Status = AppointmentStatusNoShow
	a.UpdatedAt = time.Now()

	return nil
}

// LinkSession links the appointment to a clinical session
func (a *Appointment) LinkSession(sessionID string) {
	a.SessionID = &sessionID
	a.UpdatedAt = time.Now()
}

// IsOverdue returns true if the appointment has passed and is still scheduled
func (a *Appointment) IsOverdue() bool {
	if a.Status != AppointmentStatusScheduled && a.Status != AppointmentStatusConfirmed {
		return false
	}

	// Parse start time
	startHour, startMin := parseTime(a.StartTime)

	// Create the actual appointment datetime
	apptTime := time.Date(
		a.Date.Year(), a.Date.Month(), a.Date.Day(),
		startHour, startMin, 0, 0, a.Date.Location(),
	)

	return time.Now().After(apptTime)
}

// GetStartDateTime returns the full datetime for the appointment start
func (a *Appointment) GetStartDateTime() time.Time {
	hour, min := parseTime(a.StartTime)
	return time.Date(
		a.Date.Year(), a.Date.Month(), a.Date.Day(),
		hour, min, 0, 0, a.Date.Location(),
	)
}

// GetEndDateTime returns the full datetime for the appointment end
func (a *Appointment) GetEndDateTime() time.Time {
	hour, min := parseTime(a.EndTime)
	return time.Date(
		a.Date.Year(), a.Date.Month(), a.Date.Day(),
		hour, min, 0, 0, a.Date.Location(),
	)
}

// Overlaps checks if this appointment overlaps with another
func (a *Appointment) Overlaps(other *Appointment) bool {
	// Different dates don't overlap
	if !sameDate(a.Date, other.Date) {
		return false
	}

	aStart := a.GetStartDateTime()
	aEnd := a.GetEndDateTime()
	otherStart := other.GetStartDateTime()
	otherEnd := other.GetEndDateTime()

	// Overlap occurs when one appointment starts before the other ends
	return aStart.Before(otherEnd) && aEnd.After(otherStart)
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(status string) bool {
	switch AppointmentStatus(status) {
	case AppointmentStatusScheduled, AppointmentStatusConfirmed, AppointmentStatusCompleted, AppointmentStatusCancelled, AppointmentStatusNoShow:
		return true
	}
	return false
}

// IsValidType checks if a type string is valid
func IsValidType(apptType string) bool {
	switch AppointmentType(apptType) {
	case AppointmentTypeSession, AppointmentTypeFollowup, AppointmentTypeBlocked:
		return true
	}
	return false
}

// DefaultAgendaSettings returns default agenda settings
func DefaultAgendaSettings(userID string) *AgendaSettings {
	now := time.Now()
	return &AgendaSettings{
		UserID:            userID,
		SlotDuration:      50,
		WorkStartTime:     "08:00",
		WorkEndTime:       "18:00",
		WorkDays:          []int{1, 2, 3, 4, 5}, // Monday to Friday
		BreakBetweenSlots: 10,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// Helper functions

func parseTime(timeStr string) (hour, min int) {
	// Parse time in "HH:MM" format
	t, _ := time.Parse("15:04", timeStr)
	return t.Hour(), t.Minute()
}

func sameDate(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day()
}
