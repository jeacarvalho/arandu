# TASK 20260313_221358

Requirement: req-01-00-01

Title: Criar Repository do Patient

## Objetivo

Implementar o repositório para a entidade Patient, seguindo os princípios de Clean Architecture e DDD.

## Contexto

Na tarefa anterior (20260313_215938), implementamos a entidade de domínio Patient com:
- Construtor `NewPatient()` com validação
- Geração automática de ID (UUID)
- Timestamps (CreatedAt, UpdatedAt)
- Testes de unidade

Agora precisamos garantir que o repositório está adequadamente implementado para persistir a entidade Patient.

## Análise do Estado Atual

O repositório já existe em `internal/infrastructure/repository/sqlite/patient_repository.go` e inclui:
- Métodos CRUD completos (Save, FindByID, FindAll, Update, Delete)
- Inicialização de schema (`InitSchema()`)
- Integração com SQLite

No entanto, após a refatoração da entidade de domínio, precisamos verificar:
1. Se o repositório está usando corretamente a entidade de domínio
2. Se há alguma lógica de domínio vazando para a camada de infraestrutura
3. Se os testes estão adequados

## Tarefas Específicas

### 1. Revisar Implementação do Repository
**Arquivo:** `internal/infrastructure/repository/sqlite/patient_repository.go`

Verificar:
- ✅ O método `Save()` não deve gerar ID ou timestamps (já feito na tarefa anterior)
- ✅ O método `Save()` deve apenas persistir os dados
- ✅ O método `Update()` deve atualizar apenas `UpdatedAt`
- ✅ Consultas SQL devem mapear corretamente para a struct Patient

### 2. Verificar Interface do Repository
**Arquivo:** `internal/domain/patient/patient.go`

A interface `Repository` já está definida com:
```go
type Repository interface {
	Save(patient *Patient) error
	FindByID(id string) (*Patient, error)
	FindAll() ([]*Patient, error)
	Update(patient *Patient) error
	Delete(id string) error
}
```

Verificar se:
- ✅ A interface está completa para os requisitos atuais
- ✅ Todos os métodos são necessários para o requirement req-01-00-01

### 3. Testar Integração
Criar ou atualizar testes de integração para verificar:
- Criação de paciente através do service → repository
- Persistência correta no SQLite
- Recuperação de dados

### 4. Verificar Service Layer
**Arquivo:** `internal/application/services/patient_service.go`

Após a tarefa anterior, o service já usa `NewPatient()`. Verificar:
- ✅ O método `CreatePatient()` chama `NewPatient()` e depois `repository.Save()`
- ✅ Tratamento adequado de erros
- ✅ Validações de domínio são respeitadas

### 5. Executar Testes End-to-End
Executar testes para garantir que todo o fluxo funciona:
1. Criação de paciente via service
2. Persistência no banco
3. Recuperação do banco
4. Validação dos dados

## Requisitos Técnicos

### Arquitetura
- **Clean Architecture:** Repository na camada de infraestrutura
- **DDD:** Repository implementa interface definida no domínio
- **SOLID:** Single Responsibility, Dependency Inversion

### Banco de Dados
- **SQLite:** Banco atual do projeto
- **Schema:** Tabela `patients` com campos:
  - `id TEXT PRIMARY KEY`
  - `name TEXT NOT NULL`
  - `notes TEXT`
  - `created_at DATETIME NOT NULL`
  - `updated_at DATETIME NOT NULL`

### Dependências
- `github.com/google/uuid`: Para geração de ID (já no domínio)
- `database/sql`: Para acesso ao banco
- `github.com/mattn/go-sqlite3`: Driver SQLite

## Restrições de Design

1. **Sem lógica de domínio no repository:** Apenas persistência
2. **Sem validações no repository:** Validações devem estar no domínio
3. **Interface first:** Implementar a interface definida no domínio
4. **Error handling apropriado:** Retornar erros específicos quando aplicável
5. **Transações:** Considerar uso de transações para operações complexas (futuro)

## Critérios de Aceitação

✅ Repository implementa corretamente a interface do domínio  
✅ Método `Save()` persiste paciente com ID e timestamps gerados no domínio  
✅ Método `FindByID()` recupera paciente corretamente  
✅ Método `FindAll()` retorna lista ordenada por `created_at DESC`  
✅ Método `Update()` atualiza apenas campos permitidos  
✅ Método `Delete()` remove paciente do banco  
✅ Schema da tabela está correto e é inicializado  
✅ Todos os testes passam  
✅ Projeto compila sem erros  

## Passos de Implementação

1. Revisar código atual do repository
2. Corrigir qualquer problema encontrado
3. Criar/atualizar testes de integração
4. Executar testes completos
5. Validar fluxo end-to-end
6. Documentar alterações

## Referências

- `docs/requirements/req-01-00-01-criar-paciente.md`
- `work/tasks/task_20260313_215938/implementation.md` (tarefa anterior)
- `internal/domain/patient/patient.go` (entidade de domínio)
- `internal/application/services/patient_service.go` (service layer)
