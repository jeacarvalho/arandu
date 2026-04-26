# TASK 20260426_170127
## Teste E2E Playwright — Fluxo Clínico Completo com Tenant em Memória

**Status:** PRONTO_PARA_IMPLEMENTACAO

---

## 🎯 Objetivo

Implementar um teste E2E com Playwright que valide o fluxo clínico completo do Arandu — do cadastro de paciente ao registro e consulta de sessão — usando um servidor real com banco SQLite temporário descartado ao final. O teste deve validar **layout e comportamento visual** em cada tela que percorre.

---

## 🏗️ Contexto do sistema

**Leia antes de qualquer código:**

Consulte as skills:
- `arandu-architecture` — estrutura de pastas, multi-tenancy, padrão de handlers
- `arandu-master` — protocolo de delegação e qualidade de entrega

**Stack:** Go 1.22+ · SQLite (multi-tenant) · HTMX · Templ · Playwright (TypeScript)

**Padrão de testes existente:**
- Testes E2E Go: `tests/e2e/*_test.go` usando `httptest.NewServer`
- Padrão de setup: `os.MkdirTemp` → `sqlite.NewCentralDB` → `sqlite.NewDB` (tenant) → `Migrate()` → router completo → `httptest.NewServer`
- Referência obrigatória: leia `tests/e2e/e2e_full_workflow_test.go` linhas 1–180 ANTES de escrever qualquer código

**Multi-tenancy:**
- Central DB: usuários, tenants, mapeamento user→db_path
- Tenant DB: dados clínicos (pacientes, sessões, observações, intervenções, agendamentos)
- Tenant é resolvido pelo middleware de autenticação a partir do cookie de sessão

**Autenticação:**
- Login via `POST /login` com campos `email` e `password`
- Resposta: set-cookie com session token
- O middleware injeta o tenant DB no contexto de cada request

---

## 📁 Arquivos a criar

```
playwright.config.ts                          ← config Playwright (root do projeto)
tests/e2e/playwright_runner_test.go           ← orquestrador Go: sobe servidor, roda Playwright
tests/e2e/playwright/
  clinical_workflow.spec.ts                   ← spec principal do fluxo clínico
  helpers/
    auth.ts                                   ← helper de login
    layout.ts                                 ← assertions de layout reutilizáveis
```

## 📁 Arquivos a modificar

```
Makefile                                      ← integrar playwright em "make test"
```

---

## 🔩 Arquitetura de integração

### Go como orquestrador (`playwright_runner_test.go`)

```go
// tests/e2e/playwright_runner_test.go
package e2e

import (
    "os"
    "os/exec"
    "testing"
    // ... imports do padrão existente
)

func TestPlaywrightClinicalWorkflow(t *testing.T) {
    // 1. Cria infra temporária (IGUAL ao padrão existente dos outros e2e tests)
    tmpDir := t.TempDir()
    centralDB := setupCentralDB(t, tmpDir)        // função igual à existente
    tenantDB, tenantID := setupTenantDB(t, tmpDir, centralDB) // idem

    // 2. Cria usuário de teste no central DB
    // Email/senha fixos para o teste, NÃO o arandu@test.com
    testEmail := "playwright_e2e@arandu.internal"
    testPassword := "playwright_test_2026"
    createTestUser(t, centralDB, tenantID, testEmail, testPassword)

    // 3. Monta router COMPLETO (incluindo AuthHandler)
    // Siga o padrão de tests/e2e/e2e_full_workflow_test.go mas INCLUA o AuthHandler
    router := buildFullRouter(t, centralDB, tenantDB)
    server := httptest.NewServer(router)
    t.Cleanup(server.Close)

    // 4. Verifica se Playwright está disponível
    if _, err := exec.LookPath("npx"); err != nil {
        t.Skip("npx not available — skipping Playwright tests")
    }

    // 5. Executa Playwright como subprocess
    cmd := exec.Command("npx", "playwright", "test",
        "--project=chromium",
        "--reporter=list",
        "tests/e2e/playwright/",
    )
    cmd.Env = append(os.Environ(),
        "PLAYWRIGHT_BASE_URL="+server.URL,
        "E2E_EMAIL="+testEmail,
        "E2E_PASSWORD="+testPassword,
    )
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        t.Fatalf("Playwright tests failed: %v", err)
    }
}
```

**Criação do usuário de teste:**
```go
func createTestUser(t *testing.T, centralDB *sql.DB, tenantID, email, password string) {
    t.Helper()
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    userID := uuid.New().String()
    _, err := centralDB.Exec(`
        INSERT INTO users (id, email, password_hash, tenant_id, created_at, updated_at)
        VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
    `, userID, email, string(hash), tenantID)
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
}
```

> Inspecione o schema de `users` e `tenants` no central DB (migration `0001_*.up.sql`) para garantir que os campos batem exatamente.

### `playwright.config.ts`

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e/playwright',
  timeout: 60_000,
  expect: { timeout: 10_000 },
  fullyParallel: false,         // fluxo clínico é sequencial
  retries: 0,
  reporter: [['list']],
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL ?? 'http://localhost:8080',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'off',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'], headless: true },
    },
  ],
});
```

---

## 🧪 Spec do fluxo clínico (`clinical_workflow.spec.ts`)

O teste é **sequencial** (`test.describe.serial`). Use `test.step()` para cada sub-passo. Salve o `patientID` e `sessionID` entre steps.

### Helper de layout (`helpers/layout.ts`)

Crie assertions reutilizáveis que serão chamadas em CADA tela:

```typescript
export async function assertShellLayout(page: Page) {
  // Sidebar presente e visível
  await expect(page.locator('aside.arandu-sidebar')).toBeVisible();
  // Topbar presente
  await expect(page.locator('header.sabio-topbar')).toBeVisible();
  // Main content presente
  await expect(page.locator('#main-content')).toBeVisible();
  // Sidebar NÃO sobrepõe o main content (verificação básica de z-index/layout)
  const sidebar = page.locator('aside.arandu-sidebar');
  const main = page.locator('#main-content');
  const sidebarBox = await sidebar.boundingBox();
  const mainBox = await main.boundingBox();
  if (sidebarBox && mainBox) {
    // Main content começa DEPOIS da sidebar (layout correto)
    expect(mainBox.x).toBeGreaterThan(sidebarBox.x);
  }
}

export async function assertSidebarNav(page: Page, expectedItems: string[]) {
  // Verifica itens de navegação do sidebar
  for (const item of expectedItems) {
    await expect(page.locator('.sabio-nav').getByText(item)).toBeVisible();
  }
}
```

### Fluxo completo

```typescript
import { test, expect, Page } from '@playwright/test';
import { assertShellLayout, assertSidebarNav } from './helpers/layout';

const BASE_URL = process.env.PLAYWRIGHT_BASE_URL!;
const EMAIL = process.env.E2E_EMAIL!;
const PASSWORD = process.env.E2E_PASSWORD!;

test.describe.serial('Fluxo Clínico Completo', () => {
  let page: Page;
  let patientName: string;
  let patientID: string;

  test.beforeAll(async ({ browser }) => {
    page = await browser.newPage();
  });

  test.afterAll(async () => {
    await page.close();
  });

  // ──────────────────────────────────────────────
  // LOGIN
  // ──────────────────────────────────────────────
  test('01 · Login', async () => {
    await page.goto('/login');
    await page.fill('input[name="email"]', EMAIL);
    await page.fill('input[name="password"]', PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForURL(/\/(dashboard|patients)/);
  });

  // ──────────────────────────────────────────────
  // A) CADASTRAR PACIENTE
  // ──────────────────────────────────────────────
  test('02 · Cadastrar paciente', async () => {
    patientName = `Playwright Teste ${Date.now()}`;

    await page.goto('/patients/new');
    await assertShellLayout(page);
    // Sidebar padrão: sem contexto de paciente
    await assertSidebarNav(page, ['Dashboard', 'Pacientes', 'Agenda', 'Prontuários']);

    await page.fill('input[name="name"]', patientName);
    // Preencha outros campos obrigatórios conforme o formulário de novo paciente
    await page.click('button[type="submit"]');

    // Aguarda redirect para o perfil do paciente
    await page.waitForURL(/\/patients\/[^/]+/);
    patientID = page.url().match(/\/patients\/([^/?]+)/)?.[1] ?? '';
    expect(patientID).toBeTruthy();

    await assertShellLayout(page);
    // Sidebar muda para contexto do paciente (Resumo, Anamnese, Prontuário, Plano)
    await assertSidebarNav(page, ['Resumo', 'Anamnese', 'Prontuário']);
    // Nome do paciente aparece na tela
    await expect(page.locator('h1, .patient-name, [data-testid="patient-name"]')
      .filter({ hasText: patientName })).toBeVisible();
  });

  // ──────────────────────────────────────────────
  // B) PACIENTE APARECE NO DASHBOARD
  // ──────────────────────────────────────────────
  test('03 · Dashboard exibe paciente recém-criado', async () => {
    await page.goto('/dashboard');
    await assertShellLayout(page);
    await assertSidebarNav(page, ['Dashboard', 'Pacientes', 'Agenda', 'Prontuários']);
    // Paciente deve aparecer na lista ou seção de pacientes do dashboard
    await expect(page.locator('body').getByText(patientName)).toBeVisible();
  });

  // ──────────────────────────────────────────────
  // C) BUSCA NA TOPBAR
  // ──────────────────────────────────────────────
  test('04 · Busca na topbar encontra o paciente', async () => {
    await page.goto('/dashboard');
    const searchInput = page.locator('#shell-patient-search, input[name="q"]').first();
    await searchInput.fill(patientName.split(' ')[0]); // busca pelo primeiro nome
    await searchInput.press('Enter');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('body').getByText(patientName)).toBeVisible();
  });

  // ──────────────────────────────────────────────
  // D) CRIAR AGENDAMENTO (hoje, hora atual)
  // ──────────────────────────────────────────────
  test('05 · Criar agendamento para hoje', async () => {
    await page.goto('/agenda');
    await assertShellLayout(page);

    // Clica em "Novo agendamento"
    await page.click('button:has-text("Novo agendamento"), a:has-text("Novo agendamento")');
    await page.waitForSelector('form, [id="drawer-container"] form', { timeout: 5000 });

    // Preenche o formulário
    const now = new Date();
    const dateStr = now.toISOString().split('T')[0]; // YYYY-MM-DD
    const hour = now.getHours().toString().padStart(2, '0');
    const timeStr = `${hour}:00`;

    // Campo de data
    const dateInput = page.locator('input[name="date"]');
    await dateInput.fill(dateStr);

    // Campo de paciente (pode ser select ou autocomplete)
    const patientSelect = page.locator('select[name="patient_id"]');
    if (await patientSelect.count() > 0) {
      await patientSelect.selectOption({ label: patientName });
    } else {
      const patientInput = page.locator('input[name="patient_id"], input[placeholder*="paciente"]');
      await patientInput.fill(patientName);
      await page.locator('.autocomplete-option, [data-patient]')
        .filter({ hasText: patientName }).first().click();
    }

    // Campo de horário (pode ser radio buttons ou select)
    const startTimeSelect = page.locator('select[name="start_time"], input[name="start_time"]');
    if (await startTimeSelect.count() > 0) {
      if (await startTimeSelect.evaluate(el => el.tagName) === 'SELECT') {
        await startTimeSelect.selectOption(timeStr);
      } else {
        await startTimeSelect.fill(timeStr);
      }
    }

    // Salva
    await page.click('button[type="submit"]');
    await page.waitForLoadState('networkidle');

    // Verifica que o agendamento aparece na agenda
    await page.goto('/agenda');
    await assertShellLayout(page);
    await expect(page.locator('body').getByText(patientName)).toBeVisible();
  });

  // ──────────────────────────────────────────────
  // E) CONCLUIR ATENDIMENTO VIA AGENDA → SESSÃO
  // ──────────────────────────────────────────────
  test('06 · Concluir atendimento via agenda e ir para registro de sessão', async () => {
    await page.goto('/agenda');

    // Clica no card do agendamento do paciente
    const appointmentCard = page.locator('.sabio-week-appt, .sabio-day-appt, .sabio-month-appt-pill')
      .filter({ hasText: patientName }).first();
    await appointmentCard.click();

    // Modal de detalhe do agendamento
    await page.waitForSelector('#modal-container [class*="modal"], #modal-container .card', { timeout: 5000 });

    // Clica em "Concluir" ou "Concluir com sessão"
    const completeBtn = page.locator('#modal-container')
      .getByRole('button', { name: /conclu|complete|realizar/i }).first();
    await completeBtn.click();
    await page.waitForLoadState('networkidle');

    // Deve navegar para a página de registro de sessão
    await page.waitForURL(/\/session\/[^/]+\/edit|\/session\/new/);
    await assertShellLayout(page);
    // Sidebar mostra contexto do paciente
    await assertSidebarNav(page, ['Resumo', 'Prontuário']);
  });

  // ──────────────────────────────────────────────
  // F) REGISTRAR SESSÃO: OBSERVAÇÃO + CLASSIFICAÇÕES + INTERVENÇÃO
  // ──────────────────────────────────────────────
  test('07 · Registrar observação com duas classificações', async () => {
    // Verifica layout da página de sessão
    await assertShellLayout(page);
    await expect(page.locator('h1.sabio-session-title, h1').filter({ hasText: /registro|sessão/i })).toBeVisible();
    // Colunas Escuta e Ação presentes
    await expect(page.locator('text=Observações clínicas')).toBeVisible();
    await expect(page.locator('text=Intervenções terapêuticas')).toBeVisible();

    // Digita observação
    const obsText = 'Paciente relata melhora significativa na regulação emocional esta semana.';
    await page.fill('textarea[name="content"]', obsText);
    await page.click('button:has-text("Registrar")');
    await page.waitForLoadState('networkidle');

    // Observação aparece na lista
    await expect(page.locator('#observations-list').getByText(obsText)).toBeVisible();

    // Clica no botão de classificar da observação
    const obsItem = page.locator('.sabio-notes-item').filter({ hasText: obsText });
    await obsItem.locator('button[title="Classificar"]').click();
    await page.waitForSelector('.tag-selector, [id^="selector-"]', { timeout: 5000 });

    // Seleciona 2 tags (primeira disponível de dois grupos diferentes)
    const tagLabels = page.locator('.badge.cursor-pointer, label.badge').all();
    const tags = await tagLabels;
    // Clica na primeira tag disponível
    if (tags.length > 0) await tags[0].click();
    // Clica na segunda tag de um grupo diferente (pula algumas)
    if (tags.length > 3) await tags[3].click();

    // Salva classificações
    await page.click('button[type="submit"]:has-text("Salvar")');
    await page.waitForLoadState('networkidle');

    // Verifica que badges aparecem
    const tagsContainer = page.locator(`[id$="-tags"]`).first();
    await expect(tagsContainer.locator('.badge')).toHaveCountGreaterThan(0);
  });

  test('08 · Registrar intervenção', async () => {
    const intervText = 'Técnica de reestruturação cognitiva aplicada para pensamentos automáticos negativos.';

    // Form de intervenção — segundo formulário da página
    const intervForms = page.locator('form[hx-post*="interventions"]');
    await intervForms.locator('textarea[name="content"]').fill(intervText);
    await intervForms.locator('button[type="submit"]').click();
    await page.waitForLoadState('networkidle');

    // Intervenção aparece na lista
    await expect(page.locator('#interventions-list').getByText(intervText)).toBeVisible();

    // Verifica layout final da sessão com conteúdo preenchido
    await assertShellLayout(page);
    await expect(page.locator('#observations-list .sabio-notes-item')).toHaveCount(1);
    await expect(page.locator('#interventions-list .sabio-notes-item')).toHaveCount(1);
  });

  // ──────────────────────────────────────────────
  // G) CONSULTAR PACIENTE E VALIDAR PERSISTÊNCIA
  // ──────────────────────────────────────────────
  test('09 · Consultar perfil do paciente e verificar sessão', async () => {
    await page.goto(`/patients/${patientID}`);
    await assertShellLayout(page);
    await assertSidebarNav(page, ['Resumo', 'Anamnese', 'Prontuário', 'Plano Terapêutico']);
    await expect(page.locator('body').getByText(patientName)).toBeVisible();
  });

  test('10 · Verificar observações e intervenções no prontuário', async () => {
    await page.goto(`/patients/${patientID}/history`);
    await assertShellLayout(page);

    // Sessão registrada aparece no histórico
    await expect(page.locator('body').getByText(/sessão|registro/i)).toBeVisible();

    // Navega para a sessão
    const sessionLink = page.locator('a[href*="/session/"]').first();
    if (await sessionLink.count() > 0) {
      await sessionLink.click();
      await page.waitForLoadState('networkidle');

      // Observação e intervenção registradas estão presentes
      await expect(page.locator('body')
        .getByText('Paciente relata melhora significativa')).toBeVisible();
      await expect(page.locator('body')
        .getByText('Técnica de reestruturação cognitiva')).toBeVisible();

      // Badges de classificação presentes na observação
      const tagBadges = page.locator('.badge-outline, [id$="-tags"] .badge');
      await expect(tagBadges.first()).toBeVisible();
    }

    await assertShellLayout(page);
  });
});
```

---

## 🔧 Integração com `make test`

### Makefile

Adicionar ao `Makefile` **sem quebrar o target `test` existente**:

```makefile
# Instala dependências Playwright se necessário
playwright-install:
	npx playwright install chromium --with-deps

# Roda apenas os testes Playwright (via Go runner)
test-playwright:
	go test ./tests/e2e/ -run TestPlaywrightClinicalWorkflow -v -timeout 120s

# Target "test" existente deve incluir playwright
# Verifique como "tests/run_parallel.sh" está estruturado e adicione
# a chamada ao TestPlaywrightClinicalWorkflow lá, OU adicione uma
# linha no Makefile:
test: test-unit test-playwright
```

> **IMPORTANTE**: Leia o `Makefile` e `tests/run_parallel.sh` existentes antes de modificar. Não quebre os targets existentes. Se o `test` já chama um script shell, adicione a chamada ao `TestPlaywrightClinicalWorkflow` dentro do script, não no Makefile diretamente.

---

## ✅ Critérios de aceite

**Compilação:**
- [ ] `go build ./tests/e2e/...` sem erros
- [ ] `npx playwright test --list` lista os 10 testes sem erro de configuração

**Comportamento:**
- [ ] `go test ./tests/e2e/ -run TestPlaywrightClinicalWorkflow -v` passa do início ao fim
- [ ] Cada step do fluxo passa isoladamente (não há dependência de estado externo)
- [ ] Ao final, o banco temporário é descartado (`t.TempDir()` cuida disso automaticamente)
- [ ] Não usa `arandu@test.com` nem o banco `arandu_central.db` de produção/dev
- [ ] `make test` inclui o teste Playwright e falha se ele falhar

**Layout (assertions em cada tela):**
- [ ] `aside.arandu-sidebar` visível em todas as telas pós-login
- [ ] `header.sabio-topbar` visível em todas as telas pós-login
- [ ] `#main-content` visível e posicionado à direita da sidebar
- [ ] Sidebar mostra nav global no dashboard e nav contextual (paciente) nas telas de paciente/sessão

**Fluxo:**
- [ ] Steps 01–10 passam em sequência sem intervenção manual
- [ ] Observação e intervenção criadas na sessão são encontradas na consulta do prontuário
- [ ] Badges de classificação são visíveis após registro

---

## 🚫 NÃO faça

- Não use `arandu@test.com` nem o tenant `9b33e4a0-...` de dev
- Não use `page.waitForTimeout()` para esperas fixas — use `waitForURL`, `waitForSelector`, `waitForLoadState`
- Não modifique o código de produção (handlers, migrations, domínio) para fazer os testes passarem
- Não use `test.only` — todos os 10 steps devem rodar
- Não chame `playwright test` diretamente no Makefile (o runner Go é o ponto de entrada)
- Não invente seletores CSS — inspecione o HTML real das páginas para confirmar classes e IDs

---

## 📎 Referências obrigatórias

1. **`tests/e2e/e2e_full_workflow_test.go`** — padrão de setup de DB e router para o runner Go
2. **`tests/e2e/http_patient_flow_test.go`** — padrão de construção de repositórios e services
3. **`cmd/arandu/main.go`** — como o AuthHandler é instanciado e registrado (para incluí-lo no router de teste)
4. **`internal/infrastructure/repository/sqlite/migrations/`** — prefixo numérico da última migration (para garantir que o migrator aplica todas)
5. **`web/components/session/`** — classes e IDs dos elementos da página de sessão para os seletores Playwright
6. **`web/components/layout/shell_layout.templ`** — IDs e classes do shell (sidebar, topbar, main-content) para os `assertShellLayout`
7. **`playwright.config.ts`** (a ser criado) — base URL vem de `process.env.PLAYWRIGHT_BASE_URL`

**Antes de escrever qualquer seletor Playwright**, leia o HTML da tela correspondente em busca dos IDs e classes reais. Não adivinhe seletores.
