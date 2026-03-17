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
- Contexto do projeto atualizado em `docs/agent_context/project_state.md`

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

**Saída:**
- Tarefa arquivada
- Aprendizado registrado em `docs/learnings/task_YYYYMMDD_HHMMSS.md`
- Contexto do projeto atualizado

---

### 5. **Finalizar Sessão**
```bash
./scripts/arandu_end_session.sh
```
**Descrição:** Finaliza a sessão de trabalho atual.

**Saída:**
- Sessão finalizada
- Resumo das atividades realizadas

---

## 🔄 Fluxo Completo de Exemplo

```bash
# 1. Iniciar sessão
./scripts/arandu_start_session.sh

# 2. Criar tarefa para requirement específico
./scripts/arandu_create_task.sh "Implementar criação de paciente" req-01-00-01

# 3. Processar tarefa (gerar prompt)
./scripts/arandu_process_task.sh 20260313_215938

# 4. Implementar manualmente baseado no prompt
# ... trabalho de implementação ...

# 5. Concluir tarefa com aprendizado
./scripts/arandu_conclude_task.sh 20260313_215938 "Entidade Patient com validação implementada" --success

# 6. Finalizar sessão
./scripts/arandu_end_session.sh
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
