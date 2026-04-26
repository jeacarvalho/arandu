# Auditoria UX/HTMX/Tailwind — 2026-03-29

## Diagnóstico Executivo

- O front-end tem boa cobertura de HTMX em fluxos principais (busca, timeline, formulários inline), mas há inconsistências estruturais que afetam previsibilidade de interação e acessibilidade.
- Existem problemas críticos de DOM (IDs duplicados), ausência sistemática de feedback de envio em formulários HTMX e navegação HTMX com swap agressivo no `body` sem estratégia de foco/histórico.
- O Design System está híbrido: parte em utilitários Tailwind, parte em CSS custom + estilos inline extensivos, reduzindo consistência e manutenibilidade.
- Dark mode praticamente inexistente e suporte a acessibilidade parcial (bons pontos pontuais em `role/aria-*`, porém sem padrão global).

**Notas (0–10):** UX **6.1**, Acessibilidade **5.6**, Código/Manutenibilidade **5.8**.

### Top 3 showstoppers
1. IDs duplicados (`#observations-list`, `#interventions-list`, `#session-goals-list`) no mesmo template, causando `hx-target` ambíguo e comportamento potencialmente incorreto.
2. Formulários HTMX sem `hx-indicator`/`hx-disabled-elt` em pontos críticos, abrindo espaço para “rage click”, envios repetidos e falta de feedback de estado.
3. Navegação com `hx-boost` + `hx-target="body"` + `hx-swap="innerHTML"` em ações principais, sem gestão explícita de foco e sem `hx-push-url`, prejudicando histórico/back button e a11y.

## Relatório Detalhado de Achados

### 1) IDs duplicados quebrando alvo de atualização HTMX
- **Local:** `web/components/session/edit_form.templ`
- **Categoria:** UX / A11y / Perf
- **Problema:** O template renderiza múltiplos elementos com os mesmos IDs (`observations-list`, `interventions-list`, `session-goals-list`) em áreas principal e sidebar. Isso torna `hx-target="#..."` ambíguo, podendo atualizar o bloco errado e criar inconsistência visual.
- **Solução proposta:** Garantir IDs únicos por região (`-main`, `-sidebar`) e alinhar cada formulário ao target correspondente.
- **Snippet de correção:**
```templ
<div id="observations-list-main" class="mb-md">...</div>
<form
  hx-post={ templ.URL("/session/" + data.SessionID + "/observations") }
  hx-target="#observations-list-main"
  hx-swap="beforeend"
  hx-indicator="#obs-main-loading"
  hx-disabled-elt="this button[type='submit']"
>
  ...
  <div id="obs-main-loading" class="htmx-indicator" aria-live="polite">Salvando...</div>
</form>
```

### 2) Falta de loading/disable em formulários HTMX
- **Local:** `web/components/session/observation_form_inline.templ`, `web/components/session/intervention_form_inline.templ`, `web/components/session/observation_edit_form.templ`
- **Categoria:** UX
- **Problema:** Os formulários fazem `hx-post/hx-put` sem `hx-indicator` e sem `hx-disabled-elt`. Em latência alta, usuário não recebe feedback e pode reenviar várias vezes.
- **Solução proposta:** Adicionar indicador visual + desabilitar botão durante request. Priorizar solução nativa HTMX.
- **Snippet de correção:**
```templ
<form
  hx-post={ templ.URL("/session/" + sessionID + "/observations") }
  hx-target="#observations-list"
  hx-swap="beforeend"
  hx-indicator="#obs-loading"
  hx-disabled-elt="find button[type='submit']"
  hx-on::after-request="if(event.detail.successful) this.reset()"
>
  ...
  <button type="submit" class="btn btn-primary btn-sm">Adicionar</button>
  <span id="obs-loading" class="htmx-indicator text-xs text-neutral-500" aria-live="polite">Enviando...</span>
</form>
```

### 3) Navegação HTMX agressiva no `body` sem foco/histórico
- **Local:** `web/components/patient/detail.templ` (links de ação rápida)
- **Categoria:** UX / A11y
- **Problema:** Uso recorrente de `hx-boost="true" hx-target="body" hx-swap="innerHTML"` sem `hx-push-url` e sem restauração de foco pós-swap. Isso compromete navegação por teclado e previsibilidade do botão “Voltar”.
- **Solução proposta:** Preferir target semântico (`main` ou container de página), usar `hx-push-url="true"`, e mover foco no `htmx:afterSwap` para heading primário.
- **Snippet de correção:**
```templ
<a
  href={ templ.URL("/patients/" + patient.ID + "/history") }
  hx-boost="true"
  hx-target="main"
  hx-swap="innerHTML transition:true"
  hx-push-url="true"
>
  Ver Histórico
</a>

<script>
document.body.addEventListener('htmx:afterSwap', (e) => {
  const h1 = e.target.querySelector('h1, [data-autofocus]');
  if (h1) h1.setAttribute('tabindex', '-1'), h1.focus();
});
</script>
```

### 4) Interceptação global de links depende de `e.target.tagName === 'A'`
- **Local:** `web/components/layout/layout.templ` (script de transição)
- **Categoria:** UX / Perf
- **Problema:** O listener só intercepta clique quando o target direto é `<a>`. Cliques em `<i>`/`<span>` dentro do link não entram na condição, criando comportamento inconsistente de navegação.
- **Solução proposta:** Usar `closest('a')`, filtrar modificadores e links externos.
- **Snippet de correção:**
```js
document.addEventListener('click', (e) => {
  const link = e.target.closest('a[href]');
  if (!link || link.hasAttribute('hx-boost') || e.metaKey || e.ctrlKey) return;
  if (link.origin !== window.location.origin) return;
  e.preventDefault();
  document.body.style.opacity = '0.7';
  setTimeout(() => window.location.href = link.href, 200);
});
```

### 5) Componente com alto acoplamento visual via `style="..."` inline
- **Local:** `web/components/patient/detail.templ`
- **Categoria:** Tailwind / Manutenibilidade
- **Problema:** Grande volume de estilos inline (incluindo `onmouseover/onmouseout`) dificulta consistência visual, manutenção, dark mode e reutilização.
- **Solução proposta:** Migrar para utilitários Tailwind + classes utilitárias semânticas no `@layer components` (quando repetição justificar).
- **Snippet de correção:**
```templ
<a
  href={ templ.URL("/patients/" + patient.ID + "/sessions/new") }
  class="inline-flex items-center gap-2 rounded-lg bg-arandu-primary px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-arandu-active focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-arandu-soft"
>
  <i class="fas fa-calendar-plus" aria-hidden="true"></i>
  Nova Sessão
</a>
```

### 6) Resultado de busca sem região viva consistente
- **Local:** `web/components/layout/layout.templ` + `web/components/patient/search.templ`
- **Categoria:** A11y
- **Problema:** O container em layout não define `role="status"`/`aria-live`, e resultados usam `role=listbox/option` sem estratégia completa de navegação por teclado (`aria-activedescendant`, setas, estado ativo).
- **Solução proposta:** Definir região viva no container e complementar semântica para autocomplete.
- **Snippet de correção:**
```templ
<input
  ...
  role="combobox"
  aria-expanded="true"
  aria-controls="search-results"
  aria-autocomplete="list"
/>
<div id="search-results" role="status" aria-live="polite" class="search-results-container"></div>
```

### 7) Dark mode não implementado de forma sistêmica
- **Local:** `web/components/*`, `web/static/css/input.css`
- **Categoria:** Tailwind / A11y
- **Problema:** Não há uso consistente de variantes `dark:`; paleta e componentes não contemplam contraste em tema escuro.
- **Solução proposta:** Adotar estratégia `class="dark"` no root + tokens pares (light/dark) + variantes `dark:` nos componentes críticos.
- **Snippet de correção:**
```templ
<div class="bg-white text-neutral-800 dark:bg-neutral-900 dark:text-neutral-100">
  ...
</div>
```

## Plano de Ação Priorizado

| Prioridade | Ação | Esforço Estimado |
|---|---|---|
| Alta | Eliminar IDs duplicados em `edit_form.templ` e ajustar todos os `hx-target` associados. | Médio |
| Alta | Padronizar `hx-indicator` + `hx-disabled-elt` em todos os forms HTMX de criação/edição. | Médio |
| Alta | Revisar links com `hx-target="body"` para `main/container`, adicionar `hx-push-url="true"` e política de foco pós-swap. | Médio |
| Alta | Corrigir listener global de navegação para usar `closest('a')` e evitar comportamento inconsistente. | Baixo |
| Média | Reduzir estilos inline no `patient/detail.templ`, migrando para utilitários Tailwind e classes reutilizáveis. | Alto |
| Média | Completar acessibilidade da busca com `combobox` + live region + navegação por teclado. | Médio |
| Média | Definir padrão de tratamento de erro HTMX (fragments + `hx-on::responseError`/status code UX). | Médio |
| Baixa | Organizar ordem das classes Tailwind (layout → spacing → typography → color/effects). | Baixo |
| Baixa | Introduzir dark mode progressivo em páginas de maior uso (dashboard, paciente, sessão). | Alto |
| Baixa | Revisar microinterações para reduzir layout shift e padronizar transições. | Baixo |
