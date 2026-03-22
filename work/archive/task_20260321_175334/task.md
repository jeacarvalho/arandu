Task: Engine de Agregação da Linha do Tempo (Timeline)

ID da Tarefa: task_20260321_timeline_engine

Requirement: REQ-02-01-01

Stack: Go, SQL (SQLite), templ.

🎯 Objetivo

ATENÇÃO> PARTE DESSA IMPLEMENTAÇÇÃO DA JÁ FOI REALIZADA. AVALIE O QUE FALTA, DIVIDA EM ETAPAS, APRESENTE O PLANO E AGUARDE O USUÁRIO AUTORIZAR O INÍCIO DAS ALTERAÇÕES

Criar a lógica de backend que consolida eventos de diferentes tabelas (Sessões e Metas) em uma estrutura de dados unificada chamada TimelineEvent, pronta para ser renderizada cronologicamente.

🛠️ Escopo Técnico

1. Camada de Domínio (internal/domain/patient/timeline.go)

Criar a struct TimelineEvent:

ID (string)

Type (SESSION, GOAL_ACHIEVED, MEDICATION_CHANGE)

Date (time.Time)

Title (string)

Summary (string)

Metadata (map[string]string) - para links ou IDs específicos.

2. Camada de Repositório (SQL)

Implementar uma query (ou conjunto de queries) que use UNION ALL para buscar:

Sessões: SELECT id, 'SESSION', date, summary...

Metas Concluídas: SELECT id, 'GOAL_ACHIEVED', closed_at, closure_note...

Ordenação: Garantir que o SQL retorne do mais recente para o mais antigo (DESC).

3. Camada Web (Componente templ)

TimelineView: O container principal da linha do tempo.

TimelineItem: O componente que renderiza cada evento.

Use ícones simples (Lucide: Calendar para sessões, Trophy para metas).

Fonte: Texto narrativo em Source Serif 4.

🎨 Design System (Funcional)

Estrutura: Uma linha vertical sutil à esquerda conectando os eventos.

Cores: Use apenas tons de cinza e o verde Arandu nos ícones para não "sujar" a interface antes da sprint visual.

🧪 Protocolo de Testes "Ironclad"

Paridade Cronológica: Criar uma sessão ontem e concluir uma meta hoje. A meta deve aparecer ACIMA da sessão na linha do tempo.

Isolamento: Garantir que a linha do tempo só mostre eventos do patient_id selecionado e do tenant_id logado.

Performance: Testar com um paciente fictício que tenha 50 eventos para garantir que o scroll é fluido.