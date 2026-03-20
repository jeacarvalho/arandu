REQ-07-03-05 — Provisão Automática de Tenant e Criação de DB Físico

Identificação

ID: REQ-07-03-05

Capability: CAP-07-03 — Gestão de Acesso e Multi-tenancy

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

História do utilizador

Como novo utilizador do Arandu, quero que o meu ambiente clínico exclusivo seja criado automaticamente no meu primeiro acesso, para que eu possa começar a registar os meus pacientes imediatamente, com a garantia de que os meus dados já nascem isolados e seguros.

Contexto

Este requisito descreve a "fábrica" de consultórios digitais. No momento em que um profissional se regista (via e-mail ou Google OAuth), o sistema não apenas cria um registo no Banco Central, mas deve provisionar um ficheiro físico SQLite no servidor. Este processo garante a Soberania do Dado desde o segundo zero.

Descrição funcional

O sistema deve automatizar a criação da infraestrutura de dados para novos utilizadores.

Trigger de Criação: Disparado após o sucesso do primeiro login (se o tenant_id não possuir um ficheiro associado).

Geração de Identidade: Criar um UUID único para o tenant_id.

Criação de Ficheiro: Gerar o ficheiro storage/tenants/clinical_{tenant_id}.db.

Bootstrap de Schema: Executar imediatamente o Migrator clínico para criar as tabelas de pacientes, sessões e observações no novo ficheiro.

Verificação de Saúde: Validar se a conexão com o novo banco é estável antes de redirecionar o utilizador para o Dashboard.

Lógica Técnica (SOTA)

1. O Fluxo de Provisão (Go)

func ProvisionNewTenant(ctx context.Context, userID string) (string, error) {
    tenantID := uuid.New().String()
    dbPath := fmt.Sprintf("storage/tenants/clinical_%s.db", tenantID)
    
    // 1. Criar ficheiro físico
    db, err := sql.Open("sqlite", dbPath)
    
    // 2. Executar Migrations Clínicas (go:embed)
    migrator := sqlite.NewMigrator(db)
    if err := migrator.Migrate(ctx); err != nil {
        return "", err
    }
    
    // 3. Registar no Control Plane (Banco Central)
    return tenantID, centralRepo.RegisterTenant(ctx, userID, tenantID, dbPath)
}


Interface (Tecnologia Silenciosa)

O utilizador não deve ver "barras de progresso técnicas" complexas.

Mensagem de Boas-vindas: Enquanto o banco é criado (processo que demora milissegundos), exibir: "A preparar o seu consultório seguro..." com uma animação botânica sutil.

Finalização: Redirecionamento suave para o Dashboard vazio, pronto para o primeiro paciente.

Critérios de Aceitação

CA-01: O sistema deve criar um ficheiro .db físico para cada novo utilizador registado.

CA-02: O novo banco de dados deve conter exatamente o mesmo schema (tabelas e índices) definido nas migrations clínicas.

CA-03: Se a criação do ficheiro falhar (ex: falta de espaço em disco), o sistema deve reverter o registo no Banco Central (Atomicidade).

CA-04: O ficheiro deve ser criado com as permissões de leitura/escrita restritas ao utilizador do sistema operativo que corre o Arandu.

CA-05: O processo completo de provisão não deve exceder 2 segundos de latência para o utilizador final.

Persistência

Control Plane: INSERT INTO tenants (id, db_path, status, created_at) VALUES (?, ?, 'active', ?)
FileSystem: storage/tenants/clinical_{id}.db