package sqlite

import (
	"database/sql"
	"testing"
	"time"

	"arandu/internal/domain/observation"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()

	// Create in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	sqliteDB := &DB{db}

	// Run migrations to create all tables
	if err := sqliteDB.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return sqliteDB
}

func TestObservationRepository_SaveAndFindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewObservationRepository(db)

	// Create observation
	obs := &observation.Observation{
		SessionID: "session-123",
		Content:   "Paciente demonstrou ansiedade ao falar sobre trabalho",
	}

	// Save observation
	err := repo.Save(obs)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if obs.ID == "" {
		t.Error("Save() did not generate ID")
	}

	if obs.CreatedAt.IsZero() {
		t.Error("Save() did not set CreatedAt")
	}

	// Find by ID
	found, err := repo.FindByID(obs.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindByID() returned nil")
	}

	if found.ID != obs.ID {
		t.Errorf("FindByID() ID = %v, want %v", found.ID, obs.ID)
	}

	if found.SessionID != obs.SessionID {
		t.Errorf("FindByID() SessionID = %v, want %v", found.SessionID, obs.SessionID)
	}

	if found.Content != obs.Content {
		t.Errorf("FindByID() Content = %v, want %v", found.Content, obs.Content)
	}

	if !found.CreatedAt.Equal(obs.CreatedAt) {
		t.Errorf("FindByID() CreatedAt = %v, want %v", found.CreatedAt, obs.CreatedAt)
	}
}

func TestObservationRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewObservationRepository(db)

	// Create and save observation
	obs := &observation.Observation{
		SessionID: "session-123",
		Content:   "Conteúdo original",
	}

	err := repo.Save(obs)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Wait a bit to ensure updated_at is different
	time.Sleep(10 * time.Millisecond)

	// Update observation
	originalCreatedAt := obs.CreatedAt
	obs.Content = "Conteúdo atualizado"

	err = repo.Update(obs)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Find updated observation
	updated, err := repo.FindByID(obs.ID)
	if err != nil {
		t.Fatalf("FindByID() after update error = %v", err)
	}

	if updated == nil {
		t.Fatal("FindByID() after update returned nil")
	}

	// Check content was updated
	if updated.Content != "Conteúdo atualizado" {
		t.Errorf("Update() Content = %v, want %v", updated.Content, "Conteúdo atualizado")
	}

	// Check created_at remains unchanged
	if !updated.CreatedAt.Equal(originalCreatedAt) {
		t.Errorf("Update() changed CreatedAt = %v, original = %v", updated.CreatedAt, originalCreatedAt)
	}

	// Check updated_at was set
	if updated.UpdatedAt.IsZero() {
		t.Error("Update() did not set UpdatedAt")
	}

	// Verify updated_at is after created_at
	if !updated.UpdatedAt.After(updated.CreatedAt) {
		t.Errorf("Update() UpdatedAt = %v should be after CreatedAt = %v", updated.UpdatedAt, updated.CreatedAt)
	}
}

func TestObservationRepository_FindBySessionID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewObservationRepository(db)

	// Create observations for different sessions
	obs1 := &observation.Observation{
		SessionID: "session-123",
		Content:   "Observation 1 for session 123",
	}
	obs2 := &observation.Observation{
		SessionID: "session-123",
		Content:   "Observation 2 for session 123",
	}
	obs3 := &observation.Observation{
		SessionID: "session-456",
		Content:   "Observation for session 456",
	}

	// Save all observations
	for _, obs := range []*observation.Observation{obs1, obs2, obs3} {
		if err := repo.Save(obs); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	// Find observations for session-123
	observations, err := repo.FindBySessionID("session-123")
	if err != nil {
		t.Fatalf("FindBySessionID() error = %v", err)
	}

	if len(observations) != 2 {
		t.Errorf("FindBySessionID() got %d observations, want 2", len(observations))
	}

	// Check that we got the right observations
	foundObs1 := false
	foundObs2 := false
	for _, obs := range observations {
		if obs.ID == obs1.ID {
			foundObs1 = true
		}
		if obs.ID == obs2.ID {
			foundObs2 = true
		}
	}

	if !foundObs1 || !foundObs2 {
		t.Errorf("FindBySessionID() missing expected observations, foundObs1=%v, foundObs2=%v", foundObs1, foundObs2)
	}
}

func TestObservationRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewObservationRepository(db)

	// Create and save observation
	obs := &observation.Observation{
		SessionID: "session-123",
		Content:   "Observation to delete",
	}

	err := repo.Save(obs)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete observation
	err = repo.Delete(obs.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Try to find deleted observation
	deleted, err := repo.FindByID(obs.ID)
	if err != nil {
		t.Fatalf("FindByID() after delete error = %v", err)
	}

	if deleted != nil {
		t.Error("FindByID() after delete should return nil")
	}
}

func TestObservationRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewObservationRepository(db)

	// Create multiple observations
	observations := []*observation.Observation{
		{SessionID: "session-1", Content: "Observation 1"},
		{SessionID: "session-2", Content: "Observation 2"},
		{SessionID: "session-3", Content: "Observation 3"},
	}

	// Save all observations
	for _, obs := range observations {
		if err := repo.Save(obs); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	// Find all observations
	all, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(all) != 3 {
		t.Errorf("FindAll() got %d observations, want 3", len(all))
	}
}
