# ARANDU — Prompt de Implementação
## Conteúdo das abas Observações e Intervenções em `/notes`

---

### 🎯 Objetivo

As abas "Observações" e "Intervenções" do painel de prontuário em `/notes` mostram o
count correto (ex: "Observações 20p") mas renderizam "Sem registros nesta seção."
ao serem clicadas. Fazer com que exibam a lista real dos eventos.

---

### 🐛 Causa raiz

1. `NoteRecordDetail` em `web/components/notes/types.go` tem `ObservationCount int` e
   `InterventionCount int`, mas **não tem lista** com os itens.
2. `buildNoteRecordDetail` em `internal/web/handlers/notes_handler.go` preenche os
   contadores mas **não popula nenhuma lista**.
3. O template `notes_page.templ` no bloco `else` cai sempre em `"Sem registros"`.

---

### 📐 Mudança 1 — `web/components/notes/types.go`

Adicionar campos de lista ao `NoteRecordDetail`:

```go
type NoteRecordDetail struct {
    // ... campos existentes (não remover) ...

    // Listas para as abas (adicionar estes):
    ObservationItems  []NoteEventItem
    InterventionItems []NoteEventItem
    AnamneseText      string  // patient.Notes (para aba anamnese)
}

// NoteEventItem — item de uma observação ou intervenção
type NoteEventItem struct {
    DateStr string  // "12 abr · 2026"
    Content string  // truncado em 300 chars
}
```

---

### 📐 Mudança 2 — `internal/web/handlers/notes_handler.go`

Em `buildNoteRecordDetail`, após o bloco que monta sessions/observations/interventions,
adicionar:

```go
// Montar listas para as abas
detail.ObservationItems = make([]notesComponents.NoteEventItem, 0, len(observations))
for _, e := range observations {
    detail.ObservationItems = append(detail.ObservationItems, notesComponents.NoteEventItem{
        DateStr: e.Date.Format("02 Jan · 2006"),
        Content: notesComponents.TruncateStr(e.Content, 300),
    })
}

detail.InterventionItems = make([]notesComponents.NoteEventItem, 0, len(interventions))
for _, e := range interventions {
    detail.InterventionItems = append(detail.InterventionItems, notesComponents.NoteEventItem{
        DateStr: e.Date.Format("02 Jan · 2006"),
        Content: notesComponents.TruncateStr(e.Content, 300),
    })
}

// Anamnese vem de patient.Notes (já buscamos p no início da função)
detail.AnamneseText = p.Notes
```

---

### 📐 Mudança 3 — `web/components/notes/notes_page.templ`

Substituir o bloco de renderização do conteúdo do painel (linhas com `ActiveTab`):

```templ
if detail.ActiveTab == "evolucao" && detail.LastEntryTitle != "" {
    <!-- conteúdo existente da aba evolução — não alterar -->
} else if detail.ActiveTab == "evolucao" {
    <p class="sabio-nl-empty">Nenhuma sessão registrada.</p>
} else if detail.ActiveTab == "anamnese" {
    if detail.AnamneseText != "" {
        <div class="sabio-nl-doc">
            <p class="sabio-nl-doc-body">{ detail.AnamneseText }</p>
        </div>
    } else {
        <p class="sabio-nl-empty">Anamnese não registrada.</p>
    }
} else if detail.ActiveTab == "observacoes" {
    if len(detail.ObservationItems) > 0 {
        <div class="sabio-nl-event-list">
            for _, item := range detail.ObservationItems {
                <div class="sabio-nl-event-item">
                    <span class="sabio-nl-event-date">{ item.DateStr }</span>
                    <p class="sabio-nl-event-content">{ item.Content }</p>
                </div>
            }
        </div>
    } else {
        <p class="sabio-nl-empty">Nenhuma observação registrada.</p>
    }
} else if detail.ActiveTab == "intervencoes" {
    if len(detail.InterventionItems) > 0 {
        <div class="sabio-nl-event-list">
            for _, item := range detail.InterventionItems {
                <div class="sabio-nl-event-item">
                    <span class="sabio-nl-event-date">{ item.DateStr }</span>
                    <p class="sabio-nl-event-content">{ item.Content }</p>
                </div>
            }
        </div>
    } else {
        <p class="sabio-nl-empty">Nenhuma intervenção registrada.</p>
    }
} else {
    <p class="sabio-nl-empty">Sem registros.</p>
}
```

---

### 🎨 CSS — adicionar ao bloco `NOTES LIBRARY SÁBIO` em `web/static/css/style.css`

```css
.sabio-nl-event-list {
  display: flex; flex-direction: column; gap: 16px;
}

.sabio-nl-event-item {
  padding: 14px 0;
  border-bottom: 1px dashed var(--line);
}

.sabio-nl-event-item:last-child { border-bottom: none; }

.sabio-nl-event-date {
  font-family: var(--font-mono);
  font-size: 11px; color: var(--ink-4);
  display: block; margin-bottom: 6px;
}

.sabio-nl-event-content {
  margin: 0; font-family: var(--font-serif);
  font-size: 15px; line-height: 1.6; color: var(--ink-2);
}
```

---

### ✅ Critérios de aceite

- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros
- [ ] Clicar na aba "Observações" de um paciente com observações → lista os itens
- [ ] Clicar na aba "Intervenções" de um paciente com intervenções → lista os itens
- [ ] Clicar na aba "Anamnese" → mostra `patient.Notes` ou mensagem de vazio
- [ ] Abas sem dados mostram mensagem "Nenhuma X registrada." (não "Sem registros nesta seção.")
- [ ] A aba "Evolução" continua funcionando como antes

---

### 🚫 NÃO faça

- Não alterar a estrutura do hero, chips ou lista esquerda
- Não remover campos existentes de `NoteRecordDetail`
- Não criar novas rotas — o `Detail` handler existente já recebe `?tab=` e chama `buildNoteRecordDetail`
- Não paginar os itens — exibir todos (a busca já limita a 20 eventos no `GetPatientTimeline`)

---

### 📎 Arquivos a modificar

1. `web/components/notes/types.go` — adicionar `ObservationItems`, `InterventionItems`, `AnamneseText`, `NoteEventItem`
2. `internal/web/handlers/notes_handler.go` — popular as listas em `buildNoteRecordDetail`
3. `web/components/notes/notes_page.templ` — renderizar conteúdo real nas abas
4. `web/static/css/style.css` — adicionar `.sabio-nl-event-list`, `.sabio-nl-event-item`, `.sabio-nl-event-date`, `.sabio-nl-event-content`
