package layout

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestToast_RendersMessage(t *testing.T) {
	var buf bytes.Buffer
	err := Toast("Observação adicionada", "success").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(buf.String(), "Observação adicionada") {
		t.Errorf("Toast deve renderizar a mensagem. got: %s", buf.String())
	}
}

func TestToast_HasOOBSwap(t *testing.T) {
	var buf bytes.Buffer
	err := Toast("Salvo", "success").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "hx-swap-oob") {
		t.Errorf("Toast deve ter hx-swap-oob para substituir o container. got: %s", html)
	}
	if !strings.Contains(html, "toast-container") {
		t.Errorf("Toast deve ter id='toast-container'. got: %s", html)
	}
}

func TestToast_SuccessKindHasGreenStyle(t *testing.T) {
	var buf bytes.Buffer
	err := Toast("Salvo", "success").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "toast-success") {
		t.Errorf("Toast success deve ter classe 'toast-success'. got: %s", html)
	}
}

func TestToast_ErrorKindHasErrorStyle(t *testing.T) {
	var buf bytes.Buffer
	err := Toast("Erro ao salvar", "error").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, "toast-error") {
		t.Errorf("Toast error deve ter classe 'toast-error'. got: %s", html)
	}
}

func TestToast_HasAutoDismiss(t *testing.T) {
	var buf bytes.Buffer
	err := Toast("Salvo", "success").Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()
	// Alpine.js x-init com setTimeout para auto-dismiss
	if !strings.Contains(html, "setTimeout") {
		t.Errorf("Toast deve ter auto-dismiss via setTimeout. got: %s", html)
	}
}
