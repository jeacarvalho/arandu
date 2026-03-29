#!/usr/bin/env bash
# scripts/arandu_visual_check.sh
set -e
echo "👁️  Validação Visual e de Layout"
echo "================================="

ERRORS=0

# 1. Detectar position: fixed/absolute sem z-index
echo "🔍 Verificando elementos posicionados sem z-index..."
FIXED_WITHOUT_Z=$(grep -r "position:\s*fixed" web/static/css/ --include="*.css" -l | while read f; do
    if ! grep -q "z-index:" "$f"; then echo "$f"; fi
done | wc -l)

if [ "$FIXED_WITHOUT_Z" -gt 0 ]; then
    echo "❌ $FIXED_WITHOUT_Z arquivos CSS com position:fixed sem z-index"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ OK: Elementos fixed com z-index definido"
fi

# 2. Detectar margin/padding negativos perigosos
echo "🔍 Verificando margin/padding negativos..."
NEG_MARGIN=$(grep -r "margin-top:\s*-" web/static/css/ --include="*.css" | wc -l)
if [ "$NEG_MARGIN" -gt 5 ]; then
    echo "⚠️  $NEG_MARGIN ocorrências de margin negativo (pode causar overlap)"
fi

# 3. Verificar se Top Bar tem padding-top no main-content
echo "🔍 Verificando espaçamento para header fixo..."
if ! grep -q "padding-top:\s*[678][0-9]px" web/static/css/*.css; then
    echo "❌ AVISO: Nenhum padding-top grande encontrado (conteúdo pode ficar sob header)"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ OK: Padding-top para header fixo detectado"
fi

# 4. Validar breakpoints de responsive
echo "🔍 Verificando media queries..."
MOBILE_MQ=$(grep -r "@media.*max-width.*768" web/static/css/ --include="*.css" | wc -l)
if [ "$MOBILE_MQ" -lt 3 ]; then
    echo "⚠️  Apenas $MOBILE_MQ media queries para mobile (pode estar incompleto)"
fi

echo ""
if [ "$ERRORS" -gt 0 ]; then
    echo "🔴 VALIDAÇÃO FALHOU: $ERRORS erro(s) de layout detectado(s)"
    exit 1
else
    echo "🟢 VALIDAÇÃO APROVADA"
    exit 0
fi