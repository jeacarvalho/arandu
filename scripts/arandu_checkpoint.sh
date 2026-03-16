#!/usr/bin/env bash

set -e

echo "🔐 CHECKPOINT ARANDU - Validação Arquitetural Obrigatória"
echo "=========================================================="
echo

# Data e hora
CHECKPOINT_TIME=$(date +"%Y-%m-%d %H:%M:%S")
echo "⏰ Checkpoint: $CHECKPOINT_TIME"
echo

# 1. Verificar estado do git
echo "📊 1. ESTADO DO GIT"
echo "------------------"
git status --short
echo

# 2. Validar handlers
echo "🔍 2. VALIDAÇÃO DE HANDLERS"
echo "---------------------------"
if [ -f "scripts/arandu_validate_handlers.sh" ]; then
    bash scripts/arandu_validate_handlers.sh
    HANDLER_VALIDATION=$?
else
    echo "❌ Script de validação não encontrado"
    HANDLER_VALIDATION=1
fi
echo

# 3. Verificar build
echo "🔨 3. VERIFICAÇÃO DE BUILD"
echo "--------------------------"
echo "🔎 Build dos handlers..."
if go build ./internal/web/handlers/...; then
    echo "✅ OK: Handlers compilam"
    BUILD_HANDLERS=0
else
    echo "❌ ERRO: Handlers não compilam"
    BUILD_HANDLERS=1
fi

echo "🔎 Build completo..."
if go build ./cmd/arandu; then
    echo "✅ OK: Build completo bem-sucedido"
    BUILD_FULL=0
else
    echo "❌ ERRO: Build completo falhou"
    BUILD_FULL=1
fi
echo

# 4. Verificar migrações
echo "🗄️  4. VERIFICAÇÃO DE MIGRAÇÕES"
echo "-----------------------------"
if [ -f "internal/infrastructure/repository/sqlite/migrations/0001_initial_schema.up.sql" ]; then
    echo "✅ OK: Migração inicial encontrada"
    MIGRATION_EXISTS=0
else
    echo "⚠️  AVISO: Migração inicial não encontrada"
    MIGRATION_EXISTS=1
fi

echo "🔎 Verificando InitSchema depreciado..."
# Verificar se InitSchema em repositórios retorna nil (depreciado)
# schema.go pode chamar Migrate() para backward compatibility
INITSCHEMA_PROBLEM=0
for file in internal/infrastructure/repository/sqlite/*.go; do
    if grep -q "func.*InitSchema" "$file"; then
        # schema.go pode chamar Migrate() - isso é OK
        if [[ "$file" == *"schema.go" ]]; then
            if ! grep -A5 "func.*InitSchema" "$file" | grep -q "db.Migrate()"; then
                echo "⚠️  AVISO: $file - InitSchema não chama Migrate() para backward compatibility"
                INITSCHEMA_PROBLEM=1
            fi
        else
            # Outros repositórios devem retornar nil (depreciado)
            if ! grep -A5 "func.*InitSchema" "$file" | grep -q "return nil" && \
               ! grep -A5 "func.*InitSchema" "$file" | grep -q "// Schema creation is now handled by migrations"; then
                echo "❌ ERRO: $file - InitSchema não está depreciado corretamente"
                INITSCHEMA_PROBLEM=2
            fi
        fi
    fi
done

if [ "$INITSCHEMA_PROBLEM" -eq 0 ]; then
    echo "✅ OK: InitSchema depreciado corretamente"
    INITSCHEMA_USAGE=0
elif [ "$INITSCHEMA_PROBLEM" -eq 1 ]; then
    echo "⚠️  AVISO: InitSchema com aviso de backward compatibility"
    INITSCHEMA_USAGE=0  # Apenas aviso, não erro
else
    echo "❌ ERRO: InitSchema não depreciado corretamente"
    INITSCHEMA_USAGE=1
fi
echo

# 5. Verificar anti-padrões
echo "🚫 5. VERIFICAÇÃO DE ANTI-PADRÕES"
echo "--------------------------------"
echo "🔎 Procurando HTML inline..."
HTML_INLINE_COUNT=$(grep -r "w\.Write.*<.*>" internal/web/handlers/ --include="*.go" | wc -l)
if [ "$HTML_INLINE_COUNT" -gt 0 ]; then
    echo "❌ ERRO: $HTML_INLINE_COUNT ocorrência(s) de HTML inline"
else
    echo "✅ OK: Nenhum HTML inline encontrado"
fi

echo "🔎 Procurando arquivos .html..."
HTML_FILES_COUNT=$(find . -name "*.html" -not -path "./web/static/*" -not -path "./node_modules/*" -not -path "./.git/*" | wc -l)
if [ "$HTML_FILES_COUNT" -gt 0 ]; then
    echo "⚠️  AVISO: $HTML_FILES_COUNT arquivo(s) .html fora de web/static/"
    find . -name "*.html" -not -path "./web/static/*" -not -path "./node_modules/*" -not -path "./.git/*"
else
    echo "✅ OK: Nenhum arquivo .html fora de web/static/"
fi
echo

# 6. Resumo
echo "📋 6. RESUMO DO CHECKPOINT"
echo "--------------------------"

ERROR_COUNT=0
WARNING_COUNT=0

# Contar erros
if [ "$HANDLER_VALIDATION" -ne 0 ]; then ERROR_COUNT=$((ERROR_COUNT + 1)); fi
if [ "$BUILD_HANDLERS" -ne 0 ]; then ERROR_COUNT=$((ERROR_COUNT + 1)); fi
if [ "$BUILD_FULL" -ne 0 ]; then ERROR_COUNT=$((ERROR_COUNT + 1)); fi
if [ "$INITSCHEMA_USAGE" -ne 0 ]; then ERROR_COUNT=$((ERROR_COUNT + 1)); fi
if [ "$HTML_INLINE_COUNT" -gt 0 ]; then ERROR_COUNT=$((ERROR_COUNT + HTML_INLINE_COUNT)); fi

# Contar avisos
if [ "$MIGRATION_EXISTS" -ne 0 ]; then WARNING_COUNT=$((WARNING_COUNT + 1)); fi
if [ "$HTML_FILES_COUNT" -gt 0 ]; then WARNING_COUNT=$((WARNING_COUNT + HTML_FILES_COUNT)); fi

echo "❌ Erros críticos: $ERROR_COUNT"
echo "⚠️  Avisos: $WARNING_COUNT"
echo

# 7. Decisão
if [ "$ERROR_COUNT" -gt 0 ]; then
    echo "🔴 CHECKPOINT FALHOU"
    echo "==================="
    echo "Não é possível prosseguir. Corrija os erros acima antes de continuar."
    echo
    echo "Ações recomendadas:"
    echo "1. Execute ./scripts/arandu_validate_handlers.sh para detalhes"
    echo "2. Corrija os erros de compilação"
    echo "3. Remova HTML inline e arquivos .html não autorizados"
    echo "4. Execute este checkpoint novamente"
    exit 1
elif [ "$WARNING_COUNT" -gt 0 ]; then
    echo "🟡 CHECKPOINT COM AVISOS"
    echo "======================="
    echo "Pode prosseguir, mas considere corrigir os avisos."
    echo
    echo "Avisos a considerar:"
    if [ "$MIGRATION_EXISTS" -ne 0 ]; then echo "- Migração inicial não encontrada"; fi
    if [ "$HTML_FILES_COUNT" -gt 0 ]; then echo "- Arquivos .html fora de web/static/"; fi
    echo
    echo "Deseja prosseguir mesmo com avisos? (s/n)"
    read -r RESPONSE
    if [[ "$RESPONSE" =~ ^[Ss]$ ]]; then
        echo "✅ Prosseguindo com avisos..."
        exit 0
    else
        echo "❌ Checkpoint cancelado pelo usuário"
        exit 1
    fi
else
    echo "🟢 CHECKPOINT APROVADO"
    echo "====================="
    echo "Todas as validações passaram. Pode prosseguir com segurança."
    echo
    echo "✅ Handlers validados"
    echo "✅ Build bem-sucedido"
    echo "✅ Migrações corretas"
    echo "✅ Sem anti-padrões críticos"
    exit 0
fi