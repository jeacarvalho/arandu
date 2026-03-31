package timeline

import (
	"bytes"
	"testing"
	"time"

	domainTimeline "arandu/internal/domain/timeline"
)

func TestTimelineFiltersRender(t *testing.T) {
	tests := []struct {
		name         string
		patientID    string
		activeFilter *domainTimeline.EventType
		expectInHTML []string
		notExpect    []string
	}{
		{
			name:         "renders all filter buttons",
			patientID:    "patient-123",
			activeFilter: nil,
			expectInHTML: []string{
				`hx-get="/patients/patient-123/history?filter=all"`,
				`hx-get="/patients/patient-123/history?filter=session"`,
				`hx-get="/patients/patient-123/history?filter=observation"`,
				`hx-get="/patients/patient-123/history?filter=intervention"`,
				"Todos",
				"Sessões",
				"Observações",
				"Intervenções",
				`class="timeline-filter-btn active"`,
			},
		},
		{
			name:         "highlights active session filter",
			patientID:    "patient-456",
			activeFilter: func() *domainTimeline.EventType { t := domainTimeline.EventTypeSession; return &t }(),
			expectInHTML: []string{
				"timeline-filter-btn active",
			},
		},
		{
			name:         "highlights active observation filter",
			patientID:    "patient-789",
			activeFilter: func() *domainTimeline.EventType { t := domainTimeline.EventTypeObservation; return &t }(),
			expectInHTML: []string{
				"timeline-filter-btn active",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := TimelineFilters(tt.patientID, tt.activeFilter).Render(t.Context(), &buf)
			if err != nil {
				t.Fatalf("Failed to render TimelineFilters: %v", err)
			}

			html := buf.String()

			for _, expected := range tt.expectInHTML {
				if !bytes.Contains(buf.Bytes(), []byte(expected)) {
					t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
				}
			}

			for _, notExpected := range tt.notExpect {
				if bytes.Contains(buf.Bytes(), []byte(notExpected)) {
					t.Errorf("Expected HTML NOT to contain %q, but it did.\nHTML: %s", notExpected, html)
				}
			}
		})
	}
}

func TestTimelineEventCardRender(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		event           *domainTimeline.TimelineEvent
		showSessionLink bool
		expectInHTML    []string
		notExpect       []string
	}{
		{
			name: "renders session event with all fields",
			event: &domainTimeline.TimelineEvent{
				ID:        "event-1",
				Type:      domainTimeline.EventTypeSession,
				Date:      now,
				Content:   "Conteúdo da sessão de teste",
				CreatedAt: now.Add(-time.Hour),
				Metadata:  map[string]string{"session_id": "session-123"},
			},
			showSessionLink: true,
			expectInHTML: []string{
				"timeline-event",
				"timeline-dot-session",
				"timeline-event-card-session",
				"Sessão",
				"Conteúdo da sessão de teste",
				"/session/session-123",
				"Ver sessão",
			},
		},
		{
			name: "renders observation event without session link",
			event: &domainTimeline.TimelineEvent{
				ID:        "event-2",
				Type:      domainTimeline.EventTypeObservation,
				Date:      now,
				Content:   "Observação clínica",
				CreatedAt: now.Add(-2 * time.Hour),
				Metadata:  map[string]string{},
			},
			showSessionLink: false,
			expectInHTML: []string{
				"timeline-dot-observation",
				"timeline-event-card-observation",
				"Observação",
				"Observação clínica",
			},
			notExpect: []string{"Ver sessão"},
		},
		{
			name: "renders intervention event",
			event: &domainTimeline.TimelineEvent{
				ID:        "event-3",
				Type:      domainTimeline.EventTypeIntervention,
				Date:      now,
				Content:   "Intervenção realizada",
				CreatedAt: now.Add(-3 * time.Hour),
				Metadata:  map[string]string{},
			},
			showSessionLink: true,
			expectInHTML: []string{
				"timeline-dot-intervention",
				"timeline-event-card-intervention",
				"Intervenção",
				"Intervenção realizada",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := TimelineEventCard(tt.event, tt.showSessionLink).Render(t.Context(), &buf)
			if err != nil {
				t.Fatalf("Failed to render TimelineEventCard: %v", err)
			}

			html := buf.String()

			for _, expected := range tt.expectInHTML {
				if !bytes.Contains(buf.Bytes(), []byte(expected)) {
					t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
				}
			}

			for _, notExpected := range tt.notExpect {
				if bytes.Contains(buf.Bytes(), []byte(notExpected)) {
					t.Errorf("Expected HTML NOT to contain %q, but it did.\nHTML: %s", notExpected, html)
				}
			}
		})
	}
}

func TestTimelineContentRender(t *testing.T) {
	now := time.Now()

	t.Run("renders empty state when no events", func(t *testing.T) {
		var buf bytes.Buffer
		events := domainTimeline.Timeline{}
		err := TimelineContent("patient-123", events).Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineContent: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"timeline-empty",
			"Nenhum evento clínico registrado",
			"/patients/patient-123/sessions/new",
			"Criar primeira sessão",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})

	t.Run("renders events when present", func(t *testing.T) {
		var buf bytes.Buffer
		events := domainTimeline.Timeline{
			{
				ID:        "event-1",
				Type:      domainTimeline.EventTypeSession,
				Date:      now,
				Content:   "Primeira sessão",
				CreatedAt: now,
				Metadata:  map[string]string{"session_id": "sess-1"},
			},
			{
				ID:        "event-2",
				Type:      domainTimeline.EventTypeObservation,
				Date:      now.Add(-time.Hour),
				Content:   "Observação importante",
				CreatedAt: now.Add(-time.Hour),
				Metadata:  map[string]string{},
			},
		}
		err := TimelineContent("patient-456", events).Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineContent: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"timeline-container",
			"timeline-line",
			"timeline-content",
			"Primeira sessão",
			"Observação importante",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})
}

func TestTimelineEmptyRender(t *testing.T) {
	t.Run("renders empty state with CTA", func(t *testing.T) {
		var buf bytes.Buffer
		err := TimelineEmpty("patient-123", true).Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineEmpty: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"timeline-empty",
			"Nenhum evento clínico registrado",
			"Criar primeira sessão",
			"/patients/patient-123/sessions/new",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})

	t.Run("renders empty state without CTA", func(t *testing.T) {
		var buf bytes.Buffer
		err := TimelineEmpty("patient-456", false).Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineEmpty: %v", err)
		}

		html := buf.String()
		if bytes.Contains(buf.Bytes(), []byte("Criar primeira sessão")) {
			t.Errorf("Expected HTML NOT to contain CTA, but it did.\nHTML: %s", html)
		}
	})
}

func TestTimelineSearchRender(t *testing.T) {
	t.Run("renders search input with HTMX attributes", func(t *testing.T) {
		var buf bytes.Buffer
		err := TimelineSearch("patient-123").Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineSearch: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"timeline-search",
			"timeline-search-input",
			`hx-get="/patients/patient-123/history/search"`,
			`hx-target="#timeline-content"`,
			`hx-trigger="keyup changed delay:500ms, search"`,
			`hx-swap="innerHTML"`,
			"Buscar no histórico clínico",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})
}

func TestTimelineLoadingRender(t *testing.T) {
	t.Run("renders loading indicator", func(t *testing.T) {
		var buf bytes.Buffer
		err := TimelineLoading().Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineLoading: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"htmx-indicator",
			"timeline-loading",
			"A carregar mais registros",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})
}

func TestTimelineEndRender(t *testing.T) {
	t.Run("renders end marker", func(t *testing.T) {
		var buf bytes.Buffer
		err := TimelineEnd().Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render TimelineEnd: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"timeline-end",
			"Fim dos registros históricos",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})
}

func TestFiltersAndContentRender(t *testing.T) {
	now := time.Now()
	events := domainTimeline.Timeline{
		{
			ID:        "event-1",
			Type:      domainTimeline.EventTypeSession,
			Date:      now,
			Content:   "Sessão de teste",
			CreatedAt: now,
			Metadata:  map[string]string{},
		},
	}

	t.Run("renders filters and content together", func(t *testing.T) {
		var buf bytes.Buffer
		err := FiltersAndContent("patient-123", events, nil).Render(t.Context(), &buf)
		if err != nil {
			t.Fatalf("Failed to render FiltersAndContent: %v", err)
		}

		html := buf.String()
		expectedStrings := []string{
			"timeline-filters",
			"filter=all",
			"timeline-container",
			"timeline-content",
			"Sessão de teste",
		}

		for _, expected := range expectedStrings {
			if !bytes.Contains(buf.Bytes(), []byte(expected)) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nHTML: %s", expected, html)
			}
		}
	})
}
