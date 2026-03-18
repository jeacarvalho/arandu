Task: Generalização do REQ-07-04-01 — Infinite Scroll no Histórico Clínico

ID da Tarefa: task_20260317_infinite_scroll_clinical_history

Requirement: REQ-07-04-01

Stack Técnica: Go, templ, HTMX, SQLite.

🎯 Objetivo

Expandir o padrão de carregamento progressivo (Infinite Scroll) já existente na lista de pacientes para o Histórico Clínico/Linha do Tempo (REQ-02-01-01). O objetivo é permitir que o terapeuta navegue por anos de sessões sem impactar a performance do navegador ou do banco de dados.

🛠️ Escopo Técnico

1. Camada de Infraestrutura e Aplicação

Repositório: Ajustar a query de agregação da linha do tempo para suportar LIMIT e OFFSET.

Service: O TimelineService deve agora aceitar um parâmetro de page ou offset.

Nota: Como já implementado no dashboard, mantenha a consistência de usar lotes de 20 itens.

2. Camada Web (Componentes templ)

TimelineList:

Identificar o último item do lote atual.

Inserir o gatilho HTMX: hx-get="/patients/{id}/history?offset={next_offset}", hx-trigger="revealed", hx-swap="afterend".

LoadingIndicator: Criar um componente reutilizável para o estado de "A carregar..." que seja removido/substituído após o carregamento do lote.

EndOfHistory: Adicionar um marcador visual em Source Serif 4 para quando a query retornar menos itens que o limit, indicando o fim dos registros.

3. Camada Web (Handlers)

GET /patients/{id}/history:

Detectar se a requisição possui o parâmetro offset.

Se sim, retornar apenas o fragmento de itens da linha do tempo (sem o layout ou cabeçalhos da página).

Se não, retornar a página de prontuário completa.

🎨 Design System (Tecnologia Silenciosa)

Transições: Garantir o uso de classes CSS para um fade-in suave dos novos blocos de história.

Estabilidade: O scroll não deve "pular". A inserção do novo fragmento via afterend deve ser imperceptível para o foco visual do usuário.

Tipografia: O indicador de "Fim dos registos" deve ser discreto, em itálico e fonte Serif.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Stress (Massa de Dados)

Validar o scroll contínuo em um paciente da base de teste que possua mais de 100 sessões.

Verificar: O uso de memória do Chrome/Firefox deve permanecer estável após carregar 5 ou 10 lotes.

B. Teste E2E (Playwright)

Abrir o histórico de um paciente com muitos dados.

Fazer scroll até o final da primeira visualização.

Validar: O indicador de carregamento aparece brevemente.

Validar: Novos itens da linha do tempo aparecem abaixo dos anteriores.

Validar: A data dos itens novos é anterior à data dos itens do primeiro lote (ordem cronológica correta).

Fazer scroll até o fim absoluto.

Validar: A mensagem "Fim dos registos históricos" é exibida.

🛡️ Checklist de Integridade

[ ] O componente reutiliza a lógica de offset já validada no Dashboard?

[ ] O hx-swap="afterend" está sendo aplicado no elemento correto para evitar duplicidade de IDs?

[ ] A query SQL foi otimizada para o uso de OFFSET em grandes volumes?

[ ] Executei templ generate?

[ ] O scripts/arandu_guard.sh passou?