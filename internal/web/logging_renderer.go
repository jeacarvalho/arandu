package web

import (
	"log"
	"net/http"
)

// LoggingRenderer implements TemplateRenderer interface with logging
type LoggingRenderer struct{}

// NewLoggingRenderer creates a new logging renderer
func NewLoggingRenderer() *LoggingRenderer {
	return &LoggingRenderer{}
}

// ExecuteTemplate logs a warning and returns an error (deprecated - use templ components)
func (r *LoggingRenderer) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	log.Printf("⚠️  WARNING: ExecuteTemplate called with template '%s' - This is deprecated! Use templ components instead.", name)
	log.Printf("   Stack trace: This indicates legacy template usage that should be migrated to templ components.")
	log.Printf("   Action required: Update handler to use templ.Component.Render() instead of ExecuteTemplate()")

	// Write a helpful error message to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`
		<div class="alert alert-error">
			<h3>⚠️ Template System Deprecated</h3>
			<p>The template system using ExecuteTemplate() is deprecated.</p>
			<p>Please use templ components instead. Check the logs for details.</p>
			<p><strong>Template attempted:</strong> ` + name + `</p>
		</div>
	`))

	return nil
}
