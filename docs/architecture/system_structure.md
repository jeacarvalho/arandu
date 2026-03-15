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

- **Single Responsibility**: Handlers apenas orquestram, services têm lógica de negócio, templates apenas renderizam
- **Open/Closed**: ViewModels permitem estender dados da view sem modificar domínio
- **Liskov Substitution**: Interfaces de serviço permitem trocar implementações (ex: SQLite → PostgreSQL)
- **Interface Segregation**: Cada handler define interfaces específicas com apenas métodos que usa
- **Dependency Inversion**: Handlers dependem de abstrações (interfaces), não de implementações concretas

### 4. Regras de Ouro da Camada Web

#### Independência de Domínio
Handlers não contêm lógica de negócio. Eles apenas:
1. Decodificam o Request
2. Chamam o Application Service
3. Mapeiam o resultado para um ViewModel
4. Renderizam o Template

#### Consciência de Contexto HTMX
Cada Handler verifica se a requisição é HTMX (`HX-Request`):
- **Se for HTMX**: Renderiza apenas o fragmento (bloco específico)
- **Se não for**: Renderiza a página completa com o `layout.html`

#### Tipagem Forte (ViewModels)
Nunca passe entidades de domínio diretamente para o template. Crie structs específicas de "ViewData" dentro do handler para garantir que o template tenha exatamente o que precisa e nada mais.

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
- **PatientHandler**: Gerencia operações com pacientes (listar, mostrar detalhes, criar)
- **SessionHandler**: Gerencia operações com sessões (mostrar, criar, editar, atualizar)

Características arquiteturais:
- **Injeção de Dependência**: Handlers recebem serviços via interfaces (Clean Architecture)
- **Handlers Finos**: Apenas orquestram (recebem request → chamam service → retornam response)
- **ViewModels Fortemente Tipados**: Entidades de domínio nunca são passadas diretamente para templates
- **Consciência HTMX**: Verificam header `HX-Request` para decidir entre fragmento ou página completa
- **Tratamento de Erros Contextual**: Retornam fragmentos de erro amigáveis para requisições HTMX

Padrão de implementação:
```go
// 1. Extração de Parâmetros
id := chi.URLParam(r, "id")

// 2. Chamada ao Serviço (DDD Application Layer)
patient, err := h.service.GetPatient(r.Context(), id)

// 3. Mapeamento para ViewModel (Protege o Domínio)
data := PatientViewData{
    Patient: patient,
    Insights: h.getInsights(r.Context(), id),
}

// 4. Renderização Inteligente (Full Page vs HTMX Fragment)
if r.Header.Get("HX-Request") == "true" {
    h.templates.ExecuteTemplate(w, "patient-content", data) // Só o miolo
} else {
    h.templates.ExecuteTemplate(w, "layout", data) // Layout completo + miolo
}
```

### Templates (`web/templates`)
- **layout.html**: Esqueleto base com sidebar e painel de insights (contém `{{block "content" .}}`)
- **patient.html**: Define `{{define "content"}}` (full-page) e `{{define "patient-content"}}` (fragmento HTMX)
- **patients.html**: Define `{{define "content"}}` e `{{define "patients-content"}}`
- **session.html**: Define `{{define "content"}}` e `{{define "session-content"}}`
- **session_new.html**: Define `{{define "content"}}` e `{{define "new-session-form"}}`

Princípios:
- **DRY (Don't Repeat Yourself)**: Uso de `{{template}}` para evitar duplicação
- **Nomes Únicos**: Fragmentos com nomes específicos para evitar conflitos
- **Separação Clara**: Layout é esqueleto, fragments são conteúdo injetável

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

### Fase 1: Consolidação Web (Concluída)
- ✅ Handlers com injeção de dependência via interfaces
- ✅ ViewModels fortemente tipados protegendo o domínio
- ✅ Consciência HTMX em todos handlers
- ✅ Templates modulares com fragments nomeados especificamente
- ✅ Tratamento de erros contextual (full-page vs HTMX fragment)

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