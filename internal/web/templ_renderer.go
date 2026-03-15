//go:build templ

package web

import (
	"net/http"

	"github.com/a-h/templ"
)

// TemplRenderer wraps templ components for HTTP response
type TemplRenderer struct{}

func NewTemplRenderer() *TemplRenderer {
	return &TemplRenderer{}
}

// RenderHTMX renders a templ component for HTMX requests
func (r *TemplRenderer) RenderHTMX(w http.ResponseWriter, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Render(nil, w)
}

// RenderPage wraps a templ component in the HTML layout and renders
func (r *TemplRenderer) RenderPage(w http.ResponseWriter, pageTitle string, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render component to string first, then wrap in layout
	// For now, render directly - layout will be added as component
	layout := BaseLayout(pageTitle, c)
	layout.Render(nil, w)
}
