import { Page, expect } from '@playwright/test';

export async function login(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.fill('input[name="email"]', email);
  await page.fill('input[name="password"]', password);
  await page.click('button[type="submit"]');
  await page.waitForURL(/\/(dashboard|patients)/);
}

export async function assertShellLayout(page: Page) {
  await expect(page.locator('header.sabio-topbar')).toBeVisible();
  await expect(page.locator('#main-content')).toBeVisible();
}

export async function assertSidebarNav(page: Page, expectedItems: string[]) {
  for (const item of expectedItems) {
    await expect(page.locator('.sabio-nav').getByText(item)).toBeVisible();
  }
}

export async function getPatientIDFromURL(page: Page): Promise<string> {
  const match = page.url().match(/\/patients\/([^/?]+)/);
  return match ? match[1] : '';
}