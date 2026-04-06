# TASK 

**Requirement**: REQ-02-02-01
**Title**: Linha do Tempo Clínica — Visualização Longitudinal Unificada
**Status**: PRONTO_PARA_IMPLEMENTACAO
**Capability**: CAP-02-02 — Linha do tempo clínica do paciente
**Vision**: VISION-02 — Memória clínica longitudinal

---

## 📝 O Que Construir

Implementar visão consolidada da linha do tempo clínica do paciente, agregando sessões, observações, intervenções e metas em uma única interface cronológica com filtros por tipo de evento e período.

---

## 🧩 Componentes Necessários (Grid 4x4)

| Widget | ColSpan | Descrição |
|--------|---------|-----------|
| TimelineFilters | 4 | Filtros: tipo de evento, período, busca textual |
| TimelineContainer | 4 | Container principal com scroll vertical |
| TimelineEvent (repetível) | 4 | Card de evento (sessão/obs/intervenção/meta) |
| TimelineEmpty | 4 | Estado vazio com CTA para criar primeira sessão |

---

## 📋 Regras Específicas

1. **Ordenação**: Mais recente primeiro (padrão), com opção de inverter
2. **Filtros HTMX**: `hx-trigger="change"` nos selects, `hx-target="#timeline-container"`
3. **Busca textual**: `hx-trigger="keyup changed delay:500ms"` com debounce
4. **Cada evento**: Link para edição via `hx-get` + swap no próprio card
5. **Indicador visual por tipo**:
   - 🟢 Sessão
   - 🔵 Observação
   - 🟣 Intervenção
   - 🟡 Meta
6. **Mobile**: Cards ocupam 100% largura, filtros empilham verticalmente
7. **Infinite Scroll**: `hx-trigger="revealed"` para carregamento progressivo (lotes de 20)

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
| Cores/Spacing | Apenas tokens do tailwind.config.js |
| Margins | Zero margins externas (gap via StandardGrid) |
| Padding | Apenas via WidgetWrapper |
| Alturas | Automáticas (content-based) |

### 3. HTMX

| Atributo | Uso |
|----------|-----|
| hx-get | Carregamento inicial e filtros |
| hx-target | `#timeline-container` para swap de conteúdo |
| hx-trigger | `change` para filtros, `keyup changed delay:500ms` para busca |
| Loading | Skeleton por card durante carregamento |

### 4. Responsividade

| Breakpoint | Comportamento |
|------------|---------------|
| Desktop | Respeitar col-span definido |
| Mobile (<768px) | Grid 1 coluna, filtros empilhados |

### 5. Go Structs

```go
type TimelineEvent struct {
    ID          string
    Type        EventType // session, observation, intervention, goal
    Date        time.Time
    Title       string
    Content     string
    PatientID   string
    SessionID   string // opcional
    HTMXGetURL  string // para edição inline
}

type TimelineConfig struct {
    PatientID   string
    FilterType  string // all, session, observation, etc.
    DateFrom    string
    DateTo      string
    SearchQuery string
    Limit       int
    Offset      int
}
```

---

## 🆕 SEÇÕES OBRIGATÓRIAS DE ESPECIFICAÇÃO TÉCNICA

### 6. Tailwind Config (Tokens a Criar/Verificar)

```js
// tailwind.config.js - Verificar/Adicionar:
theme: {
  extend: {
    colors: {
      'timeline-session': '#0F6E56',      // Verde para sessões
      'timeline-observation': '#1D9E75',  // Verde ativo para observações
      'timeline-intervention': '#7C3AED', // Roxo para intervenções
      'timeline-goal': '#F59E0B',         // Âmbar para metas
    },
    spacing: {
      'timeline-dot': '12px',             // Tamanho do marcador
      'timeline-line': '2px',             // Largura da linha vertical
    },
    fontFamily: {
      'clinical': ['Source Serif 4', 'serif'], // Para conteúdo clínico
    },
  },
}
```

### 7. Componentes .templ a Criar/Reutilizar

```
web/components/timeline/
├── filters.templ        // TimelineFilters (hx-trigger="change")
├── container.templ      // TimelineContainer (hx-target="#timeline-container")
├── event_card.templ     // TimelineEvent (reutilizável por tipo)
├── empty_state.templ    // TimelineEmpty
└── skeleton.templ       // Loading state por card
```

### 8. Handler Go (internal/web/handlers/timeline_handler.go)

```go
type TimelineHandler struct {
    timelineService application.TimelineService
    templates      *template.Template
}

func (h *TimelineHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
    // 1. Extrair filtros da query string
    // 2. Chamar service com TimelineConfig
    // 3. Mapear para []TimelineEvent
    // 4. Verificar HX-Request:
    //    - true: renderizar apenas #timeline-container fragment
    //    - false: renderizar página completa via ShellLayout
}
```

### 9. Validações Obrigatórias Pré-Commit

```bash
# 1. Zero inline styles
grep -o 'style="' web/components/timeline/*.templ | wc -l  # Deve ser 0

# 2. Grid responsivo
grep -q "grid-cols-1 md:grid-cols-4" web/components/timeline/container.templ

# 3. HTMX targets definidos
grep -q 'hx-target="#timeline-container"' web/components/timeline/filters.templ

# 4. Tokens do design system
grep -q "bg-timeline-session\|text-timeline-session" web/components/timeline/event_card.templ

# 5. Tipografia clínica
grep -q "font-clinical" web/components/timeline/event_card.templ
```

### 10. Critérios de Aceitação

- [ ] **CA-01**: Timeline carrega com eventos ordenados (mais recente primeiro)
- [ ] **CA-02**: Filtros atualizam conteúdo via HTMX sem recarregar página
- [ ] **CA-03**: Busca textual com debounce (500ms) funciona corretamente
- [ ] **CA-04**: Indicadores visuais por tipo de evento (cores + ícones)
- [ ] **CA-05**: Mobile: layout colapsa para 1 coluna sem quebrar
- [ ] **CA-06**: Estado vazio exibe CTA para criar primeira sessão
- [ ] **CA-07**: Loading states visíveis durante fetch de dados
- [ ] **CA-08**: Zero hardcoded values (cores, spacing, fonts via tokens)
- [ ] **CA-09**: Infinite scroll funciona (hx-trigger="revealed")
- [ ] **CA-10**: E2E audit passa sem erros SLP

---

## 📚 Documentação de Referência

| Documento | Seção | Por Que Ler |
|-----------|-------|-------------|
| `docs/PROJECT_CONTEXT.md` | Princípios de Design | Entender padrões do projeto |
| `docs/design-system.md` | Tipografia | Aplicar fontes corretas |
| `docs/architecture/standardized_layout_protocol.md` | Layout | Estrutura SLP obrigatória |
| `docs/architecture/WEB_LAYER_PATTERN.md` | HTMX | Consciência de contexto |
| `padrao_layout_grids.md` | StandardGrid | Grid 4x4 declarativo |
| `docs/requirements/req-02-02-01-linha-tempo.md` | Requisitos | Critérios de negócio |
| `docs/learnings/MASTER_LEARNINGS.md` | Anti-Padrões | O que NÃO fazer |

---

## 🚨 Anti-Padrões a Evitar

| Anti-Padrão | Por Que Evitar | Alternativa |
|-------------|----------------|-------------|
| `style="..."` inline | Falha no E2E audit | Classes CSS em style.css |
| Grid sem responsividade | Quebra em mobile | `grid-cols-1 md:grid-cols-4` |
| Entidades de domínio no template | Acoplamento | ViewModels específicos |
| Ignorar HX-Request | Quebra HTMX | Verificar header sempre |
| Hardcoded colors/spacing | Viola design system | Tokens do tailwind.config.js |
| Padding/margin externo nos widgets | Quebra grid | Gap via StandardGrid |

---

## 🧪 Validação Pós-Implementação

```bash
# 1. Build
go build ./cmd/arandu

# 2. Quality gates
./scripts/arandu_guard.sh

# 3. E2E audit
./scripts/arandu_e2e_audit.sh --routes patients

# 4. Validação visual
./scripts/arandu_visual_check.sh

# 5. Screenshot
./scripts/arandu_screenshot.sh

# 6. Verificar inline styles
grep -o 'style="' web/components/timeline/*.templ | wc -l  # Deve ser 0
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

**Instrução Final**: Esta tarefa é crítica para a experiência clínica longitudinal. Não pule validações visuais — a timeline deve ser legível e fluida em todos os dispositivos.
