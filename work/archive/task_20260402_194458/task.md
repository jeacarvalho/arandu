Implementação: Detectar Padrões Clínicos (Theme Cloud)

📝 O que construir
Sistema de análise de frequência de termos no prontuário do paciente, identificando temas recorrentes e permitindo navegação contextual para eventos onde cada tema aparece.

🧩 Componentes necessários (Grid 4x4)
| Widget | ColSpan | Descrição |
|--------|---------|-----------|
| ThemeCloud | 2 | Nuvem de temas com pesos tipográficos |
| ThemeFilters | 2 | Filtros por período (3/6/12 meses, todo histórico) |
| TimelineFiltered | 4 | Timeline filtrada por tema selecionado |
| ThemeStats | 4 | Estatísticas de frequência (top 10 termos) |

📋 Regras específicas
- Processamento via FTS5 (observations_fts + interventions_fts)
- Filtro temporal: 3 meses, 6 meses, 1 ano, todo histórico
- Click no tema → HTMX filtra timeline para eventos com aquele termo
- Stop words em português devem ser filtradas (eu, ele, para, com, etc.)
- Máximo 50 termos exibidos, ordenados por frequência
- Performance: < 500ms mesmo para pacientes com 100+ sessões

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
| hx-get | /patients/{id}/analysis/themes |
| hx-target | #theme-cloud-container |
| hx-trigger | click (nos temas) |
| Loading | Skeleton durante processamento |

4. Responsividade
| Breakpoint | Comportamento |
|------------|---------------|
| Desktop | Grid 4 colunas, nuvem completa |
| Mobile | Grid 1 coluna, nuvem simplificada |

5. Go Structs
```go
type ThemeAnalysisData struct {
    PatientID    string
    Timeframe    string        // 3m, 6m, 12m, all
    Themes       []ThemeTerm   // Termos com frequência
    TotalTerms   int
    GeneratedAt  time.Time
}

type ThemeTerm struct {
    Term      string
    Frequency int
    Weight    int  // 1-5 para styling
}
```

🆕 SEÇÕES OBRIGATÓRIAS DE ESPECIFICAÇÃO TÉCNICA

6. 🆕 Tailwind Config (Tokens a Criar/Usar)
```js
// tailwind.config.js - Adicionar/Verificar:
theme: {
  extend: {
    colors: {
      'theme-cloud-bg': '#FFFBEB',  // Fundo âmbar suave
      'theme-highlight': '#F59E0B', // Destaque para termos
    },
    fontFamily: {
      'clinical': ['Source Serif 4', 'serif'], // Para termos
    },
    fontSize: {
      'theme-1': '0.875rem',  // Menos frequente
      'theme-2': '1rem',
      'theme-3': '1.25rem',
      'theme-4': '1.5rem',
      'theme-5': '1.875rem',  // Mais frequente
    },
  },
}
```

7. 🆕 Componentes .templ a Criar/Reutilizar
```
web/components/analysis/
├── theme_cloud.templ      // Nuvem de temas
├── theme_filters.templ    // Filtros de período
├── theme_stats.templ      // Estatísticas
└── timeline_filtered.templ // Timeline filtrada
```

8. 🆕 Handler Go
```
internal/web/handlers/analysis_handler.go
- GET /patients/{id}/analysis/themes
- GET /patients/{id}/analysis/themes/{term}
```

9. 🆕 Validações Obrigatórias
```bash
# 1. Zero inline styles
grep -o 'style="' web/components/analysis/*.templ | wc -l  # Deve ser 0

# 2. FTS5 sendo usado
grep -q "observations_fts\|interventions_fts" internal/web/handlers/analysis_handler.go

# 3. Performance < 500ms
# Testar com paciente de 100+ sessões

# 4. E2E audit
./scripts/arandu_e2e_audit.sh --routes patients
```

10. 🆕 Critérios de Aceitação
- [x] CA-01: Nuvem de temas carrega em < 500ms
- [x] CA-02: Stop words em português são filtradas
- [x] CA-03: Click no tema filtra timeline via HTMX
- [x] CA-04: Filtros de período funcionam corretamente
- [x] CA-05: Mobile: nuvem simplificada mas funcional
- [x] CA-06: Zero inline styles detectados
- [x] CA-07: E2E audit passa sem erros SLP
- [x] CA-08: Termos mais frequentes têm peso visual maior

📚 Documentação de Referência
| Documento | Seção | Por Que Ler |
|-----------|-------|-------------|
| `docs/requirements/req-04-01-01-detectar-padroes.md` | Requisitos | Critérios de negócio |
| `docs/architecture/WEB_LAYER_PATTERN.md` | HTMX | Padrões de swap |
| `learnings/SQLITE_BEST_PRACTICES.md` | FTS5 | Query otimizada |
| `design-system.md` | Tipografia | Fontes clínicas |
```