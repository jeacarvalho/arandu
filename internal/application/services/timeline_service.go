package services

import (
	"context"

	"arandu/internal/domain/timeline"
	"arandu/internal/infrastructure/repository/sqlite"
)

type TimelineService struct {
	repo *sqlite.TimelineRepository
}

func NewTimelineService(repo *sqlite.TimelineRepository) *TimelineService {
	return &TimelineService{repo: repo}
}

func (s *TimelineService) GetPatientTimeline(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	events, err := s.repo.GetTimelineByPatientID(ctx, patientID, filterType, limit, offset)
	if err != nil {
		return nil, err
	}

	events.SortByDateDesc()
	return events, nil
}

func (s *TimelineService) GetPatientTimelineWithFilter(ctx context.Context, patientID string, filterType timeline.EventType, limit, offset int) (timeline.Timeline, error) {
	return s.GetPatientTimeline(ctx, patientID, &filterType, limit, offset)
}

func (s *TimelineService) GetPatientTimelineGroupedByDate(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (map[string]timeline.Timeline, error) {
	events, err := s.GetPatientTimeline(ctx, patientID, filterType, limit, offset)
	if err != nil {
		return nil, err
	}

	return events.GroupByDate(), nil
}

func (s *TimelineService) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error) {
	return s.repo.SearchInHistory(ctx, patientID, query)
}
