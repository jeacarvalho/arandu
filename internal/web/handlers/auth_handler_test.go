package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"arandu/internal/web/handlers"
)

// authHandlerForTest creates an AuthHandler without a real DB by using nil.
// This lets us test render-only paths (GET /login) without requiring SQLite.
func newTestAuthHandler(t *testing.T) *handlers.AuthHandler {
	t.Helper()
	return handlers.NewAuthHandler(nil)
}

func TestLoginPage_RendersWithoutDevButton_WhenNotDev(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	h := newTestAuthHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if strings.Contains(body, "Login Rápido") {
		t.Error("dev login button must NOT appear when APP_ENV=production")
	}
}

func TestLoginPage_RendersWithDevButton_WhenDev(t *testing.T) {
	t.Setenv("APP_ENV", "dev")

	h := newTestAuthHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Login Rápido") {
		t.Error("dev login button must appear when APP_ENV=dev")
	}
}

func TestLoginPage_RendersWithDevButton_WhenDevelopment(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	h := newTestAuthHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "Login Rápido") {
		t.Error("dev login button must appear when APP_ENV=development")
	}
}

func TestLoginPage_HTMXRequest_Redirects(t *testing.T) {
	h := newTestAuthHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	req.Header.Set("HX-Request", "true")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("HTMX login request must redirect, got %d", w.Code)
	}
	if w.Header().Get("HX-Redirect") != "/login" {
		t.Errorf("expected HX-Redirect=/login, got %q", w.Header().Get("HX-Redirect"))
	}
}

func TestLoginPage_UsesLocalHTMX(t *testing.T) {
	h := newTestAuthHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	body := w.Body.String()
	if strings.Contains(body, "unpkg.com") {
		t.Error("login page must not load HTMX from CDN (unpkg.com); use /static/js/htmx.min.js")
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("/static/js/htmx.min.js")) {
		t.Error("login page must load HTMX from /static/js/htmx.min.js")
	}
}

func TestLoginPage_MethodNotAllowed(t *testing.T) {
	h := newTestAuthHandler(t)
	req := httptest.NewRequest(http.MethodPut, "/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for PUT /login, got %d", w.Code)
	}
}
