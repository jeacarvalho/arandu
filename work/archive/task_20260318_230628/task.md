Prompt de Refatoração: Identidade Botânica e Tecnologia Silenciosa (Versão style.css)

Objetivo: Converter o Arandu de um software administrativo para um ambiente de trabalho cognitivo (SOTA), aplicando a nova Identidade Visual Botânica, removendo a dependência de Tailwind CSS e garantindo responsividade Mobile-First.

🎨 1. A Nova Paleta (Identidade Botânica)

Deve atualizar o ficheiro web/static/css/style.css. Substitua todas as definições de cores anteriores no bloco :root pelas variáveis CSS abaixo:

--arandu-primary: #0F6E56; (Verde Base)

--arandu-active: #1D9E75; (Destaque/Interação)

--arandu-soft: #9FE1CB; (Acentos Suaves)

--arandu-bg: #E1F5EE; (Fundo Papel de Seda - Anti-fadiga)

--arandu-dark: #085041; (Texto Principal - Verde Floresta Escuro)

--arandu-paper: #FFFFFF; (Superfícies de Trabalho/Cards)

🏗️ 2. Diretrizes Técnicas (Pure CSS SOTA)

Remoção do Tailwind: Elimine as classes utilitárias do Tailwind (ex: bg-white, p-8, flex) dos ficheiros .templ. Utilize as classes semânticas que irá criar no style.css.

Dualidade Tipográfica:

Interface (UI): Inter (Sans).

Clínica (Notas, Sessões, Nomes): Source Serif 4 (Serif). Use obrigatoriamente a classe .font-clinical.

Padrão "Silent UI":

Inputs: Apenas border-bottom sutil. Sem bordas laterais ou superiores.

Cards: Use a classe .clinical-card com cantos arredondados (24px) e sombras quase imperceptíveis.

📱 3. Plano de Ação por Arquivo

Passo 1: Atualizar web/static/css/style.css

Limpe o ficheiro de estilos antigos que conflitem com a nova estética.

Implemente as classes base: .clinical-card, .stat-card, .silent-input, .font-clinical.

Adicione Media Queries nativas para garantir que o layout empilhe corretamente em ecrãs de 375px.

Passo 2: Refatorar web/components/layout/layout.templ

Adapte a TopBar para ser fixa (64px) com o logo Arandu em itálico serifado.

Transforme a Sidebar num Drawer lateral (oculto por padrão no mobile, acionado por Alpine.js).

O fundo do body deve ser obrigatoriamente --arandu-bg.

Passo 3: Refatorar web/components/dashboard/dashboard.templ

Substitua as classes de grid do Tailwind por um sistema de grid nativo no CSS.

Aplique o estilo de .patient-row na lista de pacientes recentes, garantindo que o nome do paciente use Source Serif 4.

🛡️ Checklist de Integridade

[ ] O código compila sem erros após templ generate?

[ ] O fundo da aplicação é a cor Papel de Seda (#E1F5EE)?

[ ] A Sidebar funciona como um Drawer (gaveta) em visualização mobile?

[ ] As notas clínicas e relatos usam Source Serif 4 em Verde Escuro (#085041)?

[ ] O script ./scripts/arandu_guard.sh valida todas as rotas com sucesso?

Instrução de Persona: Você é um designer/desenvolvedor de elite focado em Saúde Mental. O sistema deve parecer um objeto de luxo feito de papel, transmitindo calma e autoridade técnica.