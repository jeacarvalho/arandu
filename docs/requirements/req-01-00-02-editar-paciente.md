REQ-01-00-02 — Editar Dados do Paciente

Identificação

ID: REQ-01-00-02

Capability: CAP-01-00 — Gestão de Pacientes

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

História do utilizador

Como psicólogo clínico, quero alterar as informações cadastrais ou notas iniciais de um paciente, para manter os dados atualizados e refletir mudanças ou novas percepções sobre o perfil do paciente fora do contexto de uma sessão específica.

Contexto

Os dados de um paciente não são estáticos. Mudanças de sobrenome, telefone ou a necessidade de ajustar as "Notas de Contexto" (que aparecem no perfil principal) exigem uma interface de edição rápida. Esta edição deve seguir o padrão de "Silent UI", evitando formulários pesados e interrupções bruscas na navegação.

Descrição funcional

O sistema deve permitir a modificação dos campos da entidade Patient.

Interface de Gatilho: Um botão de "Editar" no perfil do paciente (/patients/{id}).

Modo de Edição: Substituição dos dados estáticos por campos de entrada (input e textarea) via HTMX.

Campos Editáveis: Nome, Notas/Contexto Inicial, e outros campos de metadados.

Persistência: Atualização no SQLite com registro automático do timestamp updated_at.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Edição Inline: Preferencialmente, a edição ocorre no próprio local da exibição do dado, transformando o texto em input.

Tipografia: O campo "Notas" deve usar Source Serif 4 para manter a imersão na escrita clínica. O nome usa Inter (Sans).

Visual: Inputs sem bordas agressivas (border-b apenas). Fundo bg-white para destacar a área de edição sobre o fundo cinza papel.

Fluxo HTMX

O utilizador clica em "Editar" no perfil do paciente.

GET: hx-get="/patients/{id}/edit".

O sistema retorna o fragmento PatientEditForm.

O utilizador altera os dados.

PUT: hx-put="/patients/{id}" ao clicar em "Salvar".

O sistema processa a alteração no domínio, guarda no SQLite e retorna o fragmento de visualização atualizado.

Rotas Esperadas

GET /patients/{id}/edit: Retorna o formulário de edição.

PUT /patients/{id}: Processa a atualização dos dados.

Critérios de Aceitação

CA-01: O formulário de edição deve vir pré-preenchido com os dados atuais do paciente.

CA-02: A edição deve ser feita via HTMX, sem recarregamento do Layout principal (mantendo a Sidebar estável).

CA-03: O sistema deve validar que o campo "Nome" não fique vazio após a edição.

CA-04: O campo updated_at deve ser atualizado no banco de dados.

CA-05: Deve existir um botão "Cancelar" que restaura a visão de leitura sem salvar as alterações.

Persistência (SQL)

UPDATE patients 
SET name = ?, notes = ?, updated_at = ? 
WHERE id = ?;


Fora do escopo

Alteração do histórico de sessões (REQ-01-01-02).

Arquivamento/Exclusão do paciente (REQ-01-00-04).

Log de histórico de edições (quem mudou o quê).