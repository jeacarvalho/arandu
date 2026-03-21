package timeline

import (
	"testing"
	"time"
)

func TestTimeline_SortByDateDesc(t *testing.T) {
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 3, 20, 14, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 2, 10, 9, 0, 0, 0, time.UTC)

	timeline := Timeline{
		{ID: "1", Date: t1},
		{ID: "2", Date: t2},
		{ID: "3", Date: t3},
	}

	timeline.SortByDateDesc()

	if timeline[0].ID != "2" {
		t.Errorf("expected first item ID=2, got %s", timeline[0].ID)
	}
	if timeline[1].ID != "3" {
		t.Errorf("expected second item ID=3, got %s", timeline[1].ID)
	}
	if timeline[2].ID != "1" {
		t.Errorf("expected third item ID=1, got %s", timeline[2].ID)
	}
}

func TestTimeline_SortByDateDesc_EmptyTimeline(t *testing.T) {
	timeline := Timeline{}
	timeline.SortByDateDesc()
	if len(timeline) != 0 {
		t.Errorf("expected empty timeline, got %d items", len(timeline))
	}
}

func TestTimeline_SortByDateDesc_SingleItem(t *testing.T) {
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	timeline := Timeline{{ID: "1", Date: t1}}
	timeline.SortByDateDesc()
	if len(timeline) != 1 || timeline[0].ID != "1" {
		t.Errorf("expected single item unchanged")
	}
}

func TestTimeline_SortByDateDesc_AlreadySorted(t *testing.T) {
	t1 := time.Date(2024, 3, 20, 14, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	timeline := Timeline{
		{ID: "1", Date: t1},
		{ID: "2", Date: t2},
	}
	timeline.SortByDateDesc()
	if timeline[0].ID != "1" || timeline[1].ID != "2" {
		t.Errorf("expected order preserved")
	}
}

func TestTimeline_FilterByType(t *testing.T) {
	timeline := Timeline{
		{ID: "1", Type: EventTypeSession},
		{ID: "2", Type: EventTypeObservation},
		{ID: "3", Type: EventTypeSession},
		{ID: "4", Type: EventTypeIntervention},
		{ID: "5", Type: EventTypeObservation},
	}

	t.Run("filter session", func(t *testing.T) {
		filtered := timeline.FilterByType(EventTypeSession)
		if len(filtered) != 2 {
			t.Errorf("expected 2 sessions, got %d", len(filtered))
		}
		for _, e := range filtered {
			if e.Type != EventTypeSession {
				t.Errorf("expected EventTypeSession, got %s", e.Type)
			}
		}
	})

	t.Run("filter observation", func(t *testing.T) {
		filtered := timeline.FilterByType(EventTypeObservation)
		if len(filtered) != 2 {
			t.Errorf("expected 2 observations, got %d", len(filtered))
		}
	})

	t.Run("filter intervention", func(t *testing.T) {
		filtered := timeline.FilterByType(EventTypeIntervention)
		if len(filtered) != 1 {
			t.Errorf("expected 1 intervention, got %d", len(filtered))
		}
	})

	t.Run("filter non-existent", func(t *testing.T) {
		filtered := timeline.FilterByType(EventType("unknown"))
		if len(filtered) != 0 {
			t.Errorf("expected 0 items, got %d", len(filtered))
		}
	})
}

func TestTimeline_FilterByType_EmptyTimeline(t *testing.T) {
	timeline := Timeline{}
	filtered := timeline.FilterByType(EventTypeSession)
	if len(filtered) != 0 {
		t.Errorf("expected 0 items, got %d", len(filtered))
	}
}

func TestTimeline_GroupByDate(t *testing.T) {
	d1 := time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 2, 15, 14, 0, 0, 0, time.UTC)
	d3 := time.Date(2024, 2, 17, 9, 0, 0, 0, time.UTC)
	d4 := time.Date(2024, 2, 17, 11, 0, 0, 0, time.UTC)
	d5 := time.Date(2024, 2, 20, 8, 0, 0, 0, time.UTC)

	timeline := Timeline{
		{ID: "1", Date: d1},
		{ID: "2", Date: d2},
		{ID: "3", Date: d3},
		{ID: "4", Date: d4},
		{ID: "5", Date: d5},
	}

	groups := timeline.GroupByDate()

	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}

	if len(groups["2024-02-15"]) != 2 {
		t.Errorf("expected 2 items in 2024-02-15, got %d", len(groups["2024-02-15"]))
	}

	if len(groups["2024-02-17"]) != 2 {
		t.Errorf("expected 2 items in 2024-02-17, got %d", len(groups["2024-02-17"]))
	}

	if len(groups["2024-02-20"]) != 1 {
		t.Errorf("expected 1 item in 2024-02-20, got %d", len(groups["2024-02-20"]))
	}
}

func TestTimeline_GroupByDate_EmptyTimeline(t *testing.T) {
	timeline := Timeline{}
	groups := timeline.GroupByDate()
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestTimeline_GroupByDate_SameDateDifferentTimes(t *testing.T) {
	d1 := time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC)
	d2 := time.Date(2024, 2, 15, 23, 59, 59, 0, time.UTC)
	timeline := Timeline{
		{ID: "1", Date: d1},
		{ID: "2", Date: d2},
	}
	groups := timeline.GroupByDate()
	if len(groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(groups))
	}
}

func TestEventType_Constants(t *testing.T) {
	if EventTypeSession != "session" {
		t.Errorf("EventTypeSession should be 'session'")
	}
	if EventTypeObservation != "observation" {
		t.Errorf("EventTypeObservation should be 'observation'")
	}
	if EventTypeIntervention != "intervention" {
		t.Errorf("EventTypeIntervention should be 'intervention'")
	}
}
