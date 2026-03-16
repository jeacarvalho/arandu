#!/usr/bin/env bash
# scripts/arandu_guard.sh

echo "🛡️ Arandu Guard: Verificando integridade do sistema..."

ROUTES=("/dashboard" "/patients" "/patients/new")
FAILED=0

for route in "${ROUTES[@]}"; do
  STATUS=$(curl -o /dev/null -s -w "%{http_code}" http://localhost:8080${route})
  if [ "$STATUS" -eq 200 ]; then
    echo "✅ Rota ${route} está online."
  else
    echo "❌ FALHA CRÍTICA: Rota ${route} retornou status ${STATUS}"
    FAILED=1
  fi
done

if [ $FAILED -eq 1 ]; then
  echo "🚨 O sistema apresenta regressões. Corrija antes de concluir a task."
  exit 1
fi


# --- Verificação de Geração de Templates (templ) ---
echo "🔍 Verificando integridade dos componentes templ..."

TEMPL_SOURCES=$(find web/ -name "*.templ")
STALE_FILES=0

for src in $TEMPL_SOURCES; do
  # Define o nome esperado do arquivo gerado (ex: login_templ.go)
  gen="${src%.templ}_templ.go"

  if [ ! -f "$gen" ]; then
    echo "❌ ERRO: Arquivo gerado não encontrado para $src"
    STALE_FILES=1
  elif [ "$src" -nt "$gen" ]; then
    echo "❌ ERRO: O arquivo $src foi modificado mas não foi regenerado."
    echo "   (O fonte é mais novo que o código gerado em $gen)"
    STALE_FILES=1
  fi
done

if [ $STALE_FILES -eq 1 ]; then
  echo "🚨 Falha na integridade: Execute 'templ generate' antes de prosseguir."
  exit 1
else
  echo "✅ Todos os componentes templ estão atualizados."
fi