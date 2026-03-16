# REQ-01-02-01 — Adicionar observação clínica

## Identificação

* 
**ID:** REQ-01-02-01 


* 
**Capability:** CAP-01-02 Registro de observações clínicas 


* 
**Vision:** VISION-01 Registro da prática clínica 


* 
**Status:** draft 



---

## História do usuário

Como **psicólogo clínico**, quero **registrar uma percepção clínica específica durante ou após uma sessão**, para **documentar padrões e insights que serão usados na análise da evolução do paciente**.

---

## Contexto

Diferente do resumo da sessão (que é narrativo), a observação clínica é uma unidade atômica de percepção técnica. Ela deve ser registrada de forma rápida para não interromper o fluxo de pensamento do terapeuta. No banco de dados, ela pertence a uma `Session`.

---

## Descrição funcional

O sistema deve permitir a adição de múltiplas observações em uma única sessão.

* **Entrada:** Texto livre.
* 
**Comportamento HTMX:** A adição de uma observação deve ser feita via `POST` assíncrono, atualizando apenas a lista de observações na tela, sem recarregar a página.


* **Feedback:** O campo de texto deve ser limpo após o envio bem-sucedido.

---

## Dados da Observação

### Campos obrigatórios

* 
**SessionID:** Vínculo com a sessão atual.


* 
**Content:** O texto da observação clínica.



### Campos gerados automaticamente

* 
**ID:** UUID único.


* 
**CreatedAt:** Data e hora do registro.



---

## Interface (Padrão Arandu)

Seguindo o **Design System**, a interface de adição de observação deve ser minimalista:

* 
**Tipografia:** O campo de digitação (textarea) deve usar a fonte **Source Serif** para promover a imersão clínica.


* 
**Localização:** Bloco lateral ou inferior dentro da visualização da sessão.


* 
**Estilo:** "Input silent" — sem bordas pesadas, assemelhando-se a uma folha de papel.



---

## Fluxo

1. O terapeuta está na tela de uma **Sessão** ativa ou editando uma sessão passada.


2. Digita a percepção no campo "Nova Observação".
3. Pressiona "Adicionar" (ou atalho de teclado).
4. O sistema valida o conteúdo e persiste no SQLite.


5. A nova observação aparece no topo da lista de observações daquela sessão via HTMX.



---

## Rotas e Componentes (templ)

* **Rota:** `POST /sessions/{session_id}/observations`
* **Componente de Retorno:** `ObservationItem(obs domain.Observation)` (renderiza apenas a nova linha na lista).

---

## Critérios de Aceitação

* 
**CA-01:** A observação deve ser salva com sucesso vinculada ao `SessionID` correto.


* **CA-02:** Não deve ser possível salvar uma observação vazia.
* 
**CA-03:** A lista de observações deve ser atualizada instantaneamente via HTMX.


* 
**CA-04:** O campo de texto deve usar obrigatoriamente a fonte **Source Serif**.


* 
**CA-05:** O registro deve ser persistido na tabela `observations` do SQLite.



---

## Persistência (SQL)

```sql
INSERT INTO observations (id, session_id, content, created_at)
VALUES (?, ?, ?, ?);
[cite_start]
http://googleusercontent.com/immersive_entry_chip/0
