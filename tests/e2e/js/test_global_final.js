#!/usr/bin/env node
/**
 * Teste Global - Verifica se TODAS as páginas usam CSS v2
 * Após migração completa
 */

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:8080';
const EMAIL = process.env.E2E_TEST_EMAIL || 'arandu_e2e@test.com';
const PASSWORD = process.env.E2E_TEST_PASS || 'test123456';
const OUTPUT_DIR = '/home/s015533607/Documentos/desenv/arandu/tmp/final_validation';

if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

// Rotas para testar
const ROUTES = [
    { name: 'Dashboard', url: '/dashboard', type: 'main' },
    { name: 'Lista Pacientes', url: '/patients', type: 'main' },
    { name: 'Detalhe Paciente', url: '/patients/ec07a497-e9d8-4971-8213-61567766c3d0', type: 'patient' },
    { name: 'Novo Paciente', url: '/patients/new', type: 'main' },
    { name: 'Anamnese', url: '/patients/ec07a497-e9d8-4971-8213-61567766c3d0/anamnesis', type: 'patient' },
    { name: 'Prontuário', url: '/patients/ec07a497-e9d8-4971-821313-61567766c3d0/history', type: 'patient' },
    { name: 'Nova Sessão', url: '/patients/ec07a497-e9d8-4971-8213-61567766c3d0/sessions/new', type: 'patient' },
];

async function checkPageCSS(page, url) {
    await page.goto(`${BASE_URL}${url}`, { waitUntil: 'networkidle' });
    await page.waitForTimeout(1500);
    
    const cssInfo = await page.evaluate(() => {
        const cssFiles = Array.from(document.querySelectorAll('link[rel="stylesheet"]'))
            .map(link => link.href);
        return {
            files: cssFiles,
            hasV1: cssFiles.some(href => href.includes('tailwind.css') && !href.includes('v2')),
            hasV2: cssFiles.some(href => href.includes('tailwind-v2.css'))
        };
    });
    
    return cssInfo;
}

async function runGlobalTest() {
    console.log('🌍 TESTE GLOBAL - Migração CSS v2 Completa\n');
    console.log('=' .repeat(60));
    
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext({ viewport: { width: 1440, height: 900 } });
    const page = await context.newPage();
    
    const results = {
        timestamp: new Date().toISOString(),
        pages: [],
        allV2: true,
        errors: []
    };

    try {
        // Login
        console.log('\n🔐 Login...');
        await page.goto(`${BASE_URL}/login`);
        await page.fill('input[name="email"]', EMAIL);
        await page.fill('input[name="password"]', PASSWORD);
        await page.click('button[type="submit"]');
        await page.waitForURL('**/dashboard', { timeout: 10000 });
        console.log('   ✅ OK\n');

        // Testar cada rota
        console.log('📋 Verificando páginas:\n');
        
        for (const route of ROUTES) {
            process.stdout.write(`   ${route.name}... `);
            
            try {
                const cssInfo = await checkPageCSS(page, route.url);
                const usesV2 = cssInfo.hasV2 && !cssInfo.hasV1;
                
                results.pages.push({
                    name: route.name,
                    url: route.url,
                    usesV2: usesV2,
                    cssFiles: cssInfo.files.filter(f => f.includes('tailwind'))
                });
                
                if (usesV2) {
                    console.log('✅ v2');
                } else {
                    console.log('❌ v1');
                    results.allV2 = false;
                }
                
            } catch (error) {
                console.log(`❌ ERRO: ${error.message}`);
                results.pages.push({
                    name: route.name,
                    url: route.url,
                    error: error.message
                });
                results.allV2 = false;
            }
        }

        // Screenshot final
        await page.goto(`${BASE_URL}/dashboard`);
        await page.screenshot({ path: path.join(OUTPUT_DIR, 'final_dashboard.png') });
        
        // Salvar relatório
        const reportPath = path.join(OUTPUT_DIR, 'final_report.json');
        fs.writeFileSync(reportPath, JSON.stringify(results, null, 2));
        
        // Resumo
        console.log('\n' + '='.repeat(60));
        console.log('\n📊 RESUMO:\n');
        
        const v2Count = results.pages.filter(p => p.usesV2).length;
        const v1Count = results.pages.filter(p => !p.usesV2 && !p.error).length;
        const errorCount = results.pages.filter(p => p.error).length;
        
        console.log(`   Total de páginas: ${results.pages.length}`);
        console.log(`   ✅ CSS v2: ${v2Count}`);
        console.log(`   ❌ CSS v1: ${v1Count}`);
        console.log(`   ⚠️  Erros: ${errorCount}`);
        
        if (results.allV2) {
            console.log('\n🎉 SUCESSO! Todas as páginas usam CSS v2!');
        } else {
            console.log('\n⚠️  ATENÇÃO: Algumas páginas ainda usam CSS v1');
            console.log('\n   Páginas com v1:');
            results.pages.filter(p => !p.usesV2 && !p.error).forEach(p => {
                console.log(`      - ${p.name}: ${p.cssFiles.join(', ')}`);
            });
        }
        
        console.log(`\n📄 Relatório: ${reportPath}`);
        console.log(`📸 Screenshot: ${OUTPUT_DIR}/final_dashboard.png`);
        
    } catch (error) {
        console.error('\n❌ Erro fatal:', error.message);
        process.exit(1);
    } finally {
        await browser.close();
    }
}

runGlobalTest();
