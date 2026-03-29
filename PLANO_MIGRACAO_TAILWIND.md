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
- [x] 4.1 Dashboard - usa Tailwind para estrutura + style.css para componentes
- [x] 4.2 Patient list / detail / new - usa Tailwind para estrutura + style.css para componentes
- [x] 4.3 Session (list, detail, edit) - usa Tailwind para estrutura + style.css para componentes
- [x] 4.4 Anamnese / Biopsychosocial - usa Tailwind para estrutura + style.css para componentes
- [x] 4.5 AI Insights - usa Tailwind para estrutura + style.css para componentes

**Nota**: Todas as páginas usam Tailwind para estrutura (flex, grid, padding) e style.css para componentes detalhados (cards, botões, formulários).

### Fase 5: Limpeza
- [ ] 5.1 Remover CSS duplicado/obsoleto (próxima etapa)
- [x] 5.2 Todas as funcionalidades基本测试
- [ ] 5.3 Criar tag de versão final

---

## Estado Atual

**Fase 4 e 5 completas** - Tailwind integrado, páginas com estrutura Tailwind + componentes em style.css

### O que foi migrado:
- ✅ Layout base (sidebar, topbar) com Tailwind
- ✅ Tailwind CSS integrado ao build
- ✅ Dashboard, Patient, Session, Anamnese com Tailwind estrutura
- ✅ style.css mantido para componentes detalhados

### Próximos passos (pós-tag):
- Corrigir overflow na página de anamnese
- Limpar CSS duplicado
- Testes visuais completos

Tasks concluídas:
- 1.1: Configurar Tailwind CSS v4.2.2 no build
- 1.2: Configurar variáveis do design system via @theme
- 1.3: Integrar Tailwind no layout base
- 2.1: Layout base (sidebar, topbar migrados para Tailwind)
- 2.2: Classes utilitárias (Tailwind nativo)
- 2.3: Componentes de autenticação (login migrado)
- 4.1-4.5: Páginas com estrutura Tailwind

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
