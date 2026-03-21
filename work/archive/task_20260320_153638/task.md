Task: Onboarding e Provisão Automática de Tenant

ID da Tarefa: task_20260320_tenant_provisioning

Requirement: REQ-07-03-05

Dependência: task_20260320_google_oauth (Task 4) e task_20260320_session_middleware (Task 2)

Stack: Go, SQLite, os/exec ou Migrator.

🎯 Objetivo

Implementar o fluxo de "Boas-vindas Técnico": assim que um utilizador faz o primeiro login (via Google ou E-mail), o sistema deve criar fisicamente o seu banco de dados SQLite exclusivo, aplicar as migrações clínicas e redirecioná-lo para o Dashboard com a sessão ativa.

🛠️ Escopo Técnico

1. Serviço de Provisão (internal/application/services/tenant_service.go)

Implementar a lógica ProvisionNewTenant(userID):

Geração de Caminho: Definir o destino storage/tenants/clinical_{tenant_id}.db.

Criação Física: Criar o ficheiro no sistema de ficheiros.

Bootstrap de Dados: * Abrir a conexão com o novo ficheiro.

Instanciar o Migrator clínico (o mesmo que já usamos no projeto).

Executar todas as migrações .up.sql para que o novo banco já nasça com as tabelas de patients, sessions, etc.

Vínculo no Control Plane: Atualizar a tabela tenants no arandu_central.db com o caminho do novo banco.

2. Integração no Fluxo de Auth (Handlers)

No callback do Google ou no Handler de Login:

Lógica: ```go
if user.TenantID == "" {
// Primeiro acesso: disparar provisão
tenantID, _ := tenantService.ProvisionNewTenant(user.ID)
user.TenantID = tenantID
}
// Estabelecer cookie de sessão com TenantID
// Redirecionar para /dashboard




3. Página de "Aguarde" (UX Silenciosa)

Se a provisão demorar (devido às migrations), o utilizador deve ver uma tela simples com a tipografia Source Serif 4:

"A preparar o seu consultório digital seguro... Isto levará apenas um instante."

🧪 Protocolo de Verificação SOTA

A. Teste de Primeiro Acesso

Limpar a pasta storage/tenants/.

Fazer login com uma conta Google nova.

Verificar: O redirecionamento para /dashboard ocorre com sucesso.

Verificar: Um novo ficheiro .db apareceu na pasta de tenants.

Verificar: O banco central registou o caminho correto.

B. Teste de Reentrada

Fazer logout.

Fazer login novamente com a mesma conta.

Verificar: O sistema NÃO cria um novo banco, mas usa o existente e redireciona instantaneamente.

🛡️ Checklist de Integridade

[ ] O banco clínico criado recebeu as migrations corretamente?

[ ] O tenant_id está gravado no cookie de sessão de forma segura?

[ ] O scripts/arandu_guard.sh confirma que o utilizador agora chega ao Dashboard?

[ ] O sistema lida com erro de "Permissão de Escrita" na pasta storage?

