package services

import (
	"context"

	"arandu/internal/domain/patient"
	"arandu/internal/domain/timeline"
)

// TimelineRepositoryContextAware defines the interface for context-aware timeline operations
type TimelineRepositoryContextAware interface {
	GetTimelineByPatientID(ctx context.Context, patientID string, filterType *timeline.EventType, limit, offset int) (timeline.Timeline, error)
	SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error)
	SearchGlobal(ctx context.Context, query string, limit int) ([]*timeline.SearchResult, error)
}

// TimelineServiceContext is a context-aware timeline service
type TimelineServiceContext struct {
	repo          TimelineRepositoryContextAware
	patientService PatientServiceGetter
}

// PatientServiceGetter defines the interface for getting patient data
type PatientServiceGetter interface {
	GetPatient(ctx context.Context, id string) (*patient.Patient, error)
	GetPatientByID(ctx context.Context, id string) (*patient.Patient, error)
}

// SearchGlobalResult contains search results enriched with patient name
type SearchGlobalResult struct {
	timeline.SearchResult
	PatientName string
}

// NewTimelineServiceContext creates a new context-aware timeline service
func NewTimelineServiceContext(repo TimelineRepositoryContextAware, patientService PatientServiceGetter) *TimelineServiceContext {
	return &TimelineServiceContext{repo: repo, patientService: patientService}
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

func (s *TimelineServiceContext) SearchGlobal(ctx context.Context, query string) ([]*SearchGlobalResult, error) {
	results, err := s.repo.SearchGlobal(ctx, query, 50)
	if err != nil {
		return nil, err
	}

	patientNames := make(map[string]string)
	var uniquePatientIDs []string
	for _, r := range results {
		if _, ok := patientNames[r.PatientID]; !ok {
			patientNames[r.PatientID] = ""
			uniquePatientIDs = append(uniquePatientIDs, r.PatientID)
		}
	}

	for _, patientID := range uniquePatientIDs {
		p, err := s.patientService.GetPatientByID(ctx, patientID)
		if err == nil && p != nil {
			patientNames[patientID] = p.Name
		}
	}

	var enrichedResults []*SearchGlobalResult
	for _, r := range results {
		enrichedResults = append(enrichedResults, &SearchGlobalResult{
			SearchResult: *r,
			PatientName: patientNames[r.PatientID],
		})
	}

	return enrichedResults, nil
}
