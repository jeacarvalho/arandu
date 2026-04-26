# Arandu — Design System: DaisyUI + Tema Próprio

Referência obrigatória para qualquer trabalho visual no Arandu.
Se houver conflito com `tailwind-components`: **esta skill prevalece**.

---

## Stack CSS

```
DaisyUI v4          ← componentes semânticos (btn, card, badge, drawer…)
Tailwind CSS v4     ← utilities inline quando DaisyUI não cobre
style.css           ← APENAS variáveis globais e .clinical — nada mais
```

**Não usar:**
- `input-v2.css` / `tailwind-v2.css` — legacy, será removido
- `@layer components` — não criar regras custom em @layer
- Classes Tailwind concatenadas para simular componentes — usar DaisyUI

---

## Tema Arandu

Definido uma vez em `style.css`. Não alterar sem decisão arquitetural.

```css
[data-theme="arandu"] {
  color-scheme: light;

  --color-primary:         oklch(44% 0.13 162);   /* verde #0F6E56 */
  --color-primary-content: oklch(98% 0 0);

  --color-secondary:         oklch(55% 0.09 220);
  --color-secondary-content: oklch(98% 0 0);

  --color-accent:         oklch(68% 0.16 40);     /* âmbar suave */
  --color-accent-content: oklch(20% 0.04 40);

  --color-neutral:         oklch(28% 0.02 260);
  --color-neutral-content: oklch(96% 0 0);

  --color-base-100: oklch(97.5% 0.006 248);   /* fundo papel #F7F8FA */
  --color-base-200: oklch(94%   0.007 248);
  --color-base-300: oklch(90%   0.007 248);
  --color-base-content: oklch(20% 0.02 260);

  --color-success: oklch(62% 0.17 155);
  --color-warning: oklch(78% 0.18 80);
  --color-error:   oklch(65% 0.22 25);
  --color-info:    oklch(66% 0.14 230);

  --radius-box:   0.75rem;
  --radius-btn:   0.5rem;
  --radius-badge: 0.375rem;
}
```

---

## Tipografia dual — regra de ouro

Duas classes. Nada mais.

```css
/* Em style.css */
body      { font-family: 'Inter', sans-serif; }
.clinical { font-family: 'Source Serif 4', Georgia, serif; }
```

| Onde usar `.clinical` | Onde usar (padrão Inter) |
|-----------------------|--------------------------|
| Títulos de página (`<h1>`) | Labels, botões, nav |
| Conteúdo de prontuário, observações, anamnese | Metadados, timestamps |
| Resumo de sessão | Badges, stats titles |
| Nome do paciente em destaque | Formulários |

```templ
<!-- ✅ Correto -->
<h1 class="clinical text-3xl font-medium">{ patient.Name }</h1>
<p class="clinical text-xl leading-relaxed">{ observation.Content }</p>

<!-- ❌ Errado — inline style desnecessário -->
<h1 style="font-family: 'Source Serif 4'">...</h1>
```

---

## Mapeamento de componentes

### Botões

```templ
<!-- Primário -->
<button class="btn btn-primary">Salvar</button>

<!-- Secundário / ghost -->
<button class="btn btn-ghost">Cancelar</button>

<!-- Destrutivo -->
<button class="btn btn-error btn-outline">Excluir</button>

<!-- Pequeno (em tabelas, cards) -->
<button class="btn btn-primary btn-sm">Confirmar</button>

<!-- Com ícone -->
<button class="btn btn-primary gap-2">
  <i class="fas fa-plus"></i> Nova sessão
</button>
```

### Cards

```templ
<!-- Card padrão -->
<div class="card bg-base-100 border border-base-300 shadow-sm">
  <div class="card-body p-5">
    <h2 class="card-title text-base font-semibold">Título</h2>
    ...
  </div>
</div>

<!-- Card clínico — borda lateral de identidade -->
<div class="card bg-base-100 border border-base-300 border-l-4 border-l-primary shadow-sm">
  <div class="card-body p-5">
    <p class="clinical text-lg leading-relaxed">{ observation.Content }</p>
  </div>
</div>
```

### Badges de status de sessão

```templ
<!-- Mapeamento semântico DaisyUI -->
Agendada  → <div class="badge badge-ghost badge-sm">Agendada</div>
Confirmada → <div class="badge badge-primary badge-sm">Confirmada</div>
Concluída  → <div class="badge badge-success badge-sm">Concluída</div>
Cancelada  → <div class="badge badge-error badge-outline badge-sm">Cancelada</div>
Faltou     → <div class="badge badge-warning badge-sm">Faltou</div>
```

### Formulários / Inputs

```templ
<!-- Input padrão -->
<label class="form-control w-full">
  <div class="label"><span class="label-text text-sm font-medium">Nome</span></div>
  <input type="text" class="input input-bordered bg-base-100 focus:border-primary/50 focus:outline-none" placeholder="...">
</label>

<!-- Textarea clínica -->
<label class="form-control w-full">
  <div class="label"><span class="label-text text-sm font-medium">Observação clínica</span></div>
  <textarea class="textarea textarea-bordered clinical text-base leading-relaxed bg-base-100 focus:border-primary/50 focus:outline-none" rows="4" placeholder="..."></textarea>
</label>

<!-- Select -->
<select class="select select-bordered w-full bg-base-100 focus:border-primary/50 focus:outline-none">
  <option>Opção</option>
</select>
```

### Alertas e notificações

```templ
<!-- Alerta informativo -->
<div class="alert alert-info">
  <i class="fas fa-info-circle"></i>
  <span>Mensagem</span>
</div>

<!-- Alerta de ação pendente -->
<div class="alert bg-warning/10 border border-warning/25">
  <i class="fas fa-exclamation-triangle text-warning"></i>
  <div>
    <div class="font-semibold text-sm">Prontuário em aberto</div>
    <div class="text-xs text-base-content/60">Completar antes da próxima sessão</div>
  </div>
  <button class="btn btn-warning btn-sm">Completar</button>
</div>
```

### Tabelas

```templ
<div class="overflow-x-auto">
  <table class="table table-sm">
    <thead>
      <tr class="text-xs text-base-content/50 font-medium uppercase tracking-wide">
        <th>Paciente</th>
        <th>Data</th>
        <th>Status</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      for _, s := range sessions {
        <tr class="hover:bg-base-200/50">
          <td class="font-medium">{ s.PatientName }</td>
          <td class="text-sm text-base-content/60">{ s.Date }</td>
          <td>@SessionBadge(s.Status)</td>
          <td><button class="btn btn-ghost btn-xs">Ver</button></td>
        </tr>
      }
    </tbody>
  </table>
</div>
```

### Modal

```templ
<!-- Trigger -->
<button class="btn btn-primary" onclick="meu_modal.showModal()">Abrir</button>

<!-- Modal (DaisyUI dialog nativo) -->
<dialog id="meu_modal" class="modal">
  <div class="modal-box">
    <form method="dialog">
      <button class="btn btn-sm btn-circle btn-ghost absolute right-3 top-3">✕</button>
    </form>
    <h3 class="clinical text-xl font-medium mb-4">Título do modal</h3>
    <p class="text-sm text-base-content/70">Conteúdo...</p>
    <div class="modal-action">
      <form method="dialog">
        <button class="btn btn-ghost">Cancelar</button>
        <button class="btn btn-primary">Confirmar</button>
      </form>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button>fechar</button></form>
</dialog>
```

### Stats (dashboard)

```templ
<div class="stat bg-base-100 border border-base-300 rounded-xl p-4">
  <div class="stat-title text-xs text-base-content/50">Este mês</div>
  <div class="stat-value text-2xl font-bold text-primary mt-1">38</div>
  <div class="stat-desc text-xs mt-0.5">sessões realizadas</div>
</div>
```

---

## Shell: estrutura com DaisyUI Drawer

O shell usa `drawer` do DaisyUI — mobile overlay e desktop sidebar sempre-aberta em uma única estrutura.

```templ
templ Shell(config ShellConfig, content templ.Component) {
  <div class="drawer lg:drawer-open" data-theme="arandu">
    <input id="shell-drawer" type="checkbox" class="drawer-toggle">

    <!-- Conteúdo principal -->
    <div class="drawer-content flex flex-col">
      @ShellTopbar(config)
      <main class="flex-1 mt-16 p-5 md:p-8" id="main-content">
        @content
      </main>
    </div>

    <!-- Sidebar -->
    <div class="drawer-side z-40">
      <label for="shell-drawer" class="drawer-overlay"></label>
      <aside class="arandu-sidebar w-64 min-h-screen flex flex-col" id="shell-sidebar">
        @ShellSidebarContent(config)
      </aside>
    </div>
  </div>
}
```

**Topbar alinhada com a sidebar:**
```templ
templ ShellTopbar(config ShellConfig) {
  <header class="navbar fixed top-0 right-0 lg:left-64 left-0 z-50 h-16
                 bg-base-100/90 backdrop-blur-sm border-b border-base-300 px-4 gap-2">
    <!-- Hambúrguer apenas mobile -->
    <label for="shell-drawer" class="btn btn-ghost btn-sm btn-square lg:hidden">
      <!-- ícone hamburger -->
    </label>
    <div class="breadcrumbs text-sm flex-1" id="shell-breadcrumb">
      <!-- breadcrumb items -->
    </div>
    <!-- avatar, notificações -->
  </header>
}
```

**Sidebar CSS (único bloco custom):**
```css
/* Em style.css — único CSS custom necessário para a sidebar */
.arandu-sidebar {
  background: linear-gradient(180deg, #0F6E56 0%, #085041 100%);
}
.arandu-sidebar .menu a       { color: rgba(255,255,255,0.65); border-radius: 0.5rem; }
.arandu-sidebar .menu a:hover { background: rgba(255,255,255,0.10); color: white; }
.arandu-sidebar .menu a.active{ background: rgba(255,255,255,0.18); color: white; font-weight: 500; }
.arandu-sidebar .menu-title   { color: rgba(255,255,255,0.30); font-size: 10px; letter-spacing: .12em; }
```

---

## HTMX com DaisyUI

O padrão HTMX não muda — só o HTML retornado usa DaisyUI.

```go
// Handler: fragmento HTMX
if r.Header.Get("HX-Request") == "true" {
    // Retorna fragmento DaisyUI — sem shell wrapper
    components.MeuComponente(vm).Render(r.Context(), w)
    return
}
// Página completa: usa Shell
layout.Shell(config, components.MeuComponente(vm)).Render(r.Context(), w)
```

**OOB swap de breadcrumb** (padrão inalterado):
```templ
templ ShellBreadcrumb(items []string) {
  <div id="shell-breadcrumb" class="breadcrumbs text-sm flex-1" hx-swap-oob="true">
    <ul>
      for _, item := range items {
        <li>{ item }</li>
      }
    </ul>
  </div>
}
```

---

## Checklist por componente novo

```
[ ] Usou componente DaisyUI semântico? (btn, card, badge, input, table…)
[ ] Conteúdo clínico (observações, prontuário, nome em destaque) usa .clinical?
[ ] Não criou CSS custom — apenas utilities Tailwind quando DaisyUI não cobre?
[ ] Cores via classes semânticas (btn-primary, badge-success) — não hex inline?
[ ] templ generate executado após edição .templ?
```

---

## Referência visual

Protótipo funcional em:
`design_handoff_arandu_redesign/daisyui_shell_dashboard.html`

Abrir no browser para ver o tema Arandu em ação antes de implementar.
