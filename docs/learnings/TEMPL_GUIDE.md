# 🎨 Guia Templ - Arandu

**Última atualização:** $(date +"%d de %B de %Y")
**Foco:** Framework Templ para componentes UI

> 📚 **Consulte também:** [MASTER_LEARNINGS.md](./MASTER_LEARNINGS.md) | [ARCHITECTURE_PATTERNS.md](./ARCHITECTURE_PATTERNS.md)

---

## 📋 Índice

1. [Introdução ao Templ](#introdução-ao-templ)
2. [Instalação e Setup](#instalação-e-setup)
3. [Estrutura de Componentes](#estrutura-de-componentes)
4. [Sintaxe e Funcionalidades](#sintaxe-e-funcionalidades)
5. [Erros Comuns e Soluções](#erros-comuns-e-soluções)
6. [Migração de html/template](#migração-de-htmltemplate)
7. [Boas Práticas](#boas-práticas)
8. [Referências](#referências)

---

## Introdução ao Templ

### O que é Templ?

Templ é um framework de templates para Go que gera código Go em tempo de compilação. Ele oferece:

- **Tipagem forte:** Erros de compilação em vez de erros em runtime
- **Performance:** Código Go nativo otimizado
- **Segurança:** Proteção automática contra XSS
- **Componentes reutilizáveis:** Composição de UI

### Por que Templ no Arandu?

1. **Performance:** Geração de código em compile-time
2. **Segurança:** Escape automático de HTML
3. **Manutenibilidade:** Componentes organizados e tipados
4. **Integração:** Bom suporte a HTMX

---

## Instalação e Setup

### Instalação do CLI

```bash
# Instalar templ CLI
go install github.com/a-h/templ/cmd/templ@latest

# Verificar instalação
templ --version
# Deve retornar: templ v0.3.1001 (ou similar)
```

### Fluxo de Desenvolvimento

1. **Criar componente** `.templ`
2. **Gerar código** com `templ generate`
3. **Importar e usar** no handler
4. **Build nativo** com `go build`

### Comandos Essenciais

```bash
# Gerar código Go a partir de .templ
templ generate

# Gerar e observar mudanças (desenvolvimento)
templ generate --watch

# Verificar sintaxe
templ fmt

# Limpar arquivos gerados
templ generate --clean
```

---

## Estrutura de Componentes

### Organização do Projeto

```
web/components/
├── layout/
│   ├── layout.templ          # Layout base com sidebar
│   └── layout_templ.go       # Gerado automaticamente
├── patient/
│   ├── list.templ           # Lista de pacientes
│   ├── detail.templ         # Detalhe do paciente
│   ├── form.templ           # Formulário
│   ├── search.templ         # Busca
│   └── patient_templ.go     # Gerado automaticamente
├── session/
│   ├── detail.templ         # Detalhe da sessão
│   ├── list.templ           # Lista de sessões
│   ├── form.templ           # Formulário
│   └── session_templ.go     # Gerado automaticamente
└── dashboard/
    ├── dashboard.templ      # Dashboard clínico
    └── dashboard_templ.go   # Gerado automaticamente
```

### Estrutura de um Componente

```templ
// web/components/patient/detail.templ
package patient

// Importações são gerenciadas AUTOMATICAMENTE
// NUNCA adicione import "github.com/a-h/templ" manualmente

// Struct de dados do componente
type PatientDetailData struct {
    ID        string
    Name      string
    Notes     string
    CreatedAt string
    UpdatedAt string
}

// Componente principal
templ PatientDetail(data PatientDetailData) {
    <div id="patient-detail" class="p-6">
        <header class="mb-6">
            <h1 class="text-3xl font-bold text-gray-900">{ data.Name }</h1>
            <div class="mt-2 text-sm text-gray-500">
                <span>ID: { data.ID }</span>
                <span class="mx-2">•</span>
                <span>Cadastrado: { data.CreatedAt }</span>
                <span class="mx-2">•</span>
                <span>Atualizado: { data.UpdatedAt }</span>
            </div>
        </header>
        
        @if data.Notes != "" {
            <section class="mt-6">
                <h2 class="text-lg font-semibold text-gray-800 mb-2">Observações</h2>
                <div class="font-clinical text-gray-700 bg-gray-50 p-4 rounded-lg">
                    { data.Notes }
                </div>
            </section>
        }
    </div>
}

// Componente auxiliar (privado, não exportado)
templ patientNotes(notes string) {
    <div class="font-clinical text-gray-700">
        { notes }
    </div>
}
```

---

## Sintaxe e Funcionalidades

### Interpolação de Texto

```templ
// Texto simples
<p>{ userName }</p>

// Texto com escape automático (seguro contra XSS)
<p>{ userInput }</p>  // Escape automático

// Texto sem escape (use com cuidado!)
<p>{ templ.RawHTML(trustedHTML) }</p>
```

### Condicionais

```templ
// if simples
@if user.IsAdmin {
    <button>Administrar</button>
}

// if-else
@if patient.HasNotes {
    <div class="notes">{ patient.Notes }</div>
} else {
    <div class="text-gray-500">Sem observações</div>
}

// if-else if
@if score >= 90 {
    <span class="text-green-600">Excelente</span>
} else if score >= 70 {
    <span class="text-yellow-600">Bom</span>
} else {
    <span class="text-red-600">Precisa melhorar</span>
}
```

### Loops

```templ
// Loop sobre slice
<ul>
    @for _, patient := range patients {
        <li>{ patient.Name }</li>
    }
</ul>

// Loop com índice
<ol>
    @for i, medication := range medications {
        <li value={ i+1 }>{ medication.Name } - { medication.Dosage }</li>
    }
</ol>

// Loop vazio
@if len(sessions) == 0 {
    <div class="empty-state">
        Nenhuma sessão registrada
    </div>
}
```

### Atributos Dinâmicos

```templ
// Atributos condicionais
<div class={ "active" if isActive else "inactive" }>
    Conteúdo
</div>

// Múltiplas classes
<div class={ "base-class " + (additionalClass if hasAdditional else "") }>
    Conteúdo
</div>

// Atributos booleanos
<input type="checkbox" checked={ isChecked } disabled={ isDisabled }>

// Atributos com valores dinâmicos
<a href={ "/patient/" + patientID }>Ver paciente</a>
```

### CSS no Componente

```templ
// CSS scoped ao componente
css patientCard() {
    background-color: white;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
    margin-bottom: 1rem;
}

css clinicalText() {
    font-family: 'Source Serif 4', serif;
    font-size: 1.125rem;
    line-height: 1.75;
    color: #1F2937;
}

templ PatientCard(data PatientData) {
    <div class={ patientCard() }>
        <h2 class="text-xl font-bold">{ data.Name }</h2>
        <div class={ clinicalText() }>
            { data.ClinicalNotes }
        </div>
    </div>
}
```

### Slots e Composição

```templ
// Componente com slot
templ Card(title string) {
    <div class="card">
        <h3 class="card-title">{ title }</h3>
        <div class="card-content">
            { children... }
        </div>
    </div>
}

// Uso com slot
templ PatientView(patient Patient) {
    @Card("Informações do Paciente") {
        <p>Nome: { patient.Name }</p>
        <p>Idade: { patient.Age }</p>
        <p>Diagnóstico: { patient.Diagnosis }</p>
    }
}
```

---

## Erros Comuns e Soluções

### ❌ Erro: "templ redeclared in this block"

**Sintoma:**
```
templ redeclared in this block
"github.com/a-h/templ" imported and not used
```

**Causa:** Adicionou `import "github.com/a-h/templ"` manualmente no arquivo `.templ`.

**Solução:** Remova TODAS as importações manuais. O `templ generate` adiciona automaticamente.

**Arquivo CORRETO:**
```templ
package patient

// NENHUMA importação manual aqui

type PatientData struct {
    Name string
}

templ PatientComponent(data PatientData) {
    <div>{ data.Name }</div>
}
```

### ❌ Erro: Componente não encontrado

**Sintoma:** `undefined: PatientComponent` ao compilar.

**Causa:** Não executou `templ generate` após criar/modificar `.templ`.

**Solução:**
```bash
# Gerar código Go
templ generate

# Ou em modo watch (desenvolvimento)
templ generate --watch
```

### ❌ Erro: Build constraints

**Sintoma:** Não compila com `//go:build templ`.

**Causa:** Usando build constraints desnecessárias.

**Solução:** Remova `//go:build templ`. Templ gera código Go nativo que deve compilar sem flags especiais.

### ❌ Problema: CSS não aplicado

**Causa:** CSS definido no componente não está sendo aplicado.

**Solução:** Certifique-se de usar a função CSS corretamente:

```templ
// DEFINIÇÃO
css cardStyle() {
    background: white;
    padding: 1rem;
}

// USO CORRETO
<div class={ cardStyle() }>
    Conteúdo
</div>

// USO INCORRETO (não funciona)
<div class="cardStyle">
    Conteúdo
</div>
```

### ❌ Problema: HTMX não funciona

**Causa:** Componente Templ retornando HTML inválido para HTMX.

**Solução:** Verifique se está retornando apenas o fragmento, não a página completa:

```go
// NO HANDLER
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
    data := getData()
    
    // Componente Templ
    component := patient.PatientDetail(data)
    
    if r.Header.Get("HX-Request") == "true" {
        // HTMX: apenas o componente
        component.Render(r.Context(), w)
    } else {
        // Full page: componente dentro do layout
        layout.Base("Título", component).Render(r.Context(), w)
    }
}
```

---

## Migração de html/template

### Diferenças Principais

| html/template | Templ |
|--------------|-------|
| Runtime parsing | Compile-time generation |
| `{{.Field}}` | `{ data.Field }` |
| `{{if .Cond}}` | `@if cond {` |
| `{{range .Items}}` | `@for _, item := range items {` |
| Sem tipagem | Tipagem forte |
| Templates em arquivos `.html` | Componentes em `.templ` |

### Passos da Migração

1. **Instalar Templ CLI**
2. **Criar estrutura de componentes**
3. **Converter templates gradualmente**
4. **Atualizar handlers para usar componentes**
5. **Remover templates antigos**

### Exemplo de Conversão

**html/template (antigo):**
```html
<!-- patient.html -->
{{define "content"}}
<div class="patient-detail">
    <h1>{{.Name}}</h1>
    <p>{{.Notes}}</p>
    {{if .HasSessions}}
    <ul>
        {{range .Sessions}}
        <li>{{.Date}}: {{.Summary}}</li>
        {{end}}
    </ul>
    {{end}}
</div>
{{end}}
```

**Templ (novo):**
```templ
// web/components/patient/detail.templ
package patient

type PatientDetailData struct {
    Name     string
    Notes    string
    Sessions []SessionData
}

type SessionData struct {
    Date    string
    Summary string
}

templ PatientDetail(data PatientDetailData) {
    <div class="patient-detail">
        <h1>{ data.Name }</h1>
        <p>{ data.Notes }</p>
        
        @if len(data.Sessions) > 0 {
            <ul>
                @for _, session := range data.Sessions {
                    <li>{ session.Date }: { session.Summary }</li>
                }
            </ul>
        }
    </div>
}
```

**Handler atualizado:**
```go
// Antigo
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
    data := getData()
    h.templates.ExecuteTemplate(w, "content", data)
}

// Novo
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
    data := getData()
    component := patient.PatientDetail(data)
    
    if r.Header.Get("HX-Request") == "true" {
        component.Render(r.Context(), w)
    } else {
        layout.Base("Paciente", component).Render(r.Context(), w)
    }
}
```

---

## Boas Práticas

### 1. Componentes Pequenos e Focados

**RUIM:** Componente com 500+ linhas fazendo tudo.

**BOM:** Componentes pequenos (< 200 linhas) com responsabilidade única.

### 2. Dados via Props, não Global

**RUIM:** Acessar variáveis globais no componente.

**BOM:** Receber todos os dados via struct de props.

### 3. CSS Scoped no Componente

**RUIM:** CSS global que afeta outros componentes.

**BOM:** CSS definido no próprio componente com `css nome() {}`.

### 4. Testabilidade

**RUIM:** Lógica complexa no componente.

**BOM:** Lógica no handler/service, componente apenas renderiza.

### 5. Nomenclatura Consistente

- **Componentes:** `PascalCase` (ex: `PatientDetail`)
- **Arquivos:** `snake_case.templ` (ex: `patient_detail.templ`)
- **Props structs:** `PascalCaseData` (ex: `PatientDetailData`)

### 6. Documentação no Componente

```templ
// PatientDetail renderiza a página de detalhe do paciente.
//
// Props:
//   - ID: Identificador único do paciente
//   - Name: Nome completo do paciente
//   - Notes: Observações clínicas (opcional)
//   - CreatedAt: Data de cadastro formatada
//
// Uso:
//   component := patient.PatientDetail(patient.PatientDetailData{...})
//   component.Render(ctx, w)
templ PatientDetail(data PatientDetailData) {
    // ...
}
```

### 7. Integração com HTMX

```templ
// Componente compatível com HTMX
templ PatientSearchResults(patients []PatientData, query string) {
    <div id="search-results">
        @if len(patients) == 0 {
            <p class="text-gray-500">
                Nenhum paciente encontrado para "{ query }"
            </p>
        } else {
            <ul class="divide-y divide-gray-200">
                @for _, patient := range patients {
                    <li class="py-3">
                        <a 
                          href={ "/patient/" + patient.ID }
                          hx-get={ "/patient/" + patient.ID }
                          hx-target="#main-content"
                          hx-push-url="true"
                          class="text-blue-600 hover:text-blue-800">
                            { patient.Name }
                        </a>
                    </li>
                }
            </ul>
        }
    </div>
}
```

---

## Referências

### Arquivos de Exemplo no Projeto

1. **`web/components/patient/detail.templ`** - Componente completo
2. **`web/components/layout/layout.templ`** - Layout base
3. **`web/components/session/form.templ`** - Formulário com validação

### Documentação Oficial

- [Documentação do Templ](https://templ.guide)
- [GitHub do Templ](https://github.com/a-h/templ)
- [Exemplos oficiais](https://github.com/a-h/templ-examples)

### Ferramentas Relacionadas

- **`templ generate --watch`** - Desenvolvimento com hot reload
- **`templ fmt`** - Formatação automática
- **VS Code Extension** - Syntax highlighting para `.templ`

### Troubleshooting

1. **Não esqueça de rodar `templ generate`**
2. **Nunca adicione imports manualmente**
3. **Verifique se o componente está sendo importado corretamente**
4. **Teste com `go build ./...` após mudanças**

---

*Guia criado como parte da refatoração do Arandu - Baseado em aprendizados de `REQ-01-02-02.md` e `task_20260315_154905.md`*