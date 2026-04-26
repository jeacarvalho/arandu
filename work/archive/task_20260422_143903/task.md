# Task: Busca contextual global no prontuário
Requirement: REQ-07-04-02
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

O Arandu possui FTS5 indexando observações, intervenções e timeline, mas toda busca
existente é **scoped por paciente** — `SearchInHistory(ctx, patientID, query)`. Não há
busca global. O terapeuta não consegue digitar "ansiedade" e ver em quais pacientes e
sessões esse tema aparece.

Esta task implementa a busca global contextual: barra de busca no layout principal,
resultados com snippet destacado, paciente, data e tipo de registro.

### O que já existe (não recriar)

```go
// timeline.SearchResult — tipo pronto
type SearchResult struct {
    ID        string
    Type      EventType  // "observation" | "intervention"
    Date      time.Time
    Content   string
    Snippet   string     // snippet FTS5 com tags <b>termo</b>
    SessionID string
    PatientID string
}

// TimelineServiceContext — já tem SearchInHistory por paciente
func (s *TimelineServiceContext) SearchInHistory(ctx context.Context, patientID, query string) ([]*timeline.SearchResult, error)

// PatientService — busca pacientes por nome
func (s *PatientService) SearchPatients(ctx context.Context, query string, limit, offset int) ([]*patient.Patient, error)

// TimelineHandler.SearchPatientHistory — busca dentro de 1 paciente
func (h *TimelineHandler) SearchPatientHistory(w http.ResponseWriter, r *http.Request)

// Rota existente (NÃO modificar):
// /history/search — scoped por patient
```

---

## O que implementar

### 1. Novo método no TimelineRepository — busca global

**Arquivo:** `internal/infrastructure/repository/sqlite/timeline_repository.go`

Adicionar método que faz FTS5 sem filtro de `patient_id`:

```go
func (r *TimelineRepository) SearchGlobal(ctx context.Context, query string, limit int) ([]*timeline.SearchResult, error)
```

A query SQL deve usar `snippet()` do FTS5 para gerar o campo `Snippet`:

```sql
-- Adapte ao schema real do FTS. Exemplo de padrão a seguir:
SELECT
    e.id,
    e.type,
    e.date,
    e.content,
    snippet(timeline_fts, 0, '<b>', '</b>', '…', 10) AS snippet,
    e.session_id,
    e.patient_id
FROM timeline_fts
JOIN timeline_entries e ON timeline_fts.rowid = e.rowid
WHERE timeline_fts MATCH ?
ORDER BY rank
LIMIT ?
```

> **Atenção:** inspecione o schema FTS real (nome da tabela virtual, colunas indexadas)
> antes de escrever a query. Siga exatamente o padrão de `SearchInHistory` do mesmo arquivo.

O `ContextAwareTimelineRepository` em `context_wrapper.go` expõe `SearchInHistory` —
adicione `SearchGlobal` com o mesmo padrão de delegate.

### 2. Novo método no TimelineService

**Arquivo:** `internal/application/services/timeline_service_context.go`

```go
// SearchGlobal busca em todos os registros clínicos do tenant, sem filtro de paciente.
// Enriquece cada resultado com o nome do paciente via PatientService.
func (s *TimelineServiceContext) SearchGlobal(ctx context.Context, query string) ([]*SearchGlobalResult, error)
```

`SearchGlobalResult` — novo tipo no mesmo arquivo (ou em `timeline_service.go`):

```go
type SearchGlobalResult struct {
    timeline.SearchResult           // embed — mantém ID, Type, Date, Snippet, SessionID, PatientID
    PatientName string
}
```

Lógica do método:
1. Chama `s.timelineRepo.SearchGlobal(ctx, query, 50)` — limite de 50 resultados
2. Para cada resultado único `PatientID`, chama `s.patientService.GetPatient(ctx, patientID)`
   para enriquecer o `PatientName` (use um map para evitar chamadas duplicadas)
3. Retorna `[]*SearchGlobalResult`

> **Verifique** se `TimelineServiceContext` já tem acesso a um `PatientService` ou se
> precisa receber via construtor. Se não tiver, injete como novo campo.

### 3. Novo handler `SearchHandler`

**Arquivo:** `internal/web/handlers/search_handler.go` (novo arquivo)

```go
type SearchHandler struct {
    timelineService TimelineSearchServiceInterface
}

type TimelineSearchServiceInterface interface {
    SearchGlobal(ctx context.Context, query string) ([]*services.SearchGlobalResult, error)
}

func NewSearchHandler(timelineService TimelineSearchServiceInterface) *SearchHandler

// Search handles GET /search?q=termo
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request)
```

Lógica do handler:
1. Extrai `q := r.URL.Query().Get("q")` — se vazio, renderiza página com resultados vazios
2. `q` menor que 2 caracteres: renderiza página com mensagem "Digite ao menos 2 caracteres"
3. Chama `h.timelineService.SearchGlobal(ctx, q)`
4. Mapeia para `SearchResultsViewModel`
5. Se `HX-Request == "true"`: renderiza fragmento `SearchResults(vm)`
6. Senão: renderiza página completa com layout shell

### 4. ViewModel e componente templ

**Arquivo:** `web/components/search/types.go` (novo)

```go
package search

type SearchResultItem struct {
    ID          string
    PatientID   string
    PatientName string
    SessionID   string
    Type        string    // "Observação" | "Intervenção"
    Date        string    // formatado: "22 de abril de 2026"
    Snippet     string    // HTML com <b>termo</b> — usar templ.Raw()
}

type SearchResultsViewModel struct {
    Query   string
    Results []SearchResultItem
    Total   int
}
```

**Arquivo:** `web/components/search/search_results.templ` (novo)

Layout da página:
- Header: `Busca: "{query}" — {N} resultado(s)`
- Se `Total == 0`: empty state "Nenhum resultado para '{query}'"
- Lista de resultados: cada item mostra
  - Nome do paciente (link para `/patient/{id}`) — `font-sans font-semibold`
  - Tipo + data — `text-xs text-neutral-400`
  - Snippet — `font-serif text-sm` com `templ.Raw(item.Snippet)` para renderizar o `<b>`
  - Link "Ver sessão →" para `/session/{sessionID}` se `SessionID != ""`

> **Snippet seguro:** o `Snippet` vem do FTS5 com tags `<b>` geradas pelo próprio banco —
> não é input do usuário. Usar `templ.Raw()` é seguro neste caso específico.
> Documente isso com um comentário no código.

### 5. Barra de busca no layout shell

**Arquivo:** `web/components/layout/shell.templ` (ou equivalente — verifique o arquivo real)

Adicionar input de busca no header/navbar:

```templ
<form
    hx-get="/search"
    hx-target="#main-content"
    hx-push-url="true"
    hx-trigger="submit, keyup changed delay:500ms from:find input"
    class="relative"
>
    <input
        type="search"
        name="q"
        placeholder="Buscar no prontuário…"
        class="w-48 h-8 pl-8 pr-3 text-sm bg-neutral-100 border-0 rounded-lg focus:ring-2 focus:ring-arandu-primary/30 focus:w-64 transition-all"
    />
    <i class="fas fa-search absolute left-2.5 top-2 text-xs text-neutral-400"></i>
</form>
```

> O `hx-target="#main-content"` deve apontar para o container principal do layout.
> Verifique o `id` real do container em `shell.templ` antes de usar.

### 6. Registrar rota em `main.go`

```go
searchHandler := handlers.NewSearchHandler(timelineServiceAdapter)
mux.HandleFunc("/search", searchHandler.Search)
```

---

## Checklist de implementação

- [ ] `TimelineRepository.SearchGlobal` implementado, segue schema FTS real
- [ ] `ContextAwareTimelineRepository.SearchGlobal` delegando corretamente
- [ ] `SearchGlobalResult` struct criada com embed de `timeline.SearchResult`
- [ ] `TimelineServiceContext.SearchGlobal` enriquece com `PatientName`
- [ ] `SearchHandler` criado em arquivo novo
- [ ] Rota `GET /search` registrada em `main.go`
- [ ] `SearchResultsViewModel` e `SearchResultItem` em `web/components/search/types.go`
- [ ] `search_results.templ` renderiza lista com snippet, paciente e link para sessão
- [ ] `templ.Raw(snippet)` usado com comentário explicando por que é seguro
- [ ] Barra de busca adicionada ao shell layout
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build ./cmd/arandu/` sem erros
- [ ] `go test ./internal/application/services/...` continua passando

---

## Arquivos a criar

- `internal/web/handlers/search_handler.go`
- `web/components/search/types.go`
- `web/components/search/search_results.templ`

## Arquivos a modificar

- `internal/infrastructure/repository/sqlite/timeline_repository.go` — `SearchGlobal`
- `internal/infrastructure/repository/sqlite/context_wrapper.go` — delegate `SearchGlobal`
- `internal/application/services/timeline_service_context.go` — `SearchGlobal` + `SearchGlobalResult`
- `cmd/arandu/main.go` — registrar `searchHandler`
- `web/components/layout/shell.templ` (ou arquivo equivalente) — barra de busca

---

## 🔒 Privacidade

- [x] **Tier 2 — texto livre**: `Snippet` e `Content` contêm texto clínico livre.
  NUNCA logar o conteúdo dos resultados — apenas metadados (ID, tipo, data).
  O `Snippet` não sai do sistema; é renderizado apenas na UI do próprio terapeuta.
- [x] A busca é sempre scoped ao tenant do request (extração via context) —
  nunca acessa dados de outro psicólogo.

---

## 🚫 NÃO faça

- Não modifique `SearchInHistory` existente — ela continua sendo usada pela timeline de paciente
- Não use `html/template` — apenas `.templ`
- Não passe `SearchGlobalResult` direto ao template — mapeie para `SearchResultItem`
- Não logar `Snippet` ou `Content` nos handlers
- Não crie migration — as tabelas FTS5 já existem
- Não implemente paginação nesta task — limite fixo de 50 resultados é suficiente

---

## 📎 Padrão de referência

- `SearchInHistory` em `timeline_repository.go` — siga exatamente para `SearchGlobal` (mesma estrutura SQL, mesmo scan de rows)
- `ContextAwareTimelineRepository` em `context_wrapper.go` — siga o padrão de delegate para expor `SearchGlobal`
- Handler `PatientHandler.Search` em `patient_handler.go` — padrão de extração de `q` e resposta HTMX
- `web/components/patient/search.templ` — referência visual para lista de resultados (adapte para snippets)
