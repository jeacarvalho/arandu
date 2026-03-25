package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	CentralDBFilename = "arandu_central.db"
	TenantPrefix      = "clinical_"
)

func IsTestEnvironment() bool {
	return os.Getenv("APP_ENV") == "test"
}

func GetCentralDBPath() string {
	return filepath.Join("storage", CentralDBFilename)
}

func GetTenantDBPath(tenantID string) string {
	if os.Getenv("APP_ENV") == "test" {
		return fmt.Sprintf("file:clinical_%s?mode=memory&cache=shared", tenantID)
	}

	prefix := tenantID[:2]
	suffix := tenantID[2:4]
	dir := filepath.Join("storage", "tenants", prefix, suffix)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return ""
	}

	return filepath.Join(dir, fmt.Sprintf("clinical_%s.db", tenantID))
}

func HashTenantID(tenantID string) string {
	return tenantID
}
