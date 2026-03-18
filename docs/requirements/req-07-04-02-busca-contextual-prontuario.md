REQ-07-04-02 — Busca Contextual no Prontuário (FTS5)

Identificação

ID: REQ-07-04-02

Capability: CAP-07-04 — Recuperação de Informação e Performance

Vision: VISION-02 — Memória Clínica Longitudinal

Status: draft

História do utilizador

Como psicólogo clínico, quero pesquisar termos ou conceitos específicos em todo o histórico de um paciente, para localizar momentos exatos do processo terapêutico e revisar a evolução de temas sensíveis ao longo de anos.

Contexto

Diferente da busca de pacientes (por nome), esta busca vasculha o conteúdo das notas clínicas. Em um prontuário com anos de dados e centenas de observações, encontrar uma fala específica ou a primeira menção a um sintoma é impossível sem ferramentas de indexação. Este requisito utiliza o motor SQLite FTS5 para oferecer uma experiência de busca instantânea com contexto (snippets).

Descrição funcional

O sistema deve permitir a pesquisa textual em todas as observações (observations) e intervenções (interventions) vinculadas a um paciente.

Input de Pesquisa: Um campo dedicado na interface do histórico/prontuário.

Resultados com Contexto (Snippets): O sistema deve exibir o fragmento de texto onde a palavra foi encontrada, não apenas o link para a sessão.

Realce (Highlighting): O termo pesquisado deve aparecer visualmente destacado no fragmento.

Navegação Temporal: Os resultados devem ser apresentados em ordem cronológica inversa (mais recentes primeiro).

Acesso Direto: Ao clicar num resultado, o terapeuta deve ser levado para a visualização completa daquela sessão específica.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Tipografia: Os snippets de texto clínico DEVEM usar a fonte Source Serif 4 (text-lg ou text-xl).

Highlight: O realce do termo deve ser discreto (ex: fundo bg-yellow-50 ou font-bold), evitando cores agressivas que quebrem a imersão.

Metadados: Cada resultado deve exibir a data da sessão e o tipo de registro (Observação ou Intervenção) em fonte Inter (Sans) pequena.

Feedback Visual: Enquanto o utilizador digita, um indicador de "A pesquisar..." sutil deve aparecer.

Lógica Técnica (SQL FTS5)

Para garantir performance sub-segundo na base de 63.000 registos, a consulta deve utilizar as tabelas virtuais FTS5:

-- Exemplo de query para busca de observações com snippet
SELECT 
    s.date,
    s.id as session_id,
    snippet(observations_fts, 0, '[MATCH]', '[/MATCH]', '...', 20) as snippet_text
FROM observations_fts fts
JOIN observations o ON fts.rowid = o.id
JOIN sessions s ON o.session_id = s.id
WHERE fts.content MATCH ? AND s.patient_id = ?
ORDER BY s.date DESC;


Fluxo HTMX

O terapeuta interage com o campo de busca no prontuário.

Trigger: hx-trigger="keyup changed delay:500ms, search".

GET: /patients/{id}/history/search?q={query}.

Target: Substituição dinâmica do container da Linha do Tempo (#timeline-content).

Estado: Se a query for limpa, o HTMX restaura a visualização cronológica completa.

Critérios de Aceitação

CA-01: A busca deve retornar resultados tanto de observações quanto de intervenções.

CA-02: O termo pesquisado deve estar obrigatoriamente destacado no snippet retornado.

CA-03: A performance de resposta para um histórico de 2 anos deve ser inferior a 300ms.

CA-04: O sistema deve garantir que a busca é filtrada rigorosamente pelo patient_id da sessão ativa.

CA-05: A interface deve manter o layout responsivo, garantindo leitura confortável em dispositivos móveis.

Fora do escopo

Busca global entre múltiplos pacientes (limitado ao contexto de um prontuário).

Busca semântica ou por sinonímia (ex: buscar "medo" e encontrar "fobia").

Filtros avançados por sentimentos ou tons de voz.