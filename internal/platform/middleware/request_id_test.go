package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	appcontext "arandu/internal/platform/context"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica que o request_id foi adicionado ao contexto
		requestID, err := appcontext.GetRequestID(r.Context())
		if err != nil {
			t.Errorf("Expected request_id in context, got error: %v", err)
		}
		if requestID == "" {
			t.Error("Expected non-empty request_id")
		}

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware := RequestIDMiddleware(handler)
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Verifica que o X-Request-ID foi adicionado ao header
	headerRequestID := rr.Header().Get("X-Request-ID")
	if headerRequestID == "" {
		t.Error("Expected X-Request-ID header in response")
	}
}

func TestRequestIDIsUnique(t *testing.T) {
	requestIDs := make(map[string]bool)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, _ := appcontext.GetRequestID(r.Context())
		if requestIDs[requestID] {
			t.Error("Request ID should be unique")
		}
		requestIDs[requestID] = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestIDMiddleware(handler)

	// Faz 100 requisições e verifica que todas têm IDs únicos
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
	}

	if len(requestIDs) != 100 {
		t.Errorf("Expected 100 unique request IDs, got %d", len(requestIDs))
	}
}

func TestRequestIDHeaderPropagation(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware := RequestIDMiddleware(handler)
	middleware.ServeHTTP(rr, req)

	// Verifica que o header X-Request-ID foi setado
	headerRequestID := rr.Header().Get("X-Request-ID")
	if headerRequestID == "" {
		t.Error("Expected X-Request-ID header in response")
	}

	// Verifica que o header tem 32 caracteres hex (16 bytes)
	if len(headerRequestID) != 32 {
		t.Errorf("Expected request_id to be 32 hex chars, got %d", len(headerRequestID))
	}
}
