# 🛡️ Plano de Proteção para Agentes - Arandu

**Data**: 16/03/2026  
**Status**: Em implementação  
**Objetivo**: Criar processo onde qualquer agente, ao começar a trabalhar no sistema, coloque em seu contexto tudo o que precisa para não gerar novos erros ou quebrar o que já funciona.

## 🎯 Objetivos Baseados nas Respostas

### 1. Nível de Proteção Priorizar
**PROCESSO DE CONTEXTO COMPLETO** - Expandir sistema existente de `arandu_start_session.sh` para incluir:
- ✅ Checklist arquitetural obrigatório
- ✅ Validação de conhecimento prévio  
- ✅ Testes de verificação automática
- ✅ Contexto crítico de "NUNCA FAÇA"

### 2. Rigor do Bloqueio
**REVISÃO OBRIGATÓRIA (arquiteto)**:
- Pré-commit hooks que **alertam** mas não bloqueiam
- Scripts de validação que **exigem aprovação** do arquiteto
- Sistema de **checkpoints** onde agente deve parar e validar

### 3. Recursos Disponíveis
**PROCESSO SIMPLE**:
- Desenvolvimento apenas por arquiteto + agentes
- Incremental, começando com scripts simples
- Aproveitar scripts existentes

### 4. Timeline
**CORREÇÃO IMEDIATA + MECANISMOS AGORA**:
1. **Fase 1 (Imediato)**: Corrigir problemas atuais (new patient, new session)
2. **Fase 2 (Agora)**: Implementar mecanismos de proteção usando scripts existentes
3. **Fase 3 (Contínuo)**: Expandir para processos mais robustos

## 📋 Plano Detalhado de Implementação

### FASE 1: CORREÇÃO DOS PROBLEMAS ATUAIS (Imediato)

#### 1.1 Corrigir `NewPatient` Handler
- **Problema**: HTML inline, CSS incorreto, sem layout consistente
- **Solução**: Criar componente `patientComponents.NewPatientForm` e usar `layoutComponents.BaseWithContent`
- **Ações**:
  - Criar `web/components/patient/new_form.templ`
  - Atualizar `patient_handler.go` para usar componente Templ
  - Remover HTML inline

#### 1.2 Corrigir `NewSession` Handler  
- **Problema**: Usa `ExecuteTemplate` que retorna `nil` (DummyRenderer)
- **Solução**: Criar componente `sessionComponents.NewSessionForm`
- **Ações**:
  - Criar `web/components/session/new_form.templ`
  - Atualizar `session_handler.go` para usar componente Templ
  - Remover chamadas a `ExecuteTemplate`

#### 1.3 Corrigir `DummyRenderer`
- **Problema**: Retorna `nil` silenciosamente
- **Solução**: Implementar `LoggingRenderer` que:
  - Loga warning quando `ExecuteTemplate` é chamado
  - Sugere implementação correta com Templ
  - Ou implementa fallback real para compatibilidade
- **Ação**: Criar `internal/web/logging_renderer.go`

### FASE 2: MECANISMOS DE PROTEÇÃO (Agora)

#### 2.1 Expandir `arandu_start_session.sh`
- Adicionar validação de conhecimento arquitetural
- Teste de múltipla escolha sobre padrões
- Checklist arquitetural obrigatório

#### 2.2 Criar `scripts/arandu_validate_handlers.sh`
- Validar automaticamente se handlers seguem padrões
- Verificar HTML inline, CSS correto, uso de Templ
- Bloquear violações críticas

#### 2.3 Criar `scripts/arandu_checkpoint.sh`
- Checkpoints obrigatórios onde agente deve parar e validar com arquiteto
- Executar validações automáticas
- Requer aprovação do arquiteto

#### 2.4 Atualizar `scripts/arandu_guard.sh`
- Expandir para verificar consistência visual
- Verificar componentes obrigatórios (sidebar, main-content, CSS)
- Testar todas rotas críticas

### FASE 3: PROCESSO DE TRABALHO COMPLETO

#### 3.1 Fluxo de Trabalho para Agentes
```
1. INICIAR SESSÃO → bash scripts/arandu_start_session.sh
2. IMPLEMENTAR → Seguir AGENT_GUIDE.md, copiar padrões
3. CHECKPOINT → bash scripts/arandu_checkpoint.sh "task"
4. TESTAR → bash scripts/arandu_guard.sh, testes manuais
5. REVISÃO DO ARQUITETO → Apresentar mudanças
6. CONCLUIR → bash scripts/arandu_conclude_task.sh
```

#### 3.2 Documentação de "NUNCA FAÇA" Expandida
- Criar `docs/architecture/ANTI_PATTERNS.md`
- Listar violações críticas (bloqueantes) e de alerta
- Especificar consequências

#### 3.3 Template de Task com Checklist
- Atualizar `scripts/arandu_create_task.sh`
- Incluir checklist arquitetural obrigatório
- Template para aprovação do arquiteto

### FASE 4: IMPLEMENTAÇÃO GRADUAL

#### 4.1 Semana 1: Correções + Mecanismos Básicos
- [ ] Corrigir `NewPatient` handler
- [ ] Corrigir `NewSession` handler  
- [ ] Atualizar `arandu_start_session.sh` com validação
- [ ] Criar `arandu_validate_handlers.sh`

#### 4.2 Semana 2: Processo Estabelecido
- [ ] Implementar `arandu_checkpoint.sh`
- [ ] Expandir `arandu_guard.sh` para validação visual
- [ ] Criar `ANTI_PATTERNS.md`
- [ ] Treinar agentes no novo processo

#### 4.3 Semana 3: Automação e Monitoramento
- [ ] Pre-commit hooks (alerta apenas)
- [ ] Dashboard de saúde do sistema
- [ ] Métricas de conformidade arquitetural
- [ ] Sistema de aprendizado contínuo

## 🎯 Benefícios

### Para Agentes:
- Contexto completo desde o início
- Checklist claro do que fazer/não fazer
- Checkpoints que previnem erros graves
- Padrões consistentes para copiar

### Para Arquiteto:
- Revisões focadas em padrões, não em bugs básicos
- Sistema auto-documentado
- Redução drástica de regressões
- Agentes mais independentes e confiáveis

### Para o Sistema:
- Consistência arquitetural preservada
- Qualidade visual mantida
- Funcionalidades existentes protegidas
- Base sólida para crescimento

## 🚀 Próximos Passos Imediatos

1. **Corrigir `NewPatient` handler** - Criar componente Templ
2. **Corrigir `NewSession` handler** - Criar componente Templ  
3. **Implementar `LoggingRenderer`** - Para compatibilidade durante transição
4. **Expandir `arandu_start_session.sh`** - Com validação de conhecimento

---

**Última atualização**: 16/03/2026  
**Status**: Em implementação (Fase 1 iniciada)