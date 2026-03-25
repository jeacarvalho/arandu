De Formulário Administrativo para "Primeira Página de Caderno"

Contexto: A tela de Novo Paciente está funcional, mas visualmente ruidosa. Sua tarefa é purificar o design, aplicar a Identidade Botânica e incluir os campos de Identidade SOTA (REQ-01-00-01) usando apenas CSS Puro e templ.

🧹 1. Purificação e Limpeza (Prioritário)

Remoção de Inline Styles: Remova TODOS os atributos style="..." das tags HTML. Mova qualquer lógica de espaçamento para classes semânticas no base.css.

Eliminação de Gradientes: Remova os backgrounds coloridos e gradientes dos ícones. No SOTA, os ícones devem ser apenas contornos sutis em var(--arandu-primary) ou desaparecer se não agregarem valor à escrita.

🏗️ 2. Novo Layout: O Fluxo da História

Em vez de cards isolados, organize a página como uma sequência fluida de escrita:

A. Cabeçalho de Identidade

Campo Nome: Deve ser o destaque. Remova o box. Use a classe .silent-input com fonte maior, transmitindo que este é o título do prontuário que está a ser aberto.

Campos Biopsicossociais (Novos): Adicione um grid sutil (2 ou 3 colunas) abaixo do nome para:

Identidade de Gênero

Etnia/Raça

Ocupação

Escolaridade

Estilo: Estes campos devem usar a fonte Inter (Sans) por serem dados técnicos/administrativos.

B. O Corpo da Anamnese (Notas Iniciais)

Textarea: Deve ocupar a maior parte da largura (max-w-3xl).

Soberania Serif: Use obrigatoriamente font-family: var(--font-clinical) (Source Serif 4) com text-xl e leading-relaxed.

Silent UI: Remova bordas e fundos cinzas. O texto deve parecer escrito sobre o "papel de seda" (#E1F5EE).

🎨 3. Consolidação Visual (CSS)

Background: Garanta que o fundo da página é --arandu-bg (#E1F5EE).

Botões: O botão "Cadastrar" deve ser discreto. Use a classe .btn-primary (fundo verde botânico, sem sombras agressivas).

Labels: Use a classe .silent-label (Sans, small-caps, cor esmaecida).

🧪 Protocolo de Validação para o Agente

[ ] O formulário contém os campos de Gênero, Etnia e Ocupação?

[ ] As "Observações Iniciais" estão usando a fonte Source Serif 4?

[ ] Todos os atributos style inline foram removidos?

[ ] O layout permanece centralizado e legível em dispositivos móveis (375px)?

[ ] O comando templ generate foi executado?

Instrução de Persona: Você é um arquiteto focado em UX Clínico. Trate esta tela não como uma captura de dados, mas como o início de um vínculo terapêutico. O design deve ser calmo e encorajar a descrição detalhada do sujeito.