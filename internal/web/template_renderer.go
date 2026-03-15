package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// TemplateRendererAdapter adapts the old template system to the new TemplateRenderer interface
type TemplateRendererAdapter struct {
	templates *template.Template
}

// NewTemplateRendererAdapter creates a new adapter for the template system
func NewTemplateRendererAdapter(templatePath string) *TemplateRendererAdapter {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"now": func() time.Time {
			return time.Now()
		},
		"dateFormat": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
		"ToUpper": strings.ToUpper,
	}

	templates := template.New("").Funcs(funcMap)

	// Load all HTML templates
	templateFiles, err := filepath.Glob(filepath.Join(templatePath, "*.html"))
	if err != nil {
		log.Fatalf("Failed to find template files: %v", err)
	}

	if len(templateFiles) == 0 {
		log.Fatalf("No template files found in %s", templatePath)
	}

	templates, err = templates.ParseFiles(templateFiles...)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	return &TemplateRendererAdapter{
		templates: templates,
	}
}

// ExecuteTemplate implements the TemplateRenderer interface
func (a *TemplateRendererAdapter) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	return a.templates.ExecuteTemplate(w, name, data)
}
