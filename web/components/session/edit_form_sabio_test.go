package session

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func makeEditFormData() EditSessionFormData {
	return EditSessionFormData{
		SessionID:     "sess-42",
		SessionNumber: 5,
		PatientName:   "André Barbosa",
		FormData: &SessionFormValues{
			PatientID: "pat-01",
			Date:      "2026-04-02",
			Summary:   "Sessão produtiva com emergência de material relevante.",
		},
		Observations: []Observation{
			{ID: "o1", Content: "Questionamento existencial.", CreatedAt: "02/04/2026 20:18"},
		},
		Interventions: []Intervention{
			{ID: "i1", Content: "Escuta ativa e validação.", CreatedAt: "02/04/2026 20:22"},
		},
	}
}

// TestEditSessionForm_HasSabioHeader verifica o header editorial da sessão.
func TestEditSessionForm_HasSabioHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-session-header") {
		t.Error("Sessão deve ter div.sabio-session-header")
	}
	if !strings.Contains(html, "André Barbosa") {
		t.Error("Header deve conter o nome do paciente")
	}
}

// TestEditSessionForm_HeaderShowsSessionNumber verifica número da sessão no header.
func TestEditSessionForm_HeaderShowsSessionNumber(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Sessão 05") {
		t.Error("Header deve exibir 'Sessão 05' (zero-padded)")
	}
}

// TestEditSessionForm_HasNotesGrid verifica o grid de 2 colunas de notas.
func TestEditSessionForm_HasNotesGrid(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-notes-grid") {
		t.Error("Sessão deve ter div.sabio-notes-grid para as 2 colunas")
	}
}

// TestEditSessionForm_HasObservationsColumn verifica coluna de observações.
func TestEditSessionForm_HasObservationsColumn(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "Observações clínicas") {
		t.Error("Deve ter coluna 'Observações clínicas'")
	}
	if !strings.Contains(html, "Escuta") {
		t.Error("Coluna de observações deve ter eyebrow 'Escuta'")
	}
}

// TestEditSessionForm_HasInterventionsColumn verifica coluna de intervenções.
func TestEditSessionForm_HasInterventionsColumn(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "Intervenções terapêuticas") {
		t.Error("Deve ter coluna 'Intervenções terapêuticas'")
	}
	if !strings.Contains(html, "Ação") {
		t.Error("Coluna de intervenções deve ter eyebrow 'Ação'")
	}
}

// TestEditSessionForm_PreservesHtmxObservations verifica HTMX para observações.
func TestEditSessionForm_PreservesHtmxObservations(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "/session/sess-42/observations") {
		t.Error("Formulário de observações deve ter hx-post para /observations")
	}
	if !strings.Contains(html, "observations-list") {
		t.Error("Deve ter id=observations-list como target HTMX")
	}
}

// TestEditSessionForm_PreservesHtmxInterventions verifica HTMX para intervenções.
func TestEditSessionForm_PreservesHtmxInterventions(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "/session/sess-42/interventions") {
		t.Error("Formulário de intervenções deve ter hx-post para /interventions")
	}
	if !strings.Contains(html, "interventions-list") {
		t.Error("Deve ter id=interventions-list como target HTMX")
	}
}

// TestEditSessionForm_HasSynthesisCard verifica o card de síntese.
func TestEditSessionForm_HasSynthesisCard(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-synthesis-card") {
		t.Error("Sessão deve ter sabio-synthesis-card para o resumo")
	}
	if !strings.Contains(html, "Síntese") {
		t.Error("Card de síntese deve ter título 'Síntese'")
	}
}

// TestEditSessionForm_SummaryHasAutosave verifica autosave no textarea de resumo.
func TestEditSessionForm_SummaryHasAutosave(t *testing.T) {
	var buf bytes.Buffer
	if err := EditSessionForm(makeEditFormData()).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "/session/sess-42/summary") {
		t.Error("Textarea de resumo deve ter hx-patch para /summary (autosave)")
	}
	if !strings.Contains(html, "delay:") {
		t.Error("Autosave deve ter delay no hx-trigger")
	}
}
