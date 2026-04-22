# Features de IA — Especificação

## Princípios gerais (não negociáveis)

1. **IA como copiloto, não piloto**: toda saída é sugestão; psicólogo decide
2. **Anonimização obrigatória**: nenhum dado identificável vai para LLM externo
3. **Opt-in explícito**: psicólogo ativa análise de IA por paciente, não é automático
4. **Transparência**: psicólogo sempre sabe quando IA foi usada e em quais dados
5. **Auditável**: cada chamada de IA é logada com: timestamp, dados enviados (hash), modelo usado, output retornado

---

## Feature 1 — Detecção de Padrões Recorrentes

**O que faz:** analisa múltiplas sessões de um paciente e identifica temas, estados
emocionais e comportamentos que aparecem com frequência.

**Trigger:** manual (psicólogo clica em "Analisar padrões") ou automático após N sessões
(configurável, padrão: 5 sessões).

**Input para a IA:**
```
- themes[] de cada SessionForm (tags, não texto livre)
- emotional_state de cada SessionForm (escala numérica + categoria)
- techniques_used[] de cada SessionForm
- next_session_focus de cada SessionForm
- data da sessão (sem hora exata)
- NÃO inclui: observations, chief_complaint, homework, risk_indicators
- NÃO inclui: nome, CPF, data de nascimento, qualquer PII
- Identificador do paciente: token opaco rotativo (muda a cada análise)
```

**Output esperado:**
```json
{
  "recurring_themes": [
    { "theme": "relações familiares", "frequency": 7, "first_seen": "sessão 2", "last_seen": "sessão 9" }
  ],
  "emotional_trajectory": "tendência de melhora em regulação emocional nas últimas 4 sessões",
  "behavioral_patterns": [...],
  "confidence": "baixa|média|alta",
  "disclaimer": "Análise baseada em X sessões. Sugestão para reflexão clínica — não substitui avaliação profissional."
}
```

---

## Feature 2 — Sugestão de Caminhos Terapêuticos

**O que faz:** com base nos padrões identificados e na abordagem do psicólogo, sugere
possíveis direções terapêuticas a explorar.

**Input para a IA:**
```
- RecurringPatterns identificados (Feature 1)
- therapeutic_approach do Psicólogo (TCC, ACT, etc.)
- therapeutic_goals da Anamnese (anonimizado: categorias, não texto livre)
- número de sessões realizadas
- NÃO inclui nenhum dado identificável
```

**Output esperado:**
```
Sugestões apresentadas como hipóteses exploratórias:
- "Dado o padrão de [tema], pode ser útil explorar [técnica] alinhada à abordagem [X]"
- Sempre com ressalva de que o psicólogo avalia a pertinência
- Máximo 3 sugestões por análise
```

**Restrição crítica:** se houver `risk_indicators` registrados nas últimas 3 sessões,
esta feature é bloqueada automaticamente — psicólogo deve avaliar o risco antes de
pensar em direções terapêuticas.

---

## Feature 3 — Resumo de Sessão

**O que faz:** gera um resumo estruturado dos campos do SessionForm para facilitar
revisão rápida pelo psicólogo antes da próxima sessão.

**Input para a IA:**
```
- themes[] da sessão
- emotional_state da sessão
- techniques_used[]
- next_session_focus
- NÃO inclui: observations (texto livre), chief_complaint, homework
- NÃO inclui nenhum PII
```

**Output esperado:**
```
Resumo em 3-5 bullet points:
- Temas centrais da sessão
- Estado emocional observado
- Técnicas utilizadas
- Foco para próxima sessão
```

**Nota:** observations e chief_complaint são texto livre — têm alta probabilidade de
conter PII mesmo após anonimização automática. Por isso ficam fora do payload.

---

## Feature 4 — Análise de Progresso

**O que faz:** visão longitudinal da evolução do paciente, comparando estado emocional
e temas ao longo do tempo.

**Input para a IA:**
```
- Série temporal de emotional_state (apenas escala numérica)
- Série temporal de themes[] (tags)
- Marcos: número de sessões, tempo total de acompanhamento
- NÃO inclui texto livre de nenhum campo
- NÃO inclui PII
```

**Output esperado:**
```
- Gráfico de tendência emocional (dados processados localmente, não pela IA)
- Narrativa textual: "Nas últimas X sessões, observa-se..."
- Temas que aumentaram/diminuíram em frequência
- Disclaimer obrigatório
```

---

## Campos SEMPRE excluídos de qualquer payload de IA

```
❌ patient.name
❌ patient.cpf
❌ patient.birth_date
❌ patient.phone
❌ patient.email
❌ patient.address
❌ legal_guardian.*  (todos os campos)
❌ session_form.observations  (texto livre — risco de PII implícito)
❌ session_form.chief_complaint  (texto livre)
❌ session_form.homework  (texto livre)
❌ session_form.risk_indicators  (sempre excluído, sem exceção)
❌ psychologist.name
❌ psychologist.crp
❌ psychologist.email
```
