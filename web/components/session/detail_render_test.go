package session

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestSessionDetailView_NovaObservacaoButtonTargetsForm(t *testing.T) {
	detail := SessionDetail{
		ID:        "sess-1",
		PatientID: "pat-1",
		PatientName: "João Silva",
		Date:      "01/04/2026",
		Summary:   "",
		CreatedAt: "2026-04-01",
	}

	var buf bytes.Buffer
	err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	html := buf.String()

	// O botão "Nova Observação" deve ter onclick que foca o textarea do form
	if !strings.Contains(html, "observation-content") {
		t.Errorf("Botão 'Nova Observação' deve referenciar #observation-content, got: %s", html)
	}
}

// TestSessionDetailView_ObservationFormPresent verifica que o formulário de
// nova observação está presente na página (dentro do card, não no header).
// O botão flutuante do header foi removido — o form inline é suficiente.
func TestSessionDetailView_ObservationFormPresent(t *testing.T) {
	detail := SessionDetail{
		ID:        "sess-1",
		PatientID: "pat-1",
		Date:      "01/04/2026",
	}

	var buf bytes.Buffer
	err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	html := buf.String()
	// O form inline de observação deve estar presente
	if !strings.Contains(html, "observation-content") {
		t.Errorf("Formulário inline de observação deve estar presente na página")
	}
	// O form de intervenção também deve estar presente
	if !strings.Contains(html, "intervention-content") {
		t.Errorf("Formulário inline de intervenção deve estar presente na página")
	}
}

func TestSessionDetailView_ErrorState_ShowsErrorMsg(t *testing.T) {
	var buf bytes.Buffer
	err := SessionDetailView(SessionDetail{}, nil, nil, "pat-1", "Sessão não encontrada").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(buf.String(), "Sessão não encontrada") {
		t.Errorf("Estado de erro deve exibir a mensagem de erro")
	}
}

// ===========================================================================
// P2 — Visual Consistency Tests
// ===========================================================================

// TestSessionDetailView_UsesShellHeaderPattern verifica que detail usa o padrão
// Shell (back-button + page-title) em vez do padrão antigo (content-header).
func TestSessionDetailView_UsesShellHeaderPattern(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-1", Date: "01/04/2026"}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()

	if strings.Contains(html, "content-header") {
		t.Error("detail.templ não deve usar 'content-header' — usar padrão Shell atual")
	}
	if !strings.Contains(html, "page-title") {
		t.Error("detail.templ deve usar 'page-title' — padrão Shell atual")
	}
	if !strings.Contains(html, "back-button") {
		t.Error("detail.templ deve ter back-button para voltar ao perfil do paciente")
	}
}

// TestSessionDetailView_CardHeadersUseShellPattern verifica que os cards usam
// card-icon + section-title em vez de card-header + card-title.
func TestSessionDetailView_CardHeadersUseShellPattern(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-1", Date: "01/04/2026"}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()

	if strings.Contains(html, `class="card-header"`) {
		t.Error("detail.templ não deve usar 'card-header' — usar flex flex-center gap-md")
	}
	if strings.Contains(html, `class="card-title"`) {
		t.Error("detail.templ não deve usar 'card-title' — usar 'section-title'")
	}
	if !strings.Contains(html, "section-title") {
		t.Error("detail.templ deve usar 'section-title' nos headers dos cards")
	}
	if !strings.Contains(html, "card-icon") {
		t.Error("detail.templ deve usar 'card-icon' nos headers dos cards")
	}
}

// TestSessionDetailView_BackButtonLinksToPatient verifica que o back-button
// aponta para o perfil do paciente correto.
func TestSessionDetailView_BackButtonLinksToPatient(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-42", Date: "01/04/2026"}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-42", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, "/patients/pat-42") {
		t.Error("back-button deve linkar para /patients/{patientID}")
	}
}

func TestSessionDetailView_EmptySummary_ShowsPlaceholder(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-1", Date: "01/04/2026"}
	var buf bytes.Buffer
	err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(buf.String(), "Nenhum resumo registrado") {
		t.Errorf("Summary vazio deve mostrar placeholder")
	}
}
