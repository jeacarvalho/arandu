package session

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestSessionDetailView_CardsHaveFlexClass verifica que os cards de observações e
// intervenções têm a classe session-detail-card, que habilita flex-column no CSS.
// Isso garante que os formulários "Adicionar" sempre fiquem no fundo dos cards,
// alinhados entre si independente do número de itens.
func TestSessionDetailView_CardsHaveFlexClass(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-1", Date: "24/03/2026", PatientName: "João"}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()

	count := strings.Count(html, "session-detail-card")
	if count < 2 {
		t.Errorf("detail deve ter 2 cards com classe 'session-detail-card' (observações + intervenções), encontrado: %d", count)
	}
}

// TestSessionDetailView_ListHasGrowClass verifica que os containers de lista
// têm a classe session-detail-list, que permite flex-grow no CSS.
func TestSessionDetailView_ListHasGrowClass(t *testing.T) {
	detail := SessionDetail{ID: "sess-1", PatientID: "pat-1", Date: "24/03/2026"}
	var buf bytes.Buffer
	if err := SessionDetailView(detail, nil, nil, "pat-1", "").Render(context.Background(), &buf); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	html := buf.String()

	count := strings.Count(html, "session-detail-list")
	if count < 2 {
		t.Errorf("detail deve ter 2 containers com classe 'session-detail-list' (obs + intv), encontrado: %d", count)
	}
}
