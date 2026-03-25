# Scripts de Automação

Scripts utilizados para gerenciar sessões, tarefas e aprendizados no projeto Arandu.

## 📋 Fluxo de Trabalho Padrão

### 1. **Iniciar Sessão**
```bash
./scripts/arandu_start_session.sh
```
**Descrição:** Cria uma nova sessão de trabalho e atualiza o contexto do projeto.

**Saída:**
- Sessão criada com timestamp
- Contexto do projeto atualizado em `docs/agent_context/project_state.mTask: Implementação de Infraestrutura FTS5 para Análise Clínica
ID da Tarefa: task_20260317_infra_fts5

Requirement Relacionado: REQ-04-01-01

Stack Técnica: SQLite (FTS5), Go, Migrations SQL.

🎯 Objetivod`

---

### 2. **Criar Tarefa**
```bash
./scripts/arandu_create_task.sh "Título da tarefa" [REQ-ID]
```
**Parâmetros:**
- `"Título da tarefa"`: Descrição da tarefa (obrigatório)
- `[REQ-ID]`: ID do requirement relacionado (opcional, ex: `req-01-00-01`)

**Exemplos:**
```bash
# Tarefa genérica
./scripts/arandu_create_task.sh "Implementar feature X"

# Tarefa vinculada a requirement
./scripts/arandu_create_task.sh "Criar entidade Patient" req-01-00-01
```

**Saída:**
- Tarefa criada em `work/tasks/task_YYYYMMDD_HHMMSS/`
- Arquivo `task.md` com detalhes da tarefa

---

### 3. **Processar Tarefa**
```bash
./scripts/arandu_process_task.sh TASK_ID
```
**Parâmetros:**
- `TASK_ID`: ID da tarefa (ex: `20260313_215938`)

**Descrição:** Verifica se a tarefa está pronta para implementação e exibe os detalhes.

**Comportamento:**
1. Verifica se o arquivo `task.md` existe
2. Verifica se `task.md` foi preenchido além do template básico (mais de 13 linhas)
3. Se não estiver preenchido, exige que o usuário complete o arquivo primeiro
4. Se estiver preenchido, exibe o conteúdo para revisão

**Saída:**
- Status da verificação
- Conteúdo do `task.md` se estiver completo
- Mensagem de erro se precisar ser preenchido

**Nota:** A flag `--execute` está planejada para execução automática futura.

---

### 4. **Concluir Tarefa**
```bash
./scripts/arandu_conclude_task.sh TASK_ID "Descrição do aprendizado" [--success|--failure]
```
**Parâmetros:**
- `TASK_ID`: ID da tarefa (ex: `20260313_215938`)
- `"Descrição do aprendizado"`: Resumo do que foi aprendido (obrigatório)
- `[--success|--failure]`: Status da tarefa (opcional, padrão: `--success`)

**Exemplos:**
```bash
# Concluir com sucesso
./scripts/arandu_conclude_task.sh tid "" --success

# Concluir com falha
./scripts/arandu_conclude_task.sh 20260313_215938 "Erro na implementação da validação" --failure
```

**Saída:**# TASK 20260323_203720

Title: Refatoração do Layout Unificado (SLP)

## Status
**AGUARDANDO_DETALHES_DO_USUARIO** - NÃO inicie trabalho até que o usuário edite este arquivo

## Objetivo

Refatoração do Layout Unificado (SLP)

## Descrição

Tarefa criada sem requirement específico.

## Instruções para o Agente

**CRÍTICO: NÃO leia, edite ou execute qualquer ação nesta tarefa até que o usuário tenha editado este arquivo com os detalhes completos.**

1. Aguarde o usuário fornecer detalhes da tarefa neste arquivo
2. Quando o usuário editar este arquivo com os detalhes completos (removendo esta seção de instruções), inicie a implementação
3. Siga o padrão de referenciar requirements quando aplicável

**Verificação obrigatória antes de iniciar:**
- Esta seção "Instruções para o Agente" deve ter sido removida/replaceada pelo usuário
- O arquivo deve conter uma descrição detalhada da tarefa fornecida pelo usuário
- O status deve ter sido atualizado para "PRONTO_PARA_IMPLEMENTACAO" ou similar

## Checklist de Integridade (OBRIGATÓRIO)
- [ ] O componente usa .templ e herda de Layout?
- [ ] A tipografia Source Serif 4 foi aplicada ao conteúdo clínico?
- [ ] Executei 'templ generate' e o código Go compilou?
- [ ] Testei a rota atual e as rotas vizinhas (Regressão)?
- [ ] O banco de dados foi atualizado via migration .up.sql?


- Tarefa arquivada
- Aprendizado registrado em `docs/learnings/task_YYYYMMDD_HHMMSS.md`
- Contexto do projeto atualizado

---

### 5. **Finalizar Sessão**
run ./scripts/safe_deploy.sh; run ./scripts/arandu_checkpoint.sh; run ./scripts/arandu_guard.sh ; run ./scripts/arandu_validate_handlers.sh; run ./scripts/arandu_update_context.sh


```bash
run ./scripts/arandu_end_session.sh
```
**Descrição:** Finaliza a sessão de trabalho atual.

**Saída:**
- Sessão finalizada
- Resumo das atividades realizadas

---

## 🔄 Fluxo Completo de Exemplo

run ./scripts/safe_deploy.sh
```bash
# 1. Iniciar sessão
run ./scripts/arandu_start_session.sh

# 2. Criar tarefa para requirement específico
run ./scripts/arandu_create_task.sh "Implementar criação de paciente" req-01-00-01

# 3. Processar tarefa (gerar prompt)
run ./scripts/arandu_process_task.sh 20260313_215938

# 4. Implementar manualmente baseado no prompt
# ... trabalho de implementação ...

# 5. Concluir tarefa com aprendizado
run ./scripts/arandu_conclude_task.sh 20260313_215938 "Entidade Patient com validação implementada" --success

# 6. Finalizar sessão
run ./scripts/arandu_end_session.sh
```

---

## 📁 Estrutura de Diretórios

```
scripts/
├── arandu_start_session.sh      # Inicia sessão
├── arandu_create_task.sh        # Cria tarefa
├── arandu_process_task.sh       # Processa tarefa (gera prompt)
├── arandu_conclude_task.sh      # Conclui tarefa
├── arandu_end_session.sh        # Finaliza sessão
├── arandu_update_context.sh     # Atualiza contexto (uso interno)
└── README.md                    # Esta documentação

work/tasks/
└── task_YYYYMMDD_HHMMSS/        # Diretório da tarefa
    ├── task.md                  # Detalhes da tarefa
    └── implementation.md        # Documentação da implementação (opcional)
20260320_153638
docs/
├── agent_context/
│   └── project_state.md         # Estado atual do projeto
└── learnings/
    └── task_YYYYMMDD_HHMMSS.md  # Aprendizados registrados
```

---

## 🎯 Boas Práticas

1. **Sempre iniciar com sessão:** Mantém o contexto atualizado
2. **Vincular tarefas a requirements:** Facilita rastreabilidade
3. **Documentar aprendizados:** Registra conhecimento adquirido
4. **Usar status apropriado:** `--success` ou `--failure` para transparência
5. **Manter implementação documentada:** Criar `implementation.md` quando relevante

---

## ⚙️ Scripts de Suporte

### `arandu_update_context.sh`
**Uso interno:** Atualiza o contexto do projeto automaticamente após operações.

### `init_design_system.sh`
**Configuração inicial:** Inicializa o sistema de design do projeto.

---

## 📊 Monitoramento

O estado do projeto é mantido em `docs/agent_context/project_state.md` e inclui:
- Visões ativas
- Capabilities implementadas
- Requirements pendentes
- Tarefas recentes
- Aprendizados registrados
- Status do sistema
