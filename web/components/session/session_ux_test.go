package session

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"arandu/internal/domain/observation"
)

// ===========================================================================
// detail.templ — patient name + grid
// ===========================================================================

func TestSessionDetailView_ShowsPatientName(t *testing.T) {
	detail := SessionDetail{
		ID:          "sess-1",
		PatientID:   "pat-1",
		PatientName: "Maria Souza",
		Date:        "08/11/2025",
	}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Maria Souza") {
		t.Errorf("Header deve exibir o nome do paciente 'Maria Souza'. got: %s", buf.String()[:400])
	}
}

func TestSessionDetailView_PatientNameIsMainTitle(t *testing.T) {
	detail := SessionDetail{
		ID:          "sess-1",
		PatientID:   "pat-1",
		PatientName: "Ana Lima",
		Date:        "08/11/2025",
	}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	// Nome do paciente deve aparecer em page-title ou equivalente proeminente
	namIdx := strings.Index(html, "Ana Lima")
	if namIdx == -1 {
		t.Fatal("Nome do paciente não encontrado no HTML")
	}
	// Data da sessão deve aparecer como subtitle, não como título principal
	if !strings.Contains(html, "08/11/2025") {
		t.Errorf("Data da sessão deve aparecer no header")
	}
}

func TestSessionDetailView_GridExactlyTwoColumns(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-1", Date: "08/11/2025"}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	// Deve usar grid de 2 colunas fixas, não auto-fit
	if strings.Contains(html, "auto-fit") {
		t.Errorf("Grid não deve usar auto-fit — deve ser 2 colunas fixas para evitar 1 ou 4 colunas em telas largas")
	}
	if !strings.Contains(html, "grid-cols-2") && !strings.Contains(html, "session-detail-grid") {
		t.Errorf("Grid deve usar classe de 2 colunas fixas (grid-cols-2 ou session-detail-grid)")
	}
}

// ===========================================================================
// observation_item.templ — sem bloco de tag vazio
// ===========================================================================

func TestObservationItem_EmptyTags_NoInstructionText(t *testing.T) {
	obs := ObservationItemData{
		ID:        "obs-1",
		Content:   "Paciente relatou melhora no sono",
		CreatedAt: "01/04/2026 10:00",
		Tags:      nil, // nenhuma tag
	}
	var buf bytes.Buffer
	if err := ObservationItem(obs).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, "Nenhuma classificação ainda") {
		t.Errorf("Item sem tags não deve exibir texto 'Nenhuma classificação ainda' — ocupa espaço desnecessário")
	}
	if strings.Contains(html, "Clique no ícone de tags") {
		t.Errorf("Item sem tags não deve exibir instrução 'Clique no ícone de tags'")
	}
}

func TestObservationItem_EmptyTags_ContainerStillPresent(t *testing.T) {
	obs := ObservationItemData{
		ID:      "obs-42",
		Content: "Conteúdo",
		Tags:    nil,
	}
	var buf bytes.Buffer
	if err := ObservationItem(obs).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	// Container deve existir para HTMX poder injetar tags depois
	if !strings.Contains(buf.String(), "observation-obs-42-tags") {
		t.Errorf("Container #observation-{id}-tags deve existir mesmo sem tags (target HTMX)")
	}
}

func TestObservationItem_WithTags_ShowsTags(t *testing.T) {
	tagColor := "#4A7C59"
	obs := ObservationItemData{
		ID:      "obs-1",
		Content: "Conteúdo",
		Tags: []observation.ObservationTag{
			{
				ObservationID: "obs-1",
				Tag:           &observation.Tag{ID: "t1", Name: "Ansiedade", Color: tagColor},
				Intensity:     2,
			},
		},
	}
	var buf bytes.Buffer
	if err := ObservationItem(obs).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Ansiedade") {
		t.Errorf("Item com tags deve renderizar as tags")
	}
}

// ===========================================================================
// intervention_item.templ — sem bloco de tag vazio
// ===========================================================================

func TestInterventionItem_EmptyTags_NoInstructionText(t *testing.T) {
	intv := InterventionItemData{
		ID:        "intv-1",
		Content:   "TCC — reestruturação cognitiva",
		CreatedAt: "01/04/2026 10:00",
		Tags:      nil,
	}
	var buf bytes.Buffer
	if err := InterventionItem(intv).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if strings.Contains(html, "Nenhuma classificação ainda") {
		t.Errorf("InterventionItem sem tags não deve exibir 'Nenhuma classificação ainda'")
	}
	if strings.Contains(html, "Clique no ícone de tags") {
		t.Errorf("InterventionItem sem tags não deve exibir instrução de tags")
	}
}

func TestInterventionItem_EmptyTags_ContainerStillPresent(t *testing.T) {
	intv := InterventionItemData{ID: "intv-99", Content: "Conteúdo", Tags: nil}
	var buf bytes.Buffer
	if err := InterventionItem(intv).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "intervention-intv-99-tags") {
		t.Errorf("Container #intervention-{id}-tags deve existir mesmo sem tags")
	}
}

// ===========================================================================
// intervention_form.templ — botão consistente com observation_form
// ===========================================================================

func TestInterventionForm_ButtonMatchesObservationStyle(t *testing.T) {
	var bufIntv, bufObs bytes.Buffer
	if err := InterventionForm("sess-1").Render(context.Background(), &bufIntv); err != nil {
		t.Fatalf("InterventionForm Render() error: %v", err)
	}
	if err := ObservationForm("sess-1").Render(context.Background(), &bufObs); err != nil {
		t.Fatalf("ObservationForm Render() error: %v", err)
	}

	intvHTML := bufIntv.String()
	obsHTML := bufObs.String()

	// Ambos devem usar btn-primary no botão de submit
	if !strings.Contains(intvHTML, "btn-primary") {
		t.Errorf("InterventionForm deve usar 'btn-primary' no submit — atualmente usa btn-secondary, inconsistente com ObservationForm")
	}
	if !strings.Contains(obsHTML, "btn-primary") {
		t.Errorf("ObservationForm deve usar 'btn-primary' no submit")
	}
}
