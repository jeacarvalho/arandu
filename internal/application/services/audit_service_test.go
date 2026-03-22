package services

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"arandu/internal/infrastructure/repository/sqlite"
	appcontext "arandu/internal/platform/context"
)

func TestAuditService_Log(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	auditService := NewAuditService(centralDB.DB)
	defer auditService.Close()

	userID := "test-user-123"
	tenantID := "test-tenant-456"
	patientID := "patient-789"

	ctx := context.Background()
	ctx = appcontext.WithUserID(ctx, userID)
	ctx = appcontext.WithTenantID(ctx, tenantID)

	t.Run("LogPatientAccess", func(t *testing.T) {
		auditService.Log(ctx, AuditActionAccessPatient, patientID)
		time.Sleep(100 * time.Millisecond)

		logs, err := auditService.GetLogsByTenant(ctx, tenantID, 10)
		if err != nil {
			t.Fatalf("GetLogsByTenant failed: %v", err)
		}

		if len(logs) == 0 {
			t.Fatal("Expected at least one audit log")
		}

		found := false
		for _, log := range logs {
			if log.Action == AuditActionAccessPatient && log.ResourceID == patientID {
				found = true
				if log.UserID != userID {
					t.Errorf("Expected userID %s, got %s", userID, log.UserID)
				}
				if log.TenantID != tenantID {
					t.Errorf("Expected tenantID %s, got %s", tenantID, log.TenantID)
				}
				break
			}
		}

		if !found {
			t.Errorf("Expected audit log with action=%s and resourceID=%s", AuditActionAccessPatient, patientID)
		}
	})

	t.Run("LogPatientCreate", func(t *testing.T) {
		newPatientID := "new-patient-abc"
		auditService.Log(ctx, AuditActionCreatePatient, newPatientID)
		time.Sleep(100 * time.Millisecond)

		logs, err := auditService.GetLogsByTenant(ctx, tenantID, 10)
		if err != nil {
			t.Fatalf("GetLogsByTenant failed: %v", err)
		}

		found := false
		for _, log := range logs {
			if log.Action == AuditActionCreatePatient && log.ResourceID == newPatientID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected audit log with action=%s", AuditActionCreatePatient)
		}
	})
}

func TestAuditService_AsyncBehavior(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	auditService := NewAuditService(centralDB.DB)
	defer auditService.Close()

	ctx := context.Background()
	ctx = appcontext.WithUserID(ctx, "user-test")
	ctx = appcontext.WithTenantID(ctx, "tenant-test")

	t.Run("LogReturnsWithoutBlocking", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < 100; i++ {
			auditService.Log(ctx, AuditActionAccessPatient, "patient-"+string(rune(i)))
		}
		elapsed := time.Since(start)

		if elapsed > 100*time.Millisecond {
			t.Errorf("Expected log to return quickly, took %v", elapsed)
		}
	})
}

func TestAuditService_Close(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	centralDB.Migrate(nil)

	auditService := NewAuditService(centralDB.DB)

	ctx := context.Background()
	ctx = appcontext.WithUserID(ctx, "user-test")
	ctx = appcontext.WithTenantID(ctx, "tenant-test")

	for i := 0; i < 10; i++ {
		auditService.Log(ctx, AuditActionAccessPatient, "patient-"+string(rune(i)))
	}

	if err := auditService.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestAuditService_LogWithoutContext(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	centralDB.Migrate(nil)

	auditService := NewAuditService(centralDB.DB)
	defer auditService.Close()

	ctx := context.Background()

	auditService.Log(ctx, AuditActionAccessPatient, "patient-123")

	time.Sleep(100 * time.Millisecond)

	logs, err := auditService.GetLogsByTenant(ctx, "tenant-test", 10)
	if err != nil {
		t.Fatalf("GetLogsByTenant failed: %v", err)
	}

	if len(logs) > 0 {
		t.Error("Expected no logs when context is missing user_id/tenant_id")
	}
}

func TestAuditService_GetLogsByTenant(t *testing.T) {
	tempDir := t.TempDir()
	storagePath := filepath.Join(tempDir, "storage")

	centralDB, err := sqlite.NewCentralDB(storagePath)
	if err != nil {
		t.Fatalf("Failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("Failed to migrate central DB: %v", err)
	}

	auditService := NewAuditService(centralDB.DB)
	defer auditService.Close()

	ctx := context.Background()
	ctx = appcontext.WithUserID(ctx, "user-1")
	ctx = appcontext.WithTenantID(ctx, "tenant-1")

	for i := 0; i < 5; i++ {
		auditService.Log(ctx, AuditActionAccessPatient, "patient-"+string(rune(i)))
	}

	time.Sleep(200 * time.Millisecond)

	logs, err := auditService.GetLogsByTenant(ctx, "tenant-1", 3)
	if err != nil {
		t.Fatalf("GetLogsByTenant failed: %v", err)
	}

	if len(logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}

	logs, err = auditService.GetLogsByTenant(ctx, "tenant-1", 100)
	if err != nil {
		t.Fatalf("GetLogsByTenant failed: %v", err)
	}

	if len(logs) != 5 {
		t.Errorf("Expected 5 logs, got %d", len(logs))
	}
}
