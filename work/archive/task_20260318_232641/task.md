Prompt: Refatoração de Layout e Componentes de Informação (SOTA)

Objetivo: Implementar a estrutura de duas colunas e a organização de dados nos cards conforme a imagem de referência (Dashboard Mariana Alves), utilizando a Paleta Botânica e CSS Puro.

🏗️ 1. Estrutura de Layout (Main Canvas)

O componente dashboard.templ (ou o perfil do paciente) deve ser dividido em duas colunas principais no Desktop (min-width: 1024px):

Coluna Esquerda (65%): Focada no fluxo narrativo (Sessões Recentes).

Coluna Direita (35%): Focada em métricas e cronologia (Evolução e Linha do Tempo).

Mobile: Ambas as colunas empilham verticalmente (100% largura).

📊 2. Bloco de Evolução (Stats Cards)

Implementar o grid de 4 indicadores (Sessões, Meses, Padrões, Hipóteses) seguindo estas regras:

Layout: Grid 2x2.

Estilo do Card: Fundo branco, cantos arredondados (radius-xl), borda sutil --arandu-soft.

Conteúdo: - Valor Numérico: Fonte Sans, text-4xl, cor --arandu-primary.

Label: Fonte Sans, text-[10px], uppercase, cor --arandu-soft, abaixo do valor.

Espaçamento: gap-4 entre os cards.

⏳ 3. Linha do Tempo (Longitudinal)

Criar o componente de cronologia à direita:

Marcador: Uma linha vertical de 2px cor --arandu-soft.

Eventos: - Ponto (Dot): Pequeno círculo cor --arandu-active.

Texto do Evento: Fonte Sans, cor --arandu-dark.

Data: Abaixo do texto, cor --arandu-soft, tamanho pequeno.

Respiro: Margem esquerda generosa para a linha não encostar no texto.

📝 4. Cards de Sessão (Narrativa)

Refatorar a exibição das sessões para o estilo "Caderno":

Badge de Sessão: No topo, "Sessão #XX" em negrito, cor --arandu-active, fundo --arandu-bg.

Relato Clínico: Obrigatoriamente em Source Serif 4, itálico, text-xl, leading-relaxed.

Tags (Chips): Fundo --arandu-bg, texto --arandu-primary, sem bordas, cantos totalmente arredondados (rounded-full).

Indicador Lateral: Uma barra vertical grossa (4px) à esquerda do card com a cor --arandu-active.

💡 5. Painel de Insights (IA)

Posicionamento: No rodapé da área de conteúdo ou fixo na base.

Visual: Fundo --arandu-bg, borda superior --arandu-soft.

Conteúdo: Ícone de lâmpada e texto em Serif itálico com baixa opacidade.

🛠️ Instruções Técnicas para o Agente

CSS Semântico: Crie as classes .evolution-grid, .stat-card-sota, .timeline-marker e .session-notebook-card no style.css.

HTMX Awareness: Garanta que os botões de "+ Registrar Sessão" e "Ver Histórico" no header do paciente utilizem hx-boost="true" ou hx-get para trocas suaves.

Dualidade: Mantenha a regra: Dados Administrativos (Sans) vs Conteúdo Terapêutico (Serif).

Checklist de Validação:

[ ] O grid de estatísticas está em 2x2?

[ ] O relato da sessão está em itálico com a fonte Source Serif 4?

[ ] A linha do tempo vertical está visível e alinhada à direita?

[ ] O fundo de todo o canvas de trabalho é #E1F5EE?