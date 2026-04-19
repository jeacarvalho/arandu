package layout

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestLLMDrawer_HasDrawerContainer verifica o container principal do drawer.
func TestLLMDrawer_HasDrawerContainer(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-llm-drawer") {
		t.Error("LLM drawer deve ter aside.sabio-llm-drawer")
	}
}

// TestLLMDrawer_HasAlpineShowBinding verifica binding Alpine para abrir/fechar.
func TestLLMDrawer_HasAlpineShowBinding(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "llmOpen") {
		t.Error("Drawer deve ter binding Alpine para llmOpen")
	}
}

// TestLLMDrawer_HasBrandHeader verifica header com nome "Arandu" e ícone.
func TestLLMDrawer_HasBrandHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-llm-header") {
		t.Error("Drawer deve ter sabio-llm-header")
	}
	if !strings.Contains(html, "Arandu") {
		t.Error("Header deve conter o nome 'Arandu'")
	}
	if !strings.Contains(html, "Inteligência clínica") {
		t.Error("Header deve conter subtítulo 'Inteligência clínica'")
	}
}

// TestLLMDrawer_HasCloseButton verifica botão de fechar com Alpine.
func TestLLMDrawer_HasCloseButton(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "closeLLM") {
		t.Error("Drawer deve ter botão com $store.shell.closeLLM()")
	}
}

// TestLLMDrawer_HasBackdrop verifica overlay de backdrop.
func TestLLMDrawer_HasBackdrop(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-llm-backdrop") {
		t.Error("Drawer deve ter div.sabio-llm-backdrop")
	}
}

// TestLLMDrawer_HasSuggestions verifica chips de sugestões.
func TestLLMDrawer_HasSuggestions(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-llm-suggestions") {
		t.Error("Drawer deve ter sabio-llm-suggestions")
	}
}

// TestLLMDrawer_HasChatInput verifica textarea de input da conversa.
func TestLLMDrawer_HasChatInput(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-llm-input") {
		t.Error("Drawer deve ter sabio-llm-input")
	}
	if !strings.Contains(html, "<textarea") {
		t.Error("Input deve ter um textarea")
	}
}

// TestLLMDrawer_HasDisclaimer verifica aviso sobre julgamento clínico.
func TestLLMDrawer_HasDisclaimer(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "julgamento clínico") {
		t.Error("Drawer deve ter disclaimer sobre julgamento clínico")
	}
}

// TestLLMDrawer_HasCrossInsightsSection verifica seção "Padrões cruzados" com barras.
func TestLLMDrawer_HasCrossInsightsSection(t *testing.T) {
	var buf bytes.Buffer
	if err := LLMDrawer().Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-llm-insights") {
		t.Error("Drawer deve ter sabio-llm-insights (seção padrões cruzados)")
	}
	if !strings.Contains(html, "Padrões cruzados") {
		t.Error("Seção deve ter label 'Padrões cruzados'")
	}
	if !strings.Contains(html, "sabio-llm-insight-bar") {
		t.Error("Seção deve ter barras de progresso sabio-llm-insight-bar")
	}
}
