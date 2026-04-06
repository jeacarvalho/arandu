# Migração: Tela de Anamnese para Novo Layout Padrão

## Objetivo
Migrar `/patients/{id}/anamnesis` para ShellLayout + StandardGrid + WidgetWrapper.

## ⚠️ ATENÇÃO: Padrões Obrigatórios (NÃO IGNORAR)

### 1. Estrutura de Layout
| Componente | Uso |
|------------|-----|
| `ShellLayout` | **Obrigatório** - NÃO criar variant nova, usar padrão |
| `StandardGrid` | **Obrigatório** - Grid 4 colunas (1 nav + 3 conteúdo) |
| `WidgetWrapper` | **Obrigatório** - Cada seção de anamnese em um widget |

### 2. Tailwind Config (CRIAR/ATUALIZAR)
Adicione ESTES tokens no `tailwind.config.js` antes de começar:

```javascript
theme: {
  extend: {
    colors: {
      'arandu-primary': '#xxxxxx',
      'arandu-soft': '#xxxxxx',
      'arandu-bg': '#xxxxxx',
      'neutral-50': '#f9fafb',
      'neutral-200': '#e5e7eb',
    },
    fontFamily: {
      'serif-clinical': ['Source Serif 4', 'serif'],
    },
    spacing: {
      'anamnesis-nav': '240px',
    },
  },
}
```

### 3. Componentes a Criar (COM CÓDIGO EXEMPLO)

#### PageHeader
```templ
@page_header(title, subtitle, breadcrumbs) {
    <header class="mb-6">
        <nav class="breadcrumb"><!-- breadcrumbs --></nav>
        <h1 class="text-2xl font-bold">{ title }</h1>
        <p class="text-gray-600">{ subtitle }</p>
    </header>
}
```

#### AnamneseSection (COM WidgetWrapper)
```templ
@anamnesis_section(patientID, sectionID, title, content) {
    @widget_wrapper(widget.Props{
        Title: title,
        HTMXAttributes: templ.Attributes{
            "hx-patch": "/patients/" + patientID + "/anamnesis/" + sectionID,
            "hx-trigger": "keyup delay:2s, blur",
            "hx-target": "#indicator-" + sectionID,
            "hx-swap": "outerHTML",
        },
    }) {
        <textarea 
            class="anamnesis-textarea w-full min-h-[200px] font-serif-clinical"
            name="content"
        >{ content }</textarea>
        
        @autosave_indicator(sectionID)
    }
}
```

#### AnamneseNavigation (Sticky)
```templ
@anamnesis_navigation(patientID, activeSection) {
    <nav class="section-navigation sticky top-4 space-y-2">
        <a href="#queixa" class="nav-link { if activeSection == "queixa" }nav-link-active{/if}">
            <i class="fa-solid fa-stethoscope"></i> Queixa Principal
        </a>
        <!-- ... outras seções ... -->
    </nav>
}
```

### 4. CSS Classes (DEFINIR NO input-v2.css)

```css
/* Textarea clínica */
.anamnesis-textarea {
    @apply bg-neutral-50 border-neutral-200 rounded-md p-4;
    @apply focus:border-arandu-soft focus:ring-2 focus:ring-arandu-soft;
    @apply font-serif-clinical text-base leading-relaxed;
    min-height: 200px;
}

/* Navegação sticky */
.section-navigation {
    @apply bg-white p-4 rounded-lg border border-neutral-200;
    position: sticky;
    top: 24px; /* Padding da Main Canvas */
}

/* Link ativo */
.nav-link-active {
    @apply bg-arandu-bg text-arandu-primary font-semibold;
}

/* Auto-save indicator */
.autosave-indicator {
    @apply text-sm text-gray-500 flex items-center gap-2 mt-2;
}
.autosave-indicator.saving {
    @apply text-arandu-primary;
}
.autosave-indicator.saved {
    @apply text-green-600;
}
```

### 5. StandardGrid Config (EXEMPLO PRÁTICO)

```go
config := GridConfig{
    Rows: []GridRow{
        {
            Items: []GridItem{
                {Component: AnamneseNavigation(patientID, "queixa"), ColSpan: 1, ID: "nav-anamnese"},
                {Component: PageHeader("Anamnese - " + patient.Name, "História profunda...", breadcrumbs), ColSpan: 3, ID: "header-anamnese"},
            },
        },
        {
            Items: []GridItem{
                {Component: nil, ColSpan: 1, ID: "nav-spacer"}, // Navegação continua
                {Component: AnamneseSection(patientID, "queixa", "Queixa Principal", anamnesis.Queixa), ColSpan: 3, ID: "section-queixa"},
            },
        },
        // ... repetir para outras 3 seções ...
    },
}
```

### 6. HTMX Integration (ESPECÍFICO)

| Evento | Atributo | Valor |
|--------|----------|-------|
| Auto-save trigger | `hx-trigger` | `keyup delay:2s, blur` |
| Auto-save endpoint | `hx-patch` | `/patients/{id}/anamnesis/{section}` |
| Indicator target | `hx-target` | `#indicator-{sectionID}` |
| Loading state | `hx-indicator` | `.htmx-indicator` |

### 7. Responsividade (BREAKPOINTS EXPLÍCITOS)

```html
<!-- Desktop: Nav 1 col, Content 3 cols -->
<div class="grid grid-cols-4">

<!-- Tablet: Nav 1 col, Content 1 col -->
<div class="grid grid-cols-2 md:grid-cols-4">

<!-- Mobile: 1 col, nav em accordion -->
<div class="grid grid-cols-1 md:grid-cols-2">
```

### 8. PatientSidebar (CONTEXTUAL)

**IMPORTANTE:** O PatientSidebar NÃO é parte do ShellLayout. Ele deve ser:
- OU um widget dentro da StandardGrid (Col 1)
- OU parte do conteúdo da página (não da estrutura macro)

**NÃO crie** `ShellLayout variant="patient"` - isso não existe.

---

## ✅ Checklist de Validação (OBRIGATÓRIO)

Antes de entregar, valide CADA item:

- [ ] `tailwind.config.js` atualizado com tokens de cor/fonte/spacing
- [ ] `input-v2.css` com classes `.anamnesis-textarea`, `.section-navigation`, etc.
- [ ] Cada seção usa `WidgetWrapper` com props corretas
- [ ] StandardGrid com configuração 1+3 colunas (nav+conteúdo)
- [ ] HTMX `hx-trigger="keyup delay:2s, blur"` em cada textarea
- [ ] AutoSaveIndicator integrado com estados (oculto/salvando/gravado)
- [ ] Navegação sticky com `position: sticky; top: 24px`
- [ ] Fonte Source Serif 4 aplicada nas textareas
- [ ] Teste em mobile (grid colapsa para 1 coluna)
- [ ] Zero valores hardcoded (px, hex) no código

---

## 📎 Entregáveis Esperados

1. `tailwind.config.js` atualizado (trecho com novos tokens)
2. `input-v2.css` com classes clínicas
3. `anamnesis_page_v2.templ` (página principal)
4. `anamnesis_navigation.templ` (widget de navegação)
5. `anamnesis_section.templ` (widget de seção com WidgetWrapper)
6. `page_header.templ` (componente reutilizável)
7. `autosave_indicator.templ` (componente de estado)
8. Handler atualizado (`ShowAnamnesis` apontando para nova página)

---

## 🚨 Erros Comuns a Evitar

| Erro | Correção |
|------|----------|
| Criar `ShellLayout variant="patient"` | Use ShellLayout padrão + PatientSidebar como widget |
| Hardcoded `w-[240px]` na nav | Use `w-anamnesis-nav` do Tailwind Config |
| Hardcoded `#e5e7eb` nas bordas | Use `border-neutral-200` do Tailwind Config |
| Auto-save sem indicador visual | Sempre inclua AutoSaveIndicator |
| Navegação não-sticky em desktop | Use `sticky top-4` com offset correto |
| Textarea sem fonte clínica | Aplique `font-serif-clinical` |

---

## 📸 Referência Visual (Descrição)

```
┌─────────────────────────────────────────────────────────────────┐
│ Topbar (fixa)                                                   │
├─────────────┬───────────────────────────────────────────────────┤
│ Sidebar     │ Main Canvas (scrollável)                          │
│ (pacientes) │ ┌─────────────────────────────────────────────┐   │
│             │ │ PageHeader: "Anamnese - Baba"               │   │
│             │ ├─────────────────────────────────────────────┤   │
│             │ │ ┌───────────┬─────────────────────────────┐ │   │
│             │ │ │ Navigation│ Widget: Queixa Principal    │ │   │
│             │ │ │ (sticky)  │ [textarea com auto-save]    │ │   │
│             │ │ │           │ [indicador "Gravado"]       │ │   │
│             │ │ ├───────────┼─────────────────────────────┤ │   │
│             │ │ │           │ Widget: História Pessoal    │ │   │
│             │ │ │           │ [textarea com auto-save]    │ │   │
│             │ │ ├───────────┼─────────────────────────────┤ │   │
│             │ │ │           │ Widget: História Familiar   │ │   │
│             │ │ ├───────────┼─────────────────────────────┤ │   │
│             │ │ │           │ Widget: Exame Mental        │ │   │
│             │ │ └───────────┴─────────────────────────────┘ │   │
│             │ └─────────────────────────────────────────────┘   │
└─────────────┴───────────────────────────────────────────────────┘
```

---

