package web

import "net/http"

// DummyRenderer implements TemplateRenderer interface but does nothing
type DummyRenderer struct{}

// NewDummyRenderer creates a new dummy renderer
func NewDummyRenderer() *DummyRenderer {
	return &DummyRenderer{}
}

// ExecuteTemplate does nothing (compatibility layer)
func (r *DummyRenderer) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	// Return nil - handlers should use templ components directly
	return nil
}
