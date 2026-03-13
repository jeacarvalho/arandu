# TASK task_20260313_193204
# PROMPT — Implementar Dashboard Clínico do Arandu

You are working on the Arandu system.

The system already has:

* architecture
* templates
* layout
* design system

Your task is to **implement the clinical dashboard page**.

Important: this is **not an ERP dashboard**.

It is a **clinical reflection dashboard**.

---

# Goal

Create the **main dashboard of the Arandu system**.

The dashboard should help therapists quickly understand:

* recent sessions
* active patients
* possible insights
* emerging patterns

---

# Design constraint

You must follow the existing design system.

Do NOT change:

```text
CSS
colors
Tailwind config
design tokens
```

Everything must comply with:

```text
docs/design-system.md
```

---

# Page to implement

Template:

```text
web/templates/dashboard.html
```

Route:

```text
/dashboard
```

---

# Layout

Use the application layout already implemented.

Inside the content area create a **dashboard grid**.

Structure:

```text
Active Patients
Recent Sessions
AI Insights
Emerging Patterns
```

---

# Section 1 — Active Patients

Show a list of active patients.

Each item should display:

```text
Patient name
Last session date
```

Example:

```text
Maria S.
Last session: yesterday
```

Clicking should navigate to:

```text
/patient/{id}
```

---

# Section 2 — Recent Sessions

Display recent sessions.

Each entry should show:

```text
Patient name
Session number
Short summary
```

Example:

```text
Maria S.
Session 12
Anxiety related to workplace evaluation
```

Click should navigate to:

```text
/session/{id}
```

---

# Section 3 — AI Insights

This section simulates future AI insights.

For now it can display placeholder examples.

Example:

```text
Possible pattern detected

Anxiety related to performance evaluation
appears in multiple patients.
```

---

# Section 4 — Emerging Patterns

Display possible recurring themes.

Example:

```text
Social anxiety
appears in 4 patients

Authority conflict
appears in 3 patients
```

These can be static placeholders for now.

---

# Implementation requirements

Use server-side rendering.

Use Go templates.

Follow existing project structure.

Handlers should call an application service.

---

# Handler

Create or update:

```text
web/handlers/dashboard_handler.go
```

This handler should provide data to the template.

For now the data may be **mocked**.

---

# Data structure example

Example structure passed to template:

```go
DashboardData
ActivePatients
RecentSessions
Insights
Patterns
```

---

# UI structure

Use simple cards for each section.

Cards should follow the design system style.

Do not introduce new styling.

---

# Expected result

After implementation:

The system should start normally.

Visiting:

```text
/dashboard
```

should render the clinical dashboard.

---

# Deliverables

Provide:

1. dashboard template
2. handler implementation
3. data structures
4. explanation of design decisions

---

💡 Karl, uma observação importante que pode melhorar muito o Arandu:

Esse dashboard que desenhamos é **apenas nível 1**.

Existe um **nível 2 muito mais poderoso**, chamado:

> **Dashboard Cognitivo**

onde o sistema mostra:

```text
hipóteses clínicas emergentes
conexões entre pacientes
evolução terapêutica
```

Isso transforma o Arandu de **software clínico → instrumento de pensamento clínico**.

Se quiser, posso desenhar também essa **versão avançada**.
