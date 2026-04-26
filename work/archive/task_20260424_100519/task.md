# Task: CA-07 — Fix shell layout — remover regra global `header {}` e corrigir colisões
Requirement: Interno — qualidade de layout
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Objetivo

Corrigir dois sintomas visuais que têm a mesma causa raiz:

1. **Topbar "grudada"**: a topbar começa em `x=0` mesmo com `lg:left-64` no elemento — porque `style.css` tem `header { left: 0 !important }` que derruba o Tailwind
2. **Texto flutuante no topo**: `<header class="sabio-notes-col-header">` em `edit_form.templ` é fixado em `(0,0)` pela mesma regra global

---

## Causa raiz (diagnóstico preciso)

**`web/static/css/style.css` — linhas 1931–1937:**

```css
header {
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    right: 0 !important;
    z-index: 100 !important;
}
```

Esta regra aplica `position: fixed; left: 0` a **todo** elemento `<header>` no DOM, com `!important` que derruba qualquer utilidade Tailwind. Ela afeta:

- `header.navbar` em `ShellTopbar()` → `lg:left-64` do Tailwind é ignorado, topbar fica em `left:0`
- `<header class="sabio-notes-col-header">` em `edit_form.templ:124` → fixado em `(0,0)`, aparece sobre a topbar/sidebar como texto flutuante mostrando "AÇÃO / Intervenções terapêuticas / Técnicas e intervenções realizadas"

**Ordem de carregamento de CSS** (confirma a hierarquia):
```html
<link href="/static/css/tailwind-v2.css">  <!-- @layer utilities — mais fraco -->
<link href="/static/css/style.css">         <!-- unlayered + !important — mais forte -->
```

---

## Arquivos a modificar

### 1. `web/static/css/style.css` — remover o bloco global `header {}`

**Remover completamente** as linhas 1931–1937:

```css
/* REMOVER ESTE BLOCO INTEIRO: */
header {
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    right: 0 !important;
    z-index: 100 !important;
}
```

Não substituir por nada. O `header.navbar` em `ShellTopbar()` já tem todas as classes Tailwind necessárias (`fixed top-0 left-0 right-0 lg:left-64 z-50 h-16`). A regra global não serve mais para nada com o novo shell DaisyUI.

---

### 2. `web/components/session/edit_form.templ` — trocar `<header>` por `<div>`

**Linha 124** — trocar o elemento semântico `<header>` pelo `<div>` para evitar colisão com qualquer regra CSS futura que afete o elemento `header`:

```templ
<!-- ANTES: -->
<header class="sabio-notes-col-header">
    <div class="sabio-card-eyebrow">Ação</div>
    <h2 class="sabio-card-title serif">Intervenções terapêuticas</h2>
    <p class="sabio-notes-col-sub">Técnicas e intervenções realizadas</p>
</header>

<!-- DEPOIS: -->
<div class="sabio-notes-col-header">
    <div class="sabio-card-eyebrow">Ação</div>
    <h2 class="sabio-card-title serif">Intervenções terapêuticas</h2>
    <p class="sabio-notes-col-sub">Técnicas e intervenções realizadas</p>
</div>
```

---

### 3. Verificar outros `<header>` no codebase

Buscar todos os `<header>` em `web/components/` que **não sejam** `header.navbar` do ShellTopbar:

```bash
grep -rn "<header " web/components/ | grep -v "_templ.go" | grep -v "shell_layout"
```

Para cada resultado: se o `<header>` é um container interno de componente (não a topbar do shell), trocar por `<div>` mantendo as classes CSS.

---

## Após as edições

```bash
~/go/bin/templ generate ./web/components/...
go build -o arandu ./cmd/arandu/
npm run tailwind:build:v2
# reiniciar o servidor
```

---

## Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Visual — testar em `http://localhost:8080/session/s0014-146/edit`**
- [ ] CA01: O texto "Ação / Intervenções terapêuticas / Técnicas e intervenções realizadas" NÃO aparece no topo da tela — está dentro do card da coluna direita
- [ ] CA02: A topbar começa APÓS a sidebar (não grudada) — visível no desktop 1280px+
- [ ] CA03: A sidebar vai do topo ao rodapé com gradiente verde

**Visual — testar em `http://localhost:8080/dashboard`**
- [ ] CA04: Layout correto (sidebar separada da topbar)

**Testes Playwright** (se o servidor estiver rodando)
- [ ] CA05: `bash tests/e2e/shell/test_layout_geometry.sh` passa — todos os 6 critérios geométricos

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa
- [ ] `./scripts/arandu_validate_handlers.sh` passa

---

## NÃO faça

- Não reescrever o conteúdo interno do `ShellTopbar` ou `ShellSidebar`
- Não adicionar `left: 256px` ou qualquer offset manual em CSS — o Tailwind `lg:left-64` já faz isso
- Não apagar `#main-content`, `#modal-container`, `#drawer-container`
- Não tocar em outros blocos do `style.css` (timeline, widget, patient-form, anamnese)
- Não tocar em `input-v2.css` — apenas `style.css` e `edit_form.templ`

---

## Padrão de referência

- `web/components/layout/shell_layout.templ` — `ShellTopbar()` usa `header.navbar fixed top-0 left-0 right-0 lg:left-64 z-50 h-16` (correto, não mudar)
- `web/components/session/edit_form.templ:124` — trocar `<header>` por `<div>`
- `web/static/css/style.css:1931` — remover o bloco `header { position: fixed !important; ... }`
