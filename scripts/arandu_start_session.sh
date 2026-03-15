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

echo "✅ Contexto obrigatório carregado"