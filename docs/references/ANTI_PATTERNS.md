# Anti-Padrões Críticos - Sistema Arandu

## 🚫 VIOLAÇÕES CRÍTICAS (NÃO PERMITIDAS)

### 1. HTML Inline em Handlers
**❌ PROIBIDO:**
```go
// NUNCA faça isso:
w.Write([]byte("<html><body>...</body></html>"))
fmt.Fprintf(w, "<div>%s</div>", variable)
```

**✅ CORRETO:**
```go
// Use componentes Templ:
component := patientComponents.PatientDetail(data)
component.Render(r.Context(), w)
```

### 2. Uso de ExecuteTemplate ou html/template
**❌ PROIBIDO:**
```go
import "html/template"
h.templates.ExecuteTemplate(w, "template-name", data)
```

**✅ CORRETO:**
```go
import "github.com/a-h/templ"
import layoutComponents "arandu/web/components/layout"
layoutComponents.BaseWithContent("Título", component).Render(r.Context(), w)
```

### 3. Arquivos .html Soltos
**❌ PROIBIDO:**
- `templates/patients.html`
- `views/session.html`
- Qualquer arquivo `.html` fora de `web/static/`

**✅ CORRETO:**
- `web/components/patient/list.templ`
- `web/components/session/detail.templ`
- Arquivos `.html` apenas em `web/static/` (CSS, JS, imagens)

### 4. Handlers sem Componentes Templ
**❌ PROIBIDO:**
```go
func (h *Handler) SomeHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica sem componente Templ
    w.Write([]byte("conteúdo"))
}
```

**✅ CORRETO:**
```go
func (h *Handler) SomeHandler(w http.ResponseWriter, r *http.Request) {
    data := components.SomeComponentData{...}
    component := components.SomeComponent(data)
    
    // HTMX-aware rendering
    if r.Header.Get("HX-Request") == "true" {
        component.Render(r.Context(), w)
        return
    }
    
    // Full page with layout
    layoutComponents.BaseWithContent("Título", component).Render(r.Context(), w)
}
```

### 5. Interface TemplateRenderer
**❌ PROIBIDO:**
```go
type TemplateRenderer interface {
    ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error
}

type Handler struct {
    templates TemplateRenderer
}
```

**✅ CORRETO:**
```go
type Handler struct {
    service SomeService
    // Sem campo templates
}
```

## ⚠️ PADRÕES LEGADO (MIGRAR QUANDO POSSÍVEL)

### 1. DummyRenderer
**⚠️ LEGADO:**
```go
templateRenderer := web.NewDummyRenderer()
```

**✅ MODERNO:**
```go
// Remova completamente a necessidade de TemplateRenderer
// Handlers devem usar componentes Templ diretamente
```

### 2. TODO Comments para Migração
**⚠️ LEGADO:**
```go
// TODO: Migrate to templ
```

**✅ MODERNO:**
```go
// Já migrado para templ
```

## 🎯 PADRÕES ARQUITETURAIS (OBRIGATÓRIOS)

### 1. Estrutura de Componentes
```
web/components/
├── layout/           # Componentes de layout (Base, Sidebar, etc.)
│   ├── layout.templ
│   └── error.templ
├── patient/          # Componentes de paciente
│   ├── list.templ
│   ├── detail.templ
│   └── new_form.templ
└── session/          # Componentes de sessão
    ├── detail.templ
    ├── new_form.templ
    └── edit_form.templ
```

### 2. Handlers com Injeção de Dependência
```go
type PatientHandler struct {
    patientService PatientService
    sessionService SessionService
    // Sem templates, sem renderers
}

func NewPatientHandler(
    patientService PatientService,
    sessionService SessionService,
) *PatientHandler {
    return &PatientHandler{
        patientService: patientService,
        sessionService: sessionService,
    }
}
```

### 3. Renderização HTMX-aware
```go
func (h *Handler) SomeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    
    if r.Header.Get("HX-Request") == "true" {
        // Render fragment only
        component.Render(r.Context(), w)
        return
    }
    
    // Render full page with layout
    layoutComponents.BaseWithContent("Title", component).Render(r.Context(), w)
}
```

### 4. Tratamento de Erros com Componentes
```go
func (h *Handler) renderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(statusCode)
    
    isHTMX := r.Header.Get("HX-Request") == "true"
    errorData := layoutComponents.ErrorData{
        Error:  message,
        IsHTMX: isHTMX,
        Title:  "Erro",
    }
    
    if isHTMX {
        layoutComponents.ErrorFragment(errorData).Render(r.Context(), w)
        return
    }
    
    layoutComponents.BaseWithContent("Erro", 
        layoutComponents.ErrorFragment(errorData),
    ).Render(r.Context(), w)
}
```

## 🔧 FERRAMENTAS DE VALIDAÇÃO

### 1. Script de Validação
```bash
./scripts/arandu_validate_handlers.sh
```

### 2. Comandos de Build
```bash
# Verificar handlers
go build ./internal/web/handlers/...

# Build completo
go build ./cmd/arandu
```

### 3. Regeneração de Templ
```bash
# Instalar templ (se necessário)
go install github.com/a-h/templ/cmd/templ@latest

# Gerar arquivos
~/go/bin/templ generate
```

## 📚 DOCUMENTAÇÃO DE REFERÊNCIA

1. `docs/architecture/WEB_LAYER_PATTERN.md` - Padrões da camada web
2. `docs/architecture/AGENT_GUIDE.md` - Guia prático para agentes
3. `docs/design-system.md` - Sistema de design visual
4. `docs/learnings/REQ-01-02-02.md` - Aprendizado sobre imports em .templ

## 🚨 PROCESSO DE CORREÇÃO

### Passo 1: Identificar Violação
```bash
./scripts/arandu_validate_handlers.sh
```

### Passo 2: Analisar o Erro
- Verificar qual anti-padrão foi violado
- Localizar no código (`grep -r "pattern" internal/web/handlers/`)

### Passo 3: Corrigir
1. Criar componente Templ (se necessário)
2. Atualizar handler para usar componente
3. Remover código legado

### Passo 4: Validar
```bash
./scripts/arandu_validate_handlers.sh
go build ./internal/web/handlers/...
go build ./cmd/arandu
```

### Passo 5: Testar
- Executar testes existentes
- Testar manualmente a funcionalidade
- Verificar compatibilidade com HTMX

---
**Última atualização**: 16/03/2026  
**Status**: Sistema migrado para Templ  
**Próxima revisão**: Após nova funcionalidade web