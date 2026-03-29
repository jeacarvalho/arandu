# Migração para Tailwind CSS

## Contexto do Projeto

- **Stack**: Go + templ + Alpine.js + HTMX
- **Baseline**: Tag `v1.1.0-tailwind-migration-base`
- **CSS atual**: 7.096 linhas de CSS customizado (`web/static/css/style.css`)
- **Tailwind**: Configuração parcial existente em `design/tailwind.config.js`
- **Componentes**: ~37 arquivos `.templ` para migrar

## Objetivo

Migrar toda a interface do projeto de CSS customizado para Tailwind CSS, mantendo a identidade visual botânica existente.

## Design System Atual (a preservar)

### Cores
- Primary: `#0F6E56` (Verde Base)
- Active: `#1D9E75` (Destaque)
- Soft: `#9FE1CB` (Acentos)
- Background: `#E1F5EE` (Papel de Seda)
- Dark: `#085041` (Verde Floresta)
- Accent: `#F59E0B` (Âmbar para IA)

### Tipografia
- Interface: Inter (sans)
- Conteúdo clínico: Source Serif 4 (serif)

---

## Tasks

### Fase 1: Setup e Configuração
- [x] 1.1 Configurar Tailwind CSS completo no build (PostCSS)
- [x] 1.2 Configurar variáveis do design system (via @theme no Tailwind v4)
- [x] 1.3 Integrar Tailwind no layout base (substituir style.css)

### Fase 2: Componentes de Layout
- [x] 2.1 Layout base (`layout.templ` - sidebar, topbar)
- [x] 2.2 Classes utilitárias de estrutura (flex, grid, spacing) - nativas do Tailwind
- [x] 2.3 Componentes de autenticação (login)

### Fase 3: Componentes de Interface
- Manter style.css para componentes detalhados (cards, botões, formulários)
- Usar Tailwind apenas para estrutura e utilitários

### Fase 4: Páginas (migração gradual)
- [ ] 4.1 Dashboard - migração de estrutura
- [ ] 4.2 Patient list / detail / new - migração de estrutura
- [ ] 4.3 Session (list, detail, edit) - migração de estrutura
- [ ] 4.4 Anamnese / Biopsychosocial - migração de estrutura
- [ ] 4.5 AI Insights - migração de estrutura

### Fase 5: Limpeza
- [ ] 5.1 Remover CSS não utilizado
- [ ] 5.2 Testar todas as funcionalidades
- [ ] 5.3 Criar tag de versão final

---

## Estado Atual

**Fase 4 pausada** - Migração de Páginas

### O que foi migrado:
- ✅ Layout base (sidebar, topbar) com Tailwind
- ✅ Tailwind CSS integrado ao build

### O que foi aprendido:
- Migrar componente a componente causa quebra visual
- style.css deve ser preservado para componentes detalhados

### Próximos passos sugeridos:
- Usar Tailwind apenas para estrutura (flex, grid, spacing)
- Manter style.css para cards, botões, formulários
- NÃO migrar páginas até ter testes visuais adequados

Tasks concluídas:
- 1.1: Configurar Tailwind CSS v4.2.2 no build
- 1.2: Configurar variáveis do design system via @theme
- 1.3: Integrar Tailwind no layout base
- 2.1: Layout base (sidebar, topbar migrados para Tailwind)
- 2.2: Classes utilitárias (Tailwind nativo)
- 2.3: Componentes de autenticação (login migrado)

Próximo passo: **Fase 3 - Task 3.1** - Cards (patient-card, stat-card)

---

## Comandos Úteis

```bash
# Verificar estado atual
git status

# Verificar tag atual
git describe --tags

# Criar tag ao final de cada fase
git tag -a v1.x.x -m "Descrição"
```
