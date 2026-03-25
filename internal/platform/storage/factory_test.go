package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetTenantDBPath(t *testing.T) {
	tests := []struct {
		name           string
		tenantID       string
		envAPP_ENV     string
		expectedPrefix string
	}{
		{
			name:           "production - full UUID",
			tenantID:       "626eac12-0d34-4794-a100-b379afa0a1fc",
			envAPP_ENV:     "",
			expectedPrefix: "storage/tenants/62/6e/clinical_626eac12-0d34-4794-a100-b379afa0a1fc.db",
		},
		{
			name:           "production - different UUID",
			tenantID:       "abc12345-6789-0abc-def0-123456789abc",
			envAPP_ENV:     "",
			expectedPrefix: "storage/tenants/ab/c1/clinical_abc12345-6789-0abc-def0-123456789abc.db",
		},
		{
			name:           "test environment - should use in-memory",
			tenantID:       "test-1234-5678",
			envAPP_ENV:     "test",
			expectedPrefix: "file:clinical_test-1234-5678",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envAPP_ENV != "" {
				os.Setenv("APP_ENV", tt.envAPP_ENV)
				defer os.Unsetenv("APP_ENV")
			} else {
				os.Unsetenv("APP_ENV")
			}

			result := GetTenantDBPath(tt.tenantID)

			if tt.envAPP_ENV == "test" {
				expectedInMemory := "file:clinical_test-1234-5678?mode=memory&cache=shared"
				if result != expectedInMemory {
					t.Errorf("GetTenantDBPath() = %s, want %s", result, expectedInMemory)
				}
			} else {
				if result == "" {
					t.Error("GetTenantDBPath returned empty string for production")
				}
				if result != tt.expectedPrefix {
					t.Errorf("GetTenantDBPath() = %s, want %s", result, tt.expectedPrefix)
				}
			}
		})
	}
}

func TestGetTenantDBPath_Consistency(t *testing.T) {
	os.Unsetenv("APP_ENV")

	tenantID := "626eac12-0d34-4794-a100-b379afa0a1fc"

	path1 := GetTenantDBPath(tenantID)
	path2 := GetTenantDBPath(tenantID)

	if path1 != path2 {
		t.Errorf("GetTenantDBPath is not consistent: first=%s, second=%s", path1, path2)
	}

	expectedPath := filepath.Join("storage", "tenants", "62", "6e", "clinical_626eac12-0d34-4794-a100-b379afa0a1fc.db")
	if path1 != expectedPath {
		t.Errorf("GetTenantDBPath() = %s, want %s", path1, expectedPath)
	}
}

func TestIsTestEnvironment(t *testing.T) {
	tests := []struct {
		envValue    string
		expected    bool
		description string
	}{
		{"test", true, "APP_ENV=test should return true"},
		{"", false, "empty APP_ENV should return false"},
		{"production", false, "production should return false"},
		{"dev", false, "dev should return false"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if tt.envValue == "" {
				os.Unsetenv("APP_ENV")
			} else {
				os.Setenv("APP_ENV", tt.envValue)
			}
			defer os.Unsetenv("APP_ENV")

			result := IsTestEnvironment()
			if result != tt.expected {
				t.Errorf("IsTestEnvironment() = %v, want %v for APP_ENV=%s", result, tt.expected, tt.envValue)
			}
		})
	}
}
