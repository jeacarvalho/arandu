---

# CAP-01-02 — Registro de observações clínicas

## Identificação

**ID:** CAP-01-02 **Vision associada:** VISION-01 Registro da prática clínica 
**Nome:** Registro de observações clínicas
**Status:** ativo

---

# Propósito

Permitir que o profissional registre **percepções clínicas específicas** identificadas durante ou após a sessão terapêutica. As observações representam o olhar técnico do terapeuta sobre o material trazido pelo paciente.

Enquanto o resumo da sessão (CAP-01-01) oferece uma visão narrativa, as observações clínicas são **unidades atômicas de percepção** que servirão de base para a inteligência do sistema.

---

# Papel no sistema

As observações clínicas são componentes fundamentais da sessão. Elas transformam o relato subjetivo em dados estruturados para reflexão.

Exemplos de conteúdo para observação:

* Padrões emocionais recorrentes.


* Relatos significativos ou marcantes.


* Comportamentos não-verbais observados.


* Hipóteses clínicas iniciais.



Estrutura conceitual:

```text
Patient
   └── Session
         └── Observations (N)

```

---

# Problema que resolve

Evita que percepções sutis se percam na massa de texto de um prontuário comum. Sem o registro atômico de observações, o profissional enfrenta:

* 
**Fragmentação do conhecimento**: Dificuldade em isolar o que é observação técnica do que é relato do paciente.


* 
**Limitação de análise**: Impossibilidade de correlacionar percepções semelhantes entre diferentes sessões ou pacientes.



---

# Estrutura da observação

A observação deve ser mantida com foco no conteúdo clínico puro.

Campos da entidade:

* 
**ID**: Identificador único.


* 
**SessionID**: Vínculo com a sessão de origem.


* 
**Content**: O conteúdo da percepção clínica.


* 
**CreatedAt**: Registro temporal da observação.



---

# Fluxo funcional

1. O terapeuta acessa uma sessão existente ou em andamento.


2. Registra uma nova percepção no painel de observações.


3. O sistema armazena a observação como parte da memória longitudinal do paciente.



---

# Requisitos associados

* 
**REQ-01-02-01**: Adicionar observação clínica.


* 
**REQ-01-02-02**: Editar observação clínica.



---

# Relação com outras capabilities

Habilita diretamente as funções analíticas do sistema:

* 
**CAP-03-01 Organização de observações**: Permite classificar e categorizar o conhecimento.


* 
**CAP-04-01 Identificação de padrões**: Base para encontrar recorrências clínicas.


* 
**CAP-05-01 Assistente reflexivo**: Alimenta a IA com as percepções do terapeuta para gerar insights.



---

# Fora do escopo

* Registro de ações do terapeuta (pertence à **CAP-01-03 Intervenções**).


* Classificação diagnóstica ou tags (pertence à **CAP-03-01**).


* Análise automatizada (pertence à **CAP-09-01**).



---

# Resultado esperado

O Arandu permitirá que cada sessão seja enriquecida com camadas de percepção técnica. Ao final, o terapeuta terá um **whiteboard clínico digital** organizado, onde as observações são facilmente recuperáveis para análise de evolução.

---

# Observação estratégica

As observações são a "matéria-prima" da inteligência do Arandu. Elas não devem ser tratadas como burocracia, mas como um **exercício de escrita reflexiva**. A interface deve convidar o terapeuta a registrar *insights* no momento em que eles emergem, sem fricção visual.

---
