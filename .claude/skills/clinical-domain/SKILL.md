---
name: clinical-domain
description: >
  Domínio clínico de uma plataforma de psicologia para múltiplos psicólogos independentes.
  Use esta skill SEMPRE que o usuário trabalhar com qualquer conceito do sistema: paciente,
  sessão, agenda, prontuário, percurso terapêutico, anamnese, evolução, vínculo terapêutico,
  encaminhamento, ou qualquer feature de IA clínica. Também dispara ao escrever requisitos,
  modelar agregados DDD, criar prompts de IA, nomear variáveis/funções/tabelas relacionadas
  ao negócio, ou perguntar "como modelar X?", "o que é Y no sistema?", "como chamar esse
  campo?". Esta skill define a Ubiquitous Language do sistema — todo artefato gerado
  (código, requisito, prompt, migration) deve usar os termos aqui definidos.
---

# Domínio Clínico — Plataforma de Psicologia

Ubiquitous Language e regras de negócio da plataforma. Todo código, requisito, prompt
e migration deve usar exatamente os termos definidos aqui.

---

## Modelo de negócio

Plataforma **multi-tenant** onde cada **Psicólogo** é um tenant independente com seus
próprios **Pacientes**, **Agenda** e **Prontuários**. Pacientes não são compartilhados
entre psicólogos. A plataforma oferece features de IA que analisam os dados *dentro*
do tenant, nunca cruzando dados entre psicólogos.

---

## Glossário — Ubiquitous Language

### Entidades principais

**Psicólogo** (`Psychologist`)
Profissional registrado no CFP, titular de um espaço na plataforma. Possui CRP (Conselho
Regional de Psicologia) que identifica o profissional e a região. Pode atender presencial
e/ou online (telepsicologia). É o único que acessa os dados de seus pacientes.
- Campos obrigatórios: nome completo, CRP, e-mail profissional, modalidade de atendimento
- Nunca chamado de: *usuário*, *terapeuta* (termo genérico), *médico*, *doutor*

**Paciente** (`Patient`)
Pessoa que está em acompanhamento psicológico com um Psicólogo da plataforma. Pode ser
adulto (18+) ou menor de idade — neste caso requer responsável legal com dados separados.
- Identificado internamente por UUID opaco — nunca por nome em contextos de IA
- Nunca chamado de: *cliente*, *usuário*, *caso*
- Menor de idade: sempre tem `LegalGuardian` associado

**Responsável Legal** (`LegalGuardian`)
Obrigatório quando Paciente é menor de idade. Assina o Termo de Consentimento Livre e
Esclarecido (TCLE) e autoriza o tratamento de dados.

---

### Agenda e Sessões

**Agenda** (`Schedule`)
Conjunto de disponibilidades e compromissos do Psicólogo. Cada Psicólogo tem exatamente
uma Agenda. Não é um simples calendário — contém regras de recorrência, bloqueios e
políticas de cancelamento.

**Disponibilidade** (`Availability`)
Janela de tempo em que o Psicólogo aceita agendamentos. Define dias da semana, horários,
duração padrão de sessão e modalidade (presencial/online).

**Agendamento** (`Appointment`)
Reserva de um horário entre Psicólogo e Paciente. Tem ciclo de vida próprio:

```
REQUESTED → CONFIRMED → COMPLETED
                      ↘ CANCELLED
                      ↘ NO_SHOW
```

- `REQUESTED`: paciente solicitou, aguarda confirmação
- `CONFIRMED`: psicólogo confirmou
- `COMPLETED`: sessão ocorreu e foi registrada
- `CANCELLED`: cancelado (por qualquer parte) com motivo e responsável
- `NO_SHOW`: paciente não compareceu sem aviso

**Sessão** (`Session`)
Registro clínico do atendimento realizado. Criada a partir de um `Appointment` confirmado
quando o psicólogo registra a evolução. É o núcleo do prontuário.
- Uma Sessão sempre deriva de um Appointment (nunca criada avulsamente)
- Contém: data/hora, duração real, modalidade efetiva, formulário de evolução
- Nunca chamada de: *consulta*, *atendimento*, *encontro*

**Formulário de Evolução** (`SessionForm`)
Estrutura de campos preenchida pelo Psicólogo após cada Sessão. Campos:
- `chief_complaint`: queixa principal relatada na sessão
- `observations`: observações clínicas do psicólogo (texto livre)
- `emotional_state`: estado emocional observado (escala + descritivo)
- `themes`: temas abordados (lista de tags controladas + livres)
- `techniques_used`: técnicas terapêuticas aplicadas
- `homework`: tarefa/reflexão proposta para o paciente
- `risk_indicators`: indicadores de risco observados (ideação, etc.) — campo crítico
- `next_session_focus`: foco planejado para próxima sessão

---

### Prontuário

**Prontuário** (`MedicalRecord`)
Conjunto completo e cronológico de todas as Sessões de um Paciente com um Psicólogo.
É um documento legal sob responsabilidade do Psicólogo.
- Inclui: Anamnese inicial + todas as Sessões + documentos anexos
- Deve ser mantido por mínimo 5 anos após último atendimento (CFP)
- Nunca pode ser excluído — apenas arquivado

**Anamnese** (`Intake`)
Formulário de avaliação inicial preenchido na primeira sessão ou antes dela. Estabelece
o histórico do paciente. Campos principais:
- `chief_complaint`: motivo da busca por terapia
- `clinical_history`: histórico clínico relevante
- `family_history`: histórico familiar relevante
- `previous_treatments`: tratamentos anteriores (psicológicos, psiquiátricos)
- `current_medications`: medicações em uso
- `life_context`: contexto de vida atual (trabalho, relacionamentos, moradia)
- `therapeutic_goals`: objetivos terapêuticos iniciais

**Evolução** (`Progress`)
Visão longitudinal do paciente — gerada pela IA a partir das Sessões. Não é um campo
preenchido pelo psicólogo, é uma *análise computada*. Sempre apresentada como sugestão,
nunca como diagnóstico.

---

### Percurso Terapêutico

**Percurso Terapêutico** (`TherapeuticJourney`)
Conceito central das features de IA. Representa a narrativa de evolução do paciente
ao longo do tempo, identificada pela análise de múltiplas Sessões. Composto por:
- `recurring_themes`: temas que aparecem em múltiplas sessões
- `emotional_trajectory`: tendência do estado emocional ao longo do tempo
- `milestone_sessions`: sessões identificadas como pontos de virada
- `suggested_paths`: caminhos terapêuticos sugeridos pela IA (sempre como hipóteses)

**Padrão Recorrente** (`RecurringPattern`)
Elemento identificado pela IA que aparece em 3 ou mais sessões. Pode ser:
- Tema (`ThemePattern`): assunto recorrente nas sessões
- Emocional (`EmotionalPattern`): estado emocional recorrente
- Comportamental (`BehavioralPattern`): comportamento relatado recorrentemente
- Relacional (`RelationalPattern`): padrão em relacionamentos

**Indicador de Risco** (`RiskIndicator`)
Campo crítico do SessionForm. Qualquer registro aqui deve:
1. Ser armazenado com timestamp e nunca editado (apenas adicionado)
2. Gerar alerta visível no prontuário
3. **Nunca** ser processado por IA externa sem anonimização total
4. Ser excluído do payload de qualquer chamada a LLM

---

### Abordagens Terapêuticas

O sistema deve conhecer as principais abordagens para categorizar técnicas e sugestões:

| Sigla | Nome | Uso típico |
|-------|------|------------|
| TCC | Terapia Cognitivo-Comportamental | Padrões de pensamento, ansiedade, depressão |
| ACT | Terapia de Aceitação e Compromisso | Valores, flexibilidade psicológica |
| TFE | Terapia Focada nas Emoções | Regulação emocional, apego |
| Psicanálise | Psicanálise/Psicodinâmica | Inconsciente, vínculos |
| Humanista | Abordagem Humanista/Gestalt | Autorrealização, presente |
| EMDR | Eye Movement Desensitization | Trauma |
| DBT | Terapia Dialético-Comportamental | Regulação emocional intensa |

O Psicólogo registra sua(s) abordagem(ns) no perfil. A IA usa isso para contextualizar
sugestões de caminhos terapêuticos.

---

## Regras de negócio críticas

### Isolamento de dados (multi-tenant)
- Paciente pertence a exatamente **um** Psicólogo — nunca transferido
- Consultas de IA são sempre **scoped** ao tenant do Psicólogo
- Nenhuma análise cruza dados entre Psicólogos, nem anonimizados

### Menores de idade
- Todo acesso ao prontuário de menor requer registro de consentimento do LegalGuardian
- Sessões com menor: campo `guardian_present` obrigatório no SessionForm
- Conteúdo sensível de menor tem camada extra de proteção na anonimização

### Ciclo de vida do Prontuário
- Criado na primeira Sessão confirmada
- **Nunca deletado** — arquivado quando paciente encerra acompanhamento
- Retenção mínima: 5 anos após último atendimento (adultos), 5 anos após maioridade (menores)

### Indicadores de Risco
- Presença de `risk_indicators` bloqueia o envio dos dados da sessão para qualquer IA externa
- Psicólogo recebe aviso proeminente ao abrir prontuário com risco registrado
- Campo imutável após salvo (append-only)

### O sistema SUGERE — nunca DIAGNOSTICA
- Toda saída de IA deve ser apresentada como **sugestão** ou **hipótese**
- Proibido usar linguagem de certeza: "o paciente tem", "diagnóstico indica"
- Correto: "padrão observado sugere", "pode indicar", "considerar explorar"
- O Psicólogo é sempre o decisor final — a IA é uma ferramenta de apoio

---

## Entidades adicionais (Arandu atual)

**Contexto Biopsicossocial** (`PatientContext`)
Dados de identidade social do Paciente. Campos altamente sensíveis — proteção reforçada.
- `ethnicity`: etnia (padrão IBGE / saúde)
- `gender_identity`: identidade de gênero
- `sexual_orientation`: orientação sexual
- `occupation`: ocupação profissional
- `education_level`: nível de escolaridade
> ⚠️ Todos esses campos são **Tier 1-Plus** em privacidade — mais sensíveis que dados
> de saúde comuns. Nunca enviados para IA. Nunca exibidos em logs.

**Observação** (`Observation`)
Unidade atômica de registro clínico dentro de uma Sessão. Texto livre do psicólogo.
Indexada via SQLite FTS5 para busca instantânea. Campo `content` nunca enviado para IA.

**Intervenção** (`Intervention`)
Ação terapêutica aplicada durante uma Sessão. Registra técnica usada e resposta observada.

**Medicação** (`Medication`)
Histórico farmacológico do Paciente. Campos: nome, dosagem, frequência, prescritor,
status (`active` | `suspended` | `finished`), período de uso.

**Sinais Vitais e Hábitos** (`Vitals`)
Registro periódico: horas de sono, apetite (escala 1-10), peso, atividade física, notas.
Permite análise de correlação entre bem-estar físico e evolução clínica.

**Timeline** (`Timeline`)
Read model longitudinal — agrega Sessões, Observações e Intervenções em ordem cronológica
para visualização do percurso terapêutico. Gerado sob demanda, não persistido.

---

## Mapeamento DDD → Go (Arandu real)

| Conceito de Domínio | Tipo Go | Pacote no Arandu |
|---|---|---|
| Patient | Aggregate Root | `internal/domain/patient` |
| PatientContext | Entity (filho de Patient) | `internal/domain/patient` |
| Medication | Entity (filho de Patient) | `internal/domain/patient` |
| Vitals | Value Object | `internal/domain/patient` |
| Session | Aggregate Root | `internal/domain/session` |
| Observation | Entity (filho de Session) | `internal/domain/observation` |
| Intervention | Entity (filho de Session) | `internal/domain/intervention` |
| Timeline | Read Model | `internal/domain/timeline` |
| TherapeuticJourney | Read Model (IA) | `internal/domain/timeline` |
| RecurringPattern | Value Object | `internal/domain/timeline` |

> Nota: No Arandu não há `Psychologist` como entidade no tenant DB — o psicólogo
> existe no Control Plane. Dentro do tenant DB, o contexto do psicólogo é implícito
> (é o dono do banco).

---

## Referências
- `references/regulatory.md` — CFP Resolução 11/2018, LGPD aplicada, prontuário eletrônico
- `references/ai-features.md` — Especificação das features de IA: inputs, outputs, restrições
- `references/requirements-template.md` — Template para escrever requisitos neste domínio
