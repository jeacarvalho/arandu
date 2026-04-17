package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/session"
)

// mockSessionServiceForAutosave é um mock focado nos métodos necessários para autosave
type mockSessionServiceForAutosave struct {
	sessions map[string]*session.Session
	updated  map[string]string // id → summary updated
	getErr   error
	updateErr error
}

func (m *mockSessionServiceForAutosave) GetSession(ctx context.Context, id string) (*session.Session, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	sess, ok := m.sessions[id]
	if !ok {
		return nil, nil
	}
	return sess, nil
}

func (m *mockSessionServiceForAutosave) ListSessionsByPatient(ctx context.Context, patientID string) ([]*session.Session, error) {
	return nil, nil
}

func (m *mockSessionServiceForAutosave) CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error) {
	return nil, nil
}

func (m *mockSessionServiceForAutosave) UpdateSession(ctx context.Context, input services.UpdateSessionInput) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if sess, ok := m.sessions[input.ID]; ok {
		sess.Summary = input.Summary
	}
	if m.updated == nil {
		m.updated = make(map[string]string)
	}
	m.updated[input.ID] = input.Summary
	return nil
}

func newMinimalSessionHandler(svc SessionServiceInterface) *SessionHandler {
	return &SessionHandler{
		sessionService: svc,
	}
}

// TestSessionHandler_PatchSummary_Success verifica que um PATCH válido
// no summary retorna 200 com o indicador de autosave "Gravado".
func TestSessionHandler_PatchSummary_Success(t *testing.T) {
	sessDate := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	svc := &mockSessionServiceForAutosave{
		sessions: map[string]*session.Session{
			"sess-abc": {
				ID:        "sess-abc",
				PatientID: "patient-xyz",
				Date:      sessDate,
				Summary:   "Resumo anterior",
			},
		},
	}
	handler := newMinimalSessionHandler(svc)

	body := "summary=Novo+resumo+clínico+atualizado"
	req := httptest.NewRequest(http.MethodPatch, "/session/sess-abc/summary", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = setPathVar(req, "/session/", "sess-abc")
	w := httptest.NewRecorder()

	handler.PatchSummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("PatchSummary() status = %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	html := w.Body.String()
	if !strings.Contains(html, "Gravado") {
		t.Errorf("PatchSummary() response deve conter indicador 'Gravado', got: %s", html)
	}

	// Verifica que o serviço foi chamado com o novo summary
	if svc.updated["sess-abc"] != "Novo resumo clínico atualizado" {
		t.Errorf("PatchSummary() não atualizou o summary no serviço. got: %q", svc.updated["sess-abc"])
	}
}

// TestSessionHandler_PatchSummary_EmptySummaryAllowed verifica que autosave
// permite summary vazio — o terapeuta deve poder limpar o campo.
func TestSessionHandler_PatchSummary_EmptySummaryAllowed(t *testing.T) {
	sessDate := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	svc := &mockSessionServiceForAutosave{
		sessions: map[string]*session.Session{
			"sess-abc": {
				ID:        "sess-abc",
				PatientID: "patient-xyz",
				Date:      sessDate,
				Summary:   "Resumo anterior",
			},
		},
	}
	handler := newMinimalSessionHandler(svc)

	body := "summary="
	req := httptest.NewRequest(http.MethodPatch, "/session/sess-abc/summary", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = setPathVar(req, "/session/", "sess-abc")
	w := httptest.NewRecorder()

	handler.PatchSummary(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("PatchSummary() summary vazio deve retornar 200, got %d", w.Code)
	}
}

// TestSessionHandler_PatchSummary_SessionNotFound verifica que sessão inexistente
// retorna 404.
func TestSessionHandler_PatchSummary_SessionNotFound(t *testing.T) {
	svc := &mockSessionServiceForAutosave{
		sessions: map[string]*session.Session{},
	}
	handler := newMinimalSessionHandler(svc)

	body := "summary=qualquer+coisa"
	req := httptest.NewRequest(http.MethodPatch, "/session/nao-existe/summary", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = setPathVar(req, "/session/", "nao-existe")
	w := httptest.NewRecorder()

	handler.PatchSummary(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("PatchSummary() sessão inexistente deve retornar 404, got %d", w.Code)
	}
}

// TestSessionHandler_PatchSummary_ServiceError verifica que erro no serviço
// retorna 500 e não o indicador de "Gravado".
func TestSessionHandler_PatchSummary_ServiceError(t *testing.T) {
	sessDate := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	svc := &mockSessionServiceForAutosave{
		sessions: map[string]*session.Session{
			"sess-abc": {
				ID:      "sess-abc",
				Date:    sessDate,
				Summary: "Resumo",
			},
		},
		updateErr: context.DeadlineExceeded,
	}
	handler := newMinimalSessionHandler(svc)

	body := "summary=novo+resumo"
	req := httptest.NewRequest(http.MethodPatch, "/session/sess-abc/summary", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = setPathVar(req, "/session/", "sess-abc")
	w := httptest.NewRecorder()

	handler.PatchSummary(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("PatchSummary() com erro no serviço deve retornar 500, got %d", w.Code)
	}

	if strings.Contains(w.Body.String(), "Gravado") {
		t.Errorf("PatchSummary() com erro não deve conter 'Gravado'")
	}
}

// setPathVar simula a extração do ID do path no handler.
// Os handlers usam extractIDFromPath, então precisamos montar a URL corretamente.
func setPathVar(r *http.Request, prefix, id string) *http.Request {
	r2 := r.Clone(r.Context())
	r2.URL.Path = prefix + id + "/summary"
	return r2
}

// TestAutosaveResponseContainsHTMXFragment verifica que a resposta do autosave
// é um fragmento HTML adequado para swap HTMX (não contém tags html/body).
func TestAutosaveResponseContainsHTMXFragment(t *testing.T) {
	sessDate := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	svc := &mockSessionServiceForAutosave{
		sessions: map[string]*session.Session{
			"sess-abc": {
				ID:      "sess-abc",
				Date:    sessDate,
				Summary: "Resumo",
			},
		},
	}
	handler := newMinimalSessionHandler(svc)

	body := "summary=Resumo+atualizado"
	req := httptest.NewRequest(http.MethodPatch, "/session/sess-abc/summary", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = setPathVar(req, "/session/", "sess-abc")
	w := httptest.NewRecorder()

	handler.PatchSummary(w, req)

	html := w.Body.String()

	// Fragmento HTMX não deve ter estrutura de página completa
	if bytes.Contains([]byte(html), []byte("<html")) {
		t.Errorf("PatchSummary() resposta não deve conter <html> — deve ser fragmento HTMX")
	}
	if bytes.Contains([]byte(html), []byte("<body")) {
		t.Errorf("PatchSummary() resposta não deve conter <body> — deve ser fragmento HTMX")
	}
}
