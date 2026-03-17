REQ-02-01-01 — Visualizar histórico do paciente (Prontuário)

Identificação

ID: REQ-02-01-01

Capability: CAP-02-01 Histórico do paciente

Vision: VISION-02 — Memória Clínica Longitudinal

Status: draft

História do utilizador

Como psicólogo clínico, quero visualizar todos os eventos clínicos de um paciente (sessões, notas e intervenções) organizados cronologicamente, para identificar a evolução do caso e preparar-me para os próximos encontros com uma visão de longo prazo.

Contexto

Este requisito representa a transição do Arandu de um "diário de sessão" para um "sistema de inteligência clínica". O prontuário não deve ser uma lista de ficheiros, mas uma narrativa contínua onde o terapeuta faz o "scroll" pela história do paciente.

Descrição funcional

O sistema deve consolidar todos os registos atómicos de um paciente numa visão única de linha do tempo.

Agregação: Unir dados das tabelas sessions, observations e interventions.

Cronologia: Exibição do mais recente para o mais antigo.

Filtros Iniciais: Capacidade de alternar a visualização para ver "Apenas Intervenções" ou "Apenas Observações" (via HTMX).

Navegação: Cada evento na linha do tempo deve permitir o salto direto para a edição da sessão correspondente.

Interface (Padrão Arandu SOTA)

Seguindo o design de Caderno Clínico:

Layout de Linha do Tempo: Uma linha vertical subtil que liga os eventos.

Diferenciação Visual: * Observações: Renderizadas como notas de margem ou blocos de texto puro.

Intervenções: Destacadas com um marcador técnico ou cor de acento sutil (ex: o azul primário do Arandu).

Tipografia: Todo o conteúdo clínico DEVE usar Source Serif 4 em tamanho generoso (text-xl).

Performance: Carregamento progressivo (Lazy Loading) à medida que o utilizador faz scroll para baixo.

Fluxo

O utilizador seleciona um paciente.

Clica na aba ou secção "Prontuário/Histórico".

O sistema executa uma query complexa (ou múltiplas queries) para reunir os dados.

O componente .templ processa a intercalação cronológica dos eventos.

A página é renderizada com foco na legibilidade e no "respiro" visual.

Rotas Esperadas

GET /patients/{id}/history: Retorna a visão consolidada do prontuário.

Critérios de Aceitação

CA-01: Todos os eventos clínicos (sessões, observações e intervenções) devem aparecer na ordem correta.

CA-02: O sistema deve manter a imersão visual, escondendo menus desnecessários durante a leitura profunda.

CA-03: A troca de filtros (ex: ver apenas intervenções) deve ser instantânea via HTMX.

CA-04: O layout deve ser responsivo, permitindo a leitura confortável no telemóvel (Mobile-First).

CA-05: Registos sem conteúdo não devem ocupar espaço visual.

Persistência (Lógica de Query)

Embora os dados estejam em tabelas separadas, a query deve ser orquestrada pelo serviço de aplicação para criar uma lista de TimelineEvents.

-- Exemplo de lógica de união (pseudo-SQL)
SELECT 'session' as type, date as event_date, summary as content FROM sessions WHERE patient_id = ?
UNION ALL
SELECT 'observation' as type, created_at, content FROM observations WHERE session_id IN (SELECT id FROM sessions WHERE patient_id = ?)
-- ... repetir para intervenções
ORDER BY event_date DESC;


Fora do escopo

Pesquisa por palavras-chave dentro do histórico (REQ-09-01-01).

Impressão em PDF (será tratado na Vision-07).

Gráficos de evolução quantitativa.