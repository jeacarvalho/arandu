#!/usr/bin/env node
/**
 * Teste Comparativo - Layout v1 vs v2
 * Compara estilos detalhados entre páginas antigas e novas
 */

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:8080';
const EMAIL = process.env.E2E_TEST_EMAIL || 'arandu_e2e@test.com';
const PASSWORD = process.env.E2E_TEST_PASS || 'test123456';
const OUTPUT_DIR = '/home/s015533607/Documentos/desenv/arandu/tmp/style_comparison';

if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

// Extrai estilos detalhados de um elemento
async function extractDetailedStyles(page, selector, name) {
    const element = await page.$(selector);
    if (!element) {
        return { exists: false, name, selector };
    }

    const styles = await page.evaluate((sel) => {
        const el = document.querySelector(sel);
        if (!el) return null;
        
        const computed = window.getComputedStyle(el);
        const rect = el.getBoundingClientRect();
        
        return {
            // Cores
            backgroundColor: computed.backgroundColor,
            color: computed.color,
            borderColor: computed.borderColor,
            
            // Layout
            width: computed.width,
            height: computed.height,
            padding: computed.padding,
            margin: computed.margin,
            position: computed.position,
            display: computed.display,
            
            // Tipografia
            fontFamily: computed.fontFamily,
            fontSize: computed.fontSize,
            fontWeight: computed.fontWeight,
            lineHeight: computed.lineHeight,
            letterSpacing: computed.letterSpacing,
            textAlign: computed.textAlign,
            
            // Bordas
            border: computed.border,
            borderRadius: computed.borderRadius,
            borderWidth: computed.borderWidth,
            
            // Posicionamento
            top: computed.top,
            left: computed.left,
            right: computed.right,
            bottom: computed.bottom,
            zIndex: computed.zIndex,
            
            // Outros
            boxShadow: computed.boxShadow,
            opacity: computed.opacity,
            visibility: computed.visibility,
            
            // Posição real no viewport
            rect: {
                x: rect.x,
                y: rect.y,
                width: rect.width,
                height: rect.height
            }
        };
    }, selector);

    return {
        exists: true,
        name,
        selector,
        styles
    };
}

// Extrai informações dos itens da sidebar
async function extractSidebarItems(page) {
    return await page.evaluate(() => {
        const sidebar = document.querySelector('.shell-sidebar, aside, .sidebar');
        if (!sidebar) return null;
        
        const items = sidebar.querySelectorAll('a, .nav-item, .sidebar-nav-item');;
        return Array.from(items).map((item, index) => {
            const rect = item.getBoundingClientRect();
            const computed = window.getComputedStyle(item);
            return {
                index,
                text: item.textContent?.trim().substring(0, 30),
                href: item.getAttribute('href'),
                position: {
                    x: rect.x,
                    y: rect.y,
                    width: rect.width,
                    height: rect.height
                },
                styles: {
                    padding: computed.padding,
                    margin: computed.margin,
                    display: computed.display,
                    alignItems: computed.alignItems,
                    justifyContent: computed.justifyContent,
                    gap: computed.gap,
                    fontSize: computed.fontSize,
                    fontWeight: computed.fontWeight
                }
            };
        });
    });
}

// Testa a página "Novo Paciente" acessada via botão do dashboard
async function testNewPatientViaButton() {
    console.log('🧪 Testando "Novo Paciente" via botão do Dashboard...\n');
    
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext({ 
        viewport: { width: 1440, height: 900 }
    });
    const page = await context.newPage();

    const results = {
        timestamp: new Date().toISOString(),
        dashboard: {},
        newPatient: {},
        comparison: {}
    };

    try {
        // 1. Login
        console.log('🔐 Fazendo login...');
        await page.goto(`${BASE_URL}/login`);
        await page.fill('input[name="email"]', EMAIL);
        await page.fill('input[name="password"]', PASSWORD);
        await page.click('button[type="submit"]');
        await page.waitForURL('**/dashboard', { timeout: 10000 });
        console.log('   ✅ Login realizado\n');

        // 2. Capturar métricas do Dashboard (v2)
        console.log('📊 Dashboard (v2) - Métricas:');
        
        // Forçar refresh sem cache
        await page.reload({ waitUntil: 'networkidle' });
        await page.waitForTimeout(2000);
        
        results.dashboard = await capturePageMetrics(page, 'Dashboard v2');
        
        // Screenshot do dashboard
        await page.screenshot({ 
            path: path.join(OUTPUT_DIR, 'dashboard_v2.png'),
            fullPage: false 
        });

        // 3. Clicar no botão "Novo Paciente" no canvas (não na sidebar)
        console.log('\n🖱️  Clicando no botão "Novo Paciente" do Dashboard...');
        
        // Procurar botão de novo paciente no conteúdo principal (não sidebar)
        const novoPacienteBtn = await page.$('main a[href="/patients/new"], .shell-content a[href="/patients/new"], [data-testid="new-patient-btn"]');
        
        if (!novoPacienteBtn) {
            // Tentar encontrar por texto
            const buttons = await page.$$('main button, main a, .shell-content button, .shell-content a');
            for (const btn of buttons) {
                const text = await btn.textContent();
                if (text && text.toLowerCase().includes('paciente')) {
                    console.log(`   Encontrado botão: ${text.trim()}`);
                    await btn.click();
                    break;
                }
            }
        } else {
            await novoPacienteBtn.click();
        }

        // Aguardar navegação
        await page.waitForTimeout(3000);
        
        // Verificar URL atual
        const currentUrl = page.url();
        console.log(`   URL atual: ${currentUrl}`);

        // 4. Capturar métricas da página Novo Paciente
        console.log('\n📋 Página Novo Paciente - Métricas:');
        results.newPatient = await capturePageMetrics(page, 'Novo Paciente');
        
        // Screenshot
        await page.screenshot({ 
            path: path.join(OUTPUT_DIR, 'new_patient_current.png'),
            fullPage: true 
        });

        // 5. Comparar métricas
        console.log('\n🔍 COMPARAÇÃO DE ESTILOS:\n');
        results.comparison = compareMetrics(results.dashboard, results.newPatient);
        
        // 6. Salvar relatório completo
        const reportPath = path.join(OUTPUT_DIR, 'comparison_report.json');
        fs.writeFileSync(reportPath, JSON.stringify(results, null, 2));
        console.log(`\n📄 Relatório salvo: ${reportPath}`);

        // 7. Resumo dos problemas
        printSummary(results);

    } catch (error) {
        console.error('\n❌ Erro:', error.message);
        await page.screenshot({ 
            path: path.join(OUTPUT_DIR, 'error.png'),
            fullPage: true 
        });
    } finally {
        await browser.close();
    }

    return results;
}

// Captura métricas completas de uma página
async function capturePageMetrics(page, pageName) {
    const metrics = {
        pageName,
        url: page.url(),
        cssVersion: null,
        sidebar: {},
        topbar: {},
        mainContent: {},
        searchBox: {},
        sidebarItems: [],
        fonts: {},
        colors: {}
    };

    // Verificar versão CSS
    const cssInfo = await page.evaluate(() => {
        const cssFiles = Array.from(document.querySelectorAll('link[rel="stylesheet"]'))
            .map(link => link.href);
        return {
            files: cssFiles,
            hasV1: cssFiles.some(href => href.includes('tailwind.css') && !href.includes('v2')),
            hasV2: cssFiles.some(href => href.includes('tailwind-v2.css'))
        };
    });
    console.log('   CSS Files:', cssInfo.files);
    metrics.cssVersion = cssInfo.hasV2 ? 'v2' : (cssInfo.hasV1 ? 'v1' : 'unknown');
    console.log(`   CSS Version: ${metrics.cssVersion} (v1: ${cssInfo.hasV1}, v2: ${cssInfo.hasV2})`);

    // Sidebar
    metrics.sidebar = await extractDetailedStyles(
        page, 
        '.shell-sidebar, aside.sidebar, .sidebar-drawer, nav.sidebar',
        'Sidebar'
    );
    console.log(`   Sidebar: ${metrics.sidebar.exists ? '✅' : '❌'}`);
    if (metrics.sidebar.exists && metrics.sidebar.styles) {
        console.log(`      Background: ${metrics.sidebar.styles.backgroundColor}`);
        console.log(`      Width: ${metrics.sidebar.styles.width}`);
    }

    // Topbar
    metrics.topbar = await extractDetailedStyles(
        page,
        '.shell-topbar, header.topbar, .top-bar',
        'Topbar'
    );
    console.log(`   Topbar: ${metrics.topbar.exists ? '✅' : '❌'}`);

    // Logo
    const logo = await extractDetailedStyles(page, '.shell-logo, .logo, .brand', 'Logo');
    if (logo.exists) {
        console.log(`   Logo color: ${logo.styles?.color || 'N/A'}`);
    }

    // Search box
    metrics.searchBox = await extractDetailedStyles(
        page,
        '#patient-search, .search-input, input[type="search"]',
        'Search Box'
    );
    console.log(`   Search: ${metrics.searchBox.exists ? '✅' : '❌'}`);
    if (metrics.searchBox.exists && metrics.searchBox.styles) {
        console.log(`      Background: ${metrics.searchBox.styles.backgroundColor}`);
        console.log(`      Border: ${metrics.searchBox.styles.border}`);
    }

    // Main content
    metrics.mainContent = await extractDetailedStyles(
        page,
        '.shell-main, main, .main-content',
        'Main Content'
    );
    console.log(`   Main: ${metrics.mainContent.exists ? '✅' : '❌'}`);
    if (metrics.mainContent.exists && metrics.mainContent.styles) {
        console.log(`      Background: ${metrics.mainContent.styles.backgroundColor}`);
    }

    // Sidebar items
    metrics.sidebarItems = await extractSidebarItems(page);
    if (metrics.sidebarItems && metrics.sidebarItems.length > 0) {
        console.log(`   Sidebar items: ${metrics.sidebarItems.length}`);
        console.log(`   First item: "${metrics.sidebarItems[0]?.text?.substring(0, 30)}"`);
    }

    // Fontes usadas na página
    const fonts = await page.evaluate(() => {
        const elements = document.querySelectorAll('body, h1, h2, h3, p, a, button');
        const fontSet = new Set();
        elements.forEach(el => {
            const style = window.getComputedStyle(el);
            fontSet.add(style.fontFamily);
        });
        return Array.from(fontSet);
    });
    metrics.fonts = fonts;
    console.log(`   Fonts: ${fonts.slice(0, 2).join(', ')}${fonts.length > 2 ? '...' : ''}`);

    return metrics;
}

// Compara métricas entre duas páginas
function compareMetrics(dashboard, newPatient) {
    const differences = [];

    // Comparar versões CSS
    if (dashboard.cssVersion !== newPatient.cssVersion) {
        differences.push({
            type: 'CSS_VERSION',
            severity: 'HIGH',
            message: `Mudança de CSS: ${dashboard.cssVersion} → ${newPatient.cssVersion}`,
            dashboard: dashboard.cssVersion,
            newPatient: newPatient.cssVersion
        });
    }

    // Comparar backgrounds
    if (dashboard.sidebar?.styles?.backgroundColor !== newPatient.sidebar?.styles?.backgroundColor) {
        differences.push({
            type: 'SIDEBAR_BG',
            severity: 'MEDIUM',
            message: 'Cor de fundo da sidebar diferente',
            dashboard: dashboard.sidebar?.styles?.backgroundColor,
            newPatient: newPatient.sidebar?.styles?.backgroundColor
        });
    }

    // Comparar main content background
    if (dashboard.mainContent?.styles?.backgroundColor !== newPatient.mainContent?.styles?.backgroundColor) {
        differences.push({
            type: 'MAIN_BG',
            severity: 'MEDIUM',
            message: 'Cor de fundo do conteúdo diferente',
            dashboard: dashboard.mainContent?.styles?.backgroundColor,
            newPatient: newPatient.mainContent?.styles?.backgroundColor
        });
    }

    // Verificar se topbar existe em ambas
    if (dashboard.topbar?.exists !== newPatient.topbar?.exists) {
        differences.push({
            type: 'TOPBAR_MISSING',
            severity: 'HIGH',
            message: 'Topbar presente em apenas uma página',
            dashboard: dashboard.topbar?.exists,
            newPatient: newPatient.topbar?.exists
        });
    }

    // Comparar número de itens na sidebar
    if (dashboard.sidebarItems?.length !== newPatient.sidebarItems?.length) {
        differences.push({
            type: 'SIDEBAR_ITEMS_COUNT',
            severity: 'LOW',
            message: 'Número diferente de itens na sidebar',
            dashboard: dashboard.sidebarItems?.length,
            newPatient: newPatient.sidebarItems?.length
        });
    }

    return differences;
}

// Imprime resumo dos problemas
function printSummary(results) {
    console.log('\n' + '='.repeat(60));
    console.log('📊 RESUMO DA ANÁLISE\n');

    console.log('Dashboard (v2):');
    console.log(`   CSS: ${results.dashboard.cssVersion}`);
    console.log(`   Sidebar BG: ${results.dashboard.sidebar?.styles?.backgroundColor || 'N/A'}`);
    console.log(`   Main BG: ${results.dashboard.mainContent?.styles?.backgroundColor || 'N/A'}`);

    console.log('\nNovo Paciente:');
    console.log(`   CSS: ${results.newPatient.cssVersion}`);
    console.log(`   Sidebar BG: ${results.newPatient.sidebar?.styles?.backgroundColor || 'N/A'}`);
    console.log(`   Main BG: ${results.newPatient.mainContent?.styles?.backgroundColor || 'N/A'}`);

    console.log('\nDiferenças Encontradas:');
    if (results.comparison.length === 0) {
        console.log('   ✅ Nenhuma diferença significativa');
    } else {
        results.comparison.forEach(diff => {
            const icon = diff.severity === 'HIGH' ? '🔴' : (diff.severity === 'MEDIUM' ? '🟡' : '🟢');
            console.log(`   ${icon} [${diff.type}] ${diff.message}`);
            console.log(`      Dashboard: ${JSON.stringify(diff.dashboard)}`);
            console.log(`      Novo Paciente: ${JSON.stringify(diff.newPatient)}`);
        });
    }

    console.log('\n' + '='.repeat(60));
}

// Executar teste
testNewPatientViaButton();
