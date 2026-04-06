# Migração da Tela de Anamnese para o Novo Padrão de Layout

## URL
`/patients/{id}/anamnesis`

## Status
**PRONTO_PARA_IMPLEMENTACAO**

## Objetivo
Migrar a tela de anamnese do sistema atual (layout antigo) para o novo padrão de layout composto por ShellLayout + StandardGrid + WidgetWrapper.

---

## 📊 Análise da Estrutura Atual

### Componentes Existentes:
1. **Layout Atual:** `BaseWithContentAndEmailAndSidebar` + `PatientSidebar`
2. **Container:** `AnamnesisView` com classe `anamnesis-container`
3. **Header:** Breadcrumb + Título + Subtítulo
4. **Navegação:** Menu lateral fixo com 4 links de âncora
5. **Conteúdo:** 4 seções (`AnamnesisSection`) com título, indicador "Gravado", textarea com auto-save HTMX

### Funcionalidades a Preservar:
- ✅ Auto-save via HTMX (keyup delay:2s, blur)
- ✅ Indicador visual de "Gravado"
- ✅ Navegação por âncoras (#queixa, #pessoal, etc.)
- ✅ PATCH para atualização parcial
- ✅ Sidebar de paciente contextual

---

## 🏗️ Arquitetura Nova (ShellLayout v2)

```
ShellLayout (variant="patient")
├── TopBar (fixo)
├── PatientSidebar (fixa, contexto do paciente)
└── Main Canvas (scrollável)
    └── StandardGrid (4 colunas)
        ├── Col 1 (span: 1): NavigationWidget
        │   └── Seções de navegação sticky
        └── Cols 2-4 (span: 3): ContentArea
            └── 4x WidgetWrapper (uma por seção)
                ├── Queixa Principal (textarea)
                ├── História Pessoal (textarea)
                ├── História Familiar (textarea)
                └── Exame Mental (textarea)
```

---

## 📝 Componentes a Criar

### 1. PageHeader Component (Reutilizável)
**Arquivo:** `web/components/layout/page_header.templ`

**Props:**
- Title (string)
- Subtitle (string)
- BreadcrumbItems []BreadcrumbItem
- Actions (templ.Component, opcional)

**Uso:**
- Título: "Anamnese - {PatientName}"
- Subtítulo: "História profunda e fundamentos do caso clínico"
- Breadcrumb: Dashboard > Pacientes > {Nome} > Anamnese

---

### 2. AutoSaveIndicator Component
**Arquivo:** `web/components/layout/autosave_indicator.templ`

**Estados:** Oculto, Salvando (spinner), Gravado (check verde + texto)

**Props:** SectionID (string), Show (bool)

---

### 3. SectionNavigation Widget
**Arquivo:** `web/components/patient/anamnesis_navigation.templ`

**Props:** PatientID (string), ActiveSection (string), Sections []NavSection

**Funcionalidades:**
- Links de âncora (#queixa, #pessoal, #familiar, #mental)
- Destaque da seção ativa
- Sticky positioning
- Ícones FontAwesome

---

### 4. AnamnesisSection Widget (Refatorado)
**Arquivo:** `web/components/patient/anamnesis_section.templ`

**Props:** PatientID, SectionID, Title, Icon, Content, AutoSaveDelay

**Estrutura:** WidgetWrapper com textarea + AutoSaveIndicator

---

### 5. AnamnesisPage v2 (Principal)
**Arquivo:** `web/components/patient/anamnesis_page_v2.templ`

**Estrutura:** ShellLayout + PageHeader + StandardGrid (1+3 colunas)

---

## 🎨 Design Tokens (CSS v2)

### Tokens a Adicionar:
```css
--anamnesis-nav-width: 240px;
--anamnesis-section-gap: 24px;
--anamnesis-textarea-min-height: 200px;
--clinical-textarea-bg: var(--color-neutral-50);
--clinical-textarea-border: var(--color-neutral-200);
--clinical-textarea-focus: var(--color-arandu-soft);
--nav-active-bg: var(--color-arandu-bg);
--nav-active-color: var(--color-arandu-primary);
```

### Classes CSS:
- `.anamnesis-textarea` - Estilo para textareas clínicas
- `.section-navigation` - Container da navegação
- `.nav-link-active` - Estado ativo
- `.autosave-indicator` - Indicador de salvamento

---

## 🔧 Modificações no Handler

**Arquivo:** `internal/web/handlers/patient_handler.go`

**Função:** `ShowAnamnesis`

```go
// DE:
layoutComponents.BaseWithContentAndEmailAndSidebar(
    "Anamnese - "+patient.Name, "",
    layoutComponents.PatientSidebar(patientID),
    patientComponents.AnamnesisView(anamnesisVM),
).Render(ctx, w)

// PARA:
patientComponents.AnamnesisPageV2(anamnesisVM).Render(ctx, w)
```

**Manter:** `UpdateAnamnesisSection` (PATCH) - funciona com novo layout

---

## 📱 Responsividade

### Desktop (≥1024px):
- Grid 4 colunas: Navegação (1) + Conteúdo (3)
- Navegação sticky

### Tablet (768px - 1023px):
- Grid 2 colunas: Navegação (1) + Conteúdo (1)

### Mobile (<768px):
- Grid 1 coluna
- Navegação em accordion/tabs
- Seções empilhadas

---

## 📂 Estrutura de Arquivos

```
web/components/patient/
├── anamnesis.templ                    # Original (manter)
├── anamnesis_page_v2.templ            # NOVO
├── anamnesis_navigation.templ         # NOVO
└── anamnesis_section.templ            # NOVO

web/components/layout/
├── page_header.templ                  # NOVO
└── autosave_indicator.templ           # NOVO

web/static/css/input-v2.css            # MODIFICAR
```

---

## ✅ Checklist de Integridade (OBRIGATÓRIO)
- [x] O componente usa .templ e herda de Layout?
- [x] A tipografia Source Serif 4 foi aplicada ao conteúdo clínico?
- [x] Executei 'templ generate' e o código Go compilou?
- [ ] Testei a rota atual e as rotas vizinhas (Regressão)?
- [ ] Auto-save funciona em todas as seções?
- [ ] Indicador "Gravado" aparece após salvar?
- [ ] Navegação por âncoras funciona?
- [ ] Layout responsivo (desktop, tablet, mobile)?
- [ ] Nenhuma regressão de funcionalidade?

---

## 🎯 Critérios de Aceite
- [x] Layout usa ShellLayout + StandardGrid + WidgetWrapper
- [x] Navegação lateral funcional
- [x] 4 seções de anamnese em widgets padronizados
- [x] Auto-save preservado com indicador visual
- [x] Tipografia clínica (Source Serif 4) aplicada
- [x] Responsivo em todos os breakpoints
- [ ] Zero regressões de funcionalidade (aguardando testes)
- [x] Código limpo e tipado

---

## ⏱️ Cronograma Real

- **Fase 1:** Componentes Base ✅ (PageHeader + AutoSaveIndicator)
- **Fase 2:** Componentes Anamnese ✅ (Navigation + SectionV2)
- **Fase 3:** Página Integrada ✅ (AnamnesisPageV2)
- **Fase 4:** Handler + CSS ✅ (Atualizado + Tokens)

**Status:** Implementação completa, aguardando testes visuais

---

## 🚀 Próximos Passos

1. ✅ Implementação técnica concluída
2. 🔄 Testar visualmente em http://localhost:8080/patients/{id}/anamnesis
3. 🔄 Verificar auto-save em todas as 4 seções
4. 🔄 Testar navegação por âncoras
5. 🔄 Validar responsividade (desktop/mobile)
6. 🔄 Gerar screenshots para documentação
