# 📋 TASK PROMPT — Correção de Validações SLP e Design System

```markdown
# TASK 20260325_XXXXXX

**Requirement:** REQ-01-00-02, REQ-02-01-01, CAP-01-06
**Title:** Correção de Validações SLP — Sidebar, Tipografia Clínica e Inline Styles
**Status:** PRONTO_PARA_IMPLEMENTACAO

---

## 🎯 Objetivo

Corrigir 3 problemas críticos identificados pelo E2E Audit que estão causando falhas nas validações de **Standardized Layout Protocol (SLP)** e **Design System**:

1. **Sidebar sem links SLP** — Dashboard e páginas de paciente não exibem "Anamnese" e "Prontuário"
2. **Tipografia clínica ausente** — Conteúdo clínico não usa `.font-clinical` (Source Serif 4)
3. **Inline styles** — observation_edit e intervention_edit possuem `style=""` inline

---

## 📊 Problemas Identificados (E2E Audit)

### Problema 1: SLP Sidebar Missing Expected Links
```
[ERROR] dashboard - SLP sidebar missing expected links
[ERROR] patients_list - SLP sidebar missing expected links
```

**O que o script valida:**
```bash
grep -q "Anamnese" "$file" && grep -q "Prontuário" "$file"
```

**O que está faltando:**
```html
<!-- Sidebar contextual do paciente DEVE conter -->
<a href="/patients/{id}/anamnesis">Anamnese</a>
<a href="/patients/{id}/history">Prontuário</a>
```

---

### Problema 2: Clinical Typography Not Found
```
[ERROR] dashboard - Clinical typography (.font-clinical) not found
```

**O que o script valida:**
```bash
grep -q "font-clinical\|font-serif\|Source Serif" "$file"
```

**O que está faltando:**
```html
<!-- Conteúdo clínico DEVE usar -->
<div class="font-clinical">
    { patient.Notes }
</div>
```

---

### Problema 3: Inline Styles Detected
```
[ERROR] observation_edit - Inline styles: 1
[ERROR] intervention_edit - Inline styles: 1
```

**O que o script valida:**
```bash
grep -o 'style="' "$file" | wc -l
```

**O que precisa ser corrigido:**
```templ
<!-- ERRADO -->
<textarea style="margin-bottom: 1rem; font-family: serif;">

<!-- CORRETO -->
<textarea class="clinical-textarea mb-4">
```

---

## 📚 Documentação de Referência

### Leitura Obrigatória (ANTES de codificar)

| Documento | Seção | Por Que Ler |
|-----------|-------|-------------|
| `docs/architecture/standardized_layout_protocol.md` | Todas | Define anatomia obrigatória da sidebar |
| `docs/design-system.md` | Tipografia | Define `.font-clinical` e Source Serif 4 |
| `docs/learnings/MASTER_LEARNINGS.md` | UI/UX e Design System | Anti-padrões e soluções consolidadas |
| `docs/learnings/TEMPL_GUIDE.md` | CSS no Componente | Como usar CSS scoped em .templ |
| `docs/requirements/req-01-00-02-editar-paciente.md` | Interface Esperada | Silent Input e tipografia |
| `docs/requirements/req-02-01-01-visualizar-historico.md` | Interface | Tipografia clínica no prontuário |

---

## 🛠️ Escopo da Tarefa

### Arquivos para Modificar

| Arquivo | Problema | Ação Esperada |
|---------|----------|---------------|
| `web/components/layout/sidebar.templ` | Sidebar sem links SLP | Adicionar links "Anamnese" e "Prontuário" no contexto paciente |
| `web/components/dashboard/dashboard.templ` | Sem tipografia clínica | Aplicar `.font-clinical` ao conteúdo |
| `web/components/patient/detail.templ` | Sem tipografia clínica | Aplicar `.font-clinical` a notes e conteúdo |
| `web/components/session/observation_edit_form.templ` | Inline style | Mover para classe CSS |
| `web/components/session/intervention_edit_form.templ` | Inline style | Mover para classe CSS |
| `web/static/css/style.css` | Classes faltando | Adicionar `.clinical-textarea`, `.font-clinical` se necessário |

---

## ✅ Critérios de Aceitação

### CA-01: Sidebar SLP
- [ ] Sidebar no contexto paciente contém link "Anamnese"
- [ ] Sidebar no contexto paciente contém link "Prontuário"
- [ ] Links usam rotas corretas (`/patients/{id}/anamnesis`, `/patients/{id}/history`)

### CA-02: Tipografia Clínica
- [ ] Todo conteúdo clínico usa classe `.font-clinical`
- [ ] Fonte Source Serif 4 é aplicada via CSS
- [ ] UI administrativa (botões, labels) usa Inter (sans-serif)

### CA-03: Zero Inline Styles
- [ ] `grep -o 'style="' web/components/**/*.templ | wc -l` retorna 0
- [ ] Todos estilos movidos para `web/static/css/style.css`
- [ ] Classes CSS são semânticas (ex: `.clinical-textarea`, não `.m-16`)

### CA-04: Validação E2E
- [ ] `./scripts/arandu_e2e_audit.sh --routes dashboard` passa sem erros SLP
- [ ] `./scripts/arandu_e2e_audit.sh --routes patients` passa sem erros SLP
- [ ] `./scripts/arandu_e2e_audit.sh --routes observations` não detecta inline styles
- [ ] `./scripts/arandu_e2e_audit.sh --routes interventions` não detecta inline styles

### CA-05: Regressão
- [ ] `./scripts/arandu_guard.sh` passa sem erros
- [ ] `templ generate` executa sem warnings
- [ ] `go build ./cmd/arandu` compila sem erros

---

## 🔧 Implementação Esperada

### 1. Corrigir Sidebar (web/components/layout/sidebar.templ)

```templ
// ADICIONAR no contexto do paciente
if patientContext {
    <nav class="contextual-sidebar">
        <a href={ templ.URL("/patients/" + patientID + "/anamnesis") } 
           class="sidebar-link">
            <span class="icon">📋</span>
            <span>Anamnese</span>
        </a>
        <a href={ templ.URL("/patients/" + patientID + "/history") } 
           class="sidebar-link">
            <span class="icon">📖</span>
            <span>Prontuário</span>
        </a>
    </nav>
}
```

---

### 2. Corrigir Tipografia (web/static/css/style.css)

```css
/* ADICIONAR se não existir */
.font-clinical {
    font-family: 'Source Serif 4', serif;
    font-size: 1.125rem;
    line-height: 1.75;
    color: #1F2937;
}

.clinical-textarea {
    font-family: 'Source Serif 4', serif;
    font-size: 1.125rem;
    line-height: 1.75;
    margin-bottom: 1rem;
    border: none;
    border-bottom: 1px solid #E5E7EB;
    background: #F7F8FA;
    padding: 0.5rem;
    width: 100%;
}

.clinical-textarea:focus {
    outline: none;
    border-bottom-color: #0F6E56;
    background: #FFFFFF;
}
```

---

### 3. Corrigir Observation Edit Form (web/components/session/observation_edit_form.templ)

```templ
// ANTES (ERRADO)
<textarea style="margin-bottom: 1rem; font-family: serif;" 
          name="content">{ data.Content }</textarea>

// DEPOIS (CORRETO)
<textarea class="clinical-textarea" 
          name="content">{ data.Content }</textarea>
```

---

### 4. Corrigir Intervention Edit Form (web/components/session/intervention_edit_form.templ)

```templ
// ANTES (ERRADO)
<textarea style="margin-bottom: 1rem; font-family: serif;" 
          name="content">{ data.Content }</textarea>

// DEPOIS (CORRETO)
<textarea class="clinical-textarea" 
          name="content">{ data.Content }</textarea>
```

---

### 5. Aplicar Tipografia Clínica (dashboard.templ, patient/detail.templ)

```templ
// ANTES (SEM CLASSE)
<div class="patient-notes">
    { patient.Notes }
</div>

// DEPOIS (COM CLASSE)
<div class="patient-notes font-clinical">
    { patient.Notes }
</div>
```

---

## 🧪 Validação Pós-Implementação

Execute esta sequência **obrigatória** antes de concluir a tarefa:

```bash
# 1. Gerar templates
~/go/bin/templ generate

# 2. Validar inline styles (deve retornar 0)
grep -o 'style="' web/components/**/*.templ | wc -l

# 3. Validar SLP no HTML gerado
grep -q "Anamnese" web/components/layout/sidebar.templ && echo "✅ Sidebar OK"
grep -q "Prontuário" web/components/layout/sidebar.templ && echo "✅ Prontuário OK"

# 4. Validar tipografia
grep -q "font-clinical" web/components/**/*.templ && echo "✅ Tipografia OK"

# 5. Build
go build ./cmd/arandu

# 6. Guard
./scripts/arandu_guard.sh

# 7. E2E Audit (CRÍTICO)
./scripts/arandu_e2e_audit.sh --routes dashboard,patients,observations,interventions

# 8. Verificar relatório
# Esperado: 0 erros SLP, 0 inline styles
```

---

## 🚨 Anti-Padrões a Evitar

| Anti-Padrão | Por Que Evitar | Alternativa |
|-------------|----------------|-------------|
| `style="..."` em .templ | Falha no E2E audit | Classes CSS em `style.css` |
| Fonte genérica (`serif`) | Não usa Source Serif 4 | `.font-clinical` |
| Sidebar estática | Não é sensível ao contexto | Condicional `if patientContext` |
| Links hardcoded | Quebra convenções de rotas | `templ.URL()` |
| Ignorar `templ generate` | Código desatualizado | Sempre rodar após editar .templ |

---

## 📁 Entregáveis

1. [ ] `web/components/layout/sidebar.templ` — Links SLP adicionados
2. [ ] `web/static/css/style.css` — Classes `.font-clinical` e `.clinical-textarea`
3. [ ] `web/components/session/observation_edit_form.templ` — Zero inline styles
4. [ ] `web/components/session/intervention_edit_form.templ` — Zero inline styles
5. [ ] `web/components/dashboard/dashboard.templ` — Tipografia clínica aplicada
6. [ ] `web/components/patient/detail.templ` — Tipografia clínica aplicada
7. [ ] Output do `arandu_e2e_audit.sh` mostrando 0 erros SLP e 0 inline styles

---

## 📝 Checklist de Conclusão

Antes de executar `arandu_conclude_task.sh`:

- [ ] Li toda documentação de referência listada acima
- [ ] Implementei todas correções de sidebar, tipografia e inline styles
- [ ] Rodei `templ generate` e código compilou sem erros
- [ ] Rodei `arandu_guard.sh` e passou
- [ ] Rodei `arandu_e2e_audit.sh` e validações SLP passaram
- [ ] Verifiquei que não há inline styles (`grep` retorna 0)
- [ ] Testei regressão em rotas vizinhas (/patients, /sessions)
- [ ] Documentei aprendizados valiosos (se houver)

---

## 💡 Dicas de Implementação

1. **Comece pelo CSS** — Crie as classes em `style.css` antes de modificar os templates
2. **Use grep para validar** — Não confie apenas no visual, valide com os mesmos comandos do audit
3. **Teste incrementalmente** — Corrija um arquivo por vez e valide antes de prosseguir
4. **Consulte o Logbook** — `docs/learnings/MASTER_LEARNINGS.md` tem soluções para problemas similares

---

## 🔗 Referências Cruzadas

- **VISION-01** — Registro da prática clínica
- **CAP-01-06** — Avaliação Inicial e Anamnese
- **REQ-01-00-02** — Editar Dados do Paciente (Identidade SOTA)
- **REQ-02-01-01** — Visualizar Histórico do Paciente (Prontuário)
- **SLP** — `docs/architecture/standardized_layout_protocol.md`

---

**Instrução Final:** Esta tarefa é **crítica para a qualidade visual e arquitetural** do Arandu. Não pule validações. Cada erro SLP detectado pelo audit é uma violação do contrato de Design System que deve ser corrigida antes de prosseguir com novas funcionalidades.
```

---

## 📋 Como Usar Este Prompt

```bash
# 1. Criar a tarefa
./scripts/arandu_create_task.sh "Correção de Validações SLP — Sidebar, Tipografia e Inline Styles"

# 2. Editar o arquivo da tarefa
code work/tasks/task_*/task.md

# 3. Colar o prompt completo acima no arquivo

# 4. Atualizar status
# Mudar de: AGUARDANDO_DETALHES_DO_USUARIO
# Para: PRONTO_PARA_IMPLEMENTACAO
---