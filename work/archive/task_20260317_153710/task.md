Task: Implementação do REQ-01-00-03 — Busca e Localização de Pacientes

ID da Tarefa: task_20260317_search_patients

Requirement: REQ-01-00-03

Capability: CAP-07-04 — Recuperação de Informação e Performance

Stack: Go, templ, HTMX, SQLite.

🎯 Objetivo

Transformar a navegação de pacientes de uma listagem estática para um modelo de "Search First". Implementar uma barra de pesquisa global (Command Bar) na Top Bar que ofereça resultados instantâneos (autocomplete) com suporte a grandes volumes de dados.

🛠️ Escopo Técnico

1. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Implementar o método SearchByName(ctx, query, limit, offset) no repositório de pacientes.

SQL SOTA: Utilizar a cláusula LIKE com ordenação alfabética e limite de 15 resultados.

Nota: No futuro, este método será migrado para FTS5, mas agora deve garantir a compatibilidade com o schema atual.

2. Camada de Aplicação (internal/application/services)

Adicionar o método SearchPatients(ctx, query) ao PatientService.

Garantir que a busca seja resiliente a strings vazias (retornando os pacientes mais recentes ou uma lista vazia, conforme o contexto).

3. Camada Web (Componentes templ)

TopBarSearch: Implementar o campo de busca no cabeçalho global (layout.templ).

HTMX Trigger: hx-get="/patients/search", hx-trigger="keyup changed delay:500ms, search", hx-target="#search-results".

SearchResults: Componente que renderiza a lista suspensa de resultados.

Cada item deve exibir o nome do paciente e um link para /patients/{id}.

Estilo: bg-white, shadow-lg, rounded-b-xl, fonte Inter.

4. Camada Web (Handlers)

GET /patients/search:

Extrair o parâmetro q da URL.

Chamar o serviço de busca.

Retornar apenas o fragmento SearchResults via templ.

🎨 Design System e Performance

Debounce: O atraso de 500ms é obrigatório para proteger o SQLite de IO excessivo durante a dactilografia rápida.

Visual Silent UI: O campo de busca deve ser discreto. No foco, deve expandir levemente ou ganhar uma sombra interna sutil, sem anéis azuis brilhantes.

Mobile First: Em ecrãs pequenos, a barra de busca deve ocupar a largura total do cabeçalho.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Performance (Backend)

Validar que a query de busca em 500 pacientes demora menos de 50ms.

B. Teste E2E (Playwright)

Abrir qualquer página do sistema (o campo de busca é global).

Digitar "Jose" no campo de busca.

Verificar: O indicador de carregamento (se houver) aparece e desaparece.

Verificar: Uma lista com "Jose Eduardo Alves de Carvalho" (da nossa massa de dados) aparece em menos de 1 segundo.

Clicar no nome do paciente.

Verificar: O sistema navega para o perfil do paciente sem recarregar o layout.

Apagar o texto da busca.

Verificar: A lista de resultados desaparece.

🛡️ Checklist de Integridade

[ ] O campo de busca foi inserido no layout.templ (acessível em todas as telas)?

[ ] O hx-trigger inclui o delay:500ms?

[ ] A busca é case-insensitive no SQLite?

[ ] O scripts/arandu_guard.sh confirma que a nova rota /patients/search está respondendo?

[ ] O componente de resultados desaparece ao perder o foco ou limpar o campo?