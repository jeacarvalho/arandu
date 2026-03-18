Task: Implementação do REQ-07-04-02 — Busca Contextual no Prontuário (FTS5)

ID da Tarefa: task_20260318_contextual_search

Requirement: REQ-07-04-02

Capability: CAP-07-04 — Recuperação de Informação e Performance

Stack Técnica: Go, templ, HTMX, SQLite (FTS5).

🎯 Objetivo

Implementar a funcionalidade de busca textual dentro do histórico clínico de um paciente específico. O sistema deve retornar trechos (snippets) das observações e intervenções onde o termo foi encontrado, com realce visual, permitindo navegação instantânea.

🛠️ Escopo Técnico

1. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Implementar o método SearchInHistory(ctx, patientID, query):

SQL SOTA: Utilizar a função snippet() do FTS5 para gerar o contexto do texto encontrado.

Query: Unir observations_fts e interventions_fts para buscar em ambos os tipos de registro.

Highlight: Configurar o snippet para envolver o termo buscado em tags <b> ou similar.

2. Camada Web (Componentes templ)

SearchHistoryInput: Campo de pesquisa a ser inserido no topo da Linha do Tempo.

Atributos HTMX: hx-get="/patients/{id}/history/search", hx-trigger="keyup changed delay:500ms", hx-target="#timeline-content".

SearchResultItem: Componente para renderizar cada snippet encontrado.

Deve exibir: Data da Sessão, Tipo (Obs/Int) e o trecho de texto com o termo destacado.

Deve usar obrigatoriamente a fonte Source Serif 4 para o conteúdo clínico.

3. Camada Web (Handlers)

GET /patients/{id}/history/search:

Handler que recebe a query q.

Retorna o fragmento de resultados.

Se q estiver vazio, deve retornar o fragmento da Linha do Tempo cronológica padrão.

🎨 Design System (Tecnologia Silenciosa)

Snippets: O texto deve ser truncado de forma inteligente para mostrar o contexto ao redor da palavra-chave.

Highlight: O realce deve ser discreto. Evite cores neon. Use font-bold ou um fundo amarelo pálido (bg-yellow-50).

Vazio: Se nada for encontrado, exibir uma nota sutil: "Nenhum registro clínico contém este termo neste paciente".

🧪 Protocolo de Testes "Ironclad"

A. Teste de Precisão (FTS5)

Pesquisar por um termo técnico presente na massa de dados (ex: "Transferência" ou "Luto").

Verificar: Se o snippet retornado realmente contém a palavra e se o link leva à sessão correta.

B. Teste de Performance (Massa de Dados)

Executar a busca contra um paciente com mais de 200 sessões.

Critério de Aceitação: Resultados renderizados em menos de 300ms.

🛡️ Checklist de Integridade

[ ] A busca está restrita ao patient_id (não vaza dados de outros pacientes)?

[ ] O hx-trigger respeita o delay de 500ms?

[ ] O componente de resultados permite "limpar" a busca facilmente?

[ ] O realce visual está sendo aplicado via classe CSS segura?

[ ] O scripts/arandu_guard.sh passou?