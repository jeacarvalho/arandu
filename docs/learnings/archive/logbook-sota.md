Logbook SOTA: Lições Aprendidas e Padrões de Projeto

Este documento é a síntese da inteligência técnica do Arandu. Ele consolida os aprendizados de todas as tarefas, focando em padrões identificados, soluções arquiteturais e prevenção de regressões.

🏗️ 1. Arquitetura Web (Go + templ + HTMX)

🧩 Padrão de Renderização Contextual

Problema: O agente frequentemente esquece de verificar se a requisição é HTMX, quebrando o layout ou retornando a página inteira dentro de um fragmento.

Solução: Sempre verificar o header HX-Request.

true: Retornar apenas o componente fragmentado (component.Render).

false: Envolver no templates.Layout(title, component).Render.

Prevenção: Utilizar a Skill arandu-fullstack-development.

🔠 Dualidade Tipográfica

Regra: UI usa Inter (Sans). Conteúdo clínico usa Source Serif 4 (Serif).

Erro Comum: Usar fontes genéricas em notas de pacientes.

Correção: Aplicar a classe .font-clinical em todos os campos de texto narrativo e snippets de busca.

🛡️ Proteção do Domínio (ViewModels)

Aprendizado: Nunca passar entidades do pacote domain diretamente para o templ.

Padrão: Criar structs de ViewModel no Handler que contenham apenas strings formatadas (ex: datas já em PT-BR) e tipos simples necessários para a View.

💾 2. Persistência e Performance (SQLite SOTA)

🔍 Motor de Busca FTS5

Implementação: Uso de tabelas virtuais com EXTERNAL CONTENT.

Lição: Sincronizar FTS5 via Triggers SQL (INSERT, UPDATE, DELETE) para garantir que a busca nunca fique obsoleta.

Highlighting: O SQLite retorna tags <b> para realce. Para exibir no templ, é necessário o helper RawHTML para evitar o escape automático do HTML.

🚄 Escalabilidade de Dados

Big Data: O sistema foi testado com 63.000 sessões.

Solução: Implementar obrigatoriamente Infinite Scroll (hx-trigger="revealed") em listas longas e Debounce (delay:500ms) em campos de busca para evitar IO excessivo.

Driver: Migrado para modernc.org/sqlite para garantir suporte nativo a FTS5 sem dependência de CGO.

🎨 3. Design System (Tecnologia Silenciosa)

📝 O Conceito de "Silent Input"

Problema: Inputs tradicionais (caixas pretas com bordas) causam fadiga visual e distraem o terapeuta.

Solução: Remover bordas de 4 lados. Usar apenas border-b sutil e fundo bg-[#F7F8FA] (cinza papel). O foco deve ser suave, sem o ring azul padrão.

📱 Mobile-First

Aprendizado: A Sidebar deve ser um drawer (gaveta) em telas pequenas.

Padrão: Usar flex-col no mobile para que o "Painel de Insights" e "Histórico Farmacológico" fiquem abaixo do conteúdo principal, mantendo a leitura fluida.

🛠️ 4. Fluxo de Trabalho do Agente

🔄 Ciclo de Vida da Tarefa

Requisito: Ler o arquivo .md em docs/requirements/.

Schema: Criar migration .up.sql em internal/infrastructure/repository/sqlite/migrations/.

Domínio: Atualizar structs em internal/domain/.

UI: Criar/Editar .templ e obrigatoriamente rodar templ generate.

Handler: Implementar lógica de fragmento.

Guard: Rodar ./scripts/arandu_guard.sh para garantir integridade.

🚨 Anti-Padrões Identificados (NÃO REPETIR)

Anti-Padrão

Consequência

Como evitar

import "github.com/a-h/templ" manual

Erro de compilação (redeclarado)

Deixar o templ generate gerenciar os imports.

SQL Hardcoded no Go

Inconsistência entre Tenants

Usar sempre o sistema de Migrations.

hx-target genérico

Troca de elementos errados

Usar IDs específicos ou closest selectors.

Ignorar updated_at

Perda de rastreabilidade

Garantir que todo UPDATE atualize o timestamp.

Última Atualização: 18 de Março de 2026.