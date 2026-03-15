# REQ-01-01-02 — Editar sessão

## Identificação

**ID:** REQ-01-01-02 **Capability:** CAP-01-01 Registro de sessões **Vision:** VISION-01 Registro da prática clínica **Status:** draft 

---

# História do usuário

Como **psicólogo clínico**,
quero **editar uma sessão terapêutica já registrada**,
para **corrigir informações, complementar o resumo ou ajustar a data do encontro clínico**.

---

# Contexto

A sessão é um elemento dinâmico da memória clínica. Muitas vezes, o terapeuta inicia um registro durante o atendimento e precisa refiná-lo ou corrigi-lo após um período de reflexão.

Este requisito garante que o registro clínico seja flexível o suficiente para acompanhar o processo de elaboração do profissional, sem perder a integridade da relação com o paciente.

---

# Descrição funcional

O sistema deve permitir a alteração dos dados de uma sessão existente.
As modificações permitidas incluem:

* **Data da sessão**: Ajuste cronológico caso o registro tenha sido feito com data retroativa incorreta.
* **Resumo da sessão**: Atualização do conteúdo narrativo do encontro.

Ao salvar as alterações, o campo `updated_at` deve ser atualizado automaticamente para garantir a rastreabilidade da informação.

---

# Dados da sessão (Edição)

## Campos editáveis

* 
**Date**: Deve ser uma data válida.


* 
**Summary**: Campo de texto livre para o resumo clínico.



## Campos imutáveis

* 
**ID**: O identificador único da sessão não muda.


* 
**PatientID**: Uma sessão não pode ser movida de um paciente para outro para garantir a integridade do histórico.


* 
**CreatedAt**: A data original de criação do registro deve ser preservada.



---

# Interface esperada

A interface de edição deve ser idêntica ou muito similar à de criação, seguindo a **Tecnologia Silenciosa**:

1. Uso da fonte **Source Serif** na área de texto do resumo para facilitar a leitura e escrita reflexiva.


2. Destaque discreto para o nome do paciente no topo.
3. Botão "Salvar Alterações" e opção "Cancelar".

---

# Fluxo

1. O usuário acessa o **Perfil do Paciente** ou a **Visualização da Sessão**.
2. Clica no botão **"Editar"**.
3. O sistema carrega o formulário com os dados atuais da sessão.
4. O usuário realiza as modificações necessárias.
5. Clica em **"Salvar"**.
6. O sistema valida os dados, atualiza o banco e redireciona para a visualização da sessão ou perfil do paciente.

---

# Rotas esperadas

* `GET  /sessions/{id}/edit`
* `POST /sessions/{id}/update` (ou `PUT /sessions/{id}`)

---

# Critérios de aceitação

### CA-01

O sistema deve carregar os dados atuais da sessão corretamente no formulário de edição.

### CA-02

O sistema não deve permitir a alteração do `patient_id` associado à sessão.

### CA-03

Ao salvar, a data `updated_at` no banco SQLite deve ser atualizada para o momento atual.

### CA-04

As alterações devem ser refletidas imediatamente na linha do tempo e histórico do paciente.

### CA-05

Se o usuário cancelar a edição, os dados originais devem ser mantidos sem alterações.

---

# Persistência

**Operação:** `UPDATE`


**Tabela:** `sessions` 

```sql
UPDATE sessions 
SET date = ?, summary = ?, updated_at = ? 
WHERE id = ?;

```

---

# Fora do escopo

* Exclusão da sessão (pertence a um requisito de *Delete*).
* Edição de observações ou intervenções vinculadas (estas possuem requisitos próprios: REQ-01-02-02 e REQ-01-03-02).

---
