# ARANDU — Prompt de Implementação
## Biblioteca Clínica — Prontuários (`/notes`)

---

### 🎯 Objetivo

Criar a página **Prontuários** — uma biblioteca clínica global acessível via sidebar,
com split-view: lista de todos os pacientes à esquerda e painel de prontuário à direita.
Design idêntico a `design_handoff_arandu_redesign/page_notes.jsx`.

Esta é uma página **nova** — não existe nenhuma rota `/notes` hoje.

---

### 🏗️ Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX · SQLite (database-per-tenant)

**Multi-tenancy**: handler DEVE obter a conexão assim:
```go
db, err := tenant.TenantDB(r.Context())
if err != nil { http.Error(w, "unauthorized", 401); return }
```

**Serviços disponíveis** (já existem, só injetar no novo handler):
```go
patientService.ListPatients(ctx) ([]*patient.Patient, error)
patientService.GetPatientByID(ctx, id) (*patient.Patient, error)
timelineService.GetPatientTimeline(ctx, patientID string, filter *timeline.EventType, limit, offset int) (timeline.Timeline, error)
```

**HTMX**: o painel direito atualiza via HTMX quando o usuário clica em um paciente.
Verificar `r.Header.Get("HX-Request")` no handler de detalhe.

---

### 📦 ViewModels — criar em `web/components/notes/types.go`

```go
package notes

import "time"

// NotesLibraryViewModel — página completa
type NotesLibraryViewModel struct {
    Records       []NoteRecord
    FocusedRecord NoteRecordDetail
    FilterTags    []string // tags únicas dos pacientes (patient.Tag)
    TotalRecords  int
    TotalEvents   int      // soma de todos os eventos como proxy de "páginas"
    CurrentFilter string   // "all" ou valor de tag
}

// NoteRecord — item da lista esquerda
type NoteRecord struct {
    PatientID   string
    PatientName string
    Initials    string   // "AB" de "André Barbosa"
    RecordID    string   // "PR-" + últimos 4 chars do PatientID
    EventCount  int      // proxy para "páginas"
    LastUpdate  string   // patient.UpdatedAt formatado "02/01/2006"
    Tags        []string // patient.Tag (pode ser string vazia → ignorar)
    Status      string   // "EM ACOMPANHAMENTO" (fixo por enquanto)
    IsFocused   bool
}

// NoteRecordDetail — painel direito
type NoteRecordDetail struct {
    PatientID    string
    PatientName  string
    RecordID     string
    ActiveTab    string   // "evolucao" | "anamnese" | "observacoes" | "intervencoes"
    Sections     []NoteSection
    // Conteúdo da aba "Evolução" (last session)
    LastEntryDate    string
    LastEntryTitle   string
    LastEntryContent string // texto da sessão
    LastEntryQuote   string // primeiro trecho de observation dentro da sessão
    // Contadores por aba
    SessionCount      int
    ObservationCount  int
    InterventionCount int
    TotalEvents       int
}

type NoteSection struct {
    Key    string // "evolucao" | "anamnese" | "observacoes" | "intervencoes"
    Label  string
    Count  int
    Active bool
}
```

---

### 🛣️ Rotas — adicionar em `cmd/arandu/main.go`

```go
// Prontuários (biblioteca clínica global)
} else if r.URL.Path == "/notes" && r.Method == "GET" {
    notesHandler.Index(w, r)
} else if strings.HasPrefix(r.URL.Path, "/notes/detail/") && r.Method == "GET" {
    notesHandler.Detail(w, r)
```

Registrar **antes** das rotas de patient para não entrar em conflito.

---

### 📁 Arquivos a criar/modificar

**Criar:**
- `internal/web/handlers/notes_handler.go` — handler principal
- `web/components/notes/notes_page.templ` — templates da página
- `web/components/notes/types.go` — ViewModels

**Modificar:**
- `cmd/arandu/main.go` — registrar rotas + instanciar `NotesHandler`
- `web/components/layout/shell_layout.templ` — adicionar "Prontuários" no nav global
- `web/static/css/style.css` — adicionar bloco CSS `NOTES LIBRARY SÁBIO`

---

### 🔧 Handler — `internal/web/handlers/notes_handler.go`

```go
type NotesHandler struct {
    patientService  PatientServiceInterface
    timelineService TimelineServiceInterface
}

func NewNotesHandler(ps PatientServiceInterface, ts TimelineServiceInterface) *NotesHandler

// Index — GET /notes → renderiza página completa
func (h *NotesHandler) Index(w http.ResponseWriter, r *http.Request) {
    // 1. Extrair ?patient= e ?filter= da query
    focusedID := r.URL.Query().Get("patient")
    filter    := r.URL.Query().Get("filter")  // tag value ou ""

    // 2. Listar todos os pacientes
    patients, err := h.patientService.ListPatients(ctx)

    // 3. Montar NoteRecord para cada paciente
    //    - RecordID = "PR-" + último 4 chars de patient.ID
    //    - Initials  = buildInitials(patient.Name)
    //    - LastUpdate = patient.UpdatedAt.Format("02/01/2006")
    //    - Tags = []string{patient.Tag} se não vazio
    //    - IsFocused = (patient.ID == focusedID) || (i == 0 && focusedID == "")

    // 4. Coletar FilterTags: todos os patient.Tag únicos (ignorar vazios)

    // 5. Para o paciente focused, buscar timeline:
    //    timelineService.GetPatientTimeline(ctx, focusedID, nil, 20, 0)
    //    e montar NoteRecordDetail

    // 6. Renderizar
    if r.Header.Get("HX-Request") == "true" {
        notes.NotesLibraryContent(vm).Render(ctx, w)
        return
    }
    notes.NotesLibraryPage(vm).Render(ctx, w)
}

// Detail — GET /notes/detail/{patientID}?tab=evolucao
// Retorna APENAS o painel direito (HTMX swap de #notes-detail)
func (h *NotesHandler) Detail(w http.ResponseWriter, r *http.Request) {
    patientID := extractIDFromPath(r.URL.Path, "/notes/detail/")
    tab := r.URL.Query().Get("tab")
    if tab == "" { tab = "evolucao" }

    // buscar paciente + timeline, montar NoteRecordDetail
    // renderizar notes.NotesDetailPanel(detail)
}
```

**buildInitials** — reutilizar função já existente em `timeline_handler.go` (mover para um pacote utilitário ou duplicar):
```go
func buildNotesInitials(name string) string {
    parts := strings.Fields(name)
    if len(parts) == 0 { return "?" }
    r0 := []rune(parts[0])
    if len(parts) == 1 { return strings.ToUpper(string(r0[:1])) }
    rN := []rune(parts[len(parts)-1])
    return strings.ToUpper(string(r0[:1]) + string(rN[:1]))
}
```

**montarNoteRecordDetail** — lógica:
```go
func buildNoteRecordDetail(p *patient.Patient, events timeline.Timeline, tab string) NoteRecordDetail {
    sessions      := events.FilterByType(timeline.EventTypeSession)
    observations  := events.FilterByType(timeline.EventTypeObservation)
    interventions := events.FilterByType(timeline.EventTypeIntervention)

    detail := NoteRecordDetail{
        PatientID:         p.ID,
        PatientName:       p.Name,
        RecordID:          buildRecordID(p.ID),
        ActiveTab:         tab,
        SessionCount:      len(sessions),
        ObservationCount:  len(observations),
        InterventionCount: len(interventions),
        TotalEvents:       len(events),
        Sections: []NoteSection{
            {Key: "evolucao",      Label: "Evolução",      Count: len(sessions),      Active: tab == "evolucao"},
            {Key: "anamnese",      Label: "Anamnese",      Count: 0,                  Active: tab == "anamnese"},
            {Key: "observacoes",   Label: "Observações",   Count: len(observations),  Active: tab == "observacoes"},
            {Key: "intervencoes",  Label: "Intervenções",  Count: len(interventions), Active: tab == "intervencoes"},
        },
    }

    // Conteúdo da aba ativa
    if tab == "evolucao" && len(sessions) > 0 {
        last := sessions[0]  // já vem em ordem desc
        detail.LastEntryDate  = last.Date.Format("02/01/2006")
        detail.LastEntryTitle = buildTimelineTitle(last)
        detail.LastEntryContent = truncateStr(last.Content, 400)
    }
    // Para anamnese, usar patient.Notes
    // Para observacoes/intervencoes, listar os primeiros 5 (Content truncado)

    return detail
}

func buildRecordID(patientID string) string {
    id := patientID
    if len(id) > 4 { id = id[len(id)-4:] }
    return "PR-" + strings.ToUpper(id)
}
```

---

### 🎨 Template — `web/components/notes/notes_page.templ`

```
templ NotesLibraryPage(vm NotesLibraryViewModel)
  → usa layout.Shell(ShellConfig{PageTitle: "Prontuários", ActivePage: "notes"}, NotesLibraryContent(vm))

templ NotesLibraryContent(vm NotesLibraryViewModel)
  → div id="notes-content" com hero + chips + split-view

templ NotesDetailPanel(detail NoteRecordDetail)
  → apenas o painel direito — retornado pelo Detail handler via HTMX
```

#### Estrutura do template (mapeado do JSX):

```html
<!-- Hero -->
<div class="sabio-nl-hero">
  <div class="sabio-nl-hero-left">
    <div class="sabio-nl-eyebrow">Biblioteca Clínica</div>
    <h1 class="sabio-nl-title">Prontuários<em class="sabio-nl-title-dot">.</em></h1>
    <p class="sabio-nl-subtitle">{ vm.TotalRecords } registros · { vm.TotalEvents } páginas indexadas</p>
  </div>
  <div class="sabio-nl-hero-actions">
    <button class="sabio-nl-btn-outline">Filtros</button>
    <button class="sabio-nl-btn-primary">+ Novo prontuário</button>
  </div>
</div>

<!-- Filter chips -->
<div class="sabio-nl-chips">
  <button class={ chipClass("all", vm.CurrentFilter) }
    hx-get="/notes?filter=all" hx-target="#notes-content" hx-swap="outerHTML" hx-push-url="true">
    Tudo
  </button>
  for _, tag := range vm.FilterTags {
    <button class={ chipClass(tag, vm.CurrentFilter) }
      hx-get={ "/notes?filter=" + tag } hx-target="#notes-content" hx-swap="outerHTML" hx-push-url="true">
      { tag }
    </button>
  }
</div>

<!-- Split view -->
<div class="sabio-nl-split">
  <!-- Lista esquerda -->
  <section class="sabio-nl-list-panel" id="notes-list">
    <!-- Search header -->
    <header class="sabio-nl-list-header">
      <input class="sabio-nl-search"
        placeholder="Buscar por nome, tag, trecho…"
        hx-get="/notes" hx-target="#notes-list" hx-swap="innerHTML"
        hx-trigger="input changed delay:300ms"
        name="q" />
      <span class="sabio-nl-list-count">{ vm.TotalRecords }</span>
    </header>
    <!-- Items -->
    <div class="sabio-nl-list-scroll">
      for _, rec := range vm.Records {
        @NoteRecordItem(rec)
      }
    </div>
  </section>

  <!-- Painel direito -->
  <div id="notes-detail">
    @NotesDetailPanel(vm.FocusedRecord)
  </div>
</div>
```

**NoteRecordItem** (botão da lista):
```html
<button class={ noteRecordItemClass(rec.IsFocused) }
  hx-get={ "/notes/detail/" + rec.PatientID + "?tab=evolucao" }
  hx-target="#notes-detail"
  hx-swap="innerHTML"
  aria-selected={ strconv.FormatBool(rec.IsFocused) }>
  if rec.IsFocused {
    <span class="sabio-nl-focus-bar"></span>
  }
  <!-- Avatar -->
  <div class="sabio-nl-avatar">{ rec.Initials }</div>
  <!-- Info -->
  <div class="sabio-nl-record-info">
    <div class="sabio-nl-record-name-row">
      <span class="sabio-nl-record-name">{ rec.PatientName }</span>
    </div>
    <div class="sabio-nl-record-meta">
      <span class="sabio-nl-record-id">{ rec.RecordID }</span>
      <span>·</span>
      <span>{ rec.EventCount }p</span>
      <span>·</span>
      <span>atualiz. { rec.LastUpdate }</span>
    </div>
    <div class="sabio-nl-record-tags">
      for _, tag := range rec.Tags {
        <span class="sabio-nl-tag">{ tag }</span>
      }
      <span class="sabio-nl-tag sabio-nl-tag-status">{ rec.Status }</span>
    </div>
  </div>
  <!-- Chevron -->
  <span class="sabio-nl-chevron">›</span>
</button>
```

**NotesDetailPanel**:
```html
<section class="sabio-nl-detail">
  <!-- Header -->
  <header class="sabio-nl-detail-header">
    <div>
      <div class="sabio-nl-detail-eyebrow">
        <span class="sabio-nl-detail-id">{ detail.RecordID }</span> · { detail.PatientName }
      </div>
      <h2 class="sabio-nl-detail-title">
        Prontuário <em class="sabio-nl-detail-title-em">completo</em>
      </h2>
    </div>
    <div class="sabio-nl-detail-actions">
      <button class="sabio-nl-icon-btn" title="Análise IA">✦</button>
      <a href={ templ.URL("/patient/" + detail.PatientID + "/history") } class="sabio-nl-icon-btn" title="Ver histórico">↗</a>
    </div>
  </header>

  <!-- Tabs -->
  <nav class="sabio-nl-tabs">
    for _, s := range detail.Sections {
      <button class={ tabClass(s.Active) }
        hx-get={ "/notes/detail/" + detail.PatientID + "?tab=" + s.Key }
        hx-target="#notes-detail" hx-swap="innerHTML">
        { s.Label }
        <span class="sabio-nl-tab-count">{ s.Count }p</span>
      </button>
    }
  </nav>

  <!-- Conteúdo editorial -->
  <div class="sabio-nl-detail-body" id="notes-detail-body">
    if detail.ActiveTab == "evolucao" && detail.LastEntryTitle != "" {
      <div class="sabio-nl-doc">
        <div class="sabio-nl-doc-eyebrow">Evolução · última entrada em { detail.LastEntryDate }</div>
        <h3 class="sabio-nl-doc-title">{ detail.LastEntryTitle }</h3>
        <div class="sabio-nl-doc-divider"></div>
        <p class="sabio-nl-doc-body">{ detail.LastEntryContent }</p>
        if detail.LastEntryQuote != "" {
          <blockquote class="sabio-nl-doc-quote">{ detail.LastEntryQuote }</blockquote>
        }
      </div>
    } else if detail.ActiveTab == "evolucao" {
      <p class="sabio-nl-empty">Nenhuma sessão registrada.</p>
    } else if detail.ActiveTab == "anamnese" {
      <!-- Renderizar patient.Notes ou mensagem vazia -->
      <p class="sabio-nl-empty">Anamnese em construção.</p>
    } else {
      <p class="sabio-nl-empty">Sem registros nesta seção.</p>
    }
  </div>

  <!-- Footer -->
  <footer class="sabio-nl-detail-footer">
    <span class="sabio-nl-detail-pages">{ detail.TotalEvents } eventos indexados</span>
    <div class="sabio-nl-detail-footer-actions">
      <a href={ templ.URL("/patient/" + detail.PatientID + "/history") } class="sabio-nl-footer-btn">
        Ver prontuário completo
      </a>
    </div>
  </footer>
</section>
```

---

### 🎨 CSS — adicionar ao final de `web/static/css/style.css`

```css
/* ============================================
   NOTES LIBRARY SÁBIO  (sabio-nl-*)
   ============================================ */

/* Hero */
.sabio-nl-hero {
  display: grid; grid-template-columns: 1fr auto;
  gap: 32px; align-items: flex-end;
  padding-bottom: 22px; border-bottom: 1px solid var(--line);
  margin-bottom: 16px;
}
.sabio-nl-eyebrow {
  font-size: 11px; letter-spacing: 1.6px; text-transform: uppercase;
  color: var(--ink-3); font-weight: 500; margin-bottom: 10px;
}
.sabio-nl-title {
  margin: 0; font-family: var(--font-serif);
  font-size: 40px; font-weight: 400; letter-spacing: -0.8px; line-height: 1;
  color: var(--ink);
}
.sabio-nl-title-dot { font-style: italic; color: var(--accent-deep); }
.sabio-nl-subtitle { margin: 8px 0 0; font-size: 14px; color: var(--ink-3); }
.sabio-nl-hero-actions { display: flex; gap: 8px; }

.sabio-nl-btn-outline {
  padding: 9px 14px; border-radius: 10px;
  background: var(--paper-2); border: 1px solid var(--line);
  color: var(--ink-2); font-size: 13px; cursor: pointer;
  display: flex; align-items: center; gap: 8px;
}
.sabio-nl-btn-primary {
  padding: 9px 14px; border-radius: 10px;
  background: var(--ink); border: 1px solid var(--ink);
  color: var(--paper); font-size: 13px; font-weight: 500; cursor: pointer;
  display: flex; align-items: center; gap: 8px;
}
.sabio-nl-btn-primary:hover { background: var(--ink-2); }

/* Filter chips */
.sabio-nl-chips { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 16px; }

.sabio-nl-chip {
  padding: 6px 14px; border-radius: 999px; font-size: 12px; font-weight: 500;
  background: transparent; border: 1px solid var(--line);
  color: var(--ink-3); cursor: pointer;
}
.sabio-nl-chip-active {
  padding: 6px 14px; border-radius: 999px; font-size: 12px; font-weight: 500;
  background: var(--ink); border: 1px solid var(--ink);
  color: var(--paper); cursor: pointer;
}

/* Split view */
.sabio-nl-split {
  display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1.3fr);
  gap: 18px; align-items: flex-start;
}

/* Lista esquerda */
.sabio-nl-list-panel {
  background: var(--paper-2); border: 1px solid var(--line);
  border-radius: var(--radius, 12px); overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,.06);
}
.sabio-nl-list-header {
  padding: 14px 18px; border-bottom: 1px solid var(--line);
  display: flex; align-items: center; gap: 10px;
}
.sabio-nl-search {
  flex: 1; background: transparent; border: 0; outline: none;
  color: var(--ink); font-size: 13px; font-family: inherit;
}
.sabio-nl-search::placeholder { color: var(--ink-4); }
.sabio-nl-list-count { font-size: 11px; color: var(--ink-4); }
.sabio-nl-list-scroll { max-height: 640px; overflow-y: auto; }

/* Item da lista */
.sabio-nl-record-item {
  width: 100%; text-align: left;
  display: grid; grid-template-columns: auto 1fr auto;
  gap: 12px; align-items: center;
  padding: 14px 18px;
  border-bottom: 1px dashed var(--line);
  background: transparent; border-left: 3px solid transparent;
  cursor: pointer; position: relative;
  transition: background .15s;
}
.sabio-nl-record-item:last-child { border-bottom: none; }
.sabio-nl-record-item:hover { background: color-mix(in oklab, var(--ink) 4%, transparent); }

.sabio-nl-record-item-focused {
  width: 100%; text-align: left;
  display: grid; grid-template-columns: auto 1fr auto;
  gap: 12px; align-items: center;
  padding: 14px 18px;
  border-bottom: 1px dashed var(--line);
  background: color-mix(in oklab, var(--accent) 8%, transparent);
  border-left: 3px solid var(--accent);
  cursor: pointer; position: relative;
}
.sabio-nl-record-item-focused:last-child { border-bottom: none; }

.sabio-nl-avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: linear-gradient(135deg, color-mix(in oklab, var(--accent) 70%, var(--paper)), var(--accent));
  color: var(--paper); display: flex; align-items: center; justify-content: center;
  font-family: var(--font-serif); font-size: 13px; font-weight: 500;
  flex-shrink: 0;
}

.sabio-nl-record-info { min-width: 0; }
.sabio-nl-record-name-row { display: flex; align-items: center; gap: 6px; margin-bottom: 2px; }
.sabio-nl-record-name { font-size: 14px; font-weight: 500; color: var(--ink); }
.sabio-nl-record-meta {
  font-size: 11.5px; color: var(--ink-3);
  display: flex; gap: 8px; align-items: center; flex-wrap: wrap;
}
.sabio-nl-record-id { font-family: var(--font-mono); }
.sabio-nl-record-tags { margin-top: 6px; display: flex; gap: 6px; flex-wrap: wrap; }
.sabio-nl-tag {
  font-size: 10.5px; font-weight: 500; letter-spacing: .2px; text-transform: uppercase;
  padding: 2px 8px; border-radius: 999px; line-height: 1;
  background: color-mix(in oklab, var(--ink) 6%, transparent);
  color: var(--ink-2); border: 1px solid var(--line);
}
.sabio-nl-tag-status {
  background: color-mix(in oklab, var(--accent) 10%, transparent);
  color: var(--accent-deep);
  border-color: color-mix(in oklab, var(--accent) 25%, transparent);
}
.sabio-nl-chevron { font-size: 18px; color: var(--ink-4); }

/* Painel direito */
.sabio-nl-detail {
  background: var(--paper-2); border: 1px solid var(--line);
  border-radius: var(--radius, 12px); overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,.06);
  display: flex; flex-direction: column;
}
.sabio-nl-detail-header {
  padding: 18px 22px; border-bottom: 1px solid var(--line);
  display: grid; grid-template-columns: 1fr auto; gap: 16px; align-items: center;
}
.sabio-nl-detail-eyebrow {
  font-size: 10.5px; letter-spacing: 1.4px; text-transform: uppercase;
  color: var(--ink-3); font-weight: 500; margin-bottom: 4px;
}
.sabio-nl-detail-id { font-family: var(--font-mono); }
.sabio-nl-detail-title {
  margin: 0; font-family: var(--font-serif);
  font-size: 24px; font-weight: 500; letter-spacing: -0.3px; color: var(--ink);
}
.sabio-nl-detail-title-em { font-style: italic; color: var(--accent-deep); }
.sabio-nl-detail-actions { display: flex; gap: 6px; }
.sabio-nl-icon-btn {
  background: transparent; border: 1px solid var(--line);
  padding: 8px; border-radius: 8px; color: var(--ink-2);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; text-decoration: none;
  font-size: 14px;
}
.sabio-nl-icon-btn:hover { background: var(--paper-2); }

/* Tabs */
.sabio-nl-tabs {
  display: flex; border-bottom: 1px solid var(--line);
  background: color-mix(in oklab, var(--paper) 60%, var(--paper-2));
}
.sabio-nl-tab {
  padding: 12px 18px; background: transparent; border: 0;
  border-bottom: 2px solid transparent;
  color: var(--ink-3); font-size: 13px; font-weight: 400;
  display: flex; align-items: center; gap: 8px; cursor: pointer;
}
.sabio-nl-tab-active {
  padding: 12px 18px; background: transparent; border: 0;
  border-bottom: 2px solid var(--accent);
  color: var(--ink); font-size: 13px; font-weight: 500;
  display: flex; align-items: center; gap: 8px; cursor: pointer;
}
.sabio-nl-tab-count { font-size: 10px; color: var(--ink-4); }

/* Conteúdo editorial */
.sabio-nl-detail-body {
  flex: 1; padding: 28px 36px 32px;
  overflow-y: auto; background: var(--paper);
  max-height: 600px;
}
.sabio-nl-doc { max-width: 620px; margin: 0 auto; }
.sabio-nl-doc-eyebrow {
  font-size: 10.5px; letter-spacing: 1.4px; text-transform: uppercase;
  color: var(--ink-3); font-weight: 500; margin-bottom: 8px;
}
.sabio-nl-doc-title {
  margin: 0; font-family: var(--font-serif);
  font-size: 28px; font-weight: 400; letter-spacing: -0.5px; line-height: 1.15;
  color: var(--ink);
}
.sabio-nl-doc-divider {
  margin: 18px 0 22px; height: 1px;
  background: repeating-linear-gradient(90deg, var(--line) 0, var(--line) 4px, transparent 4px, transparent 8px);
}
.sabio-nl-doc-body {
  font-family: var(--font-serif); font-size: 16px;
  line-height: 1.65; color: var(--ink-2); margin-top: 0;
}
.sabio-nl-doc-quote {
  margin: 22px 0; padding: 4px 0 4px 20px;
  border-left: 2px solid var(--accent);
  font-family: var(--font-serif); font-style: italic;
  font-size: 17px; line-height: 1.55; color: var(--ink-2);
}
.sabio-nl-empty {
  color: var(--ink-4); font-size: 14px; text-align: center;
  padding: 40px 0;
}

/* Footer */
.sabio-nl-detail-footer {
  padding: 12px 22px; border-top: 1px solid var(--line);
  display: flex; align-items: center; gap: 12px;
  background: color-mix(in oklab, var(--paper) 60%, var(--paper-2));
}
.sabio-nl-detail-pages { font-family: var(--font-mono); font-size: 11px; color: var(--ink-3); }
.sabio-nl-detail-footer-actions { margin-left: auto; display: flex; gap: 8px; }
.sabio-nl-footer-btn {
  padding: 7px 12px; border-radius: 8px; font-size: 12px; font-weight: 500;
  background: var(--paper); border: 1px solid var(--line); color: var(--ink-2);
  text-decoration: none; cursor: pointer;
}
.sabio-nl-footer-btn:hover { background: var(--paper-2); }
```

---

### 🔧 Funções auxiliares Templ (no `notes_page.templ`)

```go
func chipClass(tag, current string) string {
    if tag == current || (current == "" && tag == "all") {
        return "sabio-nl-chip-active"
    }
    return "sabio-nl-chip"
}

func noteRecordItemClass(focused bool) string {
    if focused { return "sabio-nl-record-item-focused" }
    return "sabio-nl-record-item"
}

func tabClass(active bool) string {
    if active { return "sabio-nl-tab-active" }
    return "sabio-nl-tab"
}
```

---

### 🔗 Sidebar — modificar `web/components/layout/shell_layout.templ`

Em `ShellDefaultNavItems`, adicionar "Prontuários" como 4ª opção:
```go
templ ShellDefaultNavItems(activePage string) {
    @ShellNavItem("/dashboard", "dashboard", activePage, "fa-chart-line", "Dashboard")
    @ShellNavItem("/patients", "patients", activePage, "fa-users", "Pacientes")
    @ShellNavItem("/agenda", "agenda", activePage, "fa-calendar", "Agenda")
    @ShellNavItem("/notes", "notes", activePage, "fa-book-medical", "Prontuários")  // ← ADICIONAR
}
```

---

### 🔧 Instanciar handler em `cmd/arandu/main.go`

```go
notesHandler := handlers.NewNotesHandler(patientService, timelineService)
```

Seguir o padrão dos outros handlers no `main.go`.

---

### 🔒 Privacidade

- `TimelineEvent.Content` e `Patient.Notes` — **Tier 2**: nunca logar, nunca enviar para IA.
  Renderizar apenas na tela do prontuário autenticado.

---

### ✅ Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Comportamento**
- [ ] Sidebar global mostra "Prontuários" como 4ª opção e ativa ao visitar `/notes`
- [ ] GET `/notes` renderiza hero + chips de filtro + split-view com lista e painel
- [ ] Lista exibe todos os pacientes com avatar, nome, ID, tags e status
- [ ] Paciente focado (primeiro por padrão) tem borda esquerda accent + fundo levemente colorido
- [ ] Clicar em outro paciente atualiza o painel direito via HTMX (sem reload)
- [ ] Painel direito mostra header com ID do paciente + "Prontuário completo"
- [ ] Tabs funcionam via HTMX — "Evolução" mostra conteúdo da última sessão
- [ ] GET `/notes/detail/{id}?tab=observacoes` retorna o painel direito sozinho (fragmento)
- [ ] Filter chips "Tudo / Ansiedade / Burnout..." filtram a lista (apenas os pacientes com aquele tag)
- [ ] Ícone "↗" no painel direito leva ao `/patient/{id}/history`
- [ ] Página funciona sem nenhum paciente (empty state: "Nenhum prontuário encontrado.")

**Integridade**
- [ ] `go test ./...` passa
- [ ] Rota `/patients/{id}/history` continua funcionando (não alterar)

---

### 🚫 NÃO faça

- Não criar banco de dados novo ou nova tabela — tudo vem de `patientService` + `timelineService`
- Não usar `html/template` — apenas `.templ`
- Não passar `patient.Patient` diretamente ao template — usar ViewModel
- Não remover ou renomear o nav item "Pacientes" — apenas adicionar "Prontuários"
- Não implementar busca full-text agora — o campo de busca pode deixar `hx-get` mas a funcionalidade de filtrar por texto é fora de escopo desta task
- Não implementar "Novo prontuário" agora — o botão existe visualmente mas não precisa fazer nada (`href="#"`)

---

### 📎 Padrão de referência

- Handler: siga `internal/web/handlers/timeline_handler.go` (estrutura, extração de tenant DB)
- Template/CSS: siga `web/components/agenda/agenda_layout.templ` + bloco `AGENDA SÁBIO` em `style.css`
- Design completo: `design_handoff_arandu_redesign/page_notes.jsx`
