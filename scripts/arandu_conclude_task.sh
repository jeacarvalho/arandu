#!/usr/bin/env bash

# Verificar argumentos
if [ $# -lt 1 ]; then
  echo "Uso: $0 TASK_ID [--success|--failure]"
  echo "Exemplo: $0 20260313_185952 --success"
  echo ""
  echo "O script analisa a tarefa concluída e sugere aprendizados valiosos"
  echo "para adicionar ao sistema consolidado de documentação."
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

echo "🧠 Concluindo tarefa: $TASK_TITLE ($REQ)"
echo "========================================"
echo

# Analisar a tarefa para identificar aprendizados potenciais
echo "📋 Analisando tarefa para identificar aprendizados valiosos..."
echo

LEARNING_SUGGESTIONS=""
HAS_VALUABLE_CONTENT=false

# Verificar se há arquivos modificados
if [ -f "$TASK_DIR/task.md" ]; then
  # Extrair progresso da tarefa
  if grep -q "## Progresso" "$TASK_DIR/task.md"; then
    PROGRESS=$(sed -n '/## Progresso/,/^##/p' "$TASK_DIR/task.md" | head -30)
    
    # Extrair bugs corrigidos (aprendizado valioso)
    BUGS_CORRECTED=$(echo "$PROGRESS" | grep -c "✅.*BUG.*CORRIGIDO")
    if [ "$BUGS_CORRECTED" -gt 0 ]; then
      HAS_VALUABLE_CONTENT=true
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}🐛 **Bugs corrigidos:** $BUGS_CORRECTED\n"
      
      # Extrair tipos de bugs específicos
      BUG_TYPES=$(echo "$PROGRESS" | grep -o "BUG [0-9]*.*CORRIGIDO" | sed 's/BUG [0-9]* //' | sed 's/ CORRIGIDO//' | head -3)
      if [ -n "$BUG_TYPES" ]; then
        while IFS= read -r bug_type; do
          LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  • $bug_type\n"
        done <<< "$BUG_TYPES"
      fi
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}\n"
    fi
    
    # Extrair soluções implementadas (aprendizado valioso)
    SOLUTIONS=$(echo "$PROGRESS" | grep -o "Solução:.*" | sed 's/Solução: //' | head -5)
    if [ -n "$SOLUTIONS" ] && [ $(echo "$SOLUTIONS" | wc -l) -gt 0 ]; then
      HAS_VALUABLE_CONTENT=true
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}🔧 **Soluções implementadas:**\n"
      while IFS= read -r solution; do
        LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  • $solution\n"
      done <<< "$SOLUTIONS"
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}\n"
    fi
    
    # Extrair padrões arquiteturais (aprendizado muito valioso)
    if echo "$PROGRESS" | grep -qi "arquitetur\|handler\|template\|viewmodel\|htmx\|templ"; then
      HAS_VALUABLE_CONTENT=true
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}🏗️ **Mudanças arquiteturais identificadas**\n"
      ARCH_PATTERNS=$(echo "$PROGRESS" | grep -i "arquitetur\|handler\|template\|viewmodel\|htmx\|templ" | head -3)
      while IFS= read -r pattern; do
        LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  • $pattern\n"
      done <<< "$ARCH_PATTERNS"
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}\n"
    fi
    
    # Extrair mudanças em banco de dados (aprendizado valioso)
    if echo "$PROGRESS" | grep -qi "sql\|migration\|fts5\|database\|query"; then
      HAS_VALUABLE_CONTENT=true
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}💾 **Mudanças em banco de dados**\n"
      DB_CHANGES=$(echo "$PROGRESS" | grep -i "sql\|migration\|fts5\|database\|query" | head -3)
      while IFS= read -r change; do
        LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  • $change\n"
      done <<< "$DB_CHANGES"
      LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}\n"
    fi
  fi
fi

# Verificar mudanças arquiteturais reais
echo "🔍 Verificando mudanças arquiteturais..."

# Verificar mudanças em components (Templ)
if [ -d "web/components" ]; then
  NEW_COMPONENTS=$(find web/components -name "*.templ" -newer "$TASK_DIR/task.md" 2>/dev/null | wc -l)
  if [ "$NEW_COMPONENTS" -gt 0 ]; then
    HAS_VALUABLE_CONTENT=true
    echo "  ⚠️  Mudança detectada: $NEW_COMPONENTS componentes Templ criados/modificados"
    LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}🎨 **Novos componentes Templ:** $NEW_COMPONENTS\n"
    LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  Consulte TEMPL_GUIDE.md para padrões estabelecidos\n\n"
  fi
fi

# Verificar novos handlers
NEW_HANDLERS=$(find internal/web/handlers -name "*_handler.go" -newer "$TASK_DIR/task.md" 2>/dev/null | wc -l)
if [ "$NEW_HANDLERS" -gt 0 ]; then
  HAS_VALUABLE_CONTENT=true
  echo "  📝 $NEW_HANDLERS novos handlers criados"
  LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}🔄 **Novos handlers:** $NEW_HANDLERS\n"
  LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  Consulte ARCHITECTURE_PATTERNS.md para padrões\n\n"
fi

# Verificar mudanças em domain
if [ -d "internal/domain" ]; then
  DOMAIN_CHANGES=$(find internal/domain -name "*.go" -newer "$TASK_DIR/task.md" 2>/dev/null | wc -l)
  if [ "$DOMAIN_CHANGES" -gt 0 ]; then
    HAS_VALUABLE_CONTENT=true
    echo "  📝 $DOMAIN_CHANGES alterações no domínio"
    LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}🩺 **Mudanças no domínio clínico:** $DOMAIN_CHANGES\n\n"
  fi
fi

# Verificar migrations
if [ -d "internal/infrastructure/repository/sqlite/migrations" ]; then
  NEW_MIGRATIONS=$(find internal/infrastructure/repository/sqlite/migrations -name "*.sql" -newer "$TASK_DIR/task.md" 2>/dev/null | wc -l)
  if [ "$NEW_MIGRATIONS" -gt 0 ]; then
    HAS_VALUABLE_CONTENT=true
    echo "  💾 $NEW_MIGRATIONS novas migrations SQL"
    LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}📊 **Novas migrations:** $NEW_MIGRATIONS\n"
    LEARNING_SUGGESTIONS="${LEARNING_SUGGESTIONS}  Consulte SQLITE_BEST_PRACTICES.md para padrões\n\n"
  fi
fi

echo
echo "📚 SISTEMA DE APRENDIZADOS CONSOLIDADO"
echo "======================================"
echo "O projeto agora usa um sistema consolidado de aprendizados:"
echo
echo "📁 docs/learnings/"
echo "├── MASTER_LEARNINGS.md          # 📚 Arquivo principal consolidado"
echo "├── ARCHITECTURE_PATTERNS.md     # 🏗️ Padrões arquiteturais"
echo "├── TEMPL_GUIDE.md              # 🎨 Guia Templ"
echo "└── SQLITE_BEST_PRACTICES.md    # 💾 Práticas SQLite"
echo

if [ "$HAS_VALUABLE_CONTENT" = true ]; then
  echo "✅ APRENDIZADOS VALIOSOS IDENTIFICADOS"
  echo "======================================"
  echo -e "$LEARNING_SUGGESTIONS"
  echo "💡 SUGESTÃO: Adicione estes aprendizados ao MASTER_LEARNINGS.md"
  echo
  echo "📝 Formato sugerido:"
  echo "### Título do Aprendizado"
  echo ""
  echo "**Contexto:** [Tarefa $TASK_ID - $TASK_TITLE]"
  echo ""
  echo "**Problema:** [Descreva o problema encontrado]"
  echo ""
  echo "**Solução:** [Descreva a solução implementada]"
  echo ""
  echo "**Referência:** $REQ"
  echo
  echo "📍 Localização: docs/learnings/MASTER_LEARNINGS.md"
else
  echo "ℹ️  NENHUM APRENDIZADO VALIOSO IDENTIFICADO"
  echo "=========================================="
  echo "Esta tarefa não gerou aprendizados significativos para o sistema consolidado."
  echo "Apenas aprendizados valiosos e não-repetitivos devem ser adicionados."
  echo
  echo "📊 Conteúdo analisado:"
  echo "- Bugs corrigidos: $BUGS_CORRECTED"
  echo "- Mudanças arquiteturais: $([ "$HAS_VALUABLE_CONTENT" = true ] && echo "Sim" || echo "Não")"
  echo "- Valor para documentação: Baixo"
fi

echo
echo "🛡️ Executando verificações de segurança antes de arquivar..."
if [ -f "scripts/arandu_guard.sh" ]; then
  bash scripts/arandu_guard.sh || exit 1
fi

# Verificar se há mudança em .templ sem o correspondente _templ.go
if find web/ -name "*.templ" -newer "$TASK_DIR/task.md" 2>/dev/null | grep -q "."; then
    echo "🔍 Verificando se os componentes Templ foram gerados..."
    # Lógica para garantir que o agente não esqueceu de rodar templ generate
    STALE_TEMPL=$(find web/ -name "*.templ" -newer "$TASK_DIR/task.md" 2>/dev/null | head -1)
    if [ -n "$STALE_TEMPL" ]; then
        GEN_FILE="${STALE_TEMPL%.templ}_templ.go"
        if [ ! -f "$GEN_FILE" ] || [ "$STALE_TEMPL" -nt "$GEN_FILE" ]; then
            echo "❌ ERRO: Arquivo $STALE_TEMPL modificado mas não foi regenerado."
            echo "   Execute: templ generate"
            exit 1
        fi
    fi
fi

# Mover tarefa para archive
mv "$TASK_DIR" work/archive/ 2>/dev/null

if [ $? -eq 0 ]; then
  echo "✅ Tarefa arquivada em work/archive/"
else
  echo "⚠️  Aviso: Não foi possível mover a tarefa para o archive"
fi

echo
echo "🎯 CONCLUSÃO DA TAREFA"
echo "======================"
echo "Tarefa: $TASK_TITLE"
echo "Status: $STATUS"
echo "Requirement: $REQ"
echo
if [ "$HAS_VALUABLE_CONTENT" = true ]; then
  echo "💎 Esta tarefa gerou aprendizados valiosos!"
  echo "   Considere adicioná-los ao MASTER_LEARNINGS.md"
else
  echo "📝 Esta tarefa não requer atualização da documentação de aprendizados."
fi

# Atualizar contexto após conclusão de task
if [ -f "scripts/arandu_update_context.sh" ]; then
  echo
  echo "🧠 Atualizando contexto do projeto..."
  bash scripts/arandu_update_context.sh
fi

echo
echo "✨ Tarefa concluída com sucesso!"