package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"arandu/internal/domain/appointment"
)

func setupAppointmentTestDB(t *testing.T) (*DB, func()) {
	tmpfile, err := os.CreateTemp("", "testdb-appt-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func TestAppointmentRepository_SaveAndFindByID(t *testing.T) {
	db, cleanup := setupAppointmentTestDB(t)
	defer cleanup()

	repo := NewAppointmentRepository(db)
	ctx := context.Background()

	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, err := appointment.NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")
	if err != nil {
		t.Fatalf("Failed to create appointment: %v", err)
	}

	if err := repo.Save(ctx, appt); err != nil {
		t.Fatalf("Failed to save appointment: %v", err)
	}

	found, err := repo.FindByID(ctx, appt.ID)
	if err != nil {
		t.Fatalf("Failed to find appointment: %v", err)
	}
	if found == nil {
		t.Fatal("Appointment not found")
	}
	if found.PatientName != "Paciente Teste" {
		t.Errorf("Expected PatientName 'Paciente Teste', got %s", found.PatientName)
	}
	if found.Status != appointment.AppointmentStatusScheduled {
		t.Errorf("Expected status scheduled, got %s", found.Status)
	}
}

func TestAppointmentRepository_FindByDateRange(t *testing.T) {
	db, cleanup := setupAppointmentTestDB(t)
	defer cleanup()

	repo := NewAppointmentRepository(db)
	ctx := context.Background()

	date1 := time.Date(2026, 4, 21, 0, 0, 0, 0, time.Local)
	date2 := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	date3 := time.Date(2026, 4, 23, 0, 0, 0, 0, time.Local)

	for _, d := range []time.Time{date1, date2, date3} {
		appt, _ := appointment.NewAppointment("patient-1", "Paciente Teste", d, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")
		repo.Save(ctx, appt)
	}

	results, err := repo.FindByDateRange(ctx, date2, date2)
	if err != nil {
		t.Fatalf("Failed to find appointments: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 appointment, got %d", len(results))
	}
}

func TestAppointmentRepository_FindOverlapping_Conflict(t *testing.T) {
	db, cleanup := setupAppointmentTestDB(t)
	defer cleanup()

	repo := NewAppointmentRepository(db)
	ctx := context.Background()

	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	existing, _ := appointment.NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")
	repo.Save(ctx, existing)

	conflicts, err := repo.FindOverlapping(ctx, date, "10:20", "11:10", "")
	if err != nil {
		t.Fatalf("Failed to find overlapping: %v", err)
	}
	if len(conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(conflicts))
	}

	noConflicts, err := repo.FindOverlapping(ctx, date, "11:00", "11:50", "")
	if err != nil {
		t.Fatalf("Failed to find overlapping: %v", err)
	}
	if len(noConflicts) != 0 {
		t.Errorf("Expected 0 conflicts, got %d", len(noConflicts))
	}
}

func TestAppointmentRepository_Update_StatusChange(t *testing.T) {
	db, cleanup := setupAppointmentTestDB(t)
	defer cleanup()

	repo := NewAppointmentRepository(db)
	ctx := context.Background()

	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := appointment.NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")
	repo.Save(ctx, appt)

	appt.Cancel()
	repo.Update(ctx, appt)

	found, _ := repo.FindByID(ctx, appt.ID)
	if found.Status != appointment.AppointmentStatusCancelled {
		t.Errorf("Expected status cancelled, got %s", found.Status)
	}
}

func TestAppointmentRepository_FindByDate(t *testing.T) {
	db, cleanup := setupAppointmentTestDB(t)
	defer cleanup()

	repo := NewAppointmentRepository(db)
	ctx := context.Background()

	date1 := time.Date(2026, 4, 21, 0, 0, 0, 0, time.Local)
	date2 := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)

	appt1, _ := appointment.NewAppointment("patient-1", "Paciente Teste", date1, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")
	appt2, _ := appointment.NewAppointment("patient-2", "Outro Paciente", date2, "10:00", "10:50", 50, appointment.AppointmentTypeSession, "")
	repo.Save(ctx, appt1)
	repo.Save(ctx, appt2)

	results, err := repo.FindByDate(ctx, date1)
	if err != nil {
		t.Fatalf("Failed to find appointments: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 appointment, got %d", len(results))
	}
}