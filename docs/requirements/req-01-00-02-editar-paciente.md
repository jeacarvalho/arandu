REQ-01-00-02 — Editar Dados do Paciente (Identidade SOTA)

Identificação

ID: REQ-01-00-02

Capability: CAP-01-00 — Gestão de Pacientes (Padrão SOTA)

Vision: VISION-01 — Registro da Prática Clínica

Status: draft

História do Utilizador

Como psicólogo clínico, quero alterar as informações cadastrais, notas iniciais ou o contexto biopsicossocial de um paciente, para manter os dados atualizados e refletir mudanças ou novas percepções sobre a identidade do sujeito fora do contexto de uma sessão específica.

Contexto

Os dados de um paciente não são estáticos. O prontuário é um documento vivo que deve acompanhar a evolução do sujeito (mudança de ocupação, atualização de identidade de gênero, refinamento da queixa principal).

A edição deve seguir o padrão de Tecnologia Silenciosa, ocorrendo preferencialmente de forma "inline" ou via fragmento, evitando formulários pesados que interrompam a fluidez da navegação clínica.

Descrição funcional

O sistema deve permitir a modificação de todos os campos da entidade Patient e da sua respectiva extensão PatientContext.

Interface de Gatilho: Um botão ou ícone de "Editar" no perfil do paciente (/patients/{id}).

Modo de Edição: Substituição dos dados estáticos por campos de entrada (input, select e textarea) via HTMX.

Persistência: Atualização coordenada entre as tabelas patients e patient_context com registro automático do timestamp updated_at.

Dados para Edição (Identidade SOTA)

Campos Administrativos

Nome (Obrigatório)

Campos de Identidade e Contexto (SOTA/Opcionais)

Etnia/Raça: Atualização baseada em padrões de saúde.

Identidade de Gênero e Orientação Sexual: Fundamental para a clínica afirmativa.

Ocupação e Escolaridade: Atualização de determinantes sociais de saúde.

Notas de Contexto Inicial: Revisão da "Queixa Principal" ou percepções de triagem.

Interface Esperada (Tecnologia Silenciosa)

A edição deve ocorrer de forma integrada ao "papel" digital do prontuário.

Edição Inline: Ao clicar em editar, o texto transforma-se em campo de entrada no mesmo local.

Tipografia:

O campo "Notas" deve usar obrigatoriamente a fonte Source Serif 4 (text-xl) para manter a imersão.

O campo "Nome" e labels administrativos usam Inter (Sans).

Visual (Silent Input):

Inputs sem bordas agressivas (apenas border-b sutil).

Fundo bg-white para os campos ativos, destacando a área de edição sobre o fundo cinza papel do sistema.

Padding generoso para facilitar o toque e a leitura.

Fluxo HTMX

O utilizador visualiza o perfil do paciente e clica em "Editar".

GET: hx-get="/patients/{id}/edit". O sistema retorna o fragmento PatientEditForm.

O utilizador altera os dados biopsicossociais ou narrativos.

Ações:

Salvar: hx-put="/patients/{id}". O sistema processa a transação, guarda no SQLite e retorna o fragmento de visualização atualizado (Read-only).

Cancelar: hx-get="/patients/{id}". O sistema descarta as alterações e restaura a visão de leitura original.

Rotas Esperadas

GET /patients/{id}/edit -> Retorna o formulário de edição (fragmento).

PUT /patients/{id} -> Processa a atualização coordenada das tabelas.

GET /patients/{id} -> Retorna a visualização padrão (usado para cancelar ou após salvar).

Critérios de Aceitação

CA-01

O formulário de edição deve vir pré-preenchido com os dados atuais (incluindo os campos da tabela patient_context).

CA-02

A edição deve ser realizada via HTMX, atualizando apenas o "Main Canvas" do paciente, sem recarregar o Layout principal ou a Sidebar.

CA-03

O sistema deve validar que o campo "Nome" não fique vazio após a edição.

CA-04

O campo updated_at na tabela patients deve ser atualizado para o timestamp do momento do salvamento.

CA-05

A persistência deve ser atômica (transação): se a atualização de patient_context falhar, a alteração em patients não deve ser consolidada.

CA-06

O campo de Notas de Contexto deve respeitar a tipografia Source Serif 4 e o estilo de "escrita fluida".

Persistência (SQL SOTA)

O salvamento envolve a atualização de duas entidades relacionadas:

-- Transação Coordenada
BEGIN TRANSACTION;

UPDATE patients 
SET name = ?, notes = ?, updated_at = ? 
WHERE id = ?;

-- Atualiza ou Insere contexto (UPSERT)
INSERT INTO patient_context (patient_id, ethnicity, gender_identity, sexual_orientation, occupation, education_level)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(patient_id) DO UPDATE SET
    ethnicity = excluded.ethnicity,
    gender_identity = excluded.gender_identity,
    sexual_orientation = excluded.sexual_orientation,
    occupation = excluded.occupation,
    education_level = excluded.education_level;

COMMIT;


Integração com outros requisitos

REQ-01-00-01: Criar Paciente (Origem dos dados).

CAP-01-04: Contexto Biopsicossocial e Farmacológico (Consome estes dados).

Fora do Escopo

Histórico de versões de cada edição (Audit Log/Versioning).

Exclusão definitiva do paciente (REQ-01-00-04).

Alteração manual do id (UUID imutável).