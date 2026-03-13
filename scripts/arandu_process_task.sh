#!/usr/bin/env bash

# Verificar argumentos
if [ $# -lt 1 ]; then
  echo "Uso: $0 TASK_ID [--execute]"
  echo "Exemplo: $0 20260313_185952"
  echo "Exemplo: $0 20260313_185952 --execute"
  exit 1
fi

TASK_ID=$1
EXECUTE=${2:-""}

TASK_DIR="work/tasks/task_${TASK_ID}"

if [ ! -d "$TASK_DIR" ]; then
  echo "❌ Tarefa não encontrada: $TASK_ID"
  echo "   Verifique se o diretório existe: $TASK_DIR"
  exit 1
fi

echo "📋 Tarefa: ${TASK_ID}"
echo "📁 Diretório: ${TASK_DIR}"
echo ""

# Mostrar conteúdo da tarefa
if [ -f "$TASK_DIR/task.md" ]; then
  cat "$TASK_DIR/task.md"
else
  echo "⚠️  Arquivo task.md não encontrado na tarefa"
fi

echo ""
echo "⚠️  Antes de implementar, leia o requirement correspondente"

# Se o parâmetro --execute foi fornecido
if [ "$EXECUTE" = "--execute" ]; then
  echo ""
  echo "🚀 Modo execute ativado"
  echo "   Esta funcionalidade executaria a tarefa automaticamente."
  echo "   (Implementação futura)"
fi
