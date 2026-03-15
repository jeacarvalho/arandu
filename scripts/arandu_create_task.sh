#!/usr/bin/env bash

set -e

# Verificar argumentos
if [ $# -lt 1 ]; then
  echo "Uso: $0 \"titulo da tarefa\" [REQ-ID]"
  echo "Exemplo 1: $0 \"Criando as fundações\""
  echo "Exemplo 2: $0 \"Implementar feature X\" req-01-01-01"
  exit 1
fi

# Se apenas um argumento foi fornecido, é o título
if [ $# -eq 1 ]; then
  TITLE=$1
  REQ_ID=""
else
  TITLE=$1
  REQ_ID=$2
fi

TASK_ID=$(date +%Y%m%d_%H%M%S)
TASK_DIR="work/tasks/task_${TASK_ID}"

mkdir -p "$TASK_DIR"

# Criar arquivo da tarefa
if [ -n "$REQ_ID" ]; then
  cat > "$TASK_DIR/task.md" <<EOF
# TASK ${TASK_ID}

Requirement: ${REQ_ID}

Title: ${TITLE}

## Objetivo

Implementar requirement ${REQ_ID}

## Referências

docs/requirements/${REQ_ID}.md
EOF
else
  cat > "$TASK_DIR/task.md" <<EOF
# TASK ${TASK_ID}

Title: ${TITLE}

## Objetivo

${TITLE}

## Descrição

Tarefa criada sem requirement específico.

## Instruções para o Agente

1. Aguarde o usuário fornecer detalhes da tarefa neste arquivo
2. Quando o usuário editar este arquivo com os detalhes, inicie a implementação
3. Siga o padrão de referenciar requirements quando aplicável
EOF
fi

echo "✅ tarefa criada: ${TASK_ID}"
echo "📁 Tarefa em: $TASK_DIR/task.md"
echo ""
echo "📝 Instruções:"
echo "1. Edite $TASK_DIR/task.md para detalhar a tarefa"
echo "2. O agente aguardará sua descrição antes de iniciar"
