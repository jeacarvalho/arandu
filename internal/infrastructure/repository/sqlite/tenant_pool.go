package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"arandu/internal/infrastructure/repository/sqlite/migrations"

	_ "modernc.org/sqlite"
)

const (
	TenantDBPrefix    = "clinical_"
	TenantDBExtension = ".db"
	IdleTimeout       = 10 * time.Minute
)

type TenantPool struct {
	mu           sync.RWMutex
	conns        map[string]*sql.DB
	storage      string
	migrationsFs interface{}
}

func NewTenantPool(storage string, migrationsFs interface{}) *TenantPool {
	return &TenantPool{
		conns:        make(map[string]*sql.DB),
		storage:      storage,
		migrationsFs: migrationsFs,
	}
}

func (tp *TenantPool) GetConnection(tenantID string) (*sql.DB, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}

	tp.mu.RLock()
	if db, ok := tp.conns[tenantID]; ok {
		if err := db.Ping(); err == nil {
			tp.mu.RUnlock()
			return db, nil
		}
		tp.mu.RUnlock()
	} else {
		tp.mu.RUnlock()
	}

	tp.mu.Lock()
	defer tp.mu.Unlock()

	if db, ok := tp.conns[tenantID]; ok {
		if err := db.Ping(); err == nil {
			return db, nil
		}
	}

	return tp.createConnection(tenantID)
}

func (tp *TenantPool) createConnection(tenantID string) (*sql.DB, error) {
	tenantsDir := filepath.Join(tp.storage, "tenants")
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tenants directory: %w", err)
	}

	dbPath := filepath.Join(tenantsDir, TenantDBPrefix+tenantID+TenantDBExtension)

	dsn := dbPath + "?_journal_mode=WAL"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open tenant database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping tenant database: %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := tp.runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	tp.conns[tenantID] = db
	return db, nil
}

func (tp *TenantPool) runMigrations(db *sql.DB) error {
	manager, err := migrations.NewMigrationManager(db)
	if err != nil {
		return fmt.Errorf("failed to create migration manager: %w", err)
	}
	return manager.Migrate()
}

func (tp *TenantPool) CloseConnection(tenantID string) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if db, ok := tp.conns[tenantID]; ok {
		delete(tp.conns, tenantID)
		return db.Close()
	}
	return nil
}

func (tp *TenantPool) CloseAll() error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	var errors []error
	for tenantID, db := range tp.conns {
		if err := db.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection for tenant %s: %w", tenantID, err))
		}
	}
	tp.conns = make(map[string]*sql.DB)

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

func (tp *TenantPool) GetActiveCount() int {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return len(tp.conns)
}

func (tp *TenantPool) IsConnected(tenantID string) bool {
	tp.mu.RLock()
	defer tp.mu.RUnlock()

	if db, ok := tp.conns[tenantID]; ok {
		return db.Ping() == nil
	}
	return false
}
