package layout

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestShellTopbar_HasLLMButton verifica que a topbar tem o botão "Arandu"
// que abre o drawer de inteligência clínica (⌘J).
func TestShellTopbar_HasLLMButton(t *testing.T) {
	config := ShellConfig{PageTitle: "Dashboard", ActivePage: "dashboard", ShowSidebar: true}
	var buf bytes.Buffer
	if err := ShellTopbar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-llm-btn") {
		t.Errorf("Topbar deve ter botão LLM com classe sabio-llm-btn")
	}
	if !strings.Contains(html, "Arandu") {
		t.Errorf("Topbar deve ter botão LLM com texto 'Arandu'")
	}
}

// TestShellTopbar_HasSearchInput verifica que a topbar ainda tem campo de busca.
func TestShellTopbar_HasSearchInput(t *testing.T) {
	config := ShellConfig{ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellTopbar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "<input") {
		t.Errorf("Topbar deve ter campo de busca")
	}
}

// TestShellTopbar_HasBreadcrumb verifica que a topbar renderiza breadcrumb
// derivado de ActivePage quando Breadcrumb não está explícito.
func TestShellTopbar_HasBreadcrumb(t *testing.T) {
	config := ShellConfig{ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellTopbar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-breadcrumb") {
		t.Errorf("Topbar deve ter nav com classe sabio-breadcrumb")
	}
	// Dashboard deve produzir breadcrumb "Hoje"
	if !strings.Contains(html, "Hoje") {
		t.Errorf("Dashboard deve ter breadcrumb 'Hoje'")
	}
}

// TestShellTopbar_BreadcrumbFromConfig verifica que Breadcrumb explícito é usado
// quando definido no ShellConfig.
func TestShellTopbar_BreadcrumbFromConfig(t *testing.T) {
	config := ShellConfig{
		ActivePage: "patient-summary",
		Breadcrumb: []string{"Pacientes", "André Barbosa"},
	}
	var buf bytes.Buffer
	if err := ShellTopbar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "Pacientes") {
		t.Errorf("Breadcrumb explícito deve conter 'Pacientes'")
	}
	if !strings.Contains(html, "André Barbosa") {
		t.Errorf("Breadcrumb explícito deve conter 'André Barbosa'")
	}
}

// TestShellSidebar_HasBrandMarkSVG verifica que a sidebar tem o BrandMark
// como SVG inline (glifo folha+olho da identidade Sábio).
func TestShellSidebar_HasBrandMarkSVG(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "<svg") {
		t.Errorf("Sidebar deve ter BrandMark como SVG inline")
	}
}

// TestShellSidebar_HasMenuSectionLabel verifica seção "MENU" na sidebar.
func TestShellSidebar_HasMenuSectionLabel(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "MENU") {
		t.Errorf("Sidebar deve ter label de seção 'MENU'")
	}
}

// TestShellSidebar_HasAtalhosLabel verifica seção "ATALHOS" na sidebar.
func TestShellSidebar_HasAtalhosLabel(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "ATALHOS") {
		t.Errorf("Sidebar deve ter label de seção 'ATALHOS'")
	}
}

// TestShellSidebar_NavItemUsesSabioClass verifica que os itens de nav usam
// a classe sabio-nav-item (não a classe legacy shell-sidebar-nav-item verde).
func TestShellSidebar_NavItemUsesSabioClass(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "sabio-nav-item") {
		t.Errorf("Items de nav devem usar classe sabio-nav-item (paleta Sábio)")
	}
	if strings.Contains(html, "shell-sidebar-nav-item") {
		t.Errorf("Items de nav não devem mais usar shell-sidebar-nav-item (classe legacy verde)")
	}
}

// TestShellSidebar_FooterHasCollapseButton verifica que o rodapé da sidebar
// tem o botão "Recolher" para colapsar a sidebar no desktop.
func TestShellSidebar_FooterHasCollapseButton(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Recolher") {
		t.Errorf("Sidebar deve ter botão 'Recolher' no rodapé")
	}
}

// TestShellSidebar_HasClinicoBrand verifica brand area com "Arandu" + "CLÍNICO".
func TestShellSidebar_HasClinicoBrand(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "Arandu") {
		t.Errorf("Sidebar deve ter 'Arandu' na área de brand")
	}
	if !strings.Contains(html, "CLÍNICO") {
		t.Errorf("Sidebar deve ter eyebrow 'CLÍNICO' na área de brand")
	}
}

// TestShellSidebar_HasUserAvatar verifica que o rodapé tem o avatar do usuário.
func TestShellSidebar_HasUserAvatar(t *testing.T) {
	config := ShellConfig{ShowSidebar: true, ActivePage: "dashboard", UserEmail: "helena@clinic.com"}
	var buf bytes.Buffer
	if err := ShellSidebar(config).Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "sabio-avatar") {
		t.Errorf("Sidebar deve ter div.sabio-avatar no rodapé")
	}
}
