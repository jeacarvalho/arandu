# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Generate templ components (required before every build after editing .templ files)
~/go/bin/templ generate ./web/components/...

# Build the server binary — ALWAYS use ./cmd/arandu/ not . (root has a test_version.go with its own main())
go build -o arandu ./cmd/arandu/

# Or use make (runs templ generate + go build)
make build

# Run (requires .env file in project root)
./arandu

# Dev mode (templ watch + go run in parallel)
make dev

# Tests
go test ./...
make test
```

## E2E / Manual Testing

Auth is required for all routes. Use the test user:
- Email: `arandu_e2e@test.com` / Password: `test123456`
- Tenant UUID: `9b33e4a0-ee4d-4054-8709-6b731f581a2a`
- Tenant DB: `storage/tenants/9b/33/clinical_9b33e4a0-ee4d-4054-8709-6b731f581a2a.db`

Use `scripts/arandu_e2e_all.sh` for automated E2E; `scripts/e2e/config.sh` handles login and cookie generation.

## Architecture

**Multi-tenant SQLite** — one SQLite file per therapist (clinical_{uuid}.db). Control plane (`storage/arandu_central.db`) manages users and tenants. Each HTTP request carries tenant context; repositories extract the correct DB from it.

**Layers:**
- `internal/domain/` — pure domain entities, no framework deps
- `internal/application/services/` — use-case orchestration
- `internal/infrastructure/repository/sqlite/` — SQLite implementations; schema lives in `migrations/`, never inline Go code
- `internal/web/handlers/` — HTTP handlers: extract params → call service → map to ViewModel → render via templ
- `web/components/` — `.templ` files organized by feature; helper Go functions for CSS class strings live in `types.go` alongside the ViewModels

**Web rendering:** templ + HTMX 2.x + **DaisyUI v4 + Tailwind CSS**. HTMX is served locally at `/static/js/htmx.min.js`. All pages use `layout.Shell(config, content)`.

Key containers in shell:
- `#main-content` — primary HTMX swap target
- `#shell-sidebar` — sidebar OOB swap target (`hx-swap-oob="innerHTML"`)
- `#shell-breadcrumb` — breadcrumb OOB swap target (`hx-swap-oob="true"`)
- `#modal-container` — appointment detail modals
- `#drawer-container` — slide-in forms

## Critical Rules

**templ:** Run `~/go/bin/templ generate` after every `.templ` edit. HTML templates in `.html` files are forbidden.

**Tailwind v4 dynamic classes:** Tailwind scans source as plain text — it cannot detect classes built by string concatenation at runtime. All CSS class strings must be returned as complete strings from helper functions in `types.go` (never use `templ.KV()` for conditional classes). For text that must be visible regardless of Tailwind compilation, use inline `style=` attributes.

**HTMX swaps:** `#agenda-content` uses `hx-swap="outerHTML"`. Fragments returned for HTMX requests must be the component only (no shell wrapper). Check `r.Header.Get("HX-Request")` in handlers.

**Migrations:** Never create tables in Go code. Add a new SQL file to `internal/infrastructure/repository/sqlite/migrations/` — the migrator applies them on startup to all tenant DBs.

**ViewModels:** Domain entities must never reach templ templates. Handlers map domain objects to ViewModels (structs in `web/components/<feature>/types.go`) before rendering.

## Design System — DaisyUI + Tema Arandu

**Stack CSS:** DaisyUI v4 (componentes semânticos) + Tailwind utilities. Não criar CSS custom — usar classes DaisyUI.

**Regras de ouro:**
- Botões: `btn btn-primary`, `btn btn-ghost`, `btn btn-error btn-outline`
- Cards: `card bg-base-100 border border-base-300 shadow-sm`
- Badges de status: `badge badge-success`, `badge badge-warning`, `badge badge-ghost`
- Tipografia clínica: classe `.clinical` (Source Serif 4) para prontuários, observações, títulos de página
- UI text: Inter (padrão do body)
- Não editar `input-v2.css` nem `tailwind-v2.css` — são arquivos legacy

**Referência completa:** skill `daisyui-arandu`
**Protótipo visual:** `design_handoff_arandu_redesign/daisyui_shell_dashboard.html`
