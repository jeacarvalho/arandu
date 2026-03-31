package timeline

import (
	"testing"

	domainTimeline "arandu/internal/domain/timeline"
)

func TestGetEventTypeLabel(t *testing.T) {
	tests := []struct {
		name      string
		eventType domainTimeline.EventType
		expected  string
	}{
		{
			name:      "session type returns Sessão",
			eventType: domainTimeline.EventTypeSession,
			expected:  "Sessão",
		},
		{
			name:      "observation type returns Observação",
			eventType: domainTimeline.EventTypeObservation,
			expected:  "Observação",
		},
		{
			name:      "intervention type returns Intervenção",
			eventType: domainTimeline.EventTypeIntervention,
			expected:  "Intervenção",
		},
		{
			name:      "unknown type returns Evento",
			eventType: domainTimeline.EventType("unknown"),
			expected:  "Evento",
		},
		{
			name:      "empty type returns Evento",
			eventType: domainTimeline.EventType(""),
			expected:  "Evento",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEventTypeLabel(tt.eventType)
			if result != tt.expected {
				t.Errorf("GetEventTypeLabel(%q) = %q, expected %q", tt.eventType, result, tt.expected)
			}
		})
	}
}

func TestGetEventIcon(t *testing.T) {
	tests := []struct {
		name      string
		eventType domainTimeline.EventType
		expected  string
	}{
		{
			name:      "session type returns calendar-check icon",
			eventType: domainTimeline.EventTypeSession,
			expected:  "fas fa-calendar-check",
		},
		{
			name:      "observation type returns sticky-note icon",
			eventType: domainTimeline.EventTypeObservation,
			expected:  "fas fa-sticky-note",
		},
		{
			name:      "intervention type returns hand-holding-heart icon",
			eventType: domainTimeline.EventTypeIntervention,
			expected:  "fas fa-hand-holding-heart",
		},
		{
			name:      "unknown type returns calendar-day icon",
			eventType: domainTimeline.EventType("unknown"),
			expected:  "fas fa-calendar-day",
		},
		{
			name:      "empty type returns calendar-day icon",
			eventType: domainTimeline.EventType(""),
			expected:  "fas fa-calendar-day",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEventIcon(tt.eventType)
			if result != tt.expected {
				t.Errorf("GetEventIcon(%q) = %q, expected %q", tt.eventType, result, tt.expected)
			}
		})
	}
}

func TestGetEventTypeClass(t *testing.T) {
	tests := []struct {
		name      string
		eventType domainTimeline.EventType
		expected  string
	}{
		{
			name:      "session type returns session class",
			eventType: domainTimeline.EventTypeSession,
			expected:  "session",
		},
		{
			name:      "observation type returns observation class",
			eventType: domainTimeline.EventTypeObservation,
			expected:  "observation",
		},
		{
			name:      "intervention type returns intervention class",
			eventType: domainTimeline.EventTypeIntervention,
			expected:  "intervention",
		},
		{
			name:      "unknown type returns empty string",
			eventType: domainTimeline.EventType("unknown"),
			expected:  "",
		},
		{
			name:      "empty type returns empty string",
			eventType: domainTimeline.EventType(""),
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEventTypeClass(tt.eventType)
			if result != tt.expected {
				t.Errorf("GetEventTypeClass(%q) = %q, expected %q", tt.eventType, result, tt.expected)
			}
		})
	}
}

func TestGetActiveClass(t *testing.T) {
	sessionType := domainTimeline.EventTypeSession
	observationType := domainTimeline.EventTypeObservation

	tests := []struct {
		name         string
		activeFilter *domainTimeline.EventType
		buttonFilter *domainTimeline.EventType
		expected     string
	}{
		{
			name:         "both nil returns active",
			activeFilter: nil,
			buttonFilter: nil,
			expected:     "active",
		},
		{
			name:         "active filter nil returns empty",
			activeFilter: nil,
			buttonFilter: &sessionType,
			expected:     "",
		},
		{
			name:         "button filter nil returns empty",
			activeFilter: &sessionType,
			buttonFilter: nil,
			expected:     "",
		},
		{
			name:         "matching types returns active",
			activeFilter: &sessionType,
			buttonFilter: &sessionType,
			expected:     "active",
		},
		{
			name:         "different types returns empty",
			activeFilter: &sessionType,
			buttonFilter: &observationType,
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getActiveClass(tt.activeFilter, tt.buttonFilter)
			if result != tt.expected {
				t.Errorf("getActiveClass(%v, %v) = %q, expected %q", tt.activeFilter, tt.buttonFilter, result, tt.expected)
			}
		})
	}
}

func TestGetActiveClassForType(t *testing.T) {
	sessionType := domainTimeline.EventTypeSession
	observationType := domainTimeline.EventTypeObservation

	tests := []struct {
		name         string
		activeFilter *domainTimeline.EventType
		buttonFilter domainTimeline.EventType
		expected     string
	}{
		{
			name:         "nil active filter returns empty",
			activeFilter: nil,
			buttonFilter: sessionType,
			expected:     "",
		},
		{
			name:         "matching types returns active",
			activeFilter: &sessionType,
			buttonFilter: sessionType,
			expected:     "active",
		},
		{
			name:         "different types returns empty",
			activeFilter: &sessionType,
			buttonFilter: observationType,
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getActiveClassForType(tt.activeFilter, tt.buttonFilter)
			if result != tt.expected {
				t.Errorf("getActiveClassForType(%v, %q) = %q, expected %q", tt.activeFilter, tt.buttonFilter, result, tt.expected)
			}
		})
	}
}
