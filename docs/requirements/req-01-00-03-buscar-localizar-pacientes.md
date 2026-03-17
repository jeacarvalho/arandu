REQ-01-00-03 — Busca e Localização de Pacientes

Identificação

ID: REQ-01-00-03

Capability: CAP-07-04 — Recuperação de Informação e Performance

Vision: VISION-07 — Organização Operacional do Consultório

Status: implementado

História do utilizador

Como psicólogo clínico, quero localizar um paciente rapidamente através de uma barra de busca, para aceder ao seu prontuário sem ter de percorrer uma lista extensa e cansativa de nomes.

Contexto

Com o crescimento da base de dados para centenas de pacientes, a listagem total torna-se impraticável. Este requisito implementa o padrão "Command Bar" ou "Search First", onde a principal forma de navegação entre pacientes é através de uma busca ativa com feedback instantâneo (autocomplete).

Descrição funcional

O sistema deve permitir a filtragem de pacientes em tempo real.

Entrada de Busca: Um campo de texto na Top Bar (Navegação Universal).

Autocomplete: Os resultados devem ser atualizados conforme o utilizador digita.

Performance: Implementar um delay (debounce) de 500ms para evitar chamadas excessivas ao banco de dados.

Paginação de Resultados: Exibir inicialmente os 15 resultados mais relevantes. Se houver mais, permitir o carregamento via "Infinite Scroll" (REQ-07-04-01).

Navegação: Ao clicar num resultado, o sistema deve navegar para o perfil do paciente via HTMX.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Localização: O campo de busca deve estar integrado na Top Bar global, acessível em qualquer ecrã.

Estética: - Ícone de lupa sutil (Lucide search).

Fundo do input levemente acinzentado (bg-gray-100/50) que clareia no foco.

Sem bordas pesadas; foco realçado por uma sombra interna mínima.

Tipografia: Resultados da busca usam Inter (Sans) para máxima clareza na identificação de nomes.

Fluxo HTMX

O utilizador clica no campo de busca.

Digita o nome (ex: "Jose").

Trigger: hx-trigger="keyup changed delay:500ms, search".

Action: hx-get="/patients/search?q=Jose".

Target: Um container de resultados (#search-results) que aparece abaixo do input.

Swap: O fragmento de resultados substitui a lista anterior de forma suave.

Rotas Esperadas

GET /patients/search?q={query}: Retorna o fragmento HTML (templ) com a lista filtrada de pacientes.

Critérios de Aceitação

CA-01: A busca deve retornar resultados parciais (ex: "Mar" deve retornar "Maria" e "Marcos").

CA-02: O sistema deve lidar com grandes volumes (500+ pacientes) mantendo o tempo de resposta do backend inferior a 100ms.

CA-03: Caso não existam resultados, exibir uma mensagem sutil: "Nenhum paciente encontrado com este nome".

CA-04: A interface de busca deve ser totalmente funcional em dispositivos móveis, ocupando a largura total da Top Bar quando ativa.

CA-05: O uso das teclas de seta (cima/baixo) e Enter deve ser suportado para selecionar pacientes na lista de resultados (acessibilidade).

Persistência (SQL)

Inicialmente, utilizaremos LIKE para simplicidade. Futuramente, migraremos para FTS5 (Full Text Search) conforme definido na estratégia de escalabilidade.

-- Busca simples por nome
SELECT id, name FROM patients 
WHERE name LIKE ? 
ORDER BY name ASC 
LIMIT 15;


Fora do escopo

Busca por conteúdo dentro das sessões (será REQ-07-04-02).

Filtros avançados por data de nascimento ou etiquetas (tags).