# 🕵️ RELATÓRIO FORENSE — Auditoria CSS/HTMX/Cache

**Data:** 31/03/2026
**Sintoma:** CSS quebrado após navegação HTMX (resolvido apenas com hard refresh)
**Escopo:** Análise completa de vetores de inconsistência de estilo

---

## 1. SUMÁRIO EXECUTIVO

O problema de "CSS quebrado após navegação HTMX" tem **múltiplas causas raízes inter-relacionadas**:

| Criticidade | Causa Raiz |
|-------------|------------|
| **[CRÍTICO]** | Duas versões de Tailwind CSS (`tailwind.css` vs `tailwind-v2.css`) servidas por layouts diferentes |
| **[CRÍTICO]** | IDs duplicados (`patient-search`, `search-results`) em múltiplos layouts causam conflitos de alvo HTMX |
| **[ALTO]** | Handlers inconsistentes na verificação `HX-Request` — alguns retornam página completa quando deveriam retornar fragmentos |
| **[ALTO]** | Scripts de configuração HTMX e tratamento de erros duplicados 4x nos layouts, criando listeners duplicados |

**Conclusão:** A navegação HTMX entre páginas usa layouts diferentes versões de CSS. Fragmentos HTML servidos para uma página podem conter classes CSS inexistentes no cache do browser da página anterior. O hard refresh força recarregamento completo, resolvendo temporariamente.

---

## 2. ACHADOS POR VETOR

### VETOR 1 — Inclusão do CSS (stylesheet e Tailwind)

| Risco | Arquivo | Linha | Problema |
|-------|---------|-------|----------|
| **[CRÍTICO]** | `web/components/layout/shell_layout.templ` | 46-51 | Usa `tailwind-v2.css` |
| **[CRÍTICO]** | `web/components/layout/layout.templ` | 48-53 | Usa `tailwind.css` — **VERSÃO DIFERENTE** |
| **[BAIXO]** | `web/components/patient/therapeutic_plan_report.templ` | 200 | Link CSS sem versionamento dinâmico: hardcoded `"20260325_v3"` |
| **[OK]** | `cmd/arandu/main.go` | 315-334 | Cache-Control configurado corretamente com cache-busting por mtime |
| **[OK]** | Templates parciais | — | Nenhum inclui `<link>` ou `<style>` — correto |

**Detalhe do problema crítico:**
```templ
<!-- Shell Layout (NOVO) -->
<link href="/static/css/tailwind-v2.css?v=..." rel="stylesheet">

<!-- Legacy Layout (ANTIGO) -->
<link href="/static/css/tailwind.css?v=..." rel="stylesheet">
```

Classes como `bg-arandu-primary`, `text-arandu-dark` existem em `tailwind-v2.css` mas podem ter definições diferentes ou ausentes em `tailwind.css`.

---

### VETOR 2 — Estrutura dos Swaps HTMX

| Risco | Arquivo | Linha | Problema |
|-------|---------|-------|----------|
| **[CRÍTICO]** | `web/components/patient/detail.templ` | 53 | `hx-push-url="true"` sem `hx-head="merge"` |
| **[CRÍTICO]** | `web/components/dashboard/dashboard_v2.templ` | 78, 144, etc. | Mesmo padrão — URL muda mas `<head>` não é atualizado |
| **[ALTO]** | `web/components/layout/shell_layout.templ` | 411 | `hx-swap-oob="true"` pode causar race conditions |
| **[MÉDIO]** | `web/components/layout/layout.templ` | 111, 370, 645, 866 | Listeners `htmx:responseError` duplicados 4x |
| **[BAIXO]** | Global | — | Sem `hx-history` para restaurar snapshots do DOM |

**Padrão problemático identificado:**
```templ
<a hx-get="/patients/{id}"
   hx-target="#main-content"
   hx-push-url="true"
   hx-swap="innerHTML">
<!-- Falta: hx-head="merge" ou estratégia para atualizar CSS -->
```

Quando navegando de Dashboard (v2 CSS) para Patient Detail (v1 CSS), o `<head>` com CSS não é atualizado.

---

### VETOR 3 — Geração das Classes Tailwind

| Risco | Arquivo | Linha | Problema |
|-------|---------|-------|----------|
| **[MÉDIO]** | `web/static/css/input-v2.css` | 8-9 | `@source` não inclui handlers Go que geram classes dinâmicas |
| **[BAIXO]** | `internal/web/handlers/timeline_handler.go` | 232-264 | Classes `timeline-dot-{type}` geradas via concatenação — mas estão no CSS |
| **[BAIXO]** | `web/components/dashboard/dashboard_v2.templ` | 22-40 | Parâmetros `bgColor`, `textColor` dinâmicos podem gerar classes não safelistadas |
| **[INFO]** | Projeto | — | Usa Tailwind v4 com configuração via `@theme` no CSS (sem `tailwind.config.js`) |

**Configuração atual:**
```css
@source "../../../web/components/**/*.{templ,go,html}";
@source "../../../web/pages/**/*.{templ,go,html}";
/* NOTA: Não inclui internal/web/handlers/*.go */
```

---

### VETOR 4 — Servidor Go e Headers HTTP

| Risco | Arquivo | Linha | Problema |
|-------|---------|-------|----------|
| **[ALTO]** | `internal/web/handlers/dashboard_handler.go` | 100 | Não verifica `HX-Request` — sempre retorna página completa |
| **[ALTO]** | `internal/web/handlers/observation_handler.go` | 30, 63, 103 | Sem diferenciação HTMX vs página completa |
| **[ALTO]** | `internal/web/handlers/intervention_handler.go` | 30, 63, 103 | Mesmo problema |
| **[ALTO]** | `internal/web/handlers/session_handler.go` | 727, 779, 862 | `Content-Type: text/html` sem `; charset=utf-8` |
| **[MÉDIO]** | Handlers de página | — | Sem headers `Cache-Control` para diferenciação |
| **[OK]** | `cmd/arandu/main.go` | 318 | Static assets: `no-cache, no-store, must-revalidate, max-age=0` |

**Padrão ausente:**
```go
// Handler que deveria diferenciar mas não diferencia:
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
    // FALTA: isHTMX := r.Header.Get("HX-Request") == "true"
    layoutComponents.Shell(config, component).Render(ctx, w)
}
```

---

### VETOR 5 — Gerenciamento de Estado do DOM

| Risco | Arquivo | Linha | Problema |
|-------|---------|-------|----------|
| **[CRÍTICO]** | `web/components/layout/shell_layout.templ` | 283 | ID `search-results` duplicado |
| **[CRÍTICO]** | `web/components/layout/layout.templ` | 1017 | Mesmo ID `search-results` |
| **[CRÍTICO]** | `web/components/layout/shell_layout.templ` | 262 | ID `patient-search` duplicado |
| **[CRÍTICO]** | `web/components/layout/layout.templ` | 996, 1056 | Mesmo ID `patient-search` |
| **[ALTO]** | `web/components/layout/shell_layout.templ` | 184-195 | Acesso direto a `__x.$data` do Alpine.js |
| **[MÉDIO]** | `web/components/layout/layout.templ` | 54-280, 333-539, etc. | Scripts inline duplicados criam listeners múltiplos |
| **[BAIXO]** | Global | — | Sem `hx-swap="morph"` para preservar estado Alpine |

**Código problemático:**
```javascript
// Após HTMX swap, este seletor pode retornar elemento errado:
const alpineData = document.querySelector('[x-data]')?.__x?.$data;
```

---

### VETOR 6 — Estrutura dos Templates Templ

| Risco | Arquivo | Problema |
|-------|---------|----------|
| **[CRÍTICO]** | `layout.templ` vs `shell_layout.templ` | Dois sistemas de layout concorrentes com CSS diferentes |
| **[ALTO]** | Diretório `web/components/` | Sem separação de `/pages` vs `/fragments` |
| **[ALTO]** | `web/components/auth/login.templ` | 150+ linhas de `<style>` inline duplicado |
| **[ALTO]** | `web/components/auth/signup.templ` | Mesmo `<style>` inline duplicado |
| **[MÉDIO]** | `web/components/patient/therapeutic_plan_report.templ` | Layout isolado com versão CSS diferente |

**Hierarquia de templates:**
```
web/components/
├── layout/
│   ├── layout.templ         → USA: tailwind.css (LEGACY)
│   ├── shell_layout.templ  → USA: tailwind-v2.css (NEW)
│   └── ...
├── auth/
│   ├── login.templ         → inline <style> (150 linhas)
│   └── signup.templ        → inline <style> (120 linhas)
├── patient/                → Mistura: páginas + fragmentos
├── session/                → Mistura: páginas + fragmentos
└── timeline/               → Fragmentos apenas
```

---

## 3. MAPA DE CAUSA RAIZ

```
┌─────────────────────────────────────────────────────────────────────────────┐
│           CENÁRIO: CSS QUEBRADO APÓS NAVEGAÇÃO HTMX                         │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
         ┌───────────────────────────┼───────────────────────────┐
         │                           │                           │
         ▼                           ▼                           ▼
┌─────────────────────┐   ┌─────────────────────┐   ┌─────────────────────┐
│  DUAS VERSÕES CSS   │   │   IDS DUPLICADOS    │   │ HX-REQUEST NÃO      │
│  tailwind.css ≠     │   │   #search-results   │   │ VERIFICADO          │
│  tailwind-v2.css    │   │   #patient-search   │   │ EXTENSIVAMENTE      │
│                     │   │   em 3+ layouts     │   │                     │
└─────────────────────┘   └─────────────────────┘   └─────────────────────┘
         │                           │                           │
         │                           │                           │
         ▼                           ▼                           ▼
┌─────────────────────┐   ┌─────────────────────┐   ┌─────────────────────┐
│ Página com Shell    │   │ HTMX hx-target      │   │ Handler retorna     │
│ carrega v2 CSS      │   │ encontra ID errado  │   │ página completa     │
│                     │   │ após swap           │   │ sem ser fragmento   │
└─────────────────────┘   └─────────────────────┘   └─────────────────────┘
         │                           │                           │
         └───────────────────────────┼───────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  RESULTADO:                                                                 │
│  1. Classes CSS inexistentes no cache do browser                           │
│  2. Elementos HTMX alvejando containers errados                            │
│  3. Fragmento HTML servido dentro de snippet esperando layout completo     │
│  4. Hard refresh limpa cache e força recarregar CSS correto                │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Cadeia de Eventos que Causa o Problema

1. Usuário acessa `/dashboard` → carrega `shell_layout.templ` → **CSS v2 no cache**
2. Usuário navega para `/patients/{id}` via HTMX (`hx-get` + `hx-push-url`)
3. Handler retorna página com `layout.templ` → **espera CSS v1**
4. Browser mantém CSS v2 no cache (sem hard refresh)
5. Classes do CSS v1 podem não existir no v2 ou ter valores diferentes
6. **Resultado visual:** componentes sem estilo ou com layout quebrado

---

## 4. INVENTÁRIO TÉCNICO

| Elemento | Valor | Localização |
|----------|-------|-------------|
| **Versão do HTMX** | 1.9.10 (CDN) | `shell_layout.templ:48`, `layout.templ:46` |
| **Versão do Tailwind** | v4.2.2 (CLI mode) | `package.json` |
| **Arquivo de entrada CSS** | `input-v2.css` | `web/static/css/input-v2.css` |
| **Arquivo de saída CSS (v2)** | `tailwind-v2.css` | Servido via `/static/css/` |
| **Arquivo de saída CSS (v1)** | `tailwind.css` | Servido via `/static/css/` |
| **CSS custom** | `style.css` | Componentes custom `timeline-*`, `shell-*`, etc. |
| **Versionamento de CSS** | Mtime `?v={timestamp}` | `helpers.GetCSSVersion()` |
| **Serving de assets** | `http.FileServer` | `cmd/arandu/main.go:315` |
| **Cache-Control static** | `no-cache, no-store, must-revalidate, max-age=0` | Desenvolvimento |
| **Uso de `hx-boost`** | Não | Navegação via `hx-get` explícito |
| **Uso de `hx-history`** | Não | Sem configuração de snapshots |
| **Uso de `hx-push-url`** | Sim (29 ocorrências) | Sem `hx-head="merge"` |
| **Uso de `hx-swap-oob`** | 1 ocorrência | `ShellContentWrapper` |
| **Frameworks JS** | Alpine.js 3.13.5 (CDN) | `defer` no `<head>` |
| **Swap primário** | `innerHTML` (52 ocorrências) | Comum em todos os fragmentos |
| **Swap com transição** | `innerHTML transition:true` (24 ocorrências) | Dashboards e listas |

---

## 5. PADRÕES DE RISCO IDENTIFICADOS

### [RISCO 1] CSS Version Conflict Pattern

**Arquivos:** `shell_layout.templ:46-51`, `layout.templ:48-53`

```templ
<!-- Shell Layout usa v2 -->
<link href="/static/css/tailwind-v2.css?v={{ .CSSVersion }}" rel="stylesheet">

<!-- Legacy Layout usa v1 -->
<link href="/static/css/tailwind.css?v={{ .CSSVersion }}" rel="stylesheet">
```

**Impacto:** Classes definidas em uma versão podem não existir na outra ou ter valores diferentes.

**Solução recomendada:** Unificar em uma única versão de CSS.

---

### [RISCO 2] Duplicate Element ID Pattern

**Arquivos:** `shell_layout.templ:283`, `layout.templ:1017`, `search.templ:11`

```html
<!-- TRÊS layouts têm o mesmo ID -->
<div id="search-results" class="...">

<!-- TRÊS layouts têm o mesmo ID -->
<input id="patient-search" type="text" ...>
```

**Impacto:** `hx-target="#search-results"` pode gravar no container errado após navegação.

**Solução recomendada:** Usar IDs únicos com contexto: `id="dashboard-search-results"`, `id="patient-search-results"`.

---

### [RISCO 3] Missing HX-Request Check Pattern

**Arquivos:** `dashboard_handler.go:100`, `observation_handler.go:30`, `intervention_handler.go:30`

```go
// Handler NÃO diferencia HTMX de página completa
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
    layoutComponents.Shell(config, component).Render(ctx, w)
    // FALTA: if isHTMX { return fragment only }
}
```

**Impacto:** Fragmentos são servidos com `<html><head>...` completo, sujando o DOM e potencialmente injetando CSS duplicado.

**Solução recomendada:**
```go
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
    isHTMX := r.Header.Get("HX-Request") == "true"
    if isHTMX {
        component.Render(ctx, w)
        return
    }
    layoutComponents.Shell(config, component).Render(ctx, w)
}
```

---

### [RISCO 4] Alpine.js State Destruction Pattern

**Arquivo:** `shell_layout.templ:184-195`

```javascript
document.addEventListener('htmx:afterSwap', (e) => {
    const alpineData = document.querySelector('[x-data]')?.__x?.$data;
    // Após HTMX swap:
    // 1. Alpine pode ter sido re-inicializado
    // 2. __x.$data pode ser de componente diferente
    // 3. Código quebra silenciosamente
});
```

**Impacto:** Erros de JavaScript silenciosos, quebra de funcionalidade Alpine.

**Solução recomendada:** Usar Alpine.js API oficial (`Alpine.$data()`) ou passar dados via `x-model` / `Alpine.store()`.

---

### [RISCO 5] Inline Script Duplication Pattern

**Arquivo:** `layout.templ` (múltiplas funções)

```javascript
// MESMO código aparece nas linhas:
// 54-280, 333-539, 608-814, 866-1000
document.body.addEventListener('htmx:responseError', function(e) { ... });
document.body.addEventListener('htmx:swapError', function(e) { ... });
document.body.addEventListener('htmx:afterSwap', (e) => { ... });
```

**Impacto:** Múltiplos listeners para o mesmo evento → múltiplos toasts de erro, estado inconsistente.

**Solução recomendada:** Extrair para arquivo JS dedicado (`/static/js/htmx-handlers.js`), incluir uma única vez.

---

### [RISCO 6] hx-push-url Without Head Merge Pattern

**Arquivos:** `patient/detail.templ:53`, `dashboard_v2.templ:78`, etc.

```templ
<a hx-get="/patients/{id}"
   hx-target="#main-content"
   hx-push-url="true"    <!-- URL atualiza -->
   hx-swap="innerHTML">  <!-- Mas <head> NÃO atualiza -->
```

**Impacto:** O `<head>` com CSS permanece o da página anterior. Se as páginas usam layouts diferentes com CSS diferente, as classes inexistentes quebram o estilo.

**Solução recomendada:**
1. **Opção A:** Usar `hx-head="merge"` (HTMX 1.9+)
2. **Opção B:** Usar OOB swap para atualizar CSS links
3. **Opção C:** Unificar todos os layouts para usar o mesmo CSS

---

## 6. RECOMENDAÇÕES PRIORITÁRIAS

| Prioridade | Recomendação | Esforço | Impacto |
|------------|--------------|---------|---------|
| **URGENTE** | Unificar em um único layout (`Shell`) com uma única versão de CSS | Médio | Resolve 80% dos problemas |
| **URGENTE** | Gerar IDs únicos para elementos HTMX target (padrão `{context}-{element}`) | Baixo | Evita conflitos de swap |
| **ALTA** | Implementar verificação `HX-Request` consistente em todos os handlers | Médio | Fragmentos corretos |
| **ALTA** | Adicionar `hx-head="merge"` ou OOB swap para CSS links | Médio | Atualiza CSS na navegação |
| **MÉDIA** | Extrair scripts duplicados para arquivo JS dedicado | Baixo | Mais fácil manutenção |
| **MÉDIA** | Usar `hx-swap="morph:innerHTML"` para preservar estado Alpine | Médio | Melhora UX |
| **BAIXA** | Adicionar saftelist para classes Tailwind geradas dinamicamente | Baixo | Cobertura completa |
| **BAIXA** | Implementar `hx-history` para snapshots do DOM | Baixo | Navegação mais robusta |

---

## 7. ARQUIVOS AFETADOS

### Layouts (Raiz do Problema)
- `web/components/layout/layout.templ` — Layout legacy (v1 CSS)
- `web/components/layout/shell_layout.templ` — Layout novo (v2 CSS)

### Handlers com Problemas de HX-Request
- `internal/web/handlers/dashboard_handler.go`
- `internal/web/handlers/observation_handler.go`
- `internal/web/handlers/intervention_handler.go`
- `internal/web/handlers/session_handler.go`

### Componentes com IDs Duplicados
- `web/components/layout/shell_layout.templ`
- `web/components/layout/layout.templ`
- `web/components/session/search.templ`

### CSS e Configuração
- `web/static/css/input-v2.css` — Config Tailwind v4
- `web/static/css/tailwind-v2.css` — Build v2
- `web/static/css/tailwind.css` — Build v1
- `web/static/css/style.css` — CSS custom

---

## 8. PRÓXIMOS PASSOS SUGERIDOS

1. **Análise de Impacto:** Executar diff entre `tailwind.css` e `tailwind-v2.css` para identificar classes diferentes
2. **Plano de Migração:** Migrar todas as páginas para `shell_layout.templ` com CSS unificado
3. **Teste de Regressão:** Validar navegação HTMX entre todas as páginas após unificação
4. **Refactoring:** Consolidar handlers com padrão `isHTMX` check
5. **Documentação:** Atualizar `WEB_LAYER_PATTERN.md` com diretriz de layout único

---

**Fim do Relatório**