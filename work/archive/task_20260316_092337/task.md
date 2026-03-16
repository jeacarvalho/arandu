Task: Implementação do REQ-01-02-02 — Editar Observação Clínica

ID da Tarefa: task_20260316_edit_observation

Requirement Relacionado: REQ-01-02-02

Stack Técnica: Go, templ, HTMX, SQLite, Playwright.

Status: CONCLUIDA

🎯 Objetivo

Implementar a funcionalidade de edição "inline" de observações clínicas. O terapeuta deve conseguir alternar entre o modo de leitura e o modo de edição sem recarregar a página, garantindo a persistência dos dados e a fidelidade ao Design System Arandu.

🛠️ Escopo Técnico

1. Camada de Domínio (internal/domain/observation)

Garantir que a entidade Observation suporte a actualização do campo Content.

O campo UpdatedAt deve ser definido obrigatoriamente no momento da modificação.

2. Camada de Infraestrutura (internal/infrastructure/repository/sqlite)

Implementar o método Update(ctx, observation) no repositório.

A query SQL deve actualizar apenas content e updated_at filtrando pelo id.

3. Camada de Aplicação (internal/application/services)

Implementar UpdateObservation(ctx, id, content).

Validar a existência da observação antes da tentativa de actualização.

4. Camada Web (Componentes templ)

ObservationItem: Actualizar para incluir um botão/ícone de edição que dispara um hx-get para /observations/{id}/edit com hx-target="closest .observation-container".

ObservationEditForm: Criar este novo componente que contém:

Um textarea com a classe .clinical-font (Source Serif 4).

Atributo hx-put para /observations/{id}.

Atributo hx-swap="outerHTML".

Botão "Cancelar" que dispara um hx-get para /observations/{id} para restaurar o modo de leitura.

5. Camada Web (Handlers)

GET /observations/{id}/edit: Retorna o componente ObservationEditForm.

PUT /observations/{id}: Processa a alteração e retorna o componente ObservationItem actualizado.

GET /observations/{id}: Retorna o componente ObservationItem (usado para cancelar a edição).

🎨 Design System e UI/UX

Tipografia: O campo de edição deve usar obrigatoriamente a fonte Source Serif 4.

Silent Input: O formulário de edição não deve ter bordas pesadas; deve parecer uma edição directa sobre o "papel".

Feedback: Utilizar transições suaves do HTMX para a troca de componentes.

🧪 Protocolo de Testes "Ironclad"

A. Testes Unitários e de Integração

Repository: Testar o Update num banco SQLite real (em memória) e verificar se o updated_at foi modificado.

Domain: Validar que conteúdos vazios ou acima de 5000 caracteres são rejeitados.

B. Teste E2E (Go + Playwright)

O agente deve validar o fluxo completo no navegador:

Iniciar o servidor Arandu.

Aceder à página de uma sessão com observações já existentes.

Clicar no botão "Editar" de uma observação.

Validar: O texto transformou-se num campo de entrada e os botões "Salvar/Cancelar" apareceram.

Alterar o texto e clicar em "Cancelar".

Validar: O texto original voltou e o campo de edição desapareceu.

Clicar em "Editar" novamente, alterar o texto e clicar em "Salvar".

Validar:

A requisição foi PUT e não houve recarregamento total da página.

O novo texto aparece na lista.

A fonte utilizada é a Source Serif 4.

Ao actualizar a página (F5), o novo texto persiste (leitura do DB).

📋 Critérios de Aceitação

[x] Componentes .templ gerados e integrados sem erros de namespace.

[x] Fluxo HTMX (GET edit -> PUT update) funcionando inline.

[x] Persistência correcta no SQLite (campo updated_at preenchido).

[x] Teste E2E Playwright confirmando o sucesso da edição e do cancelamento.

[x] Interface respeita o Design System (Tecnologia Silenciosa).

## ✅ Implementação Concluída

### Resumo da Implementação

1. **Camada de Domínio (observation.go)**
   - Adicionado campo `UpdatedAt time.Time` à struct `Observation`
   - Interface `Repository` já incluía método `Update`

2. **Camada de Infraestrutura (observation_repository.go)**
   - Atualizado schema da tabela para incluir campo `updated_at DATETIME`
   - Implementado método `Update` que atualiza `content` e `updated_at`
   - Atualizados métodos `FindByID`, `FindBySessionID`, `FindAll` para ler `updated_at`
   - Atualizado método `Save` para definir `updated_at` como NULL inicialmente

3. **Camada de Aplicação (observation_service.go)**
   - Implementado método `UpdateObservation` com validações:
     - Conteúdo não pode ser vazio
     - Conteúdo não pode exceder 5000 caracteres
     - Observação deve existir

4. **Camada Web - Componentes Templ**
   - Criado `ObservationEditForm.templ` com:
     - Textarea usando fonte Source Serif 4 (.clinical-font)
     - Formulário com hx-put para /observations/{id}
     - Botões "Salvar" e "Cancelar" (cancelar usa hx-get para /observations/{id})
   - Atualizado `ObservationItem.templ` com:
     - Botão de edição com hx-get para /observations/{id}/edit
     - Ícone de lápis (fa-edit)

5. **Camada Web - Handlers**
   - Criado `observation_handler.go` com:
     - `GET /observations/{id}` - Retorna componente ObservationItem
     - `GET /observations/{id}/edit` - Retorna componente ObservationEditForm  
     - `PUT /observations/{id}` - Processa atualização e retorna ObservationItem atualizado
   - Atualizado `service_adapters.go` para suportar nova interface
   - Atualizado `main.go` para registrar rotas

6. **Testes**
   - Testes unitários para serviço (observation_service_test.go)
   - Testes de integração para repositório (observation_repository_test.go)
   - Testes para handlers (observation_handler_test.go)
   - Testes E2E (observation_edit_test.go)
   - Todos os testes passando

7. **Aprendizado Documentado**
   - Criado `REQ-01-02-02.md` em docs/learnings/
   - Documentado problema de importação duplicada em arquivos .templ
   - Padrão estabelecido: NUNCA adicionar imports manualmente em arquivos .templ

### Fluxo HTMX Implementado
1. Usuário clica em "Editar" em uma observação
2. `GET /observations/{id}/edit` → Retorna `ObservationEditForm` com conteúdo atual
3. Usuário edita texto e clica "Salvar"
4. `PUT /observations/{id}` → Atualiza banco e retorna `ObservationItem` atualizado
5. Usuário clica "Cancelar" → `GET /observations/{id}` restaura visualização original

### Validações Implementadas
- Conteúdo não pode ser vazio
- Conteúdo máximo de 5000 caracteres
- Observação deve existir
- Campo `updated_at` sempre atualizado na modificação
- Campo `created_at` preservado intacto