 Missão: Refatoração Estética do Dashboard (Identidade Botânica SOTA)

Contexto: O Dashboard atual está funcional, mas visualmente "administrativo". Sua tarefa é aplicar a Identidade Botânica e os padrões de "Tecnologia Silenciosa" utilizando apenas CSS Puro e componentes templ.

🎨 1. Configuração de Variáveis (Base CSS)

Certifique-se de que o base.css utiliza rigorosamente estas cores:

--arandu-bg: #E1F5EE; (Fundo Papel de Seda)

--arandu-primary: #0F6E56; (Verde Botânico para Sidebar/Botões)

--arandu-active: #1D9E75; (Acentos de interação)

--arandu-paper: #FFFFFF; (Cards brancos flutuantes)

--arandu-dark: #085041; (Texto principal)

🏗️ 2. Refatoração do Layout e Grid

Background Global: Defina body { background-color: var(--arandu-bg); }.

Evolution Grid: O grid de 4 indicadores deve ser 2x2 no mobile e 4 colunas no desktop. Os números (stat-value-sota) devem ser grandes, em var(--arandu-primary), e as labels em uppercase pequeno e cinza-suave.

Main Canvas: A divisão deve ser 65% (Esquerda: Narrativa) e 35% (Direita: Métricas). No mobile, empilhe em 100%.

📝 3. O Padrão "Silent UI"

Search Bar: Remova os estilos inline. Use um fundo rgba(255,255,255,0.5) que se torna branco puro no foco. Arredondamento total (border-radius: 24px).

Cards Clínicos: Use .clinical-card com border-radius: 24px, sem sombras pesadas (apenas uma borda de 1px sutil rgba(159, 225, 203, 0.3)).

Soberania Serif: O título "Olá, Dra. Gabriela", os nomes de pacientes na lista e o texto do Painel de Insights DEVEM usar font-family: var(--font-clinical) (Source Serif 4).

⏳ 4. Componente: Linha do Tempo Longitudinal

A linha vertical deve ser --arandu-soft.

O "Dot" do evento deve ser um círculo pequeno --arandu-active.

Datas em Sans pequeno, conteúdo em Serif.

💡 5. Painel de Insights IA

Posicione-o como um bloco de destaque (Card escuro --arandu-dark ou com fundo dourado pálido).

O texto deve ser em itálico Serif: "O sistema detectou padrões..." - deve parecer um conselho de um supervisor experiente.

🛡️ Checklist de Saída (Obrigatório)

[ ] O fundo do sistema é o verde "papel de seda" (#E1F5EE)?

[ ] Os nomes e relatos clínicos estão em Source Serif 4?

[ ] A SidebarDrawer (Mobile) fecha ao clicar no overlay?

[ ] Não restam estilos inline no HTML gerado?

[ ] O comando templ generate foi executado?