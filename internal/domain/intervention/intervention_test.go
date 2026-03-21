package intervention

import (
	"context"
	"testing"
	"time"
)

func TestIntervention_StructInitialization(t *testing.T) {
	now := time.Now()
	intervention := &Intervention{
		ID:        "interv-1",
		SessionID: "session-1",
		Content:   "Applied CBT technique",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if intervention.ID != "interv-1" {
		t.Errorf("expected ID 'interv-1', got '%s'", intervention.ID)
	}
	if intervention.SessionID != "session-1" {
		t.Errorf("expected SessionID 'session-1', got '%s'", intervention.SessionID)
	}
	if intervention.Content != "Applied CBT technique" {
		t.Errorf("unexpected content: %s", intervention.Content)
	}
}

func TestIntervention_ZeroTime(t *testing.T) {
	intervention := &Intervention{}
	if !intervention.CreatedAt.IsZero() {
		t.Error("CreatedAt should be zero initially")
	}
	if !intervention.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be zero initially")
	}
}

var _ Repository = (*MockInterventionRepository)(nil)

type MockInterventionRepository struct {
	SaveFunc            func(ctx context.Context, i *Intervention) error
	FindByIDFunc        func(ctx context.Context, id string) (*Intervention, error)
	FindBySessionIDFunc func(ctx context.Context, sessionID string) ([]*Intervention, error)
	FindAllFunc         func(ctx context.Context) ([]*Intervention, error)
	UpdateFunc          func(ctx context.Context, i *Intervention) error
	DeleteFunc          func(ctx context.Context, id string) error
}

func (m *MockInterventionRepository) Save(ctx context.Context, i *Intervention) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, i)
	}
	return nil
}

func (m *MockInterventionRepository) FindByID(ctx context.Context, id string) (*Intervention, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockInterventionRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*Intervention, error) {
	if m.FindBySessionIDFunc != nil {
		return m.FindBySessionIDFunc(ctx, sessionID)
	}
	return nil, nil
}

func (m *MockInterventionRepository) FindAll(ctx context.Context) ([]*Intervention, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockInterventionRepository) Update(ctx context.Context, i *Intervention) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, i)
	}
	return nil
}

func (m *MockInterventionRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestIntervention_RepositoryInterface(t *testing.T) {
	repo := &MockInterventionRepository{}
	var _ Repository = repo
}
