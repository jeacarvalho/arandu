# TASK 20260325_225845
**Requirement:** REQ-01-00-02, REQ-02-01-01
**Title:** Correção de Layout SLP — Sidebar e Main Canvas (Overflow e Densidade)
**Status:** PRONTO_PARA_IMPLEMENTACAO

---

## 🎯 Objetivo

Corrigir problemas de layout identificados pelo E2E Audit que estão causando:
1. **Overflow do main canvas sobre a sidebar** — Conteúdo extrapolando e cobrindo menu lateral
2. **Espaço desperdiçado no main canvas** — Baixa densidade de conteúdo, muito espaço em branco
3. **Validações SLP falhando** — Warnings de z-index e densidade de layout

---

## 📊 Problemas Identificados (E2E Audit)

### Problema 1: Sidebar sendo coberta pelo conteúdo
**Sintoma:** Main canvas está se sobrepondo à sidebar em determinadas telas
**Evidência no Audit:**
```
[WARN] dashboard - Sidebar may need explicit z-index to stay above content
```

**O que o script detecta:**
- Elementos com `position: fixed` sem `z-index` explícito
- Sidebar sem z-index maior que main-content
- Main-content sem padding-top adequado para header fixo

---

### Problema 2: Baixa densidade de layout
**Sintoma:** Muito espaço em branco não aproveitado no main canvas
**Evidência no Audit:**
```
[WARN] dashboard - Alto índice de containers vazios (100% dos elementos)
[WARN] dashboard - Considere usar grid/flexbox para melhor aproveitamento
```

**O que o script detecta:**
- Cards empilhados verticalmente em vez de usar grid
- Containers sem conteúdo útil
- Falta de classes grid/flex em layouts com múltiplos elementos

---

### Problema 3: Tipografia clínica ausente
**Sintoma:** Conteúdo clínico sem fonte Source Serif 4
**Evidência no Audit:**
```
[ERROR] dashboard - Clinical typography (.font-clinical) not found
[ERROR] patients_list - Clinical typography (.font-clinical) not found
```

**O que o script valida:**
```bash
grep -q "font-clinical\|font-serif\|Source Serif" "$file"
```

---

## 📚 Documentação de Referência

### Leitura Obrigatória (ANTES de codificar)

| Documento | Seção | Por Que Ler |
|-----------|-------|-------------|
| `docs/architecture/standardized_layout_protocol.md` | Todas | Define anatomia obrigatória (3 zonas) |
| `docs/design-system.md` | Tipografia | Define `.font-clinical` e Source Serif 4 |
| `docs/learnings/MASTER_LEARNINGS.md` | UI/UX e Design System | Anti-padrões e soluções consolidadas |
| `strategy/responsive_sota_strategy.md` | Matriz de Navegação | Comportamento mobile/desktop |
| `docs/requirements/req-01-00-02-editar-paciente.md` | Interface Esperada | Silent Input e tipografia |

---

## 🛠️ Escopo da Tarefa

### Arquivos para Modificar

| Arquivo | Problema | Ação Esperada |
|---------|----------|---------------|
| `web/static/css/style.css` | z-index e overflow | Adicionar z-index explícito para sidebar e main-content |
| `web/components/layout/layout.templ` | Estrutura base | Garantir classes corretas (app-container, sidebar, main-content) |
| `web/components/patient/profile.templ` | Densidade de layout | Usar grid para cards (2 colunas em desktop) |
| `web/components/patient/list.templ` | Tipografia | Aplicar `.font-clinical` a nomes e notas |
| `web/components/dashboard/dashboard.templ` | Tipografia + Densidade | Aplicar `.font-clinical` e grid layout |

---

## ✅ Critérios de Aceitação

### CA-01: Sidebar com z-index explícito
- [ ] Sidebar tem `z-index: 40` (ou classe `z-40`)
- [ ] Main-content tem `z-index: 1` (ou menor que sidebar)
- [ ] Header/Top-bar tem `z-index: 50` (maior que sidebar)
- [ ] E2E audit não mostra warning de z-index

### CA-02: Main canvas com padding adequado
- [ ] Main-content tem `padding-top: 80px` (ou equivalente) para header fixo
- [ ] Conteúdo não fica escondido atrás do header
- [ ] Scroll funciona sem cortar conteúdo

### CA-03: Densidade de layout otimizada
- [ ] Cards em desktop usam `grid grid-cols-2` (quando aplicável)
- [ ] Em mobile, cards empilham verticalmente (1 coluna)
- [ ] E2E audit não mostra warning de "containers vazios"

### CA-04: Tipografia clínica aplicada
- [ ] Todo conteúdo clínico usa `.font-clinical`
- [ ] Nomes de pacientes em listas usam `.font-clinical`
- [ ] Notas e observações usam `.font-clinical`
- [ ] UI administrativa (botões, labels) usa Inter (sans-serif)

### CA-05: Validação E2E
- [ ] `./scripts/arandu_e2e_audit.sh --routes dashboard` passa sem warnings de layout
- [ ] `./scripts/arandu_e2e_audit.sh --routes patients` passa sem warnings de layout
- [ ] `./scripts/arandu_guard.sh` passa sem erros
- [ ] Screenshots são gerados corretamente

### CA-06: Responsividade
- [ ] Em 375px (mobile), sidebar vira drawer
- [ ] Em 768px (tablet), layout se adapta
- [ ] Em 1440px (desktop), grid 2 colunas funciona
- [ ] Testado via DevTools responsive mode

---

## 🔧 Implementação Esperada

### 1. Corrigir z-index no CSS (`web/static/css/style.css`)

```css
/* ADICIONAR/ATUALIZAR */

/* Top Bar — Maior z-index (fica acima de tudo) */
.top-bar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 64px;
    z-index: 50;
    background: #FFFFFF;
    border-bottom: 1px solid #E5E7EB;
}

/* Sidebar — z-index médio (fica abaixo do top-bar, acima do main) */
.sidebar {
    position: fixed;
    top: 64px; /* Altura do top-bar */
    left: 0;
    width: 280px;
    height: calc(100vh - 64px);
    z-index: 40;
    background: #0F6E56;
    overflow-y: auto;
}

/* Main Content — Menor z-index (fica abaixo da sidebar) */
.main-content {
    margin-left: 280px; /* Largura da sidebar */
    padding-top: 80px; /* Altura do top-bar + espaçamento */
    min-height: 100vh;
    z-index: 1;
    background: #E1F5EE; /* Papel de seda */
}

/* Mobile — Sidebar vira drawer */
@media (max-width: 768px) {
    .sidebar {
        transform: translateX(-100%);
        transition: transform 0.3s ease;
    }
    
    .sidebar.open {
        transform: translateX(0);
    }
    
    .main-content {
        margin-left: 0;
    }
}
```

---

### 2. Corrigir Estrutura do Layout (`web/components/layout/layout.templ`)

```templ
templ BaseWithContent(pageTitle string, content templ.Component) {
    <!DOCTYPE html>
    <html>
    <head>
        <title>Arandu — { pageTitle }</title>
        <link href={ templ.URL("/static/css/style.css?v=" + getCSSVersion()) } rel="stylesheet">
    </head>
    <body>
        <!-- CRÍTICO: Classes corretas para validação SLP -->
        <div class="app-container">
            
            <!-- Top Bar (z-index: 50) -->
            <header class="top-bar">
                <!-- Logo, busca, avatar -->
            </header>
            
            <!-- Sidebar (z-index: 40) -->
            <aside class="sidebar">
                <!-- Links de navegação -->
                <!-- CRÍTICO: Deve conter "Anamnese" e "Prontuário" no contexto paciente -->
            </aside>
            
            <!-- Main Canvas (z-index: 1) -->
            <main class="main-content">
                @content
            </main>
            
        </div>
    </body>
    </html>
}
```

---

### 3. Corrigir Densidade no Profile (`web/components/patient/profile.templ`)

```templ
templ PatientProfile(data PatientProfileData) {
    <!-- CRÍTICO: Usar grid para densidade -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        <!-- Card 1: Identidade Biopsicossocial -->
        <div class="clinical-card">
            <h2 class="font-clinical text-xl">Identidade Biopsicossocial</h2>
            <div class="font-clinical">
                <p>Gênero: { data.Gender }</p>
                <p>Etnia: { data.Ethnicity }</p>
                <!-- ... -->
            </div>
        </div>
        
        <!-- Card 2: Notas de Triagem -->
        <div class="clinical-card">
            <h2 class="font-clinical text-xl">Notas de Triagem</h2>
            <div class="font-clinical">
                { data.TriageNotes }
            </div>
        </div>
        
        <!-- Card 3: Linha do Tempo (full width) -->
        <div class="clinical-card md:col-span-2">
            <h2 class="font-clinical text-xl">Linha do Tempo</h2>
            <!-- ... -->
        </div>
        
        <!-- Card 4: Ações Rápidas -->
        <div class="clinical-card md:col-span-2">
            <h2 class="font-clinical text-xl">Ações Rápidas</h2>
            <!-- ... -->
        </div>
        
    </div>
}
```

---

### 4. Aplicar Tipografia Clínica (`web/components/patient/list.templ`)

```templ
templ PatientList(patients []PatientListItem, errorMsg string) {
    <div class="content-header">
        <h1 class="font-clinical text-2xl">Pacientes</h1>
    </div>
    
    <ul class="patient-list">
        for _, patient := range patients {
            <li class="patient-item">
                <!-- CRÍTICO: Nome com font-clinical -->
                <a href={ templ.URL("/patients/" + patient.ID) } 
                   class="font-clinical text-lg hover:text-arandu-primary">
                    { patient.Name }
                </a>
                
                <!-- CRÍTICO: Notas com font-clinical -->
                if patient.Notes != "" {
                    <p class="font-clinical text-base text-gray-600 mt-1">
                        { patient.Notes }
                    </p>
                }
            </li>
        }
    </ul>
}
```

---

## 🧪 Validação Pós-Implementação

Execute esta sequência **obrigatória** antes de concluir:

```bash
# 1. Regenerar templates
~/go/bin/templ generate

# 2. Build
go build ./cmd/arandu

# 3. Guard
./scripts/arandu_guard.sh

# 4. E2E Audit (CRÍTICO)
./scripts/arandu_e2e_audit.sh --routes dashboard,patients

# 5. Verificar warnings de layout
# Esperado: ZERO warnings de z-index, overflow, densidade

# 6. Verificar screenshots
ls -la tmp/audit_screenshots/

# 7. Teste visual manual
# - Abrir dashboard em 1440px
# - Abrir perfil de paciente em 1440px
# - Testar mobile (DevTools → 375px)
# - Verificar que sidebar não é coberta
```

---

## 🚨 Anti-Padrões a Evitar

| Anti-Padrão | Por Que Evitar | Alternativa |
|-------------|----------------|-------------|
| `style="z-index: 40"` inline | Falha no E2E audit | Classes CSS em `style.css` |
| Sidebar sem z-index | Pode ficar atrás do conteúdo | `z-index: 40` explícito |
| Cards empilhados em desktop | Desperdício de espaço | `grid grid-cols-2` |
| Fonte genérica (`serif`) | Não usa Source Serif 4 | `.font-clinical` |
| Main-content sem padding-top | Conteúdo fica atrás do header | `padding-top: 80px` |
| Ignorar mobile | Quebra responsividade | Media queries `@media (max-width: 768px)` |

---

## 📁 Entregáveis

1. [ ] `web/static/css/style.css` — z-index e layout corrigidos
2. [ ] `web/components/layout/layout.templ` — Estrutura SLP correta
3. [ ] `web/components/patient/profile.templ` — Grid 2 colunas
4. [ ] `web/components/patient/list.templ` — Tipografia clínica
5. [ ] `web/components/dashboard/dashboard.templ` — Tipografia + densidade
6. [ ] Output do `arandu_e2e_audit.sh` sem warnings de layout

---

## 📝 Checklist de Conclusão

Antes de executar `arandu_conclude_task.sh`:

- [ ] Li toda documentação de referência listada acima
- [ ] Implementei correções de z-index (sidebar, main, header)
- [ ] Implementei grid layout para densidade (desktop 2 colunas)
- [ ] Apliquei `.font-clinical` a todo conteúdo clínico
- [ ] Rodei `templ generate` e código compilou sem erros
- [ ] Rodei `arandu_guard.sh` e passou
- [ ] Rodei `arandu_e2e_audit.sh` e warnings de layout sumiram
- [ ] Testei responsividade (375px, 768px, 1440px)
- [ ] Verifiquei screenshots gerados em `tmp/audit_screenshots/`
- [ ] Testei regressão em rotas vizinhas (/patients, /sessions, /dashboard)

---

## 💡 Dicas de Implementação

1. **Comece pelo CSS** — Corrija z-index e padding em `style.css` primeiro
2. **Valide com grep** — Não confie apenas no visual:
   ```bash
   grep -n "z-index" web/static/css/style.css
   grep -n "font-clinical" web/components/**/*.templ
   ```
3. **Teste incrementalmente** — Corrija um arquivo por vez e valide
4. **Use DevTools** — Inspecione z-index em tempo real (Chrome DevTools → Layers)
5. **Consulte o Logbook** — `docs/learnings/MASTER_LEARNINGS.md` tem soluções para problemas similares

---

## 🔗 Referências Cruzadas

- **VISION-01** — Registro da prática clínica
- **REQ-01-00-02** — Editar Dados do Paciente (Identidade SOTA)
- **REQ-02-01-01** — Visualizar Histórico do Paciente (Prontuário)
- **SLP** — `docs/architecture/standardized_layout_protocol.md`
- **Design System** — `docs/design-system.md`

---

**Instrução Final:** Esta tarefa é **crítica para a qualidade visual e arquitetural** do Arandu. Não pule validações. Cada warning do E2E audit é uma violação do contrato SLP que deve ser corrigida antes de prosseguir com novas funcionalidades. O layout deve seguir rigorosamente o Standardized Layout Protocol definido na documentação.
```

---

## 🚀 Como Usar Este Prompt

```bash
# 1. Criar a tarefa
./scripts/arandu_create_task.sh "Correção de Layout SLP — Sidebar e Main Canvas"

# 2. Editar o arquivo da tarefa
code work/tasks/task_*/task.md

# 3. Colar o prompt completo acima no arquivo

# 4. Atualizar status
# Mudar de: AGUARDANDO_DETALHES_DO_USUARIO
# Para: PRONTO_PARA_IMPLEMENTACAO

# 5. Processar a tarefa
./scripts/arandu_process_task.sh TASK_ID --execute
```

---
