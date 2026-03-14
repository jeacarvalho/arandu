# TASK 20260313_232014

Requirement: docs/requirements/req-01-00-01-criar-paciente.md

Title: Patient: Implementar UI/UX

## Objetivo

Implementar a interface de usuário (UI/UX) para gerenciamento de pacientes conforme especificado no requisito REQ-01-00-01, integrando com o backend já implementado.

## Contexto

O backend para criação de pacientes já foi implementado:
- Domain entity com validação
- Repository com SQLite
- Sistema de migrações
- Application service com DTOs
- HTTP handlers atualizados

Agora precisamos implementar a interface web que permitirá aos usuários:
1. Visualizar lista de pacientes
2. Criar novos pacientes
3. Visualizar detalhes de um paciente

## Escopo da Implementação

### 1. Página de Listagem de Pacientes (`GET /patients`)
- Listar todos os pacientes com paginação
- Mostrar nome, data de criação, última atualização
- Link para criar novo paciente
- Link para visualizar detalhes de cada paciente

### 2. Página de Criação de Paciente (`GET /patients/new`)
- Formulário com campos:
  - Nome (obrigatório)
  - Observações (opcional, textarea)
- Botão "Salvar paciente"
- Validação client-side básica
- Feedback visual de sucesso/erro

### 3. Página de Detalhes do Paciente (`GET /patient/{id}`)
- Exibir todas as informações do paciente
- Mostrar ID, nome, observações, datas
- Espaço reservado para futuras sessões
- Botão para voltar à lista

### 4. Integração com Backend
- Consumir endpoints REST já implementados
- Tratar respostas de sucesso/erro
- Implementar redirecionamentos conforme fluxo

## Requisitos Técnicos

### Tecnologias a Utilizar
- HTML5 semântico
- CSS3 com design responsivo
- JavaScript vanilla (sem frameworks pesados)
- Go templates para renderização server-side
- Integração com handlers existentes em `web/handlers/`

### Estrutura de Arquivos
- `web/templates/patients/` - Templates para pacientes
- `web/static/css/patients.css` - Estilos específicos
- `web/static/js/patients.js` - JavaScript para interatividade

### Design e UX
- Interface limpa e profissional
- Foco em usabilidade clínica
- Feedback visual claro para ações
- Design responsivo (mobile/desktop)
- Acessibilidade básica (labels, alt text)

## Critérios de Aceitação a Implementar

### CA-01: Criar paciente apenas com nome
- Formulário deve permitir submeter apenas com nome preenchido
- Campo observações deve ser opcional

### CA-02: Identificador único gerado
- Página de detalhes deve mostrar ID único do paciente

### CA-03: Persistência confirmada
- Após criação, paciente deve aparecer na lista
- Dados devem persistir entre sessões

### CA-04: Redirecionamento após criação
- Após salvar, usuário deve ser redirecionado para página do paciente

### CA-05: Lista de pacientes
- Paciente recém-criado deve aparecer na lista principal

## Fluxo de Implementação

1. **Estrutura de templates** - Criar diretório e arquivos base
2. **Listagem de pacientes** - Implementar página principal
3. **Formulário de criação** - Criar interface de cadastro
4. **Página de detalhes** - Mostrar informações do paciente
5. **Estilização** - Aplicar CSS para melhor UX
6. **JavaScript** - Adicionar interatividade básica
7. **Testes** - Verificar integração com backend
8. **Documentação** - Atualizar documentação do requisito

## Referências

- `docs/requirements/req-01-00-01-criar-paciente.md` - Requisito original
- `web/handlers/handler.go` - Handlers HTTP já implementados
- `internal/application/services/patient_service.go` - Service layer
- `internal/infrastructure/repository/sqlite/` - Implementação do repositório
- `work/tasks/task_20260313_224928/` - Task anterior (application service)

## Notas de Design

- Manter consistência com futuras interfaces do sistema
- Priorizar simplicidade e eficiência clínica
- Considerar que psicólogos podem usar em consultório
- Interface deve carregar rapidamente mesmo em conexões lentas

## Validações

- Testar fluxo completo: criar → listar → visualizar
- Verificar responsividade em diferentes dispositivos
- Validar acessibilidade básica
- Testar com dados reais (nomes longos, caracteres especiais)
- Garantir integração com backend existente
