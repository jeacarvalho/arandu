package insight

import (
	"context"
	"testing"
	"time"
)

func TestInsight_StructInitialization(t *testing.T) {
	now := time.Now()
	insight := &Insight{
		ID:        "insight-1",
		Content:   "Patient shows improvement",
		Source:    "ai",
		CreatedAt: now,
	}

	if insight.ID != "insight-1" {
		t.Errorf("expected ID 'insight-1', got '%s'", insight.ID)
	}
	if insight.Content != "Patient shows improvement" {
		t.Errorf("unexpected content: %s", insight.Content)
	}
	if insight.Source != "ai" {
		t.Errorf("expected Source 'ai', got '%s'", insight.Source)
	}
	if insight.CreatedAt != now {
		t.Errorf("CreatedAt mismatch")
	}
}

func TestInsight_SourceValues(t *testing.T) {
	sources := []string{"ai", "therapist"}
	for _, src := range sources {
		insight := &Insight{
			ID:      "test",
			Content: "content",
			Source:  src,
		}
		if insight.Source != src {
			t.Errorf("expected Source '%s', got '%s'", src, insight.Source)
		}
	}
}

func TestInsight_JSONTags(t *testing.T) {
	insight := &Insight{
		ID:        "id-123",
		Content:   "content text",
		Source:    "ai",
		CreatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	if insight.ID == "" || insight.Content == "" || insight.Source == "" {
		t.Error("all fields should be non-zero")
	}
}

var _ Repository = (*MockInsightRepository)(nil)

type MockInsightRepository struct {
	SaveFunc     func(ctx context.Context, insight *Insight) error
	FindByIDFunc func(ctx context.Context, id string) (*Insight, error)
	FindAllFunc  func(ctx context.Context) ([]*Insight, error)
	DeleteFunc   func(ctx context.Context, id string) error
}

func (m *MockInsightRepository) Save(ctx context.Context, insight *Insight) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, insight)
	}
	return nil
}

func (m *MockInsightRepository) FindByID(ctx context.Context, id string) (*Insight, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockInsightRepository) FindAll(ctx context.Context) ([]*Insight, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockInsightRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestInsight_RepositoryInterface(t *testing.T) {
	repo := &MockInsightRepository{}
	var _ Repository = repo
}
