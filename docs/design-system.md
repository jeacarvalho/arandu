# Design System — Arandu (Sábio)

**Última atualização:** 22/04/2026

## Filosofia

Arandu é uma ferramenta de reflexão clínica.

Princípios:
- clareza > estética
- calma > impacto visual
- conteúdo > interface
- **Tecnologia Silenciosa**: a interface desaparece; o terapeuta e o paciente permanecem

---

## Stack UI

| Tecnologia | Versão | Uso |
|-----------|--------|-----|
| Go + Templ | 0.3.x | Templating type-safe |
| HTMX | 2.x (local) | Interatividade sem SPA |
| Alpine.js | 3.13.5 | Estado mínimo no cliente |
| Tailwind CSS | **v4** | Utilitários — build via `input-v2.css` |
| CSS customizado | — | Design tokens Sábio em `style.css` |

---

## Tokens de Design (CSS Variables)

Definidos em `web/static/css/style.css`:

```css
/* Superfícies */
--color-paper:    #FAFAF8   /* fundo principal — papel levemente quente */
--color-paper-2:  #F5F4F1   /* fundo alternativo */
--color-line:     #E8E6E0   /* divisores */

/* Tipografia */
--color-ink:      #1C1917   /* texto principal */
--color-ink-2:    #44403C   /* texto secundário */
--color-ink-3:    #78716C   /* metadados, labels */
--color-ink-4:    #A8A29E   /* placeholder, desabilitado */

/* Acento */
--accent:         #2563EB   /* azul clínico (ações primárias) */
--accent-deep:    #1D4ED8   /* hover/active */

/* Semânticas */
--color-ok:       #16A34A   /* sucesso */
--color-warn:     #D97706   /* aviso */
--color-danger:   #DC2626   /* erro */

/* Fontes */
--font-serif:     'Source Serif 4', Georgia, serif
--font-sans:      'Inter', system-ui, sans-serif
--font-mono:      'Geist Mono', 'Fira Code', monospace
```

---

## Regra Tipográfica

| Contexto | Fonte | Classe Templ |
|----------|-------|-------------|
| Labels, botões, navegação, metadados | Inter | `class="sans"` ou padrão |
| Observações clínicas, notas, prontuário | Source Serif 4 | `class="serif"` |
| Timestamps, IDs, código | Geist Mono | `class="mono"` |

---

## Layout Shell

```
sabio-topbar (position: fixed, height: 64px, z-index: 50)
shell (height: 100vh)
  shell-body (margin-top: 64px, flex row)
    shell-sidebar (width: 232px / 68px colapsado)
    shell-main (flex: 1, overflow-y: auto)
      shell-canvas (flex: 1, padding: 32px)
        shell-canvas-container (mx-auto, max-width: 1440px)
          shell-content #main-content
            @content
```

**Componente:** `layout.Shell(config, content)` em `web/components/layout/shell_layout.templ`

---

## Componentes Sábio (Prefixo `sabio-`)

### Topbar
```css
.sabio-topbar       /* header fixo com blur */
.sabio-breadcrumb   /* trilha de navegação */
.sabio-topbar-spacer
```

### Sidebar
```css
.sabio-brand-header
.sabio-nav
.sabio-nav-item
.sabio-nav-item-icon / .sabio-nav-item-text
.shell-sidebar-collapsed  /* estado colapsado (Alpine store) */
```

### Botões
```css
.sabio-btn               /* base */
.sabio-btn--primary      /* ação principal */
.sabio-btn--ghost        /* ação secundária */
.sabio-btn--accent       /* gradiente azul (IA/highlight) */
```

### Pills / Tags
```css
.sabio-pill              /* base */
.sabio-pill--neutral     /* observação (cinza) */
.sabio-pill--accent      /* intervenção (azul) */
.sabio-pill--ok          /* confirmado (verde) */
.sabio-pill--warn        /* aviso (âmbar) */
```

### Cards
```css
.sabio-card
.sabio-card-header
.sabio-card-eyebrow      /* label uppercase pequeno */
.sabio-card-title        /* título serif */
```

### Sessão (Edit Form)
```css
.sabio-session-page
.sabio-session-header
.sabio-session-eyebrow
.sabio-session-title
.sabio-notes-grid        /* 2 colunas */
.sabio-notes-col
.sabio-notes-col-header
.sabio-notes-items
.sabio-notes-item
.sabio-notes-item--intervention
.sabio-notes-input
.sabio-notes-textarea
.sabio-notes-submit
```

### Dashboard
```css
.sabio-dash
.sabio-kpi-grid          /* KPIs no topo */
.sabio-kpi-card
.sabio-patient-list
.sabio-patient-item
.sabio-patient-tag       /* badge clínico do paciente */
.sabio-patient-meta      /* sessões + última data */
.sabio-patient-next      /* próximo agendamento */
```

### LLM Drawer (⌘J)
```css
.sabio-llm-drawer
.sabio-llm-backdrop
.sabio-llm-drawer--open
.sabio-llm-insights
.sabio-llm-insight-row
.sabio-llm-insight-bar
```

---

## Paleta de Cores

### Cores Sábio (identidade)
- `--color-arandu-primary: #0F6E56` — Verde Arandu
- `--color-arandu-active:  #1D9E75` — Destaque/Interação
- `--color-arandu-dark:    #085041` — Verde escuro

### Cores de superfície (Sábio design)
- `--color-paper: #FAFAF8` — Fundo principal
- Fundo global alternativo: `#F7F8FA` (herança Tailwind)
- Superfícies de cards: branco `#FFFFFF`

---

## Tailwind v4: Regras Críticas

1. **Classes dinâmicas**: Tailwind v4 não detecta classes construídas por concatenação em runtime. Use funções helper em `types.go` que retornam strings completas.

2. **Tokens `arandu-*`**: use `arandu-primary`, `arandu-active` — não `primary-*`.

3. **Texto visível sempre**: para texto que precisa ser visível independente do build do Tailwind, use `style="color: var(--color-ink)"` inline.

---

## Histórico

| Data | Alteração |
|------|-----------|
| 2026-03-01 | Versão inicial (CSS puro) |
| 2026-04-22 | Reescrita para refletir sistema Sábio + Tailwind v4 |
