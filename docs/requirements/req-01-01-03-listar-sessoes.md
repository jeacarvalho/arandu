REQ-01-01-03 — Listar sessões de um paciente

Identificação

ID: REQ-01-01-03

Capability: CAP-01-01 Registro de sessões

Vision: VISION-01 — Registro da Prática Clínica

Status: draft

História do usuário

Como psicólogo clínico, quero visualizar uma lista cronológica das sessões de um paciente específico, para ter uma visão rápida da frequência dos atendimentos e selecionar uma sessão para consulta ou edição.

Contexto

Após selecionar um paciente na lista global, o terapeuta precisa de um "índice" de encontros. Esta lista não deve ser apenas uma tabela fria, mas uma linha do tempo sutil que prepare o profissional para a imersão na memória clínica longitudinal (VISION-02).

Descrição funcional

O sistema deve listar todas as sessões associadas a um patient_id.

Ordenação: Por padrão, as sessões mais recentes aparecem primeiro (ordem cronológica inversa).

Dados exibidos: Data da sessão, um pequeno resumo (se houver) e status (ex: finalizada, em aberto).

Interação: Cada item da lista deve ser um link para o detalhe da sessão (/sessions/{id}).

HTMX: O carregamento da lista deve ocorrer via HTMX ao abrir o perfil do paciente, sem recarregar a barra lateral.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Tipografia: As datas e resumos devem usar Inter (Sans) para clareza visual de índice.

Design: Lista limpa, com muito espaço em branco (p-4 ou p-6 entre itens).

Visual: Uso de bordas inferiores sutis (border-b border-gray-50) para separar os encontros.

Botão de Ação: Deve haver um link claro de "Nova Sessão" no topo desta lista.

Fluxo

O usuário acessa o perfil do paciente.

O sistema dispara um hx-get para /patients/{id}/sessions.

O backend recupera as sessões do SQLite (via SessionRepository).

O backend renderiza o fragmento SessionList via templ.

O fragmento é injetado na área principal do perfil do paciente.

Rotas Esperadas

GET /patients/{id}/sessions: Retorna o fragmento da lista de sessões.

Critérios de Aceitação

CA-01: A lista deve exibir todas as sessões do paciente em ordem decrescente de data.

CA-02: Se o paciente não tiver sessões, exibir uma mensagem amigável ("Nenhuma sessão registrada") com um convite para criar a primeira.

CA-03: A navegação para o detalhe da sessão deve ser instantânea via HTMX ou hx-boost.

CA-04: A interface deve ser responsiva, adaptando a densidade da lista para dispositivos móveis.

CA-05: O carregamento da lista não deve quebrar o estado da Sidebar.

Persistência (SQL)

-- Busca de sessões por paciente
SELECT id, date, summary, created_at 
FROM sessions 
WHERE patient_id = ? 
ORDER BY date DESC;


Fora do escopo

Filtros por intervalo de datas (será tratado em REQ-02-02-01).

Paginação (inicialmente, carregar todas; paginação será adicionada se houver > 50 sessões).

Exportação da lista (REQ-07-02-01).