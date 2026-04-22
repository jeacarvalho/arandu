package appointment

import (
	"errors"
	"testing"
	"time"
)

func TestNewAppointment_ValidSession(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, err := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.Status != AppointmentStatusScheduled {
		t.Errorf("expected status scheduled, got %s", appt.Status)
	}
	if appt.PatientID != "patient-1" {
		t.Errorf("expected patient ID patient-1, got %s", appt.PatientID)
	}
}

func TestNewAppointment_InvalidDuration(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name     string
		duration int
	}{
		{"duration_below_30", 20},
		{"duration_above_120", 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", tt.duration, AppointmentTypeSession, "")
			if !errors.Is(err, ErrInvalidDuration) {
				t.Errorf("expected ErrInvalidDuration, got %v", err)
			}
		})
	}
}

func TestNewAppointment_MissingPatient(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	_, err := NewAppointment("", "", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	if !errors.Is(err, ErrPatientRequired) {
		t.Errorf("expected ErrPatientRequired, got %v", err)
	}
}

func TestNewAppointment_BlockedSlotNoPatient(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, err := NewAppointment("", "", date, "12:00", "13:00", 60, AppointmentTypeBlocked, "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.PatientID != "" {
		t.Errorf("expected empty patient ID for blocked slot, got %s", appt.PatientID)
	}
	if appt.Type != AppointmentTypeBlocked {
		t.Errorf("expected type blocked, got %s", appt.Type)
	}
}

func TestNewAppointment_Confirm_SetsStatus(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	err := appt.Confirm()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.Status != AppointmentStatusConfirmed {
		t.Errorf("expected status confirmed, got %s", appt.Status)
	}
}

func TestAppointment_Cancel_SetsStatus(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	err := appt.Cancel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.Status != AppointmentStatusCancelled {
		t.Errorf("expected status cancelled, got %s", appt.Status)
	}
	if appt.CancelledAt == nil {
		t.Error("expected CancelledAt to be set")
	}
}

func TestAppointment_Cancel_BlocksCompleted(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")
	appt.Complete()

	err := appt.Cancel()
	if !errors.Is(err, ErrAlreadyCompleted) {
		t.Errorf("expected ErrAlreadyCompleted, got %v", err)
	}
}

func TestAppointment_Complete_SetsStatus(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	err := appt.Complete()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.Status != AppointmentStatusCompleted {
		t.Errorf("expected status completed, got %s", appt.Status)
	}
	if appt.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestAppointment_Complete_BlocksCancelled(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")
	appt.Cancel()

	err := appt.Complete()
	if !errors.Is(err, ErrAlreadyCancelled) {
		t.Errorf("expected ErrAlreadyCancelled, got %v", err)
	}
}

func TestAppointment_MarkNoShow(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	err := appt.MarkNoShow()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if appt.Status != AppointmentStatusNoShow {
		t.Errorf("expected status no_show, got %s", appt.Status)
	}
}

func TestAppointment_Overlaps_SameTimeTrue(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt1, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")
	appt2, _ := NewAppointment("patient-2", "Outro Paciente", date, "10:00", "10:50", 50, AppointmentTypeSession, "")

	if !appt1.Overlaps(appt2) {
		t.Error("expected overlaps to be true for same time")
	}
}

func TestAppointment_Overlaps_AdjacentFalse(t *testing.T) {
	date := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	appt1, _ := NewAppointment("patient-1", "Paciente Teste", date, "10:00", "10:50", 50, AppointmentTypeSession, "")
	appt2, _ := NewAppointment("patient-2", "Outro Paciente", date, "10:50", "11:40", 50, AppointmentTypeSession, "")

	if appt1.Overlaps(appt2) {
		t.Error("expected overlaps to be false for adjacent times")
	}
}

func TestAppointment_Overlaps_DifferentDateFalse(t *testing.T) {
	date1 := time.Date(2026, 4, 22, 0, 0, 0, 0, time.Local)
	date2 := time.Date(2026, 4, 23, 0, 0, 0, 0, time.Local)
	appt1, _ := NewAppointment("patient-1", "Paciente Teste", date1, "10:00", "10:50", 50, AppointmentTypeSession, "")
	appt2, _ := NewAppointment("patient-2", "Outro Paciente", date2, "10:00", "10:50", 50, AppointmentTypeSession, "")

	if appt1.Overlaps(appt2) {
		t.Error("expected overlaps to be false for different dates")
	}
}