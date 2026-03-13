# PROMPT — Tarefa 1

## Arandu — Foundation of the System

You are starting development of a system called **Arandu**.

Arandu is a **clinical intelligence platform for mental health professionals** such as psychologists and therapists.

The system helps professionals:

* organize clinical sessions
* record observations
* explore hypotheses
* register interventions
* discover patterns over time
* reflect on cases with assistance from AI

The system is designed to become a **long-term cognitive workspace for therapists**.

This task is **not about implementing features**.

Your goal is to **create the architectural foundation of the system**.

---

# Step 1 — Understand the development model

The Arandu project follows a **spec-driven development model**.

All implementation must follow this structure:

```text
DVP
 ↓
VISION
 ↓
CAPABILITY
 ↓
REQUIREMENT
 ↓
TASK
 ↓
CODE
```

Code must always be traceable back to a **requirement**.

Do not invent features outside the documented specifications.

---

# Step 2 — Read the project documentation

Before writing any code, read the following directories and documents:

```
docs/dvp.md
docs/design-system.md
docs/vision/
docs/capabilities/
docs/requirements/
docs/learnings/
```

These documents describe:

* system goals
* architecture principles
* UI philosophy
* project constraints

The Arandu system is **not a generic CRUD application**.

It is a **knowledge system for clinical reasoning**.

---

# Step 3 — Understand the technology stack

Arandu uses a deliberately simple and robust stack.

Backend:

```
Go
HTMX
SQLite
```

Architecture principles:

```
DDD
Clean Architecture
SOLID
TDD
```

Important persistence rule:

Each professional using Arandu has **one SQLite database file**.

This design improves:

* privacy
* portability
* backup
* simplicity

---

# Step 4 — Create the project structure

Create the initial project structure using Go.

Suggested layout:

```
cmd/arandu/main.go

internal/

domain/
patient/
session/
observation/
intervention/
insight/

application/

services/

infrastructure/

repository/
sqlite/

web/

handlers/
templates/

static/

css/
js/
```

The architecture must follow these layers:

```
domain
application
infrastructure
web
```

The **domain layer must not depend on any other layer**.

---

# Step 5 — Domain foundation

Create the core domain entities.

These represent the **clinical thinking model** of the system.

Create domain packages for:

```
Patient
Session
Observation
Intervention
Insight
```

Each entity should include basic fields:

```
ID
CreatedAt
UpdatedAt
```

You may also include minimal domain attributes appropriate for each concept.

Example ideas:

Patient

```
ID
Name
Notes
CreatedAt
UpdatedAt
```

Session

```
ID
PatientID
Date
Notes
CreatedAt
UpdatedAt
```

Observation

```
ID
SessionID
Content
CreatedAt
```

Do not implement complex logic yet.

Focus on **clear domain modeling**.

---

# Step 6 — Domain rules

Each domain entity must live in its own package.

Example:

```
internal/domain/patient
internal/domain/session
```

Each package should contain:

```
entity definition
repository interface
domain types
```

Example repository interface:

```
PatientRepository
SessionRepository
```

Domain packages must **not depend on infrastructure or web code**.

---

# Step 7 — Persistence layer

Implement repository interfaces in the domain layer.

Then implement SQLite repositories inside:

```
internal/infrastructure/repository/sqlite
```

Important rule:

Do NOT use any ORM.

Do NOT use libraries like:

```
gorm
ent
```

Instead use:

```
database/sql
```

Repositories should use **simple explicit SQL queries**.

Keep repositories small and clear.

---

# Step 8 — Application layer

Create application services that orchestrate domain operations.

Example location:

```
internal/application/services
```

Example services:

```
PatientService
SessionService
```

These services will coordinate:

* domain entities
* repositories
* business flows

Application services should contain **use-case orchestration**, not persistence.

---

# Step 9 — Web server foundation

Create the base HTTP server.

Location:

```
cmd/arandu/main.go
```

Use the Go standard library HTTP server.

Create initial routes:

```
/dashboard
/patients
/patient/{id}
/session/{id}
```

---

# Step 10 — Web architecture rules

Handlers should remain **thin**.

Handlers must:

```
receive HTTP request
call application service
return response
```

Handlers must **not contain business logic**.

Application services must orchestrate domain logic.

---

# Step 11 — HTMX compatible handlers

The UI will use **HTMX**.

Design handlers so that they can return **HTML fragments** when necessary.

Example use case:

```
HTMX request loads patient session list
```

Handlers should be able to respond with:

```
full HTML page
or partial fragment
```

---

# Step 12 — Template system

Create template structure:

```
web/templates/

layout.html
dashboard.html
patients.html
patient.html
session.html
```

Templates should extend a base layout.

Keep templates minimal.

No styling is required yet beyond basic layout.

---

# Step 13 — Static assets

Prepare static asset directories:

```
web/static/css
web/static/js
```

Create a base CSS file.

Do not implement complex styling.

The design system is defined in:

```
docs/design-system.md
```

---

# Step 14 — Minimal UI skeleton

Create a minimal dashboard page.

Layout concept:

```
Sidebar: patients
Main area: session content
Right panel: insights
```

This should match the design philosophy described in:

```
docs/design-system.md
```

---

# Step 15 — System must run

At the end of this task:

The system should:

```
compile
start HTTP server
render dashboard page
```

Even if functionality is minimal.

---

# Step 16 — Documentation

After creating the system foundation, generate a document:

```
docs/architecture/system_structure.md
```

Explain:

```
project directory structure
domain model
repository design
application layer
web layer
```

---

# Step 17 — Explicit constraints

Do NOT implement the following yet:

```
AI features
pattern detection
authentication
complex UI
analytics
```

This task is strictly about **foundations**.

---

# Final objective

At the end of this task the Arandu repository must contain:

```
clean Go project
DDD domain structure
repository interfaces
SQLite repository implementation
HTMX compatible web server
HTML template skeleton
basic dashboard page
```

The codebase must be **clean, minimal, and extensible**.

---

# Important principle

Arandu is not just software.

It is a **cognitive workspace for therapists**.

The architecture must support:

```
long-term clinical knowledge
pattern discovery
reflective workflows
```

Focus on **clarity, simplicity, and extensibility**.

