# Task: Agenda — corrigir bugs de resposta HTMX e handlers de sub-rota
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

A agenda clínica está quase completa. O handler, service, repository, domínio, migration e as 3 views (dia/semana/mês) já foram implementados e o app **compila e sobe sem erros**.

Após inspeção completa do código, identificamos três bugs que impedem o uso funcional da agenda:

### Bug 1 — Cancel e Complete retornam 200 vazio para HTMX (crítico)

Em `appointment_detail.templ`, os botões "Cancelar" e "Concluir" usam:

```html
hx-post="..."
hx-target="#agenda-content"
hx-swap="outerHTML"
onclick="document.getElementById('modal-container').innerHTML=''"
```

O `onclick` fecha o modal antes da requisição HTMX disparar. O handler `Cancel` e `Complete` em `internal/web/handlers/agenda_handler.go` terminam com:

```go
w.WriteHeader(http.StatusOK)
```

Sem corpo HTML. O HTMX recebe 200 com body vazio e substitui `#agenda-content` por nada — o calendário desaparece da tela.

**Referência do padrão correto:** handler `Create` (mesma linha de código) usa:
```go
if r.Header.Get("HX-Request") == "true" {
    w.Header().Set("HX-Redirect", redirectURL)
    w.WriteHeader(http.StatusOK)
    return
}
http.Redirect(w, r, redirectURL, http.StatusSeeOther)
```

### Bug 2 — Sub-handlers `/agenda/month` e `/agenda/week` não populam appointments

Os handlers dedicados `MonthView` e `WeekView` (rotas `/agenda/month` e `/agenda/week`) criam `DayViewModel` sem os appointments. Exemplo em `MonthView` (linha ~412):

```go
days = append(days, agendaComponents.DayViewModel{
    Date:      day.Date,
    DayName:   day.Date.Format("Mon"),
    DayNumber: day.Date.Format("2"),
    IsToday:   isToday(day.Date),
    // appointments ausentes — MonthView.Appointments não é indexado aqui
})
```

Já o handler principal `View` (rota `/agenda?view=mes`) usa `viewModelForMonth` que corretamente indexa `monthView.Appointments` por data usando um map e popula `day.Appointments`.

O mesmo problema existe em `WeekView`: o `viewModel` não inclui `PrevDate` e `NextDate`, quebrando os botões de navegação.

### Bug 3 — Sem feedback de conflito no formulário de nova marcação

Quando o usuário tenta criar uma marcação em horário ocupado, o handler `Create` retorna HTTP 409. Mas o formulário em `new_appointment_form.templ` não tem `hx-on::htmx:response-error` nem nenhuma área de erro — o usuário não vê feedback algum.

---

## O que implementar

### Fix 1 — Handlers Cancel e Complete (arquivo: `internal/web/handlers/agenda_handler.go`)

Nos handlers `Cancel` e `Complete`, substituir a resposta final pelo padrão HTMX do projeto:

```go
// Em Cancel, após CancelAppointment bem-sucedido:
redirectURL := "/agenda"
if r.Header.Get("HX-Request") == "true" {
    w.Header().Set("HX-Redirect", redirectURL)
    w.WriteHeader(http.StatusOK)
    return
}
http.Redirect(w, r, redirectURL, http.StatusSeeOther)

// Em Complete, idem (redirecionar para o dia do appointment se possível, ou /agenda)
```

Para `Complete`, o redirect pode ser para `/agenda?view=dia&date={date-do-appointment}` se o appointment for encontrado antes de marcar como completo (já é carregado na linha ~793).

### Fix 2 — MonthView handler (`internal/web/handlers/agenda_handler.go`)

Substituir a lógica do handler `MonthView` para usar o método `viewModelForMonth` já existente, exatamente como `View` faz:

```go
func (h *AgendaHandler) MonthView(w http.ResponseWriter, r *http.Request) {
    // ... parse year/month params (mantém código atual)
    date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

    ctx := r.Context()
    vm, err := h.viewModelForMonth(ctx, date)   // ← usar método existente
    if err != nil {
        http.Error(w, "Failed to load month view", http.StatusInternalServerError)
        return
    }
    vm.CurrentView = "mes"

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if r.Header.Get("HX-Request") == "true" {
        agendaComponents.AgendaContent(vm).Render(ctx, w)
        return
    }
    agendaComponents.AgendaPage(vm).Render(ctx, w)
}
```

Da mesma forma, o handler `WeekView` deve usar `viewModelForWeek` e incluir `CurrentView = "semana"`.

### Fix 3 — Feedback de conflito no formulário (arquivo: `web/components/agenda/new_appointment_form.templ`)

Adicionar área de erro e handler HTMX para resposta 409. CSS class via helper em `types.go` (obrigatório — Tailwind v4):

Em `types.go`, adicionar:
```go
func ConflictAlertClasses(visible bool) string {
    if visible {
        return "mt-2 p-3 rounded-lg text-sm bg-red-50 border border-red-200 text-red-700"
    }
    return "hidden"
}
```

No form em `new_appointment_form.templ`, antes das actions:
```html
<!-- Conflict warning (shown via Alpine or HTMX error handler) -->
<div id="conflict-warning" class={ ConflictAlertClasses(false) }>
    <i class="fas fa-exclamation-triangle mr-1.5"></i>
    Horário ocupado. Escolha outro slot disponível.
</div>
```

No elemento `<form>`, adicionar o atributo para capturar erros HTMX:
```html
hx-on::htmx:response-error="
    if(event.detail.xhr.status===409){
        document.getElementById('conflict-warning').className='mt-2 p-3 rounded-lg text-sm bg-red-50 border border-red-200 text-red-700'
    }
"
```

> **Nota Tailwind v4** (skill: arandu-architecture § Tailwind v4): a função `ConflictAlertClasses` deve retornar a string completa — nunca construir classes dinamicamente. Se a visibilidade mudar pelo JS inline, as classes precisam ser strings literais (não concatenação).

---

## Checklist de implementação

- [ ] Fix Cancel: adicionar HX-Redirect pattern (`internal/web/handlers/agenda_handler.go`)
- [ ] Fix Complete: idem, redirect para `/agenda?view=dia&date={date}`
- [ ] Fix MonthView handler: usar `viewModelForMonth` existente
- [ ] Fix WeekView handler: usar `viewModelForWeek` existente, incluir PrevDate/NextDate
- [ ] Fix 3: adicionar `ConflictAlertClasses` em `web/components/agenda/types.go`
- [ ] Fix 3: adicionar área de erro e `hx-on::htmx:response-error` no form
- [ ] Rodar `~/go/bin/templ generate ./web/components/agenda/...` após editar `.templ`
- [ ] `go build ./cmd/arandu/` sem erros
- [ ] Testar manualmente: criar marcação → aparece no calendário ✓
- [ ] Testar manualmente: cancelar marcação → calendário permanece visível ✓
- [ ] Testar manualmente: tentar criar marcação em horário ocupado → exibe aviso ✓

## Arquivos a modificar

- `internal/web/handlers/agenda_handler.go` — Fix 1 (Cancel, Complete) + Fix 2 (MonthView, WeekView)
- `web/components/agenda/new_appointment_form.templ` — Fix 3 (conflito)
- `web/components/agenda/types.go` — Fix 3 (helper CSS)

**Não criar arquivos novos.** Toda a estrutura já existe.

## Skills de referência obrigatória

- **arandu-architecture** — multi-tenancy, padrão HTMX fragmento vs página completa (`r.Header.Get("HX-Request")`), regra Tailwind v4 (nunca concatenar classes dinamicamente, sempre helper function em `types.go`)
- **go-templ-htmx-ux** — padrão HX-Redirect após mutations, `hx-on::htmx:response-error` para feedback de erro, `hx-swap="outerHTML"` em containers nomeados
- **ddd-go** — handler só coordena; não duplicar lógica (usar `viewModelForMonth`/`viewModelForWeek` já existentes em vez de reescrever)
