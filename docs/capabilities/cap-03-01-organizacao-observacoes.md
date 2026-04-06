---
id: CAP-03-01
vision: VISION-03
status: implemented
---

# CAP-03-01 — Organização de observações clínicas

## Visão associada

VISION-03 — Organização do Conhecimento Clínico

## Descrição

Sistema de classificação e categorização de observações clínicas que permite ao terapeuta:

1. **Classificar observações com tags estruturadas**: Associar tags predefinidas a cada observação clínica, organizadas em 6 categorias (emoção, comportamento, cognição, relação, somático, contexto).

2. **Definir intensidade**: Atribuir nível de intensidade (1-5) a cada classificação, permitindo graduar a relevância ou magnitude do aspecto observado.

3. **Visualizar distribuição**: Acompanhar a distribuição de tags por tipo através de painéis de resumo e estatísticas visuais.

4. **Filtrar e buscar**: Localizar observações por tipo de classificação, facilitando a análise retrospectiva e identificação de padrões.

## Funcionalidades

| Funcionalidade | Descrição | Status |
|----------------|-----------|--------|
| Sistema de Tags | 6 tipos de classificação com 29 tags predefinidas | ✅ Implementado |
| Níveis de Intensidade | Escala 1-5 para graduação da observação | ✅ Implementado |
| Interface HTMX | Adição/remoção de tags sem reload | ✅ Implementado |
| Badges Visuais | Exibição colorida das classificações | ✅ Implementado |
| Painel de Resumo | Distribuição e estatísticas de tags | ✅ Implementado |
| Seletor Grid | Interface de seleção múltipla organizada por tipo | ✅ Implementado |

## Requisitos relacionados

- [REQ-03-01-01 — Classificar observações clínicas](./req-03-01-01-classificar-observacao.md) ✅

## Arquitetura

### Componentes

```
web/components/classification/
├── tag.templ              # Badge individual de tag
├── tag_list.templ         # Lista de tags
├── tag_selector.templ     # Seletor inline
├── tag_selector_grid.templ # Grid completo de seleção
├── panel.templ            # Painel lateral de navegação
└── summary.templ          # Resumo e estatísticas
```

### Handlers

- `ClassificationHandler` em `internal/web/handlers/classification_handler.go`

### Rotas

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/observations/{id}/classify` | Adicionar tag |
| DELETE | `/observations/{id}/classify/{tag_id}` | Remover tag |
| GET | `/observations/{id}/classify/edit` | Form de seleção |
| GET | `/tags?type={type}` | Listar por tipo |

### Modelos de Dados

```go
// Tag representa uma tag predefinida
type Tag struct {
    ID        string
    Name      string
    TagType   TagType  // emotion, behavior, cognition, relationship, somatic, context
    Color     string
    SortOrder int
}

// ObservationTag representa o relacionamento observation-tag
type ObservationTag struct {
    ID            string
    ObservationID string
    TagID         string
    Tag           *Tag
    Intensity     int  // 1-5
}
```

## Dependências

- HTMX para interatividade
- Tailwind CSS para estilização
- SQLite para persistência
- Font Awesome para ícones

## Observações

- Design system utiliza cores específicas por tipo de tag
- Todas as operações são assíncronas via HTMX
- Sistema preparado para extensão com novas tags
- Migração 0011 cria estrutura completa com seed data

---
**Status**: ✅ Implementado e testado
**Atualizado em**: 2026-04-04
