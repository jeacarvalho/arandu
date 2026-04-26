# Task: Agenda — testes unitários e de integração
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

A agenda está funcionalmente completa. O objetivo desta task é fechar o REQ-07-01-01 adicionando cobertura de testes nas três camadas onde ela ainda não existe.

### Estado atual dos testes de agenda

| Arquivo | Estado | Testes |
|---------|--------|--------|
| `internal/web/handlers/agenda_handler_test.go` | ✅ Existe | 15 testes, todos passando |
| `web/components/agenda/types_test.go` | ✅ Existe | 7 testes, todos passando |
| `internal/domain/appointment/appointment_test.go` | ❌ Ausente | 0 |
| `internal/application/services/agenda_service_test.go` | ❌ Ausente | 0 |
| `internal/infrastructure/repository/sqlite/appointment_repository_test.go` | ❌ Ausente | 0 |

### Gaps no handler existente

O `agenda_handler_test.go` não cobre o comportamento HTMX dos handlers `Cancel` e `Complete` corrigidos na task anterior. Os testes existentes para esses handlers só verificam que a resposta tem código não-zero — não verificam o `HX-Redirect`.

---

## O que implementar

### 1. Testes de domínio (novo arquivo)

**Arquivo:** `internal/domain/appointment/appointment_test.go`
**Pacote:** `package appointment`

Cobrir os comportamentos de negócio da entidade (sem mocks — testam a entidade diretamente):

```go
func TestNewAppointment_ValidSession(t *testing.T)         // cria com status "scheduled"
func TestNewAppointment_InvalidDuration(t *testing.T)      // duration < 30 e > 120 → ErrInvalidDuration
func TestNewAppointment_MissingPatient(t *testing.T)       // type=session, patientID="" → ErrPatientRequired
func TestNewAppointment_BlockedSlotNoPatient(t *testing.T) // type=blocked aceita patientID=""

func TestAppointment_Cancel_SetsStatus(t *testing.T)        // status → "cancelled", CancelledAt != nil
func TestAppointment_Cancel_BlocksCompleted(t *testing.T)   // completed → Cancel() → ErrAlreadyCompleted
func TestAppointment_Complete_SetsStatus(t *testing.T)      // status → "completed", CompletedAt != nil
func TestAppointment_Complete_BlocksCancelled(t *testing.T) // cancelled → Complete() → ErrAlreadyCancelled
func TestAppointment_MarkNoShow(t *testing.T)               // status → "no_show"

func TestAppointment_Overlaps_SameTimeTrue(t *testing.T)        // mesmo dia e horário → true
func TestAppointment_Overlaps_AdjacentFalse(t *testing.T)       // 10:00-10:50 e 10:50-11:40 → false
func TestAppointment_Overlaps_DifferentDateFalse(t *testing.T)  // dias diferentes → false
```

---

### 2. Testes de serviço (novo arquivo)

**Arquivo:** `internal/application/services/agenda_service_test.go`
**Pacote:** `package services`

**Padrão de mock** — use o estilo function-based de `patient_service_test.go` (já no mesmo pacote). Crie:

```go
type mockAgendaRepo struct {
    saveFunc            func(ctx context.Context, appt *appointment.Appointment) error
    findByIDFunc        func(ctx context.Context, id string) (*appointment.Appointment, error)
    findByDateRangeFunc func(ctx context.Context, start, end time.Time) ([]*appointment.Appointment, error)
    findByDateFunc      func(ctx context.Context, date time.Time) ([]*appointment.Appointment, error)
    findOverlappingFunc func(ctx context.Context, date time.Time, start, end, excludeID string) ([]*appointment.Appointment, error)
    updateFunc          func(ctx context.Context, appt *appointment.Appointment) error
    deleteFunc          func(ctx context.Context, id string) error
    findUpcomingFunc    func(ctx context.Context, from time.Time, limit int) ([]*appointment.Appointment, error)
    findByPatientFunc   func(ctx context.Context, patientID string) ([]*appointment.Appointment, error)
}
// Cada método: se o funcField != nil, chama; senão retorna nil, nil (ou zero value).
```

Testes a cobrir:

```go
// GetWeekView
func TestAgendaService_GetWeekView_StartOnMonday(t *testing.T)
// Chama GetWeekView com uma quarta-feira (ex: 2026-04-22).
// Verifica: weekView.StartDate.Weekday() == time.Monday
// Verifica: len(weekView.Days) == 7

func TestAgendaService_GetWeekView_GroupsAppointmentsByDay(t *testing.T)
// findByDateRangeFunc retorna 2 appointments na segunda e 1 na quarta da mesma semana.
// Verifica que days[0] tem 2 e days[2] tem 1 (use filterByDate interno via GetWeekView).

// GetMonthView
func TestAgendaService_GetMonthView_FullWeeks(t *testing.T)
// Chama GetMonthView(2026, 4) — abril começa na quarta.
// Verifica: len(monthView.Days) é múltiplo de 7.
// Verifica: monthView.Days[0].Date.Weekday() == time.Monday.

// CreateAppointment — conflict detection
func TestAgendaService_Create_DetectsConflict(t *testing.T)
// findOverlappingFunc retorna 1 appointment existente.
// Chama CreateAppointment → deve retornar erro contendo "conflicts".

func TestAgendaService_Create_Succeeds(t *testing.T)
// findOverlappingFunc retorna slice vazio.
// saveFunc captura o appointment.
// Verifica: salvo.Status == AppointmentStatusScheduled.

// CancelAppointment
func TestAgendaService_Cancel_SetsStatusCancelled(t *testing.T)
// findByIDFunc retorna appointment com status "scheduled".
// updateFunc captura o appointment atualizado.
// Verifica: status == AppointmentStatusCancelled.

// CompleteAppointment
func TestAgendaService_Complete_SetsStatusCompleted(t *testing.T)
// findByIDFunc retorna appointment "confirmed".
// Verifica: status == AppointmentStatusCompleted após CompleteAppointment.

func TestAgendaService_Complete_LinksSession(t *testing.T)
// Chama CompleteAppointment(id, "session-xyz").
// Verifica: appointment.SessionID == "session-xyz".
```

---

### 3. Testes de repositório (novo arquivo)

**Arquivo:** `internal/infrastructure/repository/sqlite/appointment_repository_test.go`
**Pacote:** `package sqlite`

**Padrão** — siga `session_repository_test.go` exatamente:

```go
func setupAppointmentTestDB(t *testing.T) (*DB, func()) {
    tmpfile, err := os.CreateTemp("", "testdb-appt-*.db")
    // NewDB(tmpfile.Name()), db.Migrate(), cleanup: db.Close() + os.Remove
}
```

> **Atenção (skill: arandu-architecture § Multi-tenancy):** O `AppointmentRepository` usa o context wrapper para extrair o DB. Antes de implementar, leia `context_wrapper.go` e a assinatura de `NewAppointmentRepository`. Se o construtor aceitar `*DB` diretamente, use-o. Caso contrário, injete via `context.WithValue` com a mesma key que o middleware usa.

Testes:

```go
func TestAppointmentRepository_SaveAndFindByID(t *testing.T)
// Save appointment, FindByID → verifica ID, PatientName, Status.

func TestAppointmentRepository_FindByDateRange(t *testing.T)
// Cria 3 appointments em datas diferentes (dia-1, dia, dia+1).
// FindByDateRange(dia, dia) → retorna só o do dia.

func TestAppointmentRepository_FindOverlapping_Conflict(t *testing.T)
// Save appointment 10:00-10:50 no dia X.
// FindOverlapping(X, "10:20", "11:10", "") → retorna 1 (conflito).
// FindOverlapping(X, "11:00", "11:50", "") → retorna 0 (sem conflito).

func TestAppointmentRepository_Update_StatusChange(t *testing.T)
// Save com status "scheduled".
// Chama Cancel() no domínio → Update → FindByID.
// Verifica status == "cancelled".

func TestAppointmentRepository_FindByDate(t *testing.T)
// Save em 2 dias distintos. FindByDate(dia1) → só retorna do dia1.
```

---

### 4. Gaps no handler existente

**Arquivo:** `internal/web/handlers/agenda_handler_test.go` — adicionar ao final do arquivo.

O `mockAppointmentRepository` e `mockPatientService` já existem no arquivo — reutilize-os.

```go
func TestAgendaHandler_Cancel_HTMX(t *testing.T)
// Cria appointment real via agendaService.CreateAppointment.
// Envia POST /agenda/appointments/{id}/cancel com HX-Request: true.
// Verifica: w.Code == 200 e w.Header().Get("HX-Redirect") != "".

func TestAgendaHandler_Cancel_NonHTMX(t *testing.T)
// Mesma criação.
// Envia POST sem HX-Request.
// Verifica: w.Code == 303 e Location == "/agenda".

func TestAgendaHandler_Complete_HTMX(t *testing.T)
// Cria appointment via agendaService.CreateAppointment.
// Envia POST /agenda/appointments/{id}/complete com HX-Request: true.
// Verifica: w.Code == 200 e HX-Redirect contém "/agenda?view=dia&date=".
```

> **Como obter o ID:** `agendaService.CreateAppointment(...)` retorna `*appointment.Appointment` com o ID. Use `appt.ID` diretamente — não extraia de URL.

---

### 5. Adição em types_test.go

**Arquivo:** `web/components/agenda/types_test.go` — adicionar ao final.

```go
func TestConflictAlertClasses_HiddenByDefault(t *testing.T)
// ConflictAlertClasses(false) == "hidden"

func TestConflictAlertClasses_VisibleHasErrorClasses(t *testing.T)
// ConflictAlertClasses(true) deve conter "bg-red-50" e "text-red-700"
```

---

## Ordem de implementação

1. Domínio (`appointment_test.go`) — sem dependências, feedback imediato
2. Serviço (`agenda_service_test.go`) — mock puro
3. Repositório (`appointment_repository_test.go`) — SQLite real
4. Handler gaps + types additions

## Checklist

- [ ] `internal/domain/appointment/appointment_test.go` criado, todos os testes passando
- [ ] `internal/application/services/agenda_service_test.go` criado, todos os testes passando
- [ ] `internal/infrastructure/repository/sqlite/appointment_repository_test.go` criado, todos os testes passando
- [ ] 3 testes adicionados ao `agenda_handler_test.go`, passando
- [ ] 2 testes adicionados ao `types_test.go`, passando
- [ ] `go test ./internal/domain/appointment/... ./internal/application/services/... ./internal/infrastructure/repository/sqlite/... ./internal/web/handlers/... ./web/components/agenda/...` — todos passam
- [ ] Nenhum arquivo de produção modificado

## Arquivos a criar

- `internal/domain/appointment/appointment_test.go`
- `internal/application/services/agenda_service_test.go`
- `internal/infrastructure/repository/sqlite/appointment_repository_test.go`

## Arquivos a modificar (append only)

- `internal/web/handlers/agenda_handler_test.go`
- `web/components/agenda/types_test.go`

## Skills de referência obrigatória

- **tdd-go** — estrutura Go: `t.Run` para sub-testes, `errors.Is` para tipos de erro, table-driven com `[]struct{ name, ... }`; ao testar comportamentos de estado, prefira testes focados (um comportamento por teste) a mega-testes
- **ddd-go** — domínio: testa entidade direta; serviço: mock de repository; repositório: DB real (sem mock) — cada camada tem estilo próprio
- **arandu-architecture § Multi-tenancy** — antes de implementar o repo test, leia `context_wrapper.go` para saber como o DB é extraído do context nos repositories
- **clinical-domain** — use dados fictícios (`"Paciente Teste"`, `"paciente-teste-id"`) — nunca dados reais de pacientes em testes
