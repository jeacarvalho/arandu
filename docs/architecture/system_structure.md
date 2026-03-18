Arquitetura do Sistema Arandu (SOTA)

IMPORTANTE: Este é o documento mestre de arquitetura. Qualquer agente deve seguir estes padrões para garantir a "Tecnologia Silenciosa", a integridade do sistema em larga escala e a privacidade absoluta dos dados.

1. Visão Geral

O Arandu é um sistema de inteligência clínica para profissionais de saúde mental, construído com uma arquitetura baseada em Domain-Driven Design (DDD), Clean Architecture e renderização Type-Safe. O sistema é projetado para ser resiliente a grandes massas de dados (Big Data Clínico).

2. Estrutura do Projeto (Consolidada)

arandu/
├── cmd/arandu/
│   └── main.go                 # Ponto de entrada da aplicação
├── internal/
│   ├── domain/                 # Camada de domínio (Core Business - Puro)
│   │   ├── patient/           # Entidade Paciente (Agregado Raiz)
│   │   │   ├── patient.go     # Domínio Patient
│   │   │   ├── context.go     # Domínio Contexto Biopsicossocial (SOTA)
│   │   │   ├── medication.go  # Domínio Medication (REQ-01-04-01)
│   │   │   └── vitals.go      # Domínio Vitals (REQ-01-04-01)
│   │   ├── session/           # Entidade Sessão
│   │   ├── observation/       # Entidade Observação
│   │   ├── intervention/      # Entidade Intervenção
│   │   └── timeline/          # Agregadores de visualização longitudinal
│   ├── application/           # Camada de aplicação (Services)
│   │   └── services/          # Orquestração de casos de uso
│   │       ├── patient_service.go
│   │       ├── session_service.go
│   │       ├── biopsychosocial_service.go  # REQ-01-04-01
│   │       └── ...
│   └── infrastructure/        # Camada de infraestrutura
│       └── repository/
│           └── sqlite/        # Implementações SQLite (FTS5 habilitado)
│               ├── medication_repository.go
│               ├── vitals_repository.go
│               ├── migrations/
│               │   ├── 0001_initial_schema.up.sql
│               │   ├── 0002_enable_fts5.up.sql
│               │   ├── 0003_add_biopsychosocial_tables.up.sql
│               │   └── migrator.go
├── web/                       # Camada Web SOTA
│   ├── handlers/              # Handlers HTTP (Consciência HTMX)
│   │   ├── biopsychosocial_handler.go
│   │   └── ...
│   ├── components/patient/    # Componentes .templ (Atómicos)
│   │   ├── biopsychosocial_panel.templ
│   │   ├── medication_list.templ
│   │   └── vitals_widget.templ
│   └── static/                # Arquivos estáticos (Tailwind, Inter, Serif)
├── scripts/                   # Automação, Seed e Guard (Cão de Guarda)
└── docs/                      # Requisitos, Visões e Estratégias


3. Pilares Arquiteturais SOTA

A. Persistência Declarativa e Isolamento (Multi-tenancy)

Regra: NUNCA crie tabelas via código Go hardcoded. Use migrations em internal/infrastructure/repository/sqlite/migrations/.

Arquitetura Multi-tenant: Database-per-tenant.

Control Plane: Banco central gerencia users e tenants.

Data Plane: Um arquivo SQLite individual por terapeuta (clinical_uuid.db).

Migrações: O Migrator utiliza go:embed e aplica os schemas em lote para todos os tenants no startup.

B. Type-Safe UI (templ + HTMX)

Regra: Proibido o uso de arquivos .html soltos. Use componentes .templ.

Dualidade Tipográfica:

Interface (UI): Fonte Inter (Sans-serif).

Conteúdo Clínico: Fonte Source Serif 4 (Serif).

Consciência de Contexto: Handlers retornam fragmentos específicos para HX-Request ou a página completa via templates.Layout().

C. Performance e Recuperação (Big Data Clínico)

SQLite FTS5: Implementação de busca textual completa para localização instantânea em milhares de notas.

Paginação / Infinite Scroll: Uso de hx-trigger="revealed" para carregar dados sob demanda, evitando sobrecarga do DOM.

Search Delay: Uso de delay:500ms em buscas ativas para reduzir IO no banco.

4. Modelo de Domínio e Schema (Actualizado)

O banco de dados utiliza UUIDs e relações estritas com ON DELETE CASCADE.

-- Patients: Agregado Raiz (Administrativo)
CREATE TABLE IF NOT EXISTS patients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    notes TEXT, -- Notas narrativas em Source Serif 4
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Patient Context: Identidade Biopsicossocial (SOTA - CAP-01-04)
CREATE TABLE IF NOT EXISTS patient_context (
    patient_id TEXT PRIMARY KEY,
    ethnicity TEXT,                -- IBGE / Padrão Saúde
    gender_identity TEXT,
    sexual_orientation TEXT,
    occupation TEXT,
    education_level TEXT,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Sessions: Dependente de Patient
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATETIME NOT NULL,
    summary TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Observations & Interventions: Unidades Atómicas (FTS5 Indexed)
CREATE TABLE IF NOT EXISTS observations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- FTS5 Virtual Table (Busca Instantânea)
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(content, content='observations', content_rowid='id');

-- Patient Medications: Histórico Farmacológico (REQ-01-04-01)
CREATE TABLE IF NOT EXISTS patient_medications (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    name TEXT NOT NULL,
    dosage TEXT,
    frequency TEXT,
    prescriber TEXT,
    status TEXT DEFAULT 'active' CHECK(status IN ('active', 'suspended', 'finished')),
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Patient Vitals: Sinais Vitais e Hábitos (REQ-01-04-01)
CREATE TABLE IF NOT EXISTS patient_vitals (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATE NOT NULL,
    sleep_hours REAL CHECK(sleep_hours >= 0 AND sleep_hours <= 24),
    appetite_level INTEGER CHECK(appetite_level >= 1 AND appetite_level <= 10),
    weight REAL,
    physical_activity INTEGER DEFAULT 0,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);


5. Regras de Ouro da Camada Web

Independência de Domínio

Os Handlers apenas orquestram:

Extração/Decodificação de Parâmetros.

Chamada ao Serviço de Aplicação (DDD).

Mapeamento para ViewModels (Protege o domínio de tags de template).

Renderização via templ.

Estética "Silent UI" (Anti-Fadiga)

Fundo Global: bg-[#F7F8FA] (Cinza papel).

Silent Input: Inputs sem bordas pesadas, focando na escrita fluida.

Hierarquia: Texto clínico sempre em Source Serif 4, text-xl, leading-relaxed.

6. Protocolo Anti-Regressão (Segurança)

Check de Layout: Toda página completa DEVE herdar de templates.Layout().

Check de Compilação: Executar templ generate antes de qualquer teste.

Arandu Guard: Validar rotas principais via scripts/arandu_guard.sh.

E2E Testing: Toda funcionalidade crítica deve possuir um teste Playwright que valide o fluxo HTMX.

7. Decisões de Design

Simplicidade: Go puro e SQLite.

Privacidade: Isolamento físico de bancos de dados por utilizador.

Soberania do Dado: Facilidade de exportação do arquivo .db individual.

8. Evolução (Fases)

Fase 1: Consolidação Web e DDD (Concluída).

Fase 2: Multi-tenancy e Gestão de Acesso (Em andamento).

Fase 3: Escalabilidade, FTS5 e Navegação em Larga Escala (Concluindo).

Fase 4: Contexto Biopsicossocial e Identidade SOTA (Em implementação).

Fase 5: Inteligência Assistida (Vision-05).