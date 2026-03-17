package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"arandu/internal/infrastructure/repository/sqlite/migrations"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Migrate applies all pending database migrations
func (db *DB) Migrate() error {
	manager, err := migrations.NewMigrationManager(db.DB)
	if err != nil {
		return fmt.Errorf("failed to create migration manager: %w", err)
	}

	return manager.Migrate()
}

// MigrationStatus returns the current migration status
func (db *DB) MigrationStatus() (map[string]string, error) {
	manager, err := migrations.NewMigrationManager(db.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration manager: %w", err)
	}

	return manager.Status()
}
