package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestNewPatientHandler tests the NewPatient handler
func TestNewPatientHandler(t *testing.T) {
	// Create a minimal handler just for testing NewPatient
	h := &Handler{}

	t.Run("GET_returns_correct_form", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/patients/new", nil)
		rr := httptest.NewRecorder()

		h.NewPatient(rr, req)

		// Check status
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		// Check content type
		contentType := rr.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Expected Content-Type text/html, got %s", contentType)
		}

		body := rr.Body.String()

		// CRITICAL: Must contain "Novo Paciente"
		if !strings.Contains(body, "Novo Paciente") {
			t.Error("FAIL: Response must contain 'Novo Paciente'")
			t.Logf("Response preview (first 200 chars):\n%s", safeSubstring(body, 200))
		} else {
			t.Log("PASS: Contains 'Novo Paciente'")
		}

		// CRITICAL: Must NOT contain "Nova Sessão"
		if strings.Contains(body, "Nova Sessão") {
			t.Error("FAIL: Response must NOT contain 'Nova Sessão' - template conflict!")
		} else {
			t.Log("PASS: Does not contain 'Nova Sessão'")
		}

		// CRITICAL: Must have correct form action
		if !strings.Contains(body, `action="/patients"`) {
			t.Error("FAIL: Form must have action '/patients'")
		} else {
			t.Log("PASS: Form action is '/patients'")
		}

		// CRITICAL: Must have patient name field
		if !strings.Contains(body, `name="name"`) {
			t.Error("FAIL: Form must have field 'name'")
		} else {
			t.Log("PASS: Form has 'name' field")
		}

		// CRITICAL: Must NOT have patient_id field (from session form)
		if strings.Contains(body, `name="patient_id"`) {
			t.Error("FAIL: Form must NOT have 'patient_id' field (wrong template)")
		} else {
			t.Log("PASS: Form does not have 'patient_id' field")
		}
	})

	t.Run("POST_method_not_allowed", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/patients/new", nil)
		rr := httptest.NewRecorder()

		h.NewPatient(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("FAIL: Expected status 405 for POST, got %d", rr.Code)
		} else {
			t.Log("PASS: POST method returns 405")
		}
	})
}

// TestTemplateConflictDetection tests that we can detect template conflicts
func TestTemplateConflictDetection(t *testing.T) {
	testCases := []struct {
		name        string
		html        string
		shouldPass  bool
		description string
	}{
		{
			name:        "correct_template",
			html:        `<h1>Novo Paciente</h1><form action="/patients"><input name="name"></form>`,
			shouldPass:  true,
			description: "Correct new patient template",
		},
		{
			name:        "wrong_template_session",
			html:        `<h1>Nova Sessão</h1><form action="/sessions"><input name="patient_id"></form>`,
			shouldPass:  false,
			description: "Wrong template (session form)",
		},
		{
			name:        "mixed_wrong_action",
			html:        `<h1>Novo Paciente</h1><form action="/sessions"><input name="name"></form>`,
			shouldPass:  false,
			description: "Correct title but wrong form action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hasNovoPaciente := strings.Contains(tc.html, "Novo Paciente")
			hasNovaSessao := strings.Contains(tc.html, "Nova Sessão")
			hasCorrectAction := strings.Contains(tc.html, `action="/patients"`)
			hasWrongAction := strings.Contains(tc.html, `action="/sessions"`)

			// A correct response has Novo Paciente, correct action, and no session content
			isCorrect := hasNovoPaciente && !hasNovaSessao && hasCorrectAction && !hasWrongAction

			if isCorrect != tc.shouldPass {
				t.Errorf("FAIL: %s\n  Expected correct=%v, got correct=%v", tc.description, tc.shouldPass, isCorrect)
				t.Logf("  Novo Paciente: %v, Nova Sessão: %v, Action /patients: %v, Action /sessions: %v",
					hasNovoPaciente, hasNovaSessao, hasCorrectAction, hasWrongAction)
			} else {
				t.Logf("PASS: %s", tc.description)
			}
		})
	}
}

// safeSubstring returns a substring without panicking
func safeSubstring(s string, n int) string {
	if n > len(s) {
		return s
	}
	return s[:n]
}
