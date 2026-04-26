# Plano de Correção - Sistema Arandu

## Status Atual
✅ **Migração de schemas concluída**: Sistema unificado usando migrações SQL
❌ **Handlers web quebrados**: Violações arquiteturais identificadas

## Problemas Críticos Identificados

### 1. `patient_handler.go` - Handler NewPatient
- **Problema**: HTML inline no handler (violação arquitetural)
- **Local**: `internal/web/handlers/patient_handler.go:NewPatient`
- **Solução**: Criar componente `patientComponents.NewPatientForm`

### 2. `session_handler.go` - Handler NewSession
- **Problema**: Usa `ExecuteTemplate` que retorna `nil` (DummyRenderer)
- **Local**: `internal/web/handlers/session_handler.go:NewSession`
- **Solução**: Criar componente `sessionComponents.NewSessionForm`

### 3. `dummy_renderer.go` - Renderer Problemático
- **Problema**: Retorna `nil` silenciosamente, mascarando erros
- **Solução**: Implementar `LoggingRenderer` ou fallback real

## Plano de Ação

### Fase 1: Corrigir Handlers (Prioridade Alta)
1. **Criar componente `patientComponents.NewPatientForm`**
   - Baseado em `web/components/patient/detail.templ`
   - Formulário com campos: nome, data_nascimento, contato, observacoes
   - Usar `layoutComponents.BaseWithContent`

2. **Criar componente `sessionComponents.NewSessionForm`**
   - Baseado em `web/components/session/detail.templ`
   - Formulário com campos: paciente_id, data, tipo, notas
   - Usar `layoutComponents.BaseWithContent`

3. **Atualizar handlers para usar componentes Templ**
   - Remover HTML inline de `patient_handler.go`
   - Remover chamadas `ExecuteTemplate` de `session_handler.go`
   - Usar `templ.Handler` com componentes criados

### Fase 2: Melhorar Renderer (Prioridade Média)
1. **Implementar `LoggingRenderer`**
   - Logar warnings quando `ExecuteTemplate` for chamado
   - Alertar sobre uso de padrão legado
   - Manter compatibilidade temporária

2. **Ou implementar fallback real**
   - Renderizar templates HTML básicos
   - Transição gradual para componentes Templ

### Fase 3: Mecanismos de Proteção (Prioridade Baixa)
1. **Expandir `arandu_start_session.sh`**
   - Adicionar validação de conhecimento arquitetural
   - Verificar handlers problemáticos
   - Alertar sobre violações

2. **Criar `arandu_validate_handlers.sh`**
   - Validação automática de handlers
   - Detectar HTML inline
   - Detectar uso de `ExecuteTemplate`

3. **Criar `ANTI_PATTERNS.md`**
   - Listar violações críticas
   - Exemplos de código problemático
   - Soluções recomendadas

### Fase 4: Estabelecer Processo (Prioridade Baixa)
1. **Criar `arandu_checkpoint.sh`**
   - Revisão obrigatória antes de commits
   - Checklist arquitetural
   - Validação automática

2. **Atualizar templates de task**
   - Incluir checklist arquitetural
   - Referências a documentação
   - Validações obrigatórias

## Arquivos Relevantes

### Componentes Existentes (referência)
- `web/components/layout/layout.templ` - Layout base
- `web/components/patient/detail.templ` - Detalhe do paciente
- `web/components/session/detail.templ` - Detalhe da sessão

### Handlers Problemáticos
- `internal/web/handlers/patient_handler.go`
- `internal/web/handlers/session_handler.go`
- `internal/web/dummy_renderer.go`

### Scripts Existentes
- `scripts/arandu_start_session.sh`
- `scripts/arandu_guard.sh`

### Documentação
- `docs/architecture/AGENT_GUIDE.md`
- `docs/architecture/WEB_LAYER_PATTERN.md`
- `docs/design-system.md`

## Próximos Passos Imediatos
1. Ler `patient_handler.go` para entender HTML inline atual
2. Criar `patientComponents.NewPatientForm`
3. Atualizar `patient_handler.go` para usar componente
4. Testar funcionalidade
5. Repetir para `session_handler.go`

## Notas Importantes
- **Não quebrar funcionalidade existente**
- **Manter compatibilidade com testes**
- **Seguir padrões arquiteturais estabelecidos**
- **Documentar mudanças no AGENT_GUIDE.md**

---
**Criado em**: 16/03/2026
**Última atualização**: 16/03/2026
**Status**: Em progresso