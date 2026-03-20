package sqlite

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestTenantPool_New(t *testing.T) {
	pool := NewTenantPool("storage", nil)

	if pool == nil {
		t.Error("expected non-nil pool")
	}
	if pool.conns == nil {
		t.Error("expected non-nil connections map")
	}
}

func TestTenantPool_GetActiveCount(t *testing.T) {
	pool := NewTenantPool(t.TempDir(), nil)

	count := pool.GetActiveCount()
	if count != 0 {
		t.Errorf("expected 0 active connections, got %d", count)
	}
}

func TestTenantPool_IsConnected(t *testing.T) {
	pool := NewTenantPool(t.TempDir(), nil)

	if pool.IsConnected("non-existent") {
		t.Error("expected false for non-existent tenant")
	}
}

func TestTenantPool_CloseConnection(t *testing.T) {
	pool := NewTenantPool(t.TempDir(), nil)

	err := pool.CloseConnection("non-existent")
	if err != nil {
		t.Errorf("expected no error for closing non-existent connection: %v", err)
	}
}

func TestTenantPool_CloseAll(t *testing.T) {
	pool := NewTenantPool(t.TempDir(), nil)

	err := pool.CloseAll()
	if err != nil {
		t.Errorf("expected no error: %v", err)
	}

	if pool.GetActiveCount() != 0 {
		t.Error("expected 0 connections after CloseAll")
	}
}

func TestTenantPool_GetConnection_EmptyTenantID(t *testing.T) {
	pool := NewTenantPool(t.TempDir(), nil)

	_, err := pool.GetConnection("")
	if err == nil {
		t.Error("expected error for empty tenant ID")
	}
}

func TestTenantPool_GetConnection_NewTenant(t *testing.T) {
	tmpDir := t.TempDir()

	tenantsDir := filepath.Join(tmpDir, "tenants")
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		t.Fatalf("failed to create tenants dir: %v", err)
	}

	clinicalDB := filepath.Join(tenantsDir, "clinical_test-tenant.db")
	db, err := createTestDB(clinicalDB)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	defer db.Close()

	pool := NewTenantPool(tmpDir, nil)

	conn, err := pool.GetConnection("test-tenant")
	if err != nil {
		t.Fatalf("failed to get connection: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil connection")
	}

	if pool.GetActiveCount() != 1 {
		t.Errorf("expected 1 active connection, got %d", pool.GetActiveCount())
	}

	if !pool.IsConnected("test-tenant") {
		t.Error("expected tenant to be connected")
	}
}

func TestTenantPool_GetConnection_ReuseConnection(t *testing.T) {
	tmpDir := t.TempDir()

	tenantsDir := filepath.Join(tmpDir, "tenants")
	if err := os.MkdirAll(tenantsDir, 0755); err != nil {
		t.Fatalf("failed to create tenants dir: %v", err)
	}

	clinicalDB := filepath.Join(tenantsDir, "clinical_test-tenant2.db")
	db, err := createTestDB(clinicalDB)
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}
	defer db.Close()

	pool := NewTenantPool(tmpDir, nil)

	conn1, _ := pool.GetConnection("test-tenant2")
	conn2, _ := pool.GetConnection("test-tenant2")

	if conn1 != conn2 {
		t.Error("expected same connection to be reused")
	}

	if pool.GetActiveCount() != 1 {
		t.Errorf("expected 1 active connection (reused), got %d", pool.GetActiveCount())
	}
}

func createTestDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path+"?_journal_mode=WAL")
}
