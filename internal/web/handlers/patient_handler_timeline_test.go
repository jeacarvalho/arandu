package handlers

import (
	"context"
	"testing"
	"time"

	"arandu/internal/domain/timeline"
)

// mockTimelineServicePort implementa TimelineServicePort para testes
type mockTimelineServicePort struct {
	events timeline.Timeline
	err    error
}

func (m *mockTimelineServicePort) GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	if m.err != nil {
		return nil, m.err
	}
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

func TestPatientHandler_TimelineServicePortInterface(t *testing.T) {
	t.Run("TimelineServicePort interface is satisfied by mock", func(t *testing.T) {
		// Verificar que mockTimelineServicePort implementa TimelineServicePort
		var port TimelineServicePort = &mockTimelineServicePort{}
		if port == nil {
			t.Error("mockTimelineServicePort não implementa TimelineServicePort")
		}
	})

	t.Run("TimelineServicePort returns events correctly", func(t *testing.T) {
		now := time.Now()
		events := timeline.Timeline{
			{
				ID:       "event-1",
				Type:     timeline.EventTypeSession,
				Date:     now,
				Content:  "Sessão de terapia",
				Metadata: map[string]string{"session_id": "sess-1"},
			},
			{
				ID:       "event-2",
				Type:     timeline.EventTypeObservation,
				Date:     now.Add(-time.Hour),
				Content:  "Observação importante",
				Metadata: map[string]string{"session_id": "sess-1"},
			},
		}

		mock := &mockTimelineServicePort{events: events}
		result, err := mock.GetPatientTimeline(context.Background(), "patient-123", nil, 5, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Esperado 2 eventos, obtido %d", len(result))
		}

		if result[0].Type != timeline.EventTypeSession {
			t.Errorf("Esperado tipo Session, obtido %s", result[0].Type)
		}
	})

	t.Run("TimelineServicePort handles errors gracefully", func(t *testing.T) {
		mock := &mockTimelineServicePort{err: context.DeadlineExceeded}
		_, err := mock.GetPatientTimeline(context.Background(), "patient-123", nil, 5, 0)

		if err == nil {
			t.Error("Esperado erro, mas não houve nenhum")
		}
	})

	t.Run("TimelineServicePort respects limit and offset", func(t *testing.T) {
		now := time.Now()
		events := timeline.Timeline{
			{ID: "1", Date: now, Type: timeline.EventTypeSession},
			{ID: "2", Date: now.Add(-time.Hour), Type: timeline.EventTypeObservation},
			{ID: "3", Date: now.Add(-2 * time.Hour), Type: timeline.EventTypeIntervention},
			{ID: "4", Date: now.Add(-3 * time.Hour), Type: timeline.EventTypeSession},
			{ID: "5", Date: now.Add(-4 * time.Hour), Type: timeline.EventTypeObservation},
		}

		mock := &mockTimelineServicePort{events: events}

		// Buscar apenas 2 eventos a partir do offset 1
		result, err := mock.GetPatientTimeline(context.Background(), "patient-123", nil, 2, 1)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Esperado 2 eventos, obtido %d", len(result))
		}

		if result[0].ID != "2" {
			t.Errorf("Esperado evento ID 2, obtido %s", result[0].ID)
		}

		if result[1].ID != "3" {
			t.Errorf("Esperado evento ID 3, obtido %s", result[1].ID)
		}
	})

	t.Run("TimelineServicePort handles empty events", func(t *testing.T) {
		mock := &mockTimelineServicePort{events: timeline.Timeline{}}
		result, err := mock.GetPatientTimeline(context.Background(), "patient-123", nil, 5, 0)

		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if result == nil {
			t.Error("Resultado deveria ser uma lista vazia, não nil")
		}

		if len(result) != 0 {
			t.Errorf("Esperado 0 eventos, obtido %d", len(result))
		}
	})
}
