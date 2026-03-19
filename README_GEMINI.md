# Configuração da API Gemini para o Arandu

## 🚀 Configuração Rápida

1. **Crie um arquivo `.env` na raiz do projeto:**
```bash
echo "GEMINI_API_KEY=sua_chave_aqui" > .env
```

2. **Adicione `.env` ao `.gitignore` (se ainda não estiver):**
```bash
echo ".env" >> .gitignore
```

3. **Reinicie o servidor:**
```bash
go run cmd/arandu/main.go
```

## 🔑 Obtendo a API Key

1. Acesse [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Faça login com sua conta Google
3. Clique em "Get API Key"
4. Crie uma nova chave ou use uma existente
5. Copie a chave gerada

## ⚙️ Funcionamento do Sistema

### Sem API Key:
- O sistema inicializa com um cliente "dummy"
- A funcionalidade de IA fica desabilitada
- Mensagem de aviso no log: `"Warning: GEMINI_API_KEY not set. AI features will be disabled."`
- Botão "Solicitar Reflexão IA" retorna erro ao ser clicado

### Com API Key:
- Cliente Gemini é inicializado corretamente
- Funcionalidade de síntese reflexiva disponível
- Botão "Reflexão IA" com seletor de período funciona normalmente
- Respostas são cacheadas por 24 horas para evitar chamadas repetidas

## 🧪 Testando a Configuração

1. **Verifique se a chave está sendo carregada:**
```bash
go run cmd/arandu/main.go
# Deve mostrar: "Using database: arandu.db" sem avisos sobre API Key
```

2. **Teste a funcionalidade:**
   - Acesse um paciente com dados clínicos (ex: `/patient/p0001`)
   - Na seção "Ações Rápidas", selecione um período (3 meses, 6 meses, 1 ano, todo histórico)
   - Clique em "Reflexão IA"
   - Deve gerar uma síntese reflexiva estruturada em 4 partes:
     1. Temas Dominantes
     2. Pontos de Inflexão  
     3. Correlações Sugeridas
     4. Provocação Clínica

## 🔒 Segurança

- **NUNCA** commit a chave real no repositório
- Use `.env.example` como template para configurações
- A chave é usada apenas para comunicação com a API Gemini
- Nenhum dado identificável é enviado à API
- Respostas são cacheadas localmente por 24 horas

## 💾 Sistema de Cache

O sistema implementa cache em memória para respostas da IA:

- **TTL**: 24 horas por padrão
- **Chave de cache**: Combinação de `patientID:timeframe` (hash SHA256)
- **Benefícios**:
  - Reduz chamadas à API Gemini
  - Respostas mais rápidas para análises repetidas
  - Economia de tokens da API
- **Limpeza automática**: Entradas expiradas são removidas automaticamente

## 🐛 Solução de Problemas

### "Gemini client não inicializado"
- Verifique se o arquivo `.env` existe
- Confirme que a variável `GEMINI_API_KEY` está definida
- Reinicie o servidor após criar/editar o `.env`

### "API key não pode ser vazia"
- O cliente está tentando usar uma chave vazia
- Certifique-se de que o `.env` está na raiz do projeto
- Verifique permissões do arquivo `.env`

### Erros de conexão
- Verifique sua conexão com a internet
- Confirme que a API Key é válida e não expirou
- Verifique quotas de uso na Google AI Studio

## 📝 Exemplo de `.env`
```
# Configuração do Gemini AI
GEMINI_API_KEY=AIzaSyD...sua_chave_aqui...xyz

# Configuração do banco de dados
DATABASE_PATH=arandu.db

# Configurações do servidor
PORT=8080
ENVIRONMENT=development
```