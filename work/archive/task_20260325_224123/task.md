## 📋 TASK PARA CORREÇÃO DE LAYOUT

**TASK ID:** task_20260325_224123
**Requirement:** REQ-01-00-02 (Identidade SOTA)  
**Title:** Otimização de Layout — Aproveitamento de Espaço na Main Canvas
**status:** PRONTO_PARA_IMPLEMENTACAO
---

### 🎯 Problema Identificado

A tela de perfil do paciente está com **baixa densidade de conteúdo** e muito espaço em branco desperdiçado, violando o princípio de eficiência espacial do Design System Arandu.

**Problemas específicos:**
1. Cards dispostos de forma não otimizada (muito espaço vertical entre elementos)
2. Grid não aproveita largura total disponível
3. Layout poderia usar sistema de colunas mais eficiente
4. Cards pequenos demais para o espaço disponível

---

### 📐 Solução Proposta

**Layout Otimizado (Desktop):**
```
┌─────────────────────────────────────────────────────────────┐
│ Top Bar                                                      │
├──────────┬──────────────────────────────────────────────────┤
│ Sidebar  │  Perfil: Ana                          [Nova Sessão] │
│          │                                                   │
│ Resumo   │  ┌─────────────────┐  ┌─────────────────────┐  │
│ Anamnese │  │ Identidade Bio  │  │  Notas de Triagem   │  │
│ Prontuár │  │ (65% width)     │  │  (35% width)        │  │
│ Plano    │  │ - Gênero        │  │  - gases            │  │
│          │  │ - Etnia         │  │                     │  │
│          │  │ - Ocupação      │  │                     │  │
│          │  │ - Escolaridade  │  │                     │  │
│          │  └─────────────────┘  └─────────────────────┘  │
│          │                                                   │
│          │  ┌─────────────────┐  ┌─────────────────────┐  │
│          │  │ Linha do Tempo  │  │  Ações Rápidas      │  │
│          │  │ (50% width)     │  │  (50% width)        │  │
│          │  │ - Cadastro      │  │  - Nova Sessão      │  │
│          │  │ - Histórico     │  │  - Completar Anam.  │  │
│          │  │                 │  │  - Definir Metas    │  │
│          │  └─────────────────┘  └─────────────────────┘  │
└──────────┴──────────────────────────────────────────────────┘
```

---

### 🔧 Implementação Técnica

**Arquivos para modificar:**

1. **`web/components/patient/profile.templ`**
   - Adicionar classes de grid responsivas
   - Usar `grid grid-cols-2 gap-6` para cards principais
   - Ajustar tamanhos mínimos dos cards

2. **`web/static/css/style.css`**
   ```css
   .patient-profile-grid {
       display: grid;
       grid-template-columns: repeat(2, 1fr);
       gap: var(--space-lg);
       margin-top: var(--space-xl);
   }
   
   .clinical-card {
       min-height: 200px;
       padding: var(--space-xl);
   }
   
   @media (max-width: 1024px) {
       .patient-profile-grid {
           grid-template-columns: 1fr;
       }
   }
   ```

3. **`web/components/patient/biopsychosocial_panel.templ`**
   - Expandir conteúdo para ocupar espaço disponível
   - Adicionar mais informações visíveis por padrão

---

### ✅ Critérios de Aceitação

- [ ] Main canvas utiliza pelo menos **80% da largura disponível**
- [ ] Cards usam sistema de grid 2 colunas em desktop
- [ ] Em mobile, cards empilham verticalmente (1 coluna)
- [ ] Espaçamento entre cards é consistente (gap-6)
- [ ] Altura mínima dos cards evita espaço vertical excessivo
- [ ] Validação E2E passa sem warnings de densidade
- [ ] Layout responsivo testado em 375px, 768px e 1440px

---

### 🧪 Validação Pós-Implementação

```bash
# 1. Rodar E2E audit
./scripts/arandu_e2e_audit.sh --routes patients

# 2. Verificar warnings de layout
# Não deve aparecer: "Alto índice de containers vazios"

# 3. Testar responsividade
# Acessar em diferentes larguras de tela

# 4. Verificar screenshots
ls -la tmp/audit_screenshots/patient_detail_authenticated.png
```

---

**Quer que eu implemente esta correção agora?** 🌿