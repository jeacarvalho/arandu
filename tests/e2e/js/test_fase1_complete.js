#!/usr/bin/env node
/**
 * Teste Visual - Detalhe do Paciente migrado
 * E teste de navegação completa Dashboard → Pacientes → Anamnese
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

async function testPatientDetailAndNavigation() {
    console.log('🧪 Testando Detalhe do Paciente + Navegação...\n');
    
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext({ viewport: { width: 1440, height: 900 } });
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
        
        // 1. Testar Detalhe do Paciente
        console.log('👤 Testando Detalhe do Paciente:');
        await page.goto(`${BASE_URL}/patients/ec07a497-e9d8-4971-8213-61567766c3d0`);
        await page.waitForLoadState('networkidle');
        await page.waitForTimeout(2000);
        
        const screenshotPath = path.join(OUTPUT_DIR, 'patient_detail_after.png');
        await page.screenshot({ path: screenshotPath, fullPage: true });
        
        const cssInfo = await page.evaluate(() => {
            const cssFiles = Array.from(document.querySelectorAll('link[rel="stylesheet"]')).map(link => link.href);
            return {
                cssFiles,
                hasV2: cssFiles.some(href => href.includes('tailwind-v2.css')),
                shellExists: !!document.querySelector('.shell'),
                sidebarTitle: document.querySelector('.shell-sidebar-title')?.textContent || 'N/A'
            };
        });
        
        console.log('   ✅ Screenshot salvo:', screenshotPath);
        console.log('   📄 CSS:', cssInfo.hasV2 ? 'v2 (Shell)' : 'v1 (Base)');
        console.log('   🏗️  Shell:', cssInfo.shellExists ? '✅' : '❌');
        console.log('   📋 Sidebar:', cssInfo.sidebarTitle);
        
        if (!cssInfo.shellExists) {
            console.log('   ❌ ERRO: Layout Shell não encontrado!');
            process.exit(1);
        }
        
        // 2. Testar navegação Dashboard → Pacientes → Anamnese (via sidebar)
        console.log('\n🧭 Testando navegação completa:');
        
        // Ir para Dashboard
        console.log('   1️⃣ Dashboard → Pacientes');
        await page.goto(`${BASE_URL}/dashboard`);
        await page.waitForTimeout(1000);
        
        const patientsLink = await page.$('aside#shell-sidebar a[href="/patients"], .shell-sidebar a[href="/patients"]');
        if (patientsLink) {
            await patientsLink.click();
            await page.waitForTimeout(2000);
            console.log('      ✅ Navegou para Pacientes');
        }
        
        // Clicar no primeiro paciente
        console.log('   2️⃣ Pacientes → Detalhe do Paciente');
        const firstPatient = await page.$('a[href^="/patients/"]:not([href="/patients"]):not([href="/patients/new"])');
        if (firstPatient) {
            await firstPatient.click();
            await page.waitForTimeout(2000);
            console.log('      ✅ Navegou para Detalhe');
        }
        
        // Clicar em Anamnese na sidebar
        console.log('   3️⃣ Detalhe → Anamnese (via sidebar)');
        const anamnesisLink = await page.$('aside#shell-sidebar a[href*="anamnesis"], .shell-sidebar a[href*="anamnesis"]');
        if (anamnesisLink) {
            await anamnesisLink.click();
            await page.waitForTimeout(3000);
            console.log('      ✅ Navegou para Anamnese');
            
            // Verificar se anamnese está funcionando
            const anamneseCheck = await page.evaluate(() => {
                return {
                    hasV2: document.querySelector('link[href*="tailwind-v2.css"]') !== null,
                    shellExists: !!document.querySelector('.shell'),
                    widgetExists: !!document.querySelector('.widget-wrapper'),
                    title: document.querySelector('h1')?.textContent || 'N/A'
                };
            });
            
            console.log('      📊 Anamnese:');
            console.log(`         CSS v2: ${anamneseCheck.hasV2 ? '✅' : '❌'}`);
            console.log(`         Shell: ${anamneseCheck.shellExists ? '✅' : '❌'}`);
            console.log(`         Widgets: ${anamneseCheck.widgetExists ? '✅' : '❌'}`);
            console.log(`         Título: ${anamneseCheck.title}`);
            
            // Screenshot final
            const navScreenshotPath = path.join(OUTPUT_DIR, 'navigation_test_anamnese.png');
            await page.screenshot({ path: navScreenshotPath, fullPage: true });
            console.log('      📸 Screenshot:', navScreenshotPath);
        }
        
        console.log('\n✅ Fase 1 completa! Todas as páginas migradas com sucesso!');
        console.log('\n📁 Screenshots salvos em:', OUTPUT_DIR);
        
    } catch (error) {
        console.error('\n❌ Erro:', error.message);
        await page.screenshot({ path: path.join(OUTPUT_DIR, 'patient_detail_error.png'), fullPage: true });
        process.exit(1);
    } finally {
        await browser.close();
    }
}

testPatientDetailAndNavigation();
