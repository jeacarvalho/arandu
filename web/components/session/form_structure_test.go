package session

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestObservationForm_NoFormHelp verifica que observation_form não tem div.form-help.
// Texto de ajuda causa diferença de altura entre os cards quando há quebra de linha,
// quebrando o alinhamento visual dos botões "Adicionar".
func TestObservationForm_NoFormHelp(t *testing.T) {
	var buf bytes.Buffer
	if err := ObservationForm("sess-1").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(buf.String(), "form-help") {
		t.Errorf("ObservationForm não deve ter div.form-help — cria diferença de altura que desalinha botões")
	}
}

// TestInterventionForm_NoFormHelp mesma razão — form-help desalinha os botões.
func TestInterventionForm_NoFormHelp(t *testing.T) {
	var buf bytes.Buffer
	if err := InterventionForm("sess-1").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(buf.String(), "form-help") {
		t.Errorf("InterventionForm não deve ter div.form-help — cria diferença de altura que desalinha botões")
	}
}

// TestInterventionForm_HasLoadingIndicator verifica que intervention_form tem o span
// htmx-indicator, estruturalmente idêntico ao observation_form. Sem o span, as form-actions
// têm alturas diferentes e os botões ficam visualmente desalinhados.
func TestInterventionForm_HasLoadingIndicator(t *testing.T) {
	var buf bytes.Buffer
	if err := InterventionForm("sess-1").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "htmx-indicator") {
		t.Errorf("InterventionForm deve ter span.htmx-indicator — ausência cria assimetria de altura com ObservationForm")
	}
}

// TestObservationForm_HasLoadingIndicator garante que observation_form mantém o span.
func TestObservationForm_HasLoadingIndicator(t *testing.T) {
	var buf bytes.Buffer
	if err := ObservationForm("sess-1").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "htmx-indicator") {
		t.Errorf("ObservationForm deve manter span.htmx-indicator")
	}
}
