Task: Implementação do REQ-01-02-01 — Adicionar Observação Clínica

ID da Tarefa: task_20260315_add_observation

Requirement: REQ-01-02-01

Stack: Go, templ, HTMX, SQLite, Playwright.

🎯 Objetivo

Implementar a funcionalidade de adicionar observações clínicas atômicas a uma sessão existente. A interface deve ser minimalista ("Tecnologia Silenciosa") e a atualização deve ser feita via HTMX sem recarregar a página.

🛠️ Escopo Técnico

1. Camada de Domínio (internal/domain/observation)

Criar a entidade Observation com os campos: ID (UUID), SessionID (UUID), Content (string), CreatedAt (time.Time).

Implementar validação: O conteúdo não pode ser vazio e deve ter no máximo 5000 caracteres.

2. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Criar a tabela observations via migration (se não existir).

Implementar Save(ctx, observation) e GetBySessionID(ctx, sessionID).

3. Camada de Aplicação (internal/application/services)

Criar AddObservation(ctx, sessionID, content).

Orquestrar a criação do UUID e a persistência.

4. Camada Web (Componentes templ)

ObservationItem: Componente que renderiza uma única observação (usar fonte Source Serif para o conteúdo).

ObservationList: Componente que lista as observações de uma sessão.

ObservationForm: Formulário minimalista com hx-post para /sessions/{id}/observations e hx-target para a lista. Deve usar hx-on::after-request="this.reset()" para limpar o campo.

5. Camada Web (Handlers)

POST /sessions/{id}/observations:

Receber o conteúdo.

Chamar o serviço.

Retornar apenas o componente ObservationItem renderizado (fragmento HTMX).

🎨 Design System e UX

Fonte: O campo textarea de entrada e a exibição da observação DEVEM usar a classe CSS que define a Source Serif 4.

Estética: Remover bordas pesadas. Usar um estilo de "nota de margem" para as observações já registradas.

HTMX: Usar hx-swap="afterbegin" na lista de observações para que a nota mais recente apareça no topo instantaneamente.

🧪 Protocolo de Testes "Ironclad"

A. Testes Unitários

Validar que o repositório salva e recupera observações corretamente do SQLite.

Validar que o serviço rejeita observações vazias.

B. Teste E2E (Go + Playwright)

O teste deve obrigatoriamente realizar os seguintes passos:

Iniciar o servidor Arandu.

Navegar até a página de uma sessão existente (/session/{id}).

Verificar se o formulário de observação está visível.

Digitar uma percepção clínica (ex: "Paciente demonstrou resistência ao falar sobre a infância").

Clicar em "Adicionar" (ou submeter o form).

Verificar via Playwright:

O campo de texto foi limpo automaticamente.

A página NÃO foi recarregada (verificar via interceptação de rede ou log de console).

O texto "Paciente demonstrou..." agora aparece na lista de observações na tela.

O texto está renderizado com a fonte Source Serif.

📋 Critérios de Aceitação

[ ] Código .templ compilado com sucesso.

[ ] Persistência funcional no banco de dados.

[ ] Atualização assíncrona via HTMX sem "flicker" de tela.

[ ] Teste E2E Playwright passando 100%.

[ ] Arquivo de learnings atualizado com o status da implementação.