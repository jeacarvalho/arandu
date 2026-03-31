# 🔍 RELATÓRIO DE AUDITORIA — CSS/Cache em HTMX + Templ + Tailwind

**Data da Auditoria:** 31/03/2026  
**Sintoma Reportado:** CSS quebrado após navegação HTMX (resolvido apenas com hard refresh Ctrl+Shift+R)  
**Metodologia:** Análise forense de 6 vetores de risco conforme especificação

---

## 1. SUMÁRIO EXECUTIVO

O problema de "CSS quebrado após navegação HTMX" tem **múltiplas causas raízes inter-relacionadas**, sendo a principal delas a coexistência de **dois layouts diferentes com versões distintas de CSS Tailwind**.

| Criticidade | Causa Raiz Principal |
|-------------|---------------------|
| **[CRÍTICO]** | Dois layouts (`layout.templ` e `shell_layout.templ`) servem arquivos CSS diferentes (`tailwind.css` vs `tailwind-v2.css`) |
| **[CRÍTICO]** | Navegação via `hx-target="main"` aponta para elemento inexistente no shell layout (`id="main-content"`) |
| **[ALTO]** | Handlers inconsistentes na verificação `HX-Request` — alguns retornam página completa quando deveriam retornar fragmentos |
| **[ALTO]** | Scripts inline duplicados 4x nos layouts criam listeners HTMX duplicados |

**Conclusão:** Quando o usuário navega entre páginas que usam layouts diferentes, o browser mantém em cache o CSS da primeira página. Fragmentos HTML injetados via HTMX contêm classes da nova página que podem não existir (ou ter definições diferentes) no CSS cachedo. O hard refresh força recarregamento completo do CSS, resolvendo temporariamente.

---

## 2. ACHADOS POR VETOR

### VETOR 1 — Inclusão do CSS (stylesheet e Tailwind)

| Risco | Arquivo | Linha | Descrição |
|-------|---------|-------|-----------|
| **[CRÍTICO]** | `web/components/layout/shell_layout.templ` | 46-51 | Usa `/static/css/tailwind-v2.css?v=...` |
| **[CRÍTICO]** | `web/components/layout/layout.templ` | 48-53 | Usa `/static/css/tailwind.css?v=...` — **VERSÃO DIFERENTE** |
| **[BAIXO]** | `web/components/patient/therapeutic_plan_report.templ` | 200 | Link CSS hardcoded sem versionamento dinâmico: `"20260325_v3"` |
| **[OK]** | `cmd/arandu/main.go` | 315-334 | Cache-Control configurado corretamente: `no-cache, no-store, must-revalidate, max-age=0` |
| **[OK]** | Templates parciais | — | Nenhum inclui `<link>` ou `<style>` — correto |

**Detalhe do problema crítico:**
```templ
<!-- Shell Layout (NOVO) -->
<link href="/static/css/tailwind-v2.css?v={{ .CSSVersion }}" rel="stylesheet">

<!-- Legacy Layout (ANTIGO) -->
<link href="/static/css/tailwind.css?v={{ .CSSVersion }}" rel="stylesheet">
```

Classes como `bg-arandu-primary`, `text-arandu-dark`, `font-clinical` podem ter definições diferentes entre as duas versões.

---

### VETOR 2 — Estrutura dos Swaps HTMX

| Risco | Arquivo | Linha | Descrição |
|-------|---------|-------|-----------|
| **[CRÍTICO]** | `web/components/patient/detail.templ` | 52 | `hx-target="main"` — alvo incompatível com shell_layout |
| **[CRÍTICO]** | `web/components/patient/detail.templ` | 248, 259, 270 | Mesmo padrão problemático |
| **[CRÍTICO]** | `web/components/patient/profile.templ` | 26, 140, 144, 190 | `hx-target="#main-content"` — só existe em shell_layout |
| **[ALTO]** | `web/components/layout/shell_layout.templ` | 298 | `hx-swap="innerHTML"` sem `hx-head="merge"` |
| **[ALTO]** | Global | — | **Nenhum template usa `hx-head="merge"`** — head não atualiza na navegação |
| **[MÉDIO]** | `web/components/layout/layout.templ` | 111, 370, 645, 866 | Listeners `htmx:responseError` duplicados 4x |
| **[MÉDIO]** | `web/components/layout/shell_layout.templ` | 309 | `hx-swap-oob="true"` pode causar race conditions |

**Padrão problemático identificado:**
```templ
<!-- patient/detail.templ:52 -->
<a hx-get="/patients/{id}"
   hx-target="main"              <!-- Alvo "main" não existe em shell_layout! -->
   hx-push-url="true"
   hx-swap="innerHTML transition:true">
<!-- Falta: hx-head="merge" -->
```

**Incompatibilidade de targets:**
- `layout.templ`: `<main class="app-container app-main-content">` (sem ID)
- `shell_layout.templ`: `<div class="shell-canvas" id="main-content">`

Quando navegando de uma página com legacy layout para shell layout (ou vice-versa), o HTMX não encontra o target correto.

---

### VETOR 3 — Geração das Classes Tailwind

| Risco | Arquivo | Linha | Descrição |
|-------|---------|-------|-----------|
| **[MÉDIO]** | `web/static/css/input-v2.css` | 8-9 | `@source` não inclui handlers Go que geram classes dinâmicas |
| **[BAIXO]** | `internal/web/handlers/timeline_handler.go` | 232-264 | Classes `timeline-dot-{type}` geradas via concatenação — mas estão no CSS |
| **[BAIXO]** | `web/components/dashboard/dashboard_v2.templ` | 22-40 | Parâmetros `bgColor`, `textColor` dinâmicos podem gerar classes não safelistadas |
| **[INFO]** | Projeto | — | Usa Tailwind v4 com configuração via `@theme` no CSS (sem `tailwind.config.js`) |

**Configuração atual do Tailwind v4:**
```css
/* input-v2.css */
@source "../../../web/components/**/*.{templ,go,html}";
@source "../../../web/pages/**/*.{templ,go,html}";
/* NOTA: Não inclui internal/web/handlers/*.go */
```

---

### VETOR 4 — Servidor Go e Headers HTTP

| Risco | Arquivo | Linha | Descrição |
|-------|---------|-------|-----------|
| **[ALTO]** | `internal/web/handlers/observation_handler.go` | 30-67 | **Sem verificação `HX-Request`** — sempre retorna componente puro |
| **[ALTO]** | `internal/web/handlers/intervention_handler.go` | 30-67 | **Sem verificação `HX-Request`** — mesmo problema |
| **[ALTO]** | `internal/web/handlers/ai_handler.go` | — | **Sem verificação `HX-Request`** |
| **[MÉDIO]** | `internal/web/handlers/dashboard_handler.go` | 104 | Verifica `HX-Request` mas retorna layout completo mesmo para HTMX |
| **[OK]** | `internal/web/handlers/patient_handler.go` | 181, 279, 655 | Verificação correta de `isHTMX` |
| **[OK]** | `internal/web/handlers/session_handler.go` | 145, 246, 415 | Verificação correta de `isHTMX` |
| **[OK]** | `internal/web/handlers/timeline_handler.go` | 88 | Verificação correta de `isHTMXRequest` |
| **[OK]** | `cmd/arandu/main.go` | 318 | Static assets: `Cache-Control: no-cache, no-store, must-revalidate, max-age=0` |

**Exemplo de handler problemático:**
```go
// observation_handler.go:30-67
func (h *ObservationHandler) GetObservation(w http.ResponseWriter, r *http.Request) {
    // FALTA: isHTMX := r.Header.Get("HX-Request") == "true"
    // FALTA: Lógica diferenciada para HTMX vs página completa
    
    component := sessionComponents.ObservationItem(...)
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    component.Render(r.Context(), w)
}
```

---

### VETOR 5 — Gerenciamento de Estado do DOM

| Risco | Arquivo | Linha | Descrição |
|-------|---------|-------|-----------|
| **[ALTO]** | `web/components/layout/layout.templ` | 51 | Alpine.js 3.13.5 carregado com `defer` — pode não estar pronto antes do HTMX |
| **[ALTO]** | `web/components/layout/shell_layout.templ` | 47 | Mesmo problema com Alpine.js |
| **[MÉDIO]** | `web/components/layout/layout.templ` | 131-143 | Alpine.store() para sidebar — estado pode ser perdido em swaps |
| **[MÉDIO]** | `web/components/layout/shell_layout.templ` | 184-195 | Código acessa `__x.$data` diretamente — frágil após swaps |
| **[BAIXO]** | Global | — | Sem `hx-history` configurado — snapshots do DOM não são salvos |
| **[BAIXO]** | Global | — | Sem uso de `idomorph` ou morph swap — elementos são substituídos brutalmente |

**Problema com Alpine.js:**
```javascript
// shell_layout.templ:184-195
document.addEventListener('htmx:afterSwap', (e) => {
    const alpineData = document.querySelector('[x-data]')?.__x?.$data;
    // Problema: __x pode ser null após swap, ou ser de componente diferente
});
```

---

### VETOR 6 — Estrutura dos Templates Templ

| Risco | Arquivo | Linha | Descrição |
|-------|---------|-------|-----------|
| **[CRÍTICO]** | `web/components/layout/layout.templ` | 23-284 | Função `BaseWithContentAndEmailAndSidebar` — layout legacy |
| **[CRÍTICO]** | `web/components/layout/shell_layout.templ` | 1-350+ | Função `Shell` — layout novo com estrutura HTML diferente |
| **[ALTO]** | `web/components/patient/detail.templ` | 1-280+ | Assume target `main` — incompatível com shell_layout |
| **[ALTO]** | `web/components/patient/profile.templ` | 1-200+ | Assume target `#main-content` — só existe em shell_layout |
| **[MÉDIO]** | `web/components/layout/layout.templ` | 54-280, 333-539, 608-814, 866-1000 | **Mesmo script inline duplicado 4 vezes** |
| **[MÉDIO]** | `web/components/patient/search.templ` | 11 | `id="search-results"` — ID genérico pode conflitar |

**Duplicação de scripts (layout.templ):**
O mesmo código JavaScript aparece nas linhas:
- 54-280 (função `BaseWithContentAndEmailAndSidebar`)
- 333-539 (função `BaseWithContentAndEmail`)
- 608-814 (outra variante)
- 866-1000 (outra variante)

Isso cria múltiplos listeners para os mesmos eventos HTMX.

---

## 3. MAPA DE CAUSA RAIZ

### Cadeia de Eventos que Leva ao Problema

```
1. Usuário acessa Página A (ex: Dashboard)
   → Layout: shell_layout.templ
   → CSS: tailwind-v2.css (carregado no <head>)
   → Target: id="main-content"

2. Usuário clica em link para Página B (ex: Patient Detail)
   → HTMX faz GET /patients/{id} com header HX-Request: true
   → Handler retorna fragmento HTML com classes Tailwind v2
   → HTMX injeta em hx-target="main" (não "#main-content"!)

3. PROBLEMA #1: Target mismatch
   → "main" não existe em shell_layout (que usa id="main-content")
   → HTMX pode falhar silenciosamente ou injetar no lugar errado

4. PROBLEMA #2: CSS cache conflict
   → Browser já tem tailwind-v2.css em cache
   → Se Página B foi renderizada originalmente com layout.templ (tailwind.css)
   → Classes específicas de uma versão podem não existir na outra

5. PROBLEMA #3: Head não atualiza
   → hx-push-url="true" muda URL
   → Mas <head> NÃO é atualizado (falta hx-head="merge")
   → CSS links permanecem os da página anterior

6. Resultado: Elementos aparecem sem estilo ou com layout errado

7. Hard refresh (Ctrl+Shift+R):
   → Força recarregamento completo do CSS
   → Todas as classes são重新 definidas
   → Problema resolvido temporariamente
```

---

## 4. INVENTÁRIO TÉCNICO

| Elemento | Valor | Localização |
|----------|-------|-------------|
| **Versão do HTMX** | 1.9.10 (CDN) | `shell_layout.templ:48`, `layout.templ:46` |
| **Versão do Tailwind** | v4.2.2 (CLI mode) | `package.json` |
| **Arquivo de entrada CSS** | `input-v2.css` | `web/static/css/input-v2.css` |
| **Arquivo de saída CSS (v2)** | `tailwind-v2.css` (51KB) | Servido via `/static/css/` |
| **Arquivo de saída CSS (v1)** | `tailwind.css` (51KB) | Servido via `/static/css/` |
| **CSS custom** | `style.css` (155KB) | Componentes custom `timeline-*`, `shell-*`, etc. |
| **Versionamento de CSS** | Mtime `?v={timestamp}` | `helpers.GetCSSVersion()` |
| **Serving de assets** | `http.FileServer` | `cmd/arandu/main.go:315` |
| **Cache-Control static** | `no-cache, no-store, must-revalidate, max-age=0` | Desenvolvimento |
| **Uso de `hx-boost`** | Sim (11 ocorrências) | `patient/detail.templ`, `profile.templ`, `signup.templ` |
| **Uso de `hx-history`** | Não | Sem configuração de snapshots |
| **Uso de `hx-push-url`** | Sim (29+ ocorrências) | **Sem `hx-head="merge"` em nenhum** |
| **Uso de `hx-swap-oob`** | 1 ocorrência | `ShellContentWrapper` |
| **Frameworks JS** | Alpine.js 3.13.5 (CDN) | `defer` no `<head>` |
| **Swap primário** | `innerHTML` (52+ ocorrências) | Comum em todos os fragmentos |
| **Swap com transição** | `innerHTML transition:true` (24+ ocorrências) | Dashboards e listas |
| **Outros swaps** | `outerHTML`, `afterbegin`, `beforeend`, `afterend` | Formulários e listas |

---

## 5. PADRÕES DE RISCO IDENTIFICADOS

### [RISCO 1] Dual CSS Version Pattern ⚠️ CRÍTICO

**Arquivos:** `shell_layout.templ:46-51`, `layout.templ:48-53`

```templ
<!-- Shell Layout usa v2 -->
<link href="/static/css/tailwind-v2.css?v={{ .CSSVersion }}" rel="stylesheet">

<!-- Legacy Layout usa v1 -->
<link href="/static/css/tailwind.css?v={{ .CSSVersion }}" rel="stylesheet">
```

**Impacto:** Classes definidas em uma versão podem não existir na outra ou ter valores diferentes.

**Solução recomendada:** Unificar em uma única versão de CSS e migrar todos os layouts.

---

### [RISCO 2] Incompatible HTMX Target Pattern ⚠️ CRÍTICO

**Arquivos:** `patient/detail.templ:52`, `shell_layout.templ:91`, `layout.templ:69`

```templ
<!-- patient/detail.templ assume target "main" -->
<a hx-target="main" ...>

<!-- shell_layout.templ tem id="main-content" -->
<div class="shell-canvas" id="main-content">

<!-- layout.templ tem <main> sem ID -->
<main class="app-container app-main-content">
```

**Impacto:** HTMX não encontra o target correto, injeta conteúdo no lugar errado ou falha silenciosamente.

**Solução recomendada:** Padronizar todos os templates para usar o mesmo target (ex: `id="main-content"`).

---

### [RISCO 3] Missing HX-Request Check Pattern ⚠️ ALTO

**Arquivos:** `observation_handler.go:30`, `intervention_handler.go:30`, `ai_handler.go:*`

```go
// Handler NÃO diferencia HTMX de página completa
func (h *Handler) GetObservation(w http.ResponseWriter, r *http.Request) {
    // FALTA: isHTMX := r.Header.Get("HX-Request") == "true"
    // FALTA: if isHTMX { return component only }
    
    component := sessionComponents.ObservationItem(...)
    component.Render(ctx, w)
}
```

**Impacto:** Handlers retornam componentes puros sem contexto de layout, mas também sem otimização para HTMX.

**Solução recomendada:**
```go
func (h *Handler) GetObservation(w http.ResponseWriter, r *http.Request) {
    isHTMX := r.Header.Get("HX-Request") == "true"
    
    component := sessionComponents.ObservationItem(...)
    
    if isHTMX {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        component.Render(ctx, w)
        return
    }
    
    // Para página completa, envolver com layout
    layoutComponents.Shell(config, component).Render(ctx, w)
}
```

---

### [RISCO 4] Missing hx-head="merge" Pattern ⚠️ ALTO

**Arquivos:** 29+ ocorrências de `hx-push-url="true"` sem `hx-head`

```templ
<a hx-get="/patients/{id}"
   hx-target="#main-content"
   hx-push-url="true"    <!-- URL atualiza -->
   hx-swap="innerHTML">  <!-- Mas <head> NÃO atualiza -->
<!-- FALTA: hx-head="merge" -->
```

**Impacto:** O `<head>` com CSS permanece o da página anterior. Se as páginas usam layouts diferentes com CSS diferente, as classes inexistentes quebram o estilo.

**Solução recomendada:** Adicionar `hx-head="merge"` em todos os links com `hx-push-url="true"`.

---

### [RISCO 5] Inline Script Duplication Pattern ⚠️ MÉDIO

**Arquivo:** `layout.templ` (múltiplas funções)

```javascript
// MESMO código aparece 4 vezes:
// Linhas 54-280, 333-539, 608-814, 866-1000
document.body.addEventListener('htmx:responseError', function(e) { ... });
document.body.addEventListener('htmx:swapError', function(e) { ... });
document.body.addEventListener('htmx:afterSwap', (e) => { ... });
```

**Impacto:** Múltiplos listeners para o mesmo evento → múltiplos toasts de erro, estado inconsistente, performance degradada.

**Solução recomendada:** Extrair para arquivo JS dedicado (`/static/js/htmx-handlers.js`), incluir uma única vez no layout.

---

### [RISCO 6] Alpine.js State Destruction Pattern ⚠️ MÉDIO

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

### [RISCO 7] Duplicate Search Element IDs ⚠️ BAIXO

**Arquivos:** `layout.templ:1021`, `shell_layout.templ:181`, `search.templ:11`

```html
<!-- TRÊS IDs diferentes para mesma funcionalidade -->
<div id="legacy-search-results">     <!-- layout.templ -->
<div id="shell-search-results">      <!-- shell_layout.templ -->
<div id="search-results">            <!-- search.templ -->
```

**Impacto:** Potencial confusão se algum template referenciar ID errado. Menos grave pois cada layout tem seu próprio contexto.

**Solução recomendada:** Manter convenção de nomes contextuais (já está correto).

---

## 6. RECOMENDAÇÕES PRIORITÁRIAS

| Prioridade | Recomendação | Esforço | Impacto |
|------------|--------------|---------|---------|
| **URGENTE** | Unificar em um único layout (`Shell`) com uma única versão de CSS | Médio | Resolve 80% dos problemas |
| **URGENTE** | Corrigir todos os `hx-target` para usar `id="main-content"` consistentemente | Baixo | Evita falhas de swap |
| **URGENTE** | Adicionar `hx-head="merge"` em todos os links com `hx-push-url="true"` | Baixo | Atualiza CSS na navegação |
| **ALTA** | Implementar verificação `HX-Request` consistente em observation/intervention/ai handlers | Médio | Fragmentos corretos |
| **ALTA** | Extrair scripts duplicados para `/static/js/htmx-handlers.js` | Baixo | Mais fácil manutenção |
| **MÉDIA** | Padronizar Alpine.js state management com `Alpine.store()` | Médio | Melhora robustez |
| **BAIXA** | Adicionar saftelist para classes Tailwind geradas dinamicamente | Baixo | Cobertura completa |
| **BAIXA** | Implementar `hx-history="true"` para snapshots do DOM | Baixo | Navegação mais robusta |

---

## 7. ARQUIVOS AFETADOS

### Layouts (Raiz do Problema)
- `web/components/layout/layout.templ` — Layout legacy (v1 CSS, 4x scripts duplicados)
- `web/components/layout/shell_layout.templ` — Layout novo (v2 CSS, estrutura diferente)

### Handlers com Problemas de HX-Request
- `internal/web/handlers/observation_handler.go` — Sem verificação HX-Request
- `internal/web/handlers/intervention_handler.go` — Sem verificação HX-Request
- `internal/web/handlers/ai_handler.go` — Sem verificação HX-Request
- `internal/web/handlers/dashboard_handler.go` — Verifica mas retorna layout completo

### Componentes com Targets Incompatíveis
- `web/components/patient/detail.templ` — Usa `hx-target="main"` (inexistente)
- `web/components/patient/profile.templ` — Usa `hx-target="#main-content"` (só em shell)

### CSS e Configuração
- `web/static/css/input-v2.css` — Config Tailwind v4
- `web/static/css/tailwind-v2.css` — Build v2 (51KB)
- `web/static/css/tailwind.css` — Build v1 (51KB)
- `web/static/css/style.css` — CSS custom (155KB)
- `internal/platform/helpers/css_version.go` — Versionamento por mtime

---

## 8. PRÓXIMOS PASSOS SUGERIDOS

1. **Análise de Impacto Imediata:**
   - Executar diff entre `tailwind.css` e `tailwind-v2.css` para identificar classes diferentes
   - Listar todas as páginas que usam cada layout

2. **Correções Urgentes (Sprint 1):**
   - Migrar todas as páginas para `shell_layout.templ`
   - Remover `layout.templ` e `tailwind.css` do código
   - Corrigir todos os `hx-target` para `id="main-content"`
   - Adicionar `hx-head="merge"` em links com push-url

3. **Refactoring (Sprint 2):**
   - Consolidar handlers com padrão `isHTMX` check
   - Extrair scripts para arquivo JS dedicado
   - Padronizar Alpine.js state management

4. **Teste de Regressão:**
   - Validar navegação HTMX entre todas as páginas
   - Testar back/forward do browser
   - Validar CSS em todas as rotas

5. **Documentação:**
   - Atualizar padrões de desenvolvimento web
   - Documentar convenção de HTMX targets

---

**Fim do Relatório de Auditoria**

*Este relatório serve como base para um prompt de correção cirúrgica a ser aplicado posteriormente.*
