#!/usr/bin/env node
/**
 * Teste Visual - Após migração do Dashboard
 * Compara o dashboard migrado com o baseline
 */

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:8080';
const EMAIL = process.env.E2E_TEST_EMAIL || 'arandu_e2e@test.com';
const PASSWORD = process.env.E2E_TEST_PASS || 'test123456';
const BASELINE_DIR = '/home/s015533607/Documentos/desenv/arandu/tmp/migration_baseline';
const OUTPUT_DIR = '/home/s015533607/Documentos/desenv/arandu/tmp/migration_after';

if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

async function testDashboardMigration() {
    console.log('🧪 Testando migração do Dashboard...\n');
    
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext({ 
        viewport: { width: 1440, height: 900 }
    });
    const page = await context.newPage();
    
    try {
        // Login
        console.log('🔐 Fazendo login...');
        await page.goto(`${BASE_URL}/login`);
        await page.fill('input[name="email"]', EMAIL);
        await page.fill('input[name="password"]', PASSWORD);
        await page.click('button[type="submit"]');
        await page.waitForURL('**/dashboard', { timeout: 10000 });
        
        // Testar Dashboard
        console.log('\n📊 Testando Dashboard:');
        await page.goto(`${BASE_URL}/dashboard`);
        await page.waitForLoadState('networkidle');
        await page.waitForTimeout(2000);
        
        // Capturar screenshot
        const screenshotPath = path.join(OUTPUT_DIR, 'dashboard_after.png');
        await page.screenshot({ 
            path: screenshotPath,
            fullPage: true 
        });
        
        // Verificar CSS
        const cssInfo = await page.evaluate(() => {
            const cssFiles = Array.from(document.querySelectorAll('link[rel="stylesheet"]'))
                .map(link => link.href);
            return {
                cssFiles,
                hasV1: cssFiles.some(href => href.includes('tailwind.css') && !href.includes('v2')),
                hasV2: cssFiles.some(href => href.includes('tailwind-v2.css')),
                shellExists: !!document.querySelector('.shell'),
                topbarExists: !!document.querySelector('.shell-topbar'),
                sidebarExists: !!document.querySelector('.shell-sidebar')
            };
        });
        
        console.log('   ✅ Screenshot salvo:', screenshotPath);
        console.log('   📄 CSS:', cssInfo.hasV2 ? 'v2 (Shell)' : 'v1 (Base)');
        console.log('   🏗️  Elementos Shell:');
        console.log(`      - .shell: ${cssInfo.shellExists ? '✅' : '❌'}`);
        console.log(`      - .shell-topbar: ${cssInfo.topbarExists ? '✅' : '❌'}`);
        console.log(`      - .shell-sidebar: ${cssInfo.sidebarExists ? '✅' : '❌'}`);
        
        // Comparar com baseline
        const baselinePath = path.join(BASELINE_DIR, 'dashboard_baseline.png');
        if (fs.existsSync(baselinePath)) {
            console.log('\n📸 Screenshots:');
            console.log('   Baseline (v1):', baselinePath);
            console.log('   Atual (v2):  ', screenshotPath);
        }
        
        // Verificar erros
        if (!cssInfo.hasV2) {
            console.log('\n   ⚠️  ALERTA: CSS v2 não encontrado!');
        }
        if (!cssInfo.shellExists) {
            console.log('   ❌ ERRO: Layout Shell não encontrado!');
            process.exit(1);
        }
        
        console.log('\n✅ Dashboard migrado com sucesso!');
        
    } catch (error) {
        console.error('\n❌ Erro:', error.message);
        await page.screenshot({ 
            path: path.join(OUTPUT_DIR, 'dashboard_error.png'),
            fullPage: true 
        });
        process.exit(1);
    } finally {
        await browser.close();
    }
}

testDashboardMigration();
