#!/usr/bin/env bash

set -e

SESSION_ID=$(date +%Y%m%d_%H%M%S)

SESSION_DIR="work/current_session"

echo "🚀 Iniciando sessão Arandu ${SESSION_ID}"

mkdir -p "$SESSION_DIR"
mkdir -p work/tasks
mkdir -p work/archive

echo "SESSION_ID=${SESSION_ID}" > "$SESSION_DIR/session_info"

cat > "$SESSION_DIR/agent_context.md" <<EOF
# CONTEXTO DO AGENTE — ARANDU

Sessão: ${SESSION_ID}

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
EOF

echo "✅ Sessão criada"

# Atualizar contexto do projeto
if [ -f "scripts/arandu_update_context.sh" ]; then
  echo
  echo "🧠 Gerando contexto inicial do projeto..."
  bash scripts/arandu_update_context.sh
fi

echo
echo "📚 Carregando contexto obrigatório do projeto..."

# Ler o project_state.md
if [ -f "docs/agent_context/project_state.md" ]; then
  echo "📄 Lendo project_state.md..."
  cat docs/agent_context/project_state.md
  echo
fi

# Ler todos os arquivos listados no project_state.md
echo "📚 Lendo arquivos de contexto obrigatório..."

# Ler vision files
for vision_file in docs/vision/vision-*.md; do
  if [ -f "$vision_file" ]; then
    echo "📄 $vision_file"
    head -20 "$vision_file"
    echo "..."
  fi
done

# Ler capability files
for cap_file in docs/capabilities/cap-*.md; do
  if [ -f "$cap_file" ]; then
    echo "📄 $cap_file"
    head -10 "$cap_file"
    echo "..."
  fi
done

# Ler requirement files
for req_file in docs/requirements/req-*.md; do
  if [ -f "$req_file" ]; then
    echo "📄 $req_file"
    head -10 "$req_file"
    echo "..."
  fi
done

# Ler arquitetura (nova seção)
echo
echo "🏗️  Lendo documentação de arquitetura..."
for arch_file in docs/architecture/*.md; do
  if [ -f "$arch_file" ]; then
    echo "📄 $arch_file"
    head -5 "$arch_file"
    echo "..."
  fi
done

echo "✅ Contexto obrigatório carregado"