1	# Auditoria UX/Frontend (HTMX + Tailwind) — 2026-03-29
     2	
     3	## 1. 🩺 Diagnóstico Executivo
     4	
     5	A base está funcional e já demonstra boas práticas pontuais (uso de `hx-disabled-elt`, regiões com `aria-live`, tratamento global de erro com toast). Porém, há lacunas críticas em fluxos HTMX que afetam previsibilidade de interação e acessibilidade real em teclado. O design system existe em tokens (`@theme`), mas ainda convive com blocos CSS legados e inconsistências de componente que elevam custo de manutenção.
     6	
     7	- **Nota UX:** 6.8/10  
     8	- **Nota Acessibilidade:** 6.2/10  
     9	- **Nota Código/Manutenibilidade:** 6.5/10
    10	
    11	### Top 3 Showstoppers
    12	1. **IDs inválidos nas opções de busca** quebram `aria-activedescendant` e degradam navegação por teclado do combobox.
    13	2. **Formulário de observação limpa o conteúdo em qualquer resposta** (`hx-on::after-request`), inclusive erro de rede/servidor (risco de perda de texto clínico).
    14	3. **Modal de fechamento de meta fecha mesmo em erro** (`x-on:htmx:after-request`), ocultando falha e causando falsa percepção de sucesso.
    15	
    16	---
    17	
    18	## 2. 📋 Relatório Detalhado de Achados
    19	
    20	### Achado 1
    21	- **Local:** `web/components/patient/search.templ` (linha do `id` das opções).
    22	- **Categoria:** A11y / UX
    23	- **Problema:** `id={ "search-option-" + string(rune(i)) }` produz IDs não semânticos/control chars para índices baixos, quebrando referência em `aria-activedescendant` e navegação assistiva.
    24	- **Solução Proposta:** Gerar IDs estáveis e legíveis com índice numérico (`strconv.Itoa(i)`) ou usar o próprio `r.ID`.
    25	- **Snippet de Correção:**
    26	
    27	```templ
    28	import "strconv"
    29	
    30	<a
    31	  href={ templ.URL("/patients/" + r.ID) }
    32	  id={ "search-option-" + strconv.Itoa(i) }
    33	  role="option"
    34	  class="search-result-link bg-transparent"
    35	  tabindex="-1"
    36	>
    37	```
    38	
    39	### Achado 2
    40	- **Local:** `web/components/session/observation_form.templ`.
    41	- **Categoria:** UX / HTMX / A11y
    42	- **Problema:** O form reseta em `after-request` (inclusive erro), sem `hx-indicator` e sem região local de feedback. O usuário perde texto e não tem confirmação contextual de processamento.
    43	- **Solução Proposta:** Trocar para reset apenas em sucesso (`htmx:afterOnLoad`), adicionar `hx-indicator` e `role=status` com `aria-live`.
    44	- **Snippet de Correção:**
    45	
    46	```templ
    47	<form
    48	  hx-post={ templ.URL("/session/" + sessionID + "/observations") }
    49	  hx-target="#observations-list"
    50	  hx-swap="afterbegin"
    51	  hx-disabled-elt="this button[type='submit']"
    52	  hx-indicator="#obs-loading"
    53	  hx-on::after-on-load="this.reset()"
    54	>
    55	  ...
    56	  <span id="obs-loading" class="htmx-indicator text-xs text-neutral-500" role="status" aria-live="polite">
    57	    Salvando observação...
    58	  </span>
    59	</form>
    60	```
    61	
    62	### Achado 3
    63	- **Local:** `web/components/patient/goal_closure_modal.templ`.
    64	- **Categoria:** UX / HTMX
    65	- **Problema:** `x-on:htmx:after-request="showModal = false"` fecha o modal mesmo com erro HTTP, ocultando falha e removendo contexto de correção.
    66	- **Solução Proposta:** Fechar apenas em sucesso (`afterOnLoad`) e manter modal aberto em erro com mensagem amigável.
    67	- **Snippet de Correção:**
    68	
    69	```templ
    70	<form
    71	  ...
    72	  hx-indicator="#goal-close-loading"
    73	  x-on:htmx:after-on-load="showModal = false"
    74	>
    75	  ...
    76	  <p id="goal-close-feedback" class="text-sm text-red-600" aria-live="polite"></p>
    77	  <span id="goal-close-loading" class="htmx-indicator text-xs text-neutral-500" role="status" aria-live="polite">
    78	    Concluindo meta...
    79	  </span>
    80	</form>
    81	```
    82	
    83	### Achado 4
    84	- **Local:** `web/components/patient/detail.templ` (CTA “Nova Sessão”).
    85	- **Categoria:** UX / HTMX (Histórico/URL)
    86	- **Problema:** Link com `hx-boost` e swap em `main` sem `hx-push-url`, prejudicando histórico do navegador e o botão “Voltar”.
    87	- **Solução Proposta:** Declarar `hx-push-url="true"` nos principais links boosted de navegação.
    88	- **Snippet de Correção:**
    89	
    90	```templ
    91	<a
    92	  href={ templ.URL("/patients/" + patient.ID + "/sessions/new") }
    93	  hx-boost="true"
    94	  hx-target="main"
    95	  hx-swap="innerHTML transition:true"
    96	  hx-push-url="true"
    97	>
    98	```
    99	
   100	### Achado 5
   101	- **Local:** `web/components/layout/layout.templ` (inputs de busca no TopBar/TopBarContent).
   102	- **Categoria:** UX / HTMX / Perf
   103	- **Problema:** Busca HTMX sem `hx-indicator` e sem feedback textual de estado; em latência média o campo parece “travado” e induz repetição de digitação.
   104	- **Solução Proposta:** Adicionar `hx-indicator` dedicado + região de status próxima ao input.
   105	- **Snippet de Correção:**
   106	
   107	```templ
   108	<input
   109	  id="patient-search"
   110	  hx-get="/patients/search"
   111	  hx-trigger="keyup changed delay:350ms"
   112	  hx-target="#search-results"
   113	  hx-indicator="#search-loading"
   114	  ...
   115	/>
   116	<span id="search-loading" class="htmx-indicator text-xs text-neutral-500" role="status" aria-live="polite">
   117	  Buscando pacientes...
   118	</span>
   119	```
   120	
   121	### Achado 6
   122	- **Local:** `web/static/css/input.css`.
   123	- **Categoria:** Tailwind / Manutenibilidade
   124	- **Problema:** Mistura de tokens modernos (`@theme`) com blocos extensos de CSS legado (inclusive `body` duplicado), aumentando divergência visual e custo de manutenção.
   125	- **Solução Proposta:** Consolidar base em utilitários/tokens Tailwind, remover duplicação de `@layer base body`, migrar blocos de auth para componentes utilitários progressivamente.
   126	- **Snippet de Correção:**
   127	
   128	```css
   129	@layer base {
   130	  body {
   131	    @apply bg-arandu-bg text-arandu-dark font-sans;
   132	  }
   133	}
   134	
   135	@utility auth-card {
   136	  @apply w-full max-w-sm rounded-xl bg-white p-8 shadow-sm;
   137	}
   138	```
   139	
   140	---
   141	
   142	## 3. 🚀 Plano de Ação Priorizado
   143	
   144	| Prioridade | Ação | Esforço Estimado |
   145	|---|---|---|
   146	| Alta | Corrigir IDs do combobox (`search-option-*`) para restaurar `aria-activedescendant` e navegação por teclado. | Baixo |
   147	| Alta | Evitar reset de formulário em erro no fluxo de observações; reset somente em sucesso e manter conteúdo em falha. | Baixo |
   148	| Alta | Fechar modal de meta apenas em sucesso e exibir feedback de erro no próprio modal. | Baixo |
   149	| Alta | Garantir `hx-push-url="true"` em navegações principais com `hx-boost` que trocam `<main>`. | Baixo |
   150	| Média | Inserir `hx-indicator` + `role=status` nos fluxos assíncronos críticos (busca, formulários, ações de lista). | Médio |
   151	| Média | Revisar contraste dos textos secundários em dark mode (especialmente `text-neutral-400`/`text-text-tertiary`). | Médio |
   152	| Média | Criar padrão de ordenação de classes Tailwind (layout > spacing > typo > color > state) e aplicar em componentes mais acessados. | Médio |
   153	| Baixa | Consolidar CSS legado de login e outros blocos para utilitários/tokens Tailwind. | Alto |
   154	| Baixa | Reduzir scripts duplicados no layout base para evitar drift e regressões de comportamento. | Médio |