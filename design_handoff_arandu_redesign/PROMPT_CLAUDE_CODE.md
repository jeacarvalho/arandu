# PROMPT PARA CLAUDE CODE

Cole este bloco na sua conversa com o Claude Code no VSCode, **depois** de colocar a pasta `design_handoff_arandu_redesign/` dentro do seu projeto (ou como pasta irmã).

---

Sou o desenvolvedor do **Arandu** — sistema clínico para terapeutas em Go + Templ + HTMX + TailwindCSS (PostgreSQL no back, integração com LLM para análises clínicas).

Temos um **redesign completo** que quero portar para a stack real. Os arquivos de referência estão em `./design_handoff_arandu_redesign/` (ou ajuste o caminho conforme a localização). Contém:

- `Arandu Redesign.html` — protótipo React/Babel rodando no navegador. Abra no browser para ver o design em funcionamento (Dashboard, Pacientes, Agenda, Prontuários, Inteligência, Sessão).
- `README.md` — sistema de design completo: paleta Sábio (terrosa, oklch), tipografia (Fraunces serif + Geist sans + Geist Mono), grid, sombras, raios, padrões de cards/pills/botões/formulários.
- `chrome.jsx`, `page_*.jsx`, `data*.jsx`, `icons.jsx`, `llm_panel.jsx`, `app.jsx` — todos os componentes e páginas em React/JSX. São **referência de aparência e comportamento**, não devem ser copiados como código — você vai reescrever tudo em Templ + Tailwind + HTMX.

## Tarefa

Quero que você migre o redesign para a nossa stack real (Templ + Tailwind + HTMX), **começando pela base do design system** antes das telas.

### Fase 1 — Fundação (comece por aqui)

1. Leia `design_handoff_arandu_redesign/README.md` do começo ao fim.
2. Abra `design_handoff_arandu_redesign/Arandu Redesign.html` (via `open` ou Live Server) para ver o design renderizado.
3. Proponha um plano em 3-4 passos para ajustar `tailwind.config.{js,ts}` com os tokens: cores (paper, paper-2, ink, ink-2, ink-3, line, accent, accent-deep, moss, sage, clay, gold, danger), fonts (serif-editorial Fraunces + ui-sans Geist + mono Geist Mono), radii, shadows.
4. Inclua as fontes do Google no layout base.
5. Crie componentes Templ equivalentes a: `Card`, `Pill` (tones: neutral/warn/ok/info/danger/accent), `Avatar` (iniciais em serifa), `BrandMark` (glifo Aperture-Áurea — SVG exato está em `chrome.jsx` no componente `BrandMark`), `Sidebar` com estrutura de nav.
6. Me mostre o diff proposto antes de aplicar.

### Fase 2 — Telas (uma por vez)

Depois de aprovar a fundação, vamos portar página a página, sempre nesta ordem: **Dashboard → Pacientes → Agenda → Prontuários → Inteligência → Sessão**.

Para cada tela: abra o JSX correspondente (`page_dashboard.jsx`, etc.), reescreva em Templ+Tailwind, integre HTMX para as interações (troca de tabs, abrir painéis, navegação), me mostre o resultado, e só então passe para a próxima.

### Regras

- **Fidelidade ao design**: respeite tipografia (serif nos títulos, sans na UI, mono em horas/números), paleta (use tokens do Tailwind, não hex inline), espaçamentos, sombras sutis.
- **Nada de inventar**: se uma cor/ícone/componente não estiver no handoff, pergunte antes.
- **HTMX-first**: troca de páginas, abertura de modais, envio de formulários, inserção de anotações — tudo via `hx-get`/`hx-post`/`hx-swap`. Evite JS client-side exceto o mínimo necessário (ex. toggle do painel LLM, atalhos de teclado).
- **Mantenha handlers Go limpos**: uma função por partial/página, retornando `templ.Component`.
- **Logotipo**: o BrandMark é o SVG Aperture-Áurea (triângulo em proporção áurea + barra acento) — está definido em `chrome.jsx` no componente `BrandMark`. Porte esse SVG exato para um componente Templ `@BrandMark(size int)`.

### Antes de começar

Me faça 3-5 perguntas sobre a estrutura atual do repositório antes de propor o plano: onde fica `tailwind.config`, qual versão do Tailwind/Templ, qual a estrutura de pastas (`cmd/`, `internal/`, `views/`?), se já existe um layout base, e se há convenção para nomes de partials.

---

## Depois que essa tarefa estiver pronta

- Testar responsividade (atualmente o protótipo é desktop-first, confirmar breakpoints).
- Integrar dados reais do Postgres nos lugares onde o protótipo usa `data.jsx`/`data_extra.jsx`.
- Plugar o LLM real nos pontos onde o protótipo tem `llm_panel.jsx` e os cartões de "leitura cruzada" e "temas dominantes".
