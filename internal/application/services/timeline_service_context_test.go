package services

import (
	"context"
	"testing"
	"time"

	"arandu/internal/domain/timeline"
)

// mockTimelineRepositoryContextAware implements TimelineRepositoryContextAware for testing
type mockTimelineRepositoryContextAware struct {
	events        timeline.Timeline
	searchResults []*timeline.SearchResult
	err           error
}

func (m *mockTimelineRepositoryContextAware) GetTimelineByPatientID(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func (m *mockTimelineRepositoryContextAware) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.searchResults, nil
}

func TestTimelineServiceContext_GetPatientTimeline(t *testing.T) {
	now := time.Now()
	mockEvents := timeline.Timeline{
		{
			ID:      "event-1",
			Type:    timeline.EventTypeSession,
			Date:    now,
			Content: "Sessão de terapia",
		},
		{
			ID:      "event-2",
			Type:    timeline.EventTypeObservation,
			Date:    now.Add(-time.Hour),
			Content: "Observação importante",
		},
	}

	mockRepo := &mockTimelineRepositoryContextAware{events: mockEvents}
	service := NewTimelineServiceContext(mockRepo)

	t.Run("Get timeline with events", func(t *testing.T) {
		events, err := service.GetPatientTimeline(context.Background(), "patient-123", nil, 10, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(events) != 2 {
			t.Errorf("Esperado 2 eventos, obtido %d", len(events))
		}
	})

	t.Run("Get timeline with filter", func(t *testing.T) {
		filter := timeline.EventTypeSession
		events, err := service.GetPatientTimelineWithFilter(context.Background(), "patient-123", filter, 10, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(events) != 2 {
			t.Errorf("Esperado 2 eventos, obtido %d", len(events))
		}
	})

	t.Run("Get timeline grouped by date", func(t *testing.T) {
		grouped, err := service.GetPatientTimelineGroupedByDate(context.Background(), "patient-123", nil, 10, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if grouped == nil {
			t.Error("Esperado mapa agrupado, obtido nil")
		}
	})

	t.Run("Search in history", func(t *testing.T) {
		mockResults := []*timeline.SearchResult{
			{
				ID:      "result-1",
				Content: "Resultado de busca",
				Type:    timeline.EventTypeObservation,
			},
		}
		mockRepo.searchResults = mockResults

		results, err := service.SearchInHistory(context.Background(), "patient-123", "terapia")

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Esperado 1 resultado, obtido %d", len(results))
		}
	})

	t.Run("Handle errors gracefully", func(t *testing.T) {
		errRepo := &mockTimelineRepositoryContextAware{err: context.DeadlineExceeded}
		errService := NewTimelineServiceContext(errRepo)

		_, err := errService.GetPatientTimeline(context.Background(), "patient-123", nil, 10, 0)

		if err == nil {
			t.Error("Esperado erro, mas não houve nenhum")
		}
	})
}

func TestTimelineServiceContext_EdgeCases(t *testing.T) {
	t.Run("Empty timeline", func(t *testing.T) {
		mockRepo := &mockTimelineRepositoryContextAware{events: timeline.Timeline{}}
		service := NewTimelineServiceContext(mockRepo)

		events, err := service.GetPatientTimeline(context.Background(), "patient-123", nil, 10, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(events) != 0 {
			t.Errorf("Esperado 0 eventos, obtido %d", len(events))
		}
	})

	t.Run("Nil filter type", func(t *testing.T) {
		mockRepo := &mockTimelineRepositoryContextAware{events: timeline.Timeline{}}
		service := NewTimelineServiceContext(mockRepo)

		events, err := service.GetPatientTimeline(context.Background(), "patient-123", nil, 10, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if events == nil {
			t.Error("Esperado lista vazia, não nil")
		}
	})
}
