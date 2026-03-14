# TASK 20260313_223659

Requirement: req-01-00-01

Title: Implementar queries SQL do PatientRepository

## Objetivo

Revisar, otimizar e documentar as queries SQL do PatientRepository, garantindo performance, segurança e manutenibilidade.

## Contexto

Na tarefa anterior (`20260313_221358`), implementamos o PatientRepository com operações CRUD básicas. No entanto, as queries SQL estão embutidas diretamente nos métodos do repository, o que pode levar a:

1. **Duplicação de código:** Mesmas queries em múltiplos lugares
2. **Dificuldade de manutenção:** Mudanças requerem atualização em vários métodos
3. **Falta de otimização:** Queries não estão otimizadas para performance
4. **Risco de SQL injection:** Parâmetros não estão devidamente validados
5. **Falta de documentação:** Queries não estão documentadas

Precisamos refatorar as queries SQL para:
- Centralizar em constantes ou structs
- Otimizar performance com índices apropriados
- Adicionar validação de parâmetros
- Documentar cada query
- Garantir segurança contra SQL injection

## Análise do Estado Atual

**Arquivo:** `internal/infrastructure/repository/sqlite/patient_repository.go`

**Queries atuais:**
1. `Save()`: `INSERT INTO patients (id, name, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
2. `FindByID()`: `SELECT id, name, notes, created_at, updated_at FROM patients WHERE id = ?`
3. `FindAll()`: `SELECT id, name, notes, created_at, updated_at FROM patients ORDER BY created_at DESC`
4. `Update()`: `UPDATE patients SET name = ?, notes = ?, updated_at = ? WHERE id = ?`
5. `Delete()`: `DELETE FROM patients WHERE id = ?`

**Problemas identificados:**
1. Queries hardcoded em cada método
2. Sem validação de parâmetros além do prepared statements
3. Sem queries para casos de uso específicos (busca por nome, paginação, etc.)
4. Índices criados na migration mas não otimizados para queries específicas

## Tarefas Específicas

### 1. Centralizar Queries SQL
Criar struct ou constantes para centralizar todas as queries:
```go
type patientQueries struct {
    save    string
    findByID string
    findAll  string
    update   string
    delete   string
    // Novas queries
    findByName string
    countAll   string
}
```

### 2. Otimizar Queries Existentes
Revisar e otimizar queries atuais:
- Verificar uso de índices
- Considerar `EXPLAIN QUERY PLAN` para análise
- Otimizar cláusulas `ORDER BY` e `WHERE`

### 3. Adicionar Novas Queries Úteis
Implementar queries para casos de uso comuns:
- Busca por nome (LIKE com case-insensitive)
- Contagem total de pacientes
- Paginação (LIMIT/OFFSET)
- Filtro por data de criação

### 4. Adicionar Validação de Parâmetros
Implementar validação antes de executar queries:
- Validar comprimento de strings
- Validar formatos (UUID, datas)
- Sanitizar inputs

### 5. Documentar Queries
Adicionar documentação para cada query:
- Propósito da query
- Parâmetros esperados
- Performance characteristics
- Índices utilizados

### 6. Criar Testes para Queries
Expandir testes para cobrir:
- Queries com parâmetros inválidos
- Queries que retornam resultados vazios
- Performance de queries
- Casos de borda

## Requisitos Técnicos

### Estrutura de Queries
**Opção 1:** Constantes no pacote
```go
const (
    savePatientQuery = `INSERT ...`
    findPatientByIDQuery = `SELECT ...`
)
```

**Opção 2:** Struct com inicialização
```go
type Queries struct {
    Save string
    // ...
}

func NewQueries() *Queries {
    return &Queries{
        Save: `INSERT ...`,
        // ...
    }
}
```

**Opção 3:** Métodos que retornam queries
```go
func (r *PatientRepository) saveQuery() string {
    return `INSERT ...`
}
```

### Índices Necessários
Baseado nas queries, precisamos garantir índices para:
- `id` (já é PRIMARY KEY)
- `created_at` (já tem índice para ORDER BY)
- `name` (para buscas)
- Combinações frequentes (ex: `created_at, name`)

### Validação de Parâmetros
Regras de validação:
- `id`: UUID válido, não vazio
- `name`: 1-255 caracteres, não vazio
- `notes`: opcional, até 5000 caracteres
- `created_at`, `updated_at`: datas válidas, não zero

## Restrições de Design

1. **SQL Injection Prevention:** Sempre usar prepared statements
2. **Performance:** Queries devem usar índices apropriados
3. **Manutenibilidade:** Queries centralizadas e documentadas
4. **Testabilidade:** Queries devem ser fáceis de mockar/testar
5. **Extensibilidade:** Fácil adicionar novas queries no futuro
6. **Compatibilidade:** Manter interface existente do repository

## Critérios de Aceitação

✅ Queries SQL centralizadas em estrutura organizada  
✅ Queries existentes otimizadas e documentadas  
✅ Novas queries úteis implementadas (busca por nome, contagem)  
✅ Validação de parâmetros antes da execução  
✅ Documentação clara de cada query  
✅ Testes expandidos para cobrir novos casos  
✅ Performance mantida ou melhorada  
✅ Projeto compila e todos os testes passam  
✅ Backward compatibility mantida  

## Passos de Implementação

1. Analisar queries atuais e identificar otimizações
2. Definir estrutura para centralização de queries
3. Mover queries para estrutura centralizada
4. Implementar validação de parâmetros
5. Adicionar novas queries úteis
6. Documentar todas as queries
7. Expandir testes
8. Validar performance
9. Garantir compatibilidade com código existente

## Referências

- `docs/requirements/req-01-00-01-criar-paciente.md`
- `internal/infrastructure/repository/sqlite/patient_repository.go`
- `work/tasks/task_20260313_221358/implementation.md` (tarefa anterior do repository)
- `internal/infrastructure/repository/sqlite/migrations/0001_create_patients_table.up.sql` (índices existentes)
