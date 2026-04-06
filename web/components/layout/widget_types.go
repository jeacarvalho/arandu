package layout

import "github.com/a-h/templ"

// WidgetProps define as propriedades do componente WidgetWrapper
type WidgetProps struct {
	// Title é o título opcional do widget
	Title string

	// Subtitle é a descrição curta opcional do widget
	Subtitle string

	// Actions é o slot para botões/ações no header (opcional)
	Actions templ.Component

	// Content é o conteúdo principal do widget (obrigatório)
	Content templ.Component

	// HTMXAttributes permite passar atributos HTMX para a raiz do wrapper
	// Exemplo: hx-get="/api/data", hx-target="#widget", hx-trigger="load"
	HTMXAttributes templ.Attributes

	// NoPadding remove o padding interno para conteúdo full-bleed (opcional)
	// Útil para mapas, tabelas grandes ou conteúdo que precisa tocar as bordas
	NoPadding bool

	// Compact usa padding reduzido para widgets densos
	Compact bool

	// NoHeaderBorder remove a borda inferior do header
	NoHeaderBorder bool

	// Hover adiciona efeito hover de sombra
	Hover bool

	// ID identificador único do widget (opcional)
	ID string

	// Classes adicionais CSS para o wrapper (opcional)
	Classes string

	// Footer é o conteúdo do rodapé (opcional)
	Footer templ.Component
}

// DefaultWidgetProps retorna props padrão para o widget
func DefaultWidgetProps(content templ.Component) WidgetProps {
	return WidgetProps{
		Content:        content,
		NoPadding:      false,
		Compact:        false,
		NoHeaderBorder: false,
		Hover:          false,
	}
}

// WidgetPropsBuilder helper para construir WidgetProps
type WidgetPropsBuilder struct {
	props WidgetProps
}

// NewWidgetProps cria um novo builder para WidgetProps
func NewWidgetProps(content templ.Component) *WidgetPropsBuilder {
	return &WidgetPropsBuilder{
		props: DefaultWidgetProps(content),
	}
}

// WithTitle define o título do widget
func (b *WidgetPropsBuilder) WithTitle(title string) *WidgetPropsBuilder {
	b.props.Title = title
	return b
}

// WithSubtitle define o subtítulo do widget
func (b *WidgetPropsBuilder) WithSubtitle(subtitle string) *WidgetPropsBuilder {
	b.props.Subtitle = subtitle
	return b
}

// WithActions define as ações do header
func (b *WidgetPropsBuilder) WithActions(actions templ.Component) *WidgetPropsBuilder {
	b.props.Actions = actions
	return b
}

// WithHTMX define atributos HTMX
func (b *WidgetPropsBuilder) WithHTMX(attrs templ.Attributes) *WidgetPropsBuilder {
	b.props.HTMXAttributes = attrs
	return b
}

// WithNoPadding remove o padding interno
func (b *WidgetPropsBuilder) WithNoPadding() *WidgetPropsBuilder {
	b.props.NoPadding = true
	return b
}

// WithCompact usa padding reduzido
func (b *WidgetPropsBuilder) WithCompact() *WidgetPropsBuilder {
	b.props.Compact = true
	return b
}

// WithNoHeaderBorder remove a borda do header
func (b *WidgetPropsBuilder) WithNoHeaderBorder() *WidgetPropsBuilder {
	b.props.NoHeaderBorder = true
	return b
}

// WithHover adiciona efeito hover
func (b *WidgetPropsBuilder) WithHover() *WidgetPropsBuilder {
	b.props.Hover = true
	return b
}

// WithID define o ID do widget
func (b *WidgetPropsBuilder) WithID(id string) *WidgetPropsBuilder {
	b.props.ID = id
	return b
}

// WithClasses adiciona classes CSS
func (b *WidgetPropsBuilder) WithClasses(classes string) *WidgetPropsBuilder {
	b.props.Classes = classes
	return b
}

// WithFooter adiciona conteúdo ao footer
func (b *WidgetPropsBuilder) WithFooter(footer templ.Component) *WidgetPropsBuilder {
	b.props.Footer = footer
	return b
}

// Build retorna o WidgetProps construído
func (b *WidgetPropsBuilder) Build() WidgetProps {
	return b.props
}

// HasTitle retorna true se o widget tem título
func (p WidgetProps) HasTitle() bool {
	return p.Title != ""
}

// HasActions retorna true se o widget tem ações
func (p WidgetProps) HasActions() bool {
	return p.Actions != nil
}

// HasHTMX retorna true se o widget tem atributos HTMX
func (p WidgetProps) HasHTMX() bool {
	return p.HTMXAttributes != nil && len(p.HTMXAttributes) > 0
}

// HasFooter retorna true se o widget tem footer
func (p WidgetProps) HasFooter() bool {
	return p.Footer != nil
}

// GetWrapperClasses retorna as classes CSS para o wrapper
func (p WidgetProps) GetWrapperClasses() string {
	classes := "widget-wrapper"

	if p.NoPadding {
		classes += " widget-wrapper-no-padding"
	}

	if p.Compact && !p.NoPadding {
		classes += " widget-wrapper-compact"
	}

	if p.Hover {
		classes += " widget-wrapper-hover"
	}

	if p.Classes != "" {
		classes += " " + p.Classes
	}

	return classes
}

// GetHeaderClasses retorna as classes CSS para o header
func (p WidgetProps) GetHeaderClasses() string {
	classes := "widget-header"

	if p.NoHeaderBorder {
		classes += " widget-header-no-border"
	}

	return classes
}

// GetContentClasses retorna as classes CSS para o conteúdo
func (p WidgetProps) GetContentClasses() string {
	classes := "widget-content"

	if p.NoPadding {
		classes += " widget-content-no-padding"
	}

	return classes
}
