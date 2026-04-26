# ARANDU — Prompt de Implementação
## Busca Híbrida na Topbar: Pacientes por Nome + Registros Clínicos

---

### 🎯 Objetivo

A busca da topbar (`/search`) deve retornar **dois tipos de resultado** em resposta à mesma query:

1. **Pacientes** cujo nome contenha o termo buscado (busca por substring, case-insensitive)
2. **Registros clínicos** (observações e intervenções) cujo conteúdo contenha o termo (FTS5, já existente)

A UI deve apresentar ambas as seções separadas, sem recarregar a página (HTMX fragment swap em `#main-content`).

---

### 🏗️ Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX · DaisyUI v4 · SQLite (database-per-tenant)

**Multi-tenancy**: handler DEVE obter DB via `tenant.TenantDB(r.Context())`.

**HTMX**: o form da topbar usa `hx-get="/search"`, `hx-target="#main-content"`, `hx-push-url="true"`, `hx-trigger="submit, keyup changed delay:500ms from:find input"`. O handler detecta `HX-Request: true` para retornar fragmento ou página completa.

---

### 📦 Mudanças no Domínio/Serviço

Nenhuma mudança de domínio. O `SearchHandler` precisa receber um segundo serviço (`PatientSearchServiceInterface`) além do já existente `TimelineSearchServiceInterface`.

---

### 🔧 Mudanças necessárias

#### 1. `internal/web/handlers/search_handler.go`

Adicionar campo e interface para busca de pacientes:

```go
// Nova interface (adicionar ao arquivo)
type PatientSearchServiceInterface interface {
    SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error)
}

// Struct atualizada
type SearchHandler struct {
    timelineService TimelineSearchServiceInterface
    patientService  PatientSearchServiceInterface  // ← NOVO
}

// Construtor atualizado
func NewSearchHandler(
    timelineService TimelineSearchServiceInterface,
    patientService  PatientSearchServiceInterface,  // ← NOVO
) *SearchHandler {
    return &SearchHandler{
        timelineService: timelineService,
        patientService:  patientService,
    }
}
```

No método `Search`, após obter os resultados de timeline, adicionar busca de pacientes:

```go
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")

    vm := searchComponents.SearchResultsViewModel{
        Query:    query,
        Results:  []searchComponents.SearchResultItem{},
        Patients: []searchComponents.PatientSearchResultItem{}, // ← NOVO campo
        Total:    0,
    }

    if len(query) < 2 {
        h.renderSearchResults(w, r, vm)
        return
    }

    // Busca de pacientes por nome
    if h.patientService != nil {
        patients, err := h.patientService.SearchPatients(r.Context(), query, 10, 0)
        if err == nil {
            for _, p := range patients {
                vm.Patients = append(vm.Patients, searchComponents.PatientSearchResultItem{
                    ID:   p.ID,
                    Name: p.Name,
                })
            }
        }
    }

    // Busca de registros clínicos (FTS5) — código existente, mantém igual
    results, err := h.timelineService.SearchGlobal(r.Context(), query)
    if err == nil && results != nil {
        // ... mapeamento existente, sem alteração
    }

    vm.Total = len(vm.Patients) + len(vm.Results)
    h.renderSearchResults(w, r, vm)
}
```

#### 2. `web/components/search/types.go`

Adicionar novo tipo e campo no ViewModel:

```go
// Novo struct para resultado de paciente
type PatientSearchResultItem struct {
    ID   string
    Name string
}

// SearchResultsViewModel — adicionar campo Patients
type SearchResultsViewModel struct {
    Query    string
    Results  []SearchResultItem         // registros clínicos (existente)
    Patients []PatientSearchResultItem  // ← NOVO
    Total    int
}
```

#### 3. `web/components/search/search_results.templ`

Atualizar os dois templates (`SearchResults` e `SearchResultsPage`) para renderizar a seção de pacientes antes da seção de registros clínicos.

**Estrutura da seção de pacientes** (DaisyUI, dentro do componente):

```templ
if len(data.Patients) > 0 {
    <div class="mb-4">
        <div class="text-xs font-semibold uppercase tracking-wider text-base-content/50 px-2 mb-2">
            Pacientes
        </div>
        for _, p := range data.Patients {
            <a
                href={ templ.URL("/patients/" + p.ID) }
                class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
            >
                <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                    <i class="fas fa-user text-primary text-xs"></i>
                </div>
                <span class="text-sm font-medium">{ p.Name }</span>
            </a>
        }
    </div>
}
```

Se houver tanto pacientes quanto registros clínicos, adicionar um separador entre as seções:

```templ
if len(data.Patients) > 0 && len(data.Results) > 0 {
    <div class="divider my-2 text-xs text-base-content/40">Registros Clínicos</div>
}
```

A seção de registros clínicos (existente) continua como está.

Se ambas as seções estiverem vazias e `data.Query != ""`:

```templ
if len(data.Patients) == 0 && len(data.Results) == 0 && data.Query != "" {
    <p class="text-center text-sm text-base-content/50 py-6">
        Nenhum resultado para "{ data.Query }"
    </p>
}
```

#### 4. `cmd/arandu/main.go`

Atualizar a chamada do construtor (linha onde `NewSearchHandler` é chamado):

```go
// Antes:
searchHandler := handlers.NewSearchHandler(timelineServiceAdapter)

// Depois:
searchHandler := handlers.NewSearchHandler(timelineServiceAdapter, patientServiceAdapter)
```

O `patientServiceAdapter` já existe no arquivo (é usado por `patientHandler`). Apenas passá-lo como segundo argumento.

#### 5. `tests/e2e/playwright_runner_test.go`

Atualizar o router do teste para usar o `SearchHandler` real (com ambos os serviços):

```go
// Substituir o bloco do /search handler:
// Antes:
mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    if q != "" {
        patientHandler.Search(w, r)
    } else {
        dashboardHandler.Show(w, r)
    }
})

// Depois:
searchHandler := handlers.NewSearchHandler(timelineServiceAdapter, patientServiceAdapter)
mux.HandleFunc("/search", searchHandler.Search)
```

O `patientServiceAdapter` já está disponível no escopo de `setupRouterPlaywright`.

#### 6. `tests/e2e/playwright/clinical_workflow.spec.ts` — Teste 04

Corrigir para testar o comportamento real da busca híbrida:

```typescript
test('04 · Busca na topbar encontra o paciente', async () => {
    await page.goto(BASE_URL + '/dashboard');
    const searchInput = page.locator('#shell-patient-search');

    // Busca pelo primeiro nome do paciente
    await searchInput.fill(patientName.split(' ')[0]);

    // Aguarda a resposta HTMX (keyup delay:500ms)
    const [searchResponse] = await Promise.all([
        page.waitForResponse(resp => resp.url().includes('/search') && resp.status() === 200),
        searchInput.press('Enter'),
    ]);
    expect(searchResponse.status()).toBe(200);

    // O resultado deve conter o nome do paciente na seção "Pacientes"
    await expect(page.locator('#main-content').getByText(patientName)).toBeVisible();
});
```

---

### 📁 Arquivos a modificar

| Arquivo | Mudança |
|---------|---------|
| `internal/web/handlers/search_handler.go` | Nova interface + campo + lógica de busca de pacientes |
| `web/components/search/types.go` | Novo tipo `PatientSearchResultItem` + campo `Patients` no ViewModel |
| `web/components/search/search_results.templ` | Seção de pacientes antes dos registros clínicos |
| `cmd/arandu/main.go` | Passar `patientServiceAdapter` para `NewSearchHandler` |
| `tests/e2e/playwright_runner_test.go` | Usar `SearchHandler` real com ambos os serviços |
| `tests/e2e/playwright/clinical_workflow.spec.ts` | Corrigir teste 04 para busca híbrida |

**NÃO criar novos arquivos** — todas as mudanças são em arquivos existentes.

---

### 🔒 Privacidade

- [ ] **Tier 2 — texto livre**: conteúdo de observações/intervenções não sai do sistema. A busca FTS5 ocorre no tenant DB — sem envio externo.
- [ ] **Nomes de pacientes** são Tier 1 (PII): não logar os resultados da busca.

---

### ✅ Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build ./cmd/arandu/` sem erros

**Comportamento**
- [ ] CA01: Digitar nome de paciente → seção "Pacientes" aparece com link clicável para `/patients/{id}`
- [ ] CA02: Digitar termo de observação existente → seção "Registros Clínicos" aparece com snippet e nome do paciente
- [ ] CA03: Query válida que encontra ambos → duas seções separadas por divisor
- [ ] CA04: Query que não encontra nada → mensagem "Nenhum resultado para..."
- [ ] CA05: Query com 1 caractere → resultados vazios (sem request desnecessário — min 2 chars)
- [ ] CA06: Clicar no paciente na seção de resultados → navega para `/patients/{id}`
- [ ] CA07: HTMX atualiza `#main-content` sem recarregar a página

**Testes automatizados**
- [ ] `go test ./tests/e2e/ -run TestPlaywrightClinicalWorkflow -count=1` — 10 passed
- [ ] Teste 04 verificar nome do paciente em `#main-content` após busca (não no corpo inteiro da página)

**Integridade**
- [ ] `go build ./...` compila sem erros

---

### 🚫 NÃO faça

- Não criar um novo componente templ separado para resultados de paciente — estender o existente `search_results.templ`
- Não alterar o form da topbar em `shell_layout.templ` — ele já está correto
- Não alterar o endpoint `/patients/search` — ele continua funcionando como está para outros usos
- Não logar o conteúdo das buscas ou os nomes retornados
- Não usar `html/template` — apenas `.templ`
- Não remover a seção de registros clínicos — manter ambas as seções

---

### 📎 Padrão de referência

Siga `internal/web/handlers/search_handler.go` para estrutura do handler.
Siga `web/components/search/search_results.templ` para o padrão de template existente.
A busca de pacientes usa `patientService.SearchPatients(ctx, query, 10, 0)` — mesmo método que `PatientHandler.Search` já usa.
