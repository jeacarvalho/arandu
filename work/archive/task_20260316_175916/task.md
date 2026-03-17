Task: Implementação do REQ-01-01-03 — Listar Sessões de um Paciente

ID da Tarefa: task_20260316_list_sessions

Requirement: REQ-01-01-03

Stack Técnica: Go, templ, HTMX, SQLite, Playwright.

🎯 Objetivo

Implementar a funcionalidade de listagem cronológica das sessões de um paciente específico. O objetivo é fornecer um "índice" de encontros clínico que permita a navegação rápida para os detalhes de cada sessão, seguindo o padrão de "Tecnologia Silenciosa".

🛠️ Escopo Técnico

1. Camada de Domínio (internal/domain/session)

Garantir que a entidade Session possui os campos necessários para a listagem: ID, PatientID, Date, Summary (Resumo).

2. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Implementar o método GetByPatientID(ctx, patientID) no repositório de sessões.

A query SQL deve ordenar os resultados por data de forma decrescente (ORDER BY date DESC).

3. Camada de Aplicação (internal/application/services)

Implementar o método ListSessionsByPatient(ctx, patientID) no serviço de sessão.

4. Camada Web (Componentes templ)

SessionItem: Componente que renderiza uma linha da lista.

Deve exibir a data formatada e o resumo truncado (se existir).

Deve ser um link (<a>) ou usar hx-get para /sessions/{id} com hx-push-url="true".

SessionList: Componente contentor que percorre a lista de sessões.

Deve incluir o botão "Nova Sessão" no topo.

Estilo: Espaçamento generoso (p-6) e bordas inferiores sutis entre itens.

5. Camada Web (Handlers)

GET /patients/{id}/sessions:

Recuperar o ID do paciente da URL.

Chamar o serviço de aplicação.

Renderizar e retornar o fragmento SessionList (fragmento HTMX).

🎨 Design System e UX

Tipografia: Usar a fonte Inter (Sans) para esta lista, garantindo legibilidade de índice.

Espaçamento: Priorizar o "respiro" visual com gap-4 ou gap-6 entre os cards/itens.

Estado Vazio: Se não houver sessões, exibir o componente de "Empty State" com uma mensagem encorajadora e o botão para criar a primeira sessão.

🧪 Protocolo de Testes "Ironclad"

A. Testes de Integração (Repo/Service)

Validar que a consulta ao SQLite retorna as sessões do paciente correto.

Validar que a ordenação cronológica decrescente é respeitada.

B. Teste E2E (Playwright)

O fluxo validado deve ser:

Navegar para a lista de pacientes.

Selecionar um paciente com sessões já registadas.

Verificar: A lista de sessões é carregada via HTMX no perfil do paciente.

Verificar: A primeira sessão da lista é a que possui a data mais recente.

Clicar numa sessão da lista.

Verificar: O sistema navega para os detalhes da sessão (/sessions/{id}) sem recarregar o layout principal.

🛡️ Checklist de Integridade (OBRIGATÓRIO)

[ ] O componente de lista usa .templ e é injetado no Layout via HTMX?

[ ] A ordenação SQL está explicitamente definida como DESC?

[ ] O botão "Nova Sessão" está visível e funcional no topo da lista?

[ ] Executei templ generate e o código Go compilou?

[ ] O script scripts/arandu_guard.sh passou em todas as rotas principais?