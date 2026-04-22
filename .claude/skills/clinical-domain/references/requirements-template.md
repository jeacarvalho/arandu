# Template de Requisitos — Plataforma de Psicologia

## Como usar este template

Todo requisito do sistema deve seguir este formato. O objetivo é garantir que:
- Use a Ubiquitous Language do domínio (termos do clinical-domain)
- Explicite as restrições de privacidade e conformidade
- Deixe claro se envolve IA e quais dados serão usados
- Defina os critérios de aceite de forma verificável

---

## Template padrão

```markdown
## REQ-[NÚMERO]: [Título usando Ubiquitous Language]

**Contexto**
[Quem usa, em qual situação, qual problema resolve]

**Atores**
- Primário: [Psicólogo | Paciente | Sistema | IA]
- Secundário: [se houver]

**Descrição**
Como [ator], quero [ação no vocabulário do domínio] para [benefício clínico].

**Regras de negócio**
- RN01: [regra específica do domínio]
- RN02: [...]

**Restrições de privacidade** *(obrigatório se tocar dados de paciente)*
- [ ] Dados pessoais envolvidos: [listar campos]
- [ ] Base legal LGPD: [Consentimento | Tutela da saúde]
- [ ] Requer TCLE: [sim | não | já obtido]
- [ ] Dado sensível de menor: [sim | não]

**Restrições de IA** *(obrigatório se feature usa IA)*
- [ ] Opt-in do psicólogo: [sim | automático]
- [ ] Campos enviados ao LLM: [listar explicitamente]
- [ ] Campos excluídos: [listar — deve incluir todos os PII]
- [ ] risk_indicators presentes bloqueiam: [sim | não aplicável]
- [ ] Output apresentado como: [sugestão | hipótese | análise]

**Critérios de aceite**
- CA01: Dado [contexto], quando [ação], então [resultado esperado]
- CA02: [...]
- CA03 (negativo): Dado [contexto proibido], quando [ação], então [sistema rejeita/bloqueia]

**Conformidade**
- [ ] CFP Resolução 11/2018
- [ ] LGPD Art. [número]
- [ ] Prontuário eletrônico CFP
- [ ] Sigilo profissional
```

---

## Exemplos preenchidos

### Exemplo 1 — Feature sem IA

```markdown
## REQ-012: Registrar Evolução de Sessão

**Contexto**
Após realizar uma Sessão, o Psicólogo precisa registrar as observações clínicas
enquanto ainda estão frescas. O registro compõe o Prontuário do Paciente.

**Atores**
- Primário: Psicólogo

**Descrição**
Como Psicólogo, quero preencher o Formulário de Evolução após uma Sessão Confirmada
para manter o Prontuário atualizado e fundamentar análises futuras.

**Regras de negócio**
- RN01: Formulário de Evolução só pode ser criado para Sessões com status COMPLETED
- RN02: Uma Sessão tem exatamente um Formulário de Evolução — não pode ser recriado
- RN03: Após salvo, campos podem ser editados em até 48h; após isso, apenas adicionados (append)
- RN04: Campo risk_indicators, se preenchido, gera alerta imediato no Prontuário
- RN05: Psicólogo só acessa Sessões de seus próprios Pacientes

**Restrições de privacidade**
- [x] Dados pessoais envolvidos: observations, chief_complaint (texto livre — PII implícito)
- [x] Base legal LGPD: Tutela da saúde (Art. 11, II, f)
- [x] Requer TCLE: já obtido no início do acompanhamento
- [ ] Dado sensível de menor: depende do paciente

**Critérios de aceite**
- CA01: Dado Sessão com status COMPLETED, quando Psicólogo preenche e salva o Formulário, então Sessão passa a ter Formulário associado e Prontuário é atualizado
- CA02: Dado Formulário salvo há mais de 48h, quando Psicólogo tenta editar campo existente, então sistema rejeita e exibe mensagem de política de imutabilidade
- CA03 (negativo): Dado Sessão com status CONFIRMED (não realizada), quando Psicólogo tenta criar Formulário, então sistema rejeita com erro de negócio

**Conformidade**
- [x] Prontuário eletrônico CFP
- [x] LGPD Art. 11 (dados sensíveis de saúde)
- [x] Sigilo profissional
```

### Exemplo 2 — Feature com IA

```markdown
## REQ-031: Detectar Padrões Recorrentes no Percurso Terapêutico

**Contexto**
Após acumular múltiplas Sessões, o Psicólogo quer identificar temas e estados
emocionais recorrentes para orientar o planejamento terapêutico.

**Atores**
- Primário: Psicólogo
- Secundário: IA (serviço externo de LLM)

**Descrição**
Como Psicólogo, quero solicitar a detecção de Padrões Recorrentes no Percurso
Terapêutico de um Paciente para identificar temas e trajetórias que merecem atenção.

**Regras de negócio**
- RN01: Requer mínimo de 3 Sessões com Formulário de Evolução preenchido
- RN02: Psicólogo deve ativar explicitamente a análise (opt-in por paciente)
- RN03: Se houver RiskIndicator nas últimas 3 sessões, análise é bloqueada
- RN04: Resultado é salvo como TherapeuticJourney — não substitui o Prontuário

**Restrições de IA**
- [x] Opt-in do psicólogo: sim — botão explícito "Analisar padrões"
- [x] Campos enviados ao LLM: themes[], emotional_state (escala), techniques_used[], next_session_focus, data da sessão
- [x] Campos excluídos: observations, chief_complaint, homework, risk_indicators, name, cpf, birth_date, phone, email, address, psychologist.name, psychologist.crp
- [x] risk_indicators presentes bloqueiam: sim (últimas 3 sessões)
- [x] Output apresentado como: sugestão com disclaimer obrigatório

**Critérios de aceite**
- CA01: Dado Paciente com 5+ Sessões sem RiskIndicators recentes, quando Psicólogo solicita análise, então sistema anonimiza dados, chama LLM, e apresenta RecurringPatterns com disclaimer
- CA02: Dado payload enviado ao LLM, então não deve conter nenhum dos campos da lista de exclusão (verificável em teste de integração)
- CA03 (negativo): Dado RiskIndicator registrado na última sessão, quando Psicólogo tenta solicitar análise, então sistema bloqueia com mensagem orientando avaliação de risco primeiro

**Conformidade**
- [x] LGPD Art. 11 (dados sensíveis — anonimização antes de envio externo)
- [x] CFP — sigilo profissional (dados não identificáveis saem do sistema)
- [x] Opt-in documentado no audit trail
```
