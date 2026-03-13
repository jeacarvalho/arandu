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
5 docs/learnings/

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