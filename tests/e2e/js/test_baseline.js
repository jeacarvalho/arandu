#!/usr/bin/env node
/**
 * Teste Visual - Baseline das páginas antes da migração
 * Captura screenshots das páginas atuais (CSS v1) para comparação
 */

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:8080';
const EMAIL = process.env.E2E_TEST_EMAIL || 'arandu_e2e@test.com';
const PASSWORD = process.env.E2E_TEST_PASS || 'test123456';
const OUTPUT_DIR = '/home/s015533607/Documentos/desenv/arandu/tmp/migration_baseline';

// Garantir diretório de saída
if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

const ROUTES = [
    { name: 'dashboard', path: '/dashboard', fullPage: true },
    { name: 'patients_list', path: '/patients', fullPage: true },
    { name: 'patient_detail', path: '/patients/ec07a497-e9d8-4971-8213-61567766c3d0', fullPage: true },
];

async function captureBaseline() {
    console.log('📸 Capturando baseline das páginas (CSS v1)...\n');
    
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext({ 
        viewport: { width: 1440, height: 900 },
        deviceScaleFactor: 1
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
        console.log('   ✅ Login realizado\n');
        
        // Capturar cada rota
        for (const route of ROUTES) {
            console.log(`📷 Capturando: ${route.name}`);
            
            await page.goto(`${BASE_URL}${route.path}`);
            await page.waitForLoadState('networkidle');
            await page.waitForTimeout(2000); // Esperar renderização
            
            const screenshotPath = path.join(OUTPUT_DIR, `${route.name}_baseline.png`);
            await page.screenshot({ 
                path: screenshotPath,
                fullPage: route.fullPage 
            });
            
            // Capturar informações de CSS
            const cssInfo = await page.evaluate(() => {
                const cssFiles = Array.from(document.querySelectorAll('link[rel="stylesheet"]'))
                    .map(link => link.href);
                return {
                    cssFiles,
                    hasV1: cssFiles.some(href => href.includes('tailwind.css') && !href.includes('v2')),
                    hasV2: cssFiles.some(href => href.includes('tailwind-v2.css'))
                };
            });
            
            console.log(`   ✅ Salvo: ${screenshotPath}`);
            console.log(`   📄 CSS: ${cssInfo.hasV1 ? 'v1' : ''} ${cssInfo.hasV2 ? 'v2' : ''}\n`);
            
            // Salvar info em JSON
            const infoPath = path.join(OUTPUT_DIR, `${route.name}_info.json`);
            fs.writeFileSync(infoPath, JSON.stringify({
                route: route.name,
                url: route.path,
                timestamp: new Date().toISOString(),
                cssFiles: cssInfo.cssFiles,
                layout: cssInfo.hasV2 ? 'Shell v2' : 'Base v1'
            }, null, 2));
        }
        
        console.log('\n✅ Baseline capturado com sucesso!');
        console.log(`📁 Arquivos salvos em: ${OUTPUT_DIR}`);
        
    } catch (error) {
        console.error('\n❌ Erro:', error.message);
        process.exit(1);
    } finally {
        await browser.close();
    }
}

captureBaseline();
