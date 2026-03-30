#!/usr/bin/env bash

TASK_TITLE="${1:-}"
REQ_ID="${2:-}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

MODE="auto"

if [ -z "$TASK_TITLE" ]; then
  MODE="no-task"
  echo "🎯 Cycle de Validação (sem tarefa)"
  echo "===================================="
  echo ""
else
  MODE="with-task"
  echo "🎯 Cycle de Tarefa Completo"
  echo "==========================="
  echo ""
  echo "Tarefa: $TASK_TITLE"
  [ -n "$REQ_ID" ] && echo "Requirement: $REQ_ID"
  echo ""
fi

if [ "$MODE" = "with-task" ]; then
  TASK_ID=$(date +%Y%m%d_%H%M%S)

  echo "1️⃣  Criando tarefa..."
  bash "$SCRIPT_DIR/arandu_create_task.sh" "$TASK_TITLE" $REQ_ID
  echo ""

  echo "📝 INSTRUÇÕES:"
  echo "=============="
  echo "1. Edite o arquivo de tarefa criado"
  echo "2. Remova a seção 'Instruções para o Agente'"
  echo "3. Atualize o status para 'PRONTO_PARA_IMPLEMENTACAO'"
  echo ""
  read -p "Pressione ENTER quando a tarefa estiver pronta..."

  echo ""
  echo "4️⃣  Processando tarefa..."
  bash "$SCRIPT_DIR/arandu_process_task.sh" "$TASK_ID"
  echo ""
fi

echo "2️⃣  Validando integridade (checkpoint)..."
echo "==========================================="
bash "$SCRIPT_DIR/arandu_checkpoint.sh"
if [ $? -ne 0 ]; then
  echo "⚠️ _checkpoint falhou. Corrija os erros antes de continuar."
  exit 1
fi
echo ""

echo "3️⃣  Verificando guard (rotas + templ generation)..."
echo "=================================================="
bash "$SCRIPT_DIR/arandu_guard.sh"
if [ $? -ne 0 ]; then
  echo "⚠️  Guard falhou. Corrija os erros antes de continuar."
  exit 1
fi
echo ""

echo "4️⃣  Validação visual automatizada..."
echo "====================================="
bash "$SCRIPT_DIR/arandu_visual_all.sh"
if [ $? -ne 0 ]; then
  echo "⚠️  Validação visual falhou."
  echo "   Dica: Inicie o servidor antes: ./scripts/safe_deploy.sh"
  exit 1
fi
echo ""

if [ "$MODE" = "with-task" ]; then
  echo "5️⃣  Concluindo tarefa..."
  echo "======================="
  bash "$SCRIPT_DIR/arandu_conclude_task.sh" "$TASK_ID" --success
  echo ""
fi

echo "================================"
echo "✅ VALIDAÇÃO COMPLETA PASSOU"
echo "================================"
