# TASK 

**Requirement**: REQ-03-01-01
**Title**: Classificação de Observações Clínicas
**Capability**: CAP-03-01 — Organização de Observações Clínicas
**Vision**: VISION-03 — Organização do Conhecimento Clínico
**Status**: PRONTO_PARA_IMPLEMENTACAO

---

## 📝 O Que Construir

Sistema de classificação/tagging para observações clínicas, permitindo que o terapeuta categorize percepções por tipo (emoção, comportamento, cognição, relação) e intensidade, criando estrutura para análise de padrões futuros.

---

## 🧩 Componentes Necessários (Grid 4x4)

| Widget | ColSpan | Descrição |
|--------|---------|-----------|
| ObservationFilters | 4 | Filtros por tipo de classificação (HTMX) |
| ObservationList | 4 | Lista de observações com badges de classificação |
| ClassificationSelector | 4 | Seletor de tags para cada observação (inline) |
| ClassificationSummary | 4 | Resumo visual de distribuição por tipo |

---

## 📋 Regras Específicas

1. **Tipos de Classificação** (fixos, via enum Go):
   - Emoção (😢)
   - Comportamento (👤)
   - Cognição (🧠)
   - Relação (👥)
   - Somático (💪)
   - Contexto (🌍)

2. **Intensidade**: Escala 1-5 (opcional, via slider ou selects)

3. **HTMX**:
   - Classificação inline via `hx-post="/observations/{id}/classify"`
   - Filtros via `hx-trigger="change"` nos selects
   - Target: `#observation-list`

4. **Persistência**:
   - Tabela `observation_classifications` (many-to-many)
   - Campos: `observation_id`, `classification_type`, `intensity`, `created_at`

5. **Visual**:
   - Badges coloridos por tipo (cores do Design System)
   - Tipografia clínica (`.font-clinical`) no conteúdo
   - Zero inline styles

---

## ⚠️ PADRÕES OBRIGATÓRIOS (NÃO IGNORAR)

### 1. Estrutura de Layout

| Componente | Uso |
|------------|-----|
| ShellLayout | Obrigatório (TopBar + Sidebar + Main Canvas) |
| StandardGrid | Obrigatório (grid 4 colunas) |
| WidgetWrapper | Obrigatório em cada widget |

### 2. Stack & Estilo

| Item | Regra |
|------|-------|
| Stack | Go + Templ + Tailwind + HTMX |
| Cores/Spacing | Apenas tokens do `tailwind.config.js` |
| Margins | Zero margins externas (gap via StandardGrid) |
| Padding | Apenas via WidgetWrapper |
| Alturas | Automáticas (content-based) |

### 3. HTMX

| Atributo | Uso |
|----------|-----|
| hx-post | `/observations/{id}/classify` |
| hx-target | `#observation-{id}` |
| hx-swap | `outerHTML` |
| hx-trigger | `change` (nos selects de tag) |

### 4. Responsividade

| Breakpoint | Comportamento |
|------------|---------------|
| Desktop | Respeitar col-span definido |
| Mobile (<768px) | Grid 1 coluna, filtros empilhados |

### 5. Go Structs

```go
type ObservationClassification struct {
    ID            string
    ObservationID string
    Type          ClassificationType // emotion, behavior, cognition, relationship, somatic, context
    Intensity     int                // 1-5
    CreatedAt     time.Time
}

type ClassificationData struct {
    ObservationID  string
    AvailableTypes []ClassificationType
    SelectedTypes  []ObservationClassification
}

type ClassificationType struct {
    ID    string
    Label string
    Icon  string
    Color string // token CSS
}
```

---

## 🆕 SEÇÕES OBRIGATÓRIAS DE ESPECIFICAÇÃO TÉCNICA

### 6. 🆕 Tailwind Config (Tokens a Criar/Verificar)

```js
// tailwind.config.js - Adicionar/Verificar:
theme: {
  extend: {
    colors: {
      'class-emotion': '#0F6E56',      // Verde para emoções
      'class-behavior': '#1D9E75',     // Verde ativo para comportamentos
      'class-cognition': '#7C3AED',    // Roxo para cognições
      'class-relationship': '#F59E0B', // Âmbar para relações
      'class-somatic': '#DC2626',      // Vermelho para somático
      'class-context': '#6B7280',      // Cinza para contexto
    },
    spacing: {
      'classification-badge': '28px',
    },
  },
}
```

### 7. 🆕 Componentes .templ a Criar/Reutilizar

```
web/components/classification/
├── filters.templ           // ObservationFilters (hx-trigger="change")
├── selector.templ          // ClassificationSelector (inline)
├── summary.templ           // ClassificationSummary
├── badge.templ             // ClassificationBadge (reutilizável)
└── observation_item.templ  // Atualizar existente com classificação
```

### 8. 🆕 Handler Go

```
internal/web/handlers/classification_handler.go
- POST /observations/{id}/classify
- GET /observations/{id}/classify/edit
- GET /observations/filter?classification={type}
```

### 9. 🆕 Validações Obrigatórias

```bash
# 1. Zero inline styles
grep -o 'style="' web/components/classification/*.templ | wc -l  # Deve ser 0

# 2. Grid responsivo
grep -q "grid-cols-1 md:grid-cols-4" web/components/classification/*.templ

# 3. HTMX targets definidos
grep -q 'hx-target="#observation-' web/components/classification/*.templ

# 4. Tipografia clínica
grep -q "font-clinical" web/components/classification/*.templ

# 5. E2E audit
./scripts/arandu_e2e_audit.sh --routes observations
```

### 10. 🆕 Critérios de Aceitação

- [ ] **CA-01**: Tags podem ser adicionadas/removidas via HTMX sem reload
- [ ] **CA-02**: Filtro por tipo atualiza lista de observações
- [ ] **CA-03**: Intensidade visual (cores) reflete valor 1-5
- [ ] **CA-04**: Resumo mostra distribuição de tags por tipo
- [ ] **CA-05**: Mobile: seletor de tags funciona em touch
- [ ] **CA-06**: Zero inline styles detectados
- [ ] **CA-07**: E2E audit passa sem erros SLP
- [ ] **CA-08**: Classificação persiste após refresh de página

---

## 📚 Documentação de Referência

| Documento | Seção | Por Que Ler |
|-----------|-------|-------------|
| `docs/requirements/req-03-01-01-classificar-observacao.md` | Requisitos | Critérios de negócio |
| `docs/architecture/WEB_LAYER_PATTERN.md` | HTMX | Padrões de swap |
| `docs/design-system.md` | Cores | Tokens de classificação |
| `docs/learnings/MASTER_LEARNINGS.md` | Anti-Padrões | O que NÃO fazer |
| `padrao_layout_grids.md` | StandardGrid | Grid 4x4 declarativo |
| `prompt_padrao_layout_grid_tela.md` | Template | Estrutura do prompt |

---

## 🚨 Anti-Padrões a Evitar

| Anti-Padrão | Por Que Evitar | Alternativa |
|-------------|----------------|-------------|
| `style="..."` inline | Falha no E2E audit | Classes CSS em style.css |
| Grid sem responsividade | Quebra em mobile | `grid-cols-1 md:grid-cols-4` |
| Tags hardcoded no HTML | Dificulta manutenção | Tabela SQL + seed data |
| Ignorar HX-Request | Quebra HTMX | Verificar header sempre |
| Hardcoded colors | Viola design system | Tokens do tailwind.config.js |
| Margins externas nos widgets | Quebra grid | Gap via StandardGrid |

---

## 🧪 Validação Pós-Implementação

```bash
# 1. Build
go build ./cmd/arandu

# 2. Quality gates
./scripts/arandu_guard.sh

# 3. E2E audit
./scripts/arandu_e2e_audit.sh --routes observations

# 4. Validação visual
./scripts/arandu_visual_check.sh

# 5. Screenshot
./scripts/arandu_screenshot.sh

# 6. Verificar inline styles
grep -o 'style="' web/components/classification/*.templ | wc -l  # Deve ser 0
```

---

## 📝 Checklist de Conclusão

Antes de executar `arandu_conclude_task.sh`:

- [ ] Li toda documentação obrigatória
- [ ] Implementei conforme especificação
- [ ] Rodei `go build` sem erros
- [ ] Rodei `./scripts/arandu_guard.sh` e passou
- [ ] Rodei E2E audit e passou
- [ ] Validação visual manual OK (desktop + mobile)
- [ ] Screenshot gerado e revisado
- [ ] Zero inline styles detectados
- [ ] Documentei aprendizados em `docs/learnings/MASTER_LEARNINGS.md`

---

**Instrução Final**: Esta classificação é a base para toda inteligência futura do sistema (padrões, IA, comparação de casos). Não pule validações — a estrutura de tags deve ser consistente e extensível.

# ORIENTAÇÕES GERAIS DE LAYOUT

# 📐 Layout Arandu — Desenho ASCII Completo

## 🏗️ 1. ESTRUTURA GERAL SLP (Desktop)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                              TOP BAR (fixa)                                                    │
│  z-index: 100 | height: 64px | background: #FFFFFF (arandu-paper)                                             │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                 │
│  ┌──────────────────────┐      ┌────────────────────────────────────────────────────────────────────────────┐  │
│  │                      │      │                                                                             │  │
│  │                      │      │   MAIN CANVAS (scrollável)                                                 │  │
│  │                      │      │   z-index: 1 | margin-left: 280px | padding-top: 80px                      │  │
│  │      SIDEBAR         │      │   background: #E1F5EE (arandu-bg)                                          │  │
│  │                      │      │                                                                             │  │
│  │   width: 280px       │      │   ┌─────────────────────────────────────────────────────────────────────┐  │  │
│  │   z-index: 50        │      │   │                                                                     │  │  │
│  │   fixed (desktop)    │      │   │  ┌─────────────────────────────────────────────────────────────┐   │  │  │
│  │                      │      │   │  │  PATIENT PROFILE HEADER                                      │   │  │  │
│  │   [🏠] Resumo        │      │   │  │  ┌────┐  Carolina Costa  [Em tratamento] [TAG] [22 anos]    │   │  │  │
│  │   [📋] Anamnese      │      │   │  │  │ CC │  Feminino · Branca · Estudante                      │   │  │  │
│  │   [📄] Prontuário    │      │   │  │  └────┘  [124 sessões] [2,4a em terapia]                    │   │  │  │
│  │   [🎯] Plano         │      │   │  │                                                              │   │  │  │
│  │       Terapêutico    │      │   │  └─────────────────────────────────────────────────────────────┘   │  │  │
│  │                      │      │   │                                                                     │  │  │
│  │   [←] Voltar ao      │      │   │   ┌────────────────────────────┐  ┌────────────────────────────┐  │  │  │
│  │       Dashboard      │      │   │   │                            │  │                            │  │  │  │
│  │                      │      │   │   │  📋 Identidade             │  │  📝 Notas de               │  │  │  │
│  │   [🚪] Sair          │      │   │   │     Biopsicossocial        │  │     Triagem                │  │  │  │
│  │                      │      │   │   │     (grid-card-identity)   │  │     (grid-card-notes)      │  │  │  │
│  │                      │      │   │   │                            │  │                            │  │  │  │
│  │                      │      │   │   │     Gênero: Feminino       │  │     gases                  │  │  │  │
│  │                      │      │   │   │     Etnia: Branca          │  │                            │  │  │  │
│  │                      │      │   │   │     Ocupação: Estudante    │  │                            │  │  │  │
│  │                      │      │   │   │     Escolaridade: Superior │  │                            │  │  │  │
│  │                      │      │   │   │                            │  │                            │  │  │  │
│  │                      │      │   │   └────────────────────────────┘  └────────────────────────────┘  │  │  │
│  │                      │      │   │                                                                     │  │  │
│  │                      │      │   │   ┌────────────────────────────┐  ┌────────────────────────────┐  │  │  │
│  │                      │      │   │   │                            │  │                            │  │  │  │
│  │                      │      │   │   │  📅 Linha do Tempo         │  │  ⚡ Ações Rápidas          │  │  │  │
│  │                      │      │   │   │     (grid-card-timeline)   │  │     (grid-card-actions)    │  │  │  │
│  │                      │      │   │   │                            │  │                            │  │  │  │
│  │                      │      │   │   │     • 📅 Cadastro          │  │     [📅 Nova Sessão]       │  │  │  │
│  │                      │      │   │   │       realizado            │  │                            │  │  │  │
│  │                      │      │   │   │       24/03/2026           │  │     [🔄 Completar          │  │  │  │
│  │                      │      │   │   │                            │  │         Anamnese]          │  │  │  │
│  │                      │      │   │   │     [↻ Ver histórico       │  │                            │  │  │  │
│  │                      │      │   │   │         completo]          │  │     [🎯 Definir Metas]     │  │  │  │
│  │                      │      │   │   │                            │  │                            │  │  │  │
│  │                      │      │   │   └────────────────────────────┘  └────────────────────────────┘  │  │  │
│  │                      │      │   │                                                                     │  │  │
│  │                      │      │   │   ┌─────────────────────────────────────────────────────────────┐  │  │  │
│  │                      │      │   │   │  📅 SESSÕES RECENTES                                        │  │  │  │
│  │                      │      │   │   │                                                             │  │  │  │
│  │                      │      │   │   │  18/05/2025  │ Sessão #1 │ Paciente em fase de melhora... │  │  │  │
│  │                      │      │   │   │  11/05/2025  │ Sessão #2 │ Paciente em fase em processo...│  │  │  │
│  │                      │      │   │   │  04/05/2025  │ Sessão #3 │ Paciente em fase em processo...│  │  │  │
│  │                      │      │   │   │  27/04/2025  │ Sessão #4 │ Paciente em fase consolidação..│  │  │  │
│  │                      │      │   │   │  18/04/2025  │ Sessão #5 │ Paciente em fase de melhora... │  │  │  │
│  │                      │      │   │   │                                                             │  │  │  │
│  │                      │      │   │   └─────────────────────────────────────────────────────────────┘  │  │  │
│  │                      │      │   │                                                                     │  │  │
│  └──────────────────────┘      └────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

LEGEND:
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

┌────────────────────────────┐  = Clinical Card (WidgetWrapper)
│                            │    - background: var(--arandu-paper)
│                            │    - border-radius: 24px (var(--radius-xl))
│                            │    - padding: 24px (var(--space-xl))
│                            │    - border: 1px solid var(--neutral-200)
└────────────────────────────┘    - shadow: var(--shadow-sm)

→ = Chevron/Arrow (navigation indicator)

📋 📝 📅 ⚡ 🔄 🎯 ← = Icons (FontAwesome)

══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

DIMENSIONS (Desktop > 1024px):
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

Top Bar:        height: 64px (fixed, z-index: 100)
Sidebar:        width: 280px (fixed, z-index: 50)
Main Canvas:    flex: 1, margin-left: 280px, padding-top: 80px
Grid Columns:   2 columns (50%/50% or 65%/35% depending on content)
Grid Gap:       24px (var(--space-lg))
Card Padding:   24px (var(--space-xl))
Card Radius:    24px (var(--radius-xl))

══════════════════════════════════════════════════════════════════════════════════════════════════════════════════
```

---

## 📱 2. LAYOUT MOBILE (< 768px)

```
┌─────────────────────────────────────────────────────────────────┐
│                    TOP BAR (56px)                               │
│  [☰]  Arandu              🔍              [👤]                 │
│  z-index: 100                                                   │
└─────────────────────────────────────────────────────────────────┘
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  PATIENT HEADER (100% width)                              │ │
│  │  ┌────┐  Carolina Costa                                   │ │
│  │  │ CC │  [Em tratamento] [TAG] [22 anos]                  │ │
│  │  └────┘  Feminino · Branca · Estudante                    │ │
│  │         [124 sessões] [2,4a em terapia]                    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  📋 Identidade Biopsicossocial                            │ │
│  │  (100% width - cards empilham)                            │ │
│  │                                                           │ │
│  │  Gênero: Feminino                                         │ │
│  │  Etnia: Branca                                            │ │
│  │  Ocupação: Estudante                                      │ │
│  │  Escolaridade: Superior                                   │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  📝 Notas de Triagem                                      │ │
│  │  (100% width)                                             │ │
│  │                                                           │ │
│  │  gases                                                    │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  📅 Linha do Tempo                                        │ │
│  │  (100% width)                                             │ │
│  │                                                           │ │
│  │  • 📅 Cadastro realizado 24/03/2026                       │ │
│  │  [↻ Ver histórico completo]                               │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  ⚡ Ações Rápidas                                         │ │
│  │  (100% width - botões empilham)                           │ │
│  │                                                           │ │
│  │  [📅 Nova Sessão]                                         │ │
│  │  [🔄 Completar Anamnese]                                  │ │
│  │  [🎯 Definir Metas]                                       │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  📅 SESSÕES RECENTES                                      │ │
│  │  (100% width)                                             │ │
│  │                                                           │ │
│  │  18/05/2025 │ Sessão #1 │ Paciente em fase de...         │ │
│  │  11/05/2025 │ Sessão #2 │ Paciente em fase em...         │ │
│  │  04/05/2025 │ Sessão #3 │ Paciente em fase em...         │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

MOBILE BEHAVIOR:
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

Top Bar:        height: 56px (reduced from 64px)
Sidebar:        Drawer (hidden by default, hamburger menu trigger)
Main Canvas:    width: 100%, margin-left: 0, padding: 16px
Grid Columns:   1 column (100% width) - ALL cards stack vertically
Grid Gap:       16px (reduced from 24px)
Cards:          Full width, stacked vertically
Padding:        16px (reduced from 24px)

══════════════════════════════════════════════════════════════════════════════════════════════════════════════════
```

---

## 🎯 3. GRID SYSTEM DETALHADO (Desktop)

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                         MAIN CANVAS CONTENT AREA                                               │
│                                   (padding-top: 80px | background: #E1F5EE)                                    │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

GRID 4x4 SYSTEM (Reference for all pages):
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

┌─────────────┬─────────────┬─────────────┬─────────────┐
│             │             │             │             │
│  ColSpan 1  │  ColSpan 1  │  ColSpan 1  │  ColSpan 1  │
│   (25%)     │   (25%)     │   (25%)     │   (25%)     │
│             │             │             │             │
├─────────────┴─────────────┼─────────────┴─────────────┤
│                           │                           │
│      ColSpan 2 (50%)      │      ColSpan 2 (50%)      │
│                           │                           │
├───────────────────────────┴───────────────────────────┤
│                                                       │
│               ColSpan 4 (100%)                        │
│                                                       │
└───────────────────────────────────────────────────────┘

PATIENT PROFILE GRID (2 columns on desktop):
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

┌───────────────────────────────────────────────┬───────────────────────────────────────────────┐
│                                               │                                               │
│  📋 Identidade Biopsicossocial                │  📝 Notas de Triagem                          │
│  (grid-card-identity)                         │  (grid-card-notes)                            │
│  ColSpan: 1 (50%)                             │  ColSpan: 1 (50%)                             │
│  Min-Height: 220px                            │  Min-Height: 220px                            │
│                                               │                                               │
│  Gênero, Etnia, Ocupação,                     │  gases                                        │
│  Escolaridade                                 │                                               │
│                                               │                                               │
├───────────────────────────────────────────────┼───────────────────────────────────────────────┤
│                                               │                                               │
│  📅 Linha do Tempo                            │  ⚡ Ações Rápidas                             │
│  (grid-card-timeline)                         │  (grid-card-actions)                          │
│  ColSpan: 1 (50%)                             │  ColSpan: 1 (50%)                             │
│  Min-Height: 220px                            │  Min-Height: 220px                            │
│                                               │                                               │
│  • Cadastro realizado                         │  [📅 Nova Sessão]                             │
│    24/03/2026                                 │  [🔄 Completar Anamnese]                      │
│  [↻ Ver histórico completo]                   │  [🎯 Definir Metas]                           │
│                                               │                                               │
├───────────────────────────────────────────────┴───────────────────────────────────────────────┤
│                                                                                               │
│  📅 SESSÕES RECENTES (Full Width)                                                             │
│  (grid-card-sessions)                                                                         │
│  ColSpan: 2 (100%)                                                                            │
│  Min-Height: auto                                                                             │
│                                                                                               │
│  18/05/2025 │ Sessão #1 │ Paciente em fase de melhora...                                    │
│  11/05/2025 │ Sessão #2 │ Paciente em fase em processo...                                   │
│  04/05/2025 │ Sessão #3 │ Paciente em fase em processo...                                   │
│  27/04/2025 │ Sessão #4 │ Paciente em fase consolidação...                                  │
│  18/04/2025 │ Sessão #5 │ Paciente em fase de melhora...                                    │
│                                                                                               │
└───────────────────────────────────────────────────────────────────────────────────────────────┘

GRID SPECIFICATIONS:
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

Grid Gap:       24px (var(--space-lg))
Card Padding:   24px (var(--space-xl))
Card Radius:    24px (var(--radius-xl))
Card Border:    1px solid var(--neutral-200)
Card Shadow:    var(--shadow-sm)
Card Min-Height: 220px (for top cards)

Desktop (> 1024px):   2 columns (50%/50%)
Tablet (768-1024px):  2 columns (50%/50%)
Mobile (< 768px):     1 column (100%)

══════════════════════════════════════════════════════════════════════════════════════════════════════════════════
```

---

## 🎨 4. Z-INDEX LAYERING (Camadas de Sobreposição)

```
Z-INDEX HIERARCHY (Highest to Lowest):
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  z-index: 200  │  Modals / Dialogs (Goal Closure, etc.)                                                        │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 150  │  Dropdowns / Popovers                                                                         │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 100  │  TOP BAR (fixed header)                                                                       │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 95   │  Sidebar Overlay (mobile backdrop)                                                            │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 90   │  Sidebar Drawer (mobile)                                                                      │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 50   │  SIDEBAR (desktop persistent)                                                                 │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 10   │  Cards / Widgets                                                                              │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 1    │  MAIN CANVAS (background)                                                                     │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│  z-index: 0    │  Page Background (#E1F5EE)                                                                    │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

CRITICAL RULES:
══════════════════════════════════════════════════════════════════════════════════════════════════════════════════

✅ Top Bar MUST be z-index: 100 (above sidebar and content)
✅ Sidebar MUST be z-index: 50 (below top bar, above content)
✅ Main Content MUST be z-index: 1 (below sidebar)
✅ Modals MUST be z-index: 200 (above everything)
❌ NEVER use hardcoded z-index (use CSS variables)
❌ NEVER set main-content z-index higher than sidebar

══════════════════════════════════════════════════════════════════════════════════════════════════════════════════
```

---

## 📊 5. DIMENSION SUMMARY TABLE

```
┌─────────────────────────────────┬──────────────────┬──────────────────┬──────────────────┐
│ Component                       │ Desktop (>1024px)│ Tablet (768-1024)│ Mobile (<768px)  │
├─────────────────────────────────┼──────────────────┼──────────────────┼──────────────────┤
│ Top Bar Height                  │ 64px             │ 64px             │ 56px             │
│ Sidebar Width                   │ 280px            │ 280px            │ Drawer (hidden)  │
│ Sidebar Behavior                │ Fixed            │ Fixed            │ Drawer           │
│ Main Canvas Margin-Left         │ 280px            │ 280px            │ 0                │
│ Main Canvas Padding-Top         │ 80px             │ 80px             │ 72px             │
│ Main Canvas Padding             │ 32px             │ 24px             │ 16px             │
│ Grid Columns                    │ 2 (50%/50%)      │ 2 (50%/50%)      │ 1 (100%)         │
│ Grid Gap                        │ 24px             │ 20px             │ 16px             │
│ Card Padding                    │ 24px             │ 20px             │ 16px             │
│ Card Min-Height                 │ 220px            │ 200px            │ auto             │
│ Card Border Radius              │ 24px             │ 20px             │ 16px             │
│ Font Size (Clinical Content)    │ 1.125rem         │ 1rem             │ 1rem             │
│ Font Size (Patient Name)        │ 2rem             │ 1.75rem          │ 1.5rem           │
└─────────────────────────────────┴──────────────────┴──────────────────┴──────────────────┘
```

