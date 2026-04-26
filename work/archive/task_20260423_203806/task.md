# Task: CA-05 — Migrar Shell para DaisyUI Drawer
Requirement: Interno — qualidade de layout
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Objetivo

Substituir a estrutura CSS manual do shell (`shell-body`, `shell-topbar`, `shell-sidebar`)
pela estrutura nativa **DaisyUI Drawer** (`drawer lg:drawer-open`).

Resultado esperado: sidebar fica full-height separada do topbar (sem o visual "grudado"),
sem nenhum CSS de posicionamento manual. O DaisyUI cuida de tudo.

---

## Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX 2.x · **DaisyUI v5 + Tailwind CSS v4** · Alpine.js 3
**DaisyUI já instalado** (`@plugin "daisyui"` em `web/static/css/input-v2.css`, `data-theme="arandu"` em `<html>`)
**Protótipo visual de referência**: `design_handoff_arandu_redesign/daisyui_shell_dashboard.html`

### IDs que DEVEM ser preservados (HTMX depende deles)

```
#main-content      ← swap target principal
#modal-container   ← modais de agendamento
#drawer-container  ← formulários slide-in
#toast-container   ← notificações
#shell-sidebar     ← OOB swap target da sidebar
#shell-breadcrumb  ← OOB swap target do breadcrumb
```

---

## Arquitetura atual (estado inicial)

```
div.shell  [x-data="Alpine.store('shell')"]
  header.sabio-topbar          ← topbar full-width (left:0) — causa o "grudado"
  div.shell-body
    div.shell-sidebar-overlay
    aside#shell-sidebar.shell-sidebar   ← começa em y=64px
    main.shell-main
      div.shell-canvas
        div.shell-canvas-container
          div#main-content.shell-content

div#modal-container
div#drawer-container
@LLMDrawer
div#toast-container
```

Alpine store (inline no `<head>` do shell_layout.templ):
```javascript
Alpine.store('shell', {
  sidebarOpen: false,
  sidebarCollapsed: false,
  llmOpen: false,
  init() { this.sidebarCollapsed = localStorage.getItem('arandu-sidebar-collapsed') === 'true'; },
  toggleSidebar()  { this.sidebarOpen = !this.sidebarOpen; },
  closeSidebar()   { this.sidebarOpen = false; },
  toggleCollapse() {
    this.sidebarCollapsed = !this.sidebarCollapsed;
    localStorage.setItem('arandu-sidebar-collapsed', this.sidebarCollapsed);
  },
  toggleLLM() { this.llmOpen = !this.llmOpen; },
  openLLM()   { this.llmOpen = true; },
  closeLLM()  { this.llmOpen = false; },
})
```

---

## Arquitetura alvo (DaisyUI Drawer)

```
div.drawer.lg:drawer-open  [x-data="Alpine.store('shell')"]
  input#shell-drawer.drawer-toggle

  div.drawer-content.flex.flex-col
    header.navbar.fixed.top-0.left-0.right-0.lg:left-64.z-50.h-16.bg-base-100...
      ...conteúdo do ShellTopbar (breadcrumb, search, botões)...
    main.flex-1.mt-16.overflow-y-auto
      div.p-6.md:p-8.max-w-screen-xl.mx-auto
        div#main-content

  div.drawer-side.z-40
    label.drawer-overlay[for="shell-drawer"]
    aside#shell-sidebar.arandu-sidebar.w-64.min-h-screen...
      ...conteúdo do ShellSidebar (brand, nav, footer)...

div#modal-container
div#drawer-container
@LLMDrawer
div#toast-container
```

---

## Arquivos a modificar

### 1. `web/components/layout/shell_layout.templ`

**Função `Shell()`** — substituir APENAS a estrutura do `<body>`. O `<head>` (meta, CSS, JS,
Alpine store inline) deve ficar EXATAMENTE igual. Nova estrutura do body:

```templ
<body>
    <div
        class="drawer lg:drawer-open min-h-screen"
        x-data="Alpine.store('shell')"
        x-init="$store.shell.init()"
        @keydown.meta.j.window="$store.shell.toggleLLM()"
        @keydown.ctrl.j.window="$store.shell.toggleLLM()"
    >
        // Checkbox toggle — DaisyUI drawer nativo (mobile)
        <input
            id="shell-drawer"
            type="checkbox"
            class="drawer-toggle"
            :checked="sidebarOpen"
            @change="sidebarOpen = $event.target.checked"
        >

        // Área de conteúdo (topbar + main)
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

        // Sidebar — fora do drawer-content = full-height automático
        <div class="drawer-side z-40">
            <label
                for="shell-drawer"
                aria-label="close sidebar"
                class="drawer-overlay"
                @click="closeSidebar()"
            ></label>
            if config.ShowSidebar {
                @ShellSidebar(config)
            }
        </div>
    </div>

    <div id="modal-container"></div>
    <div id="drawer-container"></div>
    @LLMDrawer()
    <div id="toast-container"></div>
    <script src={ "/static/js/htmx-handlers.js?v=" + helpers.GetJSVersion() }></script>
    <script src="/static/js/htmx-debug.js"></script>
</body>
```

**Função `ShellTopbar()`** — mudar APENAS o elemento outer (`header`):

```templ
// Antes: <header class="sabio-topbar">
// Depois:
<header class="navbar fixed top-0 left-0 right-0 lg:left-64 z-50 h-16
               bg-base-100/95 backdrop-blur-sm border-b border-base-300 px-4 gap-2
               shadow-sm">
    // Hamburger mobile — label DaisyUI (substitui o button anterior)
    <label for="shell-drawer"
           class="btn btn-ghost btn-sm btn-square lg:hidden drawer-button"
           @click="toggleSidebar()">
        <i class="fas fa-bars text-base-content/70"></i>
    </label>
    // ... manter o restante do conteúdo interno IGUAL (breadcrumb, search, botões) ...
</header>
```

**Função `ShellSidebar()`** — mudar APENAS o elemento outer (`aside`):

```templ
// Antes: <aside id="shell-sidebar" class="shell-sidebar" :class="...">
// Depois:
<aside id="shell-sidebar" class="arandu-sidebar w-64 min-h-screen flex flex-col">
    // ... manter TODO o conteúdo interno EXATAMENTE igual ...
    // (sabio-brand-header, sabio-nav, sabio-nav-item, sabio-sidebar-footer, etc.)
</aside>
```

Remover o `:class` Alpine da aside — o collapse pode ser retomado em task separada.

### 2. `web/static/css/input-v2.css`

Adicionar bloco ANTES do comentário `/* Cache bust */` no final do arquivo:

```css
/* ============================================
   ARANDU SIDEBAR — DaisyUI Drawer
   ============================================ */
.arandu-sidebar {
  background: linear-gradient(180deg, var(--arandu-primary, #0F6E56) 0%,
                                       var(--arandu-dark, #085041) 100%);
  color: white;
  overflow-y: auto;
  overflow-x: hidden;
}

.arandu-sidebar .sabio-nav-item {
  color: rgba(255, 255, 255, 0.75);
}

.arandu-sidebar .sabio-nav-item:hover {
  background: rgba(255, 255, 255, 0.12);
  color: white;
}

.arandu-sidebar .sabio-nav-item[aria-current="page"] {
  background: rgba(255, 255, 255, 0.20);
  color: white;
  font-weight: 500;
}

.arandu-sidebar .sabio-nav-section-label,
.arandu-sidebar .sabio-sidebar-title {
  color: rgba(255, 255, 255, 0.40);
}
```

Não apagar nenhum bloco existente. Não tocar em `shell-topbar`, `shell-body`, `shell-sidebar`,
timeline, widget, patient-form — esses podem coexistir sem causar dano.

### 3. Alpine store — atualizar `closeSidebar()`

No bloco `<script>` inline do `<head>` em `shell_layout.templ`, atualizar `closeSidebar`
para também fechar o checkbox DaisyUI:

```javascript
closeSidebar() {
    this.sidebarOpen = false;
    const drawerInput = document.getElementById('shell-drawer');
    if (drawerInput) drawerInput.checked = false;
},
```

---

## Após as edições

```bash
~/go/bin/templ generate ./web/components/...
go build -o arandu ./cmd/arandu/
# reiniciar o servidor
```

---

## Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Visual — testar em `http://localhost:8080/dashboard`**
- [ ] CA01: Sidebar vai do topo ao rodapé com gradiente verde (full-height)
- [ ] CA02: Topbar começa após a sidebar — sem visual "grudado"
- [ ] CA03: Dashboard, Pacientes (`/patients`), Agenda, Perfil de paciente — conteúdo visível
- [ ] CA04: Mobile — hamburger abre/fecha sidebar com overlay (DaisyUI drawer nativo)
- [ ] CA05: `#modal-container` e `#drawer-container` existem no DOM (F12 → Elements)
- [ ] CA06: Sem área branca ou mint em branco na área de conteúdo

**Regressão HTMX**
- [ ] CA07: Navegar via sidebar atualiza conteúdo sem reload de página
- [ ] CA08: Breadcrumb atualiza ao mudar de página (OOB swap)
- [ ] CA09: Modal de agendamento aparece centralizado ao clicar em marcação

**Scripts**
- [ ] `./scripts/arandu_guard.sh` passa
- [ ] `./scripts/arandu_validate_handlers.sh` passa

---

## NÃO faça

- Não reescrever o conteúdo INTERNO do ShellTopbar ou ShellSidebar — só o elemento outer muda
- Não apagar `#main-content`, `#modal-container`, `#drawer-container`, `#shell-breadcrumb`
- Não usar `position: fixed` manualmente no sidebar — DaisyUI drawer cuida disso
- Não tocar em blocos CSS de timeline, widget, patient-form, anamnese
- Não apagar o Alpine store existente — apenas adicionar o `closeSidebar` atualizado
- Não criar CSS custom de posicionamento — se precisar de ajuste, usar DaisyUI utilities

---

## Padrão de referência

- Visual: `design_handoff_arandu_redesign/daisyui_shell_dashboard.html`
- Componente a migrar: `web/components/layout/shell_layout.templ`
- CSS mínimo: `web/static/css/input-v2.css` (apenas `.arandu-sidebar`)
