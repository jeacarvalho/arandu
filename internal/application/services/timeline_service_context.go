package services

import (
	"context"

	"arandu/internal/domain/timeline"
)

// TimelineRepositoryContextAware defines the interface for context-aware timeline operations
type TimelineRepositoryContextAware interface {
	GetTimelineByPatientID(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error)
	SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error)
}

// TimelineServiceContext is a context-aware timeline service
type TimelineServiceContext struct {
	repo TimelineRepositoryContextAware
}

// NewTimelineServiceContext creates a new context-aware timeline service
func NewTimelineServiceContext(repo TimelineRepositoryContextAware) *TimelineServiceContext {
	return &TimelineServiceContext{repo: repo}
}

func (s *TimelineServiceContext) GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	events, err := s.repo.GetTimelineByPatientID(ctx, patientID, filterType, limit, offset)
	if err != nil {
		return nil, err
	}

	events.SortByDateDesc()
	return events, nil
}

func (s *TimelineServiceContext) GetPatientTimelineWithFilter(ctx context.Context, patientID string, filterType timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	return s.GetPatientTimeline(ctx, patientID, &filterType, limit, offset)
}

func (s *TimelineServiceContext) GetPatientTimelineGroupedByDate(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (map[string]timeline.Timeline, error) {
	events, err := s.GetPatientTimeline(ctx, patientID, filterType, limit, offset)
	if err != nil {
		return nil, err
	}

	return events.GroupByDate(), nil
}

func (s *TimelineServiceContext) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	return s.repo.SearchInHistory(ctx, patientID, query)
}
