package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"arandu/internal/domain/appointment"
)

// AppointmentRepository handles appointment data access
type AppointmentRepository struct {
	db *DB
}

// NewAppointmentRepository creates a new appointment repository
func NewAppointmentRepository(db *DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

// Save creates a new appointment
func (r *AppointmentRepository) Save(ctx context.Context, appt *appointment.Appointment) error {
	query := `INSERT INTO appointments (
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		appt.ID,
		appt.PatientID,
		appt.PatientName,
		appt.Date,
		appt.StartTime,
		appt.EndTime,
		appt.Duration,
		string(appt.Type),
		string(appt.Status),
		appt.Notes,
		appt.SessionID,
		appt.CreatedAt,
		appt.UpdatedAt,
	)
	return err
}

// FindByID retrieves an appointment by ID
func (r *AppointmentRepository) FindByID(ctx context.Context, id string) (*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanAppointment(row)
}

// FindByDateRange retrieves appointments within a date range
func (r *AppointmentRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments 
	WHERE date >= ? AND date <= ?
	ORDER BY date, start_time`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

// FindByPatient retrieves all appointments for a patient
func (r *AppointmentRepository) FindByPatient(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments 
	WHERE patient_id = ?
	ORDER BY date DESC, start_time DESC`

	rows, err := r.db.QueryContext(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

// FindByPatientAndDateRange retrieves appointments for a patient within a date range
func (r *AppointmentRepository) FindByPatientAndDateRange(ctx context.Context, patientID string, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments 
	WHERE patient_id = ? AND date >= ? AND date <= ?
	ORDER BY date, start_time`

	rows, err := r.db.QueryContext(ctx, query, patientID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

// FindByDate retrieves appointments for a specific date
func (r *AppointmentRepository) FindByDate(ctx context.Context, date time.Time) ([]*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments 
	WHERE date = ?
	ORDER BY start_time`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

// FindOverlapping finds appointments that overlap with the given time range
func (r *AppointmentRepository) FindOverlapping(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments 
	WHERE date = ? 
		AND status NOT IN ('cancelled')
		AND (
			(start_time < ? AND end_time > ?) OR
			(start_time >= ? AND start_time < ?) OR
			(end_time > ? AND end_time <= ?)
		)`

	args := []interface{}{
		date, endTime, startTime,
		startTime, endTime,
		startTime, endTime,
	}

	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	query += ` ORDER BY start_time`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

// Update updates an appointment
func (r *AppointmentRepository) Update(ctx context.Context, appt *appointment.Appointment) error {
	query := `UPDATE appointments SET
		patient_id = ?,
		patient_name = ?,
		date = ?,
		start_time = ?,
		end_time = ?,
		duration = ?,
		type = ?,
		status = ?,
		notes = ?,
		session_id = ?,
		updated_at = ?,
		cancelled_at = ?,
		completed_at = ?
	WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		appt.PatientID,
		appt.PatientName,
		appt.Date,
		appt.StartTime,
		appt.EndTime,
		appt.Duration,
		string(appt.Type),
		string(appt.Status),
		appt.Notes,
		appt.SessionID,
		appt.UpdatedAt,
		appt.CancelledAt,
		appt.CompletedAt,
		appt.ID,
	)
	return err
}

// Delete deletes an appointment
func (r *AppointmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM appointments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CountByDate counts appointments for a specific date
func (r *AppointmentRepository) CountByDate(ctx context.Context, date time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM appointments WHERE date = ? AND status NOT IN ('cancelled')`
	var count int
	err := r.db.QueryRowContext(ctx, query, date).Scan(&count)
	return count, err
}

// FindUpcoming retrieves upcoming appointments from a given date
func (r *AppointmentRepository) FindUpcoming(ctx context.Context, fromDate time.Time, limit int) ([]*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments 
	WHERE date >= ? AND status NOT IN ('cancelled', 'completed', 'no_show')
	ORDER BY date ASC, start_time ASC
	LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, fromDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAppointments(rows)
}

// FindBySessionID finds an appointment by its linked session ID
func (r *AppointmentRepository) FindBySessionID(ctx context.Context, sessionID string) (*appointment.Appointment, error) {
	query := `SELECT 
		id, patient_id, patient_name, date, start_time, end_time, duration,
		type, status, notes, session_id, created_at, updated_at, cancelled_at, completed_at
	FROM appointments WHERE session_id = ?`

	row := r.db.QueryRowContext(ctx, query, sessionID)
	return r.scanAppointment(row)
}

// UpdatePatientName updates the denormalized patient name
func (r *AppointmentRepository) UpdatePatientName(ctx context.Context, patientID, patientName string) error {
	query := `UPDATE appointments SET patient_name = ?, updated_at = ? WHERE patient_id = ?`
	_, err := r.db.ExecContext(ctx, query, patientName, time.Now(), patientID)
	return err
}

// Helper methods

func (r *AppointmentRepository) scanAppointment(row *sql.Row) (*appointment.Appointment, error) {
	var appt appointment.Appointment
	var patientID, patientName, notes, sessionID sql.NullString
	var cancelledAt, completedAt sql.NullTime

	err := row.Scan(
		&appt.ID,
		&patientID,
		&patientName,
		&appt.Date,
		&appt.StartTime,
		&appt.EndTime,
		&appt.Duration,
		&appt.Type,
		&appt.Status,
		&notes,
		&sessionID,
		&appt.CreatedAt,
		&appt.UpdatedAt,
		&cancelledAt,
		&completedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if patientID.Valid {
		appt.PatientID = patientID.String
	}
	if patientName.Valid {
		appt.PatientName = patientName.String
	}
	if notes.Valid {
		appt.Notes = notes.String
	}
	if sessionID.Valid {
		appt.SessionID = &sessionID.String
	}
	if cancelledAt.Valid {
		appt.CancelledAt = &cancelledAt.Time
	}
	if completedAt.Valid {
		appt.CompletedAt = &completedAt.Time
	}

	return &appt, nil
}

func (r *AppointmentRepository) scanAppointments(rows *sql.Rows) ([]*appointment.Appointment, error) {
	var appointments []*appointment.Appointment

	for rows.Next() {
		var appt appointment.Appointment
		var patientID, patientName, notes, sessionID sql.NullString
		var cancelledAt, completedAt sql.NullTime

		err := rows.Scan(
			&appt.ID,
			&patientID,
			&patientName,
			&appt.Date,
			&appt.StartTime,
			&appt.EndTime,
			&appt.Duration,
			&appt.Type,
			&appt.Status,
			&notes,
			&sessionID,
			&appt.CreatedAt,
			&appt.UpdatedAt,
			&cancelledAt,
			&completedAt,
		)
		if err != nil {
			return nil, err
		}

		if patientID.Valid {
			appt.PatientID = patientID.String
		}
		if patientName.Valid {
			appt.PatientName = patientName.String
		}
		if notes.Valid {
			appt.Notes = notes.String
		}
		if sessionID.Valid {
			appt.SessionID = &sessionID.String
		}
		if cancelledAt.Valid {
			appt.CancelledAt = &cancelledAt.Time
		}
		if completedAt.Valid {
			appt.CompletedAt = &completedAt.Time
		}

		appointments = append(appointments, &appt)
	}

	return appointments, rows.Err()
}

// AgendaSettingsRepository handles agenda settings data access
type AgendaSettingsRepository struct {
	db *DB
}

// NewAgendaSettingsRepository creates a new agenda settings repository
func NewAgendaSettingsRepository(db *DB) *AgendaSettingsRepository {
	return &AgendaSettingsRepository{db: db}
}

// Save creates or updates agenda settings
func (r *AgendaSettingsRepository) Save(ctx context.Context, settings *appointment.AgendaSettings) error {
	query := `INSERT INTO agenda_settings (
		user_id, slot_duration, work_start_time, work_end_time, work_days,
		break_between_slots, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(user_id) DO UPDATE SET
		slot_duration = excluded.slot_duration,
		work_start_time = excluded.work_start_time,
		work_end_time = excluded.work_end_time,
		work_days = excluded.work_days,
		break_between_slots = excluded.break_between_slots,
		updated_at = excluded.updated_at`

	workDays := strings.Trim(strings.Join(strings.Fields(fmt.Sprintf("%d", settings.WorkDays)), ","), "[]")

	_, err := r.db.ExecContext(ctx, query,
		settings.UserID,
		settings.SlotDuration,
		settings.WorkStartTime,
		settings.WorkEndTime,
		workDays,
		settings.BreakBetweenSlots,
		settings.CreatedAt,
		settings.UpdatedAt,
	)
	return err
}

// FindByUserID retrieves agenda settings for a user
func (r *AgendaSettingsRepository) FindByUserID(ctx context.Context, userID string) (*appointment.AgendaSettings, error) {
	query := `SELECT user_id, slot_duration, work_start_time, work_end_time, work_days,
		break_between_slots, created_at, updated_at
	FROM agenda_settings WHERE user_id = ?`

	var settings appointment.AgendaSettings
	var workDaysStr string

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&settings.UserID,
		&settings.SlotDuration,
		&settings.WorkStartTime,
		&settings.WorkEndTime,
		&workDaysStr,
		&settings.BreakBetweenSlots,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse work days
	if workDaysStr != "" {
		settings.WorkDays = parseWorkDays(workDaysStr)
	}

	return &settings, nil
}

func parseWorkDays(workDaysStr string) []int {
	var days []int
	for _, d := range strings.Split(workDaysStr, ",") {
		if day, err := strconv.Atoi(strings.TrimSpace(d)); err == nil {
			days = append(days, day)
		}
	}
	return days
}
