---
id: REQ-03-01-01
vision: VISION-03
capability: CAP-03-01
status: implemented
---

# REQ-03-01-01 — Classificar observações clínicas

## Visão associada

VISION-03 — Organização do Conhecimento Clínico

## Capability associada

CAP-03-01 — Organização de observações clínicas

## Descrição

O sistema permite que o terapeuta categorize observações clínicas utilizando tags predefinidas organizadas em 6 tipos: Emoção, Comportamento, Cognição, Relação, Somático e Contexto. Cada tag pode ter uma intensidade de 1 a 5. As classificações são persistidas em banco de dados e podem ser adicionadas/removidas via interface HTMX sem recarregar a página.

### Tipos de Classificação

| Tipo | Cor | Ícone | Descrição |
|------|-----|-------|-----------|
| Emoção | #0F6E56 (verde) | ❤️ | Estados emocionais do paciente |
| Comportamento | #1D9E75 (verde claro) | 🏃 | Padrões comportamentais observados |
| Cognição | #7C3AED (roxo) | 🧠 | Processos cognitivos e pensamentos |
| Relação | #F59E0B (âmbar) | 👥 | Dinâmicas relacionais |
| Somático | #DC2626 (vermelho) | 💓 | Manifestações físicas/corporais |
| Contexto | #6B7280 (cinza) | 🌍 | Fatores contextuais e situacionais |

### Tags Predefinidas

O sistema inclui 29 tags predefinidas distribuídas entre os tipos:
- **Emoção**: Ansiedade, Tristeza, Raiva, Alegria, Medo, Frustração
- **Comportamento**: Evitação, Confronto, Isolamento, Hiperatividade, Impassividade
- **Cognição**: Pensamento catastrófico, Perfeccionismo, Ruminar, Distorção cognitiva, Insight
- **Relação**: Conflito familiar, Dificuldade social, Limiares, Vínculo terapêutico
- **Somático**: Tensão muscular, Insônia, Sintomas físicos, Mobilidade
- **Contexto**: Evento recente, Transição de vida, Estresse ocupacional, Crise

### Endpoints

- `POST /observations/{id}/classify` - Adicionar tag a uma observação
- `DELETE /observations/{id}/classify/{tag_id}` - Remover tag da observação
- `GET /observations/{id}/classify/edit` - Formulário de seleção de tags
- `GET /tags?type={tag_type}` - Listar tags por tipo

### Componentes UI

- **Tag Badge**: Badge colorido com nome da tag e indicador de intensidade
- **Tag Selector Grid**: Grid interativo para seleção múltipla de tags com slider de intensidade
- **Tag List**: Lista horizontal de tags aplicadas com botão de remoção
- **Classification Summary**: Painel de resumo mostrando distribuição de tags por tipo
- **Classification Panel**: Painel lateral para navegação e seleção de tags

## Critérios de aceitação

- [x] **CA-01**: Terapeuta pode adicionar tags a uma observação clicando no ícone de tags
- [x] **CA-02**: Tags são exibidas como badges coloridos abaixo do conteúdo da observação
- [x] **CA-03**: Cada tag pode ter intensidade de 1 a 5 visualmente representada
- [x] **CA-04**: Terapeuta pode remover tags individualmente
- [x] **CA-05**: Seleção de tags funciona via HTMX sem recarregar a página
- [x] **CA-06**: Sistema exibe tags agrupadas por tipo no seletor
- [x] **CA-07**: Persistência das classificações no banco de dados
- [x] **CA-08**: API REST para operações de classificação

## Implementação

### Arquivos Criados/Modificados

- `internal/domain/observation/tag.go` - Modelos e tipos de tags
- `internal/web/handlers/classification_handler.go` - Handler HTTP
- `web/components/classification/*.templ` - Componentes UI
- `internal/infrastructure/repository/sqlite/observation_repository.go` - Repository
- `internal/infrastructure/repository/sqlite/migrations/0011_add_observation_tags.*.sql` - Migrações
- `cmd/arandu/main.go` - Registro de rotas

### Tabelas do Banco de Dados

```sql
-- Tabela de tags predefinidas
CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tag_type TEXT NOT NULL,
    color TEXT NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL
);

-- Tabela de relacionamento many-to-many
CREATE TABLE observation_tags (
    id TEXT PRIMARY KEY,
    observation_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    intensity INTEGER CHECK (intensity >= 1 AND intensity <= 5),
    created_at DATETIME NOT NULL,
    FOREIGN KEY (observation_id) REFERENCES observations(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE(observation_id, tag_id)
);
```

## Observações

- Implementação baseada em HTMX para atualizações assíncronas
- Design responsivo compatível com mobile e desktop
- Cores definidas via Tailwind CSS seguindo o design system
- Ícones utilizando Font Awesome
- Sistema extensível: novas tags podem ser adicionadas via migrações

---
**Implementado em**: 2026-04-04
**Status**: ✅ Completo
