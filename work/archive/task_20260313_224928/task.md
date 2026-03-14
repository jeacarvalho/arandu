# TASK 20260313_224928

Requirement: req-01-00-01

Title: TASK-05 Application Service: CreatePatient

## Objetivo

Implementar e aprimorar a camada de Application Service para a funcionalidade de criação de pacientes, seguindo os princípios de Clean Architecture e Domain-Driven Design.

## Contexto

Na arquitetura atual, temos:
- **Domínio:** Entidade Patient com regras de negócio (`NewPatient()`, `Update()`)
- **Infraestrutura:** PatientRepository com queries SQL otimizadas
- **Application Service:** PatientService com métodos básicos

No entanto, a camada de Application Service precisa ser aprimorada para:
1. **Orquestração adequada:** Coordenar domínio e infraestrutura
2. **Tratamento de erros específico:** Erros de domínio vs erros de infraestrutura
3. **Transações:** Garantir atomicidade em operações complexas
4. **Validação de aplicação:** Regras específicas da camada de aplicação
5. **DTOs/Input Models:** Estruturas de entrada validadas
6. **Testabilidade:** Services fáceis de testar com mocks

## Análise do Estado Atual

**Arquivo:** `internal/application/services/patient_service.go`

**Método `CreatePatient` atual:**
```go
func (s *PatientService) CreatePatient(name, notes string) (*patient.Patient, error) {
	p, err := patient.NewPatient(name, notes)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}
```

**Problemas identificados:**
1. **Sem DTO/Input Model:** Parâmetros soltos (`name, notes string`)
2. **Sem validação de aplicação:** Apenas validação de domínio
3. **Sem tratamento diferenciado de erros:** Todos os erros retornados igualmente
4. **Sem transações:** Operação não é atômica (problema se mais steps forem adicionados)
5. **Sem logging/observabilidade:** Não registra operações
6. **Sem testes específicos:** Testes focam no repository, não no service

## Tarefas Específicas

### 1. Criar Input Model (DTO)
Definir struct para entrada do `CreatePatient`:
```go
type CreatePatientInput struct {
    Name  string `json:"name" validate:"required,min=1,max=255"`
    Notes string `json:"notes" validate:"max=5000"`
}
```

### 2. Implementar Validação de Aplicação
Adicionar validação antes de chamar o domínio:
- Validações de formato (ex: nome não pode conter caracteres especiais)
- Validações de negócio específicas da aplicação
- Sanitização de inputs

### 3. Aprimorar Tratamento de Erros
Definir tipos de erro específicos:
- `ErrPatientAlreadyExists` (quando paciente com mesmo nome existe?)
- `ErrInvalidInput` (validação falhou)
- `ErrRepository` (erro de infraestrutura)

### 4. Adicionar Transações
Garantir atomicidade:
- Iniciar transação no início da operação
- Commit/Rollback apropriados
- Considerar context para timeout/cancelamento

### 5. Implementar Logging/Observabilidade
Adicionar logging para:
- Início/fim da operação
- Erros ocorridos
- Métricas (tempo de execução)

### 6. Criar Testes para o Service
Testes focados na camada de aplicação:
- Testes com mocks do repository
- Testes de validação
- Testes de tratamento de erros
- Testes de transações

### 7. Atualizar Interface do Service
Considerar interface mais rica:
```go
type PatientService interface {
    CreatePatient(ctx context.Context, input CreatePatientInput) (*patient.Patient, error)
    // Outros métodos com context e input models
}
```

## Requisitos Técnicos

### Arquitetura Clean Architecture
```
Input (DTO) → Application Service → Domain → Repository → Database
      ↑           ↑                    ↑          ↑
  Validation  Orchestration      Business Rules  Persistence
      ↑           ↑                    ↑          ↑
   App Rules   Error Handling     Domain Rules   SQL/NoSQL
```

### DTO/Input Model
**Propriedades:**
- Tags para validação (`validate:`)
- Tags para serialização (`json:`)
- Métodos de validação (`Validate()`)
- Métodos de sanitização (`Sanitize()`)

### Tratamento de Erros
**Hierarquia proposta:**
```
error
├── domain.Error (erros de domínio)
├── application.Error (erros de aplicação)
│   ├── validation.Error
│   ├── business.Error
│   └── infrastructure.Error
└── infrastructure.Error (erros de infraestrutura)
```

### Transações
**Abordagem:**
- Usar `context.Context` para timeout/cancelamento
- Repository suportar transações (`BeginTx`, `Commit`, `Rollback`)
- Service gerenciar ciclo de vida da transação

### Testes
**Estratégia:**
- Mock do repository usando interface
- Testes de unidade para lógica de aplicação
- Testes de integração para fluxo completo
- Table-driven tests para casos de borda

## Restrições de Design

1. **Separação de responsabilidades:** Service não deve conter lógica de domínio
2. **Dependency Injection:** Service deve receber dependências via construtor
3. **Imutabilidade:** Input models devem ser validados, não modificados
4. **Error wrapping:** Erros devem ser encapsulados com contexto
5. **Testabilidade:** Fácil de mockar e testar isoladamente
6. **Observabilidade:** Logging e métricas integrados
7. **Context propagation:** Suporte a `context.Context` para cancelamento/timeout

## Critérios de Aceitação

✅ Input Model (DTO) criado com validação  
✅ Validação de aplicação implementada  
✅ Tratamento de erros específico e rico  
✅ Transações implementadas (atomicidade)  
✅ Logging e observabilidade adicionados  
✅ Testes abrangentes para o service  
✅ Interface do service atualizada (com context)  
✅ Backward compatibility mantida onde possível  
✅ Projeto compila e todos os testes passam  
✅ Código documentado e seguindo convenções  

## Passos de Implementação

1. Analisar código atual do PatientService
2. Criar structs para input models
3. Implementar validação de aplicação
4. Adicionar tratamento de erros específico
5. Implementar suporte a transações
6. Adicionar logging e observabilidade
7. Atualizar interface do service
8. Criar testes abrangentes
9. Validar integração com camadas existentes
10. Documentar mudanças e padrões

## Referências

- `docs/requirements/req-01-00-01-criar-paciente.md`
- `internal/application/services/patient_service.go` (atual)
- `internal/domain/patient/patient.go` (domínio)
- `internal/infrastructure/repository/sqlite/patient_repository.go` (repository)
- `work/tasks/task_20260313_223659/implementation.md` (tarefa anterior de queries)
