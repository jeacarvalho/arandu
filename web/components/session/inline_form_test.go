package session

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestObservationFormInline_NoRequired verifica que o form inline NÃO tem required
// no textarea — ele é sempre aninhado dentro do form principal de edição e o required
// do form interno vaza para o form externo via HTML inválido (nested forms), bloqueando
// o "Salvar Alterações" quando o campo de observação está vazio.
func TestObservationFormInline_NoRequired(t *testing.T) {
	var buf bytes.Buffer
	if err := ObservationFormInline("sess-1", "obs-list").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, " required") {
		t.Errorf("ObservationFormInline não deve ter 'required' no textarea — está aninhado no form principal e bloqueia 'Salvar Alterações'")
	}
}

// TestInterventionFormInline_NoRequired mesma razão — inline form não deve ter required.
func TestInterventionFormInline_NoRequired(t *testing.T) {
	var buf bytes.Buffer
	if err := InterventionFormInline("sess-1", "intv-list").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, " required") {
		t.Errorf("InterventionFormInline não deve ter 'required' no textarea — está aninhado no form principal e bloqueia 'Salvar Alterações'")
	}
}

// TestObservationForm_StandaloneKeepsRequired verifica que o form standalone
// (usado na tela de detalhe, não aninhado) mantém required.
func TestObservationForm_StandaloneKeepsRequired(t *testing.T) {
	var buf bytes.Buffer
	if err := ObservationForm("sess-1").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), " required") {
		t.Errorf("ObservationForm standalone deve manter 'required' para validação client-side")
	}
}
