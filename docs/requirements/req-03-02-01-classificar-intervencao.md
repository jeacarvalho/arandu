# REQ-03-02-01 — Classificar Intervenções Terapêuticas

## Identificação
| Campo | Valor |
|-------|-------|
| **ID** | REQ-03-02-01 |
| **Vision** | VISION-03 — Organização do Conhecimento Clínico |
| **Capability** | CAP-03-02 — Organização de Intervenções Terapêuticas |
| **Status** | ✅ Implementado |
| **Prioridade** | Alta |

## Descrição
O sistema deve permitir que o terapeuta categorize intervenções terapêuticas utilizando **tags predefinidas** organizadas em tipos específicos, facilitando a organização, busca e análise de padrões das condutas técnicas aplicadas.

## Contexto
Assim como as observações clínicas (REQ-03-01-01), as intervenções terapêuticas precisam ser organizadas para permitir:
- Identificação de padrões de conduta
- Análise de eficácia por tipo de intervenção
- Busca rápida por intervenções específicas
- Geração de relatórios de prática clínica

## Funcionalidades

### Sistema de Classificação por Tags
| Tipo de Classificação | Cor | Ícone | Descrição |
|----------------------|-----|-------|-----------|
| **Técnica Cognitiva** | #7C3AED (roxo) | 🧠 | Intervenções focadas em cognição |
| **Técnica Comportamental** | #1D9E75 (verde) | 🏃 | Intervenções de modificação comportamental |
| **Técnica Emocional** | #0F6E56 (verde base) | ❤️ | Intervenções focadas em emoções |
| **Psicoeducação** | #F59E0B (âmbar) | 📚 | Orientações e informações ao paciente |
| **Exploração Narrativa** | #3B82F6 (azul) | 💬 | Técnicas de exploração de história |
| **Intervenção Corporal** | #DC2626 (vermelho) | 💓 | Técnicas que envolvem o corpo |

### Tags Predefinidas (Exemplos)
O sistema incluirá tags organizadas por tipo:

**Técnica Cognitiva:**
- Reestruturação cognitiva
- Questionamento socrático
- Identificação de distorções
- Registro de pensamentos
- Experimento comportamental cognitivo

**Técnica Comportamental:**
- Exposição gradual
- Ativação comportamental
- Treino de habilidades
- Reforço positivo
- Modelagem

**Técnica Emocional:**
- Validação emocional
- Expressão de sentimentos
- Regulação emocional
- Mindfulness emocional
- Processamento emocional

**Psicoeducação:**
- Explicação sobre transtorno
- Informações sobre medicação
- Orientação familiar
- Prevenção de recaída
- Estratégias de coping

**Exploração Narrativa:**
- Externalização
- Reautorização
- Identificação de exceções
- Perguntas circulares
- Genograma

**Intervenção Corporal:**
- Respiração diafragmática
- Relaxamento muscular
- Grounding
- Técnicas de ancoragem
- Consciência corporal

## Critérios de Aceitação

- [ ] **CA-01:** Terapeuta pode adicionar tags a uma intervenção clicando no ícone de classificação
- [ ] **CA-02:** Tags são exibidas como badges coloridos abaixo do conteúdo da intervenção
- [ ] **CA-03:** Cada tag pode ter intensidade de 1 a 5 (opcional)
- [ ] **CA-04:** Terapeuta pode remover tags individualmente
- [ ] **CA-05:** Seleção de tags funciona via HTMX sem recarregar a página
- [ ] **CA-06:** Sistema exibe tags agrupadas por tipo no seletor
- [ ] **CA-07:** Persistência das classificações no banco de dados
- [ ] **CA-08:** API REST para operações de classificação
- [ ] **CA-09:** Busca de intervenções por tipo de tag
- [ ] **CA-10:** Visualização de estatísticas de intervenções por tipo

## Arquitetura da Implementação

### Banco de Dados
```sql
-- Tabela de tags predefinidas para intervenções
CREATE TABLE intervention_tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tag_type TEXT NOT NULL, -- cognitive, behavioral, emotional, psychoeducation, narrative, body
    color TEXT NOT NULL,
    icon TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL
);

-- Tabela de relacionamento many-to-many
CREATE TABLE intervention_classifications (
    id TEXT PRIMARY KEY,
    intervention_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    intensity INTEGER CHECK (intensity >= 1 AND intensity <= 5),
    created_at DATETIME NOT NULL,
    FOREIGN KEY (intervention_id) REFERENCES interventions(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES intervention_tags(id) ON DELETE CASCADE,
    UNIQUE(intervention_id, tag_id)
);

-- Índices
CREATE INDEX idx_intervention_classifications_intervention ON intervention_classifications(intervention_id);
CREATE INDEX idx_intervention_classifications_tag ON intervention_classifications(tag_id);
CREATE INDEX idx_intervention_tags_type ON intervention_tags(tag_type);
```

### Endpoints
| Método | Rota | Descrição |
|--------|------|-----------|
| `POST` | `/interventions/{id}/classify` | Adicionar tag a uma intervenção |
| `DELETE` | `/interventions/{id}/classify/{tag_id}` | Remover tag da intervenção |
| `GET` | `/interventions/{id}/classify/edit` | Formulário de seleção de tags |
| `GET` | `/tags/interventions?type={tag_type}` | Listar tags por tipo |

### Componentes UI
1. **Tag Badge** - Badge colorido com nome da tag e intensidade
2. **Tag Selector Grid** - Grid interativo para seleção múltipla
3. **Tag List** - Lista horizontal de tags aplicadas
4. **Classification Summary** - Painel de resumo por tipo
5. **Classification Panel** - Painel lateral para navegação

## Integração
- Habilitado por: REQ-01-03-01 (Registrar Intervenção Terapêutica)
- Habilita: VISION-04 (Análise de Padrões), VISION-09 (Inteligência Clínica)

---
