REQ-01-02-02 — Editar observação clínica

Identificação

ID: REQ-01-02-02
Capability: CAP-01-02 Registro de observações clínicas
Vision: VISION-01 Registro da prática clínica
Status: draft

História do usuário

Como psicólogo clínico,
quero editar uma observação clínica registrada anteriormente,
para corrigir erros de digitação, ajustar a terminologia técnica ou complementar uma percepção clínica.

Contexto

As observações são unidades atômicas de conhecimento. Embora devam ser registradas rapidamente para não interromper o fluxo clínico, o terapeuta deve ter a liberdade de "limpar" ou expandir essas notas durante o período de reflexão pós-sessão.

A edição deve ser fluida e ocorrer preferencialmente de forma "inline", mantendo o usuário no contexto da sessão.

Descrição funcional

O sistema deve permitir a alteração do conteúdo de uma observação existente.

Interface: Ao clicar em editar, o texto da observação deve ser substituído por um campo de edição (textarea) no mesmo local.

Persistência: O sistema deve atualizar o conteúdo no SQLite e registrar o timestamp de modificação.

Comportamento HTMX: A troca entre o modo "leitura" e o modo "edição" deve ocorrer via troca de fragmentos HTMX, sem recarregamento de página.

Dados da observação (Edição)

Campos editáveis

Content: O texto técnico da observação.

Campos imutáveis

ID: O identificador único.

SessionID: A observação não pode ser movida para outra sessão por este requisito.

CreatedAt: A data original de criação deve ser preservada.

Interface esperada

Seguindo o Design System Arandu:

Tipografia: O campo de edição deve usar a fonte Source Serif 4 para manter a imersão na escrita clínica.

Modo Edição: O componente de edição deve ser discreto, removendo bordas desnecessárias (Silent Input).

Ações: Botões pequenos e claros de "Salvar" e "Cancelar".

Fluxo

O usuário visualiza a lista de observações dentro de uma sessão.

Clica no ícone ou botão de "Editar" em uma observação específica.

HTMX (GET): O sistema substitui o componente ObservationItem pelo fragmento ObservationEditForm.

O usuário altera o texto.

O usuário clica em "Salvar".

HTMX (PUT): O sistema envia os dados, atualiza o banco e retorna o fragmento ObservationItem atualizado.

A interface volta ao modo de leitura com o novo conteúdo.

Rotas esperadas

GET /observations/{id}/edit -> Retorna o fragmento do formulário de edição.

PUT /observations/{id} -> Processa a atualização e retorna o fragmento de visualização.

Critérios de aceitação

CA-01

O sistema deve carregar o conteúdo original da observação corretamente no campo de edição.

CA-02

A edição deve ser realizada via HTMX, atualizando apenas o componente específico na lista de observações.

CA-03

O campo updated_at na tabela observations deve ser atualizado no banco de dados SQLite.

CA-04

O campo de edição deve respeitar a tipografia Source Serif e o limite de caracteres definido no domínio (5000 caracteres).

CA-05

Se o usuário clicar em "Cancelar", o fragmento de edição deve ser substituído de volta pelo fragmento de visualização original sem alterações.

Persistência

Operação: UPDATE
Tabela: observations

UPDATE observations 
SET content = ?, updated_at = ? 
WHERE id = ?;


Integração com outros requisitos

Este requisito complementa:

REQ-01-02-01 Adicionar observação: Garante o ciclo de vida do registro.

VISION-05 Assistência reflexiva: A IA utilizará a versão mais refinada (editada) da observação para gerar insights.

Fora do escopo

Exclusão de observação (REQ-01-02-03).

Histórico de versões (log de alterações) da observação.

Edição de observações em lote.