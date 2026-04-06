# Role
Atue como um Desenvolvedor Senior Go/Fullstack especializado em HTMX, Templ e Tailwind CSS.

# Contexto
Já possuímos o `ShellLayout` (estrutura macro) e o `StandardGrid` (grid de 4 colunas). Agora precisamos padronizar o conteúdo *dentro* de cada célula da grid. Atualmente, os agentes criam componentes com paddings internos variados, bordas inconsistentes e estilos de fundo diferentes. Precisamos de um `WidgetWrapper` universal que envolva todo conteúdo funcional (Cards, Tabelas, Gráficos) dentro da grid.

# Objetivo
Criar o componente `WidgetWrapper` (.templ) que forneça um container padronizado com visual de "cartão" (card), padding interno consistente e suporte a cabeçalho opcional com ações.

# Requisitos Técnicos
1. **Stack:** Go (Structs), Templ (Renderização), Tailwind CSS (Estilos), HTMX (Atributos).
2. **Estrutura de Dados (Go Structs):**
   - Crie uma struct `WidgetProps` contendo:
     - `Title` (string): Opcional. Título do widget.
     - `Subtitle` (string): Opcional. Descrição curta.
     - `Actions` (templ.Component): Opcional. Slot para botões/ações no header.
     - `Content` (templ.Component): Obrigatório. O conteúdo principal do widget.
     - `HTMXAttributes` (templ.Attributes): Para permitir passar `hx-get`, `hx-target`, etc., para a raiz do wrapper.
     - `NoPadding` (bool): Opcional. Se true, remove o padding interno (para casos especiais como mapas full-bleed).
3. **Componente Templ (`widget_wrapper.templ`):**
   - Deve renderizar uma `div` container com estilos de "Card" (background, border, shadow, radius).
   - **Padding Interno:** Deve usar estritamente o token de design (ex: `p-layout-widget-padding` ou `p-6` padronizado no config).
   - **Header:** Se `Title` estiver presente, renderizar um header consistente (tamanho de fonte, cor, peso) separado do conteúdo por uma borda ou espaçamento padrão.
   - **Content:** Renderizar o `Content` dentro do container.
   - **HTMX:** Aplicar os `HTMXAttributes` na div raiz do wrapper para permitir que o widget inteiro seja atualizado via HTMX.
4. **Design Tokens (Tailwind):**
   - Utilize as variáveis do `tailwind.config.js` (definidas nos passos anteriores):
     - `bg-widget-background` (ex: white ou gray-50).
     - `border-widget-border` (ex: gray-200).
     - `shadow-widget-shadow` (ex: sm ou md).
     - `rounded-widget-radius` (ex: md ou lg).
     - `p-widget-padding` (ex: 6 ou 24px).
   - **Importante:** O wrapper deve ter `box-border` para garantir que o padding não aumente a largura definida pela grid.
5. **Responsividade:**
   - Em mobile, o padding pode ser reduzido ligeiramente para aproveitar o espaço, mas deve permanecer consistente entre todos os widgets.
6. **Restrições de Estilo:**
   - **Proibido** definir margins externas no `WidgetWrapper`. O espaçamento externo é controlado exclusivamente pelo `gap` do `StandardGrid`.
   - **Proibido** usar valores hardcoded (px/rem) para padding/border/shadow. Use as classes do Tailwind configuradas.
   - O wrapper deve ocupar 100% da largura e altura disponível na célula da grid (`w-full h-full`).

# Entregáveis
1. Definição da **Struct Go** `WidgetProps`.
2. Código do componente `widget_wrapper.templ`.
3. Atualização sugerida no `tailwind.config.js` com os tokens de widget (bg, border, shadow, padding).
4. Exemplo de uso prático: Um arquivo `kpi_card.templ` que usa o `WidgetWrapper` para envolver um conteúdo simples (ex: um número grande e um label).

# Critérios de Aceite
- [x] Todo widget envolvido por este componente tem exatamente o mesmo padding interno (24px desktop, 16px mobile).
- [x] Todo widget tem o mesmo estilo de borda, sombra e fundo (border-neutral-200, shadow-sm, bg-white).
- [x] O header (título) é visualmente idêntico em todos os widgets (text-lg font-semibold text-neutral-800).
- [x] É possível passar atributos HTMX para o wrapper (ex: refresh automático).
- [x] Não há margins externas no wrapper (não quebra o gap da grid).
- [x] O conteúdo não toca as bordas do widget (graças ao padding interno).

# Status
**CONCLUÍDO**

# Arquivos Criados

## 1. Struct Go WidgetProps
**Arquivo:** `web/components/layout/widget_types.go`

Propriedades:
- `Title` (string) - Título opcional
- `Subtitle` (string) - Subtítulo opcional
- `Actions` (templ.Component) - Ações no header
- `Content` (templ.Component) - Conteúdo principal (obrigatório)
- `HTMXAttributes` (templ.Attributes) - Atributos HTMX
- `NoPadding` (bool) - Remove padding (full-bleed)
- `Compact` (bool) - Padding reduzido
- `Hover` (bool) - Efeito hover
- `ID` (string) - ID do elemento
- `Classes` (string) - Classes CSS adicionais
- `Footer` (templ.Component) - Rodapé opcional

Builders:
- `NewWidgetProps(content)` - Cria props com builder pattern
- Métodos: `WithTitle()`, `WithSubtitle()`, `WithActions()`, `WithHTMX()`, `WithNoPadding()`, `WithCompact()`, `WithHover()`, `WithID()`, `WithClasses()`, `WithFooter()`

## 2. Componente WidgetWrapper
**Arquivo:** `web/components/layout/widget_wrapper.templ`

Componentes:
- `WidgetWrapper(props)` - Componente principal
- `Widget(title, subtitle, content)` - Helper simplificado
- `WidgetWithActions(title, subtitle, actions, content)` - Com ações
- `WidgetCompact(title, content)` - Widget compacto
- `WidgetNoPadding(title, content)` - Sem padding (full-bleed)
- `WidgetWithHTMX(title, content, attrs)` - Com HTMX
- `WidgetLoading()` - Estado de loading
- `WidgetError(message, retryURL)` - Estado de erro
- `WidgetEmpty(icon, title, description)` - Estado vazio

## 3. Design Tokens (CSS v2)
**Arquivo:** `web/static/css/input-v2.css`

Tokens adicionados:
```css
--widget-bg: var(--color-arandu-paper)
--widget-bg-elevated: #FFFFFF
--widget-border-color: var(--color-neutral-200)
--widget-shadow: var(--shadow-sm)
--widget-shadow-hover: var(--shadow-md)
--widget-radius: 0.75rem (rounded-xl)
--widget-padding: 24px
--widget-padding-mobile: 16px
--widget-padding-compact: 16px
--widget-title-size: 1.125rem
--widget-subtitle-size: 0.875rem
```

Classes CSS:
- `.widget-wrapper` - Container base
- `.widget-wrapper-no-padding` - Sem padding
- `.widget-wrapper-compact` - Padding reduzido
- `.widget-wrapper-hover` - Efeito hover
- `.widget-header` - Header com borda
- `.widget-title` - Estilo do título
- `.widget-subtitle` - Estilo do subtítulo
- `.widget-content` - Área de conteúdo

## 4. Exemplos de Uso
**Arquivo:** `web/components/dashboard/kpi_widget_examples.templ`

Exemplos criados:
- `KPIWidgetCard` - Card KPI com valor, change e trend
- `WidgetWithActionsExample` - Widget com botões no header
- `WidgetWithAutoRefresh` - Widget com HTMX auto-refresh
- `WidgetCompactExample` - Widget com padding reduzido
- `WidgetNoPaddingExample` - Widget full-bleed (mapa)
- `DashboardWithWidgets` - Dashboard completo integrando tudo
- `WidgetLoadingExample` - Estado de loading
- `WidgetErrorExample` - Estado de erro
- `WidgetEmptyExample` - Estado vazio

## 5. Build
**Arquivo:** `web/static/css/tailwind-v2.css` (52KB)

## 6. Arquivos Gerados
- `widget_types.go` (7.5KB)
- `widget_wrapper_templ.go` (25KB)
- `kpi_widget_examples_templ.go` (85KB)

# Checklist de Integridade
- [x] O componente usa .templ e herda de Layout (pacote layout)
- [x] A tipografia Source Serif 4 foi aplicada (via CSS v2)
- [x] Executei 'templ generate' e o código Go compilou
- [x] Build CSS v2 atualizado: `web/static/css/tailwind-v2.css`
- [x] Sem margins externas no wrapper (usa `box-border`)
- [x] Integração HTMX: `hx-target="#main-content"`, `hx-swap="innerHTML transition:true"`

# Uso Recomendado
```go
// Widget simples
@layout.Widget("Título", "Subtítulo", conteúdo)

// Widget com builder pattern
props := layout.NewWidgetProps(conteúdo).
    WithTitle("Título").
    WithSubtitle("Descrição").
    WithActions(ações).
    WithHover().
    Build()
@layout.WidgetWrapper(props)

// Widget com HTMX
attrs := templ.Attributes{
    "hx-get": "/api/data",
    "hx-trigger": "every 30s",
}
@layout.WidgetWithHTMX("Estatísticas", conteúdo, attrs)
```