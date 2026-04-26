# ARANDU — Prompt de Implementação
## Redesign da Tela de Histórico do Paciente (`/patients/{id}/history`)

---

### 🎯 Objetivo

Redesenhar completamente a tela de histórico do paciente para ficar idêntica ao design
em `design_handoff_arandu_redesign/page_patient.jsx`. A tela deve ter:
- **Hero editorial** com avatar, nome em serif, dados clínicos e estatísticas
- **Grid 2 colunas**: coluna principal com triagem em blockquote + timeline clínica;
  coluna lateral com ações rápidas + observações recentes
- **Timeline** com layout em grade (data | dot | conteúdo | chevron), sem as
  border-left-4 cards atuais

---

### 🏗️ Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX · SQLite (database-per-tenant)

**Estrutura de pastas**:
```
internal/domain/timeline/         ← domínio (não alterar)
internal/web/handlers/timeline_handler.go  ← handler a enriquecer
web/components/patient/           ← templates a redesenhar
web/static/css/style.css          ← adicionar classes CSS novas (sabio-pt-*)
```

**Multi-tenancy**: handler já extrai DB via `tenant.TenantDB(r.Context())`. Não alterar
essa mecânica. O `patientService.GetPatientByID(ctx, patientID)` já existe e retorna
o domínio `*patient.Patient`.

**HTMX**: o handler já verifica `HX-Request`. Manter esse padrão. O swap de filtros
deve continuar usando `hx-target="#timeline-content"` + `hx-swap="innerHTML"`.

---

### 📦 Domínio / Dados disponíveis

**PatientService** (já injetado no `TimelineHandler` via `patientService PatientServiceInterface`):
```go
GetPatientByID(ctx, id) (*patient.Patient, error)
```

**Patient campos relevantes**:
```go
type Patient struct {
    ID        string
    Name      string
    BirthDate time.Time   // para calcular idade
    CreatedAt time.Time   // para calcular "em terapia desde"
    Notes     string      // conteúdo de triagem (usar como blockquote)
}
```

**TimelineEvent (domínio)**:
```go
type TimelineEvent struct {
    ID       string
    Type     EventType   // "session" | "observation" | "intervention"
    Date     time.Time
    Content  string
    Metadata map[string]string  // "session_id", "title"
    CreatedAt time.Time
}
```

---

### 📐 ViewModels — criar em `web/components/patient/types.go`

Adicionar (não remover structs existentes):

```go
// PatientHistoryViewModel — ViewModel completo da página de histórico
type PatientHistoryViewModel struct {
    // Hero
    PatientID       string
    PatientName     string
    Initials        string   // primeiras letras de nome + sobrenome
    AgeStr          string   // "34 anos"
    Since           string   // "Desde jan/2026"

    // Stats (hero direita)
    SessionCount    int
    TherapyDuration string   // "3 meses" | "1 ano"
    Frequency       string   // "Semanal" (fixo por ora)

    // Triagem
    TriageContent   string   // patient.Notes (pode estar vazio)
    TriageDate      string   // patient.CreatedAt formatada "02/01/2006"

    // Timeline
    Events          []PatientTimelineEvent
    CurrentFilter   string  // "all" | "session" | "observation" | "intervention"
    PatientIDForURL string  // usado em URLs HTMX

    // Sidebar: observações recentes (últimas 3 do tipo observation)
    RecentObservations []PatientRecentObs
}

type PatientTimelineEvent struct {
    ID        string
    Kind      string  // "Sessão" | "Observação" | "Intervenção"
    KindTone  string  // "accent" (sessão) | "neutral" (outros)
    DateStr   string  // "12 abr · 2026" — monospace
    Title     string  // Metadata["title"] ou fallback por tipo
    Summary   string  // Content truncado em 120 chars
    IsSession bool
    Href      string  // "/session/{session_id}" ou ""
    DotAccent bool   // true = dot preenchido (sessão), false = neutro
}

type PatientRecentObs struct {
    Tag     string  // tipo abreviado: "obs"
    DateStr string  // "12 abr"
    Text    string  // Content truncado em 160 chars
}
```

---

### 🛣️ Handler — modificar `internal/web/handlers/timeline_handler.go`

**Função `ShowPatientHistory`** — enriquecer para montar `PatientHistoryViewModel`.

Lógica a adicionar:
```go
// 1. Buscar paciente
patient, err := h.patientService.GetPatientByID(ctx, patientID)
// se erro: http.Error 500 ou 404

// 2. Calcular campos derivados
initials := buildInitials(patient.Name)   // "AB" de "André Barbosa"
age := time.Now().Year() - patient.BirthDate.Year()
ageStr := fmt.Sprintf("%d anos", age)
since := "Desde " + strings.ToLower(patient.CreatedAt.Format("Jan/2006"))

// 3. Contar sessões na timeline
sessionCount := 0
for _, e := range events { if e.Type == "session" { sessionCount++ } }

// 4. Duração em terapia
months := int(time.Since(patient.CreatedAt).Hours() / 730)
therapyDuration := formatDuration(months)  // "3 meses", "1 ano e 2 meses"

// 5. Mapear events para PatientTimelineEvent
// 6. Pegar últimas 3 observations para recentObservations
// 7. Renderizar com PatientHistoryPage(vm)
```

Funções auxiliares (criar no mesmo arquivo):
```go
func buildInitials(name string) string  // pega até 2 palavras, primeira letra de cada
func formatTherapyDuration(months int) string  // "3 meses" | "1 ano" | "1 ano e 2 meses"
func truncate(s string, n int) string   // trunca em n chars, adiciona "…"
func buildTimelineTitle(e *timeline.TimelineEvent) string
// → Metadata["title"] se existir
// → "Sessão clínica" se Type==session
// → "Observação" se Type==observation  
// → "Intervenção" se Type==intervention
```

---

### 🎨 Template — redesenhar `web/components/patient/timeline.templ`

**Manter**: `TimelineContent(data)` e `FiltersAndContent(data)` para o HTMX load-more
e filtros (o id `#timeline-content` precisa continuar existindo).

**Criar nova função principal**:
```templ
templ PatientHistoryPage(vm PatientHistoryViewModel)
```
que usa `layout.Shell(...)` e renderiza `PatientHistoryContent(vm)`.

```templ
templ PatientHistoryContent(vm PatientHistoryViewModel)
```
Esta é a que o handler renderiza em HX-Request.

#### Estrutura HTML do template:

```html
<!-- Hero editorial -->
<div class="sabio-pt-hero">
  <!-- Avatar 72px com iniciais -->
  <div class="sabio-pt-avatar">{ vm.Initials }</div>
  
  <!-- Info central -->
  <div class="sabio-pt-hero-info">
    <div class="sabio-pt-hero-eyebrow">Paciente #{ vm.PatientID }</div>
    <h1 class="sabio-pt-hero-name">{ vm.PatientName }</h1>
    <div class="sabio-pt-hero-meta">{ vm.AgeStr } · { vm.Since }</div>
  </div>
  
  <!-- Stats -->
  <div class="sabio-pt-stats">
    <div class="sabio-pt-stat">
      <div class="sabio-pt-stat-value">{ vm.SessionCount }</div>
      <div class="sabio-pt-stat-label">Sessões</div>
    </div>
    <div class="sabio-pt-stat-divider"></div>
    <div class="sabio-pt-stat">
      <div class="sabio-pt-stat-value">{ vm.TherapyDuration }</div>
      <div class="sabio-pt-stat-label">Em terapia</div>
    </div>
    <div class="sabio-pt-stat-divider"></div>
    <div class="sabio-pt-stat">
      <div class="sabio-pt-stat-value">{ vm.Frequency }</div>
      <div class="sabio-pt-stat-label">Frequência</div>
    </div>
  </div>
</div>

<!-- Grid 2 colunas -->
<div class="sabio-pt-grid">
  <!-- Coluna principal (2fr) -->
  <div class="sabio-pt-main">
    <!-- Triagem (só renderiza se TriageContent não for vazio) -->
    if vm.TriageContent != "" {
      <div class="sabio-pt-card">
        <div class="sabio-pt-card-header">
          <div class="sabio-pt-eyebrow">Notas de triagem</div>
          <h3 class="sabio-pt-card-title">Primeira escuta</h3>
          <div class="sabio-pt-card-subtitle">Registrada em { vm.TriageDate }</div>
        </div>
        <div class="sabio-pt-card-body">
          <blockquote class="sabio-pt-triage">{ vm.TriageContent }</blockquote>
        </div>
      </div>
    }

    <!-- Timeline card -->
    <div class="sabio-pt-card">
      <div class="sabio-pt-card-header">
        <div>
          <div class="sabio-pt-eyebrow">Percurso</div>
          <h3 class="sabio-pt-card-title">Linha do tempo clínica</h3>
        </div>
        <!-- Filter buttons -->
        <div class="sabio-pt-filters">
          <button class={ filterBtnClass("all", vm.CurrentFilter) }
            hx-get={ "/patients/" + vm.PatientIDForURL + "/history?filter=all" }
            hx-target="#timeline-content" hx-swap="innerHTML">Tudo</button>
          <button class={ filterBtnClass("session", vm.CurrentFilter) }
            hx-get={ "/patients/" + vm.PatientIDForURL + "/history?filter=session" }
            hx-target="#timeline-content" hx-swap="innerHTML">Sessões</button>
          <button class={ filterBtnClass("observation", vm.CurrentFilter) }
            hx-get={ "/patients/" + vm.PatientIDForURL + "/history?filter=observation" }
            hx-target="#timeline-content" hx-swap="innerHTML">Notas</button>
        </div>
      </div>
      <div id="timeline-content" class="sabio-pt-card-body sabio-pt-timeline">
        for i, e := range vm.Events {
          @PatientTimelineEventRow(e, i == len(vm.Events)-1)
        }
      </div>
    </div>
  </div>

  <!-- Coluna lateral (1fr) -->
  <div class="sabio-pt-sidebar">
    <!-- Ações rápidas -->
    <div class="sabio-pt-card">
      <div class="sabio-pt-card-header">
        <div class="sabio-pt-eyebrow">Atalhos</div>
        <h3 class="sabio-pt-card-title">Ações</h3>
      </div>
      <div class="sabio-pt-card-body sabio-pt-actions">
        <a href={ templ.URL("/patient/" + vm.PatientIDForURL + "/sessions/new") }
           class="sabio-pt-action sabio-pt-action-primary">
          + Nova sessão
        </a>
        <a href="#" class="sabio-pt-action">Anamnese</a>
        <a href="#" class="sabio-pt-action">Plano terapêutico</a>
      </div>
    </div>

    <!-- Observações recentes -->
    if len(vm.RecentObservations) > 0 {
      <div class="sabio-pt-card">
        <div class="sabio-pt-card-header">
          <div class="sabio-pt-eyebrow">Percepções</div>
          <h3 class="sabio-pt-card-title">Observações recentes</h3>
        </div>
        <div class="sabio-pt-card-body sabio-pt-obs-list">
          for _, o := range vm.RecentObservations {
            <div class="sabio-pt-obs-card">
              <div class="sabio-pt-obs-meta">
                <span class="sabio-pt-obs-tag">#{ o.Tag }</span>
                <span class="sabio-pt-obs-date">{ o.DateStr }</span>
              </div>
              <p class="sabio-pt-obs-text">{ o.Text }</p>
            </div>
          }
        </div>
      </div>
    }
  </div>
</div>
```

**TimelineEventRow** (inline no template):
```html
<!-- grid: 96px data | 28px dot | 1fr conteúdo | auto chevron -->
<div class={ "sabio-pt-event", timelineEventBorder(last) }
     if e.IsSession { role="link" hx-get={e.Href} ... } >
  <span class="sabio-pt-event-date">{ e.DateStr }</span>
  <div class={ "sabio-pt-event-dot", timelineDotClass(e.DotAccent) }></div>
  <div class="sabio-pt-event-body">
    <div class="sabio-pt-event-header">
      <h4 class="sabio-pt-event-title">{ e.Title }</h4>
      <span class={ "sabio-pt-event-pill", timelinePillClass(e.KindTone) }>{ e.Kind }</span>
    </div>
    <p class="sabio-pt-event-summary">{ e.Summary }</p>
  </div>
  if e.IsSession {
    <span class="sabio-pt-event-chevron">›</span>
  }
</div>
```

---

### 🎨 CSS — adicionar em `web/static/css/style.css`

Adicionar ao final do arquivo (após a seção `AGENDA SÁBIO`), bloco `PATIENT HISTORY SÁBIO`:

```css
/* ============================================
   PATIENT HISTORY SÁBIO
   ============================================ */

/* Hero */
.sabio-pt-hero {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 28px;
  align-items: center;
  padding: 8px 0 24px;
  border-bottom: 1px solid var(--line);
  margin-bottom: 20px;
}

.sabio-pt-avatar {
  width: 72px; height: 72px; border-radius: 50%;
  background: linear-gradient(135deg, color-mix(in oklab, var(--accent) 80%, var(--paper)), var(--accent));
  color: var(--paper);
  display: flex; align-items: center; justify-content: center;
  font-family: var(--font-serif);
  font-size: 26px; font-weight: 500;
  flex-shrink: 0;
}

.sabio-pt-hero-eyebrow {
  font-size: 11px; letter-spacing: 1.4px; text-transform: uppercase;
  color: var(--ink-3); font-weight: 500; margin-bottom: 8px;
}

.sabio-pt-hero-name {
  margin: 0; font-family: var(--font-serif);
  font-size: 40px; font-weight: 400; letter-spacing: -0.8px; line-height: 1;
  color: var(--ink);
}

.sabio-pt-hero-meta {
  margin-top: 10px; font-size: 13px; color: var(--ink-3);
}

.sabio-pt-stats {
  display: flex; align-items: center; gap: 24px;
}

.sabio-pt-stat { text-align: right; }

.sabio-pt-stat-value {
  font-family: var(--font-serif);
  font-size: 26px; font-weight: 500; letter-spacing: -0.5px; line-height: 1;
  color: var(--ink);
}

.sabio-pt-stat-label {
  font-size: 10.5px; letter-spacing: 1.4px; text-transform: uppercase;
  color: var(--ink-3); margin-top: 4px; font-weight: 500;
}

.sabio-pt-stat-divider {
  width: 1px; height: 40px; background: var(--line); flex-shrink: 0;
}

/* Grid principal */
.sabio-pt-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.sabio-pt-main { display: flex; flex-direction: column; gap: 20px; }
.sabio-pt-sidebar { display: flex; flex-direction: column; gap: 20px; }

/* Cards */
.sabio-pt-card {
  background: var(--paper-2);
  border: 1px solid var(--line);
  border-radius: var(--radius, 12px);
  overflow: hidden;
}

.sabio-pt-card-header {
  padding: 16px 20px 12px;
  border-bottom: 1px solid var(--line);
  display: flex; align-items: flex-start;
  justify-content: space-between; gap: 16px;
}

.sabio-pt-eyebrow {
  font-size: 10.5px; letter-spacing: 1.4px; text-transform: uppercase;
  color: var(--ink-3); font-weight: 500; margin-bottom: 4px;
}

.sabio-pt-card-title {
  margin: 0; font-family: var(--font-serif);
  font-size: 19px; font-weight: 500; color: var(--ink);
  letter-spacing: -0.2px; line-height: 1.25;
}

.sabio-pt-card-subtitle {
  margin-top: 4px; font-size: 13px; color: var(--ink-3);
}

.sabio-pt-card-body { padding: 18px 20px; }

/* Triagem blockquote */
.sabio-pt-triage {
  margin: 0; padding: 4px 0 4px 20px;
  border-left: 2px solid var(--accent);
  font-family: var(--font-serif);
  font-size: 18px; line-height: 1.55;
  color: var(--ink-2); font-weight: 400;
  letter-spacing: -0.1px;
}

/* Filtros */
.sabio-pt-filters { display: flex; gap: 6px; flex-shrink: 0; }

.sabio-pt-filter-btn {
  padding: 5px 11px; border-radius: 20px;
  background: transparent;
  border: 1px solid var(--line);
  color: var(--ink-3);
  font-size: 11.5px; font-weight: 500; cursor: pointer;
  transition: background .15s, color .15s;
}

.sabio-pt-filter-btn-active {
  padding: 5px 11px; border-radius: 20px;
  background: var(--ink);
  border: 1px solid var(--ink);
  color: var(--paper);
  font-size: 11.5px; font-weight: 500; cursor: pointer;
}

/* Timeline */
.sabio-pt-timeline { padding: 0 20px; }

.sabio-pt-event {
  display: grid;
  grid-template-columns: 96px 28px 1fr auto;
  gap: 16px; align-items: flex-start;
  padding: 18px 0;
  border-bottom: 1px dashed var(--line);
  cursor: default; width: 100%; text-align: left;
}

.sabio-pt-event-last { border-bottom: none; }

.sabio-pt-event-clickable { cursor: pointer; }
.sabio-pt-event-clickable:hover .sabio-pt-event-title { color: var(--accent); }

.sabio-pt-event-date {
  font-family: var(--font-mono);
  font-size: 12px; color: var(--ink-3); padding-top: 2px;
}

.sabio-pt-event-dot {
  width: 28px; height: 28px; border-radius: 50%;
  background: var(--paper-2);
  border: 1px solid var(--line);
  flex-shrink: 0;
}

.sabio-pt-event-dot-accent {
  width: 28px; height: 28px; border-radius: 50%;
  background: var(--accent);
  border: 1px solid var(--accent-deep);
  flex-shrink: 0;
}

.sabio-pt-event-body {}

.sabio-pt-event-header {
  display: flex; align-items: center; gap: 8px; margin-bottom: 4px;
}

.sabio-pt-event-title {
  margin: 0; font-family: var(--font-serif);
  font-size: 16px; font-weight: 500; color: var(--ink);
  transition: color .15s;
}

.sabio-pt-event-pill {
  font-size: 10.5px; font-weight: 500; letter-spacing: .2px;
  text-transform: uppercase; border-radius: 999px; padding: 2px 8px;
  line-height: 1;
}

.sabio-pt-pill-accent {
  background: color-mix(in oklab, var(--accent) 14%, transparent);
  color: var(--accent-deep);
  border: 1px solid color-mix(in oklab, var(--accent) 30%, transparent);
}

.sabio-pt-pill-neutral {
  background: color-mix(in oklab, var(--ink) 6%, transparent);
  color: var(--ink-2);
  border: 1px solid var(--line);
}

.sabio-pt-event-summary {
  margin: 0; font-size: 13px; color: var(--ink-3); line-height: 1.5;
}

.sabio-pt-event-chevron {
  font-size: 18px; color: var(--ink-4); padding-top: 4px;
}

/* Ações rápidas */
.sabio-pt-actions { display: flex; flex-direction: column; gap: 8px; }

.sabio-pt-action {
  display: flex; align-items: center;
  padding: 11px 14px; border-radius: 10px;
  background: var(--paper);
  border: 1px solid var(--line);
  color: var(--ink-2);
  font-size: 13px; font-weight: 500;
  text-decoration: none;
  transition: background .15s;
}

.sabio-pt-action:hover { background: var(--paper-2); }

.sabio-pt-action-primary {
  background: var(--ink); border-color: var(--ink); color: var(--paper);
}
.sabio-pt-action-primary:hover { background: var(--ink-2); }

/* Observações recentes */
.sabio-pt-obs-list { display: flex; flex-direction: column; gap: 14px; }

.sabio-pt-obs-card {
  padding: 12px 14px;
  background: var(--paper);
  border: 1px solid var(--line);
  border-radius: 10px;
}

.sabio-pt-obs-meta {
  display: flex; align-items: center; gap: 8px; margin-bottom: 8px;
}

.sabio-pt-obs-tag {
  font-size: 10.5px; font-weight: 500; letter-spacing: .2px;
  text-transform: uppercase; border-radius: 999px; padding: 2px 8px;
  background: color-mix(in oklab, var(--accent) 14%, transparent);
  color: var(--accent-deep);
  border: 1px solid color-mix(in oklab, var(--accent) 30%, transparent);
}

.sabio-pt-obs-date {
  font-family: var(--font-mono); font-size: 10.5px;
  color: var(--ink-3); margin-left: auto;
}

.sabio-pt-obs-text {
  margin: 0; font-family: var(--font-serif);
  font-size: 13.5px; line-height: 1.55; color: var(--ink-2);
}
```

---

### 📁 Arquivos a criar/modificar

**Modificar:**
- `web/components/patient/types.go` — adicionar `PatientHistoryViewModel`, `PatientTimelineEvent`, `PatientRecentObs`
- `web/components/patient/timeline.templ` — redesenhar layout inteiro; manter `TimelineContent` e `FiltersAndContent` para HTMX; renomear função principal para `PatientHistoryPage` / `PatientHistoryContent`
- `internal/web/handlers/timeline_handler.go` — enriquecer `ShowPatientHistory` para montar o novo ViewModel
- `web/static/css/style.css` — adicionar bloco CSS `PATIENT HISTORY SÁBIO` ao final

**Não criar novos arquivos** — tudo vai nos arquivos existentes acima.

---

### 🔒 Privacidade

- `Patient.Notes` (triagem) e `TimelineEvent.Content` — **Tier 2**: nunca logar o conteúdo,
  nunca incluir em payload de IA. Apenas renderizar na tela.

---

### 🔧 Funções auxiliares Go a adicionar no handler

```go
func buildInitials(name string) string {
    parts := strings.Fields(name)
    if len(parts) == 0 { return "?" }
    if len(parts) == 1 { return strings.ToUpper(string([]rune(parts[0])[:1])) }
    return strings.ToUpper(string([]rune(parts[0])[:1]) + string([]rune(parts[len(parts)-1])[:1]))
}

func formatTherapyDuration(createdAt time.Time) string {
    months := int(time.Since(createdAt).Hours() / 730)
    if months < 1 { return "< 1 mês" }
    years := months / 12
    rem := months % 12
    if years == 0 { return fmt.Sprintf("%d meses", months) }
    if rem == 0 { return fmt.Sprintf("%d ano", years) }
    return fmt.Sprintf("%d ano e %d meses", years, rem)
}

func truncateStr(s string, n int) string {
    runes := []rune(s)
    if len(runes) <= n { return s }
    return string(runes[:n]) + "…"
}

func buildTimelineTitle(e *timeline.TimelineEvent) string {
    if t, ok := e.Metadata["title"]; ok && t != "" { return t }
    switch e.Type {
    case timeline.EventTypeSession:      return "Sessão clínica"
    case timeline.EventTypeObservation:  return "Observação"
    case timeline.EventTypeIntervention: return "Intervenção"
    default: return "Evento"
    }
}
```

---

### 🔧 Funções auxiliares Templ (no `timeline.templ`)

```go
func filterBtnClass(filter, current string) string {
    if filter == current { return "sabio-pt-filter-btn-active" }
    return "sabio-pt-filter-btn"
}

func timelineEventBorder(last bool) string {
    if last { return "sabio-pt-event sabio-pt-event-last" }
    return "sabio-pt-event"
}

func timelineDotClass(accent bool) string {
    if accent { return "sabio-pt-event-dot-accent" }
    return "sabio-pt-event-dot"
}

func timelinePillClass(tone string) string {
    if tone == "accent" { return "sabio-pt-event-pill sabio-pt-pill-accent" }
    return "sabio-pt-event-pill sabio-pt-pill-neutral"
}
```

---

### ✅ Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Comportamento**
- [ ] GET `/patients/p0014/history` renderiza a nova tela com hero editorial
- [ ] Hero mostra avatar com iniciais, nome em serif grande, stats (sessões, duração, frequência)
- [ ] Card de triagem aparece se `patient.Notes` não estiver vazio
- [ ] Timeline lista eventos com grid data | dot | título+pill | chevron (sessão)
- [ ] Filtros "Tudo / Sessões / Notas" alternam via HTMX sem recarregar a página
- [ ] Coluna lateral mostra "Nova sessão" como botão primário escuro + botões secundários
- [ ] Coluna lateral mostra observações recentes com pill, data mono e texto serif (se houver)
- [ ] Clicar em evento de sessão navega para a sessão (hx-get ou href)
- [ ] Página funciona com scroll (conteúdo longo não quebra layout)

**Integridade**
- [ ] `go test ./...` passa (não quebrar testes existentes)
- [ ] Rota `/patients/{id}/history` continua respondendo (não alterar `main.go`)

---

### 🚫 NÃO faça

- Não remover as funções `TimelineContent` e `FiltersAndContent` do templ — o HTMX de filtros depende delas
- Não usar `html/template` — apenas `.templ`
- Não alterar rotas em `main.go`
- Não remover structs existentes em `types.go` — apenas adicionar as novas
- Não usar `fmt.Sprintf` para construir URLs em templ — use `templ.URL()` ou concatenação com `templ.URL()`
- Não adicionar CSS inline quando já existe classe CSS equivalente no bloco acima
- Não alterar o mecanismo de multi-tenancy no handler

---

### 📎 Padrão de referência

- Handler: siga o padrão de `internal/web/handlers/agenda_handler.go` para estrutura
  e montagem de ViewModel
- Template/CSS: siga o padrão de `web/components/agenda/agenda_layout.templ`
  e o bloco CSS `AGENDA SÁBIO` em `web/static/css/style.css` como referência de estilo
- Design completo: `design_handoff_arandu_redesign/page_patient.jsx`
