# 🧩 Arandu Scripts Guide

Guia completo de todos os scripts do projeto Arandu.

---

## 📁 Estrutura de Diretórios

```
scripts/
├── arandu_*.sh              # Scripts principais de task/session
├── test_*.sh               # Scripts de teste
├── verify_*.sh             # Scripts de validação
├── safe_deploy.sh          # Deployment seguro
├── start-monitoring.sh     # Monitoramento
├── analise_densidade_layout.sh  # Análise de layout
├── e2e/                    # Testes E2E
│   ├── core.sh
│   ├── config.sh
│   ├── report.sh
│   ├── modules/            # Módulos de teste
│   └── utils/              # Utilitários E2E
├── others/                 # Utilitários diversos
│   ├── create_*.sh         # Criação de docs
│   ├── copia_docs.sh       # Merge de documentação
│   └── sobe_app.sh        # Gerenciamento de porta
└── backup/                # [DEPRECADO - não usar]
```

---

## 🎯 Scripts de Sessão e Tarefa

### 1. `arandu_start_session.sh`
Inicia uma nova sessão de trabalho.

```bash
./scripts/arandu_start_session.sh
```

**O que faz:**
- Cria `work/current_session/`
- Gera `agent_context.md` com contexto obrigatório
- Atualiza contexto do projeto

---

### 2. `arandu_create_task.sh`
Cria uma nova tarefa.

```bash
# Sem requirement
./scripts/arandu_create_task.sh "Título da tarefa"

# Com requirement
./scripts/arandu_create_task.sh "Implementar feature X" req-01-01-01
```

**O que faz:**
- Cria `work/tasks/task_YYYYMMDD_HHMMSS/task.md`
- Se sem REQ: cria template com checklist de integridade
- Se com REQ: cria referência ao requirement

---

### 3. `arandu_process_task.sh`
Valida e exibe tarefa antes de implementar.

```bash
./scripts/arandu_process_task.sh 20260329_185952
```

**O que faz:**
- Verifica se task.md existe e tem conteúdo
- Exibe detalhes da tarefa
- Impede execução se tarefa não foi detalhada

---

### 4. `arandu_conclude_task.sh`
Conclui e arquiva uma tarefa.

```bash
./scripts/arandu_conclude_task.sh 20260329_185952 --success
./scripts/arandu_conclude_task.sh 20260329_185952 --failure
```

**O que faz:**
- Executa `arandu_guard.sh` para validação
- Verifica se templ foi gerado
- Sugere aprendizados para MASTER_LEARNINGS.md
- Move tarefa para `work/archive/`
- Atualiza contexto

---

### 5. `arandu_checkpoint.sh`
Validação arquitetural obrigatória antes de concluir tarefa.

```bash
./scripts/arandu_checkpoint.sh
```

**O que faz:**
1. Verifica estado do git
2. Valida handlers (`arandu_validate_handlers.sh`)
3. Compila handlers e build completo
4. Verifica migrações SQL
5. Detecta anti-padrões (HTML inline, arquivos .html)
6. Resumo com contagem de erros

**Retorna:**
- ❌ vermelho se há erros
- 🟡 amarelo se há avisos
- 🟢 verde se tudo ok

---

### 6. `arandu_guard.sh`
Verificação de integridade do sistema (usado antes de concluir).

```bash
./scripts/arandu_guard.sh
```

**O que faz:**
- Testa rotas (/dashboard, /patients, /patients/new)
- Verifica se templ generation está atualizada
- Cria usuário de teste temporário

---

### 7. `arandu_end_session.sh`
Encerra sessão de trabalho.

```bash
./scripts/arandu_end_session.sh
```

**O que faz:**
- Verifica se há tarefas pendentes
- Gera snapshot final do projeto
- Limpa diretório da sessão

---

## 🧠 Context & Trace

### 8. `arandu_update_context.sh`
Atualiza contexto permanente do projeto.

```bash
./scripts/arandu_update_context.sh
```

**O que faz:**
- Atualiza `docs/project_context/permanent_state.md`
- Gera `docs/agent_context/project_state.md`
- Adiciona tarefas concluídas ao histórico

---

### 9. `arandu_trace.sh`
Sistema de rastreamento de execuções (biblioteca).

**Não é executado diretamente.** É sourceado por outros scripts.

```bash
# Ativar trace
./scripts/arandu_create_task.sh "Minha tarefa" --trace

# Ver traces
./scripts/arandu_trace_view.sh

# Watch em tempo real
./scripts/arandu_trace_view.sh --watch
```

---

### 10. `arandu_trace_view.sh`
Visualiza logs de trace.

```bash
./scripts/arandu_trace_view.sh        # Ver tudo
./scripts/arandu_trace_view.sh --watch  # Watch em tempo real
```

---

## ✅ Testes

### 11. `test_all.sh`
Executa todos os testes de scripts.

```bash
./scripts/test_all.sh
```

**Executa:**
- `test_checkpoint.sh`
- `test_guard.sh`
- `test_trace.sh`

---

### 12. `test_checkpoint.sh`
Testes unitários para checkpoint.

```bash
./scripts/test_checkpoint.sh
```

---

### 13. `test_guard.sh`
Testes unitários para guard.

```bash
./scripts/test_guard.sh
```

---

### 14. `test_trace.sh`
Testes unitários para trace.

```bash
./scripts/test_trace.sh
```

---

## 🔍 Validação

### 15. `arandu_validate_handlers.sh`
Valida handlers web.

```bash
./scripts/arandu_validate_handlers.sh
```

**Verifica:**
- Sem uso de ExecuteTemplate (deprecated)
- Sem import de html/template
- Sem arquivos .html
- Sem HTML inline
- Componentes Templ existem
- Build compila

---

### 16. `verify_css.sh`
Verifica CSS do projeto.

```bash
./scripts/verify_css.sh
```

---

### 17. `verify_onboarding.sh`
Verifica fluxo de onboarding.

```bash
./scripts/verify_onboarding.sh
```

---

### 18. `verify_e2e_setup.sh`
Verifica se E2E está configurado.

```bash
./scripts/verify_e2e_setup.sh
```

---

### 19. `validate_modules.sh`
Valida módulos do projeto.

```bash
./scripts/validate_modules.sh
```

---

## 🎨 Visual & Screenshots

### 20. `arandu_visual_check.sh`
Validação visual obrigatória.

```bash
./scripts/arandu_visual_check.sh
```

**Checklist:**
- [ ] Desktop (1920px)
- [ ] Mobile (375px)
- [ ] Nenhum elemento sobreposto
- [ ] Scroll funciona

---

### 21. `arandu_screenshot.sh`
Gera screenshots para comparação.

```bash
./scripts/arandu_screenshot.sh
```

---

### 22. `test_screenshot_manual.sh`
Teste manual de screenshot.

```bash
./scripts/test_screenshot_manual.sh
```

---

### 23. `analise_densidade_layout.sh`
Analisa densidade de layout HTML.

```bash
./scripts/analise_densidade_layout.sh [dir]
# Default: tmp/audit_logs/
```

---

## 🚀 Deployment & Infraestrutura

### 24. `safe_deploy.sh`
Deploy seguro com validações.

```bash
./scripts/safe_deploy.sh
```

**O que faz:**
1. Para instâncias anteriores
2. Compila Tailwind CSS
3. Gera templ
4. Compila Go
5. Inicia servidor

---

### 25. `start-monitoring.sh`
Inicia monitoramento.

```bash
./scripts/start-monitoring.sh
```

---

### 26. `debug_login.sh`
Debug de login.

```bash
./scripts/debug_login.sh
```

---

## 📄 Documentação (others/)

### 27. `create_req.sh`
Cria requirement.

```bash
./scripts/others/create_req.sh "nome-do-requisito"
```

---

### 28. `create_cap.sh`
Cria capability.

```bash
./scripts/others/create_cap.sh "nome-da-capability"
```

---

### 29. `create_vision.sh`
Cria visão.

```bash
./scripts/others/create_vision.sh "nome-da-visao"
```

---

### 30. `create_struct.sh`
Cria estrutura de documentação.

```bash
./scripts/others/create_struct.sh
```

---

### 31. `copia_docs.sh`
Merge de documentação.

```bash
./scripts/others/copia_docs.sh <dir_origem> <arquivo_saida>
./scripts/others/copia_docs.sh docs/vision visao_consolidada.md
```

---

### 32. `sobe_app.sh`
Gerencia porta e inicia app.

```bash
./scripts/others/sobe_app.sh 8080 "go run cmd/arandu/main.go"
```

---

## 🧪 E2E Testing (e2e/)

### 33. `e2e/core.sh`
Funções core para E2E.

---

### 34. `e2e/config.sh`
Configuração E2E.

---

### 35. `e2e/report.sh`
Gera relatório E2E.

```bash
./scripts/e2e/report.sh
```

---

### 36. `e2e/modules/test_patients.sh`
Testa módulo de pacientes.

```bash
./scripts/e2e/modules/test_patients.sh
```

---

### 37. `e2e/modules/test_dashboard.sh`
Testa dashboard.

```bash
./scripts/e2e/modules/test_dashboard.sh
```

---

### 38. `e2e/modules/test_sessions.sh`
Testa sessões.

```bash
./scripts/e2e/modules/test_sessions.sh
```

---

### 39. `e2e/modules/test_interventions.sh`
Testa intervenções.

```bash
./scripts/e2e/modules/test_interventions.sh
```

---

### 40. `e2e/modules/test_observations.sh`
Testa observações.

```bash
./scripts/e2e/modules/test_observations.sh
```

---

### 41. `e2e/modules/test_public.sh`
Testa páginas públicas.

```bash
./scripts/e2e/modules/test_public.sh
```

---

### 42. `e2e/modules/test_responsive.sh`
Testa responsividade.

```bash
./scripts/e2e/modules/test_responsive.sh
```

---

### 43. `arandu_e2e_audit.sh`
Auditoria E2E.

```bash
./scripts/arandu_e2e_audit.sh
```

---

## ⚠️ One-Time Scripts (Usar com Cautela)

Estes scripts foram usados para migrações/setup e geralmente não devem ser executados novamente:

- `arandu_migrate_to_hashed_storage.sh` - Migração para diretórios hashed
- `arandu_migrate_to_uuid.sh` - Migração para UUID
- `init_design_system.sh` - Setup inicial do design system

---

## 🎯 Fluxo de Trabalho Recomendado

### Iniciar Sessão
```bash
./scripts/arandu_start_session.sh
```

### Criar Tarefa
```bash
./scripts/arandu_create_task.sh "Implementar feature X" req-01-01-01
# Editar work/tasks/task_YYYYMMDD_HHMMSS/task.md
```

### Desenvolver
```bash
# Validar durante desenvolvimento
./scripts/arandu_checkpoint.sh
./scripts/arandu_guard.sh
```

### Concluir Tarefa
```bash
# 1. Visual check obrigatório
./scripts/arandu_visual_check.sh

# 2. Concluir tarefa
./scripts/arandu_conclude_task.sh 20260329_185952 --success
```

### Encerrar Sessão
```bash
./scripts/arandu_end_session.sh
```

---

## 🔗 Scripts Concentradores

Ver: `arandu_workflow.sh` - Executa fluxo completo de validação.
