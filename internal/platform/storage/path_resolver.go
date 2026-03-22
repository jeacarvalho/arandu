package storage

import (
	"os"
	"path/filepath"
)

const (
	StorageTenantsDir = "storage/tenants"
	ClinicalDBPrefix  = "clinical_"
	ClinicalDBExt     = ".db"
	DirPerms          = 0700
)

type PathResolver struct {
	storagePath string
}

func NewPathResolver(storagePath string) *PathResolver {
	return &PathResolver{
		storagePath: storagePath,
	}
}

func ResolveTenantPath(tenantID string) string {
	if len(tenantID) < 4 {
		prefix := tenantID[:1]
		subPrefix := tenantID[1:]
		return filepath.Join(StorageTenantsDir, prefix, subPrefix, ClinicalDBPrefix+tenantID+ClinicalDBExt)
	}
	prefix := tenantID[:2]
	subPrefix := tenantID[2:4]
	return filepath.Join(StorageTenantsDir, prefix, subPrefix, ClinicalDBPrefix+tenantID+ClinicalDBExt)
}

func (pr *PathResolver) ResolveTenantPath(tenantID string) string {
	return ResolveTenantPath(tenantID)
}

func (pr *PathResolver) GetTenantDir(tenantID string) string {
	if len(tenantID) < 4 {
		return StorageTenantsDir
	}
	prefix := tenantID[:2]
	subPrefix := tenantID[2:4]
	return filepath.Join(StorageTenantsDir, prefix, subPrefix)
}

func EnsureTenantDir(tenantID string) error {
	var dir string
	if len(tenantID) >= 4 {
		dir = filepath.Join(StorageTenantsDir, tenantID[:2], tenantID[2:4])
	} else {
		dir = StorageTenantsDir
	}
	return os.MkdirAll(dir, DirPerms)
}

func (pr *PathResolver) EnsureTenantDir(tenantID string) error {
	return EnsureTenantDir(tenantID)
}

func (pr *PathResolver) TenantDirExists(tenantID string) bool {
	dir := pr.GetTenantDir(tenantID)
	_, err := os.Stat(dir)
	return err == nil
}

func (pr *PathResolver) TenantDBExists(tenantID string) bool {
	dbPath := pr.ResolveTenantPath(tenantID)
	_, err := os.Stat(dbPath)
	return err == nil
}

func (pr *PathResolver) GetTenantsDir() string {
	return filepath.Join(pr.storagePath, StorageTenantsDir)
}
