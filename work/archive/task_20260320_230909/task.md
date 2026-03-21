Task: Visibilidade e Foco: Metas na Sessão Clínica

ID da Tarefa: task_20260321_goals_session_integration

Requirement: REQ-01-05-01

Dependência: task_20260321_therapeutic_plan_infra (Concluída)

Stack: Go, templ, HTMX.

🎯 Objetivo

Integrar o planeamento terapêutico na interface de registro de sessão. O terapeuta deve visualizar as metas "Em Progresso" de forma sutil enquanto escreve as notas da sessão, garantindo que a prática clínica se mantenha alinhada com os objetivos definidos.

🛠️ Escopo Técnico

1. Camada Web (Componentes templ)

SessionSidebarContext: Criar um novo componente de barra lateral (ou widget) para ser exibido dentro de web/components/session/session_form.templ.

Conteúdo: Exibir apenas as metas com status GoalStatusInProgress.

Estilo:

Fundo levemente diferenciado (ex: bg-arandu-bg/50 ou borda --arandu-soft).

Tipografia: Títulos em Inter (12px, uppercase) e as metas em Source Serif 4 (14px ou 16px).

2. Camada Web (Handlers & Orquestração)

GET /sessions/new / /sessions/{id}:

O handler deve agora buscar as metas ativas do paciente antes de renderizar a página.

Injetar a lista de TherapeuticGoal no componente de sessão.

3. Interação HTMX (Quick Check)

Permitir que o terapeuta marque uma meta como "Alcançada" diretamente da sidebar da sessão, utilizando o handler de atualização de status já criado na Task 1.

Feedback: Ao marcar como alcançada, a meta deve desaparecer da sidebar da sessão com um efeito de fade-out.

🎨 Design System (SOTA)

Não Intrusão: O painel de metas não deve competir com a área de escrita principal. Use margens generosas e cores de baixo contraste para as labels.

Acessibilidade: No mobile, este painel deve ser acessível via um pequeno ícone de "Alvo/Meta" na TopBar da sessão, abrindo um drawer lateral.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Contexto

Criar 3 metas para o "Paciente A" (2 em progresso, 1 alcançada).

Abrir uma nova sessão para o "Paciente A".

Verificar: Apenas as 2 metas "Em Progresso" aparecem no painel lateral.

B. Teste de Sincronização

Abrir a sessão em um separador e o planeamento noutro.

Adicionar uma meta no planeamento.

Atualizar a sessão e verificar se a nova meta aparece.

🛡️ Checklist de Integridade

[ ] O componente de metas na sessão respeita o isolamento por tenant_id?

[ ] A performance de carregamento da sessão foi mantida (query otimizada)?

[ ] A tipografia Source Serif 4 foi usada para o texto das metas?

[ ] O scripts/arandu_guard.sh confirma que as rotas de sessão continuam estáveis?