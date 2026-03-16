REQ-01-03-01 — Registrar intervenção terapêutica

Identificação

ID: REQ-01-03-01

Capability: CAP-01-03 Registro de intervenções terapêuticas

Vision: VISION-01 Registro da prática clínica

Status: draft

História do usuário

Como psicólogo clínico, quero registrar intervenções específicas realizadas durante o encontro, para documentar a conduta técnica e facilitar a análise posterior da evolução do paciente.

Contexto

Diferente das observações (percepções do terapeuta), as intervenções são as ações ativas realizadas pelo terapeuta (ex: "Realizada técnica de exposição", "Feito questionamento socrático sobre a crença X").

O registro deve ocorrer preferencialmente dentro do contexto da sessão, utilizando a mesma infraestrutura de fragmentos dinâmicos das observações para manter a fluidez do prontuário.

Descrição funcional

O sistema deve permitir o registro de intervenções vinculadas a uma sessão.

Entrada: Texto narrativo curto ou longo descrevendo a ação.

Comportamento HTMX: Envio assíncrono (POST) que adiciona a intervenção ao topo da lista sem recarregar a barra lateral ou o layout principal.

Feedback: Limpeza automática do campo após o sucesso.

Dados da Intervenção

Campos obrigatórios

SessionID: Vínculo com a sessão.

Content: Descrição da intervenção realizada.

Campos automáticos

ID: UUID gerado no backend.

CreatedAt / UpdatedAt: Timestamps para rastreabilidade cronológica.

Interface (Padrão Arandu SOTA)

Seguindo as regras do Canvas:

Tipografia Clínica: O campo de texto e a exibição da intervenção DEVEM usar a fonte Source Serif 4.

Visual: O bloco de intervenção deve ser visualmente distinto das observações (ex: borda lateral cinza escuro ou ícone específico), mas manter a estética minimalista.

Componentização: Implementado via componente .templ dedicado.

Fluxo

O terapeuta acessa uma sessão ativa ou histórico.

Localiza a seção "Intervenções".

Digita a conduta técnica realizada.

Submete o formulário.

HTMX: O backend processa o serviço, salva no SQLite e retorna o fragmento da nova intervenção.

A lista é atualizada via hx-swap="afterbegin".

Rotas Esperadas

POST /sessions/{session_id}/interventions: Adiciona nova intervenção.

Critérios de Aceitação

CA-01: A intervenção deve ser persistida corretamente na tabela interventions.

CA-02: O campo de texto deve ser limpo imediatamente após o sucesso do HTMX.

CA-03: O conteúdo clínico deve obrigatoriamente ser renderizado com a fonte Source Serif 4.

CA-04: O sistema deve validar que o conteúdo não está vazio antes de salvar.

CA-05: A interface deve manter o estado da sidebar (aberta/fechada) durante o processo de registro.

Persistência (SQL Migration)

INSERT INTO interventions (id, session_id, content, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);


Fora do escopo

Categorização automática de intervenções (VISION-03).

Associação de intervenções a metas terapêuticas (VISION-04).