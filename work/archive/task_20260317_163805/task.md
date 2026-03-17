Task: Implementação do Componente UI ThemeCloud (REQ-04-01-01)

ID da Tarefa: task_20260317_ui_theme_cloud

Requirement Relacionado: REQ-04-01-01

Stack Técnica: Go, templ, HTMX, CSS Custom.

🎯 Objetivo

Criar o componente visual ThemeCloud para exibir os temas recorrentes de um paciente. O componente deve ser minimalista, utilizando apenas pesos e tamanhos da fonte Source Serif 4 para indicar importância, e deve ser funcional, permitindo filtrar a linha do tempo clínica ao clicar em um tema.

🛠️ Escopo Técnico

1. Camada de Aplicação / ViewModel

Definir a struct ThemeViewModel:

Name string: O termo clínico.

Count int: Frequência de ocorrência.

WeightClass string: Classe CSS calculada (ex: theme-lv1, theme-lv3, theme-lv5).

Implementar lógica de normalização: Mapear a frequência (Count) para uma escala de 5 níveis de importância visual.

2. Camada Web (Componentes templ)

ThemeCloud:

Contentor com comportamento de "wrap" e espaçamento generoso.

Título sutil em Inter (Sans): "Temas Recorrentes no Prontuário".

ThemeItem:

Renderizar o termo usando obrigatoriamente a classe .font-clinical (Source Serif 4).

Interação HTMX:

hx-get="/patients/{id}/history?q={term}".

hx-target="#timeline-content".

hx-push-url="true".

Estilo: Texto clicável que muda de cor sutilmente no hover (ex: do cinza escuro para o azul Arandu).

3. Camada Web (Handlers)

GET /patients/{id}/analysis/themes:

Chamar o serviço de análise.

Processar os 10-15 principais temas.

Retornar o fragmento ThemeCloud.

🎨 Design System: A "Mancha Temática"

Evite bordas, boxes ou fundos coloridos. Os temas devem parecer "flutuar" no papel. As classes abaixo devem ser definidas no base.css:

Nível 1 (Raro): Classe theme-lv1 (Tamanho base, cor cinza suave).

Nível 3 (Frequente): Classe theme-lv3 (Tamanho grande, peso médio, cor cinza escuro).

Nível 5 (Núcleo): Classe theme-lv5 (Tamanho extra grande, peso negrito, cor grafite).

Background: O componente deve repousar sobre o fundo de "papel" do sistema.

🧪 Protocolo de Testes "Ironclad"

A. Teste Visual

Validar se os termos com maior contagem estão visivelmente maiores e mais "pesados" (bold) que os termos menos frequentes através das classes CSS aplicadas.

Validar se a fonte é efetivamente a Source Serif 4.

B. Teste de Interação (Playwright)

Abrir o perfil de um paciente com a massa de dados de teste.

Localizar o painel de temas.

Clicar no tema "Ansiedade" (ou similar).

Verificar:

A URL mudou para incluir o parâmetro de busca.

A linha do tempo foi atualizada via HTMX para mostrar apenas eventos com o termo.

O layout principal não foi recarregado.

🛡️ Checklist de Integridade

$$$$

 O componente usa .templ?

$$$$

 A filtragem de "Stop Words" (artigos/preposições) foi aplicada no backend?

$$$$

 O componente é responsivo (quebra linhas corretamente no mobile)?

$$$$

 O hx-target aponta para o ID correto da linha do tempo?

$$$$

 Executei templ generate?

