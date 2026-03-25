#!/usr/bin/env bash
set -e

OUTPUT_DIR="${1:-./screenshots}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DASHBOARD_URL="${DASHBOARD_URL:-http://localhost:8080}"

mkdir -p "$OUTPUT_DIR"

echo "📸 Capturando screenshots do Arandu..."
echo "   URL: $DASHBOARD_URL"
echo "   Output: $OUTPUT_DIR"

playwright screenshot \
    --browser chromium \
    --viewport-size 1440,900 \
    "$DASHBOARD_URL/login" \
    "$OUTPUT_DIR/login_${TIMESTAMP}.png" 2>&1 || true

echo ""
echo "✅ Screenshots salvos em: $OUTPUT_DIR"
ls -la "$OUTPUT_DIR"/*.png 2>/dev/null | tail -5
