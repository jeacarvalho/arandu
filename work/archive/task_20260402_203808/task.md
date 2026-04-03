
## 🛠️ Instruções Específicas por Seção

### 1. **Header do Paciente — Compactar**

**Atual:**
```templ
<h1 class="text-2xl font-clinical">Carolina Costa</h1>
<p class="text-sm">Ativa · 2a3m em terapia</p>
<p class="text-sm">Masculino · Branca · Ensino Superior</p>
```

**Target:**
```templ
<div class="flex items-center gap-3">
    <!-- Avatar -->
    <div class="w-12 h-12 rounded-full bg-primary-100 text-primary-600 flex items-center justify-center font-bold">CC</div>
    
    <!-- Nome + Info -->
    <div>
        <h1 class="text-xl font-clinical font-semibold">Carolina Costa</h1>
        <div class="flex items-center gap-2 mt-1">
            <span class="badge badge-success">Em tratamento</span>
            <span class="badge badge-neutral">22 anos</span>
            <span class="text-xs text-muted">Feminino · Branca · Estudante</span>
        </div>
    </div>
    
    <!-- Stats Cards -->
    <div class="flex gap-2 ml-auto">
        <div class="stat-compact">
            <span class="text-2xl font-bold">124</span>
            <span class="text-xs text-muted">sessões</span>
        </div>
        <div class="stat-compact">
            <span class="text-2xl font-bold">2,4a</span>
            <span class="text-xs text-muted">em terapia</span>
        </div>
    </div>
</div>
```

**CSS Adicional:**
```css
/* Badge System */
.badge {
    display: inline-flex;
    align-items: center;
    padding: 4px 10px;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
    font-family: var(--font-sans);
}

.badge-success {
    background: rgba(16, 185, 129, 0.1);
    color: #10b981;
}

.badge-neutral {
    background: rgba(107, 114, 128, 0.1);
    color: #6b7280;
}

/* Stats Compactos */
.stat-compact {
    background: var(--arandu-paper);
    border: 1px solid var(--neutral-200);
    border-radius: var(--radius-lg);
    padding: var(--space-md);
    text-align: center;
    min-width: 80px;
}
```

---

### 2. **Cards Superiores — Consolidar**

**Atual:** 3 cards separados (Notas, Ações, Sessões)

**Target:** 2 cards superiores + 1 lista inferior

**Estrutura:**
```templ
<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
    <!-- Card Esquerdo (2/3 largura) -->
    <div class="md:col-span-2 card-compact">
        <h3 class="section-title-compact">Notas de triagem</h3>
        <p class="text-sm text-clinical mb-3">Paciente de 22 anos com TAG...</p>
        
        <h4 class="section-title-compact mt-4">Observações recentes</h4>
        <ul class="observation-list-compact">
            <li class="observation-item">
                <p class="text-sm text-clinical">Cognições mais flexíveis...</p>
                <span class="text-xs text-muted">02/04/2026</span>
            </li>
        </ul>
    </div>
    
    <!-- Card Direito (1/3 largura) -->
    <div class="card-compact">
        <h3 class="section-title-compact">Acesso rápido</h3>
        <div class="quick-actions-compact">
            <a href="#" class="action-btn-compact">Anamnese clínica</a>
            <a href="#" class="action-btn-compact">Plano terapêutico</a>
            <a href="#" class="action-btn-compact">Nova sessão</a>
        </div>
    </div>
</div>
```

**CSS Adicional:**
```css
/* Card Compacto */
.card-compact {
    background: var(--arandu-paper);
    border: 1px solid var(--neutral-200);
    border-radius: var(--radius-lg);
    padding: var(--space-md);
}

/* Títulos de Seção Compactos */
.section-title-compact {
    font-family: var(--font-sans);
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--arandu-dark);
    margin-bottom: var(--space-sm);
    display: flex;
    align-items: center;
    gap: var(--space-xs);
}

/* Lista de Observações Compacta */
.observation-list-compact {
    list-style: none;
    padding: 0;
    margin: 0;
}

.observation-item {
    padding: var(--space-sm) 0;
    border-bottom: 1px solid var(--neutral-100);
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
}

.observation-item:last-child {
    border-bottom: none;
}

/* Botões de Ação Compactos */
.action-btn-compact {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
    padding: var(--space-sm) var(--space-md);
    background: var(--neutral-50);
    border-radius: var(--radius-md);
    text-decoration: none;
    color: var(--arandu-dark);
    font-family: var(--font-sans);
    font-size: 0.875rem;
    margin-bottom: var(--space-xs);
    transition: all var(--transition-fast);
}

.action-btn-compact:hover {
    background: var(--primary-50);
    color: var(--arandu-primary);
}
```

---

### 3. **Lista de Sessões — Transformar Cards em Lista**

**Atual:** Cards individuais para cada sessão

**Target:** Lista compacta com badges

```templ
<div class="card-compact">
    <div class="flex justify-between items-center mb-3">
        <h3 class="section-title-compact">Sessões recentes</h3>
        <a href="#" class="text-xs text-primary-600">Ver prontuário completo →</a>
    </div>
    
    <ul class="session-list-compact">
        <li class="session-item-compact">
            <!-- Data Badge -->
            <div class="session-date-badge">
                <span class="text-xs font-semibold">18/05/2025</span>
            </div>
            
            <!-- Session Badge -->
            <div class="session-number-badge">
                <span class="text-xs">Sessão 124</span>
            </div>
            
            <!-- Content -->
            <p class="session-content-compact">
                Paciente em fase de melhora do tratamento...
            </p>
            
            <!-- Arrow -->
            <i class="fas fa-chevron-right text-xs text-muted"></i>
        </li>
    </ul>
</div>
```

**CSS Adicional:**
```css
/* Lista de Sessões Compacta */
.session-list-compact {
    list-style: none;
    padding: 0;
    margin: 0;
}

.session-item-compact {
    display: flex;
    align-items: center;
    gap: var(--space-md);
    padding: var(--space-sm) var(--space-md);
    border-bottom: 1px solid var(--neutral-100);
    transition: background-color var(--transition-fast);
}

.session-item-compact:hover {
    background: var(--neutral-50);
}

.session-item-compact:last-child {
    border-bottom: none;
}

/* Data Badge */
.session-date-badge {
    min-width: 70px;
}

.session-date-badge span {
    font-family: var(--font-sans);
    font-size: 0.75rem;
    color: var(--neutral-600);
}

/* Session Number Badge */
.session-number-badge {
    background: var(--primary-50);
    color: var(--primary-600);
    padding: 4px 10px;
    border-radius: var(--radius-md);
    min-width: 80px;
    text-align: center;
}

.session-number-badge span {
    font-family: var(--font-sans);
    font-size: 0.75rem;
    font-weight: 600;
}

/* Session Content */
.session-content-compact {
    flex: 1;
    font-family: var(--font-clinical);
    font-size: 0.875rem;
    color: var(--neutral-700);
    line-height: 1.5;
}
```

---

### 4. **Redução Geral de Spacing**

**Adicionar ao CSS:**
```css
/* Override para layout compacto */
.patient-profile-compact {
    --space-md: 12px;
    --space-lg: 16px;
    --space-xl: 20px;
}

.patient-profile-compact .card-compact {
    padding: var(--space-md);
}

.patient-profile-compact .section-title-compact {
    margin-bottom: var(--space-sm);
}

.patient-profile-compact .observation-item,
.patient-profile-compact .session-item-compact {
    padding: var(--space-sm) var(--space-md);
}
```

---

### 5. **Tipografia — Reduzir Sizes**

**Tabela de Conversão:**

| Elemento | Atual | Target |
|----------|-------|--------|
| Nome do Paciente | `text-2xl` (24px) | `text-xl` (20px) |
| Títulos de Seção | `text-lg` (18px) | `text-sm` (14px) |
| Corpo de Texto | `text-base` (16px) | `text-sm` (14px) |
| Meta/Labels | `text-xs` (12px) | `text-xs` (12px) ✓ |
| Stats Numbers | `text-3xl` (30px) | `text-2xl` (24px) |

---

## 📋 Checklist de Implementação

```markdown
## Header do Paciente
- [ ] Avatar circular com iniciais
- [ ] Nome + Tags inline (status, idade, info)
- [ ] Stats cards (sessões, tempo terapia) na mesma linha
- [ ] Reduzir font-size do nome (2xl → xl)

## Cards Superiores
- [ ] Consolidar Notas + Observações em 1 card (2/3 largura)
- [ ] Ações Rápidas em card separado (1/3 largura)
- [ ] Reduzir padding dos cards (xl → md)
- [ ] Reduzir gap entre cards (6 → 4)

## Lista de Sessões
- [ ] Transformar de cards para lista
- [ ] Adicionar date badge à esquerda
- [ ] Adicionar session number badge
- [ ] Conteúdo inline (não separado)
- [ ] Reduzir altura de cada item

## Tipografia
- [ ] Reduzir todos font-sizes em 1 nível
- [ ] Manter font-clinical apenas em conteúdo narrativo
- [ ] Usar font-sans para labels, badges, datas

## Spacing
- [ ] Reduzir padding geral (24px → 16px)
- [ ] Reduzir gap entre elementos (6 → 4)
- [ ] Manter consistência em todo layout

## Badges
- [ ] Criar sistema de badges (status, sessão, idade)
- [ ] Usar cores semânticas (success, neutral, primary)
- [ ] Border-radius full (pill shape)
```

---

## 🎨 Cores do Arandu (Manter)

```css
/* Manter paleta botânica */
--arandu-primary: #0F6E56;
--arandu-active: #1D9E75;
--arandu-bg: #E1F5EE;
--arandu-paper: #FFFFFF;
--arandu-dark: #085041;

/* Usar para badges */
.badge-success { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.badge-primary { background: rgba(15, 110, 86, 0.1); color: #0F6E56; }
.badge-neutral { background: rgba(107, 114, 128, 0.1); color: #6b7280; }
```
