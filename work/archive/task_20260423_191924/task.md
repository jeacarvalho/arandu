# Task: CA-04 — Ações de Status no Modal de Agendamento (DaisyUI)
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Objetivo

Adicionar botões de ação ao modal de detalhe do agendamento (`AppointmentDetail`) baseados
no status atual do agendamento. Simultaneamente, migrar o componente de inline styles para
**DaisyUI v4** — esta é a primeira peça do novo design system.

Handlers e service methods já existem. O trabalho é 100% de UI + HTMX wiring.

---

## Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX 2.x · **DaisyUI v4 + Tailwind CSS** · SQLite (database-per-tenant)

**Design system**: DaisyUI v4 com tema Arandu. Consultar skill `daisyui-arandu` antes de escrever qualquer CSS.
Protótipo visual de referência: `design_handoff_arandu_redesign/daisyui_shell_dashboard.html`

**Multi-tenancy**: handler DEVE obter conexão via `tenant.TenantDB(r.Context())`.

**HTMX**: fragmento em requests com `HX-Request: true`, página completa caso contrário.

---

## O que já existe — NÃO reimplementar

```go
// AgendaServiceInterface — já implementado
ConfirmAppointment(ctx, id string) error
MarkNoShow(ctx, id string) error
CompleteAppointment(ctx, id string, sessionID string) error
CancelAppointment(ctx, id string) error

// AgendaHandler — já implementado
func (h *AgendaHandler) Confirm(w, r)   // POST /agenda/appointments/{id}/confirm
func (h *AgendaHandler) NoShow(w, r)    // POST /agenda/appointments/{id}/no-show
func (h *AgendaHandler) Cancel(w, r)    // POST /agenda/appointments/{id}/cancel
func (h *AgendaHandler) Complete(w, r)  // POST /agenda/appointments/{id}/complete
```

Verificar em `cmd/arandu/main.go` se as rotas estão registradas. Adicionar se ausentes.

---

## Domínio — Status e ações

```go
// internal/domain/appointment/appointment.go
type AppointmentStatus string

const (
    StatusScheduled AppointmentStatus = "scheduled"
    StatusConfirmed AppointmentStatus = "confirmed"
    StatusCompleted AppointmentStatus = "completed"
    StatusCancelled AppointmentStatus = "cancelled"
    StatusNoShow    AppointmentStatus = "no_show"
)
```

**Matriz de ações por status:**

| Status atual | Ações disponíveis |
|-------------|-------------------|
| `scheduled` | [Confirmar] [Cancelar] [Reagendar] |
| `confirmed`  | [Faltou] [Cancelar] [Reagendar] [Concluir] |
| `completed`  | — (somente visualização) |
| `cancelled`  | — (somente visualização) |
| `no_show`    | — (somente visualização) |

---

## Arquivos a modificar

### 1. `web/components/agenda/appointment_detail.templ`

Migrar para DaisyUI **e** adicionar ações condicionais por status.

**Estrutura esperada:**

```templ
templ AppointmentDetail(model AppointmentDetailModel) {
    <div class="modal modal-open" onclick="if(event.target===this)document.getElementById('modal-container').innerHTML=''">
        <div class="modal-box max-w-sm" onclick="event.stopPropagation()">

            <!-- Cabeçalho -->
            <div class="flex items-start gap-4 mb-5">
                <div class="avatar placeholder">
                    <div class="bg-primary text-primary-content rounded-full w-12">
                        <span class="text-lg font-bold">{ string([]rune(model.PatientName)[0:1]) }</span>
                    </div>
                </div>
                <div class="flex-1">
                    <h3 class="font-semibold text-base-content">{ model.PatientName }</h3>
                    <p class="text-sm text-base-content/60">{ model.SessionType }</p>
                </div>
                <button class="btn btn-ghost btn-sm btn-circle"
                    onclick="document.getElementById('modal-container').innerHTML=''">✕</button>
            </div>

            <!-- Badge de status -->
            @appointmentStatusBadge(model.Status)

            <!-- Infos: data, horário, duração -->
            <div class="space-y-2 my-5">
                ...
            </div>

            <!-- Ações condicionais -->
            @appointmentModalActions(model)

        </div>
    </div>
}

// Badge de status — DaisyUI
templ appointmentStatusBadge(status string) {
    switch status {
    case "scheduled":
        <div class="badge badge-ghost mb-4">Agendada</div>
    case "confirmed":
        <div class="badge badge-primary mb-4">Confirmada</div>
    case "completed":
        <div class="badge badge-success mb-4">Concluída</div>
    case "cancelled":
        <div class="badge badge-error badge-outline mb-4">Cancelada</div>
    case "no_show":
        <div class="badge badge-warning mb-4">Faltou</div>
    }
}

// Ações HTMX — só aparecem em status acionáveis
templ appointmentModalActions(model AppointmentDetailModel) {
    if model.Status == "scheduled" || model.Status == "confirmed" {
        <div class="modal-action flex-wrap gap-2 justify-start">
            if model.Status == "scheduled" {
                <button class="btn btn-primary btn-sm"
                    hx-post={ "/agenda/appointments/" + model.ID + "/confirm" }
                    hx-target="#modal-container" hx-swap="innerHTML">
                    <i class="fas fa-check"></i> Confirmar
                </button>
            }
            if model.Status == "confirmed" {
                <button class="btn btn-success btn-sm"
                    hx-post={ "/agenda/appointments/" + model.ID + "/complete" }
                    hx-target="#modal-container" hx-swap="innerHTML">
                    <i class="fas fa-check-double"></i> Concluir
                </button>
                <button class="btn btn-warning btn-sm"
                    hx-post={ "/agenda/appointments/" + model.ID + "/no-show" }
                    hx-target="#modal-container" hx-swap="innerHTML">
                    <i class="fas fa-user-times"></i> Faltou
                </button>
            }
            <button class="btn btn-ghost btn-sm"
                hx-get={ "/agenda/appointments/" + model.ID + "/reschedule" }
                hx-target="#drawer-container" hx-swap="innerHTML">
                <i class="fas fa-calendar-alt"></i> Reagendar
            </button>
            <button class="btn btn-error btn-outline btn-sm"
                hx-post={ "/agenda/appointments/" + model.ID + "/cancel" }
                hx-target="#modal-container" hx-swap="innerHTML"
                hx-confirm="Cancelar este agendamento?">
                <i class="fas fa-times"></i> Cancelar
            </button>
        </div>
    }
}
```

### 2. Handlers `Confirm`, `NoShow`, `Complete`, `Cancel` em `internal/web/handlers/agenda_handler.go`

Após executar a ação, o handler deve retornar o modal atualizado com o novo status.
Se os handlers já retornam redirect ou string vazia, adaptar para retornar fragmento HTMX:

```go
func (h *AgendaHandler) Confirm(w http.ResponseWriter, r *http.Request) {
    db, err := tenant.TenantDB(r.Context())
    if err != nil { http.Error(w, "unauthorized", 401); return }

    id := extractAppointmentID(r.URL.Path) // /agenda/appointments/{id}/confirm

    if err := h.service.ConfirmAppointment(r.Context(), id); err != nil {
        http.Error(w, err.Error(), 500); return
    }

    appt, err := h.service.GetAppointment(r.Context(), id)
    if err != nil { http.Error(w, err.Error(), 500); return }

    model := agendaComponents.AppointmentDetailModel{
        ID:          appt.ID,
        PatientName: appt.PatientName,
        Status:      string(appt.Status),
        // ... demais campos
    }
    agendaComponents.AppointmentDetail(model).Render(r.Context(), w)
}
```

Mesmo padrão para `NoShow`, `Complete`, `Cancel`.

### 3. `cmd/arandu/main.go`

Verificar e registrar se ausentes:
```
POST /agenda/appointments/{id}/confirm
POST /agenda/appointments/{id}/no-show
POST /agenda/appointments/{id}/complete
POST /agenda/appointments/{id}/cancel
```

---

## CSS — apenas DaisyUI, zero custom

Usar exclusivamente:
- `modal`, `modal-box`, `modal-action` — estrutura
- `btn btn-primary`, `btn btn-warning`, `btn btn-success`, `btn btn-error btn-outline`, `btn btn-ghost` — botões
- `badge badge-primary`, `badge badge-success`, etc. — status
- `avatar placeholder` — avatar com inicial
- `text-base-content`, `text-base-content/60` — cores de texto

**Não usar** `style=` inline, classes custom, `input-v2.css`, `tailwind-v2.css`.

---

## Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Comportamento — testar manualmente em http://localhost:8080/agenda**
- [ ] CA01: Agendamento `scheduled` → modal mostra [Confirmar] [Cancelar] [Reagendar]
- [ ] CA02: Agendamento `confirmed` → modal mostra [Concluir] [Faltou] [Cancelar] [Reagendar]
- [ ] CA03: Clicar [Confirmar] → modal atualiza com badge "Confirmada" sem reload de página
- [ ] CA04: Clicar [Faltou] → badge muda para "Faltou" (badge-warning)
- [ ] CA05: Agendamento `completed` ou `cancelled` → nenhum botão de ação
- [ ] CA06: Componente usa classes DaisyUI — sem `style=` inline

**Testes automatizados**
```go
// web/components/agenda/render_test.go — adicionar:
func TestAppointmentDetail_ScheduledShowsConfirmButton(t *testing.T)   // btn btn-primary + "Confirmar"
func TestAppointmentDetail_ConfirmedShowsCompleteAndNoShow(t *testing.T) // btn btn-success + btn btn-warning
func TestAppointmentDetail_CompletedShowsNoActions(t *testing.T)        // modal-action ausente
func TestAppointmentDetail_UsesDaisyUIModal(t *testing.T)               // "modal-box" no HTML
```

- [ ] `go test ./web/components/agenda/...` passa
- [ ] `go test ./internal/web/handlers/...` passa (sem regressão)

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa
- [ ] `./scripts/arandu_validate_handlers.sh` passa

---

## NÃO faça

- Não reimplementar `ConfirmAppointment`, `MarkNoShow` — service já tem
- Não criar CSS custom — somente DaisyUI + Tailwind utilities
- Não editar `input-v2.css` nem `tailwind-v2.css`
- Não retornar shell wrapper em requests HTMX — apenas fragmento do modal
- Não declarar concluída sem validar os 6 CAs manualmente no browser

---

## Skills de referência

- **daisyui-arandu** — componentes DaisyUI, tema Arandu, mapeamento de classes
- **arandu-architecture** — padrão HX-Request, containers do shell, multi-tenancy
- **go-templ-htmx-ux** — HTMX swap, fragmentos, targets

---

## Padrão de referência

- Handler similar: `internal/web/handlers/agenda_handler.go` → func `Cancel`
- Referência visual DaisyUI: `design_handoff_arandu_redesign/daisyui_shell_dashboard.html`
- Componente a migrar: `web/components/agenda/appointment_detail.templ`
