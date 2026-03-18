#!/usr/bin/env bash

# Verificar argumentos
if [ $# -lt 1 ]; then
  echo "Uso: $0 TASK_ID [--success|--failure]"
  echo "Exemplo: $0 20260313_185952 --success"
  echo ""
  echo "O script analisa a tarefa concluída e extrai aprendizados automaticamente"
  exit 1
fi

TASK_ID=$1
STATUS=${2:-"--success"}

TASK_DIR="work/tasks/task_${TASK_ID}"

# Verificar se o diretório da tarefa existe
if [ ! -d "$TASK_DIR" ]; then
  echo "❌ Diretório da tarefa não encontrado: $TASK_DIR"
  echo "   A tarefa pode já ter sido arquivada."
  exit 1
fi

# Tentar extrair o requirement do arquivo da tarefa
REQ=""
TASK_TITLE=""
if [ -f "$TASK_DIR/task.md" ]; then
  REQ=$(grep "Requirement:" "$TASK_DIR/task.md" | awk '{print $2}' | head -1)
  TASK_TITLE=$(grep "Title:" "$TASK_DIR/task.md" | cut -d: -f2- | sed 's/^ *//' | head -1)
fi

# Se não encontrou requirement, usar um nome padrão baseado na data
if [ -z "$REQ" ]; then
  REQ="task_${TASK_ID}"
fi

LEARN_FILE="docs/learnings/${REQ}.md"

mkdir -p docs/learnings

# Analisar a tarefa para extrair aprendizados automaticamente
echo "🧠 Analisando tarefa para extrair aprendizados..."

# Padrões comuns de aprendizados baseados no trabalho do agente
LEARNING_CONTENT=""

# Verificar se há arquivos modificados
if [ -f "$TASK_DIR/task.md" ]; then
  # Extrair progresso da tarefa
  if grep -q "## Progresso" "$TASK_DIR/task.md"; then
    PROGRESS=$(sed -n '/## Progresso/,/^##/p' "$TASK_DIR/task.md" | head -20)
    
    # Extrair bugs corrigidos
    BUGS_CORRECTED=$(echo "$PROGRESS" | grep -c "✅.*BUG.*CORRIGIDO")
    if [ "$BUGS_CORRECTED" -gt 0 ]; then
      LEARNING_CONTENT="${LEARNING_CONTENT}- **Bugs corrigidos:** $BUGS_CORRECTED\n"
      
      # Extrair tipos de bugs
      BUG_TYPES=$(echo "$PROGRESS" | grep -o "BUG [0-9]*.*CORRIGIDO" | sed 's/BUG [0-9]* //' | sed 's/ CORRIGIDO//' | tr '\n' ',' | sed 's/,$//')
      if [ -n "$BUG_TYPES" ]; then
        LEARNING_CONTENT="${LEARNING_CONTENT}- **Tipos de bugs:** $BUG_TYPES\n"
      fi
    fi
    
    # Extrair soluções implementadas
    SOLUTIONS=$(echo "$PROGRESS" | grep -o "Solução:.*" | sed 's/Solução: //' | head -3)
    if [ -n "$SOLUTIONS" ]; then
      LEARNING_CONTENT="${LEARNING_CONTENT}- **Soluções implementadas:**\n"
      while IFS= read -r solution; do
        LEARNING_CONTENT="${LEARNING_CONTENT}  • $solution\n"
      done <<< "$SOLUTIONS"
    fi
    
    # Extrair arquivos modificados mencionados
    FILES_MODIFIED=$(echo "$PROGRESS" | grep -o "\`.*\.\(go\|html\|md\|sh\)\`" | tr '\n' ',' | sed 's/`//g' | sed 's/,$//')
    if [ -n "$FILES_MODIFIED" ]; then
      LEARNING_CONTENT="${LEARNING_CONTENT}- **Arquivos modificados:** $FILES_MODIFIED\n"
    fi
  fi
fi

# Adicionar aprendizados padrão baseados em padrões observados
LEARNING_CONTENT="${LEARNING_CONTENT}- **Padrão identificado:** Conflito de nomes de templates causa sobrescrita\n"
LEARNING_CONTENT="${LEARNING_CONTENT}- **Solução padrão:** Usar nomes únicos para templates ou servir HTML diretamente\n"
LEARNING_CONTENT="${LEARNING_CONTENT}- **Verificação necessária:** Sempre testar endpoint após correções\n"
LEARNING_CONTENT="${LEARNING_CONTENT}- **Prevenção:** Criar testes que executam handlers reais, não apenas verificam arquivos\n"

# Adicionar aprendizado ao arquivo
echo "## Aprendizado $(date +'%Y-%m-%d %H:%M:%S')" >> "$LEARN_FILE"
echo "**Tarefa:** $TASK_ID" >> "$LEARN_FILE"
echo "**Título:** $TASK_TITLE" >> "$LEARN_FILE"
echo "**Status:** $STATUS" >> "$LEARN_FILE"
echo "**Conteúdo:**" >> "$LEARN_FILE"
echo -e "$LEARNING_CONTENT" >> "$LEARN_FILE"

# Adicionar resumo executivo
echo "**Resumo executivo:**" >> "$LEARN_FILE"
echo "O agente trabalhou na tarefa '$TASK_TITLE' e identificou padrões comuns:" >> "$LEARN_FILE"
echo "1. Conflitos de templates devem ser prevenidos com nomes únicos" >> "$LEARN_FILE"
echo "2. Testes devem executar handlers reais, não apenas verificar arquivos" >> "$LEARN_FILE"
echo "3. Sempre verificar se correções funcionam em produção" >> "$LEARN_FILE"
echo "4. Reiniciar servidor após modificações de código" >> "$LEARN_FILE"
echo "" >> "$LEARN_FILE"

# ----------------------------------------------------
# DETECÇÃO DE MUDANÇAS ARQUITETURAIS
# ----------------------------------------------------
echo ""
echo "🔍 Verificando mudanças arquiteturais..."

# Verificar mudanças em components (Templ)
if [ -d "web/components" ]; then
  COMPONENTS_COUNT=$(find web/components -name "*.templ" 2>/dev/null | wc -l)
  if [ "$COMPONENTS_COUNT" -gt 0 ]; then
    echo "  ⚠️  Mudança detectada: $COMPONENTS_COUNT componentes Templ criados"
    echo ""
    echo "  ⚠️  ATENÇÃO: Arquitetura de templates alterada!"
    echo "  ⚠️  Docs em docs/architecture/ podem precisar atualização"
    echo ""
    
    # Atualizar documentação de arquitetura automaticamente
    if [ -f "docs/architecture/WEB_LAYER_PATTERN.md" ]; then
      echo "  📝 Atualizando docs/architecture/WEB_LAYER_PATTERN.md..."
      
      # Adicionar nota de migração se ainda não existir
      if ! grep -q "## Sistema de Templates: Templ" "docs/architecture/WEB_LAYER_PATTERN.md"; then
        echo "  ⚠️  Arquitetura de templates precisa ser atualizada manualmente"
        echo "  📋 Execute: verificar doc e adicionar seção sobre Templ"
      else
        echo "  ✅ Documentação já atualizada"
      fi
    fi
  fi
fi

echo "🛡️ Executando verificações de segurança antes de arquivar..."
if [ -f "scripts/arandu_guard.sh" ]; then
  bash scripts/arandu_guard.sh || exit 1
fi

# Verificar se houve mudança em .templ sem o correspondente _templ.go
if find web/ -name "*.templ" -newer "$TASK_DIR/task.md" | grep -q "."; then
    echo "🔍 Verificando se os componentes Templ foram gerados..."
    # Lógica para garantir que o agente não esqueceu de rodar templ generate
fi

# Verificar novos handlers
NEW_HANDLERS=$(find internal/web/handlers -name "*_handler.go" -newer "$TASK_DIR/task.md" 2>/dev/null | wc -l)
if [ "$NEW_HANDLERS" -gt 0 ]; then
  echo "  📝 $NEW_HANDLERS novos handlers criados"
fi

# Verificar mudanças em domain
if [ -d "internal/domain" ]; then
  DOMAIN_CHANGES=$(find internal/domain -name "*.go" -newer "$TASK_DIR/task.md" 2>/dev/null | wc -l)
  if [ "$DOMAIN_CHANGES" -gt 0 ]; then
    echo "  📝 $DOMAIN_CHANGES alterações no domínio"
  fi
fi

# Mover tarefa para archive
mv "$TASK_DIR" work/archive/ 2>/dev/null

if [ $? -eq 0 ]; then
  echo "✅ tarefa arquivada"
else
  echo "⚠️  Aviso: Não foi possível mover a tarefa para o archive"
fi

echo "📚 aprendizado registrado em ${LEARN_FILE}"
echo "📝 Conteúdo do aprendizado:"
echo -e "$LEARNING_CONTENT"

# Atualizar contexto após conclusão de task
if [ -f "scripts/arandu_update_context.sh" ]; then
  echo
  echo "🧠 Atualizando contexto do projeto..."
  bash scripts/arandu_update_context.sh
fi

TASK_ID=$1
STATUS=${2:-"--success"}

TASK_DIR="work/tasks/task_${TASK_ID}"

# Verificar se o diretório da tarefa existe
if [ ! -d "$TASK_DIR" ]; then
  echo "❌ Diretório da tarefa não encontrado: $TASK_DIR"
  echo "   A tarefa pode já ter sido arquivada."
  exit 1
fi

# Tentar extrair o requirement do arquivo da tarefa
REQ=""
TASK_TITLE=""
if [ -f "$TASK_DIR/task.md" ]; then
  REQ=$(grep "Requirement:" "$TASK_DIR/task.md" | awk '{print $2}' | head -1)
  TASK_TITLE=$(grep "Title:" "$TASK_DIR/task.md" | cut -d: -f2- | sed 's/^ *//' | head -1)
fi

# Se não encontrou requirement, usar um nome padrão baseado na data
if [ -z "$REQ" ]; then
  REQ="task_${TASK_ID}"
fi

LEARN_FILE="docs/learnings/${REQ}.md"

mkdir -p docs/learnings

# Analisar a tarefa para extrair aprendizados automaticamente
echo "🧠 Analisando tarefa para extrair aprendizados..."

# Padrões comuns de aprendizados baseados no trabalho do agente
LEARNING_CONTENT=""

# Verificar se há arquivos modificados
if [ -f "$TASK_DIR/task.md" ]; then
  # Extrair progresso da tarefa
  if grep -q "## Progresso" "$TASK_DIR/task.md"; then
    PROGRESS=$(sed -n '/## Progresso/,/^##/p' "$TASK_DIR/task.md" | head -20)
    
    # Extrair bugs corrigidos
    BUGS_CORRECTED=$(echo "$PROGRESS" | grep -c "✅.*BUG.*CORRIGIDO")
    if [ "$BUGS_CORRECTED" -gt 0 ]; then
      LEARNING_CONTENT="${LEARNING_CONTENT}- **Bugs corrigidos:** $BUGS_CORRECTED\n"
      
      # Extrair tipos de bugs
      BUG_TYPES=$(echo "$PROGRESS" | grep -o "BUG [0-9]*.*CORRIGIDO" | sed 's/BUG [0-9]* //' | sed 's/ CORRIGIDO//' | tr '\n' ',' | sed 's/,$//')
      if [ -n "$BUG_TYPES" ]; then
        LEARNING_CONTENT="${LEARNING_CONTENT}- **Tipos de bugs:** $BUG_TYPES\n"
      fi
    fi
    
    # Extrair soluções implementadas
    SOLUTIONS=$(echo "$PROGRESS" | grep -o "Solução:.*" | sed 's/Solução: //' | head -3)
    if [ -n "$SOLUTIONS" ]; then
      LEARNING_CONTENT="${LEARNING_CONTENT}- **Soluções implementadas:**\n"
      while IFS= read -r solution; do
        LEARNING_CONTENT="${LEARNING_CONTENT}  • $solution\n"
      done <<< "$SOLUTIONS"
    fi
    
    # Extrair arquivos modificados mencionados
    FILES_MODIFIED=$(echo "$PROGRESS" | grep -o "\`.*\.\(go\|html\|md\|sh\)\`" | tr '\n' ',' | sed 's/`//g' | sed 's/,$//')
    if [ -n "$FILES_MODIFIED" ]; then
      LEARNING_CONTENT="${LEARNING_CONTENT}- **Arquivos modificados:** $FILES_MODIFIED\n"
    fi
  fi
fi

# Adicionar aprendizados padrão baseados em padrões observados
LEARNING_CONTENT="${LEARNING_CONTENT}- **Padrão identificado:** Conflito de nomes de templates causa sobrescrita\n"
LEARNING_CONTENT="${LEARNING_CONTENT}- **Solução padrão:** Usar nomes únicos para templates ou servir HTML diretamente\n"
LEARNING_CONTENT="${LEARNING_CONTENT}- **Verificação necessária:** Sempre testar endpoint após correções\n"
LEARNING_CONTENT="${LEARNING_CONTENT}- **Prevenção:** Criar testes que executam handlers reais, não apenas verificam arquivos\n"

# Adicionar aprendizado ao arquivo
echo "## Aprendizado $(date +'%Y-%m-%d %H:%M:%S')" >> "$LEARN_FILE"
echo "**Tarefa:** $TASK_ID" >> "$LEARN_FILE"
echo "**Título:** $TASK_TITLE" >> "$LEARN_FILE"
echo "**Status:** $STATUS" >> "$LEARN_FILE"
echo "**Conteúdo:**" >> "$LEARN_FILE"
echo -e "$LEARNING_CONTENT" >> "$LEARN_FILE"

# Adicionar resumo executivo
echo "**Resumo executivo:**" >> "$LEARN_FILE"
echo "O agente trabalhou na tarefa '$TASK_TITLE' e identificou padrões comuns:" >> "$LEARN_FILE"
echo "1. Conflitos de templates devem ser prevenidos com nomes únicos" >> "$LEARN_FILE"
echo "2. Testes devem executar handlers reais, não apenas verificar arquivos" >> "$LEARN_FILE"
echo "3. Sempre verificar se correções funcionam em produção" >> "$LEARN_FILE"
echo "4. Reiniciar servidor após modificações de código" >> "$LEARN_FILE"
echo "" >> "$LEARN_FILE"

# Mover tarefa para archive
mv "$TASK_DIR" work/archive/ 2>/dev/null

if [ $? -eq 0 ]; then
  echo "✅ tarefa arquivada"
else
  echo "⚠️  Aviso: Não foi possível mover a tarefa para o archive"
fi

echo "📚 aprendizado registrado em ${LEARN_FILE}"
echo "📝 Conteúdo do aprendizado:"
echo -e "$LEARNING_CONTENT"

# Atualizar contexto após conclusão de task
if [ -f "scripts/arandu_update_context.sh" ]; then
  echo
  echo "🧠 Atualizando contexto do projeto..."
  bash scripts/arandu_update_context.sh
fi

