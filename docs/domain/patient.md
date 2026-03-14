# Entidade de domínio — Patient

## Propósito

A entidade **Patient** representa a pessoa atendida pelo profissional de saúde mental.

Ela é a **entidade raiz do domínio clínico do Arandu**.

Todos os registros clínicos são organizados em torno do paciente.

Estrutura conceitual do domínio:

```text
Patient
  └── Sessions
        ├── Observations
        └── Interventions
```

Essa estrutura reflete o fluxo real da prática terapêutica.

---

# Papel no sistema

O paciente é o ponto de partida para:

```text
registro de sessões
registro de observações clínicas
registro de intervenções terapêuticas
análise da evolução terapêutica
identificação de padrões clínicos
```

Sem pacientes registrados não é possível construir histórico clínico.

---

# Escopo da entidade

Neste momento do projeto, a entidade **Patient** deve ser mantida **deliberadamente simples**.

O objetivo é permitir registrar e organizar a prática clínica sem introduzir complexidade desnecessária.

Campos iniciais da entidade:

```text
ID
Name
Notes
CreatedAt
UpdatedAt
```

Descrição dos campos:

**ID**
Identificador único do paciente no sistema.

---

**Name**
Nome pelo qual o paciente será identificado pelo profissional.

Este campo é obrigatório.

---

**Notes**
Campo opcional para registrar observações iniciais sobre o paciente.

Exemplos:

```text
contexto inicial do atendimento
motivo da busca por terapia
informações relevantes da primeira conversa
```

---

**CreatedAt**
Data de criação do registro do paciente.

---

**UpdatedAt**
Data da última atualização do registro.

---

# Decisões de modelagem

O Arandu **não é um sistema administrativo de consultório**, mas um sistema de inteligência clínica.

Por esse motivo, inicialmente **não serão incluídos campos administrativos** como:

```text
email
telefone
CPF
endereço
dados financeiros
```

Essas informações podem ser adicionadas futuramente caso se tornem necessárias.

---

# Relação com outras entidades

O paciente possui relacionamento com:

```text
Session
```

Estrutura:

```text
Patient (1)
   └── Session (N)
```

Cada sessão pertence a um único paciente.

---

# Evolução futura da entidade

Conforme o sistema evoluir, novos atributos poderão ser adicionados.

Possíveis extensões futuras:

```text
status do paciente
tags clínicas
linha terapêutica associada
metadados clínicos
```

Essas extensões devem ser avaliadas com cuidado para não comprometer a simplicidade do modelo.

---

# Princípios de design

A entidade Patient deve seguir os seguintes princípios:

```text
simplicidade
clareza conceitual
independência de infraestrutura
foco no domínio clínico
```

A entidade deve permanecer no pacote:

```text
internal/domain/patient
```

Sem dependência de:

```text
infraestrutura
persistência
web
frameworks
```

---

# Importância para o sistema

A entidade Patient estabelece a base estrutural do Arandu.

A partir dela será possível construir:

```text
histórico clínico longitudinal
organização do conhecimento terapêutico
identificação de padrões entre casos
assistência reflexiva por IA
```

Ela é, portanto, o **primeiro elemento da memória clínica do sistema**.
