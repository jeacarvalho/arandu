TASK 20260316_add_intervention

Requirement: REQ-01-03-01
Title: Implementação do Registo de Intervenção Terapêutica

Status

PRONTO_PARA_IMPLEMENTACAO

Objetivo

Implementar a funcionalidade de registo de intervenções técnicas realizadas durante uma sessão clínica. O sistema deve permitir o armazenamento persistente em SQLite e a atualização dinâmica da interface via HTMX utilizando componentes templ.

Referências

docs/requirements/req-01-03-01-registrar-intervencao.md

architecture_sota.md

interface_patterns_sota.md

Escopo Técnico

1. Camada de Domínio (internal/domain/intervention)

Definir a entidade Intervention com campos ID, SessionID, Content, CreatedAt e UpdatedAt.

Implementar validação de domínio: O conteúdo da intervenção não pode ser vazio e deve ter integridade técnica.

2. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Garantir que a tabela interventions existe (verificar migration 0001_initial_schema.up.sql).

Implementar métodos no repositório: Save(ctx, intervention) e GetBySessionID(ctx, sessionID).

3. Camada de Aplicação (internal/application/services)

Criar o serviço InterventionService com o método AddIntervention(ctx, sessionID, content).

Orquestrar a criação do UUID e timestamps.

4. Camada Web (Componentes templ)

InterventionItem: Componente para renderizar uma intervenção individual. Deve usar obrigatoriamente a classe .font-clinical (Source Serif 4).

InterventionForm: Formulário minimalista para inserção.

Usar hx-post para /sessions/{session_id}/interventions.

Usar hx-target apontando para o topo da lista de intervenções.

Usar hx-swap="afterbegin".

Garantir o "Silent Input" (sem bordas pesadas).

5. Camada Web (Handlers)

Implementar o handler para processar o POST de novas intervenções.

Utilizar o padrão de renderização dinâmica: retornar apenas o fragmento InterventionItem para requisições HTMX.

🧪 Protocolo de Testes "Ironclad" (Obrigatório)

A. Testes Unitários

Validar falha ao tentar salvar intervenção sem conteúdo.

Validar sucesso na persistência e recuperação de dados no SQLite (em memória).

B. Teste E2E (Playwright)

O teste deve validar o seguinte fluxo:

Abrir a página de uma sessão clínica existente.

Localizar o campo de texto de "Intervenções".

Introduzir o texto: "Aplicação de técnica de reestruturação cognitiva sobre o pensamento X".

Submeter o formulário.

Verificações:

O campo de texto foi limpo automaticamente.

A intervenção aparece no topo da lista sem recarregar o layout.

A fonte utilizada no texto é a Source Serif 4.

O servidor não caiu e as rotas vizinhas (Dashboard, Lista de Pacientes) continuam acessíveis.

Checklist de Integridade (OBRIGATÓRIO)

[ ] O componente usa .templ e herda de Layout()?

[ ] A tipografia Source Serif 4 foi aplicada ao conteúdo clínico?

[ ] Executei templ generate e o código Go compilou?

[ ] Testei a rota atual e as rotas vizinhas (Regressão)?

[ ] O banco de dados foi atualizado via migration .up.sql?

[ ] O script scripts/arandu_guard.sh passou com sucesso?