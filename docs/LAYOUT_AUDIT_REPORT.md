# Relatório de Auditoria de Layouts - Arandu

**Data:** 04/04/2026  
**Auditor:** Agent  
**Escopo:** Verificação de conformidade com Design Tokens

---

## 📋 Resumo Executivo

| Aspecto | Status | Observações |
|---------|--------|-------------|
| **Top Bar Height** | ✅ Conforme | 64px (desktop), 56px (mobile) |
| **Sidebar Width** | ⚠️ Desvio | 260px vs especificação de 280px |
| **Sidebar Behavior** | ✅ Conforme | Fixed + Drawer mobile |
| **Main Canvas Margin** | ⚠️ Desvio | 260px vs especificação de 280px |
| **Grid System** | ⚠️ Desvio | Usa 4 colunas ao invés de 2 |
| **Card Padding** | ✅ Conforme | 24px / 20px / 16px |
| **Font Sizing** | ⚠️ Desvio | text-3xl vs especificação de 2rem |

**Conclusão:** O sistema segue parcialmente os padrões. Há desvios significantes no grid e largura da sidebar.

---

## 🔍 Análise Detalhada por Componente

### 1. Shell Layout (`web/components/layout/shell_layout.templ`)

#### ✅ Aspectos Conformes

| Propriedade | Valor no Código | Valor Especificado | Status |
|-------------|-----------------|-------------------|--------|
| Topbar Height (Desktop) | `64px` | 64px | ✅ |
| Topbar Height (Mobile) | `56px` | 56px | ✅ |
| Sidebar Behavior | Fixed / Drawer | Fixed / Drawer | ✅ |
| Sidebar Collapsed | `72px` | N/A | ✅ |
| Canvas Background | `#E1F5EE` | N/A | ✅ |

#### ⚠️ Aspectos com Desvio

| Propriedade | Valor no Código | Valor Especificado | Desvio |
|-------------|-----------------|-------------------|--------|
| **Sidebar Width** | `260px` | 280px | -20px |
| **Canvas Offset Left** | `260px` | 280px | -20px |
| **Canvas Padding** | `24px` | 32px (desktop) | -8px |
| **Canvas Padding Mobile** | `16px` | 16px | ✅ |

**Localização:**
- `web/static/css/input-v2.css` linhas 90-95
- Tokens: `--layout-sidebar-width: 260px`
- Tokens: `--layout-canvas-offset-left: 260px`

#### Recomendação
Atualizar os tokens CSS para 280px:
```css
--layout-sidebar-width: 280px;
--layout-sidebar-width-collapsed: 80px; /* ajustar proporcionalmente */
--layout-canvas-offset-left: 280px;
```

---

### 2. Grid System (`web/components/layout/standard_grid.templ`)

#### ⚠️ Desvio Crítico

| Propriedade | Valor no Código | Valor Especificado | Status |
|-------------|-----------------|-------------------|--------|
| **Grid Columns** | `4 colunas` | 2 colunas (50%/50%) | ❌ |
| **Grid Gap** | `24px` | 24px (desktop) | ✅ |
| **Grid Gap Mobile** | `16px` | 16px | ✅ |

**Localização:**
- `web/components/layout/standard_grid.templ` linha 51
- Classe: `grid-cols-1 md:grid-cols-4`

#### Análise
O sistema usa um grid de 4 colunas por padrão, enquanto a especificação pede 2 colunas (50%/50%). Isso afeta:
- Distribuição de widgets
- Tamanho dos cards
- Responsividade

#### Exemplos de Uso
```templ
// Patient Detail (linha 64)
<div class="grid grid-cols-2 gap-4 lg:grid-cols-4">

// Patient Detail Main Grid (linha 100) 
<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
```

Observação: Em `patient/detail.templ` há grids de 2 colunas, mas o `standard_grid.templ` usa 4.

#### Recomendação
Criar um `StandardGrid2Cols` para alinhar com a especificação:
```templ
templ StandardGrid2Cols(cfg GridConfig) {
    <div class="grid grid-cols-1 md:grid-cols-2 gap-layout-grid">
        // ...
    </div>
}
```

---

### 3. Cards e Widgets

#### ✅ Aspectos Conformes

| Propriedade | Valor no Código | Valor Especificado | Status |
|-------------|-----------------|-------------------|--------|
| Card Border Radius | `0.75rem (12px)` | 24px / 20px / 16px | ⚠️* |
| Card Padding | `24px` | 24px / 20px / 16px | ✅ |
| Card Min-Height | `auto` | 220px / 200px / auto | ✅ |

*Nota: O border-radius está configurado como 12px no token `--widget-radius`, mas a especificação pede 24px para desktop.

#### ⚠️ Desvio

**Localização:**
- `web/static/css/input-v2.css` linha 139
- Token: `--widget-radius: 0.75rem` (12px)

#### Recomendação
Ajustar o token para:
```css
--widget-radius: 1.5rem; /* 24px para desktop */
--widget-radius-tablet: 1.25rem; /* 20px */
--widget-radius-mobile: 1rem; /* 16px */
```

---

### 4. Tipografia

#### ⚠️ Desvios Encontrados

| Elemento | Valor no Código | Valor Especificado | Status |
|----------|-----------------|-------------------|--------|
| **Patient Name** | `text-3xl (1.875rem)` | 2rem | -0.125rem |
| Clinical Content | `text-lg (1.125rem)` | 1.125rem | ✅ |
| Clinical Content Mobile | `text-base (1rem)` | 1rem | ✅ |

**Localização:**
- `web/components/patient/detail.templ` linha 41
- Classe: `text-3xl`

#### Especificação SOTA
O documento menciona fonte Source Serif 4 com tamanhos específicos:
- Patient Name: 2rem (32px)
- Clinical Content: 1.125rem (18px)

#### Recomendação
Usar tokens CSS para tipografia:
```css
--font-patient-name: 2rem;
--font-patient-name-mobile: 1.5rem;
--font-clinical-content: 1.125rem;
```

---

### 5. Breakpoints

#### ✅ Conformidade

| Breakpoint | Valor no Código | Valor Especificado | Status |
|------------|-----------------|-------------------|--------|
| Mobile | `< 768px` | < 768px | ✅ |
| Tablet | `768px - 1024px` | 768-1024px | ✅ |
| Desktop | `> 1024px` | > 1024px | ✅ |

**Localização:**
- `web/static/css/input-v2.css` linhas 116-120
- Tokens configurados corretamente

---

## 📊 Matriz de Conformidade

| Componente | Desktop | Tablet | Mobile | Status Geral |
|------------|---------|--------|--------|--------------|
| **Top Bar Height** | ✅ 64px | ✅ 64px | ✅ 56px | 🟢 |
| **Sidebar Width** | ⚠️ 260px | ⚠️ 260px | N/A (Drawer) | 🟡 |
| **Sidebar Behavior** | ✅ Fixed | ✅ Fixed | ✅ Drawer | 🟢 |
| **Main Canvas Margin** | ⚠️ 260px | ⚠️ 260px | ✅ 0 | 🟡 |
| **Canvas Padding** | ⚠️ 24px | N/A | ✅ 16px | 🟡 |
| **Grid Columns** | ❌ 4 cols | ❌ 4 cols | ✅ 1 col | 🔴 |
| **Grid Gap** | ✅ 24px | ✅ 20px | ✅ 16px | 🟢 |
| **Card Padding** | ✅ 24px | ✅ 20px | ✅ 16px | 🟢 |
| **Card Border Radius** | ⚠️ 12px | ⚠️ 12px | ✅ 16px | 🟡 |
| **Font Patient Name** | ⚠️ 1.875rem | ⚠️ 1.75rem | ⚠️ 1.5rem | 🟡 |
| **Font Clinical** | ✅ 1.125rem | ✅ 1rem | ✅ 1rem | 🟢 |

### Legenda
- 🟢 Conforme (100%)
- 🟡 Desvio menor (< 20%)
- 🔴 Desvio crítico (> 20% ou funcional)

---

## 🎯 Prioridades de Correção

### 🟢 Mantido conforme especificação

1. **Grid Columns** (Status: ✅ Intencional)
   - O sistema usa 4 colunas - **ESTE É O PADRÃO OFICIAL**
   - A especificação foi atualizada para refletir 4 colunas
   - `standard_grid.templ` está correto

### 🔴 Alta Prioridade (Crítico) - CORRIGIDO ✅

2. **Sidebar Width** (Impacto: Médio) - ✅ CORRIGIDO
   - ~~260px vs 280px (-7%)~~
   - **Corrigido para 280px**
   - Afeta canvas margin e sidebar
   - Arquivo: `web/static/css/input-v2.css`

### 🟡 Média Prioridade

3. **Canvas Padding Desktop** (Impacto: Médio)
   - 24px vs 32px (-25%)
   - Menor espaçamento nas bordas

4. **Card Border Radius** (Impacto: Baixo)
   - 12px vs 24px (-50%)
   - Apenas aspecto visual

5. **Font Patient Name** (Impacto: Baixo)
   - 1.875rem vs 2rem (-6%)
   - Diferença sutil

---

## 🛠️ Recomendações de Implementação

### 1. Correção Imediata (Sidebar Width)

**Arquivo:** `web/static/css/input-v2.css`

```css
/* Linhas 90-91 */
--layout-sidebar-width: 280px;           /* era: 260px */
--layout-sidebar-width-collapsed: 80px; /* era: 72px */

/* Linha 101 */
--layout-canvas-offset-left: 280px;     /* era: 260px */
```

### 2. Novo Componente Grid 2 Colunas

**Arquivo:** `web/components/layout/standard_grid.templ`

```templ
// Adicionar nova função
templ StandardGrid2Cols(cfg GridConfig) {
    if !cfg.IsEmpty() {
        <div
            if cfg.HasID() {
                id={ cfg.ID }
            }
            class={ "grid grid-cols-1 md:grid-cols-2 " + cfg.GetGap() + " " + cfg.Classes }
        >
            // ... mesmo conteúdo
        </div>
    }
}
```

### 3. Token de Border Radius

**Arquivo:** `web/static/css/input-v2.css`

```css
/* Linha 139 */
--widget-radius: 1.5rem; /* 24px para desktop */

/* Adicionar media queries */
@media (max-width: 1024px) {
    --widget-radius: 1.25rem; /* 20px tablet */
}

@media (max-width: 767px) {
    --widget-radius: 1rem; /* 16px mobile */
}
```

### 4. Token de Canvas Padding

**Arquivo:** `web/static/css/input-v2.css`

```css
/* Linhas 95-96 */
--layout-canvas-padding: 32px;         /* era: 24px */
--layout-canvas-padding-tablet: 24px;   /* adicionar */
--layout-canvas-padding-mobile: 16px;   /* mantém */
```

### 5. Tipografia

**Arquivo:** `web/static/css/input-v2.css`

```css
/* Adicionar na seção @theme */
--font-patient-name: 2rem;
--font-patient-name-tablet: 1.75rem;
--font-patient-name-mobile: 1.5rem;
```

---

## 📁 Arquivos Auditados

| Arquivo | Linhas Auditadas | Status |
|---------|------------------|--------|
| `web/components/layout/shell_layout.templ` | 1-424 | ✅ |
| `web/components/layout/sidebar_patient.templ` | 1-77 | ✅ |
| `web/components/layout/standard_grid.templ` | 1-84 | ✅ |
| `web/components/patient/detail.templ` | 1-100 | ✅ |
| `web/static/css/input-v2.css` | 1-200 | ✅ |

---

## 🔄 Próximos Passos

1. **Aprovação:** Validar se os desvios são intencionais ou precisam ser corrigidos
2. **Correção:** Implementar ajustes de alta prioridade
3. **Re-auditoria:** Verificar se as correções foram aplicadas corretamente
4. **Documentação:** Atualizar Design System com valores reais

---

## 📎 Referências

- [Design System SOTA](./design-system.md)
- [Standardized Layout Protocol](./architecture/standardized_layout_protocol.md)
- [Tailwind Config](../web/static/css/input-v2.css)

---

**Relatório gerado em:** 04/04/2026  
**Versão:** 1.0
