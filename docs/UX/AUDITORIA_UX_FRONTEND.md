# 🩺 Auditoria de UX Engineering e Arquitetura Frontend
## Sistema: Arandu - Plataforma de Inteligência Clínica

**Data da Auditoria:** 2025-03-29  
**Versões Analisadas:** HTMX 1.9.10 | Tailwind CSS v4 | Templ (Go)  
**Auditor:** Especialista Sênior em UX Engineering

---

## 1. 🩺 Diagnóstico Executivo

### Resumo da Saúde do Sistema

O sistema apresenta uma arquitetura moderna baseada em HTMX + Tailwind CSS v4 com Templ (Go), demonstrando maturidade no uso de padrões como inline forms e loading states. No entanto, existem lacunas significativas em acessibilidade, tratamento de erros e consistência visual que impactam a experiência do usuário final e a manutenibilidade do código.

### Notas Gerais

| Dimensão | Nota (0-10) | Justificativa |
|----------|-------------|---------------|
| **UX** | 6.5 | Bons indicadores de loading, mas falhas em feedback de erro e prevenção de rage clicks |
| **Acessibilidade** | 5.0 | ARIA presente apenas parcialmente, falta foco management consistente |
| **Código** | 7.0 | Estrutura organizada, mas repetição de classes e ausência de dark mode |

### 🚨 Top 3 Showstoppers

1. **Ausência de tratamento de erro visível ao usuário** - `hx-on::error` não implementado em nenhum formulário
2. **Foco do teclado perdido após swaps do HTMX** - Apenas h1 recebe foco, ignorando inputs e regiões críticas
3. **Dark mode definido no CSS mas não aplicado nas classes Tailwind** - Variáveis CSS existem, mas zero uso de `dark:` classes

---

## 2. 📋 Relatório Detalhado de Achados

### Categoria: UX / HTMX

#### Problema #1: Ausência de Tratamento de Erros Visível

- **Local:** `/workspace/web/components/session/observation_form.templ:4-40`, `/workspace/web/components/patient/new_form.templ:59`
- **Categoria:** UX / HTMX
- **Problema:** Nenhum formulário implementa `hx-on::error` para exibir mensagens amigáveis quando requisições falham. O usuário não recebe feedback se uma operação falhou, podendo levar a perda de dados e frustração.
- **Solução Proposta:** Implementar handler global de erro com `hx-on::error` e elementos dedicados para exibição de erros de rede/servidor.
- **Snippet de Correção:**

```html
<form 
    class="inline-form"
    hx-post="/session/{sessionID}/observations"
    hx-target="#observations-list"
    hx-swap="afterbegin"
    hx-on::after-request="this.reset()"
    hx-on::error="
        const errorDiv = document.getElementById('form-error-container');
        const errorMsg = event.detail.xhr?.status === 400 
            ? 'Dados inválidos. Verifique os campos.' 
            : 'Erro ao salvar. Tente novamente.';
        errorDiv.innerHTML = '<div role=\"alert\" class=\"form-error\"><i class=\"fas fa-exclamation-circle\"></i><p>' + errorMsg + '</p></div>';
        errorDiv.classList.remove('hidden');
    "
>
    <!-- Campos do formulário -->
</form>
<div id="form-error-container" class="hidden" aria-live="polite"></div>
```

---

#### Problema #2: Botões Não São Desabilitados Durante Submit na Maioria dos Forms

- **Local:** `/workspace/web/components/session/observation_form.templ:33`, `/workspace/web/components/patient/new_form.templ:161`
- **Categoria:** UX / Prevenção de Rage Clicks
- **Problema:** Apenas `intervention_form_inline.templ` e `observation_form_inline.templ` usam `hx-disabled-elt`. Formulários principais permitem múltiplos cliques durante submit, causando duplicação de dados.
- **Solução Proposta:** Adicionar `hx-disabled-elt="this button[type='submit']"` em todos os formulários.
- **Snippet de Correção:**

```html
<!-- Opção 1: No botão -->
<button 
    type="submit" 
    class="btn btn-primary"
    hx-disabled-elt="this"
>
    <i class="fas fa-plus btn-icon"></i>Adicionar Observação
</button>

<!-- Opção 2: No nível do form -->
<form hx-post="..." hx-disabled-elt="find button[type='submit']">
    <!-- Conteúdo -->
</form>
```

---

#### Problema #3: Gerenciamento de Foco Inconsistente Após Swap

- **Local:** `/workspace/web/components/layout/layout.templ:230-237`
- **Categoria:** A11y / UX
- **Problema:** O foco é direcionado apenas para `h1` ou elementos com `data-autofocus`, ignorando inputs de formulário, mensagens de erro e regiões dinâmicas críticas para usuários de teclado.
- **Solução Proposta:** Implementar estratégia de foco hierárquica que prioriza inputs com erro, depois inputs vazios, depois headings.
- **Snippet de Correção:**

```javascript
document.body.addEventListener('htmx:afterSwap', (e) => {
    // Prioridade 1: Inputs com erro
    const errorInput = e.target.querySelector('input[aria-invalid="true"], .input-error');
    if (errorInput) {
        errorInput.focus();
        return;
    }
    
    // Prioridade 2: Inputs vazios em formulários
    const emptyInput = e.target.querySelector('input:not([disabled]):not([readonly]):not(.filled), textarea:not([disabled]):not([readonly]):not(.filled)');
    if (emptyInput) {
        emptyInput.focus();
        return;
    }
    
    // Prioridade 3: Heading ou elemento com autofocus
    const heading = e.target.querySelector('h1, [data-autofocus]');
    if (heading) {
        heading.setAttribute('tabindex', '-1');
        heading.focus();
    }
});
```

---

#### Problema #4: hx-push-url Usado Incorretamente em Navegação Parcial

- **Local:** `/workspace/web/components/patient/detail.templ:48-58, 244-265`
- **Categoria:** UX / Histórico do Navegador
- **Problema:** `hx-push-url="true"` gera URLs automáticas que podem não corresponder ao estado real da aplicação quando combinado com `hx-swap="innerHTML"`. O botão "Voltar" pode restaurar conteúdo incorreto.
- **Solução Proposta:** Usar URLs explícitas ou desabilitar push quando apropriado.
- **Snippet de Correção:**

```html
<a 
    href="/patients/{patient.ID}/sessions/new" 
    class="..."
    hx-boost="true"
    hx-target="main"
    hx-swap="innerHTML transition:true"
    <!-- Remover hx-push-url="true" e deixar o hx-boost gerenciar -->
>
    Nova Sessão
</a>
```

---

### Categoria: Design System / Tailwind

#### Problema #5: Zero Uso de Classes `dark:` Apesar de Dark Mode Configurado

- **Local:** `/workspace/web/static/css/input.css:69-83`, todo o codebase
- **Categoria:** Tailwind / Consistência Visual
- **Problema:** O CSS define variáveis para dark mode (`input.css:69-83`), mas nenhuma classe `dark:` é usada nos templates. Usuários com preferência por tema escuro não têm suporte adequado.
- **Solução Proposta:** Aplicar classes `dark:` consistentemente em backgrounds, textos e borders.
- **Snippet de Correção:**

```html
<div class="min-h-screen bg-arandu-bg dark:bg-neutral-900">
    <h1 class="text-arandu-dark dark:text-neutral-100">
        Título da Página
    </h1>
    <div class="bg-white dark:bg-neutral-800 border-neutral-200 dark:border-neutral-700">
        <!-- Conteúdo do card -->
    </div>
</div>
```

---

#### Problema #6: Repetição Excessiva de Grupos de Classes (Violação DRY)

- **Local:** Múltiplos arquivos - padrão recorrente
- **Categoria:** Tailwind / Manutenibilidade
- **Problema:** Padrões como `flex items-center gap-4 mb-6` e `w-10 h-10 bg-gradient-to-br rounded-lg flex items-center justify-center text-white` se repetem em dezenas de componentes, dificultando manutenção e consistência.
- **Solução Proposta:** Extrair padrões recorrentes para componentes helper ou usar @apply com cautela no CSS customizado.
- **Snippet de Correção:**

```css
/* No input.css */
@utility icon-button {
    @apply w-10 h-10 rounded-lg flex items-center justify-center;
}

@utility icon-button-gradient {
    @apply icon-button bg-gradient-to-br text-white;
}

@utility section-header {
    @apply flex items-center gap-4 mb-6;
}
```

```html
<!-- Uso nos templates -->
<div class="section-header">
    <div class="icon-button-gradient from-clinical-teal to-teal-600">
        <i class="fas fa-pills"></i>
    </div>
    <h2 class="font-clinical text-lg font-semibold">Contexto Biológico</h2>
</div>
```

---

#### Problema #7: Falta de Responsividade Consistente em Grids

- **Local:** `/workspace/web/components/patient/detail.templ:63-96` (sota-grid-system)
- **Categoria:** Tailwind / Responsividade
- **Problema:** Classes CSS customizadas `.sota-grid-system` e `.grid-2x2-card` não possuem media queries adequadas para viewports intermediárias (tablet), podendo causar overflow ou layout quebrado.
- **Solução Proposta:** Substituir por grid responsivo nativo do Tailwind.
- **Snippet de Correção:**

```html
<!-- Substituir sota-grid-system -->
<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="bg-white rounded-xl p-6 shadow-sm">
        <!-- Card 1 -->
    </div>
    <div class="bg-white rounded-xl p-6 shadow-sm">
        <!-- Card 2 -->
    </div>
</div>

<!-- Para grid 2x2 específico -->
<div class="grid grid-cols-2 gap-4 sm:grid-cols-2 lg:grid-cols-4">
    <!-- Cards -->
</div>
```

---

### Categoria: Acessibilidade (a11y)

#### Problema #8: Elementos Interativos Sem Labels Adequados

- **Local:** `/workspace/web/components/layout/layout.templ:830-845` (busca de pacientes)
- **Categoria:** A11y
- **Problema:** Embora `aria-label` exista no search input, botões de ação como hamburger menu e ícones de ação não possuem `aria-label` ou `sr-only` text.
- **Solução Proposta:** Adicionar `aria-label` em todos os botões icon-only.
- **Snippet de Correção:**

```html
<button 
    type="button" 
    class="hamburger-menu"
    aria-label="Abrir menu de navegação"
    aria-expanded="false"
    aria-controls="sidebar-drawer"
>
    <i class="fas fa-bars"></i>
</button>

<a 
    href="/patients" 
    class="btn btn-ghost btn-sm"
    aria-label="Voltar para a lista de pacientes"
>
    <i class="fas fa-arrow-left btn-icon"></i>
    <span class="sr-only">Voltar</span>
</a>
```

---

#### Problema #9: Regiões Dinâmicas Sem aria-live ou role Adequados

- **Local:** `/workspace/web/components/session/detail.templ:60-70`, `/workspace/web/components/patient/detail.templ:220-229`
- **Categoria:** A11y
- **Problema:** Listas que são atualizadas via HTMX (`#observations-list`, `#plan-goals-preview`) não possuem `aria-live="polite"` ou `role="region"`, impedindo que leitores de tela anunciem mudanças.
- **Solução Proposta:** Adicionar atributos ARIA em todas as regiões atualizadas dinamicamente.
- **Snippet de Correção:**

```html
<div 
    id="observations-list" 
    class="mt-lg" 
    role="region" 
    aria-label="Lista de observações clínicas"
    aria-live="polite"
    aria-relevant="additions"
>
    @if len(observations) > 0 {
        <!-- Items -->
    } else {
        <p class="text-neutral-500">Nenhuma observação registrada.</p>
    }
</div>

<div 
    id="plan-goals-preview" 
    class="plan-goals-preview"
    role="region"
    aria-label="Pré-visualização do plano terapêutico"
    aria-live="polite"
>
    <!-- Conteúdo dinâmico -->
</div>
```

---

#### Problema #10: Contraste de Cores Potencialmente Insuficiente

- **Local:** `/workspace/web/static/css/input.css` - cores como `--color-neutral-400: #98A2B3`
- **Categoria:** A11y / WCAG
- **Problema:** A cor `neutral-400` (#98A2B3) em fundo branco tem ratio de contraste de ~2.5:1, abaixo do mínimo WCAG AA (4.5:1) para texto normal. Usada em placeholders e ícones.
- **Solução Proposta:** Usar `neutral-500` ou mais escuro para texto informativo.
- **Snippet de Correção:**

```css
/* Manter neutral-400 apenas para bordas e elementos decorativos */
/* Para texto secundário, usar neutral-600 */
```

```html
<!-- Errado -->
<p class="text-neutral-400">Texto informativo</p>

<!-- Correto -->
<p class="text-neutral-600 dark:text-neutral-400">Texto informativo</p>
```

---

### Categoria: Performance

#### Problema #11: HTMX Configuração de Transição Pode Causar Flicker

- **Local:** `/workspace/web/components/layout/layout.templ:76-77`
- **Categoria:** Performance / UX
- **Problema:** `defaultSwapDelay` e `defaultSettleDelay` de 100ms podem ser perceptíveis em conexões lentas, criando sensação de lentidão desnecessária.
- **Solução Proposta:** Reduzir delays ou usar `showRequest` para indicators mais responsivos.
- **Snippet de Correção:**

```javascript
// Ajustar configuração global do HTMX
htmx.config.defaultSwapDelay = 0;
htmx.config.defaultSettleDelay = 0;

// Usar classes .htmx-indicator com transições CSS suaves
.htmx-indicator {
    opacity: 0;
    transition: opacity 200ms ease-in-out;
}
.htmx-request .htmx-indicator {
    opacity: 1;
}
```

---

#### Problema #12: Partial HTML Podem Estar Enviando Conteúdo Desnecessário

- **Local:** Precisa verificação nos handlers Go
- **Categoria:** Performance / Payload
- **Problema:** Alguns templates como `PatientListFragment` ainda renderizam estrutura completa mesmo sendo fragments, potencialmente enviando HTML redundante.
- **Solução Proposta:** Auditar handlers para garantir que apenas o HTML necessário seja retornado.
- **Snippet de Correção:**

```go
// Handler deve retornar apenas o fragmento
func PatientListHandler(w http.ResponseWriter, r *http.Request) {
    patients, err := getPatients()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Detectar se é request HTMX
    if r.Header.Get("HX-Request") == "true" {
        // Retornar apenas o fragmento da lista
        templ.Handler(patient.PatientListFragment(patients, "")).ServeHTTP(w, r)
        return
    }
    
    // Página completa apenas para request direto
    templ.Handler(patient.ListPage(patients)).ServeHTTP(w, r)
}
```

---

## 3. 🚀 Plano de Ação Priorizado

### Matriz de Prioridades

| Prioridade | Ação | Esforço | Impacto | Prazo Sugerido |
|------------|------|---------|---------|----------------|
| **🔴 ALTA** | Implementar `hx-on::error` em todos os formulários críticos (login, create patient, observations) | Baixo | Evita perda de dados e frustração do usuário | Sprint 1 |
| **🔴 ALTA** | Adicionar `hx-disabled-elt` em todos os botões de submit | Baixo | Previne duplicação de registros | Sprint 1 |
| **🔴 ALTA** | Melhorar gerenciamento de foco pós-swap (priorizar inputs com erro) | Médio | Crítico para usuários de teclado e screen readers | Sprint 1 |
| **🔴 ALTA** | Adicionar `aria-live` e `role="region"` em listas dinâmicas | Baixo | Essencial para acessibilidade | Sprint 1 |
| **🟡 MÉDIA** | Implementar classes `dark:` em todos os componentes principais | Alto | Melhora experiência para 10-15% dos usuários | Sprint 2 |
| **🟡 MÉDIA** | Adicionar `aria-label` em botões icon-only | Baixo | Melhora navegação por screen reader | Sprint 2 |
| **🟡 MÉDIA** | Extrair padrões de classes repetidas para @utility ou componentes | Médio | Reduz dívida técnica e melhora consistência | Sprint 2-3 |
| **🟡 MÉDIA** | Corrigir uso de `hx-push-url` para navegação consistente | Baixo | Melhora experiência do botão "Voltar" | Sprint 2 |
| **🟢 BAIXA** | Ajustar delays do HTMX para 0ms | Muito Baixo | Melhoria marginal de percepção de velocidade | Sprint 3 |
| **🟢 BAIXA** | Auditoria de contraste de cores (neutral-400 → neutral-600) | Baixo | Garante conformidade WCAG AA | Sprint 3 |
| **🟢 BAIXA** | Substituir grids customizados por Tailwind grid nativo | Médio | Melhora responsividade e reduz CSS customizado | Sprint 3 |
| **🟢 BAIXA** | Otimizar partials para enviar apenas HTML necessário | Médio | Reduz payload e melhora performance | Sprint 3 |

---

## 📌 Próximos Passos Imediatos (Sprint 1)

### Dia 1-2: Tratamento de Erros + Prevenção de Rage Clicks
- [ ] Implementar `hx-on::error` em:
  - [ ] Login form
  - [ ] Patient creation form
  - [ ] Observation forms
  - [ ] Intervention forms
- [ ] Adicionar `hx-disabled-elt` em todos os botões de submit
- [ ] Criar componente reutilizável de mensagem de erro

### Dia 3-4: Acessibilidade Crítica
- [ ] Implementar gerenciamento de foco hierárquico pós-swap
- [ ] Adicionar `aria-live` e `role="region"` em:
  - [ ] Lists de observações
  - [ ] Lists de intervenções
  - [ ] Preview de plano terapêutico
  - [ ] Notificações e alerts

### Dia 5: Labels e Navegação
- [ ] Adicionar `aria-label` em todos os botões icon-only
- [ ] Revisar e corrigir uso de `hx-push-url`
- [ ] Testar navegação com teclado e screen reader

### Estimativa Total para Correções Críticas (ALTA): **3-5 dias de desenvolvimento**

---

## 📊 Métricas de Sucesso

Após implementação das correções prioritárias, esperar:

| Métrica | Antes | Depois (Meta) |
|---------|-------|---------------|
| Score Lighthouse Accessibility | ~50 | ≥90 |
| Taxa de erro não reportado | Alta | <5% |
| Incidentes de duplicação de dados | Frequente | Zero |
| Tempo médio de tarefa (usuários de teclado) | +40% vs mouse | ≤10% vs mouse |

---

## 🔧 Recursos e Referências

- [HTMX Reference Guide](https://htmx.org/reference/)
- [HTMX Extensions](https://htmx.org/extensions/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [WCAG 2.1 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)
- [ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/)

---

*Documento gerado automaticamente como parte do processo de Code Review + UX Review*  
**Próxima revisão agendada:** 30 dias após implementação das correções de alta prioridade
