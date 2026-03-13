# Arquitetura do Sistema Arandu

## Visão Geral

O Arandu é um sistema de inteligência clínica para profissionais de saúde mental, construído com uma arquitetura baseada em Domain-Driven Design (DDD) e Clean Architecture.

## Estrutura do Projeto

```
arandu/
├── cmd/arandu/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── domain/                 # Camada de domínio (core business)
│   │   ├── patient/           # Entidade Paciente
│   │   ├── session/           # Entidade Sessão
│   │   ├── observation/       # Entidade Observação
│   │   ├── intervention/      # Entidade Intervenção
│   │   └── insight/           # Entidade Insight
│   ├── application/           # Camada de aplicação
│   │   └── services/          # Serviços de aplicação
│   └── infrastructure/        # Camada de infraestrutura
│       └── repository/
│           └── sqlite/        # Implementações de repositório SQLite
├── web/                       # Camada web
│   ├── handlers/              # Handlers HTTP
│   ├── templates/             # Templates HTML
│   └── static/                # Arquivos estáticos
│       ├── css/
│       └── js/
└── docs/                      # Documentação
```

## Princípios Arquiteturais

### 1. Domain-Driven Design (DDD)
- **Domínio Rico**: Entidades com comportamento e regras de negócio
- **Agregados**: Paciente como agregado raiz, Sessão como agregado filho
- **Repositórios**: Interfaces definidas no domínio, implementadas na infraestrutura

### 2. Clean Architecture
- **Independência de Camadas**: Domínio não depende de outras camadas
- **Inversão de Dependência**: Interfaces apontam para dentro
- **Testabilidade**: Domínio puro pode ser testado sem infraestrutura

### 3. SOLID
- **Single Responsibility**: Cada classe/pacote tem uma única responsabilidade
- **Open/Closed**: Aberto para extensão, fechado para modificação
- **Liskov Substitution**: Substituição segura de implementações
- **Interface Segregation**: Interfaces específicas para cada cliente
- **Dependency Inversion**: Dependa de abstrações, não de implementações

## Modelo de Domínio

### Entidades Principais

#### 1. Paciente (`domain/patient`)
- Representa um paciente em tratamento
- Atributos: ID, Nome, Notas, Datas de criação/atualização
- Agregado raiz para sessões

#### 2. Sessão (`domain/session`)
- Representa uma sessão terapêutica
- Atributos: ID, PacienteID, Data, Notas, Datas de criação/atualização
- Pertence a um paciente

#### 3. Observação (`domain/observation`)
- Observação clínica feita durante uma sessão
- Atributos: ID, SessãoID, Conteúdo, Data de criação
- Pertence a uma sessão

#### 4. Intervenção (`domain/intervention`)
- Intervenção terapêutica realizada
- Atributos: ID, SessãoID, Conteúdo, Datas de criação/atualização
- Pertence a uma sessão

#### 5. Insight (`domain/insight`)
- Insight gerado por IA ou terapeuta
- Atributos: ID, Conteúdo, Fonte ("ai" ou "therapist"), Data de criação

## Camada de Aplicação

### Serviços de Aplicação (`application/services`)
- **PatientService**: Gerencia operações com pacientes
- **SessionService**: Gerencia operações com sessões
- **ObservationService**: Gerencia operações com observações
- **InterventionService**: Gerencia operações com intervenções
- **InsightService**: Gerencia operações com insights

Cada serviço:
- Orquestra operações de domínio
- Não contém lógica de persistência
- Implementa casos de uso específicos

## Camada de Infraestrutura

### Repositórios SQLite (`infrastructure/repository/sqlite`)
- **PatientRepository**: Implementação SQLite para repositório de pacientes
- **SessionRepository**: Implementação SQLite para repositório de sessões
- **ObservationRepository**: Implementação SQLite para repositório de observações
- **InterventionRepository**: Implementação SQLite para repositório de intervenções
- **InsightRepository**: Implementação SQLite para repositório de insights

Princípios de implementação:
- Sem ORM (usa `database/sql` diretamente)
- SQL explícito e simples
- Transações quando necessário
- Migrações manuais

## Camada Web

### Handlers HTTP (`web/handlers`)
- **Handler**: Handler principal com injeção de dependências
- **Dashboard**: Página inicial com visão geral
- **Patients**: Lista de pacientes
- **Patient**: Detalhes do paciente e histórico
- **Session**: Detalhes da sessão com observações e intervenções

Características:
- Handlers finos (apenas recebem request, chamam serviço, retornam response)
- Compatível com HTMX (retorna HTML fragments quando necessário)
- Templates separados da lógica

### Templates (`web/templates`)
- **layout.html**: Layout base com sidebar e painel de insights
- **dashboard.html**: Dashboard com pacientes e sessões recentes
- **patients.html**: Lista completa de pacientes
- **patient.html**: Perfil do paciente com histórico
- **session.html**: Detalhes da sessão com observações e intervenções

## Persistência

### Banco de Dados SQLite
- **Um arquivo por profissional**: Privacidade e portabilidade
- **Schema simples**: Tabelas normalizadas com relações
- **Backup fácil**: Copiar arquivo .db

### Schema do Banco
```sql
-- Tabela de pacientes
CREATE TABLE patients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Tabela de sessões
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATETIME NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Tabela de observações
CREATE TABLE observations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Tabela de intervenções
CREATE TABLE interventions (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Tabela de insights
CREATE TABLE insights (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    source TEXT NOT NULL,
    created_at DATETIME NOT NULL
);
```

## Fluxo de Dados

```
Request HTTP → Handler → Service → Repository → SQLite
                                     ↓
Response HTTP ← Template ← Service ← Domain
```

## Dependências Externas

### Go Modules
- `github.com/mattn/go-sqlite3`: Driver SQLite
- `github.com/google/uuid`: Geração de UUIDs

### Frontend
- **HTMX**: Interatividade sem JavaScript complexo
- **Alpine.js**: Interatividade mínima quando necessário
- **Tailwind CSS**: Utilitários CSS
- **CSS Custom**: Sistema de design Arandu

## Decisões de Design

### 1. Simplicidade Tecnológica
- Go puro (sem frameworks web pesados)
- SQLite (sem servidor de banco de dados)
- HTMX (SPA-like sem complexidade de frontend)

### 2. Privacidade por Design
- Banco local (dados nunca saem da máquina)
- Sem nuvem obrigatória
- Backup controlado pelo usuário

### 3. Extensibilidade
- Domínio independente (pode trocar infraestrutura)
- Interfaces claras (pode adicionar novos repositórios)
- Arquitetura em camadas (pode adicionar novas funcionalidades)

## Próximos Passos Arquiteturais

### Fase 2: Inteligência Assistida
- Serviço de IA como camada de aplicação
- Integração com modelos de linguagem
- Cache de embeddings locais

### Fase 3: Analytics
- Camada de queries complexas
- Indexação full-text
- Análise temporal

### Fase 4: Multi-usuário
- Autenticação e autorização
- Isolamento de dados por profissional
- Compartilhamento seguro (opcional)