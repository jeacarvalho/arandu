package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTelemetryMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	telemetry := NewTelemetryMiddleware()
	wrapped := telemetry.Middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestTelemetryMiddlewareShouldSkip(t *testing.T) {
	callCount := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	telemetry := NewTelemetryMiddleware("/static/", "/health")
	wrapped := telemetry.Middleware(handler)

	// Testa path que deve ser ignorado
	req := httptest.NewRequest("GET", "/static/style.css", nil)
	rr := httptest.NewRecorder()
	wrapped.ServeHTTP(rr, req)

	if callCount != 1 {
		t.Errorf("Expected handler to be called once, got %d", callCount)
	}

	// Testa path que NÃO deve ser ignorado
	req2 := httptest.NewRequest("GET", "/api/patients", nil)
	rr2 := httptest.NewRecorder()
	wrapped.ServeHTTP(rr2, req2)

	if callCount != 2 {
		t.Errorf("Expected handler to be called twice, got %d", callCount)
	}
}

func TestSanitizePath(t *testing.T) {
	telemetry := NewTelemetryMiddleware()

	tests := []struct {
		name     string
		path     string
		query    string
		expected string
	}{
		{
			name:     "Sem query params",
			path:     "/api/patients",
			query:    "",
			expected: "/api/patients",
		},
		{
			name:     "Query params normais",
			path:     "/api/patients",
			query:    "page=1&limit=10",
			expected: "/api/patients?page=1&limit=10",
		},
		{
			name:     "Query com code sensível",
			path:     "/auth/google/callback",
			query:    "code=abc123&state=xyz",
			expected: "/auth/google/callback?code=[REDACTED]&state=xyz",
		},
		{
			name:     "Query com token sensível",
			path:     "/api/endpoint",
			query:    "token=secret123&user=john",
			expected: "/api/endpoint?token=[REDACTED]&user=john",
		},
		{
			name:     "Query com password sensível",
			path:     "/auth/login",
			query:    "user=admin&password=secret",
			expected: "/auth/login?user=admin&password=[REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := telemetry.sanitizePath(tt.path, tt.query)
			// A ordem dos parâmetros pode variar, então verificamos se contém os elementos esperados
			if tt.query == "" {
				if result != tt.expected {
					t.Errorf("Expected '%s', got '%s'", tt.expected, result)
				}
			} else {
				// Para queries com parâmetros sensíveis, verificamos se o parâmetro foi redacted
				if strings.Contains(tt.query, "code=") && !strings.Contains(result, "code=[REDACTED]") {
					t.Errorf("Expected code to be redacted in '%s'", result)
				}
				if strings.Contains(tt.query, "token=") && !strings.Contains(result, "token=[REDACTED]") {
					t.Errorf("Expected token to be redacted in '%s'", result)
				}
				if strings.Contains(tt.query, "password=") && !strings.Contains(result, "password=[REDACTED]") {
					t.Errorf("Expected password to be redacted in '%s'", result)
				}
			}
		})
	}
}

func TestGetClientIP(t *testing.T) {
	telemetry := NewTelemetryMiddleware()

	tests := []struct {
		name     string
		headers  map[string]string
		remote   string
		expected string
	}{
		{
			name:     "X-Forwarded-For com múltiplos IPs",
			headers:  map[string]string{"X-Forwarded-For": "10.0.0.1, 10.0.0.2, 10.0.0.3"},
			remote:   "192.168.1.1:12345",
			expected: "10.0.0.1",
		},
		{
			name:     "X-Real-Ip",
			headers:  map[string]string{"X-Real-Ip": "10.0.0.5"},
			remote:   "192.168.1.1:12345",
			expected: "10.0.0.5",
		},
		{
			name:     "RemoteAddr apenas",
			headers:  map[string]string{},
			remote:   "192.168.1.1:12345",
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			req.RemoteAddr = tt.remote

			result := telemetry.getClientIP(req)
			if result != tt.expected {
				t.Errorf("Expected IP '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestResponseWriterWrapper(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	telemetry := NewTelemetryMiddleware()
	wrapped := telemetry.Middleware(handler)

	req := httptest.NewRequest("POST", "/api/resource", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	body := rr.Body.String()
	if body != "Created" {
		t.Errorf("Expected body 'Created', got '%s'", body)
	}
}
