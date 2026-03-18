# 🏗️ Padrões Arquiteturais - Arandu

**Última atualização:** $(date +"%d de %B de %Y")
**Foco:** Arquitetura Web (Go + templ + HTMX)

> 📚 **Consulte também:** [MASTER_LEARNINGS.md](./MASTER_LEARNINGS.md)

---

## 📋 Índice

1. [Estrutura de Camadas](#estrutura-de-camadas)
2. [Handlers e ViewModels](#handlers-e-viewmodels)
3. [HTMX e Renderização Contextual](#htmx-e-renderização-contextual)
4. [Injeção de Dependência](#injeção-de-dependência)
5. [Tratamento de Erros](#tratamento-de-erros)
6. [Componentes Templ](#componentes-templ)
7. [Referências de Código](#referências-de-código)

---

## Estrutura de Camadas

### Princípios Fundamentais

**Clean Architecture aplicada:**

```
┌─────────────────────────────────────────────────────────┐
│                    Request HTTP                         │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Handler (web/handlers/)                                │
│  1. Extrai parâmetros                                   │
│  2. Valida input básico                                 │
│  3. Chama Service                                       │
│  4. Mapeia para ViewModel                               │
│  5. Verifica HX-Request                                 │
│  6. Renderiza template                                  │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Service (application/services/)                        │
│  - Lógica de negócio                                    │
│  - Orquestração de entidades de domínio                 │
│  - Validações complexas                                 │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Repository (infrastructure/repository/)                │
│  - Persistência em SQLite                               │
│  - CRUD operations                                      │
└─────────────────────────────────────────────────────────┘
```

### Regras de Ouro

#### ✅ Independência de Domínio
Handlers **NUNCA** contêm lógica de negócio. Apenas orquestram.

#### ✅ Consciência HTMX
Sempre verifique `HX-Request` header:
- `true` → Renderiza fragmento (`patient-content`, `session-content`)
- `false` → Renderiza página completa (`layout`)

#### ✅ ViewModels Fortemente Tipados
**NUNCA** passe entidades de domínio para templates. Sempre crie structs específicas.

---

## Handlers e ViewModels

### Estrutura do Handler

```go
// internal/web/handlers/patient_handler.go
type PatientHandler struct {
    patientService  application.PatientService
    sessionService  application.SessionService
    insightService  application.InsightService
    templates       *template.Template
}

// Interfaces específicas (segregação)
type PatientService interface {
    GetByID(ctx context.Context, id string) (*domain.Patient, error)
    List(ctx context.Context) ([]*domain.Patient, error)
    Create(ctx context.Context, input application.CreatePatientInput) (*domain.Patient, error)
}

// Apenas o necessário para este handler
type SessionService interface {
    GetByPatientID(ctx context.Context, patientID string) ([]*domain.Session, error)
}
```

### ViewModels (Proteção do Domínio)

**Problema:** Templates acessam campos internos de entidades.

**Solução:** ViewModels intermediários:

```go
// ViewModel para página de paciente
type PatientViewData struct {
    Patient   PatientViewModel
    Sessions  []SessionViewModel
    Insights  []InsightViewModel
    Error     string
}

// ViewModel específico para paciente
type PatientViewModel struct {
    ID        string
    Name      string
    Notes     string
    CreatedAt string // Já formatada
    UpdatedAt string // Já formatada
}

// Mapeamento no handler
func (h *PatientHandler) mapToPatientViewModel(patient *domain.Patient) PatientViewModel {
    return PatientViewModel{
        ID:        patient.ID,
        Name:      patient.Name,
        Notes:     patient.Notes,
        CreatedAt: patient.CreatedAt.Format("02/01/2006 15:04"),
        UpdatedAt: patient.UpdatedAt.Format("02/01/2006 15:04"),
    }
}
```

**Benefícios:**
1. **Encapsulamento:** Domínio protegido de mudanças na UI
2. **Formatação:** Datas, números, textos formatados no ViewModel
3. **Performance:** Apenas campos necessários são passados
4. **Segurança:** Evita expor campos sensíveis

---

## HTMX e Renderização Contextual

### Verificação de Contexto HTMX

```go
func (h *PatientHandler) Show(w http.ResponseWriter, r *http.Request) {
    // ... lógica para obter dados ...
    
    data := PatientViewData{
        Patient:  patientVM,
        Sessions: sessionVMs,
        Error:    "",
    }
    
    // VERIFICAÇÃO CRÍTICA
    if r.Header.Get("HX-Request") == "true" {
        // HTMX: retorna apenas o fragmento
        h.templates.ExecuteTemplate(w, "patient-content", data)
    } else {
        // Requisição normal: retorna página completa
        h.templates.ExecuteTemplate(w, "layout", data)
    }
}
```

### Fragmentos de Template

**Estrutura de templates:**

```html
<!-- web/templates/patient.html -->
{{define "content"}}
    {{template "patient-content" .}}
{{end}}

{{define "patient-content"}}
    <!-- Fragmento HTMX reutilizável -->
    <div id="patient-detail" class="p-6">
        {{if .Error}}
            <div class="error-message">{{.Error}}</div>
        {{else}}
            <h1 class="text-2xl font-bold">{{.Patient.Name}}</h1>
            <!-- ... resto do conteúdo ... -->
        {{end}}
    </div>
{{end}}
```

### Swaps e Triggers HTMX

**Padrões estabelecidos:**

```html
<!-- Infinite Scroll -->
<div 
  hx-get="/patients/{id}/history?offset={next}"
  hx-trigger="revealed"
  hx-swap="afterend"
  hx-indicator="#loading">
</div>

<!-- Busca com debounce -->
<input 
  type="search"
  hx-get="/patients/search"
  hx-trigger="keyup changed delay:500ms"
  hx-target="#search-results"
  hx-indicator="#search-loading">

<!-- Form submission -->
<form 
  hx-post="/patients"
  hx-target="#patient-list"
  hx-swap="beforeend">
</form>
```

---

## Injeção de Dependência

### Setup no main.go

```go
// cmd/arandu/main.go
func main() {
    // Setup database
    db := setupDatabase()
    
    // Repositories
    patientRepo := repository.NewPatientRepository(db)
    sessionRepo := repository.NewSessionRepository(db)
    
    // Services
    patientService := application.NewPatientService(patientRepo)
    sessionService := application.NewSessionService(sessionRepo)
    
    // Handlers
    patientHandler := handlers.NewPatientHandler(
        patientService,
        sessionService,
        templates,
    )
    
    // Routes
    http.HandleFunc("/patients", patientHandler.List)
    http.HandleFunc("/patient/{id}", patientHandler.Show)
    // ...
}
```

### Benefícios

1. **Testabilidade:** Mock fácil de dependências
2. **Flexibilidade:** Troca de implementações sem modificar handlers
3. **Manutenção:** Dependências explícitas e documentadas
4. **Inicialização controlada:** Ordem de criação garantida

---

## Tratamento de Erros

### Erros Contextuais (HTMX vs Full Page)

```go
func (h *PatientHandler) renderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
    w.WriteHeader(statusCode)
    data := PatientViewData{Error: message}
    
    if r.Header.Get("HX-Request") == "true" {
        // HTMX: fragmento de erro amigável
        h.templates.ExecuteTemplate(w, "error-fragment", data)
        return
    }
    
    // Full page: página de erro completa
    h.templates.ExecuteTemplate(w, "layout", data)
}

// Uso no handler
func (h *PatientHandler) Show(w http.ResponseWriter, r *http.Request) {
    patient, err := h.patientService.GetByID(r.Context(), patientID)
    if err != nil {
        h.renderError(w, r, "Paciente não encontrado", http.StatusNotFound)
        return
    }
    // ...
}
```

### Fragmento de Erro HTMX

```html
<!-- web/templates/error-fragment.html -->
{{define "error-fragment"}}
<div class="error-container bg-red-50 border-l-4 border-red-500 p-4 my-4">
    <div class="flex">
        <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
            </svg>
        </div>
        <div class="ml-3">
            <p class="text-sm text-red-700">{{.Error}}</p>
        </div>
    </div>
</div>
{{end}}
```

---

## Componentes Templ

### Estrutura de Componentes

```
web/components/
├── layout/
│   ├── layout.templ          # Layout base
│   └── sidebar.templ         # Sidebar navegação
├── patient/
│   ├── list.templ           # Lista de pacientes
│   ├── detail.templ         # Detalhe do paciente
│   ├── form.templ           # Formulário
│   └── search.templ         # Componente de busca
├── session/
│   ├── detail.templ         # Detalhe da sessão
│   ├── list.templ           # Lista de sessões
│   └── form.templ           # Formulário de sessão
└── dashboard/
    └── dashboard.templ      # Dashboard clínico
```

### Exemplo de Componente

```templ
// web/components/patient/detail.templ
package patient

import "github.com/a-h/templ"

type PatientDetailData struct {
    ID        string
    Name      string
    Notes     string
    CreatedAt string
}

templ PatientDetail(data PatientDetailData) {
    <div id="patient-detail" class="p-6 bg-white rounded-lg shadow">
        <div class="mb-6">
            <h1 class="text-3xl font-bold text-gray-900">{ data.Name }</h1>
            <p class="text-sm text-gray-500 mt-1">
                Cadastrado em { data.CreatedAt }
            </p>
        </div>
        
        @if data.Notes != "" {
            <div class="mt-6">
                <h2 class="text-lg font-semibold text-gray-800 mb-2">Observações</h2>
                <div class="font-clinical text-gray-700 bg-gray-50 p-4 rounded">
                    { data.Notes }
                </div>
            </div>
        }
    </div>
}
```

### Regras para Templ

1. **NUNCA** adicione `import "github.com/a-h/templ"` manualmente
2. **SEMPRE** execute `templ generate` após modificar `.templ`
3. **NÃO** use `//go:build templ` - compile nativamente
4. **MANTENHA** componentes pequenos e focados (< 200 linhas)

---

## Referências de Código

### Handlers Modelo

1. **`internal/web/handlers/patient_handler.go`** - Handler completo com ViewModels
2. **`internal/web/handlers/session_handler.go`** - Handler com HTMX e erros
3. **`internal/web/handlers/biopsychosocial_handler.go`** - Handler recente com validações

### Templates Modelo

1. **`web/templates/patient.html`** - Template com fragments HTMX
2. **`web/templates/error-fragment.html`** - Tratamento de erros contextual
3. **`web/components/patient/detail.templ`** - Componente Templ moderno

### Documentação Relacionada

1. **`docs/architecture/WEB_LAYER_PATTERN.md`** - Padrão completo da camada web
2. **`docs/architecture/system_structure.md`** - Visão geral da arquitetura
3. **`docs/architecture/ROUTE_CONVENTIONS.md`** - Convenções de rotas

---

## Checklist de Implementação

### ✅ Antes de Criar Handler

- [ ] Verifique se já existe handler similar para reutilizar padrões
- [ ] Defina interfaces específicas para serviços necessários
- [ ] Planeje ViewModels necessários

### ✅ Durante Implementação

- [ ] Handler verifica `HX-Request` header
- [ ] ViewModels protegem domínio (nunca passar entidades)
- [ ] Tratamento de erros contextual (HTMX vs full page)
- [ ] Injeção de dependência via interfaces

### ✅ Após Implementação

- [ ] Teste ambos cenários: requisição direta e via HTMX
- [ ] Execute `arandu_guard.sh` para verificar integridade
- [ ] Valide rotas existentes não quebraram

---

*Documento parte da refatoração da arquitetura do Arandu - Março 2026*