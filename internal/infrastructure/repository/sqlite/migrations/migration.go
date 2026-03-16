package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
)

//go:embed *.sql
var migrationsFS embed.FS

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) (*MigrationManager, error) {
	manager := &MigrationManager{
		db: db,
	}

	// Load migrations from embedded filesystem
	if err := manager.loadMigrations(migrationsFS); err != nil {
		return nil, err
	}

	// Ensure schema_migrations table exists
	if err := manager.createMigrationsTable(); err != nil {
		return nil, err
	}

	return manager, nil
}

// loadMigrations loads migration files from the filesystem
func (m *MigrationManager) loadMigrations(migrationsFS fs.FS) error {
	entries, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Group files by version
	migrationFiles := make(map[string]map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		// Parse filename: 0001_create_patients_table.up.sql
		parts := strings.Split(filename, ".")
		if len(parts) != 3 {
			continue // Skip invalid filenames
		}

		version := parts[0]
		direction := parts[1] // "up" or "down"

		if migrationFiles[version] == nil {
			migrationFiles[version] = make(map[string]string)
		}

		// Read SQL content
		content, err := fs.ReadFile(migrationsFS, filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		migrationFiles[version][direction] = string(content)
	}

	// Convert to Migration structs
	for version, files := range migrationFiles {
		migration := Migration{
			Version: version,
			Name:    extractMigrationName(version),
		}

		if upSQL, ok := files["up"]; ok {
			migration.UpSQL = upSQL
		}

		if downSQL, ok := files["down"]; ok {
			migration.DownSQL = downSQL
		}

		// Only add migration if it has both up and down SQL
		if migration.UpSQL != "" && migration.DownSQL != "" {
			m.migrations = append(m.migrations, migration)
		}
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	return nil
}

// extractMigrationName extracts the name from version string
// Example: "0001_create_patients_table" -> "create_patients_table"
func extractMigrationName(version string) string {
	parts := strings.SplitN(version, "_", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return version
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func (m *MigrationManager) createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := m.db.Exec(query)
	return err
}

// Migrate applies all pending migrations
func (m *MigrationManager) Migrate() error {
	applied, err := m.appliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range m.migrations {
		if _, ok := applied[migration.Version]; ok {
			continue // Already applied
		}

		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Apply migration
		if _, err := tx.Exec(migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		// Record migration
		recordQuery := `INSERT INTO schema_migrations (version) VALUES (?)`
		if _, err := tx.Exec(recordQuery, migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Applied migration: %s\n", migration.Version)
	}

	return nil
}

// Rollback reverts a specific migration
func (m *MigrationManager) Rollback(version string) error {
	applied, err := m.appliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find the migration
	var migration *Migration
	for i := len(m.migrations) - 1; i >= 0; i-- {
		if m.migrations[i].Version == version {
			migration = &m.migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found", version)
	}

	if _, ok := applied[version]; !ok {
		return fmt.Errorf("migration %s not applied", version)
	}

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Revert migration
	if _, err := tx.Exec(migration.DownSQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to revert migration %s: %w", migration.Version, err)
	}

	// Remove migration record
	deleteQuery := `DELETE FROM schema_migrations WHERE version = ?`
	if _, err := tx.Exec(deleteQuery, migration.Version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %s: %w", migration.Version, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback %s: %w", migration.Version, err)
	}

	fmt.Printf("Reverted migration: %s\n", migration.Version)
	return nil
}

// RollbackLast reverts the last applied migration
func (m *MigrationManager) RollbackLast() error {
	applied, err := m.appliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find the last applied migration
	var lastVersion string
	for _, migration := range m.migrations {
		if _, ok := applied[migration.Version]; ok {
			lastVersion = migration.Version
		}
	}

	if lastVersion == "" {
		return fmt.Errorf("no migrations to rollback")
	}

	return m.Rollback(lastVersion)
}

// CurrentVersion returns the version of the last applied migration
func (m *MigrationManager) CurrentVersion() (string, error) {
	query := `SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1`
	var version string
	err := m.db.QueryRow(query).Scan(&version)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get current version: %w", err)
	}
	return version, nil
}

// PendingMigrations returns a list of migration versions that haven't been applied
func (m *MigrationManager) PendingMigrations() ([]string, error) {
	applied, err := m.appliedMigrations()
	if err != nil {
		return nil, err
	}

	var pending []string
	for _, migration := range m.migrations {
		if _, ok := applied[migration.Version]; !ok {
			pending = append(pending, migration.Version)
		}
	}

	return pending, nil
}

// appliedMigrations returns a map of applied migration versions
func (m *MigrationManager) appliedMigrations() (map[string]bool, error) {
	query := `SELECT version FROM schema_migrations`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

// Status returns the current migration status
func (m *MigrationManager) Status() (map[string]string, error) {
	applied, err := m.appliedMigrations()
	if err != nil {
		return nil, err
	}

	status := make(map[string]string)
	for _, migration := range m.migrations {
		if _, ok := applied[migration.Version]; ok {
			status[migration.Version] = "applied"
		} else {
			status[migration.Version] = "pending"
		}
	}

	return status, nil
}

// NewMigrationManagerFromDir is kept for backward compatibility with tests
// In production, use NewMigrationManager instead
func NewMigrationManagerFromDir(db *sql.DB, migrationsDir string) (*MigrationManager, error) {
	// For tests, still use directory-based approach
	migrationsFS := os.DirFS(migrationsDir)

	manager := &MigrationManager{
		db: db,
	}

	// Load migrations from filesystem
	if err := manager.loadMigrations(migrationsFS); err != nil {
		return nil, err
	}

	// Ensure schema_migrations table exists
	if err := manager.createMigrationsTable(); err != nil {
		return nil, err
	}

	return manager, nil
}
