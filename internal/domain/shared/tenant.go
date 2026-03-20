package shared

import (
	"errors"
	"time"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
)

type Tenant struct {
	ID        string       `json:"id"`
	DBPath    string       `json:"db_path"`
	Status    TenantStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

func NewTenant(id, dbPath string) *Tenant {
	return &Tenant{
		ID:        id,
		DBPath:    dbPath,
		Status:    TenantStatusActive,
		CreatedAt: time.Now(),
	}
}

func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

func (t *Tenant) Suspend() error {
	if t.Status == TenantStatusSuspended {
		return errors.New("tenant already suspended")
	}
	t.Status = TenantStatusSuspended
	return nil
}

func (t *Tenant) Activate() error {
	if t.Status == TenantStatusActive {
		return errors.New("tenant already active")
	}
	t.Status = TenantStatusActive
	return nil
}
