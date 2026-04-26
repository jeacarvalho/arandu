# TASK 20260424_164758
## Redesign do Painel de Classificação de Observações e Intervenções

**Status:** PRONTO_PARA_IMPLEMENTACAO  
**Requisito:** REQ-03-02-01

---

## 🎯 Objetivo

Redesenhar o painel de classificação de observações e intervenções clínicas para ser visualmente coeso com o design system Arandu (DaisyUI + Fraunces/Geist) e mais usável: o terapeuta deve conseguir selecionar uma tag em **1 clique**, sem dropdown, sem fricção.

O backend está completo. A tarefa é exclusivamente de UI.

---

## 🏗️ Contexto do sistema

**Stack:** Go 1.22+ · Templ · HTMX · DaisyUI v4 + Tailwind CSS · SQLite

**Design system:**
- Fonte UI: Geist (`font-sans`)
- Fonte clínica: Fraunces (`font-serif` / classe `.serif`)
- Tokens: `--paper`, `--ink`, `--accent`, `--line`
- Componentes DaisyUI: `badge`, `btn`, `btn-ghost`
- **Não criar CSS custom** — usar classes DaisyUI + tokens existentes inline style quando necessário

**Estrutura de dados disponível nos handlers:**

```go
// TagSelectorData — injetado nos componentes
type TagSelectorData struct {
    InterventionID string   // ou ObservationID
    AvailableTags  []*Tag   // 42 tags pré-definidas
    SelectedTags   []*InterventionClassification
}

type Tag struct {
    ID      string
    Name    string
    TagType TagType  // cognitive | behavioral | emotional | psychoeducation | narrative | body
    Color   string   // hex, ex: "#7C3AED"
    Icon    string
}

type InterventionClassification struct {
    TagID     string
    Tag       *Tag
    Intensity int  // 1-5, opcional
}

// Tipos com cores
// cognitive       → #7C3AED (roxo)
// behavioral      → #1D9E75 (verde)
// emotional       → #0F6E56 (verde base)
// psychoeducation → #F59E0B (âmbar)
// narrative       → #3B82F6 (azul)
// body            → #DC2626 (vermelho)
```

**Rótulos em PT-BR para os tipos:**
```
cognitive       → Cognitiva
behavioral      → Comportamental
emotional       → Emocional
psychoeducation → Psicoeducação
narrative       → Narrativa
body            → Corporal
```

---

## 🗂️ Arquivos a modificar

**Foco total nestes arquivos (não alterar mais nada):**

- `web/components/intervention/tag_selector_grid.templ` — painel principal que o handler retorna via `GET /interventions/{id}/classify/edit`
- Se existir equivalente para observações (ex: `web/components/classification/tag_selector_grid.templ`) — aplicar o mesmo redesign
- `web/components/intervention/tag_badge.templ` — badge de tag já selecionada (exibida abaixo do item)

**Verificar se existem e ajustar se necessário:**
- `web/components/intervention/tag_selector.templ` — pode ser componente auxiliar
- `web/components/session/intervention_item.templ` — onde os badges de tags selecionadas são exibidos

**Não alterar:**
- Handlers (`.go`)
- Migrations (`.sql`)
- Rotas (`cmd/arandu/main.go`)
- CSS files (`style.css`, `tailwind-v2.css`)

---

## 🎨 Especificação do novo design

### Painel de seleção (substitui o dropdown atual)

O painel abre inline, dentro do `#intervention-{id}-tags` ou `#observation-{id}-tags`, via HTMX `innerHTML` swap — igual ao comportamento atual.

**Layout esperado:**

```
┌─────────────────────────────────────────────────┐
│  Classificar                              [✕]    │
│                                                  │
│  COGNITIVA                                       │
│  [Reestruturação cognitiva] [Quest. socrático]   │
│  [Identificação distorções] [Registro pens.]     │
│                                                  │
│  COMPORTAMENTAL                                  │
│  [Exposição gradual] [Ativação comportamental]   │
│  ...                                             │
│                                                  │
│                              [Salvar]            │
└─────────────────────────────────────────────────┘
```

**Regras de implementação:**

1. **Sem dropdown** — mostrar todos os 6 grupos diretamente, um abaixo do outro
2. **Tags como chips clicáveis:** `badge badge-outline` DaisyUI com `cursor-pointer`. Ao selecionar: badge filled com a cor do tipo (via `style="background:{color};color:white;border-color:{color}"`). Não selecionado: outline cinza.
3. **Estado via checkboxes hidden com labels:**

```templ
<form hx-post={ fmt.Sprintf("/interventions/%s/classify", data.InterventionID) }
      hx-target={ fmt.Sprintf("#intervention-%s-tags", data.InterventionID) }
      hx-swap="innerHTML">
  for _, group := range groupedTags {
    <div class="mb-3">
      <div class="text-xs font-semibold uppercase tracking-wider mb-1"
           style={ fmt.Sprintf("color:%s", group.Color) }>
        { group.Label }
      </div>
      <div class="flex flex-wrap gap-1">
        for _, tag := range group.Tags {
          <label class={ badgeClass(isSelected(tag.ID, data.SelectedTags)) }
                 style={ badgeStyle(tag.Color, isSelected(tag.ID, data.SelectedTags)) }>
            <input type="checkbox" name="tag_ids" value={ tag.ID }
                   class="hidden"
                   if isSelected(tag.ID, data.SelectedTags) { checked }/>
            { tag.Name }
          </label>
        }
      </div>
    </div>
  }
  <div class="flex justify-end mt-3">
    <button type="submit" class="btn btn-sm btn-primary">Salvar</button>
  </div>
</form>
```

4. **Intensidade:** remover os círculos 1-5 da interface principal. Complexidade desnecessária. Pode enviar `intensity=1` como padrão no POST se o handler exigir.
5. **Botão fechar (✕):** `<button type="button" onclick="this.closest('[id$=-tags]').innerHTML=''" class="btn btn-ghost btn-xs">✕</button>`
6. **Agrupamento de tags no Go:** criar função helper que agrupa `AvailableTags` por `TagType` antes de passar ao template (ou fazer no próprio templ com loop + verificação de tipo)

### Badges de tags selecionadas (abaixo do item após salvar)

O handler já retorna os badges via `InterventionTagsWrapper`. O `tag_badge.templ` deve exibir:

```
[● Reestruturação cognitiva ×]  [● Validação emocional ×]
```

- Ponto colorido `●` (ou pequeno dot via `style="color:{tag.Color}"`)
- Nome da tag
- `×` que aciona `DELETE /interventions/{id}/classify/{tag_id}` via HTMX
- Estilo DaisyUI: `badge badge-outline gap-1` com border e color do tipo
- O `×` deve ter `hx-delete`, `hx-target` apontando para o wrapper de tags, `hx-swap="outerHTML"`

---

## 📋 Critérios de aceite

**Compilação:**
- [ ] `~/go/bin/templ generate ./web/components/...` sem erros
- [ ] `go build ./...` sem erros

**Comportamento:**
- [ ] CA-01: Clicar no ícone de tag abre o painel inline sem reload
- [ ] CA-02: Todas as 6 categorias são visíveis diretamente (sem dropdown para filtrar)
- [ ] CA-03: Clicar em uma tag a seleciona (toggle visual: badge filled vs outline)
- [ ] CA-04: Clicar em "Salvar" persiste as tags e fecha o painel, mostrando os badges
- [ ] CA-05: Badges exibidos abaixo do item com a cor do tipo e botão ×
- [ ] CA-06: Clicar × em um badge remove a tag via HTMX sem reload da sessão
- [ ] CA-07: O mesmo design vale para observações (se o componente for separado, aplicar idem)

**Qualidade:**
- [ ] Não usar CSS custom fora dos tokens existentes
- [ ] Não alterar arquivos `.go`, `.sql` ou CSS globais

---

## 🚫 NÃO faça

- Não alterar os handlers Go
- Não criar tabelas ou migrations
- Não usar `html/template` — apenas `.templ`
- Não quebrar o HTMX swap do painel (o target é `#intervention-{id}-tags` com `innerHTML`)
- Não reintroduzir o dropdown de categoria

---

## 📎 Padrão de referência

- `web/components/session/intervention_item.templ` — estrutura do item com footer
- Cores dos tipos estão na migration `0012_add_intervention_tags.up.sql` e na struct `Tag.Color`
