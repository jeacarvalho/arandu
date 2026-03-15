# TASK 20260315_131839

Title: Corrigir bugs do sistema

## Objetivo

a) Nao consigo subir o app --> 2026/03/15 13:12:21 Failed to run database migrations: failed to create migration manager: failed to read migrations directory: open .: no such file or directory

Vamos começar com essa correção e durante os testes vou inserindo novas informações aqui

## Progresso

✅ **BUG 1 CORRIGIDO**: Erro de migração resolvido
- **Problema**: O diretório de migrações estava sendo construído com caminho relativo (`.`), causando "open .: no such file or directory"
- **Solução**: Adicionado lógica para usar caminho absoluto do executável com fallback para caminho relativo
- **Arquivo corrigido**: `cmd/arandu/main.go:25-35`

✅ **BUG 2 CORRIGIDO**: Erro no template `dashboard.html`
- **Erro**: `template: dashboard.html:263:68: executing "dashboard.html" at <.SessionNumber>: can't evaluate field SessionNumber in type interface {}`
- **Problema**: O template tentava acessar campo `.SessionNumber` que não existe na estrutura de dados
- **Solução**: Removida referência a `.SessionNumber` do template
- **Arquivo corrigido**: `web/templates/dashboard.html:263`

✅ **BUG 3 CORRIGIDO**: Rota `/patients/new` carregando template errado
- **Problema**: A rota `/patients/new` estava retornando template de "Nova Sessão" em vez de "Novo Paciente"
- **Causa**: Conflito de nomes de templates - múltiplos templates definindo `{{define "content"}}`
- **Solução**: 
  1. Criado template `new_patient_page` com nome único
  2. Modificado handler `NewPatient` para servir arquivo diretamente
  3. Criado testes para garantir unicidade de nomes de templates
- **Arquivos corrigidos**: 
  - `web/templates/new_patient.html` (template completo)
  - `web/handlers/handler.go:209-225` (handler simplificado)
  - `tests/handlers/simple_template_test.go` (testes de prevenção)

🎉 **APP FUNCIONANDO**: O sistema agora sobe sem erros e está rodando em `http://localhost:8080`
🔒 **PREVENÇÃO**: Testes criados para evitar conflitos de templates no futuro

## Descrição

Tarefa criada sem requirement específico.
