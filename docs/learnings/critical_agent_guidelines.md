# Diretrizes Críticas para Agentes - ARANDU v0.9.0

**⚠️ LEIA ANTES DE QUALQUER MODIFICAÇÃO ⚠️**

## 🚨 PRINCÍPIO FUNDAMENTAL

**O SISTEMA ARANDU v0.9.0 ESTÁ FUNCIONAL E ESTÁVEL.**
**NÃO ALTERE SEM NECESSIDADE EXPLÍCITA E APROVAÇÃO DO USUÁRIO.**

## 📋 CHECKLIST OBRIGATÓRIO (ANTES DE QUALQUER MUDANÇA)

### 1. ✅ VERIFIQUE O ESTADO ATUAL
```bash
# 1.1 Sistema está rodando?
curl http://localhost:8080
# Deve redirecionar (302) para /dashboard

# 1.2 Rotas principais funcionam?
curl http://localhost:8080/dashboard
curl http://localhost:8080/patients
curl http://localhost:8080/patients/new
# Todas devem retornar 200 OK

# 1.3 Scripts de validação passam?
./scripts/arandu_guard.sh
./scripts/arandu_checkpoint.sh
# Ambos devem passar TODOS os testes
```

### 2. ✅ VERIFIQUE O ESTADO GIT
```bash
git status
# Deve mostrar "nothing to commit, working tree clean"
git describe --tags
# Deve mostrar "v0.9.0"
```

### 3. ✅ VERIFIQUE O BUILD
```bash
go build -o arandu ./cmd/arandu
# Deve compilar sem erros
```

## 🚫 O QUE NUNCA DEVE SER ALTERADO (SEM APROVAÇÃO EXPLÍCITA)

### **LAYOUT VISUAL (NÃO MEXER):**
- ✅ Cores do sistema (`--arandu-active`, escala de neutros)
- ✅ Tipografia (Inter + Source Serif 4)
- ✅ Estrutura do layout (sidebar, top bar, main content)
- ✅ Sistema de espaçamento (`--space-xs` a `--space-2xl`)

### **ARQUITETURA (NÃO MEXER):**
- ✅ Sistema de rotas (`cmd/arandu/main.go`)
- ✅ Framework Templ (já migrado, não voltar para html/template)
- ✅ Estrutura de componentes (`web/components/`)
- ✅ Sistema de database SQLite

### **COMPORTAMENTO (NÃO MEXER):**
- ✅ Redirecionamento `/` → `/dashboard`
- ✅ Navegação entre páginas
- ✅ Funcionalidades existentes que funcionam

## 🔄 PROTOCOLO DE MUDANÇA SEGURA

### **PASSO 1: COMUNICAÇÃO**
```
1. Pergunte ao usuário: "O que exatamente precisa ser alterado?"
2. Confirme: "Isso vai alterar [especificar o que será mudado]?"
3. Alerte: "Isso pode afetar [listar possíveis impactos]?"
4. Obtenha aprovação explícita: "Posso prosseguir?"
```

### **PASSO 2: BACKUP**
```bash
# Crie uma tag de backup antes de qualquer mudança
git tag backup-pre-modificacao-$(date +%Y%m%d-%H%M%S)
```

### **PASSO 3: MUDANÇA MÍNIMA**
```
1. Altere APENAS o necessário para a tarefa
2. Mantenha compatibilidade com versão anterior
3. Não "melhore" coisas que funcionam
```

### **PASSO 4: VALIDAÇÃO (OBRIGATÓRIO)**
```bash
# 4.1 Build
go build -o arandu ./cmd/arandu

# 4.2 Teste de rotas
curl http://localhost:8080
curl http://localhost:8080/dashboard
curl http://localhost:8080/patients
curl http://localhost:8080/patients/new

# 4.3 Scripts de validação
./scripts/arandu_guard.sh
./scripts/arandu_checkpoint.sh

# 4.4 Usuário deve validar visualmente
```

## 🚨 SINAIS DE ALERTA (PARAR IMEDIATAMENTE)

Se QUALQUER UM destes ocorrer:
1. ❌ Rota `/` não redireciona para `/dashboard`
2. ❌ Layout/cores alterados sem aprovação
3. ❌ `arandu_guard.sh` falha
4. ❌ `arandu_checkpoint.sh` falha
5. ❌ Build não compila
6. ❌ Usuário diz "não era isso que eu queria"

**AÇÃO IMEDIATA:**
```bash
# 1. Reverter para v0.9.0
git checkout v0.9.0
git reset --hard HEAD
git clean -fd

# 2. Recompilar
go build -o arandu ./cmd/arandu

# 3. Validar
./scripts/arandu_guard.sh
```

## 📊 SISTEMA ARANDU v0.9.0 - ESTADO CONHECIDO

### **FUNCIONALIDADES OPERACIONAIS:**
- ✅ Dashboard clínico
- ✅ Gestão de pacientes (CRUD)
- ✅ Registro de sessões
- ✅ Observações clínicas
- ✅ Intervenções terapêuticas
- ✅ Contexto biopsicossocial
- ✅ Timeline clínica
- ✅ IA assistente reflexivo
- ✅ Busca de pacientes
- ✅ Word Cloud de temas

### **ARQUITETURA ESTÁVEL:**
- ✅ Go + Templ + Alpine.js + HTMX
- ✅ SQLite com migrations
- ✅ Design system estabelecido
- ✅ Sistema de validação (`guard.sh`, `checkpoint.sh`)

### **LAYOUT PRESERVADO:**
- ✅ Cores: Azul principal (#3b82f6) + escala de neutros
- ✅ Fontes: Inter (UI) + Source Serif 4 (conteúdo)
- ✅ Layout: Sidebar + Top bar + Main content
- ✅ Responsivo: Mobile + Desktop

## 📚 REFERÊNCIAS CRÍTICAS

### **Arquivos que DEFINEM o sistema:**
- `cmd/arandu/main.go` - Rotas e ponto de entrada
- `web/components/layout/layout.templ` - Layout visual
- `web/static/css/style.css` - Estilos do sistema
- `go.mod` - Dependências

### **Scripts de VALIDAÇÃO (executar SEMPRE):**
- `./scripts/arandu_guard.sh` - Integridade do sistema
- `./scripts/arandu_checkpoint.sh` - Validação arquitetural

### **Tags de VERSÃO:**
- `v0.9.0` - Versão atual estável (SEMPRE voltar para esta se houver problemas)

## 🏁 REGRA DE OURO

**"SE FUNCIONA, NÃO CONSERTE."**

O valor do Arandu está em sua FUNCIONALIDADE CLÍNICA, não em "melhorias" arquiteturais ou visuais não solicitadas.

**PRIORIDADE MÁXIMA:** Manter o sistema funcionando exatamente como está na v0.9.0.

---
*Documento criado após incidente de regressão em 19/03/2026. Versão do sistema: v0.9.0.*
*Todo agente DEVE ler e seguir estas diretrizes antes de qualquer modificação.*