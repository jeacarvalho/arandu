# Task: Corrigir Layout Shell — Sobreposição Topbar/Sidebar/Conteúdo
Requirement: —
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Contexto e sintomas visuais

O shell do Arandu apresenta sobreposição visual entre três camadas:
1. **Topbar fixa** (`.sabio-topbar`, `position: fixed; top: 0; z-index: 50`)
2. **Sidebar de navegação** (`aside.shell-sidebar`, desktop: `position: relative` dentro de `.shell-body`)
3. **Sidebar contextual de paciente** (`sidebar_patient.templ`, renderizado dentro do conteúdo)

**Sintomas confirmados em produção (screenshots):**

| Página | Sintoma |
|--------|---------|
| Dashboard | Breadcrumb exibe texto de página anterior ("Pacientes › Novo Paciente") em vez de limpar na navegação HTMX |
| Perfil do paciente | Texto da sidebar lateral invade a área da topbar no topo esquerdo |
| Edição de sessão (`/session/{id}/edit`) | Múltiplos textos sobrepostos no canto superior esquerdo: navegação de paciente, breadcrumb e itens de sidebar empilhados sobre a topbar |

---

## Estrutura atual do shell (NÃO alterar a arquitetura — apenas corrigir CSS/lógica)

```
web/components/layout/
├── shell_layout.templ          ← layout principal
├── sidebar_patient.templ       ← sidebar contextual de paciente
├── page_header.templ           ← header de página (breadcrumb?)
├── standard_grid.templ
└── ...
```

**Variáveis CSS em vigor** (não alterar valores, apenas garantir que são aplicadas):
```css
--layout-topbar-height: 64px;
--layout-sidebar-width: 232px;
--layout-sidebar-width-collapsed: 68px;
--z-layout-topbar: 50;
--z-layout-sidebar: 40;
```

**Estrutura do shell_layout.templ:**
```templ
<div class="shell">
  @ShellTopbar(config)           ← header com class="sabio-topbar"
  <div class="shell-body">       ← margin-top: 64px
    <aside id="shell-sidebar">   ← sidebar principal, desktop: relative
      @ShellSidebarContent(config)
    </aside>
    <main class="shell-main">
      <div class="shell-canvas">
        <div class="shell-content" id="main-content">
          @content               ← conteúdo da página via HTMX
        </div>
      </div>
    </main>
  </div>
</div>
```

---

## O que investigar e corrigir

### Bug 1 — Breadcrumb não atualiza na navegação HTMX

**Causa provável:** O breadcrumb da topbar (`sabio-topbar`) não tem um `id` alvo para
HTMX out-of-band swap. Quando HTMX troca `#main-content`, o breadcrumb permanece com
o valor da requisição anterior.

**Solução:**
1. Adicionar `id="shell-breadcrumb"` ao elemento de breadcrumb na topbar (se ainda não existe)
2. Cada handler que serve fragmento HTMX deve incluir o breadcrumb atualizado via
   `hx-swap-oob="true"` — OU os handlers devem emitir `HX-Trigger` para limpar o breadcrumb
3. Verificar qual o mecanismo já existe e completar/corrigir — não duplicar mecanismos

### Bug 2 — Sidebar de paciente invade área da topbar

**Causa provável:** A `sidebar_patient.templ` é posicionada de forma que seu conteúdo
começa no `y: 0` do viewport, sem `padding-top` para compensar a topbar fixa.

Investigar:
- Se `sidebar_patient.templ` usa `position: fixed` ou `position: sticky` com `top: 0`
- Se há `overflow: visible` no `.shell-sidebar` permitindo conteúdo vazar para cima
- Se o template é renderizado dentro do `#main-content` (incorreto) em vez de ser
  injetado via out-of-band swap no `#shell-sidebar`

**Solução esperada:**
- A sidebar de paciente deve ser renderizada **dentro do `#shell-sidebar`** via HTMX OOB
  (`hx-swap-oob="innerHTML:#shell-sidebar"`) nas páginas de paciente
- Se já é renderizada no sidebar correto: adicionar `padding-top: var(--layout-topbar-height)`
  na sidebar quando em modo `position: fixed` (mobile)
- Garantir que `.shell-sidebar` tenha `overflow: hidden` no eixo Y superior

### Bug 3 — Sobreposição severa na página de edição de sessão

**Causa provável:** Combinação dos dois bugs acima. Adicionalmente:
- A rota `/session/{id}/edit` pode estar retornando o shell completo em vez de fragmento
  (duplicando sidebar na resposta HTMX), OU
- O `sidebar_patient.templ` está sendo incluído dentro do `@content` E no sidebar
  simultaneamente

**Investigar no handler `session_handler.go`:**
```go
// Verificar se em requests HTMX retorna fragmento sem shell
if r.Header.Get("HX-Request") == "true" {
    // deve retornar APENAS o conteúdo, sem shell wrapper
    components.SessionEdit(vm).Render(r.Context(), w)
    return
}
// Retorna shell completo apenas em navegação direta
layout.Shell(config, components.SessionEdit(vm)).Render(r.Context(), w)
```

---

## Checklist de implementação

- [ ] **Diagnóstico**: mapear como cada rota afetada decide retornar fragmento vs shell completo (verificar `HX-Request` em session_handler, patient_handler)
- [ ] **Diagnóstico**: verificar se `sidebar_patient.templ` está dentro do `@content` ou via OOB swap em `#shell-sidebar`
- [ ] **Fix breadcrumb**: implementar atualização do breadcrumb no HTMX (OOB ou evento HTMX)
- [ ] **Fix sidebar mobile**: garantir `padding-top: 64px` quando sidebar é `position: fixed`
- [ ] **Fix sidebar paciente**: mover para OOB swap em `#shell-sidebar` se estiver renderizada dentro de `#main-content`
- [ ] **Fix z-index**: garantir topbar sempre visível sobre qualquer conteúdo da sidebar
- [ ] `~/go/bin/templ generate ./web/components/...` após qualquer edição em `.templ`
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

---

## Arquivos a verificar/modificar

| Arquivo | O que verificar |
|---------|----------------|
| `web/components/layout/shell_layout.templ` | Estrutura do breadcrumb, IDs dos elementos |
| `web/components/layout/sidebar_patient.templ` | Posicionamento, se é OOB ou inline no content |
| `internal/web/handlers/session_handler.go` | Se retorna shell completo em HTMX request |
| `internal/web/handlers/patient_handler.go` | Se retorna shell completo em HTMX request |
| `web/static/css/` (arquivo relevante) | Regras de `padding-top` da sidebar em mobile e `overflow` |

---

## Skills de referência obrigatória

- **arandu-architecture** — padrão HTMX fragmento vs página completa (verificar `HX-Request`), containers do shell (`#modal-container`, `#drawer-container`, `#main-content`, `#shell-sidebar`), o que é e não é retornado em requests HTMX
- **go-templ-htmx-ux** — HTMX out-of-band swap (`hx-swap-oob`), como atualizar elementos fora do target principal (breadcrumb, sidebar), eventos HTMX (`htmx:afterSettle`)
- **tailwind-components** — z-index tokens, `position: fixed` com `top` correto, `overflow` management

---

## Critérios de aceite

**Compilação**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build -o arandu ./cmd/arandu/` sem erros

**Comportamento visual (validar manualmente percorrendo o fluxo completo)**

O coder deve abrir o browser em `localhost:8080` e percorrer o fluxo completo **antes de declarar conclusão**:

1. Login → Dashboard: breadcrumb mostra "Dashboard" (não texto de página anterior)
2. Clicar em "Pacientes": breadcrumb atualiza para "Pacientes"
3. Clicar em um paciente: sidebar lateral muda para nav do paciente (Resumo, Anamnese, Prontuário...) **sem invadir a topbar**
4. Clicar em uma sessão: conteúdo da sessão aparece no painel principal **sem sobreposição de texto**
5. Clicar em "Editar sessão": página de edição abre **sem textos empilhados no canto superior esquerdo**
6. Navegar de volta pelo breadcrumb: breadcrumb atualiza corretamente

**Critérios objetivos:**
- [ ] CA01: Nenhum texto de sidebar aparece sobre a topbar em nenhuma página
- [ ] CA02: Breadcrumb na topbar reflete sempre a página atual após navegação HTMX
- [ ] CA03: Sidebar de paciente (Resumo, Anamnese, Prontuário...) aparece apenas no painel lateral esquerdo
- [ ] CA04: Em desktop (1280px), tablet (768px) e mobile (375px): sem sobreposição visual
- [ ] CA05: Ao navegar via HTMX, nenhum artefato visual de renderizações anteriores persiste

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa sem erros
- [ ] `./scripts/arandu_validate_handlers.sh` passa
- [ ] `go test ./...` sem regressões

---

## NÃO faça

- Não alterar as variáveis CSS (`--layout-topbar-height`, `--layout-sidebar-width`, etc.)
- Não refatorar a arquitetura do shell — apenas corrigir os bugs de posicionamento/swap
- Não introduzir novos IDs sem verificar conflito com os existentes (`#modal-container`, `#drawer-container`, `#main-content`, `#shell-sidebar`)
- Não usar `!important` desnecessário — os bugs são de estrutura, não de especificidade CSS
- Não declarar a task concluída sem ter verificado manualmente os 6 passos do fluxo completo no browser

---

## Padrão de referência

Siga o padrão de verificação de `HX-Request` em `internal/web/handlers/agenda_handler.go`
para garantir que handlers retornam fragmento (sem shell) em requests HTMX.

Para HTMX OOB swap de breadcrumb, o padrão é incluir no fragmento retornado:
```templ
<nav id="shell-breadcrumb" hx-swap-oob="true">
    <span>Pacientes</span> › <span>Adriana Barbosa</span>
</nav>
```
Verificar se esse mecanismo já existe no projeto antes de implementar do zero.
