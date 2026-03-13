#!/usr/bin/env bash

set -e

SESSION_DIR="work/current_session"
TASKS_DIR="work/tasks"
ARCHIVE_DIR="work/archive"

echo "🔚 Encerrando sessão Arandu..."
echo

# Verificar se sessão existe
if [ ! -d "$SESSION_DIR" ]; then
  echo "❌ Nenhuma sessão ativa encontrada."
  echo "💡 Execute primeiro: scripts/arandu_start_session.sh"
  exit 1
fi

# Carregar informações da sessão
SESSION_INFO="${SESSION_DIR}/session_info"

if [ -f "$SESSION_INFO" ]; then
  source "$SESSION_INFO"
else
  SESSION_ID=$(date +%Y%m%d_%H%M%S)
  START_TIME=$(date +%s)
fi

echo "📋 Sessão: $SESSION_ID"

# Verificar tarefas pendentes
if [ -d "$TASKS_DIR" ]; then
  PENDING=$(ls "$TASKS_DIR" 2>/dev/null | wc -l)

  if [ "$PENDING" -gt 0 ]; then
    echo
    echo "❌ Existem tarefas pendentes:"
    ls "$TASKS_DIR"
    echo
    echo "💡 Conclua tarefas antes de encerrar sessão"
    exit 1
  fi
fi

# Atualizar contexto final da sessão
if [ -f "scripts/arandu_update_context.sh" ]; then
  echo
  echo "🧠 Gerando snapshot final do projeto..."
  bash scripts/arandu_update_context.sh
fi


# Calcular duração
END_TIME=$(date +%s)
if [ -n "$START_TIME" ]; then
  DURATION=$((END_TIME - START_TIME))
  HOURS=$((DURATION / 3600))
  MINUTES=$(( (DURATION % 3600) / 60 ))
  SECONDS=$((DURATION % 60))
  
  echo
  echo "⏱️  Duração da sessão: ${HOURS}h ${MINUTES}m ${SECONDS}s"
fi

# Limpar diretório da sessão
rm -rf "$SESSION_DIR"

echo
echo "✅ Sessão encerrada com sucesso!"
echo "📁 Arquivos da sessão movidos para: $ARCHIVE_DIR"
