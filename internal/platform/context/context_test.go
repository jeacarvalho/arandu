package context

import (
	"context"
	"database/sql"
	"testing"
)

func TestWithTenantDB(t *testing.T) {
	db := &sql.DB{}
	ctx := context.Background()

	result := WithTenantDB(ctx, db)
	if result == ctx {
		t.Error("WithTenantDB should return a new context")
	}

	retrieved, err := GetTenantDB(result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved != db {
		t.Error("expected same DB instance")
	}
}

func TestWithTenantDB_Nil(t *testing.T) {
	ctx := context.Background()

	result := WithTenantDB(ctx, nil)
	retrieved, err := GetTenantDB(result)

	if err != ErrNoTenantDB {
		t.Errorf("expected ErrNoTenantDB, got %v", err)
	}
	if retrieved != nil {
		t.Error("expected nil DB")
	}
}

func TestGetTenantDB_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetTenantDB(ctx)
	if err != ErrNoTenantDB {
		t.Errorf("expected ErrNoTenantDB, got %v", err)
	}
}

func TestWithTenantID(t *testing.T) {
	ctx := context.Background()
	tenantID := "test-tenant-123"

	result := WithTenantID(ctx, tenantID)
	if result == ctx {
		t.Error("WithTenantID should return a new context")
	}

	retrieved, err := GetTenantID(result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved != tenantID {
		t.Errorf("expected '%s', got '%s'", tenantID, retrieved)
	}
}

func TestWithTenantID_Empty(t *testing.T) {
	ctx := context.Background()

	result := WithTenantID(ctx, "")
	retrieved, err := GetTenantID(result)

	if err != ErrNoTenantID {
		t.Errorf("expected ErrNoTenantID, got %v", err)
	}
	if retrieved != "" {
		t.Error("expected empty string")
	}
}

func TestGetTenantID_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetTenantID(ctx)
	if err != ErrNoTenantID {
		t.Errorf("expected ErrNoTenantID, got %v", err)
	}
}

func TestWithUserID(t *testing.T) {
	ctx := context.Background()
	userID := "user-456"

	result := WithUserID(ctx, userID)
	if result == ctx {
		t.Error("WithUserID should return a new context")
	}

	retrieved, err := GetUserID(result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved != userID {
		t.Errorf("expected '%s', got '%s'", userID, retrieved)
	}
}

func TestWithUserID_Empty(t *testing.T) {
	ctx := context.Background()

	result := WithUserID(ctx, "")
	retrieved, err := GetUserID(result)

	if err != ErrNoUserID {
		t.Errorf("expected ErrNoUserID, got %v", err)
	}
	if retrieved != "" {
		t.Error("expected empty string")
	}
}

func TestGetUserID_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetUserID(ctx)
	if err != ErrNoUserID {
		t.Errorf("expected ErrNoUserID, got %v", err)
	}
}

func TestContextChaining(t *testing.T) {
	ctx := context.Background()
	tenantID := "tenant-chain"
	userID := "user-chain"
	db := &sql.DB{}

	ctx = WithTenantID(ctx, tenantID)
	ctx = WithUserID(ctx, userID)
	ctx = WithTenantDB(ctx, db)

	tid, err := GetTenantID(ctx)
	if err != nil || tid != tenantID {
		t.Errorf("tenant ID mismatch: err=%v, got=%s", err, tid)
	}

	uid, err := GetUserID(ctx)
	if err != nil || uid != userID {
		t.Errorf("user ID mismatch: err=%v, got=%s", err, uid)
	}

	retrievedDB, err := GetTenantDB(ctx)
	if err != nil || retrievedDB != db {
		t.Errorf("DB mismatch: err=%v", err)
	}
}
