import { Page, expect } from '@playwright/test';

export async function assertShellLayout(page: Page) {
  await expect(page.locator('header.sabio-topbar')).toBeVisible();
  await expect(page.locator('#main-content')).toBeVisible();
}

export async function assertSidebarItem(page: Page, item: string) {
  await expect(page.locator('.sabio-nav').getByText(item)).toBeVisible();
}

export async function assertBreadcrumb(page: Page, items: string[]) {
  for (const item of items) {
    await expect(page.locator('#shell-breadcrumb').getByText(item)).toBeVisible();
  }
}