package observation

import (
	"context"
	"testing"
	"time"
)

func TestObservation_StructInitialization(t *testing.T) {
	now := time.Now()
	obs := &Observation{
		ID:        "obs-1",
		SessionID: "session-1",
		Content:   "Patient reported anxiety",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if obs.ID != "obs-1" {
		t.Errorf("expected ID 'obs-1', got '%s'", obs.ID)
	}
	if obs.SessionID != "session-1" {
		t.Errorf("expected SessionID 'session-1', got '%s'", obs.SessionID)
	}
	if obs.Content != "Patient reported anxiety" {
		t.Errorf("unexpected content: %s", obs.Content)
	}
}

func TestObservation_ZeroTime(t *testing.T) {
	obs := &Observation{}
	if !obs.CreatedAt.IsZero() {
		t.Error("CreatedAt should be zero initially")
	}
	if !obs.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be zero initially")
	}
}

func TestObservation_LongContent(t *testing.T) {
	longContent := ""
	for i := 0; i < 5000; i++ {
		longContent += "a"
	}
	obs := &Observation{
		ID:        "obs-long",
		SessionID: "session-1",
		Content:   longContent,
	}
	if len(obs.Content) != 5000 {
		t.Errorf("expected content length 5000, got %d", len(obs.Content))
	}
}

var _ Repository = (*MockObservationRepository)(nil)

type MockObservationRepository struct {
	SaveFunc            func(ctx context.Context, o *Observation) error
	FindByIDFunc        func(ctx context.Context, id string) (*Observation, error)
	FindBySessionIDFunc func(ctx context.Context, sessionID string) ([]*Observation, error)
	FindAllFunc         func(ctx context.Context) ([]*Observation, error)
	UpdateFunc          func(ctx context.Context, o *Observation) error
	DeleteFunc          func(ctx context.Context, id string) error
}

func (m *MockObservationRepository) Save(ctx context.Context, o *Observation) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, o)
	}
	return nil
}

func (m *MockObservationRepository) FindByID(ctx context.Context, id string) (*Observation, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockObservationRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*Observation, error) {
	if m.FindBySessionIDFunc != nil {
		return m.FindBySessionIDFunc(ctx, sessionID)
	}
	return nil, nil
}

func (m *MockObservationRepository) FindAll(ctx context.Context) ([]*Observation, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockObservationRepository) Update(ctx context.Context, o *Observation) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, o)
	}
	return nil
}

func (m *MockObservationRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockObservationRepository) GetTags(ctx context.Context) ([]Tag, error) {
	return nil, nil
}

func (m *MockObservationRepository) GetTagsByType(ctx context.Context, tagType TagType) ([]Tag, error) {
	return nil, nil
}

func (m *MockObservationRepository) AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error {
	return nil
}

func (m *MockObservationRepository) RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error {
	return nil
}

func (m *MockObservationRepository) GetObservationTags(ctx context.Context, observationID string) ([]ObservationTag, error) {
	return nil, nil
}

func (m *MockObservationRepository) GetTagsSummary(ctx context.Context) ([]TagSummary, error) {
	return nil, nil
}

func (m *MockObservationRepository) GetTagsSummaryByPatient(ctx context.Context, patientID string) ([]TagSummary, error) {
	return nil, nil
}

func (m *MockObservationRepository) FindByTag(ctx context.Context, tagID string) ([]*Observation, error) {
	return nil, nil
}

func TestObservation_RepositoryInterface(t *testing.T) {
	repo := &MockObservationRepository{}
	var _ Repository = repo
}
