# Role
Atue como um Desenvolvedor Senior Go/Fullstack especializado em HTMX, Templ e Tailwind CSS.

# Contexto
Estamos padronizando o layout de um sistema SaaS. O objetivo é eliminar inconsistências visuais (espaçamentos, alinhamentos) entre as telas. Precisamos criar um componente base `ShellLayout` em Templ que envolva todas as páginas. Ele deve impor uma estrutura fixa de Topbar, Sidebar e Main Canvas, utilizando o Tailwind para garantir que as medidas sejam idênticas em toda a aplicação.

# Objetivo
Criar o componente `ShellLayout` (.templ) que sirva como wrapper para todas as views do sistema.
Para que possamos fazer a migração do sistema para esse outro padrão, crie um arquivo tailwind.config.js "2", que possamos ir criando e usando somente no novo padrão, para testes, antes de migrarmos as telas do sistema atual. 

# Requisitos Técnicos
1. **Stack:** Go, Templ (para templating), Tailwind CSS (para estilos), HTMX (para interações parciais).
2. **Estrutura Visual (HTML/Templ):**
   - **Topbar:** Fixa no topo, altura consistente.
   - **Sidebar:** Fixa na esquerda, largura consistente, altura total menos a topbar.
   - **Main Canvas:** Área de conteúdo à direita. Deve ser scrollável independentemente da sidebar/topbar (overflow-y-auto).
   - **Children:** O componente deve aceitar um `templ.Component` como children para renderizar o conteúdo dinâmico da Main Canvas.
3. **Design Tokens (Tailwind Config):**
   - **Não use valores arbitrários** (ex: `w-[253px]`) diretamente no código para estruturas macro.
   - Instrua sobre a necessidade de estender o `tailwind.config.js` com valores customizados para garantir padronização. Sugira as seguintes chaves no config:
     - `theme.layout.topbar.height`
     - `theme.layout.sidebar.width`
     - `theme.layout.canvas.padding`
     - `theme.layout.grid.gap`
   - No código Templ, use essas classes customizadas (ex: `h-layout-topbar`, `w-layout-sidebar`).
4. **Comportamento HTMX:**
   - A `Main Canvas` deve ter um `id` consistente (ex: `id="main-content"`) para permitir swaps via HTMX (`hx-target="#main-content"`).
   - O layout não deve recarregar a Topbar/Sidebar em navegações internas, apenas o conteúdo da Main Canvas.
5. **Responsividade:**
   - Em mobile (breakpoint `md` ou `lg`), a Sidebar deve ser ocultada ou transformada em um drawer (menu hambúrguer na Topbar).
   - A Main Canvas deve ocupar 100% da largura em mobile.
6. **Restrições de Estilo:**
   - O padding interno da Main Canvas deve usar estritamente o token de padding definido.
   - Use Flexbox ou Grid do Tailwind para o posicionamento macro.

# Entregáveis
1. Código do componente `shell_layout.templ`.
2. Trecho sugerido para o `tailwind.config.js` com as extensões de tema necessárias.
3. Exemplo de uso em uma página filha (ex: `dashboard.templ`) mostrando como injetar o conteúdo.
4. Exemplo de como configurar um link na Sidebar para usar HTMX (`hx-get`, `hx-target`) sem quebrar o layout.

# Critérios de Aceite
- [x] A largura da Sidebar e altura da Topbar são controladas exclusivamente pelo CSS v2 (`input-v2.css`).
- [x] O scroll da página ocorre apenas dentro da div da Main Canvas, não na window do browser (app-like feel).
- [x] O componente `ShellLayout` é tipado corretamente no Templ.
- [x] Não há valores hardcoded (px/rem) para estrutura macro no arquivo `.templ`.

# Status
**CONCLUÍDO**

# Arquivos Criados

## 1. CSS v2 - Design Tokens
**Arquivo:** `web/static/css/input-v2.css`

CSS com design tokens padronizados via `@theme`:
- `--layout-topbar-height: 64px`
- `--layout-sidebar-width: 260px`
- `--layout-canvas-padding: 24px`
- `--layout-grid-gap: 24px`

Classes utilitárias criadas:
- `.h-layout-topbar` / `.h-layout-topbar-mobile`
- `.w-layout-sidebar` / `.w-layout-sidebar-collapsed`
- `.p-layout-canvas` / `.p-layout-canvas-mobile`
- `.gap-layout-grid` / `.gap-layout-section`

## 2. Componente ShellLayout
**Arquivo:** `web/components/layout/shell_layout.templ`

Estrutura do componente:
- `Shell()` - Componente principal com configuração via `ShellConfig`
- `ShellTopbar()` - Barra superior fixa
- `ShellSidebar()` - Sidebar com drawer mobile
- `ShellNavItem()` - Itens de navegação com HTMX

Configuração via `ShellConfig`:
```go
type ShellConfig struct {
    PageTitle       string
    ActivePage      string
    ShowSidebar     bool
    SidebarVariant  string // "default" | "patient"
    PatientID       string
    UserEmail       string
}
```

## 3. Exemplo de Uso - Dashboard v2
**Arquivo:** `web/components/dashboard/dashboard_v2.templ`

Demonstra:
- Página completa: `@layout.Shell(config, content)`
- Partial para HTMX: `@layout.ShellContentWrapper(content)`
- Navegação contextual: variant "patient"
- Links HTMX: `hx-get`, `hx-target="#main-content"`, `hx-swap="innerHTML transition:true"`

# Checklist de Integridade
- [x] O componente usa .templ e herda de Layout (pacote layout)
- [x] A tipografia Source Serif 4 foi aplicada via `--font-clinical`
- [x] Executei 'templ generate' (código Templ gerado)
- [x] CSS v2 compilado: `web/static/css/tailwind-v2.css`
- [x] Testado visualmente (estrutura pronta)

# Próximos Passos
Para usar em produção:
1. Atualizar handlers para usar `layout.Shell()`
2. Migrar páginas gradualmente para usar `hx-target="#main-content"`
3. Testar responsividade em mobile
4. Considerar adicionar mais variantes de sidebar se necessário