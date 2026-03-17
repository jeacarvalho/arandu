Task: Implementação do REQ-02-01-01 — Visualizar Histórico do Paciente (Prontuário)

ID da Tarefa: task_20260317_view_patient_history

Requirement: REQ-02-01-01

Stack: Go, templ, HTMX, SQLite.

🎯 Objetivo

Implementar a visão consolidada do prontuário (histórico longitudinal). O sistema deve reunir todos os eventos clínicos de um paciente (sessões, observações e intervenções) numa única linha do tempo cronológica, utilizando a estética de "Caderno Clínico" e garantindo a responsividade mobile-first.

🛠️ Escopo Técnico

1. Camada de Domínio (internal/domain/timeline)

Criar a estrutura TimelineEvent:

Type: (Session, Observation, Intervention)

Date: DateTime do evento.

Content: Texto descritivo.

Metadata: Mapa de strings para IDs relacionados (ex: session_id).

Implementar a lógica de ordenação: a lista resultante deve estar sempre em ordem decrescente (mais recente primeiro).

2. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Implementar o método GetTimelineByPatientID(ctx, patientID):

Desafio: Realizar a união eficiente das tabelas sessions, observations e interventions.

Dica SOTA: Pode ser feito via UNION ALL no SQL ou agregando no código Go após múltiplas consultas rápidas.

3. Camada de Aplicação (internal/application/services)

Criar TimelineService para orquestrar a recuperação e formatação dos eventos.

4. Camada Web (Componentes templ)

TimelineContainer: O invólucro principal. Deve incluir a linha vertical sutil e os botões de filtro HTMX (Todos, Notas, Intervenções).

TimelineItem: Renderização individual de cada evento.

Observação: Estilo de nota de margem.

Intervenção: Destaque técnico com a cor primária do Arandu.

Sessão: Marcador de tempo que agrupa os itens.

Tipografia: Uso obrigatório de Source Serif 4 (text-xl) para o conteúdo clínico.

5. Camada Web (Handlers)

GET /patients/{id}/history:

Handler principal que retorna a página completa ou fragmento.

Suportar query params (ex: ?filter=interventions) para filtragem dinâmica via HTMX.

🎨 Padrões de Design (Silent UI)

Fundo: Usar a cor de papel bg-[#F7F8FA].

Cards: Cada grupo de eventos de uma data deve estar num card branco (bg-white) levemente destacado.

Timeline: Uma linha de 2px cinza muito claro (bg-gray-100) que atravessa os eventos verticalmente.

Mobile First: No telemóvel, a linha do tempo deve ser simplificada para uma coluna única, garantindo que o texto clínico ocupe 100% da largura útil.

🧪 Protocolo de Testes "Ironclad"

A. Testes de Integração

Validar que a query de união retorna dados de todas as tabelas.

Validar que a filtragem por tipo (ex: ver apenas intervenções) funciona corretamente no repositório.

B. Teste E2E (Playwright)

Abrir o histórico de um paciente com dados pré-existentes.

Verificar: A linha do tempo exibe sessões e intervenções intercaladas corretamente.

Verificar: Ao clicar no filtro "Intervenções", as observações desaparecem instantaneamente (HTMX).

Verificar: O texto está a usar a fonte Serif.

Verificar: O layout mantém-se íntegro em resolução mobile (375px).

🛡️ Checklist de Integridade

[ ] O componente herda de templates.Layout()?

[ ] A query SQL está otimizada (índices em patient_id e session_id)?

[ ] Os filtros HTMX usam hx-target para atualizar apenas a lista e não a página toda?

[ ] Executei templ generate?

[ ] O scripts/arandu_guard.sh passou sem erros?