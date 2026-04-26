# Task: Histórico de agendamentos no perfil do paciente
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

O perfil do paciente (`/patient/{id}`) mostra sessões clínicas mas nenhuma informação
de agendamento. O terapeuta não consegue ver, ao consultar um paciente, quantas consultas
foram feitas, quando foi a última, quais foram canceladas.

O repositório já tem `FindByPatient(ctx, patientID)` implementado. Faltam três coisas:
um método no `AgendaService`, a injeção do service no `PatientHandler`, e uma nova seção
no `profile.templ`.

### Estado atual (não recriar)

```go
// AppointmentRepository — método já existe
func (r *AppointmentRepository) FindByPatient(ctx context.Context, patientID string) ([]*appointment.Appointment, error)

// PatientHandler — struct atual (sem agendaService)
type PatientHandler struct {
    patientService         PatientService
    sessionService         SessionService
    insightService         InsightService
    biopsychosocialService BiopsychosocialService
    timelineService        TimelineServicePort
    anamnesisService       AnamnesisService
}

// PatientViewData — struct atual passada ao template
// Campos existentes: Patient, Sessions, Insights, Error
// Campo a adicionar: Appointments []AppointmentHistoryItem

// Appointment domain type — campos relevantes
type Appointment struct {
    ID          string
    Date        time.Time
    StartTime   string        // "HH:MM"
    Duration    int           // minutos
    Status      AppointmentStatus
    SessionID   *string       // nil se não vinculado
}
// Status constants: scheduled, confirmed, completed, cancelled, no_show
```

---

## O que implementar

### 1. Novo método no AgendaService

**Arquivo:** `internal/application/services/agenda_service.go`

```go
// GetPatientAppointments retorna todos os agendamentos de um paciente, ordenados
// por data descendente (mais recente primeiro).
func (s *AgendaService) GetPatientAppointments(ctx context.Context, patientID string) ([]*appointment.Appointment, error) {
    appts, err := s.apptRepo.FindByPatient(ctx, patientID)
    if err != nil {
        return nil, err
    }
    // Ordena por Date descendente
    sort.Slice(appts, func(i, j int) bool {
        return appts[i].Date.After(appts[j].Date)
    })
    return appts, nil
}
```

> Adicione `"sort"` aos imports se necessário.

### 2. ViewModel no paciente

**Arquivo:** `web/components/patient/types.go`

Adicionar ao final do arquivo:

```go
// AppointmentHistoryItem representa um agendamento no perfil do paciente.
// StatusClass e StatusLabel são pré-computados no handler — o template não faz lógica.
type AppointmentHistoryItem struct {
    ID           string
    Date         string // "22 abr 2026"
    StartTime    string // "10:00"
    Duration     int
    StatusLabel  string // "Agendada" | "Confirmada" | "Realizada" | "Cancelada" | "Não Compareceu"
    StatusClass  string // classes CSS completas para o badge
    HasSession   bool
    SessionID    string // vazio se HasSession == false
}
```

Classes CSS para `StatusClass` (retornar string completa — Tailwind não aceita concatenação):

```go
// Função auxiliar — adicionar em types.go
func AppointmentStatusBadgeClass(status string) string {
    base := "inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
    switch status {
    case "scheduled":
        return base + " bg-amber-100 text-amber-800"
    case "confirmed":
        return base + " bg-emerald-100 text-emerald-800"
    case "completed":
        return base + " bg-arandu-primary/10 text-arandu-primary"
    case "cancelled":
        return base + " bg-neutral-100 text-neutral-500"
    case "no_show":
        return base + " bg-red-100 text-red-700"
    default:
        return base + " bg-neutral-100 text-neutral-600"
    }
}
```

### 3. Injetar AgendaService no PatientHandler

**Arquivo:** `internal/web/handlers/patient_handler.go`

**a)** Declare a interface mínima no topo do arquivo (junto às outras interfaces):

```go
type AgendaServicePort interface {
    GetPatientAppointments(ctx context.Context, patientID string) ([]*appointment.Appointment, error)
}
```

> Import necessário: `"arandu/internal/domain/appointment"`

**b)** Adicione o campo à struct e ao construtor:

```go
type PatientHandler struct {
    // ... campos existentes ...
    agendaService AgendaServicePort  // novo campo — último da lista
}

// NewPatientHandler — adicionar agendaService como último parâmetro
```

**c)** No handler `Show` (GET /patient/{id}), após carregar as sessões, adicione:

```go
appts, err := h.agendaService.GetPatientAppointments(ctx, patientID)
if err != nil {
    // não falhe a página por isso — apenas logue e continue com lista vazia
    appts = nil
}

// Mapear para ViewModel — limite de 20 registros
historyItems := make([]patientComponents.AppointmentHistoryItem, 0, len(appts))
for i, a := range appts {
    if i >= 20 {
        break
    }
    sessionID := ""
    if a.SessionID != nil {
        sessionID = *a.SessionID
    }
    historyItems = append(historyItems, patientComponents.AppointmentHistoryItem{
        ID:          a.ID,
        Date:        a.Date.Format("02 jan 2006"),
        StartTime:   a.StartTime,
        Duration:    a.Duration,
        StatusLabel: statusLabelForAppointment(string(a.Status)),
        StatusClass: patientComponents.AppointmentStatusBadgeClass(string(a.Status)),
        HasSession:  sessionID != "",
        SessionID:   sessionID,
    })
}
```

Função auxiliar (adicionar no mesmo arquivo):

```go
func statusLabelForAppointment(status string) string {
    switch status {
    case "scheduled":  return "Agendada"
    case "confirmed":  return "Confirmada"
    case "completed":  return "Realizada"
    case "cancelled":  return "Cancelada"
    case "no_show":    return "Não Compareceu"
    default:           return status
    }
}
```

**d)** Adicione `Appointments` ao `PatientViewData` e passe `historyItems` ao renderizar.

### 4. Seção no template de perfil

**Arquivo:** `web/components/patient/profile.templ`

Adicionar nova seção após a lista de sessões recentes. Siga o padrão visual das seções
existentes no mesmo arquivo (card branco com borda, título de seção, lista).

```templ
<!-- Histórico de Agendamentos -->
<div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
    <div class="px-5 py-3.5 border-b border-neutral-100 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-neutral-800">Agendamentos</h3>
        <a href={ templ.URL("/agenda?patient=" + data.Patient.ID) }
           class="text-xs text-arandu-primary hover:underline">
            Ver na agenda
        </a>
    </div>
    if len(data.Appointments) == 0 {
        <div class="px-5 py-8 text-center text-sm text-neutral-400">
            Nenhum agendamento registrado
        </div>
    } else {
        <ul class="divide-y divide-neutral-50">
            for _, appt := range data.Appointments {
                <li class="px-5 py-3 flex items-center gap-3">
                    <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2">
                            <span class="text-sm font-medium text-neutral-800">
                                { appt.Date }
                            </span>
                            <span class="text-xs text-neutral-400">{ appt.StartTime }</span>
                            <span class={ appt.StatusClass }>{ appt.StatusLabel }</span>
                        </div>
                    </div>
                    if appt.HasSession {
                        <a href={ templ.URL("/session/" + appt.SessionID) }
                           class="text-xs text-arandu-primary hover:underline flex-shrink-0">
                            Ver sessão →
                        </a>
                    }
                </li>
            }
        </ul>
    }
</div>
```

> **Atenção:** verifique o nome exato do campo `data.Patient.ID` e `data.Appointments`
> no ViewModel real do `profile.templ` antes de usar — adapte se necessário.

### 5. Atualizar `main.go`

**Arquivo:** `cmd/arandu/main.go`

Passar `agendaService` ao construtor do `patientHandler`:

```go
// Antes:
patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, ...)

// Depois — adicionar agendaService como último argumento:
patientHandler := handlers.NewPatientHandler(patientServiceAdapter, sessionServiceAdapter, ..., agendaService)
```

> `agendaService` já é criado mais acima em `main.go` para o `agendaHandler` —
> reutilize a mesma variável.

---

## Checklist de implementação

- [ ] `AgendaService.GetPatientAppointments` implementado com ordenação descendente
- [ ] `AppointmentHistoryItem` e `AppointmentStatusBadgeClass` em `types.go`
- [ ] `AgendaServicePort` interface declarada em `patient_handler.go`
- [ ] `PatientHandler` tem campo `agendaService` e construtor atualizado
- [ ] Handler `Show` carrega agendamentos, mapeia para ViewModel, limita a 20
- [ ] `PatientViewData` tem campo `Appointments []AppointmentHistoryItem`
- [ ] `profile.templ` tem seção de agendamentos com badge de status e link para sessão
- [ ] `main.go` passa `agendaService` ao `patientHandler`
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build ./cmd/arandu/` sem erros
- [ ] `go test ./internal/application/services/...` continua passando

---

## Arquivos a modificar (nenhum arquivo novo)

- `internal/application/services/agenda_service.go`
- `internal/web/handlers/patient_handler.go`
- `web/components/patient/types.go`
- `web/components/patient/profile.templ`
- `cmd/arandu/main.go`

---

## 🔒 Privacidade

- [x] **Tier 1 — PII indireto**: `PatientName` já está no `Appointment` mas **não** deve
  aparecer em nenhum log do handler. Não adicione `log.Printf` com dados do agendamento.
- [x] A lista é sempre scoped ao tenant do request via context — nunca mistura pacientes
  de tenants diferentes.

---

## 🚫 NÃO faça

- Não crie migration — nenhuma mudança de schema
- Não use `html/template` — apenas `.templ`
- Não busque todos os agendamentos sem limite — máximo 20 no handler
- Não quebre o handler `Show` se `GetPatientAppointments` falhar — degrade gracefully
- Não importe `web/components/agenda` dentro de `web/components/patient` — as classes CSS
  de badge devem ficar em `patient/types.go` (já especificadas acima)
- Não altere `FindByPatient` no repositório

---

## 📎 Padrão de referência

- Como o `AgendaHandler` recebeu `sessionService` na task anterior: mesmo padrão para
  injetar `agendaService` no `PatientHandler`
- Seção de sessões recentes em `profile.templ` — siga o mesmo padrão visual para a
  seção de agendamentos
- `PatientViewData` em `patient_handler.go` — veja como `Sessions` é adicionado ao struct
  para replicar com `Appointments`
