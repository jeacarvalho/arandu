Task: Narrativa de Desfecho e Relatório de Metas

ID da Tarefa: task_20260321_goal_closure_narrative

Requirement: REQ-01-05-01

Dependência: task_20260321_goals_session_integration (Concluída)

Stack: Go, templ, HTMX.

🎯 Objetivo

Implementar a capacidade de adicionar uma "Nota de Desfecho" quando uma meta é concluída ou arquivada, e criar uma visualização de Relatório de Evolução que consolide o percurso terapêutico do paciente para revisão técnica ou supervisão.

🛠️ Escopo Técnico

1. Camada de Infraestrutura (SQL Migration)

Arquivo: internal/infrastructure/repository/sqlite/migrations/0008_add_goal_closure_notes.up.sql.

Ação: Adicionar a coluna closure_note (TEXT) e closed_at (DATETIME) à tabela therapeutic_goals.

2. Camada Web (Componentes templ)

GoalClosureModal: Um formulário sutil que aparece quando o utilizador clica para concluir uma meta, solicitando uma breve síntese do alcance desse objetivo.

TherapeuticPlanReport: Uma nova página ou vista de "Impressão" que lista:

Racional Clínico do Plano.

Metas Alcançadas (com as suas notas de desfecho).

Metas em Progresso.

ReportTrigger (Ponto de Entrada): Adicionar um botão "Gerar Relatório de Evolução" no topo da seção de Planejamento Terapêutico (dentro do Perfil do Paciente).

O botão deve abrir a rota do relatório em uma nova aba (target="_blank") ou via hx-get se for uma visualização de modal de impressão.

Estilo: Use o padrão "Print-Friendly" (fundo branco puro, tipografia Source Serif 4 dominante, sem sidebars na versão de impressão).

3. Camada Web (Handlers)

POST /goals/{id}/close: Recebe a nota de desfecho, atualiza o status para achieved e grava o timestamp.

GET /patients/{id}/plan/report: Renderiza o relatório consolidado do plano terapêutico em formato otimizado para impressão.

🎨 Design System (SOTA)

Narrativa de Desfecho: O campo de texto deve usar Source Serif 4 para encorajar uma escrita reflexiva e técnica.

Visual do Relatório: O relatório deve parecer um documento oficial e elegante. Use margens de 2cm (simuladas em CSS) e uma marca d'água sutil do Arandu no rodapé.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Conclusão Narrativa

Escolher uma meta "Em Progresso".

Clicar em "Concluir".

Verificar: Se o modal/campo de nota de desfecho aparece.

Escrever a nota e salvar.

Verificar: Se a meta agora exibe a nota de desfecho na lista de "Metas Alcançadas".

B. Teste de Relatório

Localizar o botão "Gerar Relatório" no Perfil do Paciente.

Aceder a /patients/{id}/plan/report.

Verificar: Se todas as metas (ativas e fechadas) aparecem organizadas.

Verificar: Se a tipografia Serif está correta para os textos longos e se os menus (Sidebar/TopBar) estão ocultos na versão de impressão.

🛡️ Checklist de Integridade

[ ] A migration foi aplicada no banco clínico (Tenant)?

[ ] O relatório respeita o isolamento de dados (Multi-tenancy)?

[ ] O scripts/arandu_guard.sh confirma que as novas rotas de relatório estão seguras?

[ ] O componente de fechamento de meta limpa o estado do modal após o sucesso?