import { test, expect, Page } from '@playwright/test';
import { login, assertShellLayout, getPatientIDFromURL } from './helpers/auth';
import { assertBreadcrumb } from './helpers/layout';

const BASE_URL = process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:8080';
const EMAIL = process.env.E2E_EMAIL || 'playwright_e2e@arandu.internal';
const PASSWORD = process.env.E2E_PASSWORD || 'playwright_test_2026';

test.describe.serial('Fluxo Clínico Completo', () => {
  let page: Page;
  let patientName: string;
  let patientID: string;
  let sessionID: string;

  const OBS_TEXT = 'Paciente relata melhora no humor. Sono regularizado nos últimos dias.';
  const INT_TEXT = 'Técnica de reestruturação cognitiva aplicada a padrões de pensamento negativos.';

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
  });

  test.afterAll(async () => {
    await page.close();
  });

  test('01 · Login', async () => {
    await page.goto(BASE_URL + '/login');
    await page.waitForLoadState('networkidle');
    
    await page.locator('input[name="email"]').fill(EMAIL);
    await page.locator('input[name="password"]').fill(PASSWORD);
    await page.locator('button[type="submit"]').click();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);
  });

  test('02 · Cadastrar paciente', async () => {
    patientName = `Playwright Teste ${Date.now()}`;

    await page.goto(BASE_URL + '/patients/new');
    await page.waitForLoadState('networkidle');

    await page.fill('input[name="name"]', patientName);
    await page.fill('input[name="gender"]', 'masculino');
    await page.fill('input[name="ethnicity"]', 'branca');
    await page.fill('input[name="occupation"]', 'estudante');
    await page.fill('input[name="education"]', 'superior');
    await page.fill('textarea[name="notes"]', 'Paciente de teste automatizado');
    await page.click('button[type="submit"]');

    await page.waitForURL(/\/patients\/[^/]+/, { timeout: 30000 });
    patientID = await getPatientIDFromURL(page);
    expect(patientID).toBeTruthy();

    await assertShellLayout(page);
    await expect(page.locator('h1').filter({ hasText: patientName })).toBeVisible();
  });

  test('03 · Dashboard exibe paciente recém-criado', async () => {
    await page.goto(BASE_URL + '/dashboard');
    await assertShellLayout(page);
    await expect(page.locator('body').getByText(patientName)).toBeVisible();
  });

  test('04 · Busca na topbar encontra o paciente', async () => {
    await page.goto(BASE_URL + '/dashboard');
    const searchInput = page.locator('#shell-patient-search');

    await searchInput.fill(patientName.split(' ')[0]);
    await searchInput.press('Enter');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    await expect(page.locator('#main-content').getByText(patientName)).toBeVisible();
  });

  test('05 · Criar sessão clínica para o paciente', async () => {
    // Navegar primeiro para que o fetch use o origin correto do servidor de teste
    await page.goto(BASE_URL + '/patients/' + patientID + '/sessions/new');
    await page.waitForLoadState('networkidle');

    const today = new Date().toISOString().split('T')[0];

    // Fetch com URLSearchParams (application/x-www-form-urlencoded) — compatível com r.ParseForm() do Go
    // Usa /session/ com trailing slash para evitar redirect 301 que converte POST→GET
    const result = await page.evaluate(
      async ({ patientID, date, summary }: { patientID: string; date: string; summary: string }) => {
        const body = new URLSearchParams();
        body.append('patient_id', patientID);
        body.append('date', date);
        body.append('summary', summary);
        const resp = await fetch('/session/', { method: 'POST', body });
        const text = await resp.text();
        // Buscar o session ID no HTML completo
        const sessionMatch = text.match(/\/session\/([^/"'\s]+)\/(?:observations|interventions|update)/);
        return { status: resp.status, url: resp.url, sessionIDFound: sessionMatch ? sessionMatch[1] : '', htmlHead: text.substring(0, 300) };
      },
      { patientID, date: today, summary: 'Sessão de teste via Playwright' }
    );

    sessionID = result.sessionIDFound;
    expect(sessionID, `Session ID não encontrado — status ${result.status}, url ${result.url}`).toBeTruthy();
  });

  test('06 · Sessão aparece no perfil do paciente', async () => {
    await page.goto(BASE_URL + '/patients/' + patientID);
    await assertShellLayout(page);
    await expect(page.locator('h1').filter({ hasText: patientName })).toBeVisible();
    // A sessão recém-criada deve aparecer no perfil (lista de sessões ou timeline)
    await expect(page.locator('body')).toContainText('Sessão de teste via Playwright');
  });

  test('07 · Registrar observação na sessão via API', async () => {
    // Cria observação via fetch (mesmo padrão do teste 05 — usa cookies do browser)
    const obsResult = await page.evaluate(
      async ({ sessionID, content }: { sessionID: string; content: string }) => {
        const body = new URLSearchParams();
        body.append('content', content);
        const resp = await fetch(`/session/${sessionID}/observations`, { method: 'POST', body });
        return { status: resp.status, html: (await resp.text()).substring(0, 200) };
      },
      { sessionID, content: OBS_TEXT }
    );
    expect(obsResult.status, `Observação criada — status ${obsResult.status}`).toBe(200);
  });

  test('08 · Registrar intervenção na sessão via API', async () => {
    const intResult = await page.evaluate(
      async ({ sessionID, content }: { sessionID: string; content: string }) => {
        const body = new URLSearchParams();
        body.append('content', content);
        const resp = await fetch(`/session/${sessionID}/interventions`, { method: 'POST', body });
        return { status: resp.status, html: (await resp.text()).substring(0, 200) };
      },
      { sessionID, content: INT_TEXT }
    );
    expect(intResult.status, `Intervenção criada — status ${intResult.status}`).toBe(200);
  });

  test('09 · Reload da sessão mostra registros persistidos', async () => {
    // Verifica persistência navegando à sessão e checando conteúdo via API
    const sessionHTML = await page.evaluate(
      async ({ sessionID }: { sessionID: string }) => {
        const resp = await fetch(`/session/${sessionID}/edit`);
        return resp.text();
      },
      { sessionID }
    );
    expect(sessionHTML).toContain(OBS_TEXT.substring(0, 30));
    expect(sessionHTML).toContain(INT_TEXT.substring(0, 30));
  });

  test('10 · Prontuário do paciente exibe a sessão com registros', async () => {
    await page.goto(BASE_URL + '/patients/' + patientID + '/history');
    await page.waitForLoadState('domcontentloaded');
    await assertShellLayout(page);

    // Timeline deve conter a sessão
    await expect(page.locator('body')).toContainText('Sessão de teste via Playwright');
  });
});