# Task: CA-09 — Agenda redesign Sábio
Requirement: Redesign visual — handoff `design_handoff_arandu_redesign/page_agenda.jsx`
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Objetivo

Reescrever a tela de agenda (`web/components/agenda/agenda_layout.templ`) para ficar fiel ao
design proposto em `design_handoff_arandu_redesign/page_agenda.jsx`. A referência visual exata
está em `design_handoff_arandu_redesign/Arandu Redesign.html` — abra no browser, clique em
"Agenda" na sidebar.

---

## Contexto do sistema

**Stack**: Go · Templ · HTMX · DaisyUI v5 + Tailwind CSS v4 · Alpine.js
**Paleta Sábio**: `--paper`, `--paper-2`, `--paper-3`, `--ink`, `--ink-2`, `--ink-3`, `--ink-4`,
`--line`, `--line-2`, `--accent`, `--accent-deep`, `--accent-soft`, `--moss`, `--moss-2`, `--clay`,
`--danger` — definidas em `web/static/css/style.css`
**HTMX existente**: `hx-get="/agenda?view=X&date=Y"` → `hx-target="#agenda-content"` →
`hx-swap="outerHTML"` — **preservar** este comportamento
**Modal de marcação**: `hx-get="/agenda/appointments/{id}"` → `#modal-container` — **preservar**

---

## Arquivos a modificar

```
web/components/agenda/agenda_layout.templ  ← reescrever completamente
web/components/agenda/types.go             ← adicionar campos Tone, Metrics
internal/web/handlers/agenda_handler.go   ← popular novos campos do ViewModel
web/static/css/style.css                  ← adicionar .sabio-agenda-* classes
```

---

## Parte 1 — Novos campos no ViewModel

### `web/components/agenda/types.go`

Adicionar ao `AppointmentViewModel`:
```go
Tone string // "accent" | "info" | "warn" | "danger" | "moss" | "ghost" | "neutral"
```

Adicionar ao `AgendaViewModel`:
```go
Metrics []AgendaMetric // para o hero
```

Criar struct:
```go
type AgendaMetric struct {
    Value string
    Label string
}
```

Adicionar função helper de mapeamento de status → tone:
```go
// StatusToTone mapeia o status de agendamento para o tone visual Sábio
func StatusToTone(status string) string {
    switch status {
    case "confirmed", "completed", "scheduled":
        return "accent"   // marrom — retorno
    case "first_session":
        return "info"     // azul — 1ª consulta
    case "pending":
        return "warn"     // âmbar — atenção
    case "no_show", "risk":
        return "danger"   // vermelho — risco
    case "cancelled":
        return "ghost"    // hachurado — cancelado
    default:
        return "neutral"
    }
}
```

### `internal/web/handlers/agenda_handler.go`

Na função `convertAppointmentsForWeek` (ou onde AppointmentViewModel é criado), adicionar:
```go
Tone: agendaComponents.StatusToTone(mapAppointmentStatus(string(a.Status))),
```

Nas funções `viewModelForWeek`, `viewModelForDay`, `viewModelForMonth`, popular `Metrics`:
```go
// Exemplo para semana
Metrics: []agendaComponents.AgendaMetric{
    {Value: fmt.Sprintf("%d", total), Label: "sessões esta semana"},
    {Value: fmt.Sprintf("%d", confirmed), Label: "confirmadas"},
    {Value: fmt.Sprintf("%d", pending), Label: "pendentes"},
},
```
(calcular `confirmed` e `pending` contando os appts por status)

---

## Parte 2 — Reescrever `agenda_layout.templ`

### Estrutura geral do `AgendaContent`

```templ
templ AgendaContent(vm AgendaViewModel) {
    <div id="agenda-content" aria-live="polite">
        <!-- Hero -->
        @AgendaHero(vm)

        <!-- Grid principal: calendário (1fr) + sidebar (320px) -->
        <div style="display:grid;grid-template-columns:1fr 320px;gap:20px;margin-top:22px">
            <!-- Coluna calendário -->
            <div>
                if vm.CurrentView == "dia" {
                    @AgendaDayView(vm)
                } else if vm.CurrentView == "mes" {
                    @AgendaMonthView(vm)
                } else {
                    @AgendaWeekView(vm)
                }
            </div>
            <!-- Sidebar lateral -->
            @AgendaSidebar()
        </div>
    </div>
}
```

---

### `AgendaHero` — hero editorial

```
┌─────────────────────────────────────────┬──────────────────────────────────────┐
│  eyebrow "AGENDA · SEMANA"              │  [Dia][Semana][Mês] ← → Hoje  [+ Novo] │
│  H1 serif 40px: "6 – 12 de ·abril"*    │                                       │
│  N sessões · N confirmadas · N pendentes │                                       │
└─────────────────────────────────────────┴──────────────────────────────────────┘
* "abril" em itálico, cor --accent-deep
```

**CSS a usar:**
- Eyebrow: `font-size:11px; letter-spacing:1.6px; text-transform:uppercase; color:var(--ink-3); font-weight:500`
- H1: `font-family:var(--font-serif); font-size:40px; font-weight:400; letter-spacing:-0.8px; line-height:1; margin:10px 0 0`
- Parte itálica do H1: `<em style="font-style:italic;color:var(--accent-deep)">`
- Metrics: `font-size:13px; color:var(--ink-3)`; valores em `font-weight:500; color:var(--ink-2)`
- Separador borda inferior: `padding-bottom:22px; border-bottom:1px solid var(--line)`

**View toggle pills:**
```html
<div style="padding:3px;background:var(--paper-2);border:1px solid var(--line);border-radius:10px;display:flex">
  <button style="padding:6px 12px;border-radius:7px;font-size:12px;font-weight:500;
                 background:var(--ink);color:var(--paper)">Semana</button>  <!-- ativo -->
  <button style="padding:6px 12px;border-radius:7px;font-size:12px;font-weight:500;
                 background:transparent;color:var(--ink-3)">Dia</button>
```
Com HTMX: cada botão `hx-get="/agenda?view=X&date=Y"` → `hx-target="#agenda-content"` → `hx-swap="outerHTML"`

**Botões de navegação:**
```html
<button style="background:var(--paper-2);border:1px solid var(--line);
               padding:8px;border-radius:10px;color:var(--ink-2)">‹</button>
<button style="...mesmo estilo...;padding:8px 14px;font-size:12.5px">Hoje</button>
<button style="...mesmo...">›</button>
```

**Botão "Novo agendamento":**
```html
<button style="background:var(--ink);color:var(--paper);border:1px solid var(--ink);
               padding:8px 14px;border-radius:10px;font-size:13px;font-weight:500">
    + Novo agendamento
</button>
```

---

### `AgendaWeekView` — vista semanal

**Container:**
```css
background: var(--paper-2);
border: 1px solid var(--line);
border-radius: var(--radius, 14px);
overflow: hidden;
box-shadow: var(--shadow-sm);
```

**Header de dias** (`grid-template-columns: 64px repeat(7, 1fr)`):
- Célula de hora vazia: 64px
- Cada dia: `padding:14px 10px; text-align:center; border-left:1px solid var(--line)`
- DOW: `font-size:10.5px; letter-spacing:1.4px; text-transform:uppercase; color:var(--ink-3); font-weight:500`
- Número: `font-family:var(--font-serif); font-size:22px; font-weight:500; letter-spacing:-0.3px; margin-top:3px`
- Hoje (background): `color-mix(in oklab, var(--accent) 10%, transparent)`
- Hoje (número): `color:var(--accent-deep)`

**Grid de slots** (22 slots × 32px = 704px total):
- Coluna tempo: para cada slot i, `height:32px; border-bottom: i%2==1 ? 1px solid var(--line) : none`
- Hora label (apenas slots pares): `font-family:var(--font-mono); font-size:10.5px; color:var(--ink-4)`, alinhado à direita com `padding-right:10px`
- Coluna de cada dia: `position:relative; border-left:1px solid var(--line)`
- Fundo hoje: `color-mix(in oklab, var(--accent) 4%, transparent)`

**AppointmentCard no WeekView:**
```
position: absolute
left: 4px; right: 4px
top: {event.Start * 32}px
height: {event.Dur * 32 - 2}px
border-left: 3px solid {toneBd}
border-radius: 6px
padding: 5px 8px
font-size: 11px
overflow: hidden
```

**Função `ToneStylesWeek(tone string)` → (bg, bd, fg string):**
```go
func ToneStylesWeek(tone string) (bg, bd, fg string) {
    switch tone {
    case "accent":  return "#EADFCB", "#6B4E3D", "#3E2A1E"
    case "info":    return "#DDE6EE", "#3C5C7A", "#1F3A55"
    case "warn":    return "#F3E3C4", "#B8842A", "#6B4A10"
    case "danger":  return "#EDCFC7", "#A0463A", "#6F241A"
    case "moss":    return "#D9E2D3", "#4A5D4F", "#263326"
    case "ghost":   return "repeating-linear-gradient(135deg,var(--paper-3),var(--paper-3) 4px,var(--paper-2) 4px,var(--paper-2) 8px)", "var(--ink-4)", "var(--ink-3)"
    default:        return "var(--paper)", "var(--ink-4)", "var(--ink-2)"
    }
}
```
Usar `style=fmt.Sprintf("background:%s;border-left:3px solid %s;color:%s;...", bg, bd, fg)` no AppointmentCard.

**Legenda (rodapé):**
```
padding:12px 16px; border-top:1px solid var(--line); display:flex; gap:16px; font-size:11px; color:var(--ink-3)
```
Cada dot: `width:10px; height:10px; border-radius:2px; background:{cor}`

| Rótulo | Cor |
|--------|-----|
| Retorno | `#6B4E3D` |
| 1ª consulta | `#3C5C7A` |
| Atenção | `#B8842A` |
| Risco | `#A0463A` |
| Supervisão | `#4A5D4F` |
| Cancelado | `var(--ink-4)` |

---

### `AgendaDayView` — vista de dia (editorial)

Grid `96px 1fr`, altura mínima 640px.

**Coluna tempo** (96px, `background:var(--paper)`, `border-right:1px solid var(--line)`):
- 11 linhas horárias (08–18), cada com `height:64px; border-bottom:1px dashed var(--line)`
- Hora: `font-family:var(--font-mono); font-size:11px; color:var(--ink-4)`, alinhado à direita

**Coluna de eventos** (`padding:20px 24px; display:flex; flex-direction:column; gap:10px`):
- Cada card: `display:grid; grid-template-columns:auto 1fr auto; gap:16px; align-items:center`
- Container: `padding:14px 18px; background:var(--paper); border:1px solid var(--line); border-left:4px solid {tone.bd}; border-radius:10px`
- Estado "Próxima" (próxima sessão do dia): `border-color:{tone.bd}; box-shadow:0 0 0 3px color-mix(in oklab, {tone.bd} 20%, transparent)`
- Estado concluída: `opacity:0.55`
- Coluna esquerda (hora): `font-family:var(--font-mono); font-size:13px; font-weight:500; color:{tone.bd}`; duração em `10.5px; color:var(--ink-4)`
- Coluna centro: nome `font-family:var(--font-serif); font-size:17px; font-weight:500`; tipo em `12px; color:var(--ink-3)`
- Coluna direita: pill de status

Pills de status (uso de `style=` inline com `color-mix`):
- Concluída: `background:color-mix(in oklab, var(--moss-2) 14%, transparent); color:var(--moss); border:1px solid color-mix(in oklab, var(--moss-2) 30%, transparent)`
- Próxima: `background:color-mix(in oklab, var(--accent) 12%, transparent); color:var(--accent-deep); border:...`
- Agendada: `background:color-mix(in oklab, var(--ink) 6%, transparent); color:var(--ink-2); border:1px solid var(--line)`

---

### `AgendaMonthView` — vista mensal

**Header dias da semana** (`grid-template-columns: repeat(7,1fr)`, `background:var(--paper)`):
- `font-size:10.5px; letter-spacing:1.4px; text-transform:uppercase; color:var(--ink-3); font-weight:500; padding:12px 14px`
- Dias: dom, seg, ter, qua, qui, sex, sáb

**Grid de células** (`grid-template-columns: repeat(7,1fr)`):
- Cada célula: `padding:10px 12px; border-top:1px solid var(--line); border-left:1px solid var(--line); min-height:100px`
- Hoje: `background:color-mix(in oklab, var(--accent) 8%, transparent)`
- Fora do mês: `opacity:0.35`
- Número do dia: `font-family:var(--font-serif); font-size:18px; font-weight:500; letter-spacing:-0.3px`
- Hoje (número): dentro de círculo `background:var(--accent); color:var(--paper); border-radius:50%; width:26px; height:26px; display:inline-flex; align-items:center; justify-content:center`
- Event dots: `width:5px; height:5px; border-radius:50%; background:{tone color}` + nome 10.5px truncado
- "+N mais": `font-size:10px; color:var(--ink-4)`

**Legenda**: igual ao WeekView.

---

### `AgendaSidebar` — lateral 320px

```templ
templ AgendaSidebar() {
    <div style="display:flex;flex-direction:column;gap:18px">
        <!-- Card: Solicitações -->
        <div style="background:var(--paper-2);border:1px solid var(--line);border-radius:var(--radius);box-shadow:var(--shadow-sm);overflow:hidden">
            <div style="padding:16px 20px 12px;border-bottom:1px solid var(--line)">
                <div style="font-size:10.5px;letter-spacing:1.4px;text-transform:uppercase;color:var(--ink-3);font-weight:500;margin-bottom:4px">Solicitações</div>
                <h3 style="font-family:var(--font-serif);font-size:19px;font-weight:500;color:var(--ink);letter-spacing:-0.2px;margin:0">Em aberto</h3>
                <p style="font-size:13px;color:var(--ink-3);margin:4px 0 0">Nenhuma aguardando resposta</p>
            </div>
            <div style="padding:18px 20px">
                <p style="font-family:var(--font-serif);font-style:italic;font-size:13.5px;color:var(--ink-3);text-align:center;padding:20px 0">
                    Nenhuma solicitação no momento.
                </p>
            </div>
        </div>

        <!-- Card: Disponibilidade -->
        <div style="background:var(--paper-2);border:1px solid var(--line);border-radius:var(--radius);box-shadow:var(--shadow-sm);overflow:hidden">
            <div style="padding:16px 20px 12px;border-bottom:1px solid var(--line)">
                <div style="font-size:10.5px;letter-spacing:1.4px;text-transform:uppercase;color:var(--ink-3);font-weight:500;margin-bottom:4px">Disponibilidade</div>
                <h3 style="font-family:var(--font-serif);font-size:19px;font-weight:500;color:var(--ink);letter-spacing:-0.2px;margin:0">Horários livres esta semana</h3>
            </div>
            <div style="padding:18px 20px;display:flex;flex-direction:column;gap:8px">
                <p style="font-family:var(--font-serif);font-style:italic;font-size:13.5px;color:var(--ink-3);text-align:center;padding:8px 0">
                    Configure sua disponibilidade na configuração de agenda.
                </p>
            </div>
        </div>
    </div>
}
```

> Os cards de solicitações e disponibilidade podem ser placeholders estáticos nesta task — dados reais virão em task futura.

---

## Parte 3 — CSS em `web/static/css/style.css`

Adicionar ao final (antes do `/* Cache bust */`):

```css
/* ============================================
   AGENDA SÁBIO
   ============================================ */

/* Container principal */
#agenda-content {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* Semana: grid de slots */
.sabio-week-grid {
  display: grid;
  grid-template-columns: 64px repeat(7, 1fr);
}

.sabio-week-slot {
  height: 32px;
}

.sabio-week-slot--full-hour {
  border-bottom: 1px solid var(--line);
}

/* Mês: célula */
.sabio-month-cell {
  padding: 10px 12px;
  min-height: 100px;
  border-top: 1px solid var(--line);
  border-left: 1px solid var(--line);
}

/* Pill de status inline */
.sabio-status-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  font-size: 10.5px;
  font-weight: 500;
  letter-spacing: 0.2px;
  text-transform: uppercase;
  border-radius: 999px;
  line-height: 1;
  white-space: nowrap;
}
```

---

## Após as edições

```bash
~/go/bin/templ generate ./web/components/...
npm run tailwind:build:v2
go build -o arandu ./cmd/arandu/
# reiniciar servidor
```

---

## Critérios de aceite

**Compilação**
- [ ] `templ generate` sem erros
- [ ] `go build` sem erros

**Visual — comparar com `design_handoff_arandu_redesign/Arandu Redesign.html` → aba Agenda**

- [ ] CA01: Hero com eyebrow "AGENDA · SEMANA", H1 serif 40px, parte itálica em --accent-deep
- [ ] CA02: View toggle pills (Dia/Semana/Mês) com ativo em `--ink` / inativos em `--paper-2`
- [ ] CA03: WeekView — header com DOW + número serif, slots de 32px, eventos com cores por tone
- [ ] CA04: WeekView — hoje com fundo `color-mix(accent 10%)` e número em --accent-deep
- [ ] CA05: WeekView — legenda no rodapé com dots coloridos (Retorno/1ª consulta/Atenção/Risco/Supervisão/Cancelado)
- [ ] CA06: DayView — cards com `border-left:4px solid tone`, nome em serif 17px, hora em mono
- [ ] CA07: MonthView — grid semanal com dots de eventos + número do dia hoje em círculo --accent
- [ ] CA08: Sidebar com dois cards (Solicitações + Disponibilidade) visíveis à direita do calendário
- [ ] CA09: Navegação HTMX funciona — prev/next/hoje atualizam o calendário sem reload
- [ ] CA10: Modal de marcação abre ao clicar em evento (preservado do código anterior)

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa

---

## NÃO faça

- Não remover `id="agenda-content"` nem o `hx-swap="outerHTML"` — HTMX depende disso
- Não remover o hx-get de marcações para `#modal-container`
- Não usar cores fora da paleta Sábio — sem `bg-emerald-100`, `bg-amber-100`, `bg-blue-100`
- Não usar `bg-arandu-primary` em nenhum lugar — usar `var(--accent)` / `var(--accent-deep)`
- Não criar dados reais para solicitações/disponibilidade na sidebar — placeholder é suficiente
- Não alterar o handler de modal (`/agenda/appointments/{id}`)

---

## Referências

- Design visual: `design_handoff_arandu_redesign/Arandu Redesign.html` → sidebar "Agenda"
- Componentes JSX: `design_handoff_arandu_redesign/page_agenda.jsx` (código fonte completo)
- Paleta de tones: ver seção `toneStyles` em `AgendaEvent` e `DayView` no JSX
- Implementação atual: `web/components/agenda/agenda_layout.templ`
- ViewModel atual: `web/components/agenda/types.go`
- Handler atual: `internal/web/handlers/agenda_handler.go`
