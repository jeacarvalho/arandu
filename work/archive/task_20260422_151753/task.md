# Task: Reagendamento UI e link bidirecional sessão → agenda
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

Dois gaps restantes no ciclo agenda ↔ prontuário:

**A) Reagendamento:** o handler `Reschedule` existe e funciona, mas não há formulário UI
nem botão "Reagendar" no modal. O terapeuta hoje precisa cancelar + recriar para mudar
um horário. Também falta suporte HTMX no handler (retorna redirect puro).

**B) Link sessão → agenda:** a tabela `appointments` já tem `session_id`. O link de
agenda para sessão existe (modal mostra "Ver Sessão"). O inverso — da sessão para o
agendamento de origem — não existe. Solução: consulta reversa `WHERE session_id = ?`,
sem migration.

### Estado atual (não recriar)

```go
// Handler Reschedule — já existe, lê: date, start_time, duration
func (h *AgendaHandler) Reschedule(w http.ResponseWriter, r *http.Request)
// Retorna http.Redirect sem verificar HX-Request — precisa de HTMX fix

// RescheduleAppointment — já existe no AgendaService
func (s *AgendaService) RescheduleAppointment(ctx context.Context, id string,
    newDate time.Time, newStartTime, newEndTime string) error

// Session struct — sem AppointmentID
type Session struct { ID, PatientID string; Date time.Time; Summary string; ... }

// appointments tabela — tem coluna session_id (string, nullable)
// Consulta reversa possível: SELECT * FROM appointments WHERE session_id = ?
```

---

## Parte A — Reagendamento UI

### A1. Adicionar HTMX ao handler `Reschedule`

**Arquivo:** `internal/web/handlers/agenda_handler.go`

Antes da linha `http.Redirect(w, r, "/agenda", http.StatusSeeOther)`, adicionar:

```go
redirectURL := "/agenda"
if r.Header.Get("HX-Request") == "true" {
    w.Header().Set("HX-Redirect", redirectURL)
    w.WriteHeader(http.StatusOK)
    return
}
http.Redirect(w, r, redirectURL, http.StatusSeeOther)
```

Também adicionar tratamento do conflito 409 para HTMX (antes do redirect de sucesso):
```go
// No bloco de erro de conflito, também verificar HX-Request:
if strings.Contains(err.Error(), "conflicts") {
    w.WriteHeader(http.StatusConflict)
    w.Write([]byte("Horário ocupado. Escolha outro horário."))
    return
}
```

### A2. Novo handler `RescheduleForm` (GET)

**Arquivo:** `internal/web/handlers/agenda_handler.go` — adicionar ao final.

```go
// RescheduleForm handles GET /agenda/appointments/{id}/reschedule-form
// Retorna o drawer de reagendamento (fragmento HTMX).
func (h *AgendaHandler) RescheduleForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
    id = strings.TrimSuffix(id, "/reschedule-form")

    ctx := r.Context()
    appt, err := h.agendaService.GetAppointment(ctx, id)
    if err != nil || appt == nil {
        http.NotFound(w, r)
        return
    }

    vm := agendaComponents.RescheduleFormModel{
        AppointmentID: appt.ID,
        CurrentDate:   appt.Date,
        CurrentStart:  appt.StartTime,
        Duration:      appt.Duration,
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    agendaComponents.RescheduleForm(vm).Render(ctx, w)
}
```

### A3. ViewModel e componente templ de reagendamento

**Arquivo:** `web/components/agenda/types.go` — adicionar:

```go
type RescheduleFormModel struct {
    AppointmentID string
    CurrentDate   time.Time
    CurrentStart  string
    Duration      int
}
```

**Arquivo:** `web/components/agenda/reschedule_form.templ` (novo)

Drawer lateral com formulário de reagendamento — siga o padrão de `new_appointment_form.templ`
(mesmo `id="drawer"`, mesmo header, mesma estrutura de botões Cancelar/Confirmar):

```templ
package agenda

templ RescheduleForm(model RescheduleFormModel) {
    <div id="drawer" class="fixed inset-y-0 right-0 w-96 bg-white shadow-xl z-50 flex flex-col">
        <!-- Header -->
        <div class="flex items-center justify-between px-5 py-4 border-b border-neutral-200 bg-neutral-50 flex-shrink-0">
            <h3 class="text-base font-semibold text-neutral-800">Reagendar Consulta</h3>
            <button onclick="document.getElementById('drawer-container').innerHTML=''"
                class="w-8 h-8 flex items-center justify-center rounded-lg text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 transition-colors">
                <i class="fas fa-times"></i>
            </button>
        </div>

        <!-- Form -->
        <form
            hx-post={ "/agenda/appointments/" + model.AppointmentID + "/reschedule" }
            hx-target="#agenda-content"
            hx-swap="outerHTML"
            class="flex-1 overflow-y-auto p-5 space-y-4"
        >
            <!-- Date -->
            <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1.5">Nova data</label>
                <input type="date" name="date"
                    value={ model.CurrentDate.Format("2006-01-02") }
                    class="w-full h-10 px-3 border border-neutral-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-arandu-primary/50"
                    required />
            </div>

            <!-- Start time -->
            <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1.5">Novo horário</label>
                <input type="time" name="start_time"
                    value={ model.CurrentStart }
                    class="w-full h-10 px-3 border border-neutral-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-arandu-primary/50"
                    required />
            </div>

            <!-- Duration -->
            <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1.5">Duração</label>
                <select name="duration"
                    class="w-full h-10 px-3 border border-neutral-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-arandu-primary/50">
                    <option value="30" selected?={ model.Duration == 30 }>30 min</option>
                    <option value="50" selected?={ model.Duration == 50 }>50 min</option>
                    <option value="60" selected?={ model.Duration == 60 }>60 min</option>
                </select>
            </div>

            <!-- Conflict warning -->
            <div id="reschedule-conflict" class={ ConflictAlertClasses(false) }>
                <i class="fas fa-exclamation-triangle mr-1.5"></i>
                Horário ocupado. Escolha outro horário.
            </div>

            <!-- Actions -->
            <div class="flex gap-3 pt-4">
                <button type="button"
                    onclick="document.getElementById('drawer-container').innerHTML=''"
                    class="flex-1 h-10 text-sm font-medium text-neutral-700 bg-neutral-100 hover:bg-neutral-200 rounded-lg transition-colors">
                    Cancelar
                </button>
                <button type="submit"
                    class="flex-1 h-10 text-sm font-medium text-white bg-arandu-primary hover:bg-arandu-dark rounded-lg transition-colors">
                    Confirmar
                </button>
            </div>
        </form>
    </div>
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black/30 z-40"
        onclick="document.getElementById('drawer-container').innerHTML=''"></div>
}
```

> **Nota:** `ConflictAlertClasses` já existe em `types.go` — reutilize sem recriar.

### A4. Botão "Reagendar" no modal de detalhe

**Arquivo:** `web/components/agenda/appointment_detail.templ`

Na área de ações (bloco `<div class="px-6 pb-6 flex gap-2">`), adicionar botão de
reagendamento para agendamentos com status `scheduled` ou `confirmed`:

```templ
if model.Status == "scheduled" || model.Status == "confirmed" {
    <button
        hx-get={ "/agenda/appointments/" + model.ID + "/reschedule-form" }
        hx-target="#drawer-container"
        hx-swap="innerHTML"
        onclick="document.getElementById('modal-container').innerHTML=''"
        class="h-10 px-4 text-sm font-medium rounded-xl border transition-colors"
        style="color:#374151;border-color:#d1d5db"
        onmouseover="this.style.background='#f9fafb'"
        onmouseout="this.style.background=''"
    >
        <i class="fas fa-calendar-alt text-xs mr-1"></i>Reagendar
    </button>
    // ... botões Cancel e Concluir existentes ...
}
```

### A5. Registrar nova rota em `main.go`

No bloco `mux.HandleFunc("/agenda/appointments/", ...)`, adicionar:

```go
case strings.HasSuffix(path, "/reschedule-form") && r.Method == http.MethodGet:
    agendaHandler.RescheduleForm(w, r)
```

---

## Parte B — Link sessão → agendamento de origem

### B1. Novo método no AppointmentRepository

**Arquivo:** `internal/infrastructure/repository/sqlite/appointment_repository.go`

```go
// FindBySessionID retorna o agendamento vinculado a uma sessão clínica, se existir.
func (r *AppointmentRepository) FindBySessionID(ctx context.Context, sessionID string) (*appointment.Appointment, error) {
    // Mesma query base do FindByID — apenas WHERE diferente
    // WHERE session_id = ?
    // Retornar nil, nil se não encontrado (não é erro)
}
```

Siga o padrão exato de `FindByID` do mesmo arquivo (mesmo scan de colunas).

O `ContextAwareAppointmentRepository` em `context_wrapper.go` também precisa expor
`FindBySessionID` — siga o padrão de delegate dos outros métodos.

### B2. Expor no AgendaService

**Arquivo:** `internal/application/services/agenda_service.go`

```go
func (s *AgendaService) GetAppointmentBySession(ctx context.Context, sessionID string) (*appointment.Appointment, error) {
    return s.apptRepo.FindBySessionID(ctx, sessionID)
}
```

### B3. Injetar AgendaService no SessionHandler

**Arquivo:** `internal/web/handlers/session_handler.go`

Verificar se `SessionHandler` já tem acesso a algum service de agenda. Se não:

```go
// Adicionar interface mínima
type AgendaLookupPort interface {
    GetAppointmentBySession(ctx context.Context, sessionID string) (*appointment.Appointment, error)
}

// Adicionar campo ao struct
type SessionHandler struct {
    // ... campos existentes ...
    agendaService AgendaLookupPort  // novo — último campo
}
```

No handler `Show` (GET /session/{id}), após carregar a sessão:

```go
originAppt, _ := h.agendaService.GetAppointmentBySession(ctx, session.ID)
// Ignorar erro — não deve quebrar a sessão se não houver agendamento vinculado
```

Passar ao ViewModel:

```go
// Se originAppt != nil:
vm.OriginDate  = originAppt.Date.Format("02 jan 2006")
vm.OriginTime  = originAppt.StartTime
vm.OriginApptID = originAppt.ID
vm.HasOriginAppt = true
```

### B4. ViewModel e template da sessão

**Arquivo:** `web/components/session/types.go` (ou equivalente) — adicionar campos ao
ViewModel da sessão:

```go
HasOriginAppt bool
OriginDate    string
OriginTime    string
OriginApptID  string
```

**Arquivo:** `web/components/session/detail.templ` — no header, após o botão "Voltar",
adicionar linha discreta quando `HasOriginAppt == true`:

```templ
if data.HasOriginAppt {
    <a href={ templ.URL("/agenda?view=dia&date=" + data.OriginDate) }
       class="text-xs text-neutral-400 hover:text-arandu-primary transition-colors">
        <i class="fas fa-calendar-check mr-1"></i>
        Agendada { data.OriginDate } às { data.OriginTime }
    </a>
}
```

### B5. Atualizar `main.go`

Passar `agendaService` ao `sessionHandler` (mesmo objeto já criado para agendaHandler).

---

## Checklist de implementação

**Parte A — Reagendamento:**
- [ ] Handler `Reschedule` retorna `HX-Redirect` para HTMX requests
- [ ] Handler `RescheduleForm` (GET) criado
- [ ] `RescheduleFormModel` em `types.go`
- [ ] `reschedule_form.templ` criado seguindo padrão do `new_appointment_form.templ`
- [ ] Botão "Reagendar" adicionado ao `appointment_detail.templ`
- [ ] Rota `/reschedule-form` registrada em `main.go`

**Parte B — Link sessão → agenda:**
- [ ] `AppointmentRepository.FindBySessionID` implementado
- [ ] `ContextAwareAppointmentRepository.FindBySessionID` delegado
- [ ] `AgendaService.GetAppointmentBySession` exposto
- [ ] `SessionHandler` injeta `agendaService` e carrega `originAppt`
- [ ] ViewModel da sessão tem campos `HasOriginAppt`, `OriginDate`, `OriginTime`
- [ ] `detail.templ` exibe link de origem quando `HasOriginAppt == true`
- [ ] `main.go` passa `agendaService` ao `sessionHandler`

**Qualidade:**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build ./cmd/arandu/` sem erros
- [ ] `go test ./internal/...` continua passando

---

## Arquivos a criar

- `web/components/agenda/reschedule_form.templ`

## Arquivos a modificar

- `internal/web/handlers/agenda_handler.go`
- `internal/web/handlers/session_handler.go`
- `internal/infrastructure/repository/sqlite/appointment_repository.go`
- `internal/infrastructure/repository/sqlite/context_wrapper.go`
- `internal/application/services/agenda_service.go`
- `web/components/agenda/types.go`
- `web/components/agenda/appointment_detail.templ`
- `web/components/session/detail.templ` (e types.go equivalente)
- `cmd/arandu/main.go`

---

## 🔒 Privacidade

- [x] **Tier 1**: `PatientName` presente no `Appointment` — não logar em nenhum handler novo.
- [x] Consulta reversa (`FindBySessionID`) é sempre scoped ao tenant via context.

---

## 🚫 NÃO faça

- Não crie migration — `session_id` já existe na tabela `appointments`
- Não modifique `RescheduleAppointment` no service — ele já funciona
- Não quebre o handler `Show` de sessão se `GetAppointmentBySession` retornar nil
- Não use `html/template` — apenas `.templ`
- O link "Agendada em..." na sessão deve ser **discreto** — metadado, não destaque

---

## 📎 Padrão de referência

- `new_appointment_form.templ` — estrutura exata do drawer para `reschedule_form.templ`
- `CompleteWithSession` em `agenda_handler.go` — padrão de injeção de service e HX-Redirect
- `FindByID` em `appointment_repository.go` — padrão de query + scan para `FindBySessionID`
- `context_wrapper.go` — padrão de delegate para expor o novo método
