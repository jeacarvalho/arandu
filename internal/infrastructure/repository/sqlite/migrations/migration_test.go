package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, string) {
	t.Helper()

	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "migration-test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()

	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		os.Remove(tmpfile.Name())
		t.Fatalf("Failed to open database: %v", err)
	}

	return db, tmpfile.Name()
}

func cleanupTestDB(t *testing.T, db *sql.DB, dbPath string) {
	t.Helper()

	if db != nil {
		db.Close()
	}
	if dbPath != "" {
		os.Remove(dbPath)
	}
}

func TestMigrationManager(t *testing.T) {
	db, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, db, dbPath)

	// Get the migrations directory path
	migrationsDir := filepath.Join("..", "..", "..", "..", "..", "internal", "infrastructure", "repository", "sqlite", "migrations")
	migrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	t.Run("Create migration manager", func(t *testing.T) {
		manager, err := NewMigrationManagerFromDir(db, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to create migration manager: %v", err)
		}

		if manager == nil {
			t.Fatal("Migration manager is nil")
		}
	})

	t.Run("Check initial status", func(t *testing.T) {
		manager, err := NewMigrationManagerFromDir(db, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to create migration manager: %v", err)
		}

		// Check that schema_migrations table was created
		var tableExists bool
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='schema_migrations'`
		err = db.QueryRow(query).Scan(&tableExists)
		if err != nil {
			t.Fatalf("Failed to check schema_migrations table: %v", err)
		}

		if !tableExists {
			t.Error("schema_migrations table was not created")
		}

		// Check current version (should be empty)
		version, err := manager.CurrentVersion()
		if err != nil {
			t.Fatalf("Failed to get current version: %v", err)
		}

		if version != "" {
			t.Errorf("Expected empty version, got %q", version)
		}

		// Check pending migrations
		pending, err := manager.PendingMigrations()
		if err != nil {
			t.Fatalf("Failed to get pending migrations: %v", err)
		}

		if len(pending) == 0 {
			t.Error("Expected at least one pending migration")
		}

		// Check status
		status, err := manager.Status()
		if err != nil {
			t.Fatalf("Failed to get status: %v", err)
		}

		if len(status) == 0 {
			t.Error("Expected migration status")
		}
	})

	t.Run("Apply migrations", func(t *testing.T) {
		manager, err := NewMigrationManagerFromDir(db, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to create migration manager: %v", err)
		}

		// Apply migrations
		if err := manager.Migrate(); err != nil {
			t.Fatalf("Failed to apply migrations: %v", err)
		}

		// Check current version
		version, err := manager.CurrentVersion()
		if err != nil {
			t.Fatalf("Failed to get current version: %v", err)
		}

		if version != "0001_initial_schema" {
			t.Errorf("Expected version 0001_initial_schema, got %q", version)
		}

		// Check pending migrations (should be empty)
		pending, err := manager.PendingMigrations()
		if err != nil {
			t.Fatalf("Failed to get pending migrations: %v", err)
		}

		if len(pending) != 0 {
			t.Errorf("Expected no pending migrations, got %v", pending)
		}

		// Check that patients table was created
		var tableExists bool
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='patients'`
		err = db.QueryRow(query).Scan(&tableExists)
		if err != nil {
			t.Fatalf("Failed to check patients table: %v", err)
		}

		if !tableExists {
			t.Error("patients table was not created")
		}

		// Check that indexes were created
		var indexCount int
		query = `SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name LIKE 'idx_patients_%'`
		err = db.QueryRow(query).Scan(&indexCount)
		if err != nil {
			t.Fatalf("Failed to check indexes: %v", err)
		}

		if indexCount < 2 {
			t.Errorf("Expected at least 2 indexes, got %d", indexCount)
		}
	})

	t.Run("Rollback migration", func(t *testing.T) {
		manager, err := NewMigrationManagerFromDir(db, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to create migration manager: %v", err)
		}

		// Rollback the migration
		if err := manager.Rollback("0001_initial_schema"); err != nil {
			t.Fatalf("Failed to rollback migration: %v", err)
		}

		// Check that patients table was removed
		var tableExists bool
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='patients'`
		err = db.QueryRow(query).Scan(&tableExists)
		if err != nil {
			t.Fatalf("Failed to check patients table: %v", err)
		}

		if tableExists {
			t.Error("patients table should have been removed")
		}

		// Check current version (should be empty)
		version, err := manager.CurrentVersion()
		if err != nil {
			t.Fatalf("Failed to get current version: %v", err)
		}

		if version != "" {
			t.Errorf("Expected empty version after rollback, got %q", version)
		}
	})

	t.Run("Re-apply migrations after rollback", func(t *testing.T) {
		manager, err := NewMigrationManagerFromDir(db, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to create migration manager: %v", err)
		}

		// Apply migrations again
		if err := manager.Migrate(); err != nil {
			t.Fatalf("Failed to re-apply migrations: %v", err)
		}

		// Check that patients table was re-created
		var tableExists bool
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='patients'`
		err = db.QueryRow(query).Scan(&tableExists)
		if err != nil {
			t.Fatalf("Failed to check patients table: %v", err)
		}

		if !tableExists {
			t.Error("patients table should have been re-created")
		}
	})

	t.Run("Rollback last migration", func(t *testing.T) {
		manager, err := NewMigrationManagerFromDir(db, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to create migration manager: %v", err)
		}

		// Rollback last migration
		if err := manager.RollbackLast(); err != nil {
			t.Fatalf("Failed to rollback last migration: %v", err)
		}

		// Check that patients table was removed
		var tableExists bool
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='patients'`
		err = db.QueryRow(query).Scan(&tableExists)
		if err != nil {
			t.Fatalf("Failed to check patients table: %v", err)
		}

		if tableExists {
			t.Error("patients table should have been removed by RollbackLast")
		}
	})
}

func TestMigrationManager_ErrorCases(t *testing.T) {
	db, dbPath := setupTestDB(t)
	defer cleanupTestDB(t, db, dbPath)

	migrationsDir := filepath.Join("..", "..", "..", "..", "..", "internal", "infrastructure", "repository", "sqlite", "migrations")
	migrationsDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	manager, err := NewMigrationManagerFromDir(db, migrationsDir)
	if err != nil {
		t.Fatalf("Failed to create migration manager: %v", err)
	}

	t.Run("Rollback non-existent migration", func(t *testing.T) {
		err := manager.Rollback("non-existent")
		if err == nil {
			t.Error("Expected error when rolling back non-existent migration")
		}
	})

	t.Run("Rollback not-applied migration", func(t *testing.T) {
		// Don't apply migrations first
		err := manager.Rollback("0001_initial_schema")
		if err == nil {
			t.Error("Expected error when rolling back not-applied migration")
		}
	})

	t.Run("Migrate twice (idempotency)", func(t *testing.T) {
		// Apply migrations
		if err := manager.Migrate(); err != nil {
			t.Fatalf("Failed to apply migrations: %v", err)
		}

		// Apply again - should not error
		if err := manager.Migrate(); err != nil {
			t.Fatalf("Second migrate should be idempotent, got error: %v", err)
		}
	})
}
