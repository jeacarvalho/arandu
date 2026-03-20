package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"arandu/internal/infrastructure/repository/sqlite"
)

func TestAuthHandler_LoginGET(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Skip("Cannot create central DB for test")
	}
	defer centralDB.Close()

	authHandler := NewAuthHandler(centralDB)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rec := httptest.NewRecorder()

	authHandler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected 'text/html; charset=utf-8', got '%s'", contentType)
	}

	body := rec.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty body")
	}
}

func TestAuthHandler_LoginPOST(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Skip("Cannot create central DB for test")
	}
	defer centralDB.Close()

	authHandler := NewAuthHandler(centralDB)

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()

	authHandler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAuthHandler_LoginPOST_EmptyCredentials(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Skip("Cannot create central DB for test")
	}
	defer centralDB.Close()

	authHandler := NewAuthHandler(centralDB)

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.PostForm = map[string][]string{
		"email":    {""},
		"password": {""},
	}
	rec := httptest.NewRecorder()

	authHandler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if len(body) == 0 {
		t.Error("expected non-empty body with error message")
	}
}

func TestAuthHandler_InvalidMethod(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Skip("Cannot create central DB for test")
	}
	defer centralDB.Close()

	authHandler := NewAuthHandler(centralDB)

	req := httptest.NewRequest(http.MethodPut, "/login", nil)
	rec := httptest.NewRecorder()

	authHandler.Login(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}

func TestAuthHandler_GoogleLogin_NoConfig(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Skip("Cannot create central DB for test")
	}
	defer centralDB.Close()

	authHandler := NewAuthHandler(centralDB)

	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	rec := httptest.NewRecorder()

	authHandler.GoogleLogin(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("expected redirect, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/login?error=google_not_configured" {
		t.Errorf("expected redirect to /login?error=google_not_configured, got '%s'", location)
	}
}

func TestAuthHandler_ServeHTTP(t *testing.T) {
	centralDB, err := sqlite.NewCentralDB(t.TempDir())
	if err != nil {
		t.Skip("Cannot create central DB for test")
	}
	defer centralDB.Close()

	authHandler := NewAuthHandler(centralDB)

	tests := []struct {
		path       string
		method     string
		wantStatus int
	}{
		{"/login", http.MethodGet, http.StatusOK},
		{"/auth/google", http.MethodGet, http.StatusFound},
		{"/logout", http.MethodGet, http.StatusFound},
		{"/invalid", http.MethodGet, http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rec := httptest.NewRecorder()

		authHandler.ServeHTTP(rec, req)

		if rec.Code != tt.wantStatus {
			t.Errorf("path=%s method=%s: got %d, want %d", tt.path, tt.method, rec.Code, tt.wantStatus)
		}
	}
}
