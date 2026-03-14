

# PROMPT — Implementar REQ-01-01-01 Criar sessão

You are implementing a feature in the **Arandu** system.

Arandu is a clinical intelligence platform for mental health professionals.

The system already supports:

```text
patient registration
patient persistence
patient listing
```

Now we will implement the capability:

```text
CAP-01-01 Registro de sessões
```
Leia o arquivo docs/capabilities/cap-01-01-registro-sessoes.md

Requirement:

```text
REQ-01-01-01 Criar sessão
```
Leia o arquivo docs/requirements/req-01-01-01-criar-sessao.md

This feature allows therapists to **register a therapy session for an existing patient**.

---

# Architecture

The system uses:

```text
Go
Clean Architecture
DDD
SQLite
HTMX
Server-side HTML rendering
```

Layer structure:

```text
domain
application
infrastructure
web
```

Dependencies must follow:

```text
web → application → domain
infrastructure → domain
```

Domain must remain independent.

---

# Important HTMX rule

This system uses **HTML over the wire**.

Handlers must return:

```text
HTML pages
HTML fragments
```

Handlers must **not return JSON APIs**.

HTMX interactions must be handled by **server-rendered templates**.

---

# TASK-01 — Create Session domain entity

Create directory:

```text
internal/domain/session
```

Create file:

```text
internal/domain/session/entity.go
```

Define Session entity.

Example:

```go
type Session struct {
	ID        string
	PatientID string
	Date      time.Time
	Summary   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

Add constructor:

```go
func NewSession(patientID string, date time.Time, summary string) *Session
```

Responsibilities:

```text
generate session ID
set CreatedAt
set UpdatedAt
```

Use:

```text
github.com/google/uuid
```

---

# TASK-02 — Create Session repository interface

Create file:

```text
internal/domain/session/repository.go
```

Define repository:

```go
type Repository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	ListByPatient(ctx context.Context, patientID string) ([]*Session, error)
}
```

---

# TASK-03 — Create SQLite migration

Create file:

```text
migrations/002_create_sessions_table.sql
```

SQL:

```sql
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATETIME NOT NULL,
    summary TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id)
);

CREATE INDEX IF NOT EXISTS idx_sessions_patient
ON sessions(patient_id);
```

---

# TASK-04 — Implement SQLite repository

Create file:

```text
internal/infrastructure/repository/sqlite/session_repository.go
```

Define struct:

```go
type SessionRepository struct {
	db *sql.DB
}
```

Constructor:

```go
func NewSessionRepository(db *sql.DB) *SessionRepository
```

Implement:

```text
Create
GetByID
ListByPatient
```

Example INSERT:

```sql
INSERT INTO sessions (
    id,
    patient_id,
    date,
    summary,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?)
```

Use:

```text
ExecContext
QueryRowContext
QueryContext
```

---

# TASK-05 — Create application service

Create file:

```text
internal/application/services/create_session_service.go
```

Define service:

```go
type CreateSessionService struct {
	repo session.Repository
}
```

Constructor:

```go
func NewCreateSessionService(repo session.Repository) *CreateSessionService
```

Input:

```go
type CreateSessionInput struct {
	PatientID string
	Date      time.Time
	Summary   string
}
```

Execute method:

```go
func (s *CreateSessionService) Execute(ctx context.Context, input CreateSessionInput) (*session.Session, error)
```

Responsibilities:

```text
validate input
create session entity
persist session
return session
```

---

# TASK-06 — Create HTTP handler (HTMX compatible)

Create file:

```text
internal/web/handlers/session_handler.go
```

Handler responsibilities:

```text
render session creation page
process form submission
call CreateSessionService
redirect to patient page
```

Routes:

```text
GET  /patients/{id}/sessions/new
POST /sessions
```

Important rule:

Handlers must return **HTML templates**, not JSON.

---

# TASK-07 — Create HTML template

Create template:

```text
web/templates/session_new.html
```

Example:

```html
<h1>Nova sessão</h1>

<form 
    method="POST"
    action="/sessions"
    hx-post="/sessions"
    hx-swap="none"
>

<input type="hidden" name="patient_id" value="{{ .PatientID }}">

<label>Data</label>
<input type="date" name="date">

<label>Resumo</label>
<textarea name="summary"></textarea>

<button type="submit">Salvar sessão</button>

</form>
```

---

# Important HTMX behaviour

Form submission must:

```text
send POST request
server processes session creation
server redirects to patient page
```

Do not return JSON.

---

# Expected result

After implementation the system must allow:

```text
open patient
↓
click "new session"
↓
fill session form
↓
save session
↓
session appears in patient history
```

---

# Deliverables

Provide:

```text
full directory tree
Session entity code
repository code
migration SQL
application service
handler
template
```

---

# Important constraints

Do not implement:

```text
observations
interventions
AI features
analytics
```

Those belong to later capabilities.

---

# Goal

After this implementation the Arandu system must support:

```text
Patient
   └── Sessions
```

which establishes the **first level of clinical memory**.

---