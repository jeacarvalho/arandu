package agenda_test

import (
	"bytes"
	"testing"
	"time"

	"arandu/web/components/agenda"
)

func strPtr(s string) *string {
	return &s
}

func baseDetailModel(id string) agenda.AppointmentDetailModel {
	return agenda.AppointmentDetailModel{
		ID:          id,
		PatientName: "Paciente Teste",
		Date:        time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		StartTime:   "10:00",
		EndTime:     "10:50",
		Duration:    50,
		SessionType: "Sessão individual",
		Status:      "scheduled",
		SessionID:   nil,
	}
}

func TestAppointmentDetail_ScheduledShowsConfirmButton(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "scheduled"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`hx-post="/agenda/appointments/test-id/confirm"`)) {
		t.Errorf("expected HTML to contain confirm button for scheduled status\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Confirmar")) {
		t.Errorf("expected HTML to contain 'Confirmar' text\nHTML: %s", html)
	}
}

func TestAppointmentDetail_ScheduledHasReagendarButton(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "scheduled"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`hx-get="/agenda/appointments/test-id/reschedule-form"`)) {
		t.Errorf("expected HTML to contain reschedule-form button for scheduled\nHTML: %s", html)
	}
}

func TestAppointmentDetail_ConfirmedShowsConcluirPanel(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "confirmed"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`id="confirm-panel-test-id"`)) {
		t.Errorf("expected HTML to contain confirm-panel for confirmed\nHTML: %s", html)
	}
}

func TestAppointmentDetail_ConfirmedWithSessionShowsVerSessao(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "confirmed"
	model.SessionID = strPtr("session-abc")

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("/session/session-abc")) {
		t.Errorf("expected HTML to contain session link\nHTML: %s", html)
	}
}

func TestAppointmentDetail_CancelledStatus_NoActionButtons(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "cancelled"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`hx-post="/agenda/appointments/test-id/cancel"`)) {
		t.Errorf("expected HTML NOT to contain cancel button for cancelled status\nHTML: %s", html)
	}
}

func TestAppointmentDetail_ScheduledDoesNotShowFaltouButton(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "scheduled"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`/no-show`)) {
		t.Errorf("expected HTML NOT to contain no-show button for scheduled status\nHTML: %s", html)
	}
	if bytes.Contains(buf.Bytes(), []byte("Faltou")) {
		t.Errorf("expected HTML NOT to contain 'Faltou' for scheduled status\nHTML: %s", html)
	}
}

func TestAppointmentDetail_ConfirmedShowsFaltouButton(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "confirmed"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`hx-post="/agenda/appointments/test-id/no-show"`)) {
		t.Errorf("expected HTML to contain no-show button for confirmed status\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Faltou")) {
		t.Errorf("expected HTML to contain 'Faltou' text\nHTML: %s", html)
	}
}

func TestAppointmentDetail_ConfirmedDoesNotShowConfirmButton(t *testing.T) {
	model := baseDetailModel("test-id")
	model.Status = "confirmed"

	var buf bytes.Buffer
	err := agenda.AppointmentDetail(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`/confirm`)) {
		t.Errorf("expected HTML NOT to contain confirm button for confirmed status\nHTML: %s", html)
	}
}

func TestRescheduleForm_FormHasCorrectHTMXPost(t *testing.T) {
	model := agenda.RescheduleFormModel{
		AppointmentID: "appt-xyz",
		CurrentDate:   time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		CurrentStart:  "10:00",
		Duration:     50,
	}

	var buf bytes.Buffer
	err := agenda.RescheduleForm(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`hx-post="/agenda/appointments/appt-xyz/reschedule"`)) {
		t.Errorf("expected HTML to contain hx-post route\nHTML: %s", html)
	}
	if bytes.Contains(buf.Bytes(), []byte("htmx-post")) {
		t.Errorf("expected HTML NOT to contain htmx-post typo\nHTML: %s", html)
	}
}

func TestRescheduleForm_FormTarget(t *testing.T) {
	model := agenda.RescheduleFormModel{
		AppointmentID: "appt-xyz",
		CurrentDate:   time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		CurrentStart:  "10:00",
		Duration:     50,
	}

	var buf bytes.Buffer
	err := agenda.RescheduleForm(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`hx-target="#agenda-content"`)) {
		t.Errorf("expected HTML to contain hx-target\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`hx-swap="outerHTML"`)) {
		t.Errorf("expected HTML to contain hx-swap\nHTML: %s", html)
	}
}

func TestRescheduleForm_HasRequiredInputs(t *testing.T) {
	model := agenda.RescheduleFormModel{
		AppointmentID: "appt-xyz",
		CurrentDate:   time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		CurrentStart:  "10:00",
		Duration:     50,
	}

	var buf bytes.Buffer
	err := agenda.RescheduleForm(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`name="date"`)) {
		t.Errorf("expected HTML to contain date input\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`name="start_time"`)) {
		t.Errorf("expected HTML to contain start_time input\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`name="duration"`)) {
		t.Errorf("expected HTML to contain duration input\nHTML: %s", html)
	}
}

func TestRescheduleForm_PreFillsCurrentValues(t *testing.T) {
	model := agenda.RescheduleFormModel{
		AppointmentID: "appt-xyz",
		CurrentDate:   time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		CurrentStart:  "10:00",
		Duration:     50,
	}

	var buf bytes.Buffer
	err := agenda.RescheduleForm(model).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`value="2026-04-22"`)) {
		t.Errorf("expected HTML to contain current date value\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`value="10:00"`)) {
		t.Errorf("expected HTML to contain current start time value\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`value="50"`)) {
		t.Errorf("expected HTML to contain duration 50 as selected option\nHTML: %s", html)
	}
}

func TestNewAppointmentForm_FormHasCorrectHTMXPost(t *testing.T) {
	data := agenda.NewAppointmentFormData{
		Date:     time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		Slots:    nil,
		Patients: nil,
	}

	var buf bytes.Buffer
	err := agenda.NewAppointmentForm(data).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`hx-post="/agenda/appointments"`)) {
		t.Errorf("expected HTML to contain hx-post route\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`hx-target="#agenda-content"`)) {
		t.Errorf("expected HTML to contain hx-target\nHTML: %s", html)
	}
}

func TestNewAppointmentForm_ConflictWarningHiddenByDefault(t *testing.T) {
	data := agenda.NewAppointmentFormData{
		Date:     time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		Slots:    nil,
		Patients: nil,
	}

	var buf bytes.Buffer
	err := agenda.NewAppointmentForm(data).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`id="conflict-warning"`)) {
		t.Errorf("expected HTML to contain conflict-warning id\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte("hidden")) {
		t.Errorf("expected HTML to contain hidden class\nHTML: %s", html)
	}
}

func TestNewAppointmentForm_HasRequiredInputs(t *testing.T) {
	data := agenda.NewAppointmentFormData{
		Date:  time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
		Slots: []agenda.TimeSlotViewModel{
			{StartTime: "10:00", EndTime: "10:50", Available: true},
		},
		Patients: nil,
	}

	var buf bytes.Buffer
	err := agenda.NewAppointmentForm(data).Render(t.Context(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	html := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte(`name="date"`)) {
		t.Errorf("expected HTML to contain date input\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`name="start_time"`)) {
		t.Errorf("expected HTML to contain start_time input (as radio button)\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`name="duration"`)) {
		t.Errorf("expected HTML to contain duration input\nHTML: %s", html)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`name="patient_id"`)) {
		t.Errorf("expected HTML to contain patient_id input\nHTML: %s", html)
	}
}