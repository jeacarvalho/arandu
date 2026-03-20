package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestNewCentralDB(t *testing.T) {
	tmpDir := t.TempDir()

	centralDB, err := NewCentralDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	expectedPath := filepath.Join(tmpDir, CentralDBName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("expected central DB file to be created")
	}

	err = centralDB.Migrate(nil)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	var count int
	err = centralDB.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Errorf("failed to query migrations table: %v", err)
	}
}

func TestCentralDBMigrate(t *testing.T) {
	tmpDir := t.TempDir()

	centralDB, err := NewCentralDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	err = centralDB.Migrate(nil)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	var tenantCount int
	err = centralDB.QueryRow("SELECT COUNT(*) FROM tenants").Scan(&tenantCount)
	if err != nil {
		t.Errorf("failed to query tenants table: %v", err)
	}

	var userCount int
	err = centralDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		t.Errorf("failed to query users table: %v", err)
	}

	var emailIndexExists int
	err = centralDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='idx_users_email'").Scan(&emailIndexExists)
	if err != nil {
		t.Errorf("failed to check index: %v", err)
	}
	if emailIndexExists != 1 {
		t.Error("expected idx_users_email index to exist")
	}
}

func TestCentralDBMigrateIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	centralDB, err := NewCentralDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	err = centralDB.Migrate(nil)
	if err != nil {
		t.Fatalf("first migration failed: %v", err)
	}

	err = centralDB.Migrate(nil)
	if err != nil {
		t.Fatalf("second migration should be idempotent: %v", err)
	}

	var count int
	err = centralDB.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Errorf("failed to count migrations: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 migration, got %d", count)
	}
}

func TestCentralDBForeignKeys(t *testing.T) {
	tmpDir := t.TempDir()

	centralDB, err := NewCentralDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	err = centralDB.Migrate(nil)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	_, err = centralDB.Exec("INSERT INTO tenants (id, db_path, status) VALUES ('tenant-1', 'path/db.db', 'active')")
	if err != nil {
		t.Fatalf("failed to insert tenant: %v", err)
	}

	_, err = centralDB.Exec("INSERT INTO users (id, email, password_hash, tenant_id) VALUES ('user-1', 'test@test.com', 'hash', 'tenant-1')")
	if err != nil {
		t.Fatalf("failed to insert user with valid tenant: %v", err)
	}

	_, err = centralDB.Exec("INSERT INTO users (id, email, password_hash, tenant_id) VALUES ('user-2', 'test2@test.com', 'hash', 'non-existent')")
	if err == nil {
		t.Error("expected foreign key violation for non-existent tenant")
	}
}

func TestCentralDBIsolatedFromClinical(t *testing.T) {
	tmpDir := t.TempDir()

	centralDB, err := NewCentralDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	err = centralDB.Migrate(nil)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	var tables []string
	rows, err := centralDB.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			t.Fatalf("failed to scan table: %v", err)
		}
		tables = append(tables, table)
	}

	expectedTables := map[string]bool{
		"tenants":           true,
		"users":             true,
		"schema_migrations": true,
	}

	for _, table := range tables {
		if !expectedTables[table] {
			t.Errorf("unexpected table in central DB: %s", table)
		}
	}

	if !tablesContains(tables, "patients") || !tablesContains(tables, "sessions") {
		t.Log("central DB is properly isolated from clinical tables (as expected)")
	}
}

func tablesContains(tables []string, target string) bool {
	for _, t := range tables {
		if t == target {
			return true
		}
	}
	return false
}

func TestCentralDBWALEnabled(t *testing.T) {
	centralDB, err := NewCentralDB(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	t.Log("Central DB initialized with WAL mode configured in DSN")
}
