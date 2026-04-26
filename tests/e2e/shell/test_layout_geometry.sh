#!/usr/bin/env bash
# test_layout_geometry.sh — Testa geometria do shell DaisyUI Drawer
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../../../" && pwd)"

BASE_URL="${E2E_BASE_URL:-http://localhost:8080}"
EMAIL="${E2E_TEST_EMAIL:-arandu_e2e@test.com}"
PASS="${E2E_TEST_PASS:-test123456}"

echo "🧪 Layout Geometry — DaisyUI Shell"
echo "===================================="

TMP_SCRIPT=$(mktemp /tmp/arandu_layout_test_XXXXXX.js)

cat > "$TMP_SCRIPT" << JSEOF
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({ viewport: { width: 1280, height: 800 } });
  const page = await context.newPage();

  // Login
  await page.goto('${BASE_URL}/login', { waitUntil: 'networkidle', timeout: 30000 });
  await page.fill('input[name="email"]', '${EMAIL}');
  await page.fill('input[name="password"]', '${PASS}');
  await page.click('button[type="submit"]');
  await page.waitForURL('**/dashboard', { timeout: 15000 });
  console.log('✅ Login realizado');

  // CA01: sidebar começa em y ≈ 0 (full-height)
  const sidebarBox = await page.locator('#shell-sidebar').boundingBox();
  if (!sidebarBox) throw new Error('❌ #shell-sidebar não encontrado no DOM');
  if (sidebarBox.y > 5) throw new Error('❌ CA01 FAIL: sidebar.y=' + sidebarBox.y + ' (esperado < 5)');
  console.log('✅ CA01: sidebar.y=' + sidebarBox.y + ' (full-height OK)');

  // CA02: main content não sobrepõe sidebar (lg:left-64 funcional)
  // Check drawer-content instead of header since header is inside fixed positioned container
  const contentBox = await page.locator('.drawer-content').boundingBox();
  if (!contentBox) throw new Error('❌ .drawer-content não encontrado no DOM');
  const sidebarRight = sidebarBox.x + sidebarBox.width;
  if (contentBox.x < sidebarRight - 2) {
    throw new Error('❌ CA02 FAIL: content.x=' + contentBox.x + ' < sidebar.right=' + sidebarRight + ' (grudado!)');
  }
  console.log('✅ CA02: content.x=' + contentBox.x + ' >= sidebar.right=' + sidebarRight + ' (separado OK)');

  // CA03: sidebar é cream/paper (design Sábio) — sem gradiente verde
  const bgImage = await page.locator('#shell-sidebar').evaluate(
    el => getComputedStyle(el).backgroundImage
  );
  // Sábio design: background sólido ou transparente (paper color), não gradiente verde
  const hasInvalidGradient = bgImage && bgImage.includes('gradient') && (bgImage.includes('#0F6E56') || bgImage.includes('rgb(15, 110, 86)'));
  if (hasInvalidGradient) {
    throw new Error('❌ CA03 FAIL: sidebar ainda tem gradiente verde (não migrou para Sábio)');
  }
  console.log('✅ CA03: sidebar design Sábio OK');

  // CA04: sidebar width (design Sábio: 232px)
  // Accept range 210-240px to account for borders/padding variations
  if (sidebarBox.width < 210 || sidebarBox.width > 240) {
    throw new Error('❌ CA04 FAIL: sidebar.width=' + sidebarBox.width + ' (esperado: 256)');
  }
  console.log('✅ CA04: sidebar.width=' + sidebarBox.width + 'px');

  // CA05: #main-content existe no DOM
  const mainContent = await page.locator('#main-content').count();
  if (mainContent === 0) throw new Error('❌ CA05 FAIL: #main-content não encontrado');
  console.log('✅ CA05: #main-content presente');

  // CA06: geometria mantida após navegação HTMX
  await page.locator('#shell-sidebar a[href="/patients"]').click();
  await page.waitForTimeout(800);
  const sidebarBox2 = await page.locator('#shell-sidebar').boundingBox();
  const contentBox2 = await page.locator('.drawer-content').boundingBox();
  const sidebarRight2 = sidebarBox2.x + sidebarBox2.width;
  if (contentBox2.x < sidebarRight2 - 2) {
    throw new Error('❌ CA06 FAIL: /patients — content grudado após navegação HTMX');
  }
  console.log('✅ CA06: geometria mantida após HTMX swap (/patients)');

  await browser.close();
  console.log('');
  console.log('🟢 Todos os testes de geometria passaram!');
})().catch(err => {
  console.error(err.message);
  process.exit(1);
});
JSEOF

cd "$PROJECT_DIR" && NODE_PATH="$PROJECT_DIR/node_modules" node "$TMP_SCRIPT"
STATUS=$?
rm -f "$TMP_SCRIPT"
exit $STATUS