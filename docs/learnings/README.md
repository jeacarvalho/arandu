# 📚 Aprendizados do Projeto Arandu

**Sistema Refatorado - Março 2026**

> ⚠️ **NOTA DE MIGRAÇÃO:** Este diretório foi refatorado para consolidar ~45 arquivos individuais em um sistema organizado de documentação.

---

## 🎯 Objetivo

Este diretório contém **aprendizados permanentes** coletados durante o desenvolvimento do projeto Arandu. Os aprendizados são organizados por tema para facilitar consulta e manutenção.

---

## 📁 Estrutura Atual

```
docs/learnings/
├── MASTER_LEARNINGS.md          # 📚 Arquivo principal consolidado
├── ARCHITECTURE_PATTERNS.md     # 🏗️ Padrões arquiteturais detalhados
├── TEMPL_GUIDE.md              # 🎨 Guia específico para Templ
├── SQLITE_BEST_PRACTICES.md    # 💾 Práticas para SQLite/FTS5
├── archive/                    # 📦 Arquivos originais (backup)
│   ├── task_20260315_174000.md
│   ├── REQ-01-02-02.md
│   └── ...
└── README.md                   # 📖 Este arquivo
```

---

## 📚 Arquivos Principais

### 1. MASTER_LEARNINGS.md
**Arquivo principal** que consolida todos os aprendizados valiosos organizados por categoria:
- 🏗️ Arquitetura Web (Go + templ + HTMX)
- 💾 Banco de Dados (SQLite, FTS5, Migrations)
- 🩺 Domínio Clínico e Validações
- 🎨 UI/UX e Design System
- 🧪 Testes e Qualidade
- 🔄 Fluxo de Trabalho e Scripts
- 🚨 Anti-Padrões e Erros Comuns

### 2. ARCHITECTURE_PATTERNS.md
**Guia detalhado** de padrões arquiteturais:
- Estrutura de camadas (Clean Architecture)
- Handlers e ViewModels
- HTMX e renderização contextual
- Injeção de dependência
- Tratamento de erros
- Componentes Templ

### 3. TEMPL_GUIDE.md
**Guia completo** para o framework Templ:
- Instalação e setup
- Estrutura de componentes
- Sintaxe e funcionalidades
- Erros comuns e soluções
- Migração de html/template
- Boas práticas

### 4. SQLITE_BEST_PRACTICES.md
**Melhores práticas** para SQLite:
- Driver e configuração
- Sistema de migrations
- Full-Text Search (FTS5)
- Performance e otimização
- Transações e concorrência
- Backup e recuperação

---

## 🔄 Como Usar

### Para Desenvolvedores/Agentes

1. **Iniciar sessão:** Leia `MASTER_LEARNINGS.md` para contexto geral
2. **Tema específico:** Consulte arquivo especializado se necessário
3. **Referência cruzada:** Use links entre arquivos para detalhes

### Para Adicionar Novos Aprendizados

1. **Avalie se é valioso:** O aprendizado é específico, útil e não repetitivo?
2. **Escolha o arquivo:**
   - Geral → `MASTER_LEARNINGS.md`
   - Arquitetura → `ARCHITECTURE_PATTERNS.md`
   - Templ → `TEMPL_GUIDE.md`
   - SQLite → `SQLITE_BEST_PRACTICES.md`
3. **Siga o formato:** Use a estrutura estabelecida no arquivo

### Formato Sugerido para Novos Aprendizados

```markdown
### Título do Aprendizado

**Contexto:** [Breve descrição do contexto]

**Problema:** [O que deu errado ou poderia ser melhor]

**Solução:** [Como foi resolvido ou melhorado]

**Código de exemplo (se aplicável):**
```go
// Código relevante
```

**Referência:** [Tarefa ou requirement relacionado]
```

---

## 📊 Histórico da Refatoração

### Situação Anterior (Problemática)
- **45 arquivos** individuais em `docs/learnings/`
- **18 arquivos** com conteúdo idêntico/repetitivo
- **Conteúdo desorganizado** e difícil de encontrar
- **Scripts geravam lixo** automático (`arandu_conclude_task.sh`)

### Mudanças Implementadas

1. **Consolidação:** 45 arquivos → 4 arquivos principais + archive
2. **Organização:** Conteúdo categorizado por tema
3. **Qualidade:** Remoção de conteúdo repetitivo e de baixo valor
4. **Manutenibilidade:** Sistema fácil de atualizar e expandir

### Arquivos Deletados/Arquivados
- **18 arquivos** com conteúdo repetitivo sobre "Conflito de templates"
- **Arquivos com menos de 20 linhas** e conteúdo genérico
- **Arquivos problemáticos**: `.md`, `Teste.md`

### Arquivos Preservados (em archive/)
- `task_20260315_174000.md` (207 linhas) - Refatoração arquitetura
- `REQ-01-02-02.md` (121 linhas) - Problemas com Templ
- `logbook-sota.md` (113 linhas) - Consolidação anterior
- Outros arquivos com conteúdo valioso

---

## 🔗 Integração com Scripts

### Scripts Atualizados

1. **`arandu_start_session.sh`** (linha 30)
   - **Antes:** "5 docs/learnings/" (instruía ler pasta inteira)
   - **Depois:** Instrui ler `MASTER_LEARNINGS.md`

2. **`arandu_conclude_task.sh`**
   - **Antes:** Gerava conteúdo repetitivo automático
   - **Depois:** Sugere adicionar aprendizado valioso aos arquivos principais

3. **`arandu_update_context.sh`** e **`_v2.sh`**
   - **Antes:** Listava 5 arquivos mais recentes
   - **Depois:** Referencia `MASTER_LEARNINGS.md`

### Compatibilidade
- **Archive mantido:** Arquivos originais em `docs/learnings/archive/`
- **Scripts funcionam:** Todos atualizados para novo sistema
- **Histórico preservado:** Git mantém histórico completo

---

## 🎯 Benefícios do Novo Sistema

### Para Desenvolvedores/Agentes
1. **Encontra rápido:** Conteúdo organizado por tema
2. **Qualidade superior:** Aprendizados reais e valiosos
3. **Contexto rico:** Exemplos de código e referências
4. **Atualizado:** Conteúdo revisado e consolidado

### Para Manutenção do Projeto
1. **Fácil atualizar:** Apenas 4 arquivos principais
2. **Sem duplicação:** Conteúdo centralizado
3. **Escalável:** Adicionar novos temas é simples
4. **Consistente:** Formato padronizado

### Para Qualidade do Código
1. **Prevenção de erros:** Anti-padrões documentados
2. **Padrões estabelecidos:** Boas práticas consolidadas
3. **Referência rápida:** Soluções para problemas comuns

---

## 📈 Estatísticas Pós-Refatoração

### Quantitativas
- **Arquivos:** 45 → 4 principais + archive
- **Redução:** ~90% menos arquivos para gerenciar
- **Linhas:** Conteúdo consolidado e organizado

### Qualitativas
- **Organização:** Conteúdo categorizado por tema
- **Acessibilidade:** Fácil encontrar o que precisa
- **Manutenibilidade:** Fácil atualizar e expandir
- **Qualidade:** Conteúdo curado e valioso

---

## ❓ FAQ

### Onde estão os arquivos antigos?
Todos os arquivos originais foram movidos para `docs/learnings/archive/` para referência histórica.

### Como adicionar um novo aprendizado?
1. Avalie se é valioso e não repetitivo
2. Adicione à seção apropriada de `MASTER_LEARNINGS.md`
3. Ou crie nova seção se for um tema novo

### E se precisar do conteúdo original?
Todos os arquivos originais estão em `archive/` com timestamps preservados.

### Os scripts ainda funcionam?
Sim, todos os scripts foram atualizados para o novo sistema.

### Como consultar aprendizados sobre um tema específico?
Use o índice no `MASTER_LEARNINGS.md` ou consulte o arquivo especializado:
- Arquitetura → `ARCHITECTURE_PATTERNS.md`
- Templ → `TEMPL_GUIDE.md`
- SQLite → `SQLITE_BEST_PRACTICES.md`

---

## 🔄 Manutenção Futura

### Revisões Periódicas
1. **Trimestralmente:** Revisar conteúdo por obsolescência
2. **Após grandes mudanças:** Atualizar com novos aprendizados
3. **Sempre:** Remover conteúdo desatualizado

### Expansão do Sistema
1. **Novos temas:** Criar arquivos especializados quando necessário
2. **Melhor organização:** Refinar categorias conforme projeto evolui
3. **Integração:** Manter sincronia com documentação de arquitetura

### Contribuição
1. **Todos podem contribuir:** Desenvolvedores, agentes, usuários
2. **Formato padrão:** Seguir estrutura estabelecida
3. **Qualidade sobre quantidade:** Apenas aprendizados valiosos

---

*Sistema refatorado em Março de 2026 para melhorar organização, acessibilidade e manutenibilidade dos aprendizados do projeto Arandu.*