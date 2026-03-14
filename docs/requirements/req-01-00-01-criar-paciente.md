# REQ-01-00-01 — Criar paciente

## Identificação

**ID:** REQ-01-00-01
**Capability:** CAP-01-00 Gestão de pacientes
**Vision:** VISION-01 Registro da prática clínica
**Status:** draft

---

# História do usuário

Como **psicólogo clínico**,
quero **registrar um novo paciente no sistema**,
para **organizar meu acompanhamento terapêutico e registrar sessões associadas a esse paciente**.

---

# Contexto

Paciente é a **entidade raiz do domínio clínico do Arandu**.

Todos os registros clínicos dependem da existência de um paciente.

Estrutura conceitual:

```text
Paciente
  └── Sessões
        ├── Observações
        └── Intervenções
```

Sem paciente não é possível:

```text
registrar sessões
construir histórico clínico
analisar evolução terapêutica
```

Portanto, a criação de paciente é o **primeiro requisito funcional do sistema**.

---

# Descrição funcional

O sistema deve permitir que o profissional registre um novo paciente contendo **as informações mínimas necessárias para iniciar acompanhamento clínico**.

Após a criação:

* o paciente deve ser persistido no banco
* o paciente deve aparecer na lista de pacientes
* o paciente deve permitir criação de sessões futuras

---

# Dados do paciente

Inicialmente o sistema deve solicitar apenas os campos essenciais.

## Campo obrigatório

```text
Nome
```

## Campos opcionais

```text
Observações iniciais
```

## Campos gerados automaticamente

```text
ID
CreatedAt
UpdatedAt
```

---

# Interface esperada

Tela simples de cadastro.

Exemplo:

```text
Novo paciente

Nome
[________________________]

Observações
[________________________]
[________________________]

[ Salvar paciente ]
```

---

# Fluxo

```text
Usuário abre lista de pacientes
↓
Clica "Novo paciente"
↓
Preenche nome
↓
Clica "Salvar"
↓
Sistema cria o paciente
↓
Usuário é redirecionado para página do paciente
```

---

# Rotas esperadas

```text
GET  /patients/new
POST /patients
GET  /patient/{id}
```

---

# Critérios de aceitação

### CA-01

O sistema deve permitir criar um paciente informando apenas o nome.

---

### CA-02

O sistema deve gerar automaticamente um identificador único para o paciente.

---

### CA-03

O paciente deve ser persistido no banco SQLite.

---

### CA-04

Após a criação, o usuário deve ser redirecionado para a página do paciente.

---

### CA-05

O paciente recém-criado deve aparecer na lista de pacientes.

---

# Persistência

Tabela:

```sql
patients
```

Campos:

```text
id TEXT PRIMARY KEY
name TEXT NOT NULL
notes TEXT
created_at DATETIME
updated_at DATETIME
```

---

# Integração com outros requisitos

Este requisito habilita diretamente:

```text
REQ-01-01-01 Criar sessão
REQ-01-01-03 Listar sessões
REQ-02-01-01 Visualizar histórico
```

---

# Fora do escopo

Este requisito **não inclui**:

```text
edição de paciente
exclusão de paciente
agenda
registro de sessões
classificação clínica
IA
```

Essas funcionalidades pertencem a requisitos posteriores.

---

# Resultado esperado

Após a implementação deste requisito, o sistema deve permitir:

```text
registrar pacientes
```

Isso estabelece a **base mínima necessária para registrar sessões clínicas no Arandu**.

---
