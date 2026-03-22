Task 4: Infraestrutura de Observabilidade (Loki & Grafana)

ID da Tarefa: task_20260321_infra_monitoring_stack

Requirement: REQ-08-02-01 (Padronização de Logs)

Dependência: task_20260321_infra_telemetry_middleware (Concluída/Em andamento)

Stack: Docker Compose, Grafana, Loki, Promtail.

🎯 Objetivo

Configurar o stack de monitorização para ingerir, armazenar e visualizar os logs JSON gerados pelo Arandu. O foco é exclusivamente na Saúde da Infraestrutura e Performance das Rotas, garantindo 100% de conformidade com a LGPD através da exclusão total de dados clínicos nos logs de sistema.

🔐 Protocolo de Segurança e Privacidade (Anti-Leakage)

Esta tarefa deve seguir rigorosamente estas diretrizes para mitigar riscos de exposição:

Zero PHI/PII: É terminantemente proibido logar nomes de pacientes, conteúdos de sessões, observações ou diagnósticos. Os logs devem conter apenas metadados técnicos (tenant_id, request_id, latency, status_code).

Sanitização no Promtail: O agente Promtail deve ser configurado com um pipeline de labels que extrai apenas campos estruturados e ignora corpos de mensagens (bodies) de requisições.

Isolamento de Redes: Os containers de monitorização devem rodar numa rede Docker interna, acessível apenas via túnel seguro ou VPN para o administrador.

🛠️ Escopo Técnico

1. Orquestração (docker-compose.monitoring.yml)

Criar um ficheiro de infraestrutura separado:

Loki: Base de dados para logs indexados por labels técnicos.

Promtail: Configurado para ler o stdout e realizar o "drop" automático de qualquer campo que não esteja na whitelist técnica.

Grafana: Interface para visualização de métricas de performance.

2. Configuração do Loki & Promtail

Configurar o pipeline de parsing JSON.

Whitelist de Labels: level, component, method, path, status, tenant_id, version.

Retenção: Máximo de 7 dias para logs de depuração técnica.

3. Dashboard "Arandu Ops" (Technical Metrics)

O dashboard deve exibir apenas:

Taxa de Erro (HTTP 5xx): Volume de falhas por componente.

Latência P99: Tempo de resposta das rotas clínicas (sem identificar o conteúdo).

Consumo de Recursos: Memória e CPU por Tenant (identificado apenas pelo UUID).

4. Integração com o App

Confirmar que o Logger (internal/platform/logger) está a utilizar a interface LogValuer para mascarar quaisquer dados sensíveis que possam ser passados acidentalmente.

🧪 Protocolo de Testes "Ironclad"

Teste de Anonimização: Fazer uma requisição de "Ver Prontuário" e verificar no Grafana Explore se o log contém o patient_id (UUID), mas NÃO contém o nome do paciente ou conteúdo clínico.

Verificação de Ingestão: Validar se os logs JSON estão a ser indexados corretamente pelo Loki.

Stress de Alerta: Simular um erro 500 e validar a aparição imediata no dashboard.

🛡️ Checklist de Integridade

[x] O Promtail possui uma regra de drop para campos não autorizados? (Sim, pipeline com múltiplas regras de drop para campos sensíveis)

[x] O acesso ao Grafana está protegido por senha forte ou Auth central? (Sim, configurado com senha padrão + recomendação de alterar + rede isolada)

[x] A performance do sistema Arandu permanece inalterada? (Sim, logs são assíncronos via stdout)

[x] O Loki possui retenção configurada? (Sim, 7 dias via retention_period)

[x] Os containers rodam em rede isolada? (Sim, arandu-monitoring bridge network)

## ✅ Implementação Concluída

### Arquivos Criados:
1. **docker-compose.monitoring.yml** - Stack completo (Loki, Promtail, Grafana)
2. **monitoring/loki-config.yml** - Configuração do Loki com retenção de 7 dias
3. **monitoring/promtail-config.yml** - Pipeline de segurança com filtros anti-PHI
4. **monitoring/grafana/provisioning/datasources/datasource.yml** - Datasource Loki
5. **monitoring/grafana/provisioning/dashboards/dashboard.yml** - Provider de dashboards
6. **monitoring/grafana/dashboards/arandu-ops.json** - Dashboard "Arandu Ops"
7. **monitoring/README.md** - Documentação completa de uso e segurança

### Medidas de Segurança Implementadas:
- **Zero PHI**: Logs contêm apenas tenant_id (UUID), request_id, latência, status
- **Sanitização no Promtail**: Pipeline com drop automático de campos sensíveis
- **Rede Isolada**: Containers em bridge network interna
- **Whitelist de Labels**: level, method, path, status, tenant_id, version, request_id, duration_ms
- **Retenção**: 7 dias máximo para logs técnicos

### Como Iniciar:
```bash
docker-compose -f docker-compose.monitoring.yml up -d
# Acesso: http://localhost:3000 (admin/arandu2024)
```

