#!/bin/bash
# scripts/analyze_layout_density.sh

echo "📐 Análise de Densidade de Layout - Arandu"
echo "==========================================="

AUDIT_DIR="${1:-tmp/audit_logs}"

for file in "$AUDIT_DIR"/*.html; do
    [ -f "$file" ] || continue
    
    filename=$(basename "$file")
    
    # Contar elementos
    total_divs=$(grep -c '<div' "$file" 2>/dev/null || echo "0")
    empty_divs=$(grep -cE '<div[^>]*>\s*</div>' "$file" 2>/dev/null || echo "0")
    grid_usage=$(grep -cE 'class=".*grid' "$file" 2>/dev/null || echo "0")
    flex_usage=$(grep -cE 'class=".*flex' "$file" 2>/dev/null || echo "0")
    
    # Calcular métricas
    if [ "$total_divs" -gt 0 ]; then
        waste_pct=$((empty_divs * 100 / total_divs))
    else
        waste_pct=0
    fi
    
    # Alertas
    status="✅"
    if [ "$waste_pct" -gt 30 ]; then
        status="⚠️ "
    fi
    
    echo ""
    echo "$status $filename"
    echo "   Total divs: $total_divs"
    echo "   Containers vazios: $empty_divs ($waste_pct%)"
    echo "   Grid: $grid_usage | Flex: $flex_usage"
    
    if [ "$waste_pct" -gt 30 ]; then
        echo "   ⚠️  WARNING: Alto desperdício de espaço detectado"
    fi
done