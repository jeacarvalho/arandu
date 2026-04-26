# Task: CA-08 — Shell redesign fiel ao Sábio
Requirement: Redesign visual — handoff `design_handoff_arandu_redesign/`
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Objetivo

Fazer o shell (sidebar + topbar + brandmark) corresponder exatamente ao design "Sábio" proposto
em `design_handoff_arandu_redesign/chrome.jsx` e `README.md`.

Hoje o shell tem sidebar verde (diverge do design) e topbar com cores erradas (`bg-base-100`
DaisyUI ≠ `--paper` Sábio). Esta task corrige a fundação visual — todas as outras telas
dependem desta estar certa.

**Referência**: abra `design_handoff_arandu_redesign/Arandu Redesign.html` no browser para ver
o alvo visual exato.

---

## Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX · DaisyUI v5 + Tailwind CSS v4 · Alpine.js 3
**Paleta Sábio**: tokens já definidos em `web/static/css/style.css` (--paper, --ink, --accent, etc.)
**DaisyUI drawer** já implementado em `web/components/layout/shell_layout.templ` (CA-05)
**CA-07** (pendente ou concluída): removeu `header { position: fixed !important }` de style.css

---

## Arquivos a modificar

```
web/static/css/input-v2.css          ← DaisyUI theme + classes .serif/.mono
web/components/layout/shell_layout.templ  ← topbar (sticky), sidebar (cream), BrandMark
web/static/css/style.css             ← ajuste fino de .sabio-sidebar, .sabio-topbar (se necessário)
```

---

## Mudança 1 — DaisyUI theme: alinhar com paleta Sábio

Em `web/static/css/input-v2.css`, substituir o bloco `[data-theme="arandu"]` atual pelo seguinte
(troca oklch verde por valores Sábio terrosos):

```css
[data-theme="arandu"] {
  color-scheme: light;

  /* Base — papel cru Sábio */
  --color-base-100:    #F5EFE6;   /* --paper */
  --color-base-200:    #EDE5D8;   /* --paper-2 */
  --color-base-300:    #E3D9C6;   /* --paper-3 */
  --color-base-content: #1F1A15;  /* --ink */

  /* Primary — acento marrom terroso */
  --color-primary:         #4A3527;   /* --accent-deep */
  --color-primary-content: #F5EFE6;  /* --paper */

  /* Secondary / Accent */
  --color-secondary:         #6B4E3D;  /* --accent */
  --color-secondary-content: #F5EFE6;
  --color-accent:            #A67C52;  /* --accent-soft */
  --color-accent-content:    #F5EFE6;

  /* Status */
  --color-success:  #4A5D4F;  /* --moss-2 */
  --color-warning:  #C67B5C;  /* --clay */
  --color-error:    #A0463A;  /* --danger */
  --color-info:     #6B5F52;  /* --ink-3 */

  /* Neutral */
  --color-neutral:         #1F1A15;
  --color-neutral-content: #F5EFE6;

  /* Raios */
  --radius-box:   0.875rem;  /* 14px */
  --radius-btn:   0.5rem;    /* 8px */
  --radius-badge: 999px;
}
```

---

## Mudança 2 — Adicionar classes `.serif` e `.mono` em `web/static/css/style.css`

Adicionar ao final do arquivo (antes do `/* Cache bust */`):

```css
/* ============================================
   TIPOGRAFIA — classes utilitárias Sábio
   ============================================ */
.serif {
  font-family: var(--font-serif, 'Fraunces', Georgia, serif);
}

.mono {
  font-family: var(--font-mono, 'Geist Mono', 'Courier New', monospace);
  font-variant-numeric: tabular-nums;
}
```

---

## Mudança 3 — Sidebar: trocar gradiente verde por `--paper` cream

### 3a. CSS em `web/static/css/input-v2.css`

Substituir o bloco `.arandu-sidebar` (que tem `background: linear-gradient(... #0F6E56 ...)`):

```css
/* ============================================
   ARANDU SIDEBAR — design Sábio
   ============================================ */
.arandu-sidebar {
  background: var(--paper, #F5EFE6);
  border-right: 1px solid var(--line, #D9CDB8);
  color: var(--ink, #1F1A15);
  overflow-y: auto;
  overflow-x: hidden;
}

/* Nav items */
.arandu-sidebar .sabio-nav-item {
  color: var(--ink-2, #3A332B);
  border-radius: 10px;
  transition: background .15s ease, color .15s ease;
}

.arandu-sidebar .sabio-nav-item:hover {
  background: color-mix(in oklab, var(--accent, #6B4E3D) 8%, transparent);
  color: var(--ink, #1F1A15);
}

.arandu-sidebar .sabio-nav-item[aria-current="page"] {
  background: color-mix(in oklab, var(--accent, #6B4E3D) 14%, transparent);
  border: 1px solid color-mix(in oklab, var(--accent, #6B4E3D) 22%, transparent);
  color: var(--accent-deep, #4A3527);
  font-weight: 500;
}

/* Seção labels (MENU / ATALHOS) */
.arandu-sidebar .sabio-nav-section-label {
  color: var(--ink-4, #9B8E7E);
}

/* Brand header */
.arandu-sidebar .sabio-brand-header {
  border-bottom: 1px solid var(--line, #D9CDB8);
}

/* Footer */
.arandu-sidebar .sabio-sidebar-footer {
  border-top: 1px solid var(--line, #D9CDB8);
}

.arandu-sidebar .sabio-collapse-btn {
  border: 1px dashed var(--line-2, #C6B89F);
  color: var(--ink-3, #6B5F52);
}

/* Avatar iniciais */
.arandu-sidebar .sabio-avatar {
  font-family: var(--font-serif, 'Fraunces', serif);
  background: linear-gradient(135deg,
    color-mix(in oklab, var(--accent, #6B4E3D) 80%, var(--paper)),
    var(--accent, #6B4E3D));
  color: var(--paper, #F5EFE6);
}
```

### 3b. BrandMark SVG — atualizar em `shell_layout.templ`

O design usa o glifo "Aperture Áurea" (triângulo proporcional + barra accent). Substituir a
função `sabioBrandMark()` atual pelo SVG exato do design:

```templ
// sabioBrandMark glifo Aperture Áurea — triângulo + contra-forma + barra accent
templ sabioBrandMark() {
    <div class="sabio-brandmark">
        <svg width="30" height="30" viewBox="0 0 96 96" fill="none" aria-hidden="true">
            <path d="M48 6 L86 82 L10 82 Z" fill="currentColor"/>
            <path d="M48 32 L70 76 L26 76 Z" fill="var(--paper)"/>
            <rect x="34" y="60" width="28" height="3" fill="var(--accent)"/>
        </svg>
    </div>
}
```

E o brand header exibe "Arandu" em serif + "CLÍNICO" uppercase. Verificar se o `.sabio-brand-text`
já aplica isso corretamente — se não, ajustar as classes no template.

---

## Mudança 4 — Topbar: `fixed + lg:left-64` → `sticky top-0`

Esta é a mudança mais importante para eliminar definitivamente o "grudado".

O design usa `position: sticky` na topbar (não `fixed`). Com sticky dentro do `drawer-content`,
a topbar fica naturalmente dentro dos limites do conteúdo — sem precisar de `lg:left-64`.

### 4a. Alterar `ShellTopbar()` em `shell_layout.templ`

```templ
// ANTES:
<header class="navbar fixed top-0 left-0 right-0 lg:left-64 z-50 h-16 bg-base-100/95 backdrop-blur-sm border-b border-base-300 px-4 gap-2 shadow-sm">

// DEPOIS:
<header class="sabio-topbar">
```

O CSS de `.sabio-topbar` já existe em `style.css` (linha ~8021) mas precisa ser ajustado
para usar `position: sticky` e `background` Sábio:

### 4b. Atualizar `.sabio-topbar` em `web/static/css/style.css`

```css
.sabio-topbar {
    position: sticky;
    top: 0;
    height: 64px;
    background: color-mix(in oklab, var(--paper) 88%, transparent);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    border-bottom: 1px solid var(--line);
    display: flex;
    align-items: center;
    padding: 0 28px;
    gap: 12px;
    z-index: 20;
    flex-shrink: 0;
}
```

### 4c. Atualizar `drawer-content` em `Shell()` em `shell_layout.templ`

Com `sticky`, o scroll deve estar no `drawer-content`, não no `main`:

```templ
// ANTES:
<div class="drawer-content flex flex-col">
    @ShellTopbar(config)
    <main class="flex-1 mt-16 overflow-y-auto">
        <div class="p-6 md:p-8 max-w-screen-xl mx-auto">
            <div id="main-content" hx-history="false">
                @content
            </div>
        </div>
    </main>
</div>

// DEPOIS (sticky → drawer-content é o scroll container):
<div class="drawer-content flex flex-col h-screen overflow-y-auto">
    @ShellTopbar(config)
    <main class="flex-1">
        <div class="p-6 md:p-8 max-w-screen-xl mx-auto">
            <div id="main-content" hx-history="false">
                @content
            </div>
        </div>
    </main>
</div>
```

Remover `mt-16` do `main` — não é mais necessário com sticky.

---

## Após as edições

```bash
~/go/bin/templ generate ./web/components/...
npm run tailwind:build:v2
go build -o arandu ./cmd/arandu/
# reiniciar o servidor
```

---

## Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros
- [ ] `npm run tailwind:build:v2` sem erros

**Visual — comparar com `design_handoff_arandu_redesign/Arandu Redesign.html`**

- [ ] CA01: Sidebar é creme/quente (--paper), sem gradiente verde
- [ ] CA02: Nav item ativo tem fundo marrom translúcido (`color-mix accent 14%`) e borda marrom sutil
- [ ] CA03: Topbar tem fundo semitransparente quente com blur (não azulado)
- [ ] CA04: BrandMark é o triângulo Aperture Áurea (não a folha antiga)
- [ ] CA05: Topbar NÃO sobrepõe a sidebar no desktop — "grudado" eliminado
- [ ] CA06: Scroll funciona corretamente nas páginas de conteúdo longo (`/session/{id}/edit`)
- [ ] CA07: Dashboard, Pacientes, Sessão — layout correto sem área branca

**Testes Playwright** (se disponíveis)
- [ ] CA08: `bash tests/e2e/shell/test_layout_geometry.sh` passa

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa

---

## NÃO faça

- Não criar novo CSS custom de posicionamento para sidebar/topbar — usar apenas tokens Sábio
- Não remover `#main-content`, `#modal-container`, `#drawer-container`, `#shell-breadcrumb`
- Não alterar o conteúdo interno do `ShellSidebarContent` (nav items, brand text) — apenas o `aside` wrapper e CSS
- Não alterar páginas de conteúdo (dashboard, sessão, paciente) — apenas o chrome/shell
- Não usar `position: fixed` na topbar — usar `sticky`
- Não inventar cores que não estejam nos tokens Sábio do README

---

## Referências

- Design visual: `design_handoff_arandu_redesign/Arandu Redesign.html` (abrir no browser)
- Sidebar design: `design_handoff_arandu_redesign/chrome.jsx` linhas 97–175
- Topbar design: `design_handoff_arandu_redesign/chrome.jsx` linhas 232–310
- BrandMark SVG: `design_handoff_arandu_redesign/chrome.jsx` linhas 183–210
- Tokens Sábio: `design_handoff_arandu_redesign/README.md` seção "Paleta a implementar"
- Implementação atual: `web/components/layout/shell_layout.templ`
