# Documentação de Arquitetura - Arandu

Este diretório contém a documentação arquitetural do sistema Arandu.

## Documentos Principais

### 1. [system_structure.md](./system_structure.md)
Visão geral completa da arquitetura do sistema, incluindo:
- Estrutura do projeto
- Princípios arquiteturais (DDD, Clean Architecture, SOLID)
- Modelo de domínio
- Camadas de aplicação e infraestrutura
- Camada web (handlers e templates)
- Persistência (SQLite)
- Fluxo de dados

### 2. [WEB_LAYER_PATTERN.md](./WEB_LAYER_PATTERN.md) ⭐ **NOVO**
Documento de referência detalhado sobre o padrão da camada web:
- **Regras de Ouro** (Independência de Domínio, Consciência HTMX, Tipagem Forte)
- Estrutura modular de templates
- Tratamento de erros contextual
- Injeção de dependência via interfaces
- Fluxo completo de requisições
- Checklist de implementação
- Anti-padrões a evitar

### 3. Decisões de Design
- Simplicidade tecnológica (Go puro, SQLite, HTMX)
- Privacidade por design (dados locais)
- Extensibilidade (interfaces claras, camadas independentes)

## Padrão Web em Resumo

```
┌─────────────────────────────────────────────────────────┐
│                    Request HTTP                         │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Handler (web/handlers/)                                │
│  1. Extrai parâmetros                                   │
│  2. Valida input básico                                 │
│  3. Chama Service                                       │
│  4. Mapeia para ViewModel                               │
│  5. Verifica HX-Request                                 │
│  6. Renderiza template                                  │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Service (application/services/)                        │
│  - Lógica de negócio                                    │
│  - Orquestração de entidades de domínio                 │
│  - Validações complexas                                 │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Repository (infrastructure/repository/)                │
│  - Persistência em SQLite                               │
│  - CRUD operations                                      │
└─────────────────────────────────────────────────────────┘
```

## Regras de Ouro (Resumo)

### ✅ Independência de Domínio
Handlers **NUNCA** contêm lógica de negócio. Apenas orquestram.

### ✅ Consciência HTMX
Sempre verifique `HX-Request` header:
- `true` → Renderiza fragmento (`patient-content`, `session-content`)
- `false` → Renderiza página completa (`layout`)

### ✅ ViewModels Fortemente Tipados
**NUNCA** passe entidades de domínio para templates. Sempre crie structs específicas.

## Templates

Estrutura esperada:
```
web/templates/
├── layout.html           # Esqueleto base
├── error-fragment.html   # Erros HTMX
├── patients.html         # Lista de pacientes
├── patient.html          # Detalhes do paciente
├── session.html          # Detalhes da sessão
└── session_new.html      # Nova sessão
```

Cada template define:
- `{{define "content"}}` - Para full-page rendering
- `{{define "<specific>-content"}}` - Para HTMX fragments

## Atualizações Recentes (Março 2026)

### Fase 1: Consolidação Web ✅ CONCLUÍDA

- [x] Handlers com injeção de dependência via interfaces
- [x] ViewModels fortemente tipados protegendo o domínio
- [x] Consciência HTMX em todos handlers
- [x] Templates modulares com fragments nomeados especificamente
- [x] Tratamento de erros contextual (full-page vs HTMX fragment)

**Arquivos criados/modificados:**
- `internal/web/handlers/patient_handler.go` (novo)
- `internal/web/handlers/session_handler.go` (refatorado)
- `web/templates/error-fragment.html` (novo)
- `web/templates/patient.html` (atualizado)
- `web/templates/patients.html` (atualizado)
- `web/templates/session.html` (atualizado)
- `web/templates/session_new.html` (atualizado)
- `docs/architecture/WEB_LAYER_PATTERN.md` (novo)
- `docs/learnings/task_20260315_174000.md` (novo)

---

**Última atualização:** Março 2026  
**Responsável:** Arquitetura de Sistemas
