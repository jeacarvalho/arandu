#!/usr/bin/env python3
import asyncio
from playwright.async_api import async_playwright

cookies = [
    {
        "name": "arandu_session",
        "value": "c9f8257b-d5e9-4b52-b78d-c549f71e9c5c",
        "domain": "localhost",
        "path": "/",
    }
]


async def capture_dashboard():
    async with async_playwright() as p:
        browser = await p.chromium.launch()
        context = await browser.new_context(viewport={"width": 1440, "height": 900})
        await context.add_cookies(cookies)
        page = await context.new_page()
        await page.goto("http://localhost:8080/dashboard")
        await page.screenshot(
            path="screenshots/dashboard_authenticated.png", full_page=True
        )
        print("Screenshot saved: screenshots/dashboard_authenticated.png")
        await browser.close()


if __name__ == "__main__":
    asyncio.run(capture_dashboard())
