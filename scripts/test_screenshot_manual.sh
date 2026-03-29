#!/bin/bash
# scripts/test_screenshot_manual.sh
set -e

echo "🧪 Teste manual de screenshot com Playwright..."

# Criar script Node.js temporário
TMP_SCRIPT=$(mktemp)
cat > "$TMP_SCRIPT" << 'EOF'
const { chromium } = require('playwright');

(async () => {
    console.log("🚀 Iniciando browser...");
    const browser = await chromium.launch({ headless: true });
    console.log("✅ Browser iniciado");
    
    const context = await browser.newContext({ viewport: { width: 1440, height: 900 } });
    const page = await context.newPage();
    
    console.log("🌐 Navegando para http://localhost:8080/login...");
    await page.goto('http://localhost:8080/login', { waitUntil: 'networkidle', timeout: 30000 });
    
    console.log("📸 Capturando screenshot...");
    await page.screenshot({ path: 'tmp/test_screenshot.png', fullPage: true });
    
    await browser.close();
    console.log("✅ Screenshot salvo em tmp/test_screenshot.png");
})();
EOF

# Executar
node "$TMP_SCRIPT"
rm -f "$TMP_SCRIPT"

# Verificar resultado
if [ -f "tmp/test_screenshot.png" ]; then
    echo "✅ Screenshot gerado com sucesso!"
    ls -lh tmp/test_screenshot.png
else
    echo "❌ Screenshot não foi gerado"
    exit 1
fi