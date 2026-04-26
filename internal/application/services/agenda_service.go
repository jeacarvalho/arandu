package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"arandu/internal/domain/appointment"
)

// AgendaServiceInterface defines the interface for agenda operations
type AgendaServiceInterface interface {
	CreateAppointment(ctx context.Context, patientID, patientName string, date time.Time, startTime, endTime string, duration int, apptType appointment.AppointmentType, notes string) (*appointment.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*appointment.Appointment, error)
	UpdateAppointment(ctx context.Context, id string, patientID, patientName string, date time.Time, startTime, endTime string, duration int, notes string) error
	CancelAppointment(ctx context.Context, id string) error
	ConfirmAppointment(ctx context.Context, id string) error
	MarkNoShow(ctx context.Context, id string) error
	CompleteAppointment(ctx context.Context, id string, sessionID string) error
	GetDayView(ctx context.Context, date time.Time) (*DayView, error)
	GetWeekView(ctx context.Context, date time.Time) (*WeekView, error)
	GetMonthView(ctx context.Context, year, month int) (*MonthView, error)
	GetAvailableSlots(ctx context.Context, date time.Time) ([]appointment.TimeSlot, error)
	CheckConflicts(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error)
	GetUpcomingAppointments(ctx context.Context, limit int) ([]*appointment.Appointment, error)
}

// AppointmentRepository defines the repository interface for appointments
type AppointmentRepository interface {
	Save(ctx context.Context, appt *appointment.Appointment) error
	FindByID(ctx context.Context, id string) (*appointment.Appointment, error)
	FindBySessionID(ctx context.Context, sessionID string) (*appointment.Appointment, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*appointment.Appointment, error)
	FindByPatient(ctx context.Context, patientID string) ([]*appointment.Appointment, error)
	FindByDate(ctx context.Context, date time.Time) ([]*appointment.Appointment, error)
	FindOverlapping(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error)
	Update(ctx context.Context, appt *appointment.Appointment) error
	Delete(ctx context.Context, id string) error
	FindUpcoming(ctx context.Context, fromDate time.Time, limit int) ([]*appointment.Appointment, error)
}

// DayView represents a single day's appointments
type DayView struct {
	Date         time.Time
	Appointments []*appointment.Appointment
	Slots        []appointment.TimeSlot
}

// WeekView represents a week's appointments
type WeekView struct {
	StartDate    time.Time
	EndDate      time.Time
	Days         []DayView
	Appointments []*appointment.Appointment
}

// MonthView represents a month's appointments
type MonthView struct {
	Year         int
	Month        time.Month
	StartDate    time.Time
	EndDate      time.Time
	Days         []DayInfo
	Appointments []*appointment.Appointment
}

// DayInfo represents a day in month view
type DayInfo struct {
	Date         time.Time
	Appointments int
	HasEvents    bool
}

// AgendaService implements agenda business logic
type AgendaService struct {
	apptRepo AppointmentRepository
}

// NewAgendaService creates a new agenda service
func NewAgendaService(apptRepo AppointmentRepository) *AgendaService {
	return &AgendaService{
		apptRepo: apptRepo,
	}
}

// CreateAppointment creates a new appointment with conflict checking
func (s *AgendaService) CreateAppointment(ctx context.Context, patientID, patientName string, date time.Time, startTime, endTime string, duration int, apptType appointment.AppointmentType, notes string) (*appointment.Appointment, error) {
	// Check for conflicts
	conflicts, err := s.apptRepo.FindOverlapping(ctx, date, startTime, endTime, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check conflicts: %w", err)
	}

	if len(conflicts) > 0 {
		return nil, fmt.Errorf("time slot conflicts with existing appointment")
	}

	// Create appointment
	appt, err := appointment.NewAppointment(patientID, patientName, date, startTime, endTime, duration, apptType, notes)
	if err != nil {
		return nil, err
	}

	// Save to database
	if err := s.apptRepo.Save(ctx, appt); err != nil {
		return nil, fmt.Errorf("failed to save appointment: %w", err)
	}

	return appt, nil
}

// GetAppointment retrieves an appointment by ID
func (s *AgendaService) GetAppointment(ctx context.Context, id string) (*appointment.Appointment, error) {
	return s.apptRepo.FindByID(ctx, id)
}

// UpdateAppointment updates an appointment with conflict checking
func (s *AgendaService) UpdateAppointment(ctx context.Context, id string, patientID, patientName string, date time.Time, startTime, endTime string, duration int, notes string) error {
	// Find existing appointment
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find appointment: %w", err)
	}
	if appt == nil {
		return fmt.Errorf("appointment not found")
	}

	// Check for conflicts (excluding current appointment)
	conflicts, err := s.apptRepo.FindOverlapping(ctx, date, startTime, endTime, id)
	if err != nil {
		return fmt.Errorf("failed to check conflicts: %w", err)
	}

	if len(conflicts) > 0 {
		return fmt.Errorf("time slot conflicts with existing appointment")
	}

	// Update appointment
	if err := appt.Update(patientID, patientName, date, startTime, endTime, duration, notes); err != nil {
		return err
	}

	// Save changes
	if err := s.apptRepo.Update(ctx, appt); err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	return nil
}

// CancelAppointment cancels an appointment
func (s *AgendaService) CancelAppointment(ctx context.Context, id string) error {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find appointment: %w", err)
	}
	if appt == nil {
		return fmt.Errorf("appointment not found")
	}

	if err := appt.Cancel(); err != nil {
		return err
	}

	return s.apptRepo.Update(ctx, appt)
}

// ConfirmAppointment confirms an appointment
func (s *AgendaService) ConfirmAppointment(ctx context.Context, id string) error {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find appointment: %w", err)
	}
	if appt == nil {
		return fmt.Errorf("appointment not found")
	}

	if err := appt.Confirm(); err != nil {
		return err
	}

	return s.apptRepo.Update(ctx, appt)
}

// CompleteAppointment marks an appointment as completed and links to session
func (s *AgendaService) CompleteAppointment(ctx context.Context, id string, sessionID string) error {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find appointment: %w", err)
	}
	if appt == nil {
		return fmt.Errorf("appointment not found")
	}

	if err := appt.Complete(); err != nil {
		return err
	}

	if sessionID != "" {
		appt.LinkSession(sessionID)
	}

	return s.apptRepo.Update(ctx, appt)
}

// GetDayView retrieves appointments for a specific day
func (s *AgendaService) GetDayView(ctx context.Context, date time.Time) (*DayView, error) {
	appointments, err := s.apptRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get day appointments: %w", err)
	}

	slots := s.generateDaySlots(date, appointments)

	return &DayView{
		Date:         date,
		Appointments: appointments,
		Slots:        slots,
	}, nil
}

// GetWeekView retrieves appointments for a week
func (s *AgendaService) GetWeekView(ctx context.Context, date time.Time) (*WeekView, error) {
	// Calculate week start (Monday) and end (Sunday)
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday as 7
	}
	startDate := date.AddDate(0, 0, -(weekday - 1)) // Monday
	endDate := startDate.AddDate(0, 0, 6)           // Sunday

	appointments, err := s.apptRepo.FindByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get week appointments: %w", err)
	}

	// Group by day
	days := make([]DayView, 7)
	for i := 0; i < 7; i++ {
		dayDate := startDate.AddDate(0, 0, i)
		dayAppointments := s.filterByDate(appointments, dayDate)
		days[i] = DayView{
			Date:         dayDate,
			Appointments: dayAppointments,
			Slots:        s.generateDaySlots(dayDate, dayAppointments),
		}
	}

	return &WeekView{
		StartDate:    startDate,
		EndDate:      endDate,
		Days:         days,
		Appointments: appointments,
	}, nil
}

// GetMonthView retrieves appointments for a month
func (s *AgendaService) GetMonthView(ctx context.Context, year, month int) (*MonthView, error) {
	// Calculate month boundaries
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// Expand to include full weeks
	weekday := int(firstOfMonth.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startDate := firstOfMonth.AddDate(0, 0, -(weekday - 1))

	weekday = int(lastOfMonth.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	endDate := lastOfMonth.AddDate(0, 0, (7 - weekday))

	appointments, err := s.apptRepo.FindByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get month appointments: %w", err)
	}

	// Build day info
	var days []DayInfo
	current := startDate
	for !current.After(endDate) {
		dayAppointments := s.filterByDate(appointments, current)
		days = append(days, DayInfo{
			Date:         current,
			Appointments: len(dayAppointments),
			HasEvents:    len(dayAppointments) > 0,
		})
		current = current.AddDate(0, 0, 1)
	}

	return &MonthView{
		Year:         year,
		Month:        time.Month(month),
		StartDate:    startDate,
		EndDate:      endDate,
		Days:         days,
		Appointments: appointments,
	}, nil
}

// GetAvailableSlots retrieves available time slots for a date
func (s *AgendaService) GetAvailableSlots(ctx context.Context, date time.Time) ([]appointment.TimeSlot, error) {
	appointments, err := s.apptRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get day appointments: %w", err)
	}

	slots := s.generateDaySlots(date, appointments)
	return slots, nil
}

// CheckConflicts checks for overlapping appointments
func (s *AgendaService) CheckConflicts(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	return s.apptRepo.FindOverlapping(ctx, date, startTime, endTime, excludeID)
}

// GetUpcomingAppointments retrieves upcoming appointments
func (s *AgendaService) GetUpcomingAppointments(ctx context.Context, limit int) ([]*appointment.Appointment, error) {
	return s.apptRepo.FindUpcoming(ctx, time.Now(), limit)
}

// MarkNoShow marks an appointment as no-show
func (s *AgendaService) MarkNoShow(ctx context.Context, id string) error {
	appt, err := s.apptRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find appointment: %w", err)
	}
	if appt == nil {
		return fmt.Errorf("appointment not found")
	}

	if err := appt.MarkNoShow(); err != nil {
		return err
	}

	return s.apptRepo.Update(ctx, appt)
}

// RescheduleAppointment reschedules an appointment
func (s *AgendaService) RescheduleAppointment(ctx context.Context, id string, newDate time.Time, newStartTime, newEndTime string) error {
	return s.UpdateAppointment(ctx, id, "", "", newDate, newStartTime, newEndTime, 0, "")
}

// Helper methods

func (s *AgendaService) generateDaySlots(date time.Time, appointments []*appointment.Appointment) []appointment.TimeSlot {
	// Generate slots from 08:00 to 20:00 in 30-minute intervals
	var slots []appointment.TimeSlot
	startHour := 8
	endHour := 20

	for hour := startHour; hour < endHour; hour++ {
		for _, minute := range []int{0, 30} {
			slotStart := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
			slotEnd := slotStart.Add(30 * time.Minute)

			slot := appointment.TimeSlot{
				Start:     slotStart,
				End:       slotEnd,
				Available: true,
			}

			// Check if slot overlaps with any appointment
			for _, appt := range appointments {
				if appt.Status == appointment.AppointmentStatusCancelled {
					continue
				}

				apptStart := appt.GetStartDateTime()
				apptEnd := appt.GetEndDateTime()

				// Check overlap
				if slotStart.Before(apptEnd) && slotEnd.After(apptStart) {
					slot.Available = false
					slot.Appointment = appt
					break
				}
			}

			slots = append(slots, slot)
		}
	}

	return slots
}

func (s *AgendaService) filterByDate(appointments []*appointment.Appointment, date time.Time) []*appointment.Appointment {
	var filtered []*appointment.Appointment
	for _, appt := range appointments {
		if sameDate(appt.Date, date) {
			filtered = append(filtered, appt)
		}
	}
	return filtered
}

func sameDate(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day()
}

func (s *AgendaService) GetPatientAppointments(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	appts, err := s.apptRepo.FindByPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}
	sort.Slice(appts, func(i, j int) bool {
		return appts[i].Date.After(appts[j].Date)
	})
	return appts, nil
}

func (s *AgendaService) GetAppointmentBySession(ctx context.Context, sessionID string) (*appointment.Appointment, error) {
	return s.apptRepo.FindBySessionID(ctx, sessionID)
}
