package shared

import (
	"testing"
)

func TestNewTenant(t *testing.T) {
	tenant := NewTenant("test-uuid", "storage/tenants/clinical_test.db")

	if tenant.ID != "test-uuid" {
		t.Errorf("expected ID 'test-uuid', got '%s'", tenant.ID)
	}
	if tenant.DBPath != "storage/tenants/clinical_test.db" {
		t.Errorf("expected DBPath 'storage/tenants/clinical_test.db', got '%s'", tenant.DBPath)
	}
	if tenant.Status != TenantStatusActive {
		t.Errorf("expected status 'active', got '%s'", tenant.Status)
	}
	if tenant.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestTenantIsActive(t *testing.T) {
	tenant := NewTenant("id", "path")

	if !tenant.IsActive() {
		t.Error("expected active tenant to return true")
	}

	tenant.Status = TenantStatusSuspended
	if tenant.IsActive() {
		t.Error("expected suspended tenant to return false")
	}
}

func TestTenantSuspend(t *testing.T) {
	tenant := NewTenant("id", "path")

	err := tenant.Suspend()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tenant.Status != TenantStatusSuspended {
		t.Errorf("expected status 'suspended', got '%s'", tenant.Status)
	}

	err = tenant.Suspend()
	if err == nil {
		t.Error("expected error when suspending already suspended tenant")
	}
}

func TestTenantActivate(t *testing.T) {
	tenant := NewTenant("id", "path")
	tenant.Status = TenantStatusSuspended

	err := tenant.Activate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tenant.Status != TenantStatusActive {
		t.Errorf("expected status 'active', got '%s'", tenant.Status)
	}

	err = tenant.Activate()
	if err == nil {
		t.Error("expected error when activating already active tenant")
	}
}
