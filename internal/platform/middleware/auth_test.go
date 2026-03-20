package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"arandu/internal/infrastructure/repository/sqlite"

	_ "modernc.org/sqlite"
)

func TestIsPublicRoute(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/login", true},
		{"/logout", true},
		{"/static/css/style.css", true},
		{"/test", true},
		{"/", true},
		{"/dashboard", false},
		{"/patients", false},
		{"/patient/123", false},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("GET", tt.path, nil)
		result := isPublicRoute(r)
		if result != tt.expected {
			t.Errorf("isPublicRoute(%s) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}

func TestGetTenantID(t *testing.T) {
	ctx := context.WithValue(context.Background(), TenantIDKey, "test-tenant-id")

	tenantID, err := GetTenantID(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tenantID != "test-tenant-id" {
		t.Errorf("expected 'test-tenant-id', got '%s'", tenantID)
	}
}

func TestGetTenantID_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetTenantID(ctx)
	if err == nil {
		t.Error("expected error when tenant ID not in context")
	}
}

func TestGetTenantDB(t *testing.T) {
	db := &sql.DB{}
	ctx := context.WithValue(context.Background(), TenantDBKey, db)

	result, err := GetTenantDB(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != db {
		t.Error("expected same DB instance")
	}
}

func TestGetTenantDB_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetTenantDB(ctx)
	if err == nil {
		t.Error("expected error when DB not in context")
	}
}

func TestGetUserID(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, "user-123")

	userID, err := GetUserID(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("expected 'user-123', got '%s'", userID)
	}
}

func TestGetUserID_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := GetUserID(ctx)
	if err == nil {
		t.Error("expected error when user ID not in context")
	}
}

func TestAuthMiddleware_PublicRoute(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	pool := sqlite.NewTenantPool(t.TempDir(), nil)
	auth := NewAuthMiddleware(centralDB, pool)

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/login", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200 for public route, got %d", rec.Code)
	}
}

func TestAuthMiddleware_NoSessionCookie(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	pool := sqlite.NewTenantPool(t.TempDir(), nil)
	auth := NewAuthMiddleware(centralDB, pool)

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/dashboard", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("expected redirect to /login (302), got %d", rec.Code)
	}

	if rec.Header().Get("Location") != "/login" {
		t.Errorf("expected redirect to /login, got '%s'", rec.Header().Get("Location"))
	}
}

func TestAuthMiddleware_ExpiredSession(t *testing.T) {
	tmpDir := t.TempDir()
	centralDB, err := sqlite.NewCentralDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create central DB: %v", err)
	}
	defer centralDB.Close()

	if err := centralDB.Migrate(nil); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	tenantID := "tenant-test"
	dbPath := tmpDir + "/tenants/clinical_" + tenantID + ".db"
	_, err = sql.Open("sqlite", dbPath)
	if err != nil {
		t.Logf("creating tenant db: %v", err)
	}
	centralDB.Exec("INSERT INTO tenants (id, db_path, status) VALUES (?, ?, 'active')", tenantID, dbPath)

	userID := "user-test"
	centralDB.Exec("INSERT INTO users (id, email, password_hash, tenant_id) VALUES (?, ?, '', ?)", userID, "test@test.com", tenantID)

	sessionID := "expired-session"
	expiredTime := time.Now().Add(-1 * time.Hour).Unix()
	centralDB.Exec("INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES (?, ?, ?, ?)", sessionID, userID, tenantID, expiredTime)

	pool := sqlite.NewTenantPool(tmpDir, nil)
	auth := NewAuthMiddleware(centralDB, pool)

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: sessionID})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("expected redirect for expired session, got %d", rec.Code)
	}
}
