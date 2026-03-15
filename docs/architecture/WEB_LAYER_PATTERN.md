# Padrão da Camada Web - Arandu

**Documento de Referência Arquitetural**

---

## Visão Geral

Este documento descreve o padrão arquitetural adotado para a camada web do projeto Arandu, seguindo os princípios de Clean Architecture, Domain-Driven Design (DDD) e integração consciente com HTMX.

---

## Regras de Ouro (Inquestionáveis)

### 1. Independência de Domínio

**Handlers não contêm lógica de negócio.** Eles apenas:

1. **Decodificam o Request** - Extraem parâmetros, query strings e body
2. **Chamam o Application Service** - Delegam a lógica para a camada de aplicação
3. **Mapeiam o resultado para um ViewModel** - Transformam entidades em dados de view
4. **Renderizam o Template** - Executam o template apropriado

```go
// ❌ ERRADO: Handler com lógica de negócio
func (h *Handler) CreatePatient(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    if name == "" {
        // Validação de negócio no handler!
        http.Error(w, "Nome é obrigatório", http.StatusBadRequest)
        return
    }
    // ...
}

// ✅ CORRETO: Handler apenas orquestra
func (h *Handler) CreatePatient(w http.ResponseWriter, r *http.Request) {
    input := services.CreatePatientInput{
        Name: r.FormValue("name"),
        Notes: r.FormValue("notes"),
    }
    
    patient, err := h.service.CreatePatient(r.Context(), input)
    if err != nil {
        h.renderError(w, r, err.Error(), http.StatusBadRequest)
        return
    }
    
    data := PatientViewData{Patient: mapToViewModel(patient)}
    h.templates.ExecuteTemplate(w, "layout", data)
}
```

---

### 2. Consciência de Contexto HTMX

**Cada Handler deve verificar se a requisição é HTMX** (`HX-Request` header):

- **Se for HTMX**: Renderiza apenas o fragmento (bloco específico)
- **Se não for**: Renderiza a página completa com o `layout.html`

```go
func (h *Handler) ShowPatient(w http.ResponseWriter, r *http.Request) {
    // ... obter dados ...
    
    data := PatientViewData{Patient: patient}
    
    if r.Header.Get("HX-Request") == "true" {
        // Requisição via HTMX - retorna apenas o fragmento
        h.templates.ExecuteTemplate(w, "patient-content", data)
        return
    }
    
    // Requisição direta - retorna página completa
    h.templates.ExecuteTemplate(w, "layout", data)
}
```

**Benefícios:**
- Mesma endpoint serve ambos cenários
- Redução de duplicação de código
- Experiência SPA-like sem complexidade de JavaScript

---

### 3. Tipagem Forte (ViewModels)

**Nunca passe entidades de domínio diretamente para o template.**

Crie structs específicas de "ViewData" dentro do handler para garantir que:
- O template tenha exatamente o que precisa
- Entidades de domínio permaneçam isoladas
- Mudanças na UI não afetem o domínio (e vice-versa)

```go
// ❌ ERRADO: Entidade de domínio exposta ao template
type Patient struct {
    ID        string
    Name      string
    CreatedAt time.Time  // Template precisa formatar?
}
h.templates.ExecuteTemplate(w, "layout", patient)

// ✅ CORRETO: ViewModel específico para a view
type PatientViewModel struct {
    ID        string
    Name      string
    CreatedAt string  // Já formatado para exibição
}

type PatientViewData struct {
    Patient  *PatientViewModel
    Sessions []*SessionViewModel
    Insights []InsightViewModel
    Error    string
}

data := PatientViewData{
    Patient: &PatientViewModel{
        ID: patient.ID,
        Name: patient.Name,
        CreatedAt: patient.CreatedAt.Format("Jan 2006"),
    },
}
h.templates.ExecuteTemplate(w, "layout", data)
```

---

## Estrutura de Templates

### Hierarquia Modular

```
web/templates/
├── layout.html           # Esqueleto base (head, sidebar, scripts)
├── error-fragment.html   # Fragmento de erro para HTMX
├── patients.html         # Define "content" + "patients-content"
├── patient.html          # Define "content" + "patient-content"
├── session.html          # Define "content" + "session-content"
└── session_new.html      # Define "content" + "new-session-form"
```

### Padrão de Definição

Cada template de página define **dois blocos**:

```html
{{/* patient.html */}}

{{/* 1. Full-page rendering - injeta fragmento no layout */}}
{{define "content"}}
{{template "patient-content" .}}
{{end}}

{{/* 2. HTMX fragment - renderizável standalone */}}
{{define "patient-content"}}
{{if .Error}}
    <div class="error">{{.Error}}</div>
{{else if .Patient}}
    <!-- Conteúdo do paciente -->
{{end}}
{{end}}
```

### Layout Base

```html
{{/* layout.html */}}
<!DOCTYPE html>
<html>
<head>
    <!-- Meta tags, CSS, etc. -->
</head>
<body>
    <aside class="sidebar">
        <!-- Navegação, insights -->
    </aside>
    
    <main class="content">
        {{block "content" .}}Default content{{end}}
    </main>
    
    <script src="/js/htmx.min.js"></script>
</body>
</html>
```

---

## Tratamento de Erros

### Princípios

1. **Status codes apropriados** - 400, 404, 500 conforme necessário
2. **Fragmentos de erro para HTMX** - HTML amigável, não erro bruto
3. **Mensagens claras** - Úteis para usuário final

### Implementação

```go
func (h *Handler) renderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(statusCode)
    
    data := ErrorViewData{
        Message: message,
        RetryURL: "/patients",
    }
    
    if r.Header.Get("HX-Request") == "true" {
        // HTMX recebe fragmento estilizado
        h.templates.ExecuteTemplate(w, "error-fragment", data)
        return
    }
    
    // Página completa com layout
    h.templates.ExecuteTemplate(w, "layout", data)
}
```

### Template de Erro

```html
{{/* error-fragment.html */}}
{{define "error-fragment"}}
<div class="bg-red-50 border border-red-200 rounded-lg p-6">
    <p class="text-red-700 font-medium">{{.Message}}</p>
    <button onclick="window.location.reload()">Recarregar</button>
</div>
{{end}}
```

---

## Injeção de Dependência

### Interfaces por Handler

Cada handler define **interfaces específicas** com apenas os métodos que usa:

```go
type PatientService interface {
    GetPatientByID(ctx context.Context, id string) (*patient.Patient, error)
    ListPatients(ctx context.Context) ([]*patient.Patient, error)
    CreatePatient(ctx context.Context, input services.CreatePatientInput) (*patient.Patient, error)
}

type TemplateRenderer interface {
    ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error
}

type PatientHandler struct {
    patientService PatientService
    sessionService SessionService
    insightService InsightService
    templates      TemplateRenderer
}

func NewPatientHandler(
    patientService PatientService,
    sessionService SessionService,
    insightService InsightService,
    templates TemplateRenderer,
) *PatientHandler {
    return &PatientHandler{
        patientService: patientService,
        sessionService: sessionService,
        insightService: insightService,
        templates:      templates,
    }
}
```

**Benefícios:**
- Facilita testes unitários (mocking simples)
- Reduz acoplamento
- Segregação de interfaces (ISP do SOLID)

---

## Fluxo Completo

### Exemplo: Criar Paciente

```
1. Request HTTP POST /patients
   ↓
2. Handler decodifica form data
   ↓
3. Handler chama PatientService.CreatePatient()
   ↓
4. Service valida e cria entidade de domínio
   ↓
5. Service persiste via Repository
   ↓
6. Handler mapeia entidade → PatientViewModel
   ↓
7. Handler monta PatientViewData
   ↓
8. Handler verifica HX-Request
   ↓
9. Handler renderiza template (fragmento ou full-page)
   ↓
10. Response HTML para browser/HTMX
```

---

## Checklist de Implementação

Ao criar um novo handler, verifique:

- [ ] Handler não contém lógica de negócio
- [ ] ViewModels criados para todos dados da view
- [ ] Verificação de `HX-Request` implementada
- [ ] Fragmento HTMX definido no template
- [ ] Tratamento de erros com `renderError()`
- [ ] Interfaces definidas para dependências
- [ ] Injeção de dependência via construtor
- [ ] Status codes apropriados
- [ ] Mensagens de erro amigáveis

---

## Anti-Padrões a Evitar

### ❌ Lógica de Negócio no Handler

```go
// NÃO FAÇA ISSO
func (h *Handler) CreatePatient(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("name")
    
    // Validação de negócio no handler!
    if len(name) < 3 {
        http.Error(w, "Nome deve ter pelo menos 3 caracteres", http.StatusBadRequest)
        return
    }
    
    // Formatação de data no handler!
    date, _ := time.Parse("2006-01-02", r.FormValue("date"))
    
    // Persistência direta no handler!
    db.Exec("INSERT INTO patients...", name, date)
}
```

### ❌ Entidade de Domínio no Template

```go
// NÃO FAÇA ISSO
patient, _ := h.service.GetPatient(id)
h.templates.ExecuteTemplate(w, "layout", patient)  // Expõe domínio!
```

### ❌ Ignorar HX-Request

```go
// NÃO FAÇA ISSO
func (h *Handler) ShowPatient(w http.ResponseWriter, r *http.Request) {
    // Sempre retorna página completa, mesmo para HTMX!
    h.templates.ExecuteTemplate(w, "layout", data)
}
```

### ❌ Templates Genéricos Demais

```html
<!-- NÃO FAÇA ISSO -->
{{define "content"}}
<!-- Conteúdo específico do paciente -->
{{end}}

<!-- Use nomes específicos para fragments -->
{{define "patient-content"}}
<!-- Assim evita conflitos -->
{{end}}
```

---

## Referências

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://domainlanguage.com/ddd/)
- [HTMX Documentation](https://htmx.org/docs/)
- [Go Templates - Official Docs](https://pkg.go.dev/html/template)

---

**Última Atualização:** Março 2026  
**Status:** Em produção  
**Responsável:** Arquitetura de Sistemas
