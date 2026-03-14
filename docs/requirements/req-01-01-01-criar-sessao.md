# REQ-01-01-01 — Criar sessão

## Identificação

**ID:** REQ-01-01-01
**Capability:** CAP-01-01 Registro de sessões
**Vision:** VISION-01 Registro da prática clínica
**Status:** draft

---

# História do usuário

Como **psicólogo clínico**,
quero **registrar uma nova sessão terapêutica para um paciente**,
para **documentar o encontro clínico e preservar informações relevantes para acompanhamento terapêutico**.

---

# Contexto

A sessão representa um **encontro terapêutico entre profissional e paciente**.

Ela é o principal elemento temporal do registro clínico no Arandu.

Cada sessão registra o contexto onde ocorreram:

```text
observações clínicas
intervenções terapêuticas
reflexões do terapeuta
```

Estrutura conceitual:

```text
Patient
  └── Session
        ├── Observations
        └── Interventions
```

Cada sessão pertence a **um único paciente**.

Um paciente pode possuir **múltiplas sessões ao longo do tempo**.

---

# Descrição funcional

O sistema deve permitir que o profissional registre uma nova sessão associada a um paciente.

A criação da sessão deve registrar:

```text
paciente
data da sessão
resumo da sessão
```

Após a criação, a sessão deve aparecer no **histórico de sessões do paciente**.

---

# Dados da sessão

## Campos obrigatórios

```text
PatientID
Date
```

## Campos opcionais

```text
Summary
```

## Campos gerados automaticamente

```text
ID
CreatedAt
UpdatedAt
```

---

# Estrutura da entidade

```text
Session

ID
PatientID
Date
Summary
CreatedAt
UpdatedAt
```

---

# Interface esperada

A criação de sessão ocorre a partir da página do paciente.

Exemplo:

```text
Paciente: Maria

[ Nova sessão ]

Data
[ 20/03/2026 ]

Resumo da sessão
[__________________________________]
[__________________________________]

[ Salvar sessão ]
```

---

# Fluxo

```text
Usuário abre página do paciente
↓
Clica "Nova sessão"
↓
Preenche data da sessão
↓
Opcionalmente registra resumo
↓
Clica "Salvar sessão"
↓
Sistema cria sessão
↓
Sessão aparece no histórico do paciente
```

---

# Rotas esperadas

```text
GET  /patients/{id}/sessions/new
POST /sessions
GET  /patients/{id}
```

---

# Critérios de aceitação

### CA-01

O sistema deve permitir criar uma sessão associada a um paciente existente.

---

### CA-02

A sessão deve conter a data em que ocorreu.

---

### CA-03

O sistema deve gerar automaticamente um identificador único para a sessão.

---

### CA-04

Após salvar, a sessão deve ser persistida no banco SQLite.

---

### CA-05

A sessão criada deve aparecer no histórico de sessões do paciente.

---

# Persistência

Tabela:

```sql
sessions
```

Campos:

```text
id TEXT PRIMARY KEY
patient_id TEXT NOT NULL
date DATETIME NOT NULL
summary TEXT
created_at DATETIME NOT NULL
updated_at DATETIME NOT NULL
```

Relacionamento:

```text
sessions.patient_id → patients.id
```

---

# Integração com outros requisitos

Este requisito habilita diretamente:

```text
REQ-01-02-01 adicionar-observacao
REQ-01-03-01 registrar-intervencao
REQ-02-01-01 visualizar-historico
```

Porque observações e intervenções serão registradas **dentro de sessões**.

---

# Fora do escopo

Este requisito **não inclui**:

```text
edição de sessão
exclusão de sessão
agenda de atendimentos
análise automática de sessões
IA
```

Essas funcionalidades serão implementadas em requisitos posteriores.

---

# Resultado esperado

Após a implementação deste requisito, o sistema deve permitir:

```text
registrar sessões terapêuticas associadas a pacientes
```

Isso estabelece a base para registrar **conteúdo clínico estruturado dentro das sessões**.

---
