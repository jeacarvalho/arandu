# CAP-01-01 — Registro de sessões

## Identificação

**ID:** CAP-01-01
**Vision associada:** VISION-01 Registro da prática clínica
**Nome:** Registro de sessões
**Status:** ativo

---

# Propósito

Permitir que o profissional registre **sessões terapêuticas realizadas com pacientes**, preservando o conteúdo essencial de cada encontro clínico.

A sessão é o **elemento central da prática clínica registrada no Arandu**.

É dentro da sessão que o terapeuta registra:

```text
observações clínicas
intervenções terapêuticas
reflexões sobre o processo terapêutico
```

---

# Papel no sistema

A sessão representa um **encontro terapêutico entre profissional e paciente**.

Ela estabelece o contexto temporal e clínico onde ocorrem:

```text
observações clínicas
intervenções terapêuticas
análise da evolução do paciente
```

Estrutura conceitual:

```text
Patient
   └── Session
         ├── Observations
         └── Interventions
```

Cada sessão pertence a **um único paciente**.

Um paciente pode possuir **múltiplas sessões**.

---

# Problema que resolve

Sem um registro estruturado de sessões, o profissional depende de:

```text
memória pessoal
anotações dispersas
cadernos físicos
documentos isolados
```

Isso gera dificuldades como:

```text
perda de informações clínicas importantes
dificuldade de recuperar evolução terapêutica
baixa capacidade de identificar padrões ao longo do tempo
```

O registro estruturado de sessões cria a base para:

```text
memória clínica longitudinal
organização do conhecimento terapêutico
descoberta de padrões clínicos
assistência reflexiva por IA
```

---

# Estrutura conceitual da sessão

Neste momento do projeto a sessão será mantida **simples e focada**.

Campos iniciais da sessão:

```text
ID
PatientID
Date
Summary
CreatedAt
UpdatedAt
```

Descrição dos campos:

---

### ID

Identificador único da sessão.

---

### PatientID

Identificador do paciente ao qual a sessão pertence.

---

### Date

Data em que a sessão ocorreu.

---

### Summary

Resumo livre da sessão.

Esse campo permite registrar rapidamente o conteúdo essencial do encontro terapêutico.

---

### CreatedAt

Data de criação do registro da sessão.

---

### UpdatedAt

Data da última modificação da sessão.

---

# Fluxo funcional

Fluxo típico de uso:

```text
Terapeuta abre página do paciente
↓
Clica "Nova sessão"
↓
Preenche informações da sessão
↓
Salva sessão
↓
Sessão passa a fazer parte do histórico do paciente
```

---

# Interface esperada

Exemplo simplificado:

```text
Nova sessão — Maria

Data
[ 20/03/2026 ]

Resumo da sessão
[________________________________]
[________________________________]

[ Salvar sessão ]
```

Após salvar, a sessão aparece no histórico do paciente.

---

# Requisitos associados

Esta capability é implementada através dos seguintes requisitos:

```text
REQ-01-01-01 Criar sessão
REQ-01-01-02 Editar sessão
REQ-01-01-03 Listar sessões
```

---

# Relação com outras capabilities

Esta capability habilita diretamente:

```text
CAP-01-02 Observações clínicas
CAP-01-03 Intervenções terapêuticas
```

Porque observações e intervenções são registradas **dentro de sessões**.

Também habilita:

```text
CAP-02-01 Histórico do paciente
CAP-02-02 Linha do tempo clínica
```

---

# Fora do escopo desta capability

Esta capability **não inclui**:

```text
análise automática de padrões
assistência por IA
comparação entre casos
agenda de atendimentos
```

Essas funcionalidades pertencem a outras visions e capabilities.

---

# Resultado esperado

Quando esta capability estiver implementada, o Arandu permitirá:

```text
registrar sessões terapêuticas
organizar sessões por paciente
visualizar histórico de sessões
```

Isso estabelece o **primeiro nível de memória clínica estruturada do sistema**.

---

# Observação estratégica

A sessão é o **principal contêiner de informação clínica no Arandu**.

A qualidade do registro de sessões determina diretamente a qualidade futura de:

```text
memória clínica longitudinal
identificação de padrões
assistência reflexiva por IA
```

Por isso, o sistema deve tornar o registro de sessões **rápido, simples e natural para o terapeuta**.

