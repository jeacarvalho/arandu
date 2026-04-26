package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/timeline"
)

type mockTimelineService struct {
	events        timeline.Timeline
	searchResults []*timeline.SearchResult
}

func (m *mockTimelineService) GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	// Simular paginação
	start := offset
	if start > len(m.events) {
		start = len(m.events)
	}
	end := start + limit
	if end > len(m.events) {
		end = len(m.events)
	}
	return m.events[start:end], nil
}

func (m *mockTimelineService) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	return m.searchResults, nil
}

type mockNotesPatientServiceForTimeline struct {
	patient *patient.Patient
}

func (m *mockNotesPatientServiceForTimeline) GetPatientByID(ctx context.Context, id string) (*patient.Patient, error) {
	if m.patient == nil {
		return &patient.Patient{ID: id, Name: "Test Patient"}, nil
	}
	return m.patient, nil
}

func (m *mockNotesPatientServiceForTimeline) ListPatients(ctx context.Context) ([]*patient.Patient, error) {
	return nil, nil
}

func (m *mockNotesPatientServiceForTimeline) ListPatientsPaginated(ctx context.Context, page, pageSize int) ([]*patient.Patient, int, error) {
	return nil, 0, nil
}

func (m *mockNotesPatientServiceForTimeline) CreatePatient(ctx context.Context, input services.CreatePatientInput) (*patient.Patient, error) {
	return nil, nil
}

func (m *mockNotesPatientServiceForTimeline) SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error) {
	return nil, nil
}

func (m *mockNotesPatientServiceForTimeline) GetThemeFrequency(ctx context.Context, patientID string, limit int) ([]map[string]interface{}, error) {
	return nil, nil
}

func (m *mockNotesPatientServiceForTimeline) ListForDashboard(ctx context.Context, limit int) ([]*patient.DashboardSummary, error) {
	return nil, nil
}

func TestTimelineHandler_ShowPatientHistory(t *testing.T) {
	now := time.Now()

	service := &mockTimelineService{
		events: timeline.Timeline{
			{ID: "1", Type: timeline.EventTypeSession, Date: now, Content: "Session 1", Metadata: map[string]string{"session_id": "sess-1"}},
			{ID: "2", Type: timeline.EventTypeObservation, Date: now.Add(-time.Hour), Content: "Observation 1", Metadata: map[string]string{"session_id": "sess-1"}},
			{ID: "3", Type: timeline.EventTypeIntervention, Date: now.Add(-2 * time.Hour), Content: "Intervention 1", Metadata: map[string]string{"session_id": "sess-1"}},
		},
	}

	ps := &mockNotesPatientServiceForTimeline{}
	handler := NewTimelineHandler(service, ps)

	t.Run("GET /patients/{id}/history returns timeline", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/patient-123/history", nil)
		rec := httptest.NewRecorder()

		handler.ShowPatientHistory(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "Session 1") {
			t.Error("Expected body to contain session content")
		}
		if !strings.Contains(body, "Observation 1") {
			t.Error("Expected body to contain observation content")
		}
		if !strings.Contains(body, "Intervention 1") {
			t.Error("Expected body to contain intervention content")
		}
	})

	t.Run("GET with filter parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/patient-123/history?filter=session", nil)
		rec := httptest.NewRecorder()

		handler.ShowPatientHistory(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})

	t.Run("GET with pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/patient-123/history?offset=20", nil)
		rec := httptest.NewRecorder()

		handler.ShowPatientHistory(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})
}

func TestTimelineHandler_LoadMoreEvents(t *testing.T) {
	now := time.Now()

	// Criar 50 eventos para testar paginação
	events := make(timeline.Timeline, 50)
	for i := 0; i < 50; i++ {
		events[i] = &timeline.TimelineEvent{
			ID:       fmt.Sprintf("event-%d", i),
			Type:     timeline.EventTypeObservation,
			Date:     now.Add(-time.Duration(i) * time.Hour),
			Content:  fmt.Sprintf("Content %d", i),
			Metadata: map[string]string{"session_id": fmt.Sprintf("sess-%d", i)},
		}
	}

	service := &mockTimelineService{events: events}
	ps := &mockNotesPatientServiceForTimeline{}
	handler := NewTimelineHandler(service, ps)

	t.Run("Load more events with HTMX", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/patient-123/history/load-more?offset=20", nil)
		req.Header.Set("HX-Request", "true")
		rec := httptest.NewRecorder()

		handler.LoadMoreEvents(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		body := rec.Body.String()
		// Verificar que retornou os eventos paginados
		if !strings.Contains(body, "timeline") {
			t.Error("Expected timeline content in response")
		}
	})

	t.Run("Load more without patient ID returns 400", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients//history/load-more", nil)
		rec := httptest.NewRecorder()

		handler.LoadMoreEvents(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rec.Code)
		}
	})
}

func TestTimelineHandler_SearchPatientHistory(t *testing.T) {
	now := time.Now()

	service := &mockTimelineService{
		searchResults: []*timeline.SearchResult{
			{
				ID:        "1",
				Type:      timeline.EventTypeObservation,
				Date:      now,
				Content:   "Patient showed anxiety",
				Snippet:   "Patient showed <b>anxiety</b>",
				SessionID: "sess-1",
				PatientID: "patient-123",
			},
		},
	}

	ps := &mockNotesPatientServiceForTimeline{}
	handler := NewTimelineHandler(service, ps)

	t.Run("Search with query returns results", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/patient-123/history/search?q=anxiety", nil)
		rec := httptest.NewRecorder()

		handler.SearchPatientHistory(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "anxiety") {
			t.Error("Expected body to contain search results")
		}
	})

	t.Run("Empty query redirects to timeline", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/patient-123/history/search?q=", nil)
		rec := httptest.NewRecorder()

		handler.SearchPatientHistory(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})
}

func TestParseFilterType(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		expected *timeline.EventType
	}{
		{
			name:     "observation filter",
			filter:   "observation",
			expected: func() *timeline.EventType { f := timeline.EventTypeObservation; return &f }(),
		},
		{
			name:     "intervention filter",
			filter:   "intervention",
			expected: func() *timeline.EventType { f := timeline.EventTypeIntervention; return &f }(),
		},
		{
			name:     "session filter",
			filter:   "session",
			expected: func() *timeline.EventType { f := timeline.EventTypeSession; return &f }(),
		},
		{
			name:     "all filter",
			filter:   "all",
			expected: nil,
		},
		{
			name:     "empty filter",
			filter:   "",
			expected: nil,
		},
		{
			name:     "unknown filter",
			filter:   "unknown",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFilterType(tt.filter)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", *result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected %v, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("Expected %v, got %v", *tt.expected, *result)
				}
			}
		})
	}
}
