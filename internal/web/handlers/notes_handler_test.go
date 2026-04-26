package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/timeline"
)

type mockNotesPatientService struct {
	patients []*patient.Patient
}

func (m *mockNotesPatientService) GetPatientByID(ctx context.Context, id string) (*patient.Patient, error) {
	for _, p := range m.patients {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockNotesPatientService) ListPatients(ctx context.Context) ([]*patient.Patient, error) {
	return m.patients, nil
}

func (m *mockNotesPatientService) ListPatientsPaginated(ctx context.Context, page, pageSize int) ([]*patient.Patient, int, error) {
	return m.patients, len(m.patients), nil
}

func (m *mockNotesPatientService) CreatePatient(ctx context.Context, input services.CreatePatientInput) (*patient.Patient, error) {
	return nil, nil
}

func (m *mockNotesPatientService) SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error) {
	return nil, nil
}

func (m *mockNotesPatientService) GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error) {
	return nil, nil
}

func (m *mockNotesPatientService) ListForDashboard(ctx context.Context, limit int) ([]*patient.DashboardSummary, error) {
	return nil, nil
}

type mockNotesTimelineService struct {
	events timeline.Timeline
}

func (m *mockNotesTimelineService) GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	return m.events, nil
}

func (m *mockNotesTimelineService) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	return nil, nil
}

func TestNotesHandlerIndex(t *testing.T) {
	patients := []*patient.Patient{
		{ID: "p001", Name: "Ana Silva", Tag: "Ansiedade", UpdatedAt: time.Now()},
		{ID: "p002", Name: "Bruno Costa", Tag: "", UpdatedAt: time.Now()},
	}

	events := timeline.Timeline{
		{ID: "e1", Type: timeline.EventTypeSession, Content: "Sessão 1", Date: time.Now()},
		{ID: "e2", Type: timeline.EventTypeObservation, Content: "Observação 1", Date: time.Now()},
	}

	ps := &mockNotesPatientService{patients: patients}
	ts := &mockNotesTimelineService{events: events}

	handler := NewNotesHandler(ps, ts)

	req := httptest.NewRequest("GET", "/notes", nil)
	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotesHandlerIndexWithFilter(t *testing.T) {
	patients := []*patient.Patient{
		{ID: "p001", Name: "Ana Silva", Tag: "Ansiedade"},
		{ID: "p002", Name: "Bruno Costa", Tag: "Burnout"},
	}

	ps := &mockNotesPatientService{patients: patients}
	ts := &mockNotesTimelineService{events: nil}

	handler := NewNotesHandler(ps, ts)

	req := httptest.NewRequest("GET", "/notes?filter=Ansiedade", nil)
	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotesHandlerDetail(t *testing.T) {
	patients := []*patient.Patient{
		{ID: "p001", Name: "Ana Silva"},
	}

	events := timeline.Timeline{
		{ID: "e1", Type: timeline.EventTypeSession, Content: "Sessão de acompanhamento", Date: time.Now()},
		{ID: "e2", Type: timeline.EventTypeObservation, Content: "Paciente apresenta melhora", Date: time.Now()},
	}

	ps := &mockNotesPatientService{patients: patients}
	ts := &mockNotesTimelineService{events: events}

	handler := NewNotesHandler(ps, ts)

	req := httptest.NewRequest("GET", "/notes/detail/p001?tab=evolucao", nil)
	w := httptest.NewRecorder()

	handler.Detail(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotesHandlerDetailWithTab(t *testing.T) {
	patients := []*patient.Patient{
		{ID: "p001", Name: "Ana Silva"},
	}

	ps := &mockNotesPatientService{patients: patients}
	ts := &mockNotesTimelineService{events: nil}

	handler := NewNotesHandler(ps, ts)

	req := httptest.NewRequest("GET", "/notes/detail/p001?tab=observacoes", nil)
	w := httptest.NewRecorder()

	handler.Detail(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}