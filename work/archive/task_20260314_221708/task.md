---

# Task: Implementação do REQ-01-01-02 — Edição de Sessão Clínica

**ID da Tarefa:** `task_20260314_session_edit`

**Requirement Relacionado:** `REQ-01-01-02`

Leia os arquivos docs/requirements/req-01-01-02-editar-sessao.md e docs/capabilities/cap-01-01-registro-sessoes.md

**Contexto de Arquitetura:** DDD, Clean Architecture, Go, SQLite, HTMX. 

## 🎯 Objetivo

Implementar a funcionalidade completa de edição de sessões clínicas, garantindo que o terapeuta possa corrigir datas e resumos, mantendo a integridade do domínio e a persistência em SQLite. 

## 🛠️ Escopo Técnico

### 1. Camada de Domínio (`internal/domain/session`)

* Garantir que a entidade `Session` possua lógica de validação para atualizações (data não futura, conteúdo válido). 


* O campo `UpdatedAt` deve ser atualizado no momento da modificação. 



### 2. Camada de Infraestrutura (`internal/infrastructure/repository/sqlite`)

* Implementar o método `Update(ctx, session)` no `SessionRepository`. 


* Assegurar que apenas os campos `date`, `summary` e `updated_at` sejam alterados no SQL. 



### 3. Camada de Aplicação (`internal/application/services`)

* Criar o método `UpdateSession(ctx, input)` no `SessionService`. 


* Tratar erros de "sessão não encontrada" e erros de validação de domínio. 



### 4. Camada Web (`web/handlers` & `web/templates`)

* **Rota GET `/sessions/{id}/edit**`: Retornar o template `session_edit.html`.
* **Rota POST/PUT `/sessions/{id}**`: Processar a atualização e redirecionar para `/patient/{patient_id}` ou `/session/{id}`.
* 
**UI (Tecnologia Silenciosa)**: O formulário de edição deve usar obrigatoriamente a classe `.clinical-font` (Source Serif) no campo de resumo para manter a imersão. 



---

## 🧪 Protocolo de Testes "Ironclad" (Obrigatório)

A tarefa **NÃO** será considerada concluída se os testes abaixo falharem ou forem omitidos.

### A. Testes Unitários (Cobertura de Lógica)

* **Repository**: Testar o `Update` real em um banco SQLite em memória, verificando se os dados mudaram e se o `updated_at` foi alterado.
* **Service**: Testar a orquestração do serviço, garantindo que ele chama o repositório corretamente e valida os dados de entrada.

### B. Teste E2E (Go + Playwright) — Validação de Caminho Feliz

Você deve criar/rodar um teste E2E que execute os seguintes passos reais no navegador (headless):

1. **Boot Up**: Iniciar o servidor Arandu em uma porta de teste. O teste deve falhar imediatamente se o servidor não subir.
2. **Preparação**: Criar (via DB ou API) um paciente e uma sessão de teste.
3. **Navegação**: Acessar a URL `/sessions/{id}/edit`.
4. **Interação**:
* Verificar se o campo de texto contém o resumo original.
* Alterar a data para um dia anterior.
* Alterar o texto do resumo.
* Clicar em "Salvar".


5. **Verificação de UI**: Confirmar se o sistema redirecionou para a página correta (Status 200) e se os novos dados aparecem na tela sem erros de layout.
6. **Verificação de DB**: Consultar o SQLite diretamente para confirmar que os dados persistidos batem com o que foi digitado.

---

## 📋 Critérios de Aceitação (Checklist de Saída)

* [ ] O código compila sem warnings.
* [ ] Testes unitários com >80% de cobertura nas camadas afetadas.
* [ ] Teste E2E Playwright confirmando que o usuário consegue editar e salvar sem erros de "tela branca" ou "404".
* [ ] O campo de resumo na UI está usando a tipografia correta (Serif).
* [ ] Arquivo de `learning` atualizado com a conclusão da tarefa e o status dos testes.

---
