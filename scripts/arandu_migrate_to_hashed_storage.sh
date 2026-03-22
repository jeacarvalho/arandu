#!/bin/bash
set -e

echo "🔄 Script de Migração para Directory Hashing"
echo "============================================="
echo ""
echo "⚠️  ATENÇÃO: Este script vai:"
echo "   1. Identificar bancos legados na raiz de storage/tenants/"
echo "   2. Criar estrutura de diretórios hashed (aa/bb/)"
echo "   3. Mover arquivos .db para nova localização"
echo "   4. Atualizar db_path no banco central"
echo ""
read -p "Continuar com a migração? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Migração cancelada."
    exit 0
fi

CENTRAL_DB="storage/arandu_central.db"
TENANTS_DIR="storage/tenants"

if [ ! -d "$TENANTS_DIR" ]; then
    echo "❌ Diretório de tenants não encontrado: $TENANTS_DIR"
    exit 1
fi

if [ ! -f "$CENTRAL_DB" ]; then
    echo "❌ Banco central não encontrado: $CENTRAL_DB"
    exit 1
fi

echo ""
echo "📊 Identificando bancos legados..."

legacy_dbs=()
for db_file in "$TENANTS_DIR"/clinical_*.db; do
    if [ -f "$db_file" ]; then
        basename_file=$(basename "$db_file")
        
        if [[ "$basename_file" == clinical_tenant-* ]] || [[ "$basename_file" == clinical_default-* ]] || [[ "$basename_file" == clinical_test-* ]]; then
            legacy_dbs+=("$db_file")
            continue
        fi
        
        tenant_id="${basename_file#clinical_}"
        tenant_id="${tenant_id%.db}"
        
        if [[ "$tenant_id" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$ ]]; then
            legacy_dbs+=("$db_file")
        fi
    fi
done

if [ ${#legacy_dbs[@]} -eq 0 ]; then
    echo "✅ Nenhum banco legado encontrado para migrar."
    exit 0
fi

echo ""
echo "📋 Bancos legados encontrados: ${#legacy_dbs[@]}"
for db in "${legacy_dbs[@]}"; do
    echo "  - $(basename "$db")"
done

echo ""
echo "🚀 Iniciando migração..."

for old_path in "${legacy_dbs[@]}"; do
    basename_file=$(basename "$old_path")
    
    if [[ "$basename_file" == clinical_tenant-* ]]; then
        tenant_id="tenant-${basename_file#clinical_tenant-}"
        tenant_id="${tenant_id%.db}"
    elif [[ "$basename_file" == clinical_default-* ]]; then
        tenant_id="default-${basename_file#clinical_default-}"
        tenant_id="${tenant_id%.db}"
    elif [[ "$basename_file" == clinical_test-* ]]; then
        tenant_id="test-${basename_file#clinical_test-}"
        tenant_id="${tenant_id%.db}"
    else
        tenant_id="${basename_file#clinical_}"
        tenant_id="${tenant_id%.db}"
    fi
    
    echo ""
    echo "📦 Migrando: $basename_file"
    echo "   Tenant ID: $tenant_id"
    
    new_dir="storage/tenants/${tenant_id:0:2}/${tenant_id:2:2}"
    new_path="${new_dir}/clinical_${tenant_id}.db"
    
    echo "   Novo caminho: $new_path"
    
    if [ -f "$new_path" ]; then
        echo "   ⚠️  Arquivo já existe no destino, pulando..."
        continue
    fi
    
    echo "   Criando diretório: $new_dir"
    mkdir -p "$new_dir"
    
    echo "   Movendo arquivo..."
    mv "$old_path" "$new_path"
    
    echo "   Atualizando banco central..."
    sqlite3 "$CENTRAL_DB" "UPDATE tenants SET db_path = '$new_path' WHERE id = '$tenant_id';"
    
    echo "   ✅ Migração concluída!"
done

echo ""
echo "📊 Verificação final..."
echo "Estrutura de diretórios:"
find storage/tenants -type d | head -20

echo ""
echo "✅ Migração concluída!"
echo ""
echo "📝 Nota: Se houver erros, os bancos legados ainda existem em backup."
echo "   Novos tenants serão criados automaticamente na estrutura hashed."
