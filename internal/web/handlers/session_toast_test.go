package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"arandu/internal/domain/observation"
	interventionDomain "arandu/internal/domain/intervention"
)

// mockObservationService implementa ObservationServiceInterface para testes
type mockObservationService struct {
	obs *observation.Observation
	err error
}

func (m *mockObservationService) CreateObservation(_ context.Context, sessionID, content string) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.obs != nil {
		return m.obs, nil
	}
	return &observation.Observation{
		ID:        "obs-1",
		SessionID: sessionID,
		Content:   content,
		CreatedAt: time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
	}, nil
}

func (m *mockObservationService) GetObservationsBySession(_ context.Context, _ string) ([]interface{}, error) {
	return nil, nil
}

// mockInterventionService implementa InterventionServiceInterface para testes
type mockInterventionService struct {
	intv *interventionDomain.Intervention
	err  error
}

func (m *mockInterventionService) CreateIntervention(_ context.Context, sessionID, content string) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.intv != nil {
		return m.intv, nil
	}
	return &interventionDomain.Intervention{
		ID:        "intv-1",
		SessionID: sessionID,
		Content:   content,
		CreatedAt: time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
	}, nil
}

func (m *mockInterventionService) GetInterventionsBySession(_ context.Context, _ string) ([]interface{}, error) {
	return nil, nil
}

func newHandlerWithObsIntv(obs ObservationServiceInterface, intv InterventionServiceInterface) *SessionHandler {
	return &SessionHandler{
		observationService:  obs,
		interventionService: intv,
	}
}

// TestCreateObservation_ResponseIncludesToast verifica que a resposta inclui
// o fragmento OOB do toast após criação bem-sucedida.
func TestCreateObservation_ResponseIncludesToast(t *testing.T) {
	handler := newHandlerWithObsIntv(&mockObservationService{}, nil)

	form := url.Values{}
	form.Set("content", "Paciente apresentou melhora na regulação emocional")
	req := httptest.NewRequest(http.MethodPost, "/session/sess-1/observations", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.Path = "/session/sess-1/observations"
	w := httptest.NewRecorder()

	handler.CreateObservation(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("CreateObservation() status = %d, want 200. Body: %s", w.Code, w.Body.String())
	}
	html := w.Body.String()
	if !strings.Contains(html, "toast-container") {
		t.Errorf("CreateObservation() resposta deve conter fragmento OOB do toast. got: %s", html)
	}
	if !strings.Contains(html, "Observação adicionada") {
		t.Errorf("CreateObservation() toast deve conter 'Observação adicionada'. got: %s", html)
	}
}

// TestCreateIntervention_ResponseIncludesToast verifica que a resposta inclui
// o fragmento OOB do toast após criação bem-sucedida.
func TestCreateIntervention_ResponseIncludesToast(t *testing.T) {
	handler := newHandlerWithObsIntv(nil, &mockInterventionService{})

	form := url.Values{}
	form.Set("content", "Aplicação de TCC — reestruturação cognitiva")
	req := httptest.NewRequest(http.MethodPost, "/session/sess-1/interventions", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.Path = "/session/sess-1/interventions"
	w := httptest.NewRecorder()

	handler.CreateIntervention(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("CreateIntervention() status = %d, want 200. Body: %s", w.Code, w.Body.String())
	}
	html := w.Body.String()
	if !strings.Contains(html, "toast-container") {
		t.Errorf("CreateIntervention() resposta deve conter fragmento OOB do toast. got: %s", html)
	}
	if !strings.Contains(html, "Intervenção adicionada") {
		t.Errorf("CreateIntervention() toast deve conter 'Intervenção adicionada'. got: %s", html)
	}
}

// TestCreateObservation_NoToastOnError verifica que erro não gera toast de sucesso.
func TestCreateObservation_NoToastOnError(t *testing.T) {
	handler := newHandlerWithObsIntv(&mockObservationService{}, nil)

	// Conteúdo vazio → deve retornar 400
	form := url.Values{}
	form.Set("content", "")
	req := httptest.NewRequest(http.MethodPost, "/session/sess-1/observations", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.Path = "/session/sess-1/observations"
	w := httptest.NewRecorder()

	handler.CreateObservation(w, req)

	if w.Code == http.StatusOK {
		t.Error("CreateObservation() com conteúdo vazio não deve retornar 200")
	}
	if strings.Contains(w.Body.String(), "Observação adicionada") {
		t.Error("CreateObservation() com erro não deve conter toast de sucesso")
	}
}
