package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"arandu/internal/infrastructure/repository/sqlite/migrations"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Tenant struct {
	ID        string
	DBPath    string
	Status    string
	CreatedAt time.Time
}

type TenantService struct {
	centralDB *sql.DB
	storage   string
}

func NewTenantService(centralDB *sql.DB, storage string) *TenantService {
	return &TenantService{
		centralDB: centralDB,
		storage:   storage,
	}
}

func (s *TenantService) ProvisionNewTenant(ctx context.Context, userID string) (string, error) {
	tenantID := uuid.New().String()
	dbPath := filepath.Join(s.storage, "tenants", fmt.Sprintf("clinical_%s.db", tenantID))

	tenantsDir := filepath.Join(s.storage, "tenants")
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create tenants directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return "", fmt.Errorf("failed to open tenant database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return "", fmt.Errorf("failed to ping tenant database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return "", fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	migrator, err := migrations.NewMigrationManager(db)
	if err != nil {
		return "", fmt.Errorf("failed to create migration manager: %w", err)
	}

	if err := migrator.Migrate(); err != nil {
		return "", fmt.Errorf("failed to run migrations: %w", err)
	}

	tx, err := s.centralDB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	tenantQuery := `INSERT INTO tenants (id, db_path, status, created_at) VALUES (?, ?, 'active', ?)`
	if _, err := tx.ExecContext(ctx, tenantQuery, tenantID, dbPath, time.Now()); err != nil {
		return "", fmt.Errorf("failed to register tenant: %w", err)
	}

	userQuery := `UPDATE users SET tenant_id = ? WHERE id = ?`
	if _, err := tx.ExecContext(ctx, userQuery, tenantID, userID); err != nil {
		return "", fmt.Errorf("failed to update user tenant: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tenantID, nil
}

func (s *TenantService) GetTenantByUserID(ctx context.Context, userID string) (*Tenant, error) {
	query := `
		SELECT t.id, t.db_path, t.status, t.created_at
		FROM tenants t
		INNER JOIN users u ON t.id = u.tenant_id
		WHERE u.id = ?
	`

	var tenant Tenant
	err := s.centralDB.QueryRowContext(ctx, query, userID).Scan(
		&tenant.ID,
		&tenant.DBPath,
		&tenant.Status,
		&tenant.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return &tenant, nil
}

func (s *TenantService) TenantExists(ctx context.Context, tenantID string) (bool, error) {
	query := `SELECT COUNT(*) FROM tenants WHERE id = ?`
	var count int
	err := s.centralDB.QueryRowContext(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check tenant existence: %w", err)
	}
	return count > 0, nil
}

func (s *TenantService) GetTenantDBPath(tenantID string) string {
	return filepath.Join(s.storage, "tenants", fmt.Sprintf("clinical_%s.db", tenantID))
}

func (s *TenantService) ValidateTenantDB(ctx context.Context, tenantID string) error {
	dbPath := s.GetTenantDBPath(tenantID)

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("tenant database file does not exist: %s", dbPath)
	}

	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to open tenant database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping tenant database: %w", err)
	}

	var tableCount int
	query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`
	if err := db.QueryRowContext(ctx, query).Scan(&tableCount); err != nil {
		return fmt.Errorf("failed to check tables: %w", err)
	}

	if tableCount == 0 {
		return fmt.Errorf("tenant database has no tables")
	}

	return nil
}
