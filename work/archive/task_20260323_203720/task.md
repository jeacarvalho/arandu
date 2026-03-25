Refatoração do Layout Unificado (SLP)

ID da Tarefa: task_20260323_refactor_layout_slp
Documento de Referência: docs/architecture/standardized_layout_protocol.md e docs/strategy/responsive_sota_strategy.md
Stack: Go, templ, Alpine.js, CSS Puro.

🎯 Objetivo

Transformar o layout base do sistema numa estrutura de 3 zonas rígidas (Top Bar, Sidebar Contextual, Main Canvas), garantindo que a navegação seja sensível ao domínio (Dashboard vs. Paciente vs. Sessão) e 100% responsiva (Mobile-First).

🛠️ Escopo Técnico

1. Novo Componente de Layout (web/components/layout/layout.templ)

Deve implementar a estrutura definida no SLP:

Zona 1 (Top Bar): Fixa no topo. Contém o logótipo, busca central (.silent-search) e perfil do médico.

Zona 2 (Sidebar): Injetada como templ.Component. Fixa no Desktop (280px).

Zona 3 (Main Canvas): Área central com fundo #E1F5EE.

2. Lógica de Contexto

O layout deve aceitar sidebarComponent e contextID.

Se estiver num paciente, a Sidebar deve exibir ações específicas (Anamnese, Metas, etc.).

3. Responsividade Dual (Desktop vs. Mobile)

No Mobile (< 768px):

Sidebar Drawer: A barra lateral deve tornar-se um Drawer (gaveta) lateral esquerdo, controlado por Alpine.js (mobileMenuOpen).

Bottom Navigation: Implementar uma barra fixa no rodapé (.bottom-nav) com os ícones de: Dashboard, Pacientes e Pesquisar.

Top Bar: O campo de busca pode ser colapsado para um ícone ou reduzido. O botão "Hambúrguer" deve aparecer à esquerda.

No Desktop (>= 768px):

Sidebar Persistente: Fixa à esquerda.

Bottom Nav: Deve ter display: none.

4. Refatoração de CSS (web/static/css/base.css)

Implementar as classes .top-bar, .contextual-sidebar, .main-canvas e .bottom-nav.

Usar @media (max-width: 768px) para alternar a visibilidade entre a Sidebar lateral e a Bottom Bar.

Importante: Eliminar qualquer atributo style="..." no HTML gerado.

🧪 Protocolo de Testes "Ironclad"

Teste de Desktop: Validar se a Sidebar ocupa o lado esquerdo e a busca está centralizada.

Teste de Mobile (375px): * Verificar se a Bottom Nav apareceu.

Verificar se ao clicar no Hambúrguer, a Sidebar desliza sobre o conteúdo.

Mudança de Contexto: Entrar num paciente e ver se as opções da Sidebar e da Bottom Nav (opcionalmente) se adaptam.

🛡️ Checklist de Integridade

[ ] O componente utiliza a Bottom Nav no mobile?

[ ] A Sidebar é um Drawer acionado por Alpine.js?

[ ] O fundo do Main Canvas é --arandu-bg (#E1F5EE)?

[ ] O scripts/arandu_guard.sh valida a presença das 3 zonas?