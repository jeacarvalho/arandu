# Task: Agenda — Concluir agendamento com criação de sessão vinculada
Requirement: REQ-07-02-01
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto

Ao clicar "Concluir" no modal de detalhe de um agendamento, o sistema apenas muda o
status para `completed` — mas não cria nem abre uma sessão clínica para documentação.
O campo `SessionID *string` na entidade `Appointment` e a lógica `LinkSession` já existem,
mas nunca são preenchidos na prática.

Esta task fecha esse ciclo: clicar "Concluir" agora mostra um painel de confirmação inline
com duas opções — criar sessão e ir documentar, ou concluir sem sessão. O primeiro caminho
cria a `Session`, vincula ao `Appointment`, e redireciona para `/session/{id}`.

### Estado atual (leia antes de qualquer mudança)

**`appointment_detail.templ`** — o botão "Concluir" faz POST direto:
```templ
<button
    hx-post={ "/agenda/appointments/" + model.ID + "/complete" }
    hx-target="#agenda-content"
    hx-swap="outerHTML"
    onclick="document.getElementById('modal-container').innerHTML=''"
    class="flex-1 h-10 ...">
    Concluir
</button>
```

**`AgendaHandler`** — struct atual:
```go
type AgendaHandler struct {
    agendaService  AgendaServiceInterface
    patientService PatientServiceInterface
}
func NewAgendaHandler(agendaService AgendaServiceInterface, patientService PatientServiceInterface) *AgendaHandler
```

**`AgendaService.CompleteAppointment`**:
```go
func (s *AgendaService) CompleteAppointment(ctx context.Context, id string, sessionID string) error
// Já chama appt.LinkSession(sessionID) se sessionID != ""
```

**`SessionService.CreateSession`**:
```go
func (s *SessionService) CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error)
```

**`main.go`** — injeções relevantes:
```go
agendaService := services.NewAgendaService(appointmentRepo)
agendaHandler := handlers.NewAgendaHandler(agendaService, patientServiceAdapter)
// sessionServiceAdapter já existe e é injetado no sessionHandler
```

---

## O que implementar

### 1. Adicionar `sessionService` ao `AgendaHandler`

**Arquivo:** `internal/web/handlers/agenda_handler.go`

```go
type AgendaHandler struct {
    agendaService  AgendaServiceInterface
    patientService PatientServiceInterface
    sessionService SessionServiceInterface  // novo campo
}

func NewAgendaHandler(
    agendaService AgendaServiceInterface,
    patientService PatientServiceInterface,
    sessionService SessionServiceInterface,  // novo param
) *AgendaHandler {
    return &AgendaHandler{agendaService, patientService, sessionService}
}
```

`SessionServiceInterface` — declare no arquivo `agenda_handler.go` apenas com o método
necessário:
```go
type SessionServiceInterface interface {
    CreateSession(ctx context.Context, patientID string, date time.Time, summary string) (*session.Session, error)
}
```

> **Import necessário:** `"arandu/internal/domain/session"`

### 2. Novo handler `CompleteWithSession`

**Arquivo:** `internal/web/handlers/agenda_handler.go` — adicionar ao final.

```go
// CompleteWithSession handles POST /agenda/appointments/{id}/complete-with-session
// Cria uma Session vinculada ao Appointment e redireciona para ela.
func (h *AgendaHandler) CompleteWithSession(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id := extractIDFromPath(r.URL.Path, "/agenda/appointments/")
    id = strings.TrimSuffix(id, "/complete-with-session")

    ctx := r.Context()

    appt, err := h.agendaService.GetAppointment(ctx, id)
    if err != nil || appt == nil {
        http.Error(w, "Appointment not found", http.StatusNotFound)
        return
    }

    sess, err := h.sessionService.CreateSession(ctx, appt.PatientID, appt.Date, "")
    if err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return
    }

    if err := h.agendaService.CompleteAppointment(ctx, id, sess.ID); err != nil {
        http.Error(w, "Failed to complete appointment", http.StatusInternalServerError)
        return
    }

    redirectURL := "/session/" + sess.ID
    if r.Header.Get("HX-Request") == "true" {
        w.Header().Set("HX-Redirect", redirectURL)
        w.WriteHeader(http.StatusOK)
        return
    }
    http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
```

### 3. Registrar nova rota em `main.go`

**Arquivo:** `cmd/arandu/main.go`

**a)** Injetar `sessionServiceAdapter` no `AgendaHandler`:
```go
// Antes:
agendaHandler := handlers.NewAgendaHandler(agendaService, patientServiceAdapter)

// Depois:
agendaHandler := handlers.NewAgendaHandler(agendaService, patientServiceAdapter, sessionServiceAdapter)
```

**b)** No bloco `mux.HandleFunc("/agenda/appointments/", ...)`, adicionar o roteamento
para o novo sufixo `/complete-with-session` — no mesmo switch/if que trata `/cancel`
e `/complete`. Siga o padrão existente exatamente:

```go
case strings.HasSuffix(path, "/complete-with-session"):
    agendaHandler.CompleteWithSession(w, r)
```

### 4. Painel de confirmação inline no modal

**Arquivo:** `web/components/agenda/appointment_detail.templ`

Substituir o botão "Concluir" (bloco `else` dentro do `if model.SessionID != nil`) pelo
seguinte padrão: botão que revela um painel de confirmação via JavaScript puro (sem
round-trip HTTP). O painel fica oculto por default e substituí a área de ações.

```templ
} else {
    <!-- Botão inicial: mostra painel de confirmação -->
    <button
        id={ "btn-concluir-" + model.ID }
        onclick={ templ.SafeScript("document.getElementById('confirm-panel-" + model.ID + "').classList.remove('hidden'); this.classList.add('hidden')") }
        class="flex-1 h-10 text-sm font-medium text-white bg-arandu-primary hover:bg-arandu-dark rounded-xl transition-colors"
    >
        Concluir
    </button>

    <!-- Painel de confirmação (oculto inicialmente) -->
    <div id={ "confirm-panel-" + model.ID } class="hidden w-full">
        <p class="text-xs text-center mb-3" style="color:#6b7280">
            Deseja registrar o atendimento no prontuário?
        </p>
        <div class="flex gap-2">
            <button
                onclick={ templ.SafeScript("document.getElementById('modal-container').innerHTML=''") }
                hx-post={ "/agenda/appointments/" + model.ID + "/complete-with-session" }
                hx-target="body"
                hx-swap="none"
                class="flex-1 h-10 text-sm font-medium text-white bg-arandu-primary hover:bg-arandu-dark rounded-xl transition-colors"
            >
                <i class="fas fa-file-medical text-xs mr-1"></i>Sim, documentar
            </button>
            <button
                hx-post={ "/agenda/appointments/" + model.ID + "/complete" }
                hx-target="#agenda-content"
                hx-swap="outerHTML"
                onclick={ templ.SafeScript("document.getElementById('modal-container').innerHTML=''") }
                class="flex-1 h-10 text-sm font-medium rounded-xl border transition-colors"
                style="color:#374151;border-color:#d1d5db"
            >
                Só concluir
            </button>
        </div>
    </div>
}
```

> **Atenção — `templ.SafeScript`:** em templates Go templ, use `templ.SafeScript(...)` para
> atributos `onclick` que contêm interpolação de variáveis. Nunca concatene strings Go
> diretamente em atributos de evento. Verifique a versão do templ no projeto — se
> `templ.SafeScript` não existir, use `templ.ComponentScript` ou atributo como string
> com `templ.SafeURL`.
>
> **Alternativa segura** — se `templ.SafeScript` não compilar:
> ```templ
> <div id={ "confirm-panel-" + model.ID } ...>
> ```
> E para o onclick do botão inicial, use um atributo `data-panel-id={ model.ID }` e um
> script global que faça o toggle pelo data attribute — evitando interpolação inline.

---

## Checklist de implementação

- [ ] `AgendaHandler` tem campo `sessionService SessionServiceInterface`
- [ ] `NewAgendaHandler` aceita 3 parâmetros (compila sem erros)
- [ ] `CompleteWithSession` handler implementado e retorna HX-Redirect para HTMX
- [ ] Rota `/complete-with-session` registrada em `main.go`
- [ ] `sessionServiceAdapter` injetado no `agendaHandler` em `main.go`
- [ ] Botão "Concluir" no modal abre painel de confirmação (sem POST imediato)
- [ ] "Sim, documentar" → chama `/complete-with-session` → redireciona para `/session/{id}`
- [ ] "Só concluir" → chama `/complete` → comportamento existente
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build ./cmd/arandu/` sem erros
- [ ] `go test ./internal/web/handlers/...` continua passando (testes existentes não quebram)

---

## Arquivos a modificar (não criar)

- `internal/web/handlers/agenda_handler.go` — novo campo, novo handler
- `web/components/agenda/appointment_detail.templ` — painel de confirmação
- `cmd/arandu/main.go` — injeção + nova rota

---

## 🔒 Privacidade

- [x] **Tier 1 — PII indireto**: `appt.PatientID` e `appt.PatientName` transitam no handler
  mas NÃO devem aparecer em logs. Não adicione nenhum `log.Printf` com dados do paciente.
- [x] A `Session` criada começa com `summary = ""` — nenhum dado clínico é gerado
  automaticamente. O psicólogo preencherá manualmente ao documentar.

---

## 🚫 NÃO faça

- Não mova lógica de criação de sessão para dentro de `AgendaService` — o handler
  orquestra os dois services; cada service permanece com sua responsabilidade
- Não altere `AgendaService.CompleteAppointment` — ela já aceita `sessionID string`
- Não crie migration — nenhuma mudança de schema
- Não modifique o handler `Complete` existente — ele continua funcionando para "só concluir"
- Não use `html/template` — apenas `.templ`
- Não quebre os testes existentes em `agenda_handler_test.go`

---

## 📎 Padrões de referência

- Estrutura do novo handler: siga `Cancel` e `Complete` em `agenda_handler.go` —
  mesmo padrão de extração de ID, verificação de método, HX-Redirect
- Injeção no handler: siga como `sessionHandler` recebe `sessionServiceAdapter` em `main.go`
- Toggle JS no templ: siga o padrão `onclick="document.getElementById(...)"` já usado
  nos botões de fechar do modal (`appointment_detail.templ` linha 27 e 78)
