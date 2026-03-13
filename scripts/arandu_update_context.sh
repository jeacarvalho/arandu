#!/usr/bin/env bash

set -e

echo "🧠 Atualizando contexto do agente para o projeto Arandu..."
echo

CONTEXT_DIR="docs/agent_context"
OUTPUT_FILE="${CONTEXT_DIR}/project_state.md"

mkdir -p "$CONTEXT_DIR"

echo "# 🧠 Estado Atual do Projeto Arandu" > "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "**Gerado em:** $(date)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ------------------------------------------------
# VISÃO DO SISTEMA
# ------------------------------------------------

echo "## 🌳 Estrutura de Visões" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [ -d "docs/vision" ]; then
  for V in docs/vision/*.md; do
    [ -f "$V" ] || continue
    echo "- $(basename "$V")" >> "$OUTPUT_FILE"
  done
else
  echo "Nenhuma visão encontrada." >> "$OUTPUT_FILE"
fi

echo "" >> "$OUTPUT_FILE"

# ------------------------------------------------
# CAPABILITIES
# ------------------------------------------------

echo "## ⚙️ Capabilities" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [ -d "docs/capabilities" ]; then
  for C in docs/capabilities/*.md; do
    [ -f "$C" ] || continue
    echo "- $(basename "$C")" >> "$OUTPUT_FILE"
  done
else
  echo "Nenhuma capability encontrada." >> "$OUTPUT_FILE"
fi

echo "" >> "$OUTPUT_FILE"

# ------------------------------------------------
# REQUIREMENTS
# ------------------------------------------------

echo "## 📋 Requirements" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

REQ_TOTAL=0

if [ -d "docs/requirements" ]; then
  for R in docs/requirements/*.md; do
    [ -f "$R" ] || continue
    echo "- $(basename "$R")" >> "$OUTPUT_FILE"
    REQ_TOTAL=$((REQ_TOTAL+1))
  done
fi

echo "" >> "$OUTPUT_FILE"
echo "**Total de requirements:** $REQ_TOTAL" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ------------------------------------------------
# TAREFAS RECENTES
# ------------------------------------------------

echo "## 🛠️ Tarefas recentes" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [ -d "work/archive" ]; then
  ls -t work/archive | head -5 | while read SESSION
  do
    echo "- $SESSION" >> "$OUTPUT_FILE"
  done
else
  echo "Nenhuma sessão arquivada ainda." >> "$OUTPUT_FILE"
fi

echo "" >> "$OUTPUT_FILE"

# ------------------------------------------------
# APRENDIZADOS RECENTES
# ------------------------------------------------

echo "## 📚 Aprendizados recentes" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [ -d "docs/learnings" ]; then
  ls -t docs/learnings | head -5 | while read FILE
  do
    echo "- $FILE" >> "$OUTPUT_FILE"
  done
else
  echo "Nenhum aprendizado registrado." >> "$OUTPUT_FILE"
fi

echo "" >> "$OUTPUT_FILE"

# ------------------------------------------------
# STATUS DO SISTEMA
# ------------------------------------------------

echo "## 📊 Status do Sistema" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "- Visions: $(ls docs/vision 2>/dev/null | wc -l)" >> "$OUTPUT_FILE"
echo "- Capabilities: $(ls docs/capabilities 2>/dev/null | wc -l)" >> "$OUTPUT_FILE"
echo "- Requirements: $(ls docs/requirements 2>/dev/null | wc -l)" >> "$OUTPUT_FILE"
echo "- Learnings: $(ls docs/learnings 2>/dev/null | wc -l)" >> "$OUTPUT_FILE"

echo "" >> "$OUTPUT_FILE"

echo "✅ Contexto do projeto atualizado:"
echo "$OUTPUT_FILE"
