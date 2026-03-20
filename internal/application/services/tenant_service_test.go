package services

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"arandu/internal/infrastructure/repository/sqlite"
)

func TestTenantService_ProvisionNewTenant(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	// Run central migrations
	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	// Create test user with NULL tenant_id
	userID := "test-user-123"
	_, err = centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, NULL)`,
		userID, "test@example.com")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	service := NewTenantService(centralDB.DB, storagePath)

	// Test
	ctx := context.Background()
	tenantID, err := service.ProvisionNewTenant(ctx, userID)
	if err != nil {
		t.Fatalf("ProvisionNewTenant failed: %v", err)
	}

	// Verify tenant was created in central DB
	var dbPath, status string
	err = centralDB.QueryRow(`SELECT db_path, status FROM tenants WHERE id = ?`, tenantID).Scan(&dbPath, &status)
	if err != nil {
		t.Fatalf("Failed to query tenant: %v", err)
	}

	if status != "active" {
		t.Errorf("Expected status 'active', got %s", status)
	}

	expectedPath := filepath.Join(storagePath, "tenants", "clinical_"+tenantID+".db")
	if dbPath != expectedPath {
		t.Errorf("Expected db_path %s, got %s", expectedPath, dbPath)
	}

	// Verify tenant DB file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Tenant DB file does not exist: %s", dbPath)
	}

	// Verify user was updated with tenant_id
	var userTenantID string
	err = centralDB.QueryRow(`SELECT tenant_id FROM users WHERE id = ?`, userID).Scan(&userTenantID)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	if userTenantID != tenantID {
		t.Errorf("Expected user tenant_id %s, got %s", tenantID, userTenantID)
	}

	// Verify tenant DB has tables
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		t.Fatalf("Failed to open tenant DB: %v", err)
	}
	defer db.Close()

	var tableCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&tableCount)
	if err != nil {
		t.Fatalf("Failed to count tables: %v", err)
	}

	if tableCount == 0 {
		t.Error("Tenant DB has no tables")
	}
}

func TestTenantService_GetTenantByUserID(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	// Run central migrations
	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	// Create test tenant and user
	tenantID := "test-tenant-123"
	dbPath := filepath.Join(storagePath, "tenants", "clinical_"+tenantID+".db")

	_, err = centralDB.Exec(`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
		tenantID, dbPath)
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}

	userID := "test-user-456"
	_, err = centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, ?)`,
		userID, "user@example.com", tenantID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	service := NewTenantService(centralDB.DB, storagePath)

	// Test
	ctx := context.Background()
	tenant, err := service.GetTenantByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetTenantByUserID failed: %v", err)
	}

	if tenant == nil {
		t.Fatal("Expected tenant, got nil")
	}

	if tenant.ID != tenantID {
		t.Errorf("Expected tenant ID %s, got %s", tenantID, tenant.ID)
	}

	if tenant.DBPath != dbPath {
		t.Errorf("Expected DB path %s, got %s", dbPath, tenant.DBPath)
	}

	if tenant.Status != "active" {
		t.Errorf("Expected status 'active', got %s", tenant.Status)
	}
}

func TestTenantService_TenantExists(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	// Run central migrations
	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	// Create test tenant
	tenantID := "test-tenant-789"
	dbPath := filepath.Join(storagePath, "tenants", "clinical_"+tenantID+".db")

	_, err = centralDB.Exec(`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
		tenantID, dbPath)
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}

	service := NewTenantService(centralDB.DB, storagePath)

	// Test existing tenant
	ctx := context.Background()
	exists, err := service.TenantExists(ctx, tenantID)
	if err != nil {
		t.Fatalf("TenantExists failed: %v", err)
	}

	if !exists {
		t.Error("Expected tenant to exist")
	}

	// Test non-existing tenant
	exists, err = service.TenantExists(ctx, "non-existing-tenant")
	if err != nil {
		t.Fatalf("TenantExists failed for non-existing tenant: %v", err)
	}

	if exists {
		t.Error("Expected non-existing tenant to not exist")
	}
}

func TestTenantService_ValidateTenantDB(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	// Run central migrations
	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	service := NewTenantService(centralDB.DB, storagePath)

	// Create a tenant DB manually
	tenantID := "test-tenant-validate"
	dbPath := service.GetTenantDBPath(tenantID)

	// Create directory
	tenantsDir := filepath.Join(storagePath, "tenants")
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		t.Fatalf("Failed to create tenants directory: %v", err)
	}

	// Create empty DB file
	file, err := os.Create(dbPath)
	if err != nil {
		t.Fatalf("Failed to create DB file: %v", err)
	}
	file.Close()

	// Test validation should fail because DB has no tables
	ctx := context.Background()
	err = service.ValidateTenantDB(ctx, tenantID)
	if err == nil {
		t.Error("Expected validation to fail for empty DB")
	}

	// Now create a proper tenant
	_, err = centralDB.Exec(`INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')`,
		tenantID, dbPath)
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}

	// Provision the tenant to create tables
	userID := "test-user-validate"
	_, err = centralDB.Exec(`INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, NULL, NULL)`,
		userID, "validate@example.com")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	provisionedTenantID, err := service.ProvisionNewTenant(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to provision tenant: %v", err)
	}

	// Test validation should now succeed for the provisioned tenant
	err = service.ValidateTenantDB(ctx, provisionedTenantID)
	if err != nil {
		t.Errorf("Validation failed for proper DB: %v", err)
	}
}
