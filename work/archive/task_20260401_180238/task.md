# TASK 20260401_180238

## Status
**PRONTO_PARA_IMPLEMENTACAO**

## 🏆 REQUISITO: REQ-03-01-01 — Classificar Observações Clínicas**

Implementação: Classificação de Observações Clínicas

📝 O que construir
Sistema de classificação/tagging para observações clínicas, permitindo que o terapeuta categorize percepções por tipo (emoção, comportamento, cognição, relação) e intensidade, criando estrutura para análise de padrões futuros.

🧩 Componentes necessários (Grid 4x4)
| Widget | ColSpan | Descrição |
|--------|---------|-----------|
| ClassificationPanel | 1 | Painel lateral com tags disponíveis |
| ObservationList | 3 | Lista de observações com classificação inline |
| TagSelector | 4 | Seletor de tags para cada observação (HTMX) |
| ClassificationSummary | 4 | Resumo visual de distribuição de tags |

📋 Regras específicas
- Classificação via HTMX (hx-trigger="change" nos selects)
- Tags pré-definidas: Emoção, Comportamento, Cognição, Relação, Soma, Contexto
- Intensidade: 1-5 escala visual (cores)
- Múltiplas tags por observação
- Filtro por tag na lista de observações
- Persistência em tabela observation_tags (many-to-many)

⚠️ PADRÕES OBRIGATÓRIOS (NÃO IGNORAR)

1. Estrutura de Layout
| Componente | Uso |
|------------|-----|
| ShellLayout | Obrigatório (TopBar + Sidebar + Main Canvas) |
| StandardGrid | Obrigatório (grid 4 colunas) |
| WidgetWrapper | Obrigatório em cada widget |

2. Stack & Estilo
| Item | Regra |
|------|-------|
| Stack | Go + Templ + Tailwind + HTMX |
| Cores/Spacing | Apenas tokens do tailwind.config.js |
| Margins | Zero margins externas (gap via StandardGrid) |
| Padding | Apenas via WidgetWrapper |
| Alturas | Automáticas (content-based) |

3. HTMX
| Atributo | Uso |
|----------|-----|
| hx-post | /observations/{id}/classify |
| hx-target | #observation-{id} |
| hx-swap | outerHTML |
| hx-trigger | change (nos selects de tag) |

4. Responsividade
| Breakpoint | Comportamento |
|------------|---------------|
| Desktop | Respeitar col-span definido |
| Mobile (<768px) | Grid 1 coluna, filtros empilhados |

5. Go Structs
```go
type ObservationTag struct {
    ID            string
    ObservationID string
    TagType       string  // emotion, behavior, cognition, relationship, somatic, context
    Intensity     int     // 1-5
    CreatedAt     time.Time
}

type ClassificationData struct {
    ObservationID string
    AvailableTags []TagOption
    SelectedTags  []ObservationTag
}
```

🆕 SEÇÕES OBRIGATÓRIAS DE ESPECIFICAÇÃO TÉCNICA

6. 🆕 Tailwind Config (Tokens a Criar/Verificar)
```js
// tailwind.config.js - Adicionar/Verificar:
theme: {
  extend: {
    colors: {
      'tag-emotion': '#0F6E56',      // Verde para emoções
      'tag-behavior': '#1D9E75',     // Verde ativo para comportamentos
      'tag-cognition': '#7C3AED',    // Roxo para cognições
      'tag-relationship': '#F59E0B', // Âmbar para relações
      'tag-somatic': '#DC2626',      // Vermelho para somático
      'tag-context': '#6B7280',      // Cinza para contexto
    },
    spacing: {
      'tag-height': '28px',
    },
  },
}
```

7. 🆕 Componentes .templ a Criar/Reutilizar
```
web/components/classification/
├── panel.templ           // ClassificationPanel (sidebar)
├── tag_selector.templ    // TagSelector (inline)
├── summary.templ         // ClassificationSummary
└── observation_item.templ // Atualizar com tags
```

8. 🆕 Handler Go
```
internal/web/handlers/classification_handler.go
- POST /observations/{id}/classify
- GET /observations/{id}/classify/edit
- GET /observations/filter?tag={type}
```

9. 🆕 Validações Obrigatórias
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

10. 🆕 Critérios de Aceitação
- [ ] CA-01: Tags podem ser adicionadas/removidas via HTMX sem reload
- [ ] CA-02: Filtro por tag atualiza lista de observações
- [ ] CA-03: Intensidade visual (cores) reflete valor 1-5
- [ ] CA-04: Resumo mostra distribuição de tags por tipo
- [ ] CA-05: Mobile: seletor de tags funciona em touch
- [ ] CA-06: Zero inline styles detectados
- [ ] CA-07: E2E audit passa sem erros SLP
- [ ] CA-08: Classificação persiste após refresh de página

📚 Documentação de Referência
| Documento | Seção | Por Que Ler |
|-----------|-------|-------------|
| `docs/PROJECT_CONTEXT.md` | Princípios de Design | Entender padrões do projeto |
| `docs/design-system.md` | Tipografia | Aplicar fontes corretas |
| `docs/architecture/standardized_layout_protocol.md` | Layout | Estrutura SLP obrigatória |
| `docs/architecture/WEB_LAYER_PATTERN.md` | HTMX | Consciência de contexto |
| `padrao_layout_grids.md` | StandardGrid | Grid 4x4 declarativo |
| `docs/requirements/req-03-01-01-classificar-observacao.md` | Requisitos | Critérios de negócio |
| `docs/learnings/MASTER_LEARNINGS.md` | Anti-Padrões | O que NÃO fazer |

🚨 Anti-Padrões a Evitar
| Anti-Padrão | Por Que Evitar | Alternativa |
|-------------|----------------|-------------|
| `style="..."` inline | Falha no E2E audit | Classes CSS em style.css |
| Grid sem responsividade | Quebra em mobile | `grid-cols-1 md:grid-cols-4` |
| Tags hardcoded no HTML | Dificulta manutenção | Tabela SQL + seed data |
| Ignorar HX-Request | Quebra HTMX | Verificar header sempre |
| Hardcoded colors | Viola design system | Tokens do tailwind.config.js |

🧪 Validação Pós-Implementação
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

📝 Checklist de Conclusão
- [ ] Li toda documentação obrigatória
- [ ] Implementei conforme especificação
- [ ] Rodei `go build` sem erros
- [ ] Rodei `./scripts/arandu_guard.sh` e passou
- [ ] Rodei E2E audit e passou
- [ ] Validação visual manual OK (desktop + mobile)
- [ ] Screenshot gerado e revisado
- [ ] Zero inline styles detectados
- [ ] Documentei aprendizados em `docs/learnings/MASTER_LEARNINGS.md`

**Instrução Final**: Esta classificação é a base para toda inteligência futura do sistema. Não pule validações — a estrutura de tags deve ser consistente e extensível.
