REQ-07-03-02 — Orquestração de Conexão Dinâmica

Identificação

ID: REQ-07-03-02

Capability: CAP-07-03 Gestão de Acesso e Multi-tenancy

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

História do usuário

Como psicólogo clínico, quero que o sistema selecione e abra automaticamente meu banco de dados exclusivo após eu me logar, para que eu possa trabalhar com meus dados de forma isolada e segura, sem risco de mistura com dados de outros profissionais.

Contexto

Diferente de aplicações multi-tenant tradicionais (que filtram uma tabela única por ID), o Arandu utiliza a estratégia de Database-per-Tenant. Isso exige um mecanismo que identifique o profissional autenticado, localize seu arquivo físico SQLite (.db) e injete essa conexão ativa no fluxo de execução da requisição HTTP.

Descrição funcional

O sistema deve implementar um gerenciador de conexões dinâmicas (Connection Manager).

Identificação: O TenantID (obtido no REQ-07-03-01) é utilizado para construir o caminho do arquivo.

Abertura Sob Demanda: O banco clínico é aberto apenas quando o usuário faz a primeira requisição após o login.

Cache de Conexão: Conexões ativas devem ser mantidas em memória para evitar o custo de abertura a cada clique do HTMX, respeitando limites de recursos.

Middleware: Um middleware de orquestração deve injetar a conexão no context.Context de cada requisição protegida.

Lógica Técnica (SOTA)

1. Caminho do Arquivo

O caminho deve seguir um padrão seguro definido na configuração:
PATH_STORAGE_CLINICAL/clinical_{tenant_id}.db

2. Middleware de Tenant

func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Extrair TenantID da Sessão/JWT
        tenantID := getTenantIDFromSession(r)
        
        // 2. Obter ou Abrir Conexão via ConnectionManager
        db, err := connectionManager.GetDB(tenantID)
        
        // 3. Injetar no Contexto
        ctx := context.WithValue(r.Context(), "tenant_db", db)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}


3. Gerenciamento de Ciclo de Vida

Auto-Migrate: Ao abrir uma conexão pela primeira vez, o Migrator deve ser executado naquele banco específico para garantir que o schema clínico está atualizado.

Idle Timeout: Conexões inativas por mais de X minutos devem ser fechadas para liberar descritores de arquivo no SO.

Interface (Padrão Arandu SOTA)

Embora seja um requisito majoritariamente de backend, a interface deve refletir o sucesso dessa orquestração:

Indicador de Conexão: Exibir sutilmente no perfil que o "Consultório Privado" está montado e seguro.

Tratamento de Erros: Caso o arquivo do banco esteja corrompido ou inacessível, exibir uma tela de erro amigável (via HTMX) instruindo o usuário a contatar o suporte, sem expor caminhos de diretórios.

Fluxo

O Middleware intercepta a requisição.

O sistema verifica se a conexão para o TenantID atual já está no "Pool de Tenants".

Se não estiver: - Localiza o arquivo físico.

Abre a conexão SQLite.

Executa migrator.Migrate(db).

Adiciona ao Pool.

O Repository recupera o *sql.DB do contexto e executa a query clínica.

Critérios de Aceitação

CA-01: O sistema deve lançar um erro fatal se uma requisição tentar acessar dados clínicos sem um TenantID válido no contexto.

CA-02: Operações em um banco clínico (ex: INSERT patient) não devem ser visíveis em outros bancos de outros TenantIDs.

CA-03: O sistema deve suportar a abertura concorrente de múltiplos bancos (terapeutas diferentes acessando ao mesmo tempo).

CA-04: Ao reiniciar o servidor, o sistema deve ser capaz de reabrir as conexões conforme a demanda.

CA-05: O Migrator deve garantir a paridade de schema entre todos os bancos de usuários ativos.

Persistência

Control Plane: SELECT path FROM tenants WHERE id = ?
Data Plane: Conexão direta via driver SQLite (file:path_to_db).

Fora do escopo

Sharding de bancos em servidores diferentes.

Criptografia de disco (at-rest) gerenciada pela aplicação (delegado ao SO/File System).