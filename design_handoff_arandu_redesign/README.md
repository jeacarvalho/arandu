# Handoff: Arandu — Redesign do sistema clínico

## Visão geral
Proposta de redesign visual para o Arandu — sistema de gestão clínica para terapeutas e psicólogos (pacientes, sessões, agenda, observações clínicas integradas a LLM para análises cruzadas). O redesign mantém a estrutura funcional existente (sidebar + topbar + canvas) mas substitui a identidade visual atual (verde saturado, gradientes claros) por um sistema editorial acolhedor chamado **"Sábio"** — paleta terrosa, tipografia com serifa para conteúdo clínico, e painel de IA integrado como drawer lateral.

## Sobre os arquivos
Os arquivos neste pacote são **referências de design criadas em HTML + React (via Babel standalone)**. Eles são protótipos que demonstram aparência, hierarquia e comportamento pretendidos — **não são código de produção para copiar diretamente**.

O projeto original é **HTMX + Tempo + Tailwind** (server-rendered). A tarefa é **recriar estes designs no ambiente existente do Arandu**, usando HTMX para interatividade e Tailwind para estilização, aproveitando os padrões e convenções já estabelecidos no codebase. Os componentes React aqui são apenas uma forma de organizar o protótipo — traduza-os para parciais HTML + classes Tailwind (ou componentes Tempo).

## Fidelidade
**Alta fidelidade (hifi).** Cores, tipografia, espaçamentos e interações estão finalizados. Recrie pixel-perfect, adaptando para HTMX/Tailwind. Todos os valores (hex, tamanhos em px, pesos de fonte) estão no CSS do `Arandu Redesign.html` e nos componentes JSX — consulte-os como fonte de verdade.

---

## Paleta a implementar: **Sábio (day)**

Esta é a paleta escolhida. Ignore as variantes `moss`, `clay` e `night` — elas existem no protótipo como exploração via Tweaks, mas não são parte do entregável.

### Design tokens (CSS custom properties → pode mapear para `tailwind.config.js` em `theme.extend.colors`)

| Token            | Valor      | Uso                                         |
|------------------|------------|---------------------------------------------|
| `--paper`        | `#F5EFE6`  | background principal (papel cru)            |
| `--paper-2`      | `#EDE5D8`  | cards e superfícies elevadas                |
| `--paper-3`      | `#E3D9C6`  | superfícies mais elevadas / hover           |
| `--ink`          | `#1F1A15`  | texto principal, botões primários           |
| `--ink-2`        | `#3A332B`  | texto secundário                            |
| `--ink-3`        | `#6B5F52`  | texto terciário, labels                     |
| `--ink-4`        | `#9B8E7E`  | texto placeholder, ícones inativos          |
| `--line`         | `#D9CDB8`  | bordas padrão                               |
| `--line-2`       | `#C6B89F`  | bordas de ênfase / tracejadas               |
| `--accent`       | `#6B4E3D`  | marrom terroso — acento principal           |
| `--accent-deep`  | `#4A3527`  | acento escuro (hover, strong)               |
| `--accent-soft`  | `#A67C52`  | caramelo — acento suave                     |
| `--moss`         | `#2E3A2C`  | verde profundo (status ok/concluído)        |
| `--moss-2`       | `#4A5D4F`  | verde médio                                 |
| `--sage`         | `#8FA68E`  | verde sálvia (decoração)                    |
| `--clay`         | `#C67B5C`  | terracota (warning suave)                   |
| `--gold`         | `#B8925A`  | dourado (highlights)                        |
| `--ok`           | `#4A5D4F`  | sucesso                                     |
| `--warn`         | `#C67B5C`  | atenção                                     |
| `--danger`       | `#A0463A`  | erro / risco                                |

### Raios
- `--radius-sm`: **8px** (botões pequenos, chips)
- `--radius`: **14px** (cards, inputs)
- `--radius-lg`: **22px** (modais)
- Pills / badges: **999px** (full rounded)

### Sombras
- `--shadow-sm`: `0 1px 0 rgba(31,26,21,.04), 0 1px 2px rgba(31,26,21,.04)` (cards)
- `--shadow-md`: `0 1px 0 rgba(31,26,21,.04), 0 8px 24px -8px rgba(31,26,21,.08)` (FAB, popovers)
- `--shadow-lg`: `0 1px 0 rgba(31,26,21,.05), 0 24px 48px -16px rgba(31,26,21,.14)` (drawer LLM, modais)

---

## Tipografia

| Família        | Google Fonts                                                         | Pesos usados         | Uso                                                    |
|----------------|----------------------------------------------------------------------|----------------------|--------------------------------------------------------|
| **Fraunces**   | `Fraunces:ital,opsz,wght@0,9..144,400;0,9..144,500;0,9..144,600;1,9..144,400` | 400, 500, 600, 400i  | Títulos, KPIs numéricos, citações, conteúdo clínico    |
| **Geist**      | `Geist:wght@300;400;500;600;700`                                     | 300–700              | UI: labels, botões, corpo, navegação                   |
| **Geist Mono** | `Geist+Mono:wght@400;500`                                            | 400, 500             | Horários, IDs, números tabulares, kbd                  |

Classes utilitárias no protótipo: `.serif` (Fraunces) e `.mono` (Geist Mono). Tudo que não tem classe usa Geist (`--font-ui`).

### Escala tipográfica (usada nos mocks)
- Display serifa (H1 dashboard/paciente): **40–44px**, weight 400, letter-spacing −0.8px, line-height 1.05
- H1 sessão: **32px**, weight 400, letter-spacing −0.6px
- H3 card title (serifa): **19px**, weight 500, letter-spacing −0.2px
- Corpo serifa (blockquote/observações): **14.5–17px**, line-height 1.55–1.65
- Corpo UI: **13–14px**, line-height 1.5
- Eyebrow / label: **10.5–11px**, uppercase, letter-spacing 1.4px, weight 500, cor `--ink-3`
- Mono (timestamps): **11–13px**, `tabular-nums`
- KPI grande: **42px** serifa, weight 400, letter-spacing −1px

---

## Layout geral

```
┌─────────┬──────────────────────────────────────────────┐
│         │  Topbar (64px, sticky, blur)                 │
│ Sidebar │├─────────────────────────────────────────────┤
│ 232px   │                                              │
│ (68px   │  Main (max-width 1320px, padding 32/40/60)   │
│ quando  │                                              │
│ colap.) │                                              │
└─────────┴──────────────────────────────────────────────┘

+ LLMPanel: drawer à direita, 460px, trigger ⌘J / Ctrl+J
+ Tweaks FAB: canto inferior direito (só no protótipo; remover do produto)
```

### Sidebar (`chrome.jsx` → `Sidebar`)
- Fundo `--paper` (mesmo do body), borda direita `1px solid --line`.
- **Brand header** (22px padding vertical, borda inferior): glyph `BrandMark` (gradiente `--accent` → `--accent-deep`, 30px quadrado com raio 7–8px, glifo SVG de folha + olho) + "Arandu" em serifa 20px + subtítulo "CLÍNICO" uppercase 10.5px.
- **Nav**: items com ícone 18px + label 13.5px, padding 9×12, radius 10.
  - Ativo: background `color-mix(in oklab, --accent 14%, transparent)`, border `color-mix(--accent 22%)`, cor `--accent-deep`, weight 500.
  - Inativo: transparente, cor `--ink-2`.
- Seções agrupadas com header "MENU" / "ATALHOS" (uppercase 10px letter-spacing 1.4px).
- **Rodapé**: avatar 32px + nome/CRP truncado + botão "Recolher" (dashed border 1px `--line-2`).

### Topbar (`chrome.jsx` → `Topbar`)
- 64px, sticky, backdrop-filter blur 8px, background `color-mix(--paper 88%, transparent)`.
- **Breadcrumb** à esquerda: separadores `>`, último item em serifa 15px weight 500.
- **Search** ao centro-direita: 340px, fundo `--paper-2`, border `--line`, ícone search 15px, placeholder "Buscar paciente, sessão ou observação…", kbd `⌘K`.
- **Botão LLM**: gradiente `--accent-deep → --accent`, texto "Arandu" + ícone sparkles + kbd `⌘J`.
- Botão settings (gear): `1px solid --line`, 9px padding, radius 10.

---

## Telas

### 1. Dashboard (`page_dashboard.jsx`)
**Propósito:** visão geral do dia para o clínico ao chegar.

**Layout:**
1. **Hero editorial** (grid 1fr auto, borda inferior):
   - Esquerda: eyebrow data atual (ex. "terça, 19 de abril") + H1 serifa "Bom dia, *Helena*." (itálico em "Helena", cor `--accent-deep`) + subtítulo resumindo o dia.
   - Direita: botão secundário "Abrir agenda" + botão primário preto "Nova sessão".
2. **KPIs** (grid de 4 colunas):
   - **Primeiro card em fundo escuro `--ink`, texto `--paper`**, com blob decorativo blur 30px no canto superior direito (`--accent` 40% opacidade).
   - Demais cards: `--paper-2` + border `--line`.
   - Cada card: eyebrow label + número 42px serifa + delta line com ícone (trend/alert).
3. **Grid principal** (1.1fr 1fr):
   - Coluna A: **Agenda do dia** (Card com lista de `ScheduleItem`s — grid `64px 1fr auto`, mono para horário, pill de status: "Concluída" verde, "Próxima" accent, "Agendada" neutral; items concluídos com opacity 0.55 e horário com line-through).
   - Coluna B stack:
     - **Pacientes ativos** (lista de 5, cada item: avatar 34px + nome + pill de tag + metadata "N sessões · última X · próxima Y"; clicável).
     - **Últimas sessões** (bloco lateral esquerdo 4px gradient `--accent` → fade + eyebrow patient+date+theme-pill + citação itálica serifa 14px).

### 2. Perfil do paciente (`page_patient.jsx`)
**Propósito:** visão 360° de um paciente.

**Layout:**
1. **Hero** (grid `auto 1fr auto`, padding 24px bottom, borda inferior):
   - Avatar 72px + iniciais serifa.
   - Centro: eyebrow "PACIENTE #P0446" + H1 serifa 40px (nome) + linha de metadata (idade, pronomes, marcadores sociais, since).
   - Direita: 3 `StatBlock`s separados por divisores verticais 1px × 40px — "5 Sessões", "2,3a Em terapia", "Quinzenal Frequência" (número serifa 26px + label uppercase).
2. **Corpo** (grid 2fr 1fr):
   - **Coluna principal:**
     - **Triagem** tratada como blockquote: border-left `2px solid --accent`, padding-left 20px, texto serifa 18px line-height 1.55.
     - **Timeline clínica**: filter pills ("Tudo / Sessões / Notas") + lista `TimelineItem`s (grid `96px auto 1fr auto` — data mono, dot 28px redondo preenchido accent para sessões/paper-2 para notas, título serifa + pill kind, summary 13px, chevron se clicável).
   - **Coluna lateral:**
     - **Ações** — 4 botões full-width stacked, primeiro primário preto "Nova sessão", resto `--paper` + border.
     - **Observações recentes** — cards pequenos `--paper` com pill de tag + timestamp mono + texto serifa 13.5px.

### 3. Sessão (`page_session.jsx`)
**Propósito:** registro clínico da sessão com observações e intervenções paralelas.

**Layout:**
1. **Cabeçalho** (grid `auto 1fr auto`, borda inferior):
   - Esquerda: botão "Voltar ao paciente" com seta.
   - Centro: eyebrow "Sessão 05 · André Barbosa" + H1 serifa 32px com uma palavra em itálico `--accent-deep` + metadata (data mono + duração + pill "Em rascunho").
   - Direita: botão "Ditar" (mic) + botão primário "Finalizar sessão" (check).
2. **Duas colunas paralelas** (grid 1fr 1fr) — cada `NotesColumn`:
   - Header: eyebrow ("Escuta" / "Ação") + título serifa + subtítulo.
   - Lista de notas: cada card `--paper` com texto serifa 14.5px + footer (clock ícone + timestamp mono + pill de tag + botão edit).
   - Input embutido no fundo (borda superior tracejada): eyebrow "Nova observação/intervenção" + textarea serifa itálico quando vazio + footer (hashtag pill + kbd `⌘↵ registrar`).
3. **Síntese** — card único, parágrafo serifa 17px line-height 1.65, com dois botões: primário gradient accent "Gerar síntese com Arandu" + secundário "Editar manualmente".

---

## Painel de Inteligência (LLM) — drawer (`llm_panel.jsx`)

**Trigger:** botão "Arandu" na topbar ou atalho `⌘J` / `Ctrl+J`.
**Animação:** slide from right, 460px wide, `transform: translateX(100% → 0)`, `.28s cubic-bezier(.4,0,.2,1)`. Backdrop com blur atrás.

**Estrutura:**
1. **Header**: glyph sparkles 36px gradient accent + "Arandu" serifa 18px + subtítulo dinâmico "Inteligência clínica · analisando *\<contexto\>*" + botão close.
2. **Padrões cruzados** (seção fundo `--paper-2`, borda tracejada): lista de 3 insights como barras de progresso — label uppercase + barra 4px arredondada com fill colorido (accent / gold / sage) proporcional ao peso + valor em serifa à direita.
3. **Thread**: mensagens user (balão `--paper-2` com radius 12/12/12/4) e assistant (ícone sparkles 26px quadrado + texto serifa 14px com markdown bold substituído por `<strong>` accent-deep + chips de citação "S1 · 29/01" mono 10.5px com ícone link).
4. **Sugestões rápidas** (chips serifa itálica, fundo `--paper-2`, radius 20).
5. **Input**: textarea + botão send preto. Disclaimer "Respostas de IA são auxiliares e não substituem julgamento clínico." 10.5px com ícone alert.

---

## Ícones
Todos os ícones são SVG inline 24×24 viewBox, stroke 1.5px (1.8 quando ativo), linecap/linejoin round, fill none. Ver `icons.jsx` para os 30+ paths. Biblioteca equivalente recomendada: **Lucide** (quase 1:1 com o que foi desenhado aqui). Use `lucide` ou `@heroicons/react/24/outline` no codebase real.

---

## Componentes reutilizáveis (mapa → implementação)

| Componente     | Descrição                                                          | Arquivo         |
|----------------|--------------------------------------------------------------------|-----------------|
| `Card`         | Superfície base: header (eyebrow, title serifa, subtitle, action) + body | chrome.jsx  |
| `Pill`         | Badge arredondado com 6 tones (neutral, warn, ok, info, danger, accent); usa `color-mix` para fundos translúcidos | chrome.jsx |
| `Kbd`          | Tecla de atalho mono com border-bottom 2px                         | chrome.jsx      |
| `Avatar`       | Círculo com iniciais em serifa, tamanho configurável, gradiente    | chrome.jsx      |
| `BrandMark`    | Glyph quadrado com gradient accent + SVG folha-olho                | chrome.jsx      |
| `StatBlock`    | Número serifa grande + label uppercase                             | page_patient    |
| `ScheduleItem` | Item de agenda com mono time + patient + pill                      | page_dashboard  |
| `TimelineItem` | Item de timeline com dot + título + summary                        | page_patient    |
| `NotesColumn`  | Coluna de notas com header + lista + input                         | page_session    |

---

## Interações e comportamento

- **Navegação**: estado `page` persistido em `localStorage['arandu:page']`. No protótipo é state React; no HTMX deve virar rota real (`/dashboard`, `/patients/:id`, `/sessions/:id`).
- **Sidebar recolhível**: toggle muda largura 232 ↔ 68px, transition 0.2s ease. Items viram só ícones centralizados, tooltip no title.
- **LLM drawer**: abre/fecha via ⌘J global listener + botão. Backdrop clicável fecha. Transform transition.
- **Hover states**: botões de nav ativos com background accent-14%; items de lista com leve background `--paper-3` no hover (adicionar se não houver ainda).
- **Filters**: pill ativa preta `--ink` / inativa transparente border `--line`.
- **Input textareas**: sem borda, fundo transparente, família serifa quando é conteúdo clínico, italic quando vazio.

---

## Estados a implementar (HTMX)

- **Sessão em rascunho vs finalizada**: pill muda "Em rascunho" (ok tone) → "Finalizada" (neutral); botões de ação viram read-only.
- **LLM loading**: streaming de resposta do assistant (mensagem aparece progressivamente) — usar `hx-sse` ou equivalente.
- **Observações pendentes**: KPI `pending` com tone warn + ícone alert quando > 0.
- **Próxima sessão**: item da agenda com pill accent "Próxima" — derivar de comparação com hora atual.

---

## Content / copy (português BR)
Todo o conteúdo está em pt-BR. Copywriting é neutro e técnico mas acolhedor. Termos-chave:
- "Inteligência clínica" (não "IA")
- "Reflexão terapêutica", "Escuta", "Percurso", "Síntese"
- Eyebrows em UPPERCASE, letter-spacing generoso
- Evitar linguagem de CRM ("leads", "conversões")

---

## Assets
- Google Fonts: Fraunces, Geist, Geist Mono (ver `<link>` em `Arandu Redesign.html`)
- Ícones: todos inline SVG (ver `icons.jsx`). Recomendação de produção: migrar para Lucide/Heroicons.
- Nenhuma imagem raster — tudo vetorial.

---

## Arquivos neste pacote

| Arquivo                    | Conteúdo                                                    |
|----------------------------|-------------------------------------------------------------|
| `Arandu Redesign.html`     | Shell HTML + CSS tokens + Google Fonts                      |
| `app.jsx`                  | App root: roteamento entre páginas, LLM toggle, atalhos     |
| `data.jsx`                 | Dados mockados (pacientes, sessões, KPIs, LLM threads)      |
| `icons.jsx`                | Biblioteca de ícones SVG inline                             |
| `chrome.jsx`               | Sidebar, Topbar, Card, Pill, Avatar, BrandMark, Kbd         |
| `page_dashboard.jsx`       | Dashboard                                                   |
| `page_patient.jsx`         | Perfil do paciente                                          |
| `page_session.jsx`         | Sessão (observações + intervenções)                         |
| `llm_panel.jsx`            | Drawer de Inteligência Clínica                              |
| `tweaks.jsx`               | Painel de Tweaks (apenas protótipo — não portar)            |

---

## Dicas de implementação no stack HTMX + Tempo + Tailwind

1. **Tokens** → extenda `tailwind.config.js` com as cores/fontes acima. Use nomes semânticos (`paper`, `ink`, `accent`) em vez de `orange-600`.
2. **Serifa em títulos** → classe `font-serif` (Fraunces) + `tracking-tight`.
3. **Layouts** → CSS Grid onde descrito; Tailwind `grid-cols-[1.1fr_1fr]` funciona.
4. **Sidebar** → parcial server-rendered, toggle via HTMX + `hx-swap-oob` ou simples classe via Alpine.
5. **LLM drawer** → conteúdo server-streamed via `hx-sse`; trigger com `hx-get` + `hx-target="#llm-drawer"` + `hx-swap="innerHTML"` + classes Tailwind para slide-in.
6. **Atalhos de teclado** → Alpine.js ou listener JS curto global (`@keydown.meta.j.window`).
7. **Persistência de página** → não é necessária se as URLs forem reais (voltar no navegador já resolve).
8. **Remover**: painel Tweaks, FAB de settings (ou converter gear em menu real do usuário).

## Como rodar o protótipo localmente
Basta abrir `Arandu Redesign.html` em um servidor HTTP qualquer (ex. `npx serve .` ou `python -m http.server`). Os `.jsx` são carregados via `<script type="text/babel">` e transpilados pelo Babel standalone in-browser — sem build step necessário.
