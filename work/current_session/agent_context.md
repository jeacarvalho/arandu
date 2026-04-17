# CONTEXTO DO AGENTE — ARANDU

Sessão: 20260416_145127

## PASSOS OBRIGATÓRIOS

Antes de qualquer implementação leia:

1 docs/dvp.md
2 docs/vision/
3 docs/capabilities/
4 docs/requirements/
 5 docs/learnings/MASTER_LEARNINGS.md (sistema consolidado)

# CONTEXTO CRÍTICO — ARANDU SOTA

## 🛡️ LEIS DE PROTEÇÃO (NÃO NEGOCIÁVEIS)
1. NUNCA crie arquivos .html soltos. Use componentes .templ.
2. TODA página deve herdar de templates.Layout().
3. CONTEÚDO CLÍNICO deve usar obrigatoriamente .font-clinical (Source Serif 4).
4. ROTAS EXISTENTES não podem quebrar. Verifique /patients e /sessions antes de concluir.
5. 🚨 NUNCA use hx-swap-oob com IDs que não existem no DOM base (veja docs/learnings/MASTER_LEARNINGS.md seção 0)

## PASSOS OBRIGATÓRIOS
Leia antes de qualquer código:
- architecture_sota.md (Padrões de backend e DB)
- interface_patterns_sota.md (Padrões de UI e UX)
- docs/requirements/ (O requirement da tarefa)

## ARQUITETURA WEB (PR #1 INTEGRADO)

Para implementações na camada web, CONSULTE OBRIGATORIAMENTE:

6 docs/architecture/WEB_LAYER_PATTERN.md
7 docs/architecture/system_structure.md
8 docs/architecture/ROUTE_CONVENTIONS.md
9 docs/architecture/AGENT_GUIDE.md (guia prático para agentes)

Referências de código modelo:
- internal/web/handlers/patient_handler.go
- internal/web/handlers/session_handler.go
- web/templates/patients.html

## Regra de implementação

Toda tarefa deve referenciar um REQUIREMENT.

Formato:

REQ-XX-YY-ZZ

# Adicionar no agent_context.md gerado
cat >> "work/current_session/agent_context.md" <<EOF

⚠️  ATENÇÃO: VALIDAÇÃO VISUAL OBRIGATÓRIA
Antes de marcar tarefa como concluída:
1. Rode: ./scripts/arandu_visual_check.sh
2. Rode: ./scripts/arandu_screenshot_compare.sh
3. Revise screenshots em ./screenshots/current/
4. Preencha checklist em docs/visual_validation_checklist.md
5. Só então execute: ./scripts/arandu_conclude_task.sh

