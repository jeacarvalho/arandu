# Role
Atue como um Desenvolvedor Senior Go/Fullstack especializado em HTMX, Templ e Tailwind CSS.

# Contexto
Já possuímos o componente `ShellLayout` (Topbar + Sidebar + Main Canvas). Agora precisamos padronizar o conteúdo interno da **Main Canvas**. O objetivo é criar um componente `StandardGrid` que force um layout de 4 colunas rigoroso. O layout deve ser **declarativo**: o desenvolvedor não monta o HTML da grid manualmente, ele preenche uma Struct em Go e o Templ renderiza a grid correta.

# Objetivo
Criar o componente `StandardGrid` (.templ) e as Structs Go associadas para gerenciar o layout de 4 colunas dentro da Main Canvas.

# Requisitos Técnicos
1. **Stack:** Go (Structs), Templ (Renderização), Tailwind CSS (Estilos), HTMX (Compatibilidade).
2. **Estrutura de Dados (Go Structs):**
   - Crie structs que permitam definir o layout via código Go. Sugestão:
     - `GridItem`: Deve conter `Component` (templ.Component), `ColSpan` (int 1-4), e `ID` (string).
     - `GridRow`: Deve conter uma slice de `[]GridItem`.
     - `GridConfig`: Deve conter uma slice de `[]GridRow`.
   - O `ColSpan` define quantas das 4 colunas o item ocupa.
3. **Componente Templ (`standard_grid.templ`):**
   - Deve receber `GridConfig` como parâmetro.
   - Deve renderizar um container `div` com `grid grid-cols-4` (do Tailwind).
   - Deve iterar sobre as Rows e Items, aplicando classes dinâmicas de `col-span-{n}` baseadas no `ColSpan` da struct.
   - **Importante:** Não use margins nos items individuais. O espaçamento deve ser gerido pelo `gap` do container grid.
4. **Design Tokens (Tailwind):**
   - Utilize estritamente as variáveis definidas no `tailwind.config.js` (criadas no passo anterior do ShellLayout):
     - `gap-layout-grid` para o espaçamento entre itens.
     - `padding-layout-canvas` para o padding externo do container grid (se não estiver herdado do Shell).
   - Use classes utilitárias do Tailwind para spans: `col-span-1`, `col-span-2`, `col-span-3`, `col-span-4`.
5. **Responsividade:**
   - Em desktop: Respeite o `ColSpan` (4 colunas).
   - Em mobile (breakpoint `md`): A grid deve colapsar para 1 coluna (`grid-cols-1`). Todos os itens devem ocupar 100% da largura (`col-span-1` forçado via media query ou classe responsiva `md:col-span-{n}`).
6. **Integração HTMX:**
   - Os `GridItem` podem conter componentes com diretivas HTMX (`hx-get`, `hx-target`).
   - O container da grid não deve interferir nos swaps do HTMX.
   - O container deve ter um ID opcional para permitir refresh parcial da grid inteira se necessário.
7. **Restrições:**
   - Proibido usar CSS inline ou valores arbitrários (ex: `w-[30%]`).
   - A soma dos `ColSpan` em uma linha não deve exceder 4 (validação lógica ou visual).
   - A altura dos itens é automática (content-based), não force alturas fixas.

# Entregáveis
1. Definição das **Structs Go** (`GridConfig`, `GridRow`, `GridItem`).
2. Código do componente `standard_grid.templ`.
3. Atualização sugerida no `tailwind.config.js` (caso falte alguma classe utilitária para col-span).
4. Exemplo de uso prático: Um arquivo `dashboard_page.templ` que usa o `ShellLayout` + `StandardGrid` com 3 linhas de exemplos (ex: 4 itens de 1 col, 2 itens de 2 col, 1 item de 4 col).

# Critérios de Aceite
- [x] O layout é definido 100% pelas Structs Go, não pelo HTML manual.
- [x] O espaçamento entre os componentes é idêntico em toda a grid (controlado pelo `gap`).
- [x] Em mobile, a grid se transforma em uma coluna única verticalmente.
- [x] Não há quebra de layout se um componente tiver conteúdo maior (a altura da linha se ajusta).
- [x] O código está tipado e seguro (Templ).

# Status
**CONCLUÍDO**

# Arquivos Criados

## 1. Structs Go
**Arquivo:** `web/components/layout/grid_types.go`

Structs definidas:
- `GridItem` - Componente templ, ColSpan (1-4), ID, Classes
- `GridRow` - Slice de GridItem, ID, Classes, AlignItems
- `GridConfig` - Slice de GridRow, ID, Classes, Gap, ResponsiveCols
- Builders: `NewGridConfig()`, `NewGridRow()`, `NewGridItem()`
- Helpers: `ColSpanClass()`, `AlignClass()`, `GetGridColsClass()`
- Validação: `Validate()` - verifica se soma dos spans <= 4

## 2. Componente StandardGrid
**Arquivo:** `web/components/layout/standard_grid.templ`

Componentes:
- `StandardGrid(cfg GridConfig)` - Grid principal
- `StandardGridWithValidation(cfg GridConfig)` - Com validação
- `StandardGridRow(row GridRow, cfg GridConfig)` - Linha individual
- `GridContainer(id, classes string)` - Container genérico
- `GridCell(span int, classes string)` - Célula individual

## 3. Exemplo de Uso
**Arquivo:** `web/components/dashboard/dashboard_grid_example.templ`

Exemplo DashboardWithGrid demonstra:
- Linha 1: 4 StatCards (1 coluna cada)
- Linha 2: 2 cards principais (2 colunas cada)
- Linha 3: 3 ações rápidas (1+1+2 colunas)
- Linha 4: Lembretes (4 colunas)

Uso do builder pattern:
```go
gridConfig := layout.NewGridConfig().
    WithID("dashboard-grid").
    AddRow(layout.NewGridRow().
        AddItem(layout.NewGridItem(StatCard(...)).Span(1).Build()).
        Build()).
    Build()
```

## 4. Design Tokens
Classes Tailwind usadas (já existem no CSS v2):
- `.gap-layout-grid` - Espaçamento entre itens
- `grid-cols-1 md:grid-cols-4` - Grid responsiva
- `col-span-1 md:col-span-{n}` - Colunas responsivas

# Checklist de Integridade
- [x] O componente usa .templ e herda de Layout (pacote layout)
- [x] A tipografia Source Serif 4 foi aplicada (via CSS v2)
- [x] Executei 'templ generate' e o código Go compilou
- [x] Não há valores hardcoded (usa CSS variables e Tailwind)
- [x] Integração HTMX: hx-target="#main-content", hx-swap="innerHTML transition:true"

# Uso Recomendado
```go
// No handler
cfg := layout.NewGridConfig().
    AddRow(layout.NewGridRow().
        AddItem(layout.NewGridItem(componente1).Span(2).Build()).
        AddItem(layout.NewGridItem(componente2).Span(2).Build()).
        Build()).
    Build()

// No templ
@layout.StandardGrid(cfg)
```