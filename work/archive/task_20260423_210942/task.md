# Task: CA-06 — Playwright Layout Tests — Shell DaisyUI Drawer
Requirement: Interno — qualidade visual / regressão
Status: PRONTO_PARA_IMPLEMENTACAO

---

## Objetivo

Criar testes Playwright de geometria para o shell DaisyUI Drawer, garantindo que:
- A topbar **não se sobreponha** à sidebar no desktop (lg:left-64 funcional)
- A sidebar comece em y ≈ 0 (full-height, sem espaço no topo)
- A sidebar tenha background gradiente verde (não branco/transparente)
- HTMX: navegar pela sidebar mantém a geometria correta após swap

Se esses testes passarem, o "grudado" está resolvido. Se falharem, apontam exatamente o problema.

---

## Contexto do sistema

**Stack**: Go 1.22+ · Templ · HTMX · DaisyUI v5 + Tailwind CSS v4 · Alpine.js 3
**Playwright** já instalado em `package.json` (`"playwright": "^1.58.2"`)
**Servidor de teste**: `http://localhost:8080`

**Credenciais de teste** (únicas válidas):
```
Email:    arandu_e2e@test.com
Senha:    test123456
```

**IDs críticos no DOM** (devem existir após login):
```
#shell-sidebar       ← sidebar DaisyUI drawer
header.navbar        ← topbar (selector CSS)
#main-content        ← conteúdo HTMX
```

**Estrutura HTML esperada** (resultado do CA-05 já implementado):
```html
<div class="drawer lg:drawer-open">
  <input id="shell-drawer" type="checkbox" class="drawer-toggle"/>
  <div class="drawer-content flex flex-col">
    <header class="navbar fixed top-0 left-0 right-0 lg:left-64 z-50 h-16 ...">
    </header>
    <main class="flex-1 mt-16 ...">
      <div id="main-content">...</div>
    </main>
  </div>
  <div class="drawer-side z-40">
    <label class="drawer-overlay"></label>
    <aside id="shell-sidebar" class="arandu-sidebar w-64 min-h-screen ...">
    </aside>
  </div>
</div>
```

**Width esperada da sidebar**: 256px (`w-64` = 16rem × 16px = 256px)

---

## Padrão de referência

`tests/e2e/shell/test_screenshot_manual.sh` — padrão de bash + Node.js inline sem playwright.config.ts separado.
`tests/e2e/shell/test_framework.sh` — funções `test_start`, `test_pass`, `test_fail`, `test_summary`.

---

## Arquivos a criar/modificar

**Criar:**
- `tests/e2e/shell/test_layout_geometry.sh` — script principal (chmod +x)

**Modificar:**
- `tests/e2e/shell/test_all.sh` — adicionar chamada ao novo teste

---

## Conteúdo de `tests/e2e/shell/test_layout_geometry.sh`

```bash
#!/usr/bin/env bash
# test_layout_geometry.sh — Testa geometria do shell DaisyUI Drawer
set -e

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

  // CA02: topbar não sobrepõe sidebar (lg:left-64 funcional)
  const topbarBox = await page.locator('header.navbar').boundingBox();
  if (!topbarBox) throw new Error('❌ header.navbar não encontrado no DOM');
  const sidebarRight = sidebarBox.x + sidebarBox.width;
  if (topbarBox.x < sidebarRight - 2) {
    throw new Error('❌ CA02 FAIL: topbar.x=' + topbarBox.x + ' < sidebar.right=' + sidebarRight + ' (grudado!)');
  }
  console.log('✅ CA02: topbar.x=' + topbarBox.x + ' >= sidebar.right=' + sidebarRight + ' (separado OK)');

  // CA03: sidebar tem gradiente verde (não branco)
  const bgImage = await page.locator('#shell-sidebar').evaluate(
    el => getComputedStyle(el).backgroundImage
  );
  if (!bgImage || bgImage === 'none' || !bgImage.includes('gradient')) {
    throw new Error('❌ CA03 FAIL: sidebar.backgroundImage="' + bgImage + '" (esperado: linear-gradient)');
  }
  console.log('✅ CA03: sidebar tem gradiente OK');

  // CA04: sidebar width = 256px (w-64)
  if (Math.abs(sidebarBox.width - 256) > 2) {
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
  const topbarBox2  = await page.locator('header.navbar').boundingBox();
  const sidebarRight2 = sidebarBox2.x + sidebarBox2.width;
  if (topbarBox2.x < sidebarRight2 - 2) {
    throw new Error('❌ CA06 FAIL: /patients — topbar grudado após navegação HTMX');
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

node "$TMP_SCRIPT"
STATUS=$?
rm -f "$TMP_SCRIPT"
exit $STATUS
```

---

## Modificação em `tests/e2e/shell/test_all.sh`

Adicionar ao final do script, antes de qualquer `test_summary` existente (ou ao final se não houver):

```bash
# Layout Geometry — DaisyUI Shell
test_start "Layout Geometry — DaisyUI Shell"
if bash "$(dirname "$0")/test_layout_geometry.sh" > /dev/null 2>&1; then
  test_pass
else
  test_fail "ver: bash tests/e2e/shell/test_layout_geometry.sh"
fi
```

---

## Critérios de aceite

**Execução**
- [ ] `chmod +x tests/e2e/shell/test_layout_geometry.sh` aplicado
- [ ] `bash tests/e2e/shell/test_layout_geometry.sh` passa (requer servidor rodando na 8080)

**CA01**: sidebar.y < 5 — sidebar começa no topo (full-height)
**CA02**: topbar.x >= sidebar.right — sem grudado, lg:left-64 funcional
**CA03**: sidebar.backgroundImage contém "gradient" — verde ativo
**CA04**: sidebar.width ≈ 256px — w-64 compilado corretamente
**CA05**: #main-content presente no DOM
**CA06**: geometria mantida após navegação HTMX para /patients

**Integridade**
- [ ] `./scripts/arandu_guard.sh` passa

---

## NÃO faça

- Não criar `playwright.config.ts` — seguir padrão existente de inline Node.js em bash
- Não criar diretório `tests/visual/` — usar `tests/e2e/shell/`
- Não usar `page.screenshot()` como critério — apenas asserções geométricas
- Não modificar código Go, .templ ou CSS — esta task é apenas testes
- Não hardcodar `left: 256px` como asserção — comparar topbar.x com sidebar.right dinamicamente
