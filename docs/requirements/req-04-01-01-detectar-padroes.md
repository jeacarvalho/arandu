REQ-04-01-01 — Detectar Padrões de Comportamento e Temas Recorrentes

Identificação

ID: REQ-04-01-01

Capability: CAP-04-01 Identificação de Padrões Clínicos

Vision: VISION-04 — Descoberta de Padrões Clínicos

Status: draft

História do utilizador

Como psicólogo clínico, quero identificar os temas e termos mais frequentes nos registros de um paciente, para compreender os núcleos de sofrimento recorrentes e validar minha percepção intuitiva sobre a evolução do caso.

Contexto

Com o acúmulo de anos de terapia (como na nossa base de 63.000 sessões), temas importantes podem se perder na cronologia. Este requisito implementa um motor de análise de frequência que varre Observações e Intervenções para destacar "nuvens de conceitos" que definem a jornada daquele paciente específico.

Descrição funcional

O sistema deve realizar uma análise estatística de termos sobre o prontuário do paciente.

Filtro Temporal: O terapeuta pode escolher o período (ex: "últimos 6 meses", "todo o histórico").

Processamento de Texto: O motor deve filtrar stop words (artigos, preposições, conjunções) para focar em substantivos e adjetivos clínicos (ex: "ansiedade", "mãe", "conflito", "trabalho").

Ranking de Relevância: Exibir os 10 termos mais citados.

Interatividade: Ao clicar num tema, o sistema deve filtrar a linha do tempo (REQ-02-01-01) para mostrar apenas os eventos onde aquele termo aparece.

Lógica Técnica (SOTA)

1. Motor de Busca (FTS5)

Utilizar a tabela virtual observations_fts e interventions_fts para uma contagem ultra-rápida.
O uso de fts5vocabulary permite extrair a lista de termos únicos e suas frequências sem percorrer as tabelas pesadas linha por linha.

2. Algoritmo de Extração

Buscar todos os tokens associados ao patient_id.

Cruzar com uma lista de "Stop Words Clínicas" em Português.

Agrupar por radical (stemming) ou palavra exata.

Ordenar por contagem descrescente.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Visualização: Evitar "gráficos de pizza" barulhentos. Usar uma lista elegante de "Temas Recorrentes" com pesos tipográficos sutis (tamanho da fonte ligeiramente maior para termos mais frequentes).

Tipografia: Os temas devem ser exibidos em Source Serif 4 para manter a conexão com o conteúdo clínico.

Localização: Este componente deve aparecer como um "Painel de Insights" lateral ou no topo do histórico do paciente.

Fluxo

O terapeuta acessa o prontuário do paciente.

O sistema executa em background (ou via trigger HTMX) a query de frequência.

O componente ThemeCloud renderiza os termos mais relevantes.

O terapeuta clica em "Ansiedade".

HTMX: A linha do tempo é filtrada instantaneamente para exibir o contexto de uso do termo.

Rotas Esperadas

GET /patients/{id}/analysis/themes: Retorna o fragmento com o ranking de temas.

Critérios de Aceitação

CA-01: O processamento deve ser sub-segundo mesmo em pacientes com mais de 100 sessões.

CA-02: Termos comuns (eu, ele, para, com) não devem aparecer no ranking.

CA-03: A interface deve ser minimalista, sem cores berrantes, integrando-se organicamente ao fundo "cinza papel".

CA-04: O sistema deve permitir a navegação da análise para o texto bruto de forma fluida.

CA-05: A análise deve respeitar o isolamento do banco de dados (multi-tenancy).

Persistência (SQL FTS5)

-- Exemplo de contagem rápida via FTS5
SELECT term, count(*) as freq 
FROM (SELECT * FROM observations_fts WHERE observations_fts MATCH ?)
GROUP BY term 
ORDER BY freq DESC 
LIMIT 10;


Fora do escopo

Análise de sentimento (IA complexa).

Identificação de padrões entre pacientes diferentes (CAP-10-01).

Exportação de relatórios estatísticos.