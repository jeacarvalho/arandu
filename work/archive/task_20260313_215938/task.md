# TASK-01

## Criar entidade Patient no domínio

Prompt para enviar ao agente:

---

# TASK — Domain: Patient Entity

You are working on the **Arandu** system.

Arandu is a clinical intelligence platform for mental health professionals.

This task is **small and focused**.

Your goal is to create the **domain entity for Patient**.

Do not implement UI, handlers, or HTTP routes yet.

---

# Requirement

```text
REQ-01-00-01 Criar paciente
```

This is the first functional requirement of the system.

Before creating the UI for patient registration, we must establish the **domain entity**.

---

# Architecture context

The system follows these principles:

```
Go
DDD
Clean Architecture
SOLID
```

Layer structure:

```
domain
application
infrastructure
web
```

Domain must not depend on other layers.

---

# Step 1 — Create domain package

Create directory:

```
internal/domain/patient
```

---

# Step 2 — Create entity file

Create file:

```
internal/domain/patient/entity.go
```

---

# Step 3 — Implement Patient entity

Define the Patient struct.

Example:

```go
type Patient struct {
	ID        string
	Name      string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

Rules:

* ID must be string
* Name must be required
* timestamps must exist
* no persistence logic inside the entity

---

# Step 4 — Create constructor

Add a constructor function.

Example:

```go
func NewPatient(name string, notes string) *Patient
```

Responsibilities:

* generate ID
* set CreatedAt
* set UpdatedAt

ID generation can use:

```
github.com/google/uuid
```

---

# Step 5 — Validation rule

Ensure patient name is not empty.

Example:

```
name cannot be empty
```

Return an error if validation fails.

---

# Step 6 — Keep domain pure

The domain entity must not import:

```
database/sql
net/http
framework libraries
```

Only standard library packages allowed.

---

# Step 7 — Expected directory structure

After this task the project must contain:

```
internal/domain/patient/entity.go
```

---

# Step 8 — Output

Provide:

1. created directory structure
2. full content of entity.go
3. short explanation of design decisions

---

# Important constraint

Do not implement:

```
repositories
database migrations
HTTP handlers
templates
HTMX
```

Those will be implemented in later tasks.

This task is strictly about the **domain entity**.

---

# Goal

At the end of this task the Arandu system must contain a **clean domain entity for Patient** that future tasks can build upon.

---
