package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTenantPath(t *testing.T) {
	tests := []struct {
		tenantID string
		expected string
	}{
		{
			tenantID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			expected: "storage/tenants/a1/b2/clinical_a1b2c3d4-e5f6-7890-abcd-ef1234567890.db",
		},
		{
			tenantID: "fe9beee6-b407-4bd8-bef3-fbca8f0d319d",
			expected: "storage/tenants/fe/9b/clinical_fe9beee6-b407-4bd8-bef3-fbca8f0d319d.db",
		},
		{
			tenantID: "sh",
			expected: "storage/tenants/s/h/clinical_sh.db",
		},
		{
			tenantID: "ab",
			expected: "storage/tenants/a/b/clinical_ab.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tenantID, func(t *testing.T) {
			result := ResolveTenantPath(tt.tenantID)
			if result != tt.expected {
				t.Errorf("ResolveTenantPath(%s) = %s, want %s", tt.tenantID, result, tt.expected)
			}
		})
	}
}

func TestEnsureTenantDir(t *testing.T) {
	tenantID := "abcd1234-ef56-7890-abcd-ef1234567890"

	err := EnsureTenantDir(tenantID)
	if err != nil {
		t.Fatalf("EnsureTenantDir failed: %v", err)
	}

	expectedDir := filepath.Join(StorageTenantsDir, "ab", "cd")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Directory %s was not created", expectedDir)
	}

	os.RemoveAll(filepath.Join(StorageTenantsDir, "ab"))
}

func TestPathResolver_Integration(t *testing.T) {
	tempDir := t.TempDir()
	pr := NewPathResolver(tempDir)

	tenantID := "test-uuid-1234-5678-abcd-ef1234567890"

	t.Run("ResolveAndEnsure", func(t *testing.T) {
		dbPath := pr.ResolveTenantPath(tenantID)
		expectedPath := "storage/tenants/te/st/clinical_" + tenantID + ".db"

		if dbPath != expectedPath {
			t.Errorf("ResolveTenantPath = %s, want %s", dbPath, expectedPath)
		}

		err := pr.EnsureTenantDir(tenantID)
		if err != nil {
			t.Fatalf("EnsureTenantDir failed: %v", err)
		}

		if !pr.TenantDirExists(tenantID) {
			t.Error("TenantDirExists returned false after EnsureTenantDir")
		}

		if pr.TenantDBExists(tenantID) {
			t.Error("TenantDBExists returned true before DB was created")
		}
	})
}
