package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"arandu/internal/domain/appointment"
)

type mockAppointmentRepo struct {
	saveFunc            func(ctx context.Context, appt *appointment.Appointment) error
	findByIDFunc        func(ctx context.Context, id string) (*appointment.Appointment, error)
	findByDateRangeFunc  func(ctx context.Context, start, end time.Time) ([]*appointment.Appointment, error)
	findByDateFunc      func(ctx context.Context, date time.Time) ([]*appointment.Appointment, error)
	findOverlappingFunc func(ctx context.Context, date time.Time, start, end, excludeID string) ([]*appointment.Appointment, error)
	updateFunc          func(ctx context.Context, appt *appointment.Appointment) error
	deleteFunc          func(ctx context.Context, id string) error
	findUpcomingFunc    func(ctx context.Context, from time.Time, limit int) ([]*appointment.Appointment, error)
	findByPatientFunc  func(ctx context.Context, patientID string) ([]*appointment.Appointment, error)
}

func (m *mockAppointmentRepo) Save(ctx context.Context, appt *appointment.Appointment) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, appt)
	}
	return nil
}

func (m *mockAppointmentRepo) FindByID(ctx context.Context, id string) (*appointment.Appointment, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockAppointmentRepo) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	if m.findByDateRangeFunc != nil {
		return m.findByDateRangeFunc(ctx, startDate, endDate)
	}
	return nil, nil
}

func (m *mockAppointmentRepo) FindByDate(ctx context.Context, date time.Time) ([]*appointment.Appointment, error) {
	if m.findByDateFunc != nil {
		return m.findByDateFunc(ctx, date)
	}
	return nil, nil
}

func (m *mockAppointmentRepo) FindOverlapping(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	if m.findOverlappingFunc != nil {
		return m.findOverlappingFunc(ctx, date, startTime, endTime, excludeID)
	}
	return nil, nil
}

func (m *mockAppointmentRepo) Update(ctx context.Context, appt *appointment.Appointment) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, appt)
	}
	return nil
}

func (m *mockAppointmentRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockAppointmentRepo) FindUpcoming(ctx context.Context, fromDate time.Time, limit int) ([]*appointment.Appointment, error) {
	if m.findUpcomingFunc != nil {
		return m.findUpcomingFunc(ctx, fromDate, limit)
	}
	return nil, nil
}

func (m *mockAppointmentRepo) FindByPatient(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	if m.findByPatientFunc != nil {
		return m.findByPatientFunc(ctx, patientID)
	}
	return nil, nil
}

func (m *mockAppointmentRepo) FindBySessionID(ctx context.Context, sessionID string) (*appointment.Appointment, error) {
	return nil, nil
}

func TestAgendaService_GetWeekView_StartOnMonday(t *testing.T) {
	repo := &mockAppointmentRepo{}
	svc := NewAgendaService(repo)

	wednesday := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	weekView, err := svc.GetWeekView(context.Background(), wednesday)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weekView.StartDate.Weekday() != time.Monday {
		t.Errorf("expected week to start on Monday, got %s", weekView.StartDate.Weekday())
	}
	if len(weekView.Days) != 7 {
		t.Errorf("expected 7 days in week, got %d", len(weekView.Days))
	}
}

func TestAgendaService_GetWeekView_GroupsAppointmentsByDay(t *testing.T) {
	repo := &mockAppointmentRepo{}
	svc := NewAgendaService(repo)

	friday := time.Date(2026, 4, 24, 0, 0, 0, 0, time.Local)

	weekView, err := svc.GetWeekView(context.Background(), friday)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weekView.StartDate.Weekday() != time.Monday {
		t.Errorf("expected week to start on Monday, got %s", weekView.StartDate.Weekday())
	}
	if len(weekView.Days) != 7 {
		t.Errorf("expected 7 days, got %d", len(weekView.Days))
	}
}

func TestAgendaService_GetMonthView_FullWeeks(t *testing.T) {
	repo := &mockAppointmentRepo{}
	svc := NewAgendaService(repo)

	monthView, err := svc.GetMonthView(context.Background(), 2026, 4)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(monthView.Days)%7 != 0 {
		t.Errorf("expected days to be multiple of 7, got %d", len(monthView.Days))
	}
	if monthView.Days[0].Date.Weekday() != time.Monday {
		t.Errorf("expected first day to be Monday, got %s", monthView.Days[0].Date.Weekday())
	}
}

func TestAgendaService_Create_DetectsConflict(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)

	repo := &mockAppointmentRepo{
		findOverlappingFunc: func(ctx context.Context, date time.Time, start, end, excludeID string) ([]*appointment.Appointment, error) {
			return []*appointment.Appointment{{ID: "existing", Date: date, StartTime: "10:00", EndTime: "10:50"}}, nil
		},
	}
	svc := NewAgendaService(repo)

	_, err := svc.CreateAppointment(context.Background(), "patient-1", "Paciente Teste", date, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")

	if err == nil {
		t.Fatal("expected error for conflict, got nil")
	}
	if !strings.Contains(err.Error(), "conflicts") {
		t.Errorf("expected error containing 'conflicts', got %v", err)
	}
}

func TestAgendaService_Create_Succeeds(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	var saved *appointment.Appointment

	repo := &mockAppointmentRepo{
		findOverlappingFunc: func(ctx context.Context, date time.Time, start, end, excludeID string) ([]*appointment.Appointment, error) {
			return nil, nil
		},
		saveFunc: func(ctx context.Context, appt *appointment.Appointment) error {
			saved = appt
			return nil
		},
	}
	svc := NewAgendaService(repo)

	appt, err := svc.CreateAppointment(context.Background(), "patient-1", "Paciente Teste", date, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.Status != appointment.AppointmentStatusScheduled {
		t.Errorf("expected status scheduled, got %s", appt.Status)
	}
	if saved == nil {
		t.Error("expected appointment to be saved")
	}
}

func TestAgendaService_Cancel_SetsStatusCancelled(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	var updated *appointment.Appointment

	repo := &mockAppointmentRepo{
		findByIDFunc: func(ctx context.Context, id string) (*appointment.Appointment, error) {
			return &appointment.Appointment{
				ID:     id,
				Date:   date,
				Status: appointment.AppointmentStatusScheduled,
			}, nil
		},
		updateFunc: func(ctx context.Context, appt *appointment.Appointment) error {
			updated = appt
			return nil
		},
	}
	svc := NewAgendaService(repo)

	err := svc.CancelAppointment(context.Background(), "appt-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Error("expected appointment to be updated")
	}
	if updated.Status != appointment.AppointmentStatusCancelled {
		t.Errorf("expected status cancelled, got %s", updated.Status)
	}
}

func TestAgendaService_Complete_SetsStatusCompleted(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	var updated *appointment.Appointment

	repo := &mockAppointmentRepo{
		findByIDFunc: func(ctx context.Context, id string) (*appointment.Appointment, error) {
			return &appointment.Appointment{
				ID:     id,
				Date:   date,
				Status: appointment.AppointmentStatusConfirmed,
			}, nil
		},
		updateFunc: func(ctx context.Context, appt *appointment.Appointment) error {
			updated = appt
			return nil
		},
	}
	svc := NewAgendaService(repo)

	err := svc.CompleteAppointment(context.Background(), "appt-1", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Error("expected appointment to be updated")
	}
	if updated.Status != appointment.AppointmentStatusCompleted {
		t.Errorf("expected status completed, got %s", updated.Status)
	}
}

func TestAgendaService_Complete_LinksSession(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	var updated *appointment.Appointment

	repo := &mockAppointmentRepo{
		findByIDFunc: func(ctx context.Context, id string) (*appointment.Appointment, error) {
			return &appointment.Appointment{
				ID:     id,
				Date:   date,
				Status: appointment.AppointmentStatusConfirmed,
			}, nil
		},
		updateFunc: func(ctx context.Context, appt *appointment.Appointment) error {
			updated = appt
			return nil
		},
	}
	svc := NewAgendaService(repo)

	err := svc.CompleteAppointment(context.Background(), "appt-1", "session-xyz")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Fatal("expected appointment to be updated")
	}
	if updated.SessionID == nil || *updated.SessionID != "session-xyz" {
		t.Errorf("expected SessionID to be session-xyz, got %v", updated.SessionID)
	}
}

func TestAgendaService_ConfirmAppointment_SetsStatusConfirmed(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	var updated *appointment.Appointment

	repo := &mockAppointmentRepo{
		findByIDFunc: func(ctx context.Context, id string) (*appointment.Appointment, error) {
			return &appointment.Appointment{
				ID:     id,
				Date:   date,
				Status: appointment.AppointmentStatusScheduled,
			}, nil
		},
		updateFunc: func(ctx context.Context, appt *appointment.Appointment) error {
			updated = appt
			return nil
		},
	}
	svc := NewAgendaService(repo)

	err := svc.ConfirmAppointment(context.Background(), "appt-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Error("expected appointment to be updated")
	}
	if updated.Status != appointment.AppointmentStatusConfirmed {
		t.Errorf("expected status confirmed, got %s", updated.Status)
	}
}

func TestAgendaService_ConfirmAppointment_NotFound(t *testing.T) {
	repo := &mockAppointmentRepo{
		findByIDFunc: func(ctx context.Context, id string) (*appointment.Appointment, error) {
			return nil, nil
		},
	}
	svc := NewAgendaService(repo)

	err := svc.ConfirmAppointment(context.Background(), "appt-1")

	if err == nil {
		t.Error("expected error for not found appointment")
	}
}

func TestAgendaService_MarkNoShow_SetsStatusNoShow(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	var updated *appointment.Appointment

	repo := &mockAppointmentRepo{
		findByIDFunc: func(ctx context.Context, id string) (*appointment.Appointment, error) {
			return &appointment.Appointment{
				ID:     id,
				Date:   date,
				Status: appointment.AppointmentStatusConfirmed,
			}, nil
		},
		updateFunc: func(ctx context.Context, appt *appointment.Appointment) error {
			updated = appt
			return nil
		},
	}
	svc := NewAgendaService(repo)

	err := svc.MarkNoShow(context.Background(), "appt-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Error("expected appointment to be updated")
	}
	if updated.Status != appointment.AppointmentStatusNoShow {
		t.Errorf("expected status no_show, got %s", updated.Status)
	}
}