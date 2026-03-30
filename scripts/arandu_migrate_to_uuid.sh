#!/bin/bash
# ⚠️  ONE-TIME MIGRATION SCRIPT - Run once only, then archive or delete
# This script was used to migrate IDs to UUID v4 format
set -e

echo "🔄 Script de Migração para UUID v4"
echo "=================================="

CENTRAL_DB="storage/arandu_central.db"
STORAGE_DIR="storage/tenants"

if [ ! -f "$CENTRAL_DB" ]; then
    echo "❌ Banco central não encontrado: $CENTRAL_DB"
    exit 1
fi

echo "📊 Estado atual:"
echo "Tenants:"
sqlite3 "$CENTRAL_DB" "SELECT id FROM tenants;" 2>/dev/null || echo "  Nenhum"
echo "Users:"
sqlite3 "$CENTRAL_DB" "SELECT id, email FROM users;" 2>/dev/null || echo "  Nenhum"

echo ""
echo "⚠️  ATENÇÃO: Este script vai:"
echo "   1. Gerar novos UUIDs para tenants, users e sessions"
echo "   2. Renomear os arquivos de banco de dados"
echo "   3. Atualizar o banco central"
echo ""
read -p "Continuar com a migração? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Migração cancelada."
    exit 0
fi

echo ""
echo "🚀 Iniciando migração..."

for tenant_row in $(sqlite3 "$CENTRAL_DB" "SELECT id || '|' || db_path FROM tenants WHERE id LIKE 'tenant-%' OR id LIKE 'default-%';"); do
    old_tenant_id=$(echo "$tenant_row" | cut -d'|' -f1)
    old_db_path=$(echo "$tenant_row" | cut -d'|' -f2)
    new_tenant_id=$(uuidgen | tr '[:upper:]' '[:lower:]')
    
    echo ""
    echo "📦 Migrando tenant: $old_tenant_id -> $new_tenant_id"
    
    if [ -f "$old_db_path" ]; then
        new_db_path="storage/tenants/clinical_${new_tenant_id}.db"
        mv "$old_db_path" "$new_db_path"
        echo "  ✓ Arquivo renomeado: $(basename $old_db_path) -> $(basename $new_db_path)"
    fi
    
    sqlite3 "$CENTRAL_DB" "UPDATE tenants SET id = '$new_tenant_id', db_path = '$new_db_path' WHERE id = '$old_tenant_id';"
    
    for user_row in $(sqlite3 "$CENTRAL_DB" "SELECT id || '|' || email FROM users WHERE tenant_id = '$old_tenant_id' AND (id LIKE 'user-%' OR id LIKE 'test-%');"); do
        old_user_id=$(echo "$user_row" | cut -d'|' -f1)
        new_user_id=$(uuidgen | tr '[:upper:]' '[:lower:]')
        
        echo "  📧 Migrando usuário: $old_user_id -> $new_user_id"
        sqlite3 "$CENTRAL_DB" "UPDATE users SET id = '$new_user_id', tenant_id = '$new_tenant_id' WHERE id = '$old_user_id';"
        sqlite3 "$CENTRAL_DB" "UPDATE sessions SET user_id = '$new_user_id' WHERE user_id = '$old_user_id';"
    done
    
    for session_id in $(sqlite3 "$CENTRAL_DB" "SELECT id FROM sessions WHERE tenant_id = '$old_tenant_id' AND id LIKE 'google-session-%';"); do
        new_session_id=$(uuidgen | tr '[:upper:]' '[:lower:]')
        echo "  🔑 Migrando sessão: $session_id -> $new_session_id"
        sqlite3 "$CENTRAL_DB" "UPDATE sessions SET id = '$new_session_id' WHERE id = '$session_id';"
    done
    
    echo "  ✅ Tenant migrado com sucesso"
done

echo ""
echo "📊 Resultado da migração:"
echo "Tenants:"
sqlite3 "$CENTRAL_DB" "SELECT id, db_path FROM tenants;"
echo ""
echo "Users:"
sqlite3 "$CENTRAL_DB" "SELECT id, email, tenant_id FROM users;"

echo ""
echo "✅ Migração concluída!"
echo ""
echo "📝 Nota: Os usuários existentes precisam fazer logout e login novamente para obter novas sessões UUID."
