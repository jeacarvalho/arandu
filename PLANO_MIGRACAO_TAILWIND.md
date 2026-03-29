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
- [ ] 1.3 Integrar Tailwind no layout base (substituir style.css)

### Fase 2: Componentes de Layout
- [ ] 2.1 Layout base (`layout.templ` - sidebar, topbar)
- [ ] 2.2 Classes utilitárias de estrutura (flex, grid, spacing)
- [ ] 2.3 Componentes de autenticação (login)

### Fase 3: Componentes de Interface
- [ ] 3.1 Cards (patient-card, stat-card, clinical-card)
- [ ] 3.2 Botões (btn-primary, btn-secondary, btn-ghost)
- [ ] 3.3 Formulários (form-control, form-label, form-error)
- [ ] 3.4 Timeline
- [ ] 3.5 Search e dropdowns

### Fase 4: Páginas
- [ ] 4.1 Dashboard
- [ ] 4.2 Patient list / detail / new
- [ ] 4.3 Session (list, detail, edit)
- [ ] 4.4 Anamnese / Biopsychosocial
- [ ] 4.5 AI Insights

### Fase 5: Limpeza
- [ ] 5.1 Remover CSS não utilizado
- [ ] 5.2 Testar todas as funcionalidades
- [ ] 5.3 Criar tag de versão final

---

## Estado Atual

**Fase 1 em andamento** - Tailwind CSS configurado e compilando

Última task concluída:
- 1.1: Configurar Tailwind CSS v4.2.2 no build
- 1.2: Configurar variáveis do design system via @theme

Próximo passo: **Fase 1 - Task 1.3** - Integrar Tailwind no layout base

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
