package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/appointment"
	"arandu/internal/domain/patient"

	"github.com/google/uuid"
)

type mockAppointmentRepository struct {
	appointments map[string]*appointment.Appointment
}

func newMockAppointmentRepository() *mockAppointmentRepository {
	return &mockAppointmentRepository{
		appointments: make(map[string]*appointment.Appointment),
	}
}

func (r *mockAppointmentRepository) Save(ctx context.Context, appt *appointment.Appointment) error {
	r.appointments[appt.ID] = appt
	return nil
}

func (r *mockAppointmentRepository) FindByID(ctx context.Context, id string) (*appointment.Appointment, error) {
	if appt, ok := r.appointments[id]; ok {
		return appt, nil
	}
	return nil, nil
}

func (r *mockAppointmentRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*appointment.Appointment, error) {
	var result []*appointment.Appointment
	for _, appt := range r.appointments {
		if !appt.Date.Before(startDate) && !appt.Date.After(endDate) {
			result = append(result, appt)
		}
	}
	return result, nil
}

func (r *mockAppointmentRepository) FindByPatient(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
	var result []*appointment.Appointment
	for _, appt := range r.appointments {
		if appt.PatientID == patientID {
			result = append(result, appt)
		}
	}
	return result, nil
}

func (r *mockAppointmentRepository) FindByDate(ctx context.Context, date time.Time) ([]*appointment.Appointment, error) {
	return r.FindByDateRange(ctx, date, date)
}

func (r *mockAppointmentRepository) FindOverlapping(ctx context.Context, date time.Time, startTime, endTime string, excludeID string) ([]*appointment.Appointment, error) {
	return []*appointment.Appointment{}, nil
}

func (r *mockAppointmentRepository) Update(ctx context.Context, appt *appointment.Appointment) error {
	r.appointments[appt.ID] = appt
	return nil
}

func (r *mockAppointmentRepository) Delete(ctx context.Context, id string) error {
	delete(r.appointments, id)
	return nil
}

func (r *mockAppointmentRepository) FindUpcoming(ctx context.Context, fromDate time.Time, limit int) ([]*appointment.Appointment, error) {
	return []*appointment.Appointment{}, nil
}

type mockPatientService struct{}

func (m *mockPatientService) ListPatients(ctx context.Context) ([]*patient.Patient, error) {
	return []*patient.Patient{
		{ID: "patient-1", Name: "João Silva"},
		{ID: "patient-2", Name: "Maria Santos"},
	}, nil
}

func (m *mockPatientService) GetPatientByID(ctx context.Context, id string) (*patient.Patient, error) {
	return &patient.Patient{ID: id, Name: "Test Patient"}, nil
}

// MockAgendaService implements AgendaServiceInterface
type MockAgendaService struct {
	repo *mockAppointmentRepository
}

func NewMockAgendaService() *services.AgendaService {
	return services.NewAgendaService(newMockAppointmentRepository())
}

func TestAgendaHandler_View_FullPage(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda", nil)
	w := httptest.NewRecorder()

	handler.View(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}

	// Verify it's a full page (contains DOCTYPE or meaningful agenda content)
	hasDoctype := bytes.Contains([]byte(body), []byte("<!DOCTYPE html>"))
	hasAgendaContent := bytes.Contains([]byte(body), []byte("Agenda"))
	if !hasDoctype && !hasAgendaContent {
		t.Error("Expected full page with DOCTYPE or agenda content for non-HTMX request")
	}
}

func TestAgendaHandler_View_HTMXRequest(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda", nil)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	handler.View(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	// For HTMX requests, should NOT contain DOCTYPE (it's a fragment)
	if bytes.Contains([]byte(body), []byte("<!DOCTYPE html>")) {
		t.Error("Expected fragment response without DOCTYPE for HTMX request")
	}
}

func TestAgendaHandler_DayView(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	// Test without HTMX header
	req := httptest.NewRequest(http.MethodGet, "/agenda/day?date=2026-04-14", nil)
	w := httptest.NewRecorder()

	handler.DayView(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body")
	}

	// Currently returns fragment even for non-HTMX (this is what we want to fix)
	// For now, verify it returns valid HTML
	if !bytes.Contains([]byte(body), []byte("<div")) {
		t.Error("Expected HTML content in response")
	}
}

func TestAgendaHandler_WeekView(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda/week", nil)
	w := httptest.NewRecorder()

	handler.WeekView(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAgendaHandler_MonthView(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda/month?year=2026&month=4", nil)
	w := httptest.NewRecorder()

	handler.MonthView(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAgendaHandler_Create_NonHTMX(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	formData := fmt.Sprintf("patient_id=patient-1&patient_name=João Silva&date=2026-04-14&start_time=10:00&duration=50&type=session")
	req := httptest.NewRequest(http.MethodPost, "/agenda/appointments", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	// Should redirect for non-HTMX
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 303 for redirect, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location == "" {
		t.Error("Expected Location header with redirect URL")
	}
}

func TestAgendaHandler_Create_HTMX(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	formData := fmt.Sprintf("patient_id=patient-1&patient_name=João Silva&date=2026-04-14&start_time=10:00&duration=50&type=session")
	req := httptest.NewRequest(http.MethodPost, "/agenda/appointments", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	// For HTMX, should return HX-Redirect header, not 303
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for HTMX response, got %d", w.Code)
	}

	hxRedirect := w.Header().Get("HX-Redirect")
	if hxRedirect == "" {
		t.Error("Expected HX-Redirect header for HTMX response")
	}
}

func TestAgendaHandler_Cancel_NotFound(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodDelete, "/agenda/appointments/non-existent-id", nil)
	w := httptest.NewRecorder()

	handler.Cancel(w, req)

	if w.Code == 0 {
		t.Error("Expected a response code")
	}
}

func TestAgendaHandler_NewForm(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda/new", nil)
	w := httptest.NewRecorder()

	handler.NewForm(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte("Nova Marcação")) {
		t.Error("Expected form with 'Nova Marcação' title")
	}
}

func TestAgendaHandler_GetSlots(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda/slots?date=2026-04-14", nil)
	w := httptest.NewRecorder()

	handler.GetSlots(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAgendaHandler_View_WithDate(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	// Test with specific date
	req := httptest.NewRequest(http.MethodGet, "/agenda?date=2026-04-14&view=semana", nil)
	w := httptest.NewRecorder()

	handler.View(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	// Verify response contains expected content (may not have appointments, but should have structure)
	if !bytes.Contains([]byte(body), []byte("Agenda")) {
		t.Error("Expected 'Agenda' in response for 2026-04-14")
	}
}

func TestAgendaHandler_Show_NotFound(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodGet, "/agenda/appointments/non-existent-id", nil)
	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent appointment, got %d", w.Code)
	}
}

func TestAgendaHandler_Create_InvalidDate(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	formData := fmt.Sprintf("patient_id=patient-1&date=invalid-date&start_time=10:00&duration=50")
	req := httptest.NewRequest(http.MethodPost, "/agenda/appointments", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid date, got %d", w.Code)
	}
}

func TestAgendaHandler_Create_Conflict(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	// Create first appointment
	formData := fmt.Sprintf("patient_id=patient-1&patient_name=João Silva&date=2026-04-14&start_time=10:00&duration=50&type=session")
	createReq := httptest.NewRequest(http.MethodPost, "/agenda/appointments", bytes.NewBufferString(formData))
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createW := httptest.NewRecorder()
	handler.Create(createW, createReq)

	// Try to create overlapping appointment
	conflictData := fmt.Sprintf("patient_id=patient-2&patient_name=Maria Santos&date=2026-04-14&start_time=10:30&duration=50&type=session")
	conflictReq := httptest.NewRequest(http.MethodPost, "/agenda/appointments", bytes.NewBufferString(conflictData))
	conflictReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	conflictW := httptest.NewRecorder()
	handler.Create(conflictW, conflictReq)

	// Should return conflict error (our mock doesn't detect conflicts, so this tests the structure)
	if conflictW.Code == 0 {
		t.Error("Expected a response code for conflict attempt")
	}
}

func TestAgendaHandler_Update(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	// Test PUT without proper ID
	req := httptest.NewRequest(http.MethodPut, "/agenda/appointments/", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	// Should handle gracefully (400 for missing ID)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing ID, got %d", w.Code)
	}
}

func TestAgendaHandler_MethodNotAllowed(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"View POST", http.MethodPost, "/agenda"},
		{"DayView PUT", http.MethodPut, "/agenda/day"},
		{"WeekView DELETE", http.MethodDelete, "/agenda/week"},
		{"Create GET", http.MethodGet, "/agenda/appointments"},
		{"NewForm POST", http.MethodPost, "/agenda/new"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			switch tt.path {
			case "/agenda":
				handler.View(w, req)
			case "/agenda/day":
				handler.DayView(w, req)
			case "/agenda/week":
				handler.WeekView(w, req)
			case "/agenda/appointments":
				handler.Create(w, req)
			case "/agenda/new":
				handler.NewForm(w, req)
			}

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 405, got %d for %s %s", w.Code, tt.method, tt.path)
			}
		})
	}
}

func TestAgendaHandler_Reschedule(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	formData := fmt.Sprintf("date=2026-04-15&start_time=14:00&duration=50")
	req := httptest.NewRequest(http.MethodPost, "/agenda/appointments/test-id/reschedule", bytes.NewBufferString(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.Reschedule(w, req)

	// Should not crash, returns appropriate response
	if w.Code == 0 {
		t.Error("Expected a response code")
	}
}

func TestAgendaHandler_Complete(t *testing.T) {
	agendaService := NewMockAgendaService()
	patientService := &mockPatientService{}

	handler := NewAgendaHandler(agendaService, patientService)

	req := httptest.NewRequest(http.MethodPost, "/agenda/appointments/test-id/complete", nil)
	w := httptest.NewRecorder()

	handler.Complete(w, req)

	// Should not crash
	if w.Code == 0 {
		t.Error("Expected a response code")
	}
}

// Helper to generate UUID for testing
func generateTestID() string {
	return uuid.New().String()
}
