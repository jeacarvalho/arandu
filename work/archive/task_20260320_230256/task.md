Task: Infraestrutura do Plano Terapêutico e Metas

ID da Tarefa: task_20260321_therapeutic_plan_infra

Requirement: REQ-01-05-01

Stack: Go, SQLite, templ.

🎯 Objetivo

Implementar a base de dados e a camada de domínio para o Planeamento Terapêutico. O sistema deve permitir que o terapeuta comece a estruturar os objetivos clínicos de cada paciente.

🛠️ Escopo Técnico

1. Camada de Persistência (Migration)

Arquivo: internal/infrastructure/repository/sqlite/migrations/0007_add_therapeutic_goals.up.sql.

Ação: Criar a tabela therapeutic_goals conforme definido no requisito.

2. Camada de Domínio (internal/domain/patient)

Criar a entidade TherapeuticGoal e os tipos de estado (GoalStatusInProgress, GoalStatusAchieved, GoalStatusArchived).

Implementar validações: O título da meta não pode ser vazio.

3. Camada Web (Componentes templ)

GoalList: Um container que lista as metas separadas por status.

GoalItem: O componente individual. Metas "Alcançadas" devem ter uma estilização sutil de sucesso (ex: ícone de check verde e texto levemente esmaecido).

AddGoalForm: Formulário "Silent Input" para adicionar rapidamente uma meta ao plano.

🎨 Design System (SOTA)

Tipografia: - Títulos: Inter 14px Semi-bold.

Conteúdo Clínico: Source Serif 4 18px.

Interação: Use hx-swap="afterbegin" para que novas metas apareçam instantaneamente no topo da lista "Em Progresso".

🧪 Protocolo de Testes "Ironclad"

Criação: Adicionar uma meta "Trabalhar autonomia financeira" e verificar se aparece no banco clínico do tenant.

Status: Clicar em "Alcançada" e verificar se o componente se move para a seção de metas concluídas via HTMX.

Isolamento: Garantir que o Plano Terapêutico do Paciente X não é visível quando navegamos para o Paciente Y.