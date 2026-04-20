# Guia do Agente - Implementação no Arandu

**Para agentes de IA que vão trabalhar no projeto Arandu**

## 🚀 Quick Start para Novas Features

### 1. Identificar o Tipo de Task

**A. Novo Handler (CRUD de entidade)**
- Ex: "Adicionar CRUD para Observations"
- **Solução**: Criar novo handler seguindo padrão estabelecido

**B. Nova Funcionalidade em Handler Existente**
- Ex: "Adicionar filtro por data na lista de pacientes"
- **Solução**: Estender handler existente

**C. Correção de Bug**
- Ex: "Erro ao criar sessão com data inválida"
- **Solução**: Analisar fluxo e corrigir no local apropriado

### 2. Estrutura de Referência Rápida

```
📁 ESTRUTURA CHAVE PARA IMPLEMENTAÇÃO:
internal/web/handlers/
├── patient_handler.go     # MODELO PARA NOVOS HANDLERS
├── session_handler.go     # OUTRO MODELO
└── (seu_novo_handler.go) # VOCÊ CRIARÁ AQUI

web/templates/
├── layout.html           # LAYOUT BASE
├── patients.html         # MODELO DE TEMPLATE
├── patient.html          # OUTRO MODELO
└── (seu_novo_template.html) # VOCÊ CRIARÁ AQUI

cmd/arandu/main.go        # REGISTRO DE ROTAS
```

### 3. Checklist para Novo Handler

```go
// PASSO 1: Criar arquivo handler
// File: internal/web/handlers/observation_handler.go

package handlers

import (
    "context"
    "net/http"
    "time"
    
    "arandu/internal/application/services"
    "arandu/internal/domain/observation"
)

// PASSO 2: Definir ViewModels (PROTEÇÃO DO DOMÍNIO)
type ObservationViewModel struct {
    ID        string
    Content   string
    CreatedAt string  // Já formatado!
}

type ObservationViewData struct {
    Observation *ObservationViewModel
    Error       string
    FormData    *ObservationFormValues  // Para forms
}

// PASSO 3: Definir interfaces (INJEÇÃO DE DEPENDÊNCIA)
type ObservationService interface {
    GetObservation(ctx context.Context, id string) (*observation.Observation, error)
    CreateObservation(ctx context.Context, input services.CreateObservationInput) (*observation.Observation, error)
}

// PASSO 4: Criar handler struct
type ObservationHandler struct {
    observationService ObservationService
    templates          TemplateRenderer
}

// PASSO 5: Construtor com DI
func NewObservationHandler(
    observationService ObservationService,
    templates TemplateRenderer,
) *ObservationHandler {
    return &ObservationHandler{
        observationService: observationService,
        templates:          templates,
    }
}

// PASSO 6: Métodos HTTP (SIGA O PADRÃO!)
func (h *ObservationHandler) Show(w http.ResponseWriter, r *http.Request) {
    // 1. Extrair parâmetros
    id := strings.TrimPrefix(r.URL.Path, "/observation/")
    
    // 2. Chamar serviço
    obs, err := h.observationService.GetObservation(r.Context(), id)
    if err != nil {
        h.renderError(w, r, err.Error(), http.StatusNotFound)
        return
    }
    
    // 3. Mapear para ViewModel
    data := ObservationViewData{
        Observation: &ObservationViewModel{
            ID:        obs.ID,
            Content:   obs.Content,
            CreatedAt: obs.CreatedAt.Format("02/01/2006 15:04"),
        },
    }
    
    // 4. Renderizar com consciência HTMX
    if r.Header.Get("HX-Request") == "true" {
        h.templates.ExecuteTemplate(w, "observation-content", data)
        return
    }
    
    h.templates.ExecuteTemplate(w, "layout", data)
}

// PASSO 7: Método de erro (REUTILIZAR)
func (h *ObservationHandler) renderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(statusCode)
    
    data := ErrorViewData{
        Message: message,
        RetryURL: "/observations",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        h.templates.ExecuteTemplate(w, "error-fragment", data)
        return
    }
    
    h.templates.ExecuteTemplate(w, "layout", data)
}
```

### 4. Checklist para Novo Template

```html
<!-- File: web/templates/observation.html -->

{{/* 1. DEFINIÇÃO PARA FULL-PAGE RENDERING */}}
{{define "content"}}
{{template "observation-content" .}}
{{end}}

{{/* 2. DEFINIÇÃO PARA HTMX FRAGMENT */}}
{{define "observation-content"}}
{{if .Error}}
    {{template "error-fragment" .}}
{{else if .Observation}}
<div class="max-w-4xl mx-auto">
    <header class="mb-8">
        <h1 class="text-2xl font-bold text-arandu-text">Observação</h1>
        <p class="text-arandu-text-secondary mt-2">{{.Observation.CreatedAt}}</p>
    </header>
    
    <div class="prose max-w-none">
        <p>{{.Observation.Content}}</p>
    </div>
</div>
{{end}}
{{end}}
```

### 5. Registrar no Main.go

```go
// File: cmd/arandu/main.go

// Adicionar import
"arandu/internal/web/handlers"

// Criar serviço (se necessário)
observationService := services.NewObservationService(observationRepo)

// Criar adaptador (se necessário)
observationServiceAdapter := web.NewObservationServiceAdapter(observationService)

// Criar handler
observationHandler := handlers.NewObservationHandler(observationServiceAdapter, templateRenderer)

// Registrar rotas
mux.HandleFunc("/observation/", observationHandler.Show)
mux.HandleFunc("/observations/new", observationHandler.New)
mux.HandleFunc("/observations", observationHandler.Create)
```

### 6. Padrões OBRIGATÓRIOS

#### ✅ SEMPRE FAZER:
1. **ViewModels**: Sempre criar structs específicas para templates
2. **HTMX Awareness**: Sempre verificar `r.Header.Get("HX-Request")`
3. **Error Handling**: Usar `renderError` ou método similar
4. **Dependency Injection**: Handlers recebem interfaces, não concretões
5. **Template Fragments**: Templates devem definir fragments nomeados

#### ❌ NUNCA FAZER:
1. **Expor Domínio**: Nunca passar entidades de domínio diretamente para templates
2. **Lógica de Negócio em Handlers**: Validações, cálculos, etc. devem estar em services
3. **Ignorar HTMX**: Sempre considerar ambos cenários (full-page e fragment)
4. **Hardcode Dependencies**: Usar `NewHandler(service, repo, db)` - usar interfaces!

### 7. Testes Rápidos

```bash
# 1. Compilar
go build ./cmd/arandu

# 2. Executar
go run ./cmd/arandu &

# 3. Testar endpoint
curl http://localhost:8080/seu-novo-endpoint

# 4. Testar HTMX
curl -H "HX-Request: true" http://localhost:8080/seu-novo-endpoint

# 5. Parar servidor
pkill -f "go run ./cmd/arandu"
```

### 8. Exemplos Práticos

#### Exemplo 1: Adicionar Filtro em Lista Existente
```go
// Em patient_handler.go, modificar ListPatients:
func (h *PatientHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
    // Novo: extrair parâmetro de filtro
    filter := r.URL.Query().Get("search")
    
    // Chamar serviço com filtro (precisa atualizar interface!)
    patients, err := h.patientService.ListPatients(r.Context(), filter)
    // ... resto do código
}
```

#### Exemplo 2: Adicionar Campo em Formulário
```go
// 1. Adicionar campo no ViewModel
type PatientViewModel struct {
    ID        string
    Name      string
    Email     string  // NOVO CAMPO
    CreatedAt string
}

// 2. Adicionar no template
// 3. Adicionar no service (se necessário validação)
// 4. Adicionar no repositório (se persistir)
```

### 9. Troubleshooting Comum

#### Problema: "template: X: function "now" not defined"
**Solução**: Adicionar função ao TemplateRenderer:
```go
funcMap := template.FuncMap{
    "now": func() time.Time { return time.Now() },
    "dateFormat": func(t time.Time, layout string) string {
        return t.Format(layout)
    },
}
```

#### Problema: "multiple registrations for /path/"
**Solução**: Verificar rotas duplicadas em `main.go`

#### Problema: Handler não compila - "interface mismatch"
**Solução**: Criar adaptador em `internal/web/service_adapters.go`

### 10. Referências Rápidas

| Arquivo | Propósito | Modelo a Seguir |
|---------|-----------|-----------------|
| `internal/web/handlers/patient_handler.go` | Handler completo | ✅ MELHOR MODELO |
| `web/templates/patients.html` | Template com fragments | ✅ |
| `internal/web/service_adapters.go` | Adaptadores de interfaces | ✅ |
| `internal/web/template_renderer.go` | Sistema de templates | ✅ |
| `docs/architecture/WEB_LAYER_PATTERN.md` | Documentação detalhada | 📚 |

### 11. Fluxo de Decisão para Agentes

```
NOVA TASK RECEBIDA
    ↓
É CRUD de nova entidade?
    ├── SIM → Criar novo handler (seção 3)
    └── NÃO → É extensão de funcionalidade existente?
            ├── SIM → Modificar handler existente
            └── NÃO → É correção de bug?
                    ├── SIM → Analisar e corrigir
                    └── NÃO → Perguntar ao usuário
```

### 12. Dicas de Produtividade

1. **Copiar e Adaptar**: Use `patient_handler.go` como template base
2. **Testar Incrementalmente**: Compile após cada mudança significativa
3. **Verificar Logs**: Servidor mostra erros de template em tempo real
4. **Usar Adaptadores Existentes**: Reutilize `SessionServiceAdapter` como modelo
5. **Manter ViewModels Simples**: Apenas dados necessários para a view

---

### 13. Estrutura de Testes

O Arandu consolidou todos os testes em `tests/`:

```
tests/
├── run_unit.sh              # Entry point: unit tests Go
├── run_e2e.sh              # Entry point: E2E (Go + JS + Shell)
└── e2e/
    ├── js/                  # Testes Playwright (Node.js)
    │   └── run_js_tests.sh
    ├── shell/              # Testes de scripts (Shell)
    │   └── run_shell_tests.sh
    └── *.go               # Testes E2E em Go
```

#### Executando Testes

```bash
# Unit tests (Makefile)
make test-unit
tests/run_unit.sh           # com coverage: tests/run_unit.sh -cover

# E2E (Makefile)
make test-e2e

# E2E detalhado
tests/run_e2e.sh go        # apenas Go E2E
tests/run_e2e.sh js        # apenas JS/Playwright
tests/run_e2e.sh shell     # apenas scripts shell
tests/run_e2e.sh all       # tudo
```

#### Makefile Targets

```bash
make test     # unit + e2e
make test-unit
make test-e2e
make test-all # alias para test
```

---

**Última Atualização**: Abril 2026  
**Status**: Em produção  
**Para Agentes**: Este guia deve ser consultado ANTES de iniciar qualquer implementation