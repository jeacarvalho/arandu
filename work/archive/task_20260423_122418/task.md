# Task: CA-04 — Confirmar Agendamento e Registrar Falta (HTMX inline)
Requirement: REQ-07-01-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

A agenda do Arandu já possui os fluxos de **Cancelar** e **Concluir** agendamentos.
Faltam dois fluxos de uso diário do terapeuta:

1. **Confirmar** — antes da sessão, o terapeuta confirma que o paciente vai comparecer
   (`scheduled → confirmed`).
2. **Registrar Falta** — após o horário, o terapeuta marca que o paciente não compareceu
   (`scheduled | confirmed → no_show`).

### O que já existe (não reescreva)

**Domínio** (`internal/domain/appointment/appointment.go`):
- `AppointmentStatus` enum completo: `scheduled`, `confirmed`, `completed`, `cancelled`, `no_show`
- `Appointment.Confirm() error` — regra de domínio pronta
- `Appointment.MarkNoShow() error` — regra de domínio pronta

**Service** (`internal/application/services/agenda_service.go`):
- `MarkNoShow(ctx context.Context, id string) error` — **já existe**, apenas falta handler
- `ConfirmAppointment` — **NÃO EXISTE**, precisa ser criado

**Handler** (`internal/web/handlers/agenda_handler.go`):
- `Cancel`, `Complete`, `CompleteWithSession` — existem como referência de padrão
- `Confirm`, `NoShow` — **precisam ser criados**

**Rotas registradas** (`cmd/arandu/main.go`):
- `POST /agenda/appointments/{id}/cancel` → `Cancel`
- `POST /agenda/appointments/{id}/complete` → `Complete`
- `POST /agenda/appointments/{id}/complete-with-session` → `CompleteWithSession`
- `/confirm` e `/no-show` — **não registradas ainda**

**UI** (`web/components/agenda/appointment_detail.templ`):
- Botões condicionais existentes: Cancelar, Reagendar, Concluir
- Botões **Confirmar** e **Faltou** — ausentes, precisam ser adicionados

**Helpers de status** (`web/components/agenda/types.go`):
- `StatusLabel(status string) string` — traduz status para PT-BR
- `GetAppointmentStatusClass(status string) string` — retorna classes Tailwind do badge
- Adicionar helpers se necessário; nunca interpolação dinâmica de classes no template

---

## O que implementar

### Parte A — Service: ConfirmAppointment

Em `internal/application/services/agenda_service.go`:

1. Adicionar `ConfirmAppointment(ctx context.Context, id string) error` à interface
   `AgendaServiceInterface`
2. Implementar:
   ```go
   func (s *AgendaService) ConfirmAppointment(ctx context.Context, id string) error {
       appt, err := s.repo.FindByID(ctx, id)
       if err != nil { return err }
       if err := appt.Confirm(); err != nil { return err }
       return s.repo.Update(ctx, appt)
   }
   ```

### Parte B — Handlers: Confirm e NoShow

Em `internal/web/handlers/agenda_handler.go`, seguindo exatamente o padrão de `Cancel`:

**Confirm** — `POST /agenda/appointments/{id}/confirm`:
```
1. Extrair id da URL
2. Chamar h.agendaService.ConfirmAppointment(r.Context(), id)
3. Se erro: http.Error(w, ..., 500)
4. Buscar agendamento atualizado: h.agendaService.GetAppointment(ctx, id)
5. Mapear para AppointmentDetailModel
6. Se HX-Request: renderizar fragmento agenda.AppointmentDetail(vm)
7. Senão: redirect GET /agenda/appointments/{id}
```

**NoShow** — `POST /agenda/appointments/{id}/no-show`:
```
1. Extrair id da URL
2. Chamar h.agendaService.MarkNoShow(r.Context(), id)
3. Buscar agendamento atualizado
4. Mapear para AppointmentDetailModel
5. Se HX-Request: renderizar fragmento agenda.AppointmentDetail(vm)
6. Senão: redirect
```

### Parte C — Rotas (main.go)

No bloco switch/if que despacha rotas de `/agenda/appointments/`:

```go
case strings.HasSuffix(path, "/confirm") && r.Method == http.MethodPost:
    agendaHandler.Confirm(w, r)
case strings.HasSuffix(path, "/no-show") && r.Method == http.MethodPost:
    agendaHandler.NoShow(w, r)
```

### Parte D — UI: botões de ação por status

Em `web/components/agenda/appointment_detail.templ`, adicionar botões dentro do
bloco de ações condicionais. **Matriz de ações por status:**

| Status | Botões visíveis |
|--------|----------------|
| `scheduled` | **Confirmar** · Cancelar · Reagendar |
| `confirmed` | **Faltou** · Cancelar · Reagendar · Concluir |
| `completed` | (read-only — mostra link para sessão se existir) |
| `cancelled` | (read-only — badge cinza) |
| `no_show` | (read-only — badge amarelo/laranja) |

**Botão Confirmar** (apenas quando `model.Status == "scheduled"`):
```templ
<button
    hx-post={ string(templ.URL("/agenda/appointments/" + model.ID + "/confirm")) }
    hx-target="#modal-container"
    hx-swap="innerHTML"
    hx-confirm="Confirmar presença do paciente?"
    class="..."
>
    Confirmar
</button>
```

**Botão Faltou** (apenas quando `model.Status == "confirmed"`):
```templ
<button
    hx-post={ string(templ.URL("/agenda/appointments/" + model.ID + "/no-show")) }
    hx-target="#modal-container"
    hx-swap="innerHTML"
    hx-confirm="Registrar falta do paciente?"
    class="..."
>
    Faltou
</button>
```

> Use classes Tailwind retornadas por helpers em `types.go`. Nunca concatenação
> dinâmica de strings em classes inline. Se precisar de nova helper, adicione em
> `types.go` como função Go retornando string completa.

> Use `style=` para cores críticas que Tailwind v4 não detecta estaticamente.

### Parte E — Testes

**1. Teste de domínio** (`internal/domain/appointment/appointment_test.go` — se não existir, criar):
```
TestAppointment_Confirm_HappyPath    — status scheduled → confirmed, sem erro
TestAppointment_Confirm_AlreadyCancelled — retorna erro ao confirmar cancelado
TestAppointment_MarkNoShow_HappyPath — status scheduled → no_show, sem erro
```

**2. Teste de service** (arquivo `_test.go` ao lado de `agenda_service.go`):
- Use o padrão de mock existente no projeto (struct com campos `*Func` ou interface mock)
- Verificar pattern dos testes existentes antes de criar o mock
```
TestAgendaService_ConfirmAppointment_Success
TestAgendaService_ConfirmAppointment_NotFound
TestAgendaService_MarkNoShow_Success
```

**3. Teste de handler** (arquivo `agenda_handler_test.go` ou padrão existente):
```
TestAgendaHandler_Confirm_ReturnsFragmentOnHXRequest
TestAgendaHandler_Confirm_ReturnsRedirectOnNormalRequest
TestAgendaHandler_NoShow_ReturnsFragmentOnHXRequest
```

**4. Render test do componente** (`web/components/agenda/render_test.go`):
```
TestAppointmentDetail_ScheduledShowsConfirmButton
TestAppointmentDetail_ScheduledShowsNoFaltouButton  — Faltou não aparece em scheduled
TestAppointmentDetail_ConfirmedShowsFaltouButton
TestAppointmentDetail_ConfirmedShowsNoCancelButton  — validar matrix completa
TestAppointmentDetail_CompletedIsReadOnly
TestAppointmentDetail_NoShowIsReadOnly
```

---

## Arquivos a criar/modificar

| Ação | Arquivo |
|------|---------|
| Modificar | `internal/application/services/agenda_service.go` — adicionar ConfirmAppointment |
| Modificar | `internal/web/handlers/agenda_handler.go` — adicionar Confirm e NoShow |
| Modificar | `cmd/arandu/main.go` — registrar rotas /confirm e /no-show |
| Modificar | `web/components/agenda/appointment_detail.templ` — botões Confirmar e Faltou |
| Modificar | `web/components/agenda/types.go` — helpers de status se necessário |
| Criar/modificar | `internal/domain/appointment/appointment_test.go` |
| Criar/modificar | `internal/application/services/agenda_service_test.go` |
| Criar/modificar | `web/components/agenda/render_test.go` |

**NÃO criar migration SQL** — nenhuma mudança de schema.

---

## Skills de referência obrigatória

- **arandu-architecture** — multi-tenancy (extração do tenant DB do context), convenção de rotas, padrão HTMX fragmento vs página completa, containers do shell (`#modal-container`)
- **ddd-go** — padrão de application service: busca → regra de domínio → persiste; mock de repository em testes
- **go-templ-htmx-ux** — `templ.URL()` para URLs, `hx-confirm` para ações destrutivas, `hx-target` + `hx-swap`, retorno de fragmento em request HTMX
- **tailwind-components** — classes de botão de ação (primário, destrutivo), estados de badge por status clínico
- **clinical-domain** — semântica de status: `confirmed` = paciente confirmou presença; `no_show` = faltou sem aviso; impacto no histórico clínico

---

## Regras críticas

- Rodar `~/go/bin/templ generate ./web/components/...` após **qualquer** edição em `.templ`
- Nunca passar domain struct para o template — sempre mapear para `AppointmentDetailModel`
- Verificar `r.Header.Get("HX-Request") == "true"` antes de decidir fragmento vs redirect
- URLs no template: sempre `templ.URL(...)`, nunca string literal
- Classes Tailwind: retornar de helpers em `types.go`, nunca concatenar inline
- Sem dados clínicos em logs — apenas IDs (appointment_id) e timestamps

---

## Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Testes**
- [ ] `go test ./internal/domain/appointment/...` passa
- [ ] `go test ./internal/application/services/...` passa
- [ ] `go test ./web/components/agenda/...` passa (render tests)
- [ ] `go test ./...` sem regressões

**Comportamento (validar manualmente)**
- [ ] CA01: Agendamento `scheduled` → botão "Confirmar" visível; ao clicar, status muda para `confirmed` e modal atualiza inline sem reload
- [ ] CA02: Agendamento `confirmed` → botão "Faltou" visível; ao clicar, status muda para `no_show` e modal atualiza inline
- [ ] CA03: Agendamento `scheduled` → botão "Faltou" **não** aparece
- [ ] CA04: Agendamento `completed` ou `cancelled` → nenhum botão de ação exibido (read-only)
- [ ] CA05: `hx-confirm` exibe diálogo de confirmação antes de executar a ação

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa sem erros
- [ ] `./scripts/arandu_validate_handlers.sh` passa

---

## NÃO faça

- Não criar migration SQL — nenhuma mudança de schema nesta task
- Não reescrever handlers existentes (Cancel, Complete) — apenas adicionar os novos
- Não usar `templ.KV()` para classes condicionais — usar helpers em `types.go`
- Não logar o nome do paciente — apenas appointment_id
- Não retornar página completa em request HTMX — apenas fragmento

---

## Padrão de referência

- Handler: siga `Cancel` em `internal/web/handlers/agenda_handler.go` (extração de ID, call de service, verificação HX-Request, render de fragmento)
- Service: siga `CancelAppointment` em `internal/application/services/agenda_service.go`
- Render test: siga padrão dos testes existentes em `web/components/`
- Mock de service: verifique o padrão usado nos testes existentes antes de criar — struct com campos `*Func` ou interface mock
