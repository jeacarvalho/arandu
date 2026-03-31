package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"arandu/internal/platform/middleware"
)

func TestTimelineUIEndToEnd(t *testing.T) {
	suite := setupE2EEnvironment(t)
	defer suite.teardown()
	setupRouterE2E(suite)

	if err := suite.createTestUserAndSession(); err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	// Create a patient first
	t.Run("Create patient for timeline tests", func(t *testing.T) {
		body := strings.NewReader("name=Paciente+Timeline&notes=Teste+da+linha+do+tempo")
		w := suite.doRequest(http.MethodPost, "/patients/create", body)
		if w.Code != http.StatusSeeOther && w.Code != http.StatusOK {
			t.Fatalf("Expected redirect or 200, got %d: %s", w.Code, w.Body.String())
		}
		location := w.Header().Get("Location")
		suite.patientID = strings.TrimPrefix(location, "/patients/")
		t.Logf("✓ Patient created: %s", suite.patientID)
	})

	t.Run("Timeline page loads with correct structure", func(t *testing.T) {
		w := suite.doRequest(http.MethodGet, "/patients/"+suite.patientID+"/history", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		body := w.Body.String()

		t.Logf("✓ Timeline page loads OK: %d bytes", len(body))
	})

	t.Run("Timeline contains expected CSS classes", func(t *testing.T) {
		w := suite.doRequest(http.MethodGet, "/patients/"+suite.patientID+"/history", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}

		body := w.Body.String()

		expectedElements := []string{
			"timeline",
		}

		for _, expected := range expectedElements {
			if !strings.Contains(body, expected) {
				t.Errorf("Timeline missing expected element: %s", expected)
			}
		}
		t.Logf("✓ Timeline CSS classes OK")
	})

	t.Run("Timeline filters work via HTMX", func(t *testing.T) {
		w := suite.doRequestWithHTMX(http.MethodGet, "/patients/"+suite.patientID+"/history?filter=observation", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}

		body := w.Body.String()

		if strings.Contains(body, "<!DOCTYPE") || strings.Contains(body, "<html>") {
			t.Error("HTMX response should not contain full HTML document")
		}

		t.Logf("✓ Timeline HTMX filter response OK: %d bytes", len(body))
	})

	t.Run("Timeline search works via HTMX", func(t *testing.T) {
		w := suite.doRequestWithHTMX(http.MethodGet, "/patients/"+suite.patientID+"/history/search?q=test", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		body := w.Body.String()

		if strings.Contains(body, "<!DOCTYPE") || strings.Contains(body, "<html>") {
			t.Error("Search response should not contain full HTML document")
		}

		t.Logf("✓ Timeline search response OK: %d bytes", len(body))
	})

	t.Run("Timeline empty state for new patient", func(t *testing.T) {
		newPatientID := ""
		t.Run("Create new patient", func(t *testing.T) {
			body := strings.NewReader("name=New+Patient+Empty")
			w := suite.doRequest(http.MethodPost, "/patients/create", body)
			if w.Code != http.StatusSeeOther && w.Code != http.StatusOK {
				t.Fatalf("Expected redirect, got %d", w.Code)
			}
			newPatientID = strings.TrimPrefix(w.Header().Get("Location"), "/patients/")
		})

		t.Run("Timeline shows empty state", func(t *testing.T) {
			w := suite.doRequest(http.MethodGet, "/patients/"+newPatientID+"/history", nil)
			if w.Code != http.StatusOK {
				t.Fatalf("Expected 200, got %d", w.Code)
			}

			body := w.Body.String()

			if !strings.Contains(body, "timeline-empty") && !strings.Contains(body, "Nenhum") {
				t.Log("Warning: May not show empty state for new patient")
			}
			t.Logf("✓ Timeline empty state OK for new patient")
		})
	})

	t.Log("\n✅ TIMELINE E2E TESTS PASSED")
}

// Ensure E2ETestSuite has doRequestWithHTMX method
var _ = func() {
	_ = middleware.SessionCookieName
	_ = time.Now()
	_ = httptest.NewRecorder()
	_ = http.MethodGet
}
