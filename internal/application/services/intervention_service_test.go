package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"arandu/internal/domain/intervention"
	"github.com/google/uuid"
)

var errInterventionRepoFailed = errors.New("intervention repository failed")

type mockInterventionRepoForService struct {
	interventions map[string]*intervention.Intervention
	saveErr       error
	findByIDErr   error
	updateErr     error
	deleteErr     error
}

func (m *mockInterventionRepoForService) Save(ctx context.Context, i *intervention.Intervention) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}
	m.interventions[i.ID] = i
	return nil
}

func (m *mockInterventionRepoForService) FindByID(ctx context.Context, id string) (*intervention.Intervention, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.interventions[id], nil
}

func (m *mockInterventionRepoForService) FindBySessionID(ctx context.Context, sessionID string) ([]*intervention.Intervention, error) {
	var result []*intervention.Intervention
	for _, i := range m.interventions {
		if i.SessionID == sessionID {
			result = append(result, i)
		}
	}
	return result, nil
}

func (m *mockInterventionRepoForService) FindAll(ctx context.Context) ([]*intervention.Intervention, error) {
	var result []*intervention.Intervention
	for _, i := range m.interventions {
		result = append(result, i)
	}
	return result, nil
}

func (m *mockInterventionRepoForService) Update(ctx context.Context, i *intervention.Intervention) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.interventions[i.ID] = i
	return nil
}

func (m *mockInterventionRepoForService) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.interventions, id)
	return nil
}

func TestInterventionService_CreateIntervention_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockInterventionRepoForService{
		interventions: make(map[string]*intervention.Intervention),
		saveErr:       errInterventionRepoFailed,
	}
	service := NewInterventionService(repo)

	interv, err := service.CreateIntervention(ctx, "session-123", "Applied CBT technique")

	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
	if interv != nil {
		t.Error("Expected nil intervention on error")
	}
}

func TestInterventionService_GetIntervention_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockInterventionRepoForService{
		interventions: make(map[string]*intervention.Intervention),
	}
	service := NewInterventionService(repo)

	interv, _ := service.GetIntervention(ctx, "non-existent")

	// Service returns nil, nil for not found (common pattern)
	// This test verifies that behavior
	if interv != nil {
		t.Error("Expected nil intervention for non-existent")
	}
	// No error expected for not found - this is the current behavior
}

func TestInterventionService_GetIntervention_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockInterventionRepoForService{
		interventions: make(map[string]*intervention.Intervention),
		findByIDErr:   errInterventionRepoFailed,
	}
	service := NewInterventionService(repo)

	interv, err := service.GetIntervention(ctx, "any-id")

	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
	if interv != nil {
		t.Error("Expected nil intervention on error")
	}
}

func TestInterventionService_UpdateIntervention_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockInterventionRepoForService{
		interventions: map[string]*intervention.Intervention{
			"interv-1": {ID: "interv-1", SessionID: "session-123", Content: "Old content"},
		},
		updateErr: errInterventionRepoFailed,
	}
	service := NewInterventionService(repo)

	err := service.UpdateIntervention(ctx, "interv-1", "New content")

	if err == nil {
		t.Error("Expected error from update, got nil")
	}
}

func TestInterventionService_DeleteIntervention_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockInterventionRepoForService{
		interventions: map[string]*intervention.Intervention{
			"interv-1": {ID: "interv-1", SessionID: "session-123"},
		},
		deleteErr: errInterventionRepoFailed,
	}
	service := NewInterventionService(repo)

	err := service.DeleteIntervention(ctx, "interv-1")

	if err == nil {
		t.Error("Expected error from delete, got nil")
	}
}