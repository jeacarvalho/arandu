# RELATÓRIO FINAL — CORREÇÃO CIRÚRGICA

**Data:** 31/03/2026
**Tarefa:** Aplicação de correções cirúrgicas para inconsistência de CSS após navegação HTMX

---

## Resumo das Alterações

### FASE 1 — UNIFICAÇÃO DE CSS ✅

| Arquivo | Alteração |
|---------|-----------|
| `web/components/layout/layout.templ` | 4 ocorrências de `tailwind.css` alteradas para `tailwind-v2.css` com `GetCSSVersionV2()` |

**Resultado:** Ambos os layouts (`shell_layout.templ` e `layout.templ`) agora usam o mesmo arquivo CSS canônico (`tailwind-v2.css`).

---

### FASE 2 — CORREÇÃO DOS HANDLERS GO ✅

| Arquivo | Função | Alteração |
|---------|--------|-----------|
| `internal/web/handlers/dashboard_handler.go` | `Show` | Adicionado `HX-Request` check para retornar fragmento |
| `internal/web/handlers/observation_handler.go` | `GetObservation`, `GetObservationEditForm`, `UpdateObservation` | Content-Type: `text/html; charset=utf-8` |
| `internal/web/handlers/intervention_handler.go` | `GetIntervention`, `GetInterventionEditForm`, `UpdateIntervention` | Content-Type: `text/html; charset=utf-8` |
| `internal/web/handlers/session_handler.go` | `CreateObservation`, `CreateIntervention`, `CloseGoalWithNote` | Content-Type: `text/html; charset=utf-8` |
| `internal/web/handlers/ai_handler.go` | `GeneratePatientSynthesis`, `renderError` | Content-Type: `text/html; charset=utf-8` |

**Resultado:** Handlers agora diferenciam requisições HTMX de requisições de página completa, e todos retornam `Content-Type` consistente com charset.

---

### FASE 3 — CORREÇÃO DE IDs DUPLICADOS ✅

| Layout | IDs Renomeados |
|--------|-----------------|
| `shell_layout.templ` | `patient-search` → `shell-patient-search`, `search-results` → `shell-search-results`, `search-loading` → `shell-search-loading` |
| `layout.templ` | `patient-search` → `legacy-patient-search`, `search-results` → `legacy-search-results`, `search-loading` → `legacy-search-loading` |

**Resultado:** IDs únicos para cada contexto, evitando conflitos de target HTMX.

---

### FASE 4 — ELIMINAÇÃO DE SCRIPTS DUPLICADOS ✅

| Arquivo | Alteração |
|---------|-----------|
| `web/static/js/htmx-handlers.js` | **NOVO** — Handlers centralizados com guard `window.__htmxHandlersLoaded` |
| `web/components/layout/shell_layout.templ` | Substituído inline script por referência externa |
| `web/components/layout/layout.templ` | Adicionada referência externa em 4 pontos |

**Resultado:** Scripts de HTMX event handlers centralizados em arquivo único, evitando múltiplos event listeners.

---

### FASE 5 — CORREÇÃO DO ALPINE.JS ✅

| Arquivo | Alteração |
|---------|-----------|
| `web/static/js/htmx-handlers.js` | `document.querySelector('[x-data]')?.__x?.$data` → `document.querySelector('.shell')?.__x?.$data` com verificação de existência |

**Resultado:** Acesso ao estado Alpine.js agora feito através de seletor específico (`.shell`) com verificação de existência de `__x`.

---

### FASE 6 — ATUALIZAR @SOURCE DO TAILWIND ✅

| Arquivo | Alteração |
|---------|-----------|
| `web/static/css/input-v2.css` | Adicionado `@source "../../../internal/web/handlers/**/*.go"` |
| `web/static/css/tailwind-v2.css` | **REBUILD** — de 66KB para 78KB |

**Resultado:** Classes CSS geradas dinamicamente em handlers Go agora são detectadas pelo Tailwind.

---

## Arquivos Modificados

```
web/components/layout/layout.templ          [CSS unificado, IDs renomeados, script externo]
web/components/layout/layout_templ.go      [gerado automaticamente]
web/components/layout/shell_layout.templ    [CSS unificado, IDs renomeados, script externo]
web/components/layout/shell_layout_templ.go [gerado automaticamente]
web/static/css/input-v2.css                  [@source handlers]
web/static/css/tailwind-v2.css               [rebuild]
web/static/js/htmx-handlers.js               [NOVO]
internal/web/handlers/dashboard_handler.go   [HX-Request check]
internal/web/handlers/observation_handler.go [Content-Type charset]
internal/web/handlers/intervention_handler.go [Content-Type charset]
internal/web/handlers/session_handler.go      [Content-Type charset]
internal/web/handlers/ai_handler.go           [Content-Type charset]
```

---

## Itens Pendentes

**Nenhum.** Todas as 6 fases foram concluídas com sucesso.

---

## Teste de Validação Sugerido

### Cenário 1: Navegação Entre Páginas com Layouts Diferentes

1. Acessar Dashboard (`/dashboard`) — usa `shell_layout.templ` com CSS v2
2. Navegar para Patient Detail via HTMX (`hx-get` + `hx-push-url`)
3. Retornar ao Dashboard — CSS mantido, sem necessidade de Ctrl+Shift+R

### Cenário 2: Busca de Pacientes

1. No Dashboard, digitar na busca de pacientes
2. Verificar que resultados aparecem em `#shell-search-results`
3. Navegar para página que usa `layout.templ`
4. Digitar na busca — verificar que resultados aparecem em `#legacy-search-results`
5. Retornar ao Dashboard — busca continua funcionando

### Cenário 3: Handlers HTMX

1. Abrir DevTools → Network
2. Fazer requisição HTMX (ex: clicar em um paciente)
3. Verificar Headers: `HX-Request: true` returna fragmento sem `<html><head>...`
4. Acessar mesma URL diretamente no browser — retorna página completa

---

## Riscos Residuais (Não cobertos)

| Risco | Status | Observação |
|-------|--------|------------|
| `hx-push-url` sem `hx-head="merge"` | MANTIDO | 29 ocorrências — mitigado pela unificação de CSS |
| `hx-history` desabilitado | MANTIDO | Risco baixo — snapshots não preservados |
| Classes dinâmicas não safelistadas | PARCIAL | `@source` handlers adicionado ao Tailwind |
| Inline styles em `login.templ`/`signup.templ` | MANTIDO | Fora do escopo do prompt |
| `style="padding-top: 80px; padding-left: 260px;"` em shell_layout.templ:112 | MANTIDO | Fora do escopo do prompt |

---

## Comandos de Validação

```bash
# Build completo
go build ./cmd/arandu

# Testes de handlers
go test ./internal/web/handlers/... -v

# Testes de componentes
go test ./web/components/... -v

# Rebuild do CSS (se necessário)
npm run tailwind:build:v2

# Geração de templates
~/go/bin/templ generate
```

---

## Contato

Para dúvidas sobre as correções aplicadas, consulte o relatório forense em `CSS_HTMX_FORENSIC_REPORT.md`.