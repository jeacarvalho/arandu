# Task: Render tests para componentes de agenda
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

Um bug de `htmx-post` (typo) no `reschedule_form.templ` passou por revisão de código,
compilação e testes unitários sem ser detectado — só foi descoberto testando manualmente
no browser. Testes unitários verificam lógica Go; testes de render verificam o HTML
gerado pelo templ, incluindo atributos HTMX críticos.

Esta task adiciona render tests para os três componentes interativos da agenda. O padrão
já existe no projeto: `web/components/patient/profile_render_test.go` e
`web/components/timeline/timeline_render_test.go`.

### Padrão de render test (não inventar — seguir exatamente)

```go
package agenda_test  // sufixo _test, import externo

import (
    "bytes"
    "testing"
    "arandu/web/components/agenda"
)

func TestNome(t *testing.T) {
    var buf bytes.Buffer
    err := agenda.ComponentName(params).Render(t.Context(), &buf)
    if err != nil {
        t.Fatalf("render failed: %v", err)
    }
    html := buf.String()
    if !bytes.Contains(buf.Bytes(), []byte("expected")) {
        t.Errorf("expected HTML to contain %q\nHTML: %s", "expected", html)
    }
}
```

---

## O que implementar

**Arquivo único a criar:** `web/components/agenda/render_test.go`
**Pacote:** `agenda_test`

---

### Bloco 1 — AppointmentDetail: rotas e IDs críticos

Usar este modelo base para os testes (status `scheduled`, sem sessão vinculada):

```go
func baseDetailModel(id string) agenda.AppointmentDetailModel {
    return agenda.AppointmentDetailModel{
        ID:          id,
        PatientName: "Paciente Teste",
        Date:        time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
        StartTime:   "10:00",
        EndTime:     "10:50",
        Duration:    50,
        SessionType: "Sessão individual",
        Status:      "scheduled",
        SessionID:   nil,
    }
}
```

**Testes a implementar:**

```go
func TestAppointmentDetail_ReagendarButton_HasCorrectRoute(t *testing.T)
// Render com status="scheduled"
// Verifica: HTML contém hx-get="/agenda/appointments/test-id/reschedule-form"
// Verifica: HTML contém hx-target="#drawer-container"

func TestAppointmentDetail_CancelButton_HasCorrectRoute(t *testing.T)
// Render com status="scheduled"
// Verifica: HTML contém hx-post="/agenda/appointments/test-id/cancel"

func TestAppointmentDetail_ConcluirButton_HasConfirmPanel(t *testing.T)
// Render com status="scheduled", SessionID=nil
// Verifica: HTML contém id="confirm-panel-test-id"
// Verifica: HTML contém data-panel-id="test-id" (para o toggle JS funcionar)
// Verifica: HTML contém hx-post="/agenda/appointments/test-id/complete-with-session"
// Verifica: HTML contém hx-post="/agenda/appointments/test-id/complete"

func TestAppointmentDetail_WithSession_ShowsVerSessaoLink(t *testing.T)
// Render com SessionID = ptr("session-abc")
// Verifica: HTML contém /session/session-abc
// Verifica: HTML NÃO contém "complete-with-session" (botão Concluir não aparece quando há sessão)

func TestAppointmentDetail_CompletedStatus_NoMutationButtons(t *testing.T)
// Render com status="completed"
// Verifica: HTML NÃO contém hx-post="/agenda/appointments/test-id/cancel"
// Verifica: HTML NÃO contém hx-post="/agenda/appointments/test-id/complete"
```

Helper para ponteiro de string (adicionar no arquivo):
```go
func strPtr(s string) *string { return &s }
```

---

### Bloco 2 — RescheduleForm: atributos HTMX e inputs

```go
func TestRescheduleForm_FormHasCorrectHTMXPost(t *testing.T)
// Render com AppointmentID="appt-xyz"
// Verifica: HTML contém hx-post="/agenda/appointments/appt-xyz/reschedule"
// Verifica: HTML NÃO contém "htmx-post" (garante que o typo não volta)

func TestRescheduleForm_FormTarget(t *testing.T)
// Verifica: HTML contém hx-target="#agenda-content"
// Verifica: HTML contém hx-swap="outerHTML"

func TestRescheduleForm_HasRequiredInputs(t *testing.T)
// Verifica: HTML contém name="date"
// Verifica: HTML contém name="start_time"
// Verifica: HTML contém name="duration"

func TestRescheduleForm_PreFillsCurrentValues(t *testing.T)
// Render com CurrentDate=2026-04-22, CurrentStart="10:00", Duration=50
// Verifica: HTML contém value="2026-04-22"
// Verifica: HTML contém value="10:00"
// Verifica: HTML contém value="50" (option selecionada)
```

Modelo base:
```go
agenda.RescheduleFormModel{
    AppointmentID: "appt-xyz",
    CurrentDate:   time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
    CurrentStart:  "10:00",
    Duration:      50,
}
```

---

### Bloco 3 — NewAppointmentForm: atributos HTMX e conflict warning

```go
func TestNewAppointmentForm_FormHasCorrectHTMXPost(t *testing.T)
// Render com dados mínimos (Date preenchida, Slots e Patients vazios)
// Verifica: HTML contém hx-post="/agenda/appointments"
// Verifica: HTML contém hx-target="#agenda-content"

func TestNewAppointmentForm_ConflictWarningHiddenByDefault(t *testing.T)
// Verifica: HTML contém id="conflict-warning"
// Verifica: HTML contém "hidden" (classe aplicada por ConflictAlertClasses(false))

func TestNewAppointmentForm_HasRequiredInputs(t *testing.T)
// Verifica: HTML contém name="date"
// Verifica: HTML contém name="start_time"
// Verifica: HTML contém name="duration"
// Verifica: HTML contém name="patient_id"
```

Modelo base:
```go
agenda.NewAppointmentFormData{
    Date:     time.Date(2026, 4, 22, 0, 0, 0, 0, time.UTC),
    Slots:    nil,
    Patients: nil,
}
```

---

## Checklist de implementação

- [ ] Arquivo `web/components/agenda/render_test.go` criado
- [ ] Pacote `agenda_test` com imports `bytes`, `testing`, `time`, `arandu/web/components/agenda`
- [ ] Helper `strPtr` definido
- [ ] Helper `baseDetailModel` definido
- [ ] 5 testes de `AppointmentDetail` implementados e passando
- [ ] 4 testes de `RescheduleForm` implementados e passando
- [ ] 3 testes de `NewAppointmentForm` implementados e passando
- [ ] `TestRescheduleForm_FormHasCorrectHTMXPost` verifica AUSÊNCIA de `htmx-post`
- [ ] `go test ./web/components/agenda/...` — todos passando
- [ ] Nenhum arquivo de produção modificado

---

## Arquivo a criar

- `web/components/agenda/render_test.go`

## Arquivos a modificar

Nenhum.

---

## 🚫 NÃO faça

- Não modifique nenhum arquivo `.templ` ou `.go` de produção
- Não use `strings.Contains` — use `bytes.Contains(buf.Bytes(), []byte(...))` como o padrão do projeto
- Não teste CSS ou texto visível — teste apenas atributos funcionais (`hx-*`, `id`, `name`, `value`)
- Não crie helpers além de `strPtr` e `baseDetailModel`

---

## 📎 Padrão de referência

Siga `web/components/patient/profile_render_test.go` para estrutura geral.
Siga `web/components/agenda/types_test.go` para package name e imports do pacote agenda.
