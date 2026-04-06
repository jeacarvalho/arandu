package services

import (
	"context"
	"testing"
	"time"

	"arandu/internal/domain/observation"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type mockObservationRepository struct {
	observations map[string]*observation.Observation
}

func (m *mockObservationRepository) Save(ctx context.Context, o *observation.Observation) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	m.observations[o.ID] = o
	return nil
}

func (m *mockObservationRepository) FindByID(ctx context.Context, id string) (*observation.Observation, error) {
	return m.observations[id], nil
}

func (m *mockObservationRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*observation.Observation, error) {
	var result []*observation.Observation
	for _, obs := range m.observations {
		if obs.SessionID == sessionID {
			result = append(result, obs)
		}
	}
	return result, nil
}

func (m *mockObservationRepository) FindAll(ctx context.Context) ([]*observation.Observation, error) {
	var result []*observation.Observation
	for _, obs := range m.observations {
		result = append(result, obs)
	}
	return result, nil
}

func (m *mockObservationRepository) Update(ctx context.Context, o *observation.Observation) error {
	m.observations[o.ID] = o
	return nil
}

func (m *mockObservationRepository) Delete(ctx context.Context, id string) error {
	delete(m.observations, id)
	return nil
}

func (m *mockObservationRepository) GetTags(ctx context.Context) ([]observation.Tag, error) {
	return nil, nil
}

func (m *mockObservationRepository) GetTagsByType(ctx context.Context, tagType observation.TagType) ([]observation.Tag, error) {
	return nil, nil
}

func (m *mockObservationRepository) AddTagToObservation(ctx context.Context, observationID, tagID string, intensity int) error {
	return nil
}

func (m *mockObservationRepository) RemoveTagFromObservation(ctx context.Context, observationID, tagID string) error {
	return nil
}

func (m *mockObservationRepository) GetObservationTags(ctx context.Context, observationID string) ([]observation.ObservationTag, error) {
	return nil, nil
}

func (m *mockObservationRepository) GetTagsSummary(ctx context.Context) ([]observation.TagSummary, error) {
	return nil, nil
}

func (m *mockObservationRepository) GetTagsSummaryByPatient(ctx context.Context, patientID string) ([]observation.TagSummary, error) {
	return nil, nil
}

func (m *mockObservationRepository) FindByTag(ctx context.Context, tagID string) ([]*observation.Observation, error) {
	return nil, nil
}

func TestObservationService_CreateObservation(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name:      "valid observation",
			sessionID: "session-123",
			content:   "Paciente demonstrou resistência ao falar sobre a infância",
			wantErr:   false,
		},
		{
			name:        "empty content",
			sessionID:   "session-123",
			content:     "",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "content too long",
			sessionID:   "session-123",
			content:     string(make([]byte, 5001)),
			wantErr:     true,
			errContains: "cannot exceed 5000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &mockObservationRepository{
				observations: make(map[string]*observation.Observation),
			}
			service := NewObservationService(repo)

			obs, err := service.CreateObservation(ctx, tt.sessionID, tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateObservation() expected error, got nil")
				}
				if tt.errContains != "" && err != nil && err.Error() != tt.errContains && !contains(err.Error(), tt.errContains) {
					t.Errorf("CreateObservation() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateObservation() unexpected error: %v", err)
				return
			}

			if obs == nil {
				t.Errorf("CreateObservation() returned nil observation")
				return
			}

			if obs.SessionID != tt.sessionID {
				t.Errorf("CreateObservation() SessionID = %v, want %v", obs.SessionID, tt.sessionID)
			}

			if obs.Content != tt.content {
				t.Errorf("CreateObservation() Content = %v, want %v", obs.Content, tt.content)
			}

			if obs.ID == "" {
				t.Errorf("CreateObservation() ID is empty")
			}

			if obs.CreatedAt.IsZero() {
				t.Errorf("CreateObservation() CreatedAt is zero")
			}

			// Verify it was saved
			saved, err := repo.FindByID(ctx, obs.ID)
			if err != nil {
				t.Errorf("FindByID() error: %v", err)
			}
			if saved == nil {
				t.Errorf("Observation was not saved to repository")
			}
			if diff := cmp.Diff(obs, saved); diff != "" {
				t.Errorf("Saved observation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestObservationService_ListObservationsBySession(t *testing.T) {
	ctx := context.Background()
	repo := &mockObservationRepository{
		observations: map[string]*observation.Observation{
			"obs-1": {
				ID:        "obs-1",
				SessionID: "session-123",
				Content:   "Observation 1",
			},
			"obs-2": {
				ID:        "obs-2",
				SessionID: "session-123",
				Content:   "Observation 2",
			},
			"obs-3": {
				ID:        "obs-3",
				SessionID: "session-456",
				Content:   "Observation 3",
			},
		},
	}
	service := NewObservationService(repo)

	observations, err := service.ListObservationsBySession(ctx, "session-123")
	if err != nil {
		t.Errorf("ListObservationsBySession() error: %v", err)
		return
	}

	if len(observations) != 2 {
		t.Errorf("ListObservationsBySession() got %d observations, want 2", len(observations))
	}

	// Check that we got the right observations
	foundObs1 := false
	foundObs2 := false
	for _, obs := range observations {
		if obs.ID == "obs-1" {
			foundObs1 = true
		}
		if obs.ID == "obs-2" {
			foundObs2 = true
		}
	}

	if !foundObs1 || !foundObs2 {
		t.Errorf("ListObservationsBySession() missing expected observations, foundObs1=%v, foundObs2=%v", foundObs1, foundObs2)
	}
}

func TestObservationService_UpdateObservation(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid update",
			id:      "obs-1",
			content: "Conteúdo atualizado da observação",
			wantErr: false,
		},
		{
			name:        "empty content",
			id:          "obs-1",
			content:     "",
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "content too long",
			id:          "obs-1",
			content:     string(make([]byte, 5001)),
			wantErr:     true,
			errContains: "cannot exceed 5000 characters",
		},
		{
			name:        "observation not found",
			id:          "non-existent",
			content:     "Conteúdo válido",
			wantErr:     true,
			errContains: "observation not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &mockObservationRepository{
				observations: map[string]*observation.Observation{
					"obs-1": {
						ID:        "obs-1",
						SessionID: "session-123",
						Content:   "Conteúdo original",
						CreatedAt: time.Now(),
					},
				},
			}
			service := NewObservationService(repo)

			err := service.UpdateObservation(ctx, tt.id, tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateObservation() expected error, got nil")
				}
				if tt.errContains != "" && err != nil && err.Error() != tt.errContains && !contains(err.Error(), tt.errContains) {
					t.Errorf("UpdateObservation() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateObservation() unexpected error: %v", err)
				return
			}

			// Verify the observation was updated
			updated, err := repo.FindByID(ctx, tt.id)
			if err != nil {
				t.Errorf("FindByID() error: %v", err)
				return
			}

			if updated == nil {
				t.Errorf("Observation was not found after update")
				return
			}

			if updated.Content != tt.content {
				t.Errorf("UpdateObservation() Content = %v, want %v", updated.Content, tt.content)
			}

			// Verify other fields remain unchanged
			if updated.SessionID != "session-123" {
				t.Errorf("UpdateObservation() changed SessionID = %v, want %v", updated.SessionID, "session-123")
			}

			if updated.ID != "obs-1" {
				t.Errorf("UpdateObservation() changed ID = %v, want %v", updated.ID, "obs-1")
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
