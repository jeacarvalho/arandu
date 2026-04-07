#!/usr/bin/env bash
# scripts/arandu_guard.sh

echo "🛡️ Arandu Guard: Verificando integridade do sistema..."

# Create test session directly in central.db
COOKIE_FILE=$(mktemp)
CENTRAL_DB="storage/arandu_central.db"

if [ -f "$CENTRAL_DB" ]; then
	# Usar um tenant existente válido
	TENANT_ID=$(sqlite3 "$CENTRAL_DB" "SELECT id FROM tenants WHERE status = 'active' LIMIT 1;")
	
	if [ -z "$TENANT_ID" ]; then
		echo "❌ Nenhum tenant ativo encontrado"
		exit 1
	fi
	
	echo "📋 Usando tenant existente: $TENANT_ID"
	
	# Create/update test user
	sqlite3 "$CENTRAL_DB" "INSERT OR REPLACE INTO users (id, email, password_hash, tenant_id, created_at) VALUES ('test-user-001', 'test@test.com', '', '$TENANT_ID', datetime('now'));"
	
	# Create session
	SESSION_ID="guard-test-$(date +%s)"
	EXPIRES=$(date -d "+7 days" +%s 2>/dev/null || date -v+7d +%s)
	sqlite3 "$CENTRAL_DB" "DELETE FROM sessions WHERE id LIKE 'guard-test%';"
	sqlite3 "$CENTRAL_DB" "INSERT INTO sessions (id, user_id, tenant_id, expires_at) VALUES ('$SESSION_ID', 'test-user-001', '$TENANT_ID', $EXPIRES);"
	
	# Create cookie file com formato correto (Netscape - tabs como delimitadores)
	echo "# Netscape HTTP Cookie File" > "$COOKIE_FILE"
	printf ".localhost\tTRUE\t/\tFALSE\t%s\tarandu_session\t%s\n" "$EXPIRES" "$SESSION_ID" >> "$COOKIE_FILE"
fi

ROUTES=("/dashboard" "/patients" "/patients/new")
FAILED=0

for route in "${ROUTES[@]}"; do
	STATUS=$(curl -o /dev/null -s -b "$COOKIE_FILE" -w "%{http_code}" http://localhost:8080${route})
	if [ "$STATUS" -eq 200 ]; then
		echo "✅ Rota ${route} está online."
	else
		echo "❌ FALHA CRÍTICA: Rota ${route} retornou status ${STATUS}"
		FAILED=1
	fi
done

rm -f "$COOKIE_FILE"

if [ $FAILED -eq 1 ]; then
	echo "🚨 O sistema apresenta regressões. Corrija antes de concluir a task."
	exit 1
fi

echo "🔍 Verificando integridade dos componentes templ..."

TEMPL_SOURCES=$(find web/ -name "*.templ")
STALE_FILES=0

for src in $TEMPL_SOURCES; do
	gen="${src%.templ}_templ.go"

	if [ ! -f "$gen" ]; then
		echo "❌ ERRO: Arquivo gerado não encontrado para $src"
		STALE_FILES=1
	elif [ "$src" -nt "$gen" ]; then
		echo "❌ ERRO: O arquivo $src foi modificado mas não foi regenerado."
		STALE_FILES=1
	fi
done

if [ $STALE_FILES -eq 1 ]; then
	echo "🚨 Falha na integridade: Execute 'templ generate' antes de prosseguir."
	exit 1
else
	echo "✅ Todos os componentes templ estão atualizados."
fi
