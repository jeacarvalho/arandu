const { chromium } = require('playwright');
(async () => {
  const browser = await chromium.launch();
  const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 } });
  await ctx.addCookies([{ name: 'arandu_session', value: '0289b071-f3e1-4efe-951a-bc49a60e32b2', domain: 'localhost', path: '/' }]);
  const page = await ctx.newPage();
  await page.goto('http://localhost:8080/dashboard', { waitUntil: 'networkidle' });
  await page.screenshot({ path: 'tmp/test_dashboard.png', fullPage: true });
  await browser.close();
  console.log('✅ Done');
})();
