# Visão Geral da Arquitetura - Arandu

**Versão:** 1.1  
**Data:** 22/04/2026  
**Status:** Atualizado

---

## 📐 Visão Arquitetural

O Arandu segue uma arquitetura em camadas inspirada em Clean Architecture e Domain-Driven Design, com ênfase em:

- **Separação de responsabilidades** entre camadas
- **Independência de frameworks** na camada de domínio
- **Testabilidade** através de interfaces e injeção de dependência
- **Multi-tenancy** para isolamento de dados

---

## 🏗️ Diagrama de Camadas

```mermaid
C4Container
    title Arandu - Arquitetura de Container (C4)
    
    Person(terapeuta, "Terapeuta", "Profissional de saúde mental")
    
    Container_Boundary(web, "Camada Web") {
        Container(router, "HTTP Router", "Go/net/http", "Roteamento central")
        Container(middleware, "Middleware", "Go", "Auth, Telemetry, HTMX")
        Container(handlers, "Handlers", "Go", "15 handlers HTTP")
        Container(components, "Components", "Templ", "65+ componentes UI")
    }
    
    Container_Boundary(app, "Camada Aplicação") {
        Container(services, "Services", "Go", "16 services")
        Container(ports, "Ports", "Go", "Interfaces de domínio")
    }
    
    Container_Boundary(domain, "Camada Domínio") {
        Container(entities, "Entities", "Go", "Modelos de domínio")
        Container(interfaces, "Interfaces", "Go", "Repository interfaces")
    }
    
    Container_Boundary(data, "Camada Dados") {
        Container(repos, "Repositories", "Go", "22 repositories")
        Container(migrations, "Migrations", "SQL", "Schema versioning")
        ContainerDb(db, "SQLite", "SQLite", "Dados tenant")
        ContainerDb(central, "Central DB", "SQLite", "Auth, Audit")
    }
    
    Container_Ext(gemini, "Google Gemini", "API Externa", "Geração de insights")
    
    Rel(terapeuta, components, "Usa", "HTTPS")
    Rel(components, handlers, "Renderiza", "Go")
    Rel(handlers, services, "Chama", "Go")
    Rel(handlers, middleware, "Usa", "Go")
    Rel(services, entities, "Manipula", "Go")
    Rel(services, repos, "Persiste", "Go")
    Rel(repos, db, "Query", "SQL")
    Rel(services, gemini, "API", "HTTPS")
```

---

## 📋 Camadas Detalhadas

### 1. Camada Web (Presentation)

**Responsabilidade:** Receber requisições HTTP e renderizar respostas

```
internal/web/
├── handlers/           # 13 handlers
│   ├── patient_handler.go
│   ├── session_handler.go
│   ├── observation_handler.go
│   ├── intervention_handler.go
│   ├── classification_handler.go
│   ├── intervention_classification_handler.go
│   ├── timeline_handler.go
│   ├── analysis_handler.go
│   ├── ai_handler.go
│   ├── auth_handler.go
│   ├── dashboard_handler.go
│   ├── agenda_handler.go
│   └── biopsychosocial_handler.go
└── service_adapters.go  # Adaptadores para injeção de dependência
```

**Componentes UI:**
```
web/components/
├── patient/           # 15 componentes (profile, anamnesis, goals, biopsychosocial...)
├── session/           # 10 componentes (edit_form Sábio, observation_item, intervention_item...)
├── classification/    # 5 componentes (tags para observações)
├── intervention/      # 3 componentes (tags para intervenções)
├── timeline/          # 5 componentes
├── analysis/          # 5 componentes (theme_cloud, patterns)
├── layout/            # 8 componentes (shell_layout, llm_drawer, toast...)
├── agenda/            # 4 componentes (agenda_layout, appointment_detail, new_form)
├── dashboard/         # 2 componentes
├── ai/                # 2 componentes
└── auth/              # 2 componentes
```

**Tecnologias:**
- Go 1.21+
- Templ 0.3.x (templating type-safe)
- HTMX **2.x** (servido localmente em `/static/js/htmx.min.js`)
- Tailwind CSS **v4** (build via `input-v2.css`)
- Alpine.js 3.13.5 (estado mínimo do cliente)
- CSS Sábio design system (`style.css`)

---

### 2. Camada Aplicação (Application)

**Responsabilidade:** Orquestrar casos de uso e regras de negócio

```
internal/application/services/
├── patient_service.go
├── session_service.go
├── create_session_service.go   # caso de uso isolado
├── observation_service.go
├── intervention_service.go
├── timeline_service.go
├── timeline_service_context.go # contexto biopsicossocial na timeline
├── biopsychosocial_service.go
├── goal_service.go
├── insight_service.go
├── ai_service.go
├── audit_service.go
├── tenant_service.go
└── agenda_service.go
```

**Padrão:** Cada service implementa uma interface (Port) definida no handler

```go
// Exemplo de interface (Port)
type PatientService interface {
    GetPatientByID(ctx context.Context, id string) (*patient.Patient, error)
    ListPatients(ctx context.Context) ([]*patient.Patient, error)
    CreatePatient(ctx context.Context, input CreatePatientInput) (*patient.Patient, error)
    // ...
}
```

---

### 3. Camada Domínio (Domain)

**Responsabilidade:** Modelar entidades e regras de domínio

```
internal/domain/
├── patient/
│   ├── patient.go       # Aggregate Root + DashboardSummary
│   ├── goal.go          # Meta terapêutica
│   ├── anamnesis.go     # Anamnese multidimensional
│   ├── medication.go
│   └── vitals.go
├── session/
│   ├── session.go
│   ├── entity.go
│   └── repository.go    # Interface Repository
├── observation/
│   ├── observation.go
│   └── tag.go           # ObservationTag + classificação
├── intervention/
│   ├── intervention.go
│   └── classification.go # InterventionClassification
├── appointment/
│   └── appointment.go   # Entidade de agenda
├── insight/
│   └── insight.go       # Insights gerados por IA
├── timeline/
│   └── timeline.go      # Read model longitudinal
└── shared/
    ├── tenant.go         # Entidade de tenant
    └── user.go           # Entidade de usuário
```

**Princípio:** Camada independente de frameworks e infraestrutura

---

### 4. Camada Infraestrutura (Infrastructure)

**Responsabilidade:** Implementar detalhes técnicos (DB, APIs externas)

```
internal/infrastructure/
├── repository/
│   └── sqlite/
│       ├── patient_repository.go
│       ├── session_repository.go
│       ├── observation_repository.go
│       ├── intervention_repository.go
│       ├── timeline_repository.go
│       ├── goal_repository.go
│       ├── medication_repository.go
│       ├── vitals_repository.go
│       ├── appointment_repository.go
│       ├── insight_repository.go
│       ├── context_wrapper.go   # Extrai tenant DB do context
│       ├── tenant_pool.go       # Pool de conexões por tenant
│       ├── central_db.go        # Control plane DB
│       ├── db.go                # Abstração de conexão
│       └── migrations/          # 0001..0014 arquivos .up.sql
├── ai/
│   ├── gemini_client.go
│   └── cache.go
└── auth/
    └── google_provider.go

internal/platform/       # Utilitários transversais
├── context/             # Injeção de tenant DB no context
├── env/                 # Variáveis de ambiente (IsDev, etc.)
├── helpers/             # CSS versioning hash
├── logger/              # Logger estruturado
├── middleware/          # Auth, telemetry, cache, request_id
└── version/             # Versão do binário
```

---

## 🔄 Fluxo de Dados

```mermaid
sequenceDiagram
    actor User as Terapeuta
    participant UI as Component Templ
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant DB as SQLite
    participant AI as Gemini API
    
    User->>UI: Ação (clique/submit)
    UI->>H: HTTP Request (HTMX)
    H->>H: Validar entrada
    H->>S: Chamar método
    S->>S: Regras de negócio
    S->>R: Persistir
    R->>DB: SQL Query
    DB-->>R: Resultado
    R-->>S: Entidade
    S->>AI: Gerar insight (opcional)
    AI-->>S: Insight
    S-->>H: Resposta
    H->>UI: HTML (HTMX swap)
    UI-->>User: Atualização
```

---

## 🗄️ Modelo de Dados

### Entidades Principais

```mermaid
erDiagram
    PATIENT ||--o{ SESSION : possui
    SESSION ||--o{ OBSERVATION : contem
    SESSION ||--o{ INTERVENTION : contem
    SESSION ||--o{ GOAL : define
    PATIENT ||--o{ MEDICATION : usa
    PATIENT ||--o{ VITALS : registra
    PATIENT ||--o{ ANAMNESIS : possui
    OBSERVATION ||--o{ OBSERVATION_TAG : classificado
    TAG ||--o{ OBSERVATION_TAG : utilizado
    
    PATIENT {
        string id PK
        string name
        string notes
        datetime created_at
        datetime updated_at
    }
    
    SESSION {
        string id PK
        string patient_id FK
        date date
        string summary
        datetime created_at
        datetime updated_at
    }
    
    OBSERVATION {
        string id PK
        string session_id FK
        text content
        datetime created_at
        datetime updated_at
    }
    
    INTERVENTION {
        string id PK
        string session_id FK
        text content
        datetime created_at
        datetime updated_at
    }
    
    TAG {
        string id PK
        string name
        string tag_type
        string color
        int sort_order
    }
    
    OBSERVATION_TAG {
        string id PK
        string observation_id FK
        string tag_id FK
        int intensity
        datetime created_at
    }
    
    GOAL {
        string id PK
        string patient_id FK
        string title
        text description
        string status
        datetime target_date
        datetime created_at
    }
    
    MEDICATION {
        string id PK
        string patient_id FK
        string name
        string dosage
        string status
        datetime start_date
    }
    
    VITALS {
        string id PK
        string patient_id FK
        float weight
        float heart_rate
        float blood_pressure_systolic
        float blood_pressure_diastolic
        datetime recorded_at
    }
```

---

## 🏢 Multi-Tenancy

### Arquitetura

```mermaid
C4Container
    title Multi-Tenancy - Arandu
    
    Person(user, "Usuário", "Terapeuta")
    
    Container_Boundary(app, "Aplicação") {
        Container(middleware, "Tenant Middleware", "Go", "Extrai tenant do contexto")
        Container(pool, "Tenant Pool", "Go", "Gerencia conexões")
    }
    
    ContainerDb(central, "Central DB", "SQLite", "Auth, Users, Audit")
    
    ContainerDb(tenant1, "Tenant DB 1", "SQLite", "Dados clínicos")
    ContainerDb(tenant2, "Tenant DB 2", "SQLite", "Dados clínicos")
    ContainerDb(tenant3, "Tenant DB 3", "SQLite", "Dados clínicos")
    
    Rel(user, middleware, "Requisição com JWT")
    Rel(middleware, central, "Valida usuário")
    Rel(middleware, pool, "Obtém conexão")
    Rel(pool, tenant1, "Conexão isolada", "opcional")
    Rel(pool, tenant2, "Conexão isolada", "opcional")
    Rel(pool, tenant3, "Conexão isolada", "opcional")
```

### Componentes

| Componente | Arquivo | Descrição |
|------------|---------|-----------|
| Tenant Pool | `internal/infrastructure/repository/sqlite/tenant_pool.go` | Pool de conexões |
| Context Wrapper | `internal/infrastructure/repository/sqlite/context_wrapper.go` | Extrai tenant do context |
| Central DB | `internal/infrastructure/repository/sqlite/central_db.go` | DB centralizado |

---

## 🔐 Segurança

### Autenticação

```mermaid
flowchart LR
    User[Usuário] -->|1. Login| Auth[Auth Handler]
    Auth -->|2. Validar| Google[Google OAuth]
    Google -->|3. Token| Auth
    Auth -->|4. JWT| User
    User -->|5. Request + JWT| Middleware[Auth Middleware]
    Middleware -->|6. Validar| Central[Central DB]
    Central -->|7. OK| Middleware
    Middleware -->|8. Context| Handler[Handler]
```

### Componentes
- `internal/web/handlers/auth_handler.go` - Login/logout
- `internal/platform/middleware/auth.go` - JWT validation

---

## 📡 APIs e Integrações

### APIs Internas

| Endpoint | Método | Descrição |
|----------|--------|-----------|
| `/patients` | GET/POST | CRUD pacientes |
| `/session/{id}` | GET/PUT | CRUD sessões |
| `/observations/{id}/classify` | POST | Classificar |
| `/tags` | GET | Listar tags |

### APIs Externas

| Serviço | Uso | Arquivo |
|---------|-----|---------|
| Google Gemini | Insights IA | `internal/infrastructure/ai/gemini_client.go` |
| Google OAuth | Autenticação | `internal/infrastructure/auth/google_provider.go` |

---

## 🧪 Testes

### Estrutura

```
tests/
├── e2e/                    # Testes end-to-end
│   ├── http_patient_flow_test.go
│   └── e2e_full_workflow_test.go
├── integration/              # Testes de integração
└── unit/                   # Testes unitários
```

### Cobertura

| Camada | Cobertura | Status |
|--------|-----------|--------|
| Handlers | - | 🟡 Aumentar |
| Services | - | 🟡 Aumentar |
| Repositories | - | 🟡 Aumentar |

---

## 🚀 Deployment

### Requisitos

- Go 1.21+
- SQLite 3
- Acesso à internet (Gemini API opcional)

### Estrutura de Diretórios

```
arandu/
├── cmd/arandu/            # Entry point
├── internal/              # Código privado
├── web/
│   ├── components/        # Templates Templ
│   └── static/           # Assets
├── storage/              # Dados SQLite
├── migrations/           # Migrações SQL
└── docs/                 # Documentação
```

---

## 📊 Tecnologias

| Camada | Tecnologia | Versão |
|--------|------------|--------|
| Linguagem | Go | 1.21+ |
| Template | Templ | 0.3.x |
| CSS | Tailwind | **v4** |
| CSS Design | Sábio (`style.css`) | — |
| HTMX | HTMX | **2.x** (local) |
| Alpine.js | Alpine | 3.13.5 |
| DB | SQLite | 3.x |
| Auth | Cookie session + OAuth2 | — |
| AI | Google Gemini | API |

---

## 🎯 Princípios de Design

### 1. Clean Architecture

```
Domain (Independente)
    ↑
Application (Regras de uso)
    ↑
Infrastructure (Detalhes)
    ↑
Web (Framework)
```

### 2. Dependency Inversion

```go
// Handler depende de interface, não de implementação
type PatientService interface {
    GetPatientByID(ctx context.Context, id string) (*patient.Patient, error)
}

// Service implementa a interface
type PatientServiceImpl struct {
    repo PatientRepository
}
```

### 3. Multi-Tenancy

- Isolamento por database
- Context propagation
- Connection pooling

---

## 📈 Escalabilidade

### Horizontal
- Stateless handlers
- SQLite por tenant
- Pode migrar para PostgreSQL

### Vertical
- Go routines
- Connection pooling
- Caching (preparado)

---

## 🔗 Links Relacionados

- [Índice de Implementação](./IMPLEMENTATION_INDEX.md)
- [Roadmap](./ROADMAP.md)
- [Documentação de APIs](./architecture/ROUTE_CONVENTIONS.md)
- [Padrões de Layout](./architecture/standardized_layout_protocol.md)

---

## 📅 Histórico

| Data | Versão | Alterações |
|------|--------|------------|
| 04/04/2026 | 1.0 | Criação do documento |
| 22/04/2026 | 1.1 | Stack atualizada (HTMX 2.x, Tailwind v4, Alpine), handlers e domínios adicionados (agenda, insight, shared), platform/ documentada |

---

**Arquitetura mantida por:** Arandu Team  
**Próxima revisão:** Mensal ou em mudanças significativas
