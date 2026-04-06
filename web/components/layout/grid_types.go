package layout

import (
	"fmt"

	"github.com/a-h/templ"
)

// GridItem representa um item individual dentro da grid
type GridItem struct {
	// Component é o conteúdo Templ a ser renderizado
	Component templ.Component

	// ColSpan define quantas colunas o item ocupa (1-4)
	// Em desktop: 1, 2, 3 ou 4 colunas
	// Em mobile: sempre 1 coluna
	ColSpan int

	// ID identificador único do item (opcional, para referências HTMX)
	ID string

	// Classes adicionais CSS para o item (opcional)
	Classes string
}

// GridRow representa uma linha da grid
// Cada linha contém múltiplos GridItem
// A soma dos ColSpan em uma linha não deve exceder 4
type GridRow struct {
	// Items são os componentes nesta linha
	Items []GridItem

	// ID identificador único da linha (opcional)
	ID string

	// Classes adicionais CSS para a linha (opcional)
	Classes string

	// AlignItems define o alinhamento vertical dos itens na linha (opcional)
	// Valores: "start", "center", "end", "stretch" (padrão: "stretch")
	AlignItems string
}

// GridConfig configuração completa da grid
type GridConfig struct {
	// Rows são as linhas da grid
	Rows []GridRow

	// ID identificador único da grid (opcional, útil para HTMX)
	ID string

	// Classes adicionais CSS para o container (opcional)
	Classes string

	// Gap define o espaçamento entre itens (opcional, padrão: "gap-layout-grid")
	// Valores possíveis: "gap-0", "gap-2", "gap-4", "gap-layout-grid", etc.
	Gap string

	// ResponsiveCols define configuração de colunas responsiva (opcional)
	// Padrão: 4 colunas desktop, 1 coluna mobile
	ResponsiveCols ResponsiveColsConfig
}

// ResponsiveColsConfig configuração de colunas responsivas
type ResponsiveColsConfig struct {
	// Desktop número de colunas em desktop (padrão: 4)
	Desktop int

	// Tablet número de colunas em tablet (padrão: 2)
	Tablet int

	// Mobile número de colunas em mobile (padrão: 1)
	Mobile int
}

// DefaultGridConfig retorna uma configuração padrão para a grid
func DefaultGridConfig() GridConfig {
	return GridConfig{
		ID:  "",
		Gap: "gap-layout-grid",
		ResponsiveCols: ResponsiveColsConfig{
			Desktop: 4,
			Tablet:  2,
			Mobile:  1,
		},
	}
}

// Validate verifica se a configuração da grid é válida
// Retorna erro se a soma dos ColSpan em alguma linha exceder 4
func (cfg GridConfig) Validate() error {
	for rowIdx, row := range cfg.Rows {
		totalSpan := 0
		for itemIdx, item := range row.Items {
			if item.ColSpan < 1 || item.ColSpan > 4 {
				return fmt.Errorf("linha %d, item %d: ColSpan deve ser entre 1 e 4, recebeu %d",
					rowIdx+1, itemIdx+1, item.ColSpan)
			}
			totalSpan += item.ColSpan
		}
		if totalSpan > 4 {
			return fmt.Errorf("linha %d: soma dos ColSpan (%d) excede 4 colunas",
				rowIdx+1, totalSpan)
		}
	}
	return nil
}

// IsEmpty retorna true se a grid não tem linhas
func (cfg GridConfig) IsEmpty() bool {
	return len(cfg.Rows) == 0
}

// HasID retorna true se a grid tem um ID definido
func (cfg GridConfig) HasID() bool {
	return cfg.ID != ""
}

// GetGap retorna o valor do gap ou o padrão
func (cfg GridConfig) GetGap() string {
	if cfg.Gap == "" {
		return "gap-layout-grid"
	}
	return cfg.Gap
}

// GridItemBuilder helper para construir GridItem
type GridItemBuilder struct {
	item GridItem
}

// NewGridItem cria um novo builder para GridItem
func NewGridItem(component templ.Component) *GridItemBuilder {
	return &GridItemBuilder{
		item: GridItem{
			Component: component,
			ColSpan:   1, // padrão: 1 coluna
		},
	}
}

// Span define o número de colunas (1-4)
func (b *GridItemBuilder) Span(cols int) *GridItemBuilder {
	if cols < 1 {
		cols = 1
	} else if cols > 4 {
		cols = 4
	}
	b.item.ColSpan = cols
	return b
}

// WithID define o ID do item
func (b *GridItemBuilder) WithID(id string) *GridItemBuilder {
	b.item.ID = id
	return b
}

// WithClasses adiciona classes CSS
func (b *GridItemBuilder) WithClasses(classes string) *GridItemBuilder {
	b.item.Classes = classes
	return b
}

// Build retorna o GridItem construído
func (b *GridItemBuilder) Build() GridItem {
	return b.item
}

// GridRowBuilder helper para construir GridRow
type GridRowBuilder struct {
	row GridRow
}

// NewGridRow cria um novo builder para GridRow
func NewGridRow() *GridRowBuilder {
	return &GridRowBuilder{
		row: GridRow{
			Items: []GridItem{},
		},
	}
}

// AddItem adiciona um item à linha
func (b *GridRowBuilder) AddItem(item GridItem) *GridRowBuilder {
	b.row.Items = append(b.row.Items, item)
	return b
}

// AddItems adiciona múltiplos itens à linha
func (b *GridRowBuilder) AddItems(items ...GridItem) *GridRowBuilder {
	b.row.Items = append(b.row.Items, items...)
	return b
}

// WithID define o ID da linha
func (b *GridRowBuilder) WithID(id string) *GridRowBuilder {
	b.row.ID = id
	return b
}

// WithClasses adiciona classes CSS à linha
func (b *GridRowBuilder) WithClasses(classes string) *GridRowBuilder {
	b.row.Classes = classes
	return b
}

// Align define o alinhamento vertical
func (b *GridRowBuilder) Align(align string) *GridRowBuilder {
	b.row.AlignItems = align
	return b
}

// Build retorna o GridRow construído
func (b *GridRowBuilder) Build() GridRow {
	return b.row
}

// GridConfigBuilder helper para construir GridConfig
type GridConfigBuilder struct {
	cfg GridConfig
}

// NewGridConfig cria um novo builder para GridConfig
func NewGridConfig() *GridConfigBuilder {
	return &GridConfigBuilder{
		cfg: DefaultGridConfig(),
	}
}

// WithID define o ID da grid
func (b *GridConfigBuilder) WithID(id string) *GridConfigBuilder {
	b.cfg.ID = id
	return b
}

// WithClasses adiciona classes CSS ao container
func (b *GridConfigBuilder) WithClasses(classes string) *GridConfigBuilder {
	b.cfg.Classes = classes
	return b
}

// WithGap define o espaçamento entre itens
func (b *GridConfigBuilder) WithGap(gap string) *GridConfigBuilder {
	b.cfg.Gap = gap
	return b
}

// AddRow adiciona uma linha à grid
func (b *GridConfigBuilder) AddRow(row GridRow) *GridConfigBuilder {
	b.cfg.Rows = append(b.cfg.Rows, row)
	return b
}

// WithResponsiveCols configura colunas responsivas
func (b *GridConfigBuilder) WithResponsiveCols(desktop, tablet, mobile int) *GridConfigBuilder {
	b.cfg.ResponsiveCols = ResponsiveColsConfig{
		Desktop: desktop,
		Tablet:  tablet,
		Mobile:  mobile,
	}
	return b
}

// Build retorna o GridConfig construído
func (b *GridConfigBuilder) Build() GridConfig {
	return b.cfg
}

// MustBuild retorna o GridConfig ou panic em caso de erro de validação
func (b *GridConfigBuilder) MustBuild() GridConfig {
	cfg := b.Build()
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	return cfg
}

// ColSpanClass retorna a classe CSS apropriada para o span em desktop
// Em mobile, todos os itens ocupam 1 coluna (via classe responsiva)
func ColSpanClass(span int) string {
	switch span {
	case 1:
		return "col-span-1 md:col-span-1"
	case 2:
		return "col-span-1 md:col-span-2"
	case 3:
		return "col-span-1 md:col-span-3"
	case 4:
		return "col-span-1 md:col-span-4"
	default:
		return "col-span-1"
	}
}

// AlignClass retorna a classe CSS para alinhamento vertical
func AlignClass(align string) string {
	switch align {
	case "start":
		return "items-start"
	case "center":
		return "items-center"
	case "end":
		return "items-end"
	case "stretch", "":
		return "items-stretch"
	default:
		return "items-stretch"
	}
}

// GetGridColsClass retorna a classe CSS para o número de colunas responsivas
func GetGridColsClass(cfg GridConfig) string {
	cols := cfg.ResponsiveCols

	// Se não configurado, usa padrão 4/1
	if cols.Desktop == 0 && cols.Mobile == 0 {
		return "grid-cols-1 md:grid-cols-4"
	}

	// Monta a classe responsiva
	var classes string

	// Mobile (padrão: 1)
	if cols.Mobile == 0 {
		classes = "grid-cols-1"
	} else {
		classes = fmt.Sprintf("grid-cols-%d", cols.Mobile)
	}

	// Tablet (se configurado)
	if cols.Tablet > 0 {
		classes += fmt.Sprintf(" sm:grid-cols-%d", cols.Tablet)
	}

	// Desktop
	if cols.Desktop > 0 {
		classes += fmt.Sprintf(" md:grid-cols-%d", cols.Desktop)
	}

	return classes
}
