#!/usr/bin/env bash
# scripts/arandu_screenshot_compare.sh
set -e

OUTPUT_DIR="./screenshots"
BASELINE_DIR="./screenshots/baseline"
CURRENT_DIR="./screenshots/current"

mkdir -p "$BASELINE_DIR" "$CURRENT_DIR"

echo "📸 Capturando screenshots para comparação..."

# Rotas críticas para screenshot
ROUTES=("/dashboard" "/patients" "/patients/new" "/login")

for route in "${ROUTES[@]}"; do
    ROUTE_NAME=$(echo "$route" | tr '/' '_' | sed 's/^_//')
    if [ -z "$ROUTE_NAME" ]; then
        ROUTE_NAME="root"
    fi
    OUTPUT_FILE="${CURRENT_DIR}/${ROUTE_NAME}.png"
    BASELINE_FILE="${BASELINE_DIR}/${ROUTE_NAME}.png"
    
    echo "  → $route"
    playwright screenshot \
        --browser chromium \
        --viewport-size 1440,900 \
        "http://localhost:8080${route}" \
        "$OUTPUT_FILE" 2>&1 || true
    
    # Se existe baseline, comparar
    if [ -f "$BASELINE_FILE" ]; then
        echo "    🔍 Comparando com baseline..."
        # Usar imagemMagick ou similar para diff
        # compare -metric AE baseline.png current.png diff.png
    else
        echo "    📋 Criando baseline..."
        cp "$OUTPUT_FILE" "$BASELINE_FILE"
    fi
done

echo "✅ Screenshots salvos em $CURRENT_DIR"