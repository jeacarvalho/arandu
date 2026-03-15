# Prompt de Desenvolvimento: Refatoração de Layout Mestre (Master Layout)

**Contexto:** Refatorar a interface do Arandu para adotar um padrão de "Single Layout" com barra lateral (Sidebar) retrátil e consistência visual global em todas as páginas.

## 1. Estrutura de Templates (Go `html/template`)

Você deve implementar um arquivo `web/templates/layout.html` que servirá de casca única.

* O layout deve usar o estado do **Alpine.js** (`x-data="{ sidebarOpen: true }"`) para controlar a visibilidade da sidebar.
* Use a tag `{{block "content" .}}{{end}}` no centro da área principal para injetar o conteúdo de cada página específica.

## 2. Requisitos da Sidebar (Navegação)

* **Posição:** Lateral esquerda (conforme `dashAtual.html`).
* **Comportamento:** Deve possuir um botão de "Toggle" (ícone de hambúrguer ou seta) que recolhe a barra.
* **Estado Recolhido:** Quando recolhida, a sidebar deve mostrar apenas os ícones ou desaparecer completamente (conforme preferência do design silencioso), expandindo a área de conteúdo principal.
* **Persistência:** O menu deve estar presente em `/dashboard`, `/patients`, `/patient/{id}` e `/session/{id}`.

## 3. Padronização Estética (CSS Variables)

Centralize todos os estilos no arquivo `web/static/css/style.css` usando variáveis CSS baseadas no **Design System Arandu**:

* **Cores:** `--arandu-primary: #1E3A5F;`, `--arandu-secondary: #3A7D6B;`, `--arandu-insight: #D4A84F;`, `--arandu-bg: #F7F8FA;`.
* **Tipografia:** Aplicar `Inter` para a UI (menus, botões) e `Source Serif 4` para todo o conteúdo clínico (notas de pacientes, resumos de sessão).
* **Cards e Botões:** Todos os cards devem seguir o padrão `card-hover` do `dashAtual.html`, e os botões devem ter estados `hover` e `active` consistentes.

## 4. Refatoração dos Handlers (Go)

Para evitar que o sistema se perca, todos os handlers de visualização devem agora:

1. Carregar os dados necessários (Patients, Sessions, Insights).
2. Chamar o template principal: `tmpl.ExecuteTemplate(w, "layout", data)`.
3. Garantir que o `data` contenha as informações para o "miolo" da página e metadados para marcar o item ativo no menu (ex: `CurrentPage: "patients"`).

## 5. Protocolo de Verificação (Testes)

* **Check de Consistência:** Ao navegar de "Dashboard" para "Paciente", a sidebar NÃO deve recarregar visualmente de forma brusca; apenas o bloco central deve mudar.
* **Teste de Responsividade:** A sidebar deve se comportar bem em telas menores, preferencialmente escondendo-se automaticamente.
* **E2E Playwright:** Incluir um teste que:
1. Abre o Dashboard.
2. Clica no botão de recolher sidebar.
3. Verifica se a largura da sidebar diminuiu e o conteúdo expandiu.
4. Navega para a lista de pacientes e confirma que a sidebar permanece no estado esperado.