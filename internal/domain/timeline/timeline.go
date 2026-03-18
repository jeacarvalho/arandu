package timeline

import (
	"context"
	"sort"
	"time"
)

type EventType string

const (
	EventTypeSession      EventType = "session"
	EventTypeObservation  EventType = "observation"
	EventTypeIntervention EventType = "intervention"
)

type TimelineEvent struct {
	ID        string
	Type      EventType
	Date      time.Time
	Content   string
	Metadata  map[string]string
	CreatedAt time.Time
}

type Timeline []*TimelineEvent

func (t Timeline) SortByDateDesc() {
	sort.Slice(t, func(i, j int) bool {
		return t[i].Date.After(t[j].Date)
	})
}

func (t Timeline) FilterByType(eventType EventType) Timeline {
	var filtered Timeline
	for _, event := range t {
		if event.Type == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func (t Timeline) GroupByDate() map[string]Timeline {
	groups := make(map[string]Timeline)

	for _, event := range t {
		dateKey := event.Date.Format("2006-01-02")
		groups[dateKey] = append(groups[dateKey], event)
	}

	return groups
}

type Repository interface {
	GetTimelineByPatientID(ctx context.Context, patientID string, filterType *EventType, limit, offset int) (Timeline, error)
}
