package e2e

import (
	"context"
	"strings"
	"testing"

	"arandu/internal/application/services"
	"arandu/internal/domain/observation"
	"arandu/internal/infrastructure/repository/sqlite"
)

func TestObservationEditFlow(t *testing.T) {
	ctx := context.Background()

	// Setup in-memory database
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize database with migrations
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	observationRepo := sqlite.NewObservationRepository(db)

	// Create services
	observationService := services.NewObservationService(observationRepo)

	// Create a test observation
	obs, err := observationService.CreateObservation(ctx, "session-test-123", "Observação clínica inicial para teste")
	if err != nil {
		t.Fatalf("Failed to create test observation: %v", err)
	}

	// Test 1: Get observation item (view mode)
	t.Run("GetObservationItem", func(t *testing.T) {
		// Test the service layer directly
		retrieved, err := observationService.GetObservation(ctx, obs.ID)
		if err != nil {
			t.Errorf("GetObservation failed: %v", err)
		}

		if retrieved == nil {
			t.Error("GetObservation returned nil")
		}

		if retrieved.Content != "Observação clínica inicial para teste" {
			t.Errorf("GetObservation content = %v, want %v", retrieved.Content, "Observação clínica inicial para teste")
		}
	})

	// Test 2: Update observation
	t.Run("UpdateObservation", func(t *testing.T) {
		err := observationService.UpdateObservation(ctx, obs.ID, "Observação clínica atualizada após edição")
		if err != nil {
			t.Errorf("UpdateObservation failed: %v", err)
		}

		// Verify update
		updated, err := observationService.GetObservation(ctx, obs.ID)
		if err != nil {
			t.Errorf("GetObservation after update failed: %v", err)
		}

		if updated.Content != "Observação clínica atualizada após edição" {
			t.Errorf("UpdateObservation content = %v, want %v", updated.Content, "Observação clínica atualizada após edição")
		}

		// Verify updated_at was set
		if updated.UpdatedAt.IsZero() {
			t.Error("UpdateObservation didn't set UpdatedAt")
		}
	})

	// Test 3: Update with empty content (should fail)
	t.Run("UpdateObservationEmptyContent", func(t *testing.T) {
		err := observationService.UpdateObservation(ctx, obs.ID, "")
		if err == nil {
			t.Error("UpdateObservation with empty content should fail")
		}

		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("UpdateObservation with empty content error = %v, want error containing 'cannot be empty'", err)
		}
	})

	// Test 4: Update with content too long (should fail)
	t.Run("UpdateObservationContentTooLong", func(t *testing.T) {
		longContent := strings.Repeat("a", 5001)
		err := observationService.UpdateObservation(ctx, obs.ID, longContent)
		if err == nil {
			t.Error("UpdateObservation with content too long should fail")
		}

		if !strings.Contains(err.Error(), "cannot exceed 5000 characters") {
			t.Errorf("UpdateObservation with content too long error = %v, want error containing 'cannot exceed 5000 characters'", err)
		}
	})

	// Test 5: Update non-existent observation (should fail)
	t.Run("UpdateNonExistentObservation", func(t *testing.T) {
		err := observationService.UpdateObservation(ctx, "non-existent-id", "Conteúdo válido")
		if err == nil {
			t.Error("UpdateObservation with non-existent ID should fail")
		}

		if !strings.Contains(err.Error(), "observation not found") {
			t.Errorf("UpdateObservation with non-existent ID error = %v, want error containing 'observation not found'", err)
		}
	})

	// Test 6: Verify original created_at is preserved after update
	t.Run("CreatedAtPreservedAfterUpdate", func(t *testing.T) {
		// Get original observation
		original, err := observationService.GetObservation(ctx, obs.ID)
		if err != nil {
			t.Fatalf("Failed to get observation: %v", err)
		}

		originalCreatedAt := original.CreatedAt

		// Update observation
		err = observationService.UpdateObservation(ctx, obs.ID, "Outro conteúdo atualizado")
		if err != nil {
			t.Fatalf("Failed to update observation: %v", err)
		}

		// Get updated observation
		updated, err := observationService.GetObservation(ctx, obs.ID)
		if err != nil {
			t.Fatalf("Failed to get updated observation: %v", err)
		}

		// Verify created_at is unchanged
		if !updated.CreatedAt.Equal(originalCreatedAt) {
			t.Errorf("CreatedAt changed after update: original = %v, updated = %v", originalCreatedAt, updated.CreatedAt)
		}

		// Verify updated_at is after created_at
		if !updated.UpdatedAt.After(updated.CreatedAt) {
			t.Errorf("UpdatedAt should be after CreatedAt: CreatedAt = %v, UpdatedAt = %v", updated.CreatedAt, updated.UpdatedAt)
		}
	})
}

func TestObservationRepositorySchema(t *testing.T) {
	ctx := context.Background()

	// Test that the schema includes updated_at field
	db, err := sqlite.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize database with migrations
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	repo := sqlite.NewObservationRepository(db)

	// Create observation
	obs := &observation.Observation{
		SessionID: "test-session",
		Content:   "Test content",
	}

	if err := repo.Save(ctx, obs); err != nil {
		t.Fatalf("Failed to save observation: %v", err)
	}

	// Update observation
	obs.Content = "Updated content"
	if err := repo.Update(ctx, obs); err != nil {
		t.Fatalf("Failed to update observation: %v", err)
	}

	// Retrieve and verify updated_at was set
	updated, err := repo.FindByID(ctx, obs.ID)
	if err != nil {
		t.Fatalf("Failed to find observation: %v", err)
	}

	if updated.UpdatedAt.IsZero() {
		t.Error("updated_at field was not set after update")
	}
}
