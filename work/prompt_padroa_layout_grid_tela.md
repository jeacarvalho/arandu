# 📐 Prompt Mestre de Implementação v2.0
**(Template Genérico para Todas as Features)**

**Copie, preencha TODOS os campos `[ ]` e envie:**

---

# Implementação: [NOME DA FUNCIONALIDADE]

## 📝 O que construir
[Descreva em 1-2 frases o que deve ser implementado]

## 🧩 Componentes necessários
[Liste os widgets/componentes e quantas colunas cada um ocupa na grid 4x4]

## 📋 Regras específicas
[Liste regras de negócio, interações HTMX, validações, etc. Se não houver, escreva "Nenhuma"]

---

## ⚠️ PADRÕES OBRIGATÓRIOS (NÃO IGNORAR)

### 1. Estrutura de Layout
| Componente | Uso |
|------------|-----|
| `ShellLayout` | **Obrigatório** em todas as páginas (Topbar + Sidebar + Main Canvas) |
| `StandardGrid` | **Obrigatório** para organizar conteúdo na Main Canvas (grid 4 colunas) |
| `WidgetWrapper` | **Obrigatório** em cada componente visual dentro da grid |

### 2. Stack & Estilo
| Item | Regra |
|------|-------|
| **Stack** | Go + Templ + Tailwind + HTMX |
| **Cores/Spacing** | Apenas tokens do `tailwind.config.js` (nada de `w-[px]`, `#hex`, `p-7`) |
| **Margins** | Zero margins externas nos widgets (gap controlado pelo StandardGrid) |
| **Padding** | Apenas via WidgetWrapper (nada de padding hardcoded nos componentes) |
| **Alturas** | Automáticas (content-based), nada de `h-[px]` fixo |

### 3. HTMX
| Atributo | Uso |
|----------|-----|
| `hx-get` | Para carregamento/refresh de dados |
| `hx-target` | Sempre apontar para ID do widget (ex: `#widget-vendas`) |
| `hx-trigger` | Para auto-refresh ou eventos (ex: `every 60s`, `click`) |
| **Loading** | Todo widget deve ter estado de loading (skeleton/spinner) |

### 4. Responsividade
| Breakpoint | Comportamento |
|------------|---------------|
| Desktop | Respeitar `col-span-{1-4}` definido |
| Mobile (`md:`) | Grid colapsa para 1 coluna, todos widgets 100% largura |

### 5. Go Structs
- Structs tipadas para dados de cada componente
- Separar dados de view (nada de lógica de negócio no .templ)
- Props dos componentes devem ser structs claras e documentadas

---

## 🆕 SEÇÕES OBRIGATÓRIAS DE ESPECIFICAÇÃO TÉCNICA

### 6. 🆕 Tailwind Config (Tokens a Criar/Usar)
**Liste TODOS os tokens necessários para esta feature:**

```javascript
// tailwind.config.js - Adicionar/Verificar:
theme: {
  extend: {
    colors: {
      // Ex: 'arandu-primary': '#xxxxxx',
      [LISTE AS CORES ESPECÍFICAS DESTA FEATURE]
    },
    spacing: {
      // Ex: 'anamnesis-nav': '240px',
      [LISTE OS SPACINGS ESPECÍFICOS DESTA FEATURE]
    },
    fontFamily: {
      // Ex: 'serif-clinical': ['Source Serif 4', 'serif'],
      [LISTE AS FONTES ESPECÍFICAS DESTA FEATURE]
    },
  },
}