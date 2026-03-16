package sqlite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDB(t *testing.T) {
	t.Run("Creates database successfully", func(t *testing.T) {
		tmpfile, err := os.CreateTemp("", "testdb-*.db")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpfile.Close()
		defer os.Remove(tmpfile.Name())

		db, err := NewDB(tmpfile.Name())
		if err != nil {
			t.Fatalf("Failed to create database: %v", err)
		}
		defer db.Close()

		if db == nil {
			t.Error("Expected database to be created")
		}
	})

	t.Run("Creates directory if it doesn't exist", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "testdir-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		dbPath := filepath.Join(tmpDir, "subdir", "test.db")
		db, err := NewDB(dbPath)
		if err != nil {
			t.Fatalf("Failed to create database in non-existent directory: %v", err)
		}
		defer db.Close()

		// Verify directory was created
		if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
			t.Error("Expected directory to be created")
		}
	})

	t.Run("Returns error for invalid database path", func(t *testing.T) {
		// Try to create database with an invalid path
		// Using a path that's likely to fail (root directory without permissions)
		_, err := NewDB("/")
		if err == nil {
			t.Error("Expected error for invalid database path")
		}
		// Note: Empty path might not error immediately, so we test with root
	})
}

func TestDB_Close(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}

	// Verify database is closed by trying to ping it
	err = db.DB.Ping()
	if err == nil {
		t.Error("Expected error when pinging closed database")
	}
}

func TestDB_Migrate(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test migration with embedded migrations
	err = db.Migrate()
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Verify tables were created
	tables := []string{"patients", "sessions", "observations", "interventions", "insights"}
	for _, table := range tables {
		var tableExists bool
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`
		err = db.QueryRow(query, table).Scan(&tableExists)
		if err != nil {
			t.Fatalf("Failed to check %s table: %v", table, err)
		}
		if !tableExists {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestDB_Migrate_InvalidDirectory(t *testing.T) {
	// This test is no longer needed since Migrate() doesn't take a directory parameter
	// Keeping it as a placeholder to show the method signature changed
	t.Skip("Migrate() no longer takes a directory parameter - using embedded migrations")
}

func TestDB_MigrationStatus(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	db, err := NewDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Get migration status with embedded migrations
	status, err := db.MigrationStatus()
	if err != nil {
		t.Fatalf("Migration status failed: %v", err)
	}

	// Status should be a map
	if status == nil {
		t.Error("Migration status should not be nil")
	}

	// Before migration, status should show pending migrations
	if len(status) == 0 {
		t.Error("Expected migration status entries")
	}

	// Run migration
	if err := db.Migrate(); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Get status after migration
	status, err = db.MigrationStatus()
	if err != nil {
		t.Fatalf("Migration status after migration failed: %v", err)
	}

	// Should have at least one applied migration
	hasApplied := false
	for _, s := range status {
		if s == "applied" {
			hasApplied = true
			break
		}
	}
	if !hasApplied {
		t.Error("Expected at least one applied migration after Migrate()")
	}
}

func TestDB_MigrationStatus_InvalidDirectory(t *testing.T) {
	// This test is no longer needed since MigrationStatus() doesn't take a directory parameter
	// Keeping it as a placeholder to show the method signature changed
	t.Skip("MigrationStatus() no longer takes a directory parameter - using embedded migrations")
}
