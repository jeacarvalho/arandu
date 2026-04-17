package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/domain/session"
)

// mockPatientServiceForUpdate implementa PatientServiceInterface minimal
type mockPatientServiceForUpdate struct {
	p *patient.Patient
}

func (m *mockPatientServiceForUpdate) GetPatientByID(_ context.Context, _ string) (*patient.Patient, error) {
	if m.p != nil {
		return m.p, nil
	}
	return &patient.Patient{ID: "pat-1", Name: "Ana Lima"}, nil
}
func (m *mockPatientServiceForUpdate) ListPatients(_ context.Context) ([]*patient.Patient, error) {
	return nil, nil
}

// mockSessionServiceForUpdate implementa SessionServiceInterface
type mockSessionServiceForUpdate struct {
	sessions  map[string]*session.Session
	updated   map[string]services.UpdateSessionInput
	updateErr error
}

func (m *mockSessionServiceForUpdate) GetSession(_ context.Context, id string) (*session.Session, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, nil
}
func (m *mockSessionServiceForUpdate) ListSessionsByPatient(_ context.Context, _ string) ([]*session.Session, error) {
	return nil, nil
}
func (m *mockSessionServiceForUpdate) CreateSession(_ context.Context, _ string, _ time.Time, _ string) (*session.Session, error) {
	return nil, nil
}
func (m *mockSessionServiceForUpdate) UpdateSession(_ context.Context, input services.UpdateSessionInput) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if m.updated == nil {
		m.updated = make(map[string]services.UpdateSessionInput)
	}
	m.updated[input.ID] = input
	return nil
}

func newUpdateSessionHandler(svc SessionServiceInterface) *SessionHandler {
	return &SessionHandler{
		sessionService: svc,
		patientService: &mockPatientServiceForUpdate{},
	}
}

// TestUpdateSession_ExtractsIDFromURLPath verifica que o handler extrai o session ID
// da URL path (/session/{id}/update) e não do form body.
func TestUpdateSession_ExtractsIDFromURLPath(t *testing.T) {
	sessDate := time.Date(2025, 11, 8, 0, 0, 0, 0, time.UTC)
	svc := &mockSessionServiceForUpdate{
		sessions: map[string]*session.Session{
			"s0353-124": {ID: "s0353-124", PatientID: "pat-1", Date: sessDate},
		},
	}
	handler := newUpdateSessionHandler(svc)

	form := url.Values{}
	form.Set("date", "2025-11-08")
	form.Set("summary", "Sessão atualizada")
	// INTENCIONALMENTE sem session_id no form body

	req := httptest.NewRequest(http.MethodPost, "/session/s0353-124/update", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.Path = "/session/s0353-124/update"
	w := httptest.NewRecorder()

	handler.UpdateSession(w, req)

	// Não deve retornar 400 "ID da sessão é obrigatório"
	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "ID da sessão") {
		t.Errorf("UpdateSession() não deve retornar 400 quando session_id está na URL path. Body: %s", w.Body.String())
	}

	// Deve ter chamado UpdateSession no serviço com o ID correto
	if _, ok := svc.updated["s0353-124"]; !ok {
		t.Errorf("UpdateSession() deve chamar serviço com ID 's0353-124' extraído da URL path. updated: %v, code: %d, body: %s",
			svc.updated, w.Code, w.Body.String())
	}
}

// TestUpdateSession_RedirectsToDetailAfterSave verifica que após salvar com sucesso
// o handler redireciona para a página de detalhe (/session/{id}), não de volta para /edit.
func TestUpdateSession_RedirectsToDetailAfterSave(t *testing.T) {
	sessDate := time.Date(2025, 11, 8, 0, 0, 0, 0, time.UTC)
	svc := &mockSessionServiceForUpdate{
		sessions: map[string]*session.Session{
			"s0353-124": {ID: "s0353-124", PatientID: "pat-1", Date: sessDate},
		},
	}
	handler := newUpdateSessionHandler(svc)

	form := url.Values{}
	form.Set("date", "2025-11-08")
	form.Set("summary", "Atualizado")

	req := httptest.NewRequest(http.MethodPost, "/session/s0353-124/update", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.Path = "/session/s0353-124/update"
	w := httptest.NewRecorder()

	handler.UpdateSession(w, req)

	// Deve redirecionar (3xx)
	if w.Code != http.StatusSeeOther && w.Code != http.StatusFound {
		t.Errorf("UpdateSession() deve redirecionar após salvar, got %d. Body: %s", w.Code, w.Body.String())
		return
	}

	location := w.Header().Get("Location")
	// Deve ir para o DETALHE, não para /edit
	if strings.HasSuffix(location, "/edit") {
		t.Errorf("UpdateSession() não deve redirecionar para /edit — deve ir para detalhe /session/{id}. Location: %s", location)
	}
	if !strings.Contains(location, "/session/s0353-124") {
		t.Errorf("UpdateSession() deve redirecionar para /session/s0353-124. Location: %s", location)
	}
}
