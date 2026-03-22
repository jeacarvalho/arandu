# Infraestrutura de Observabilidade - Arandu

## Visão Geral

Este diretório contém a infraestrutura de monitoramento e observabilidade do sistema Arandu, configurada para operar em conformidade total com a LGPD e boas práticas de segurança em saúde.

## ⚠️ IMPORTANTE - Protocolo de Segurança

**NENHUM DADO CLÍNICO OU PHI É LOGADO OU MONITORADO**

Esta infraestrutura foi projetada especificamente para monitorar APENAS:
- Métricas de performance (latência, throughput)
- Códigos de status HTTP
- IDs técnicos (tenant_id como UUID, request_id)
- Informações de infraestrutura (CPU, memória)

**Dados que NUNCA são coletados:**
- Nomes de pacientes
- Conteúdos de sessões terapêuticas
- Observações clínicas
- Diagnósticos ou prognósticos
- Dados de medicação ou sinais vitais
- Qualquer informação identificável do paciente (PII/PHI)

## Componentes

### Loki
- Banco de dados de logs otimizado para JSON
- Retenção: 7 dias para logs de desenvolvimento
- Indexação apenas por labels técnicos

### Promtail
- Agente de coleta de logs com filtragem ativa
- Pipeline de segurança que remove automaticamente:
  - Campos não autorizados
  - Conteúdo de mensagens (msg)
  - Dados que correspondam a padrões clínicos

### Grafana
- Interface de visualização (acesso via VPN/túnel seguro)
- Dashboard "Arandu Ops" com métricas técnicas apenas
- Autenticação: admin/arandu2024 (altere em produção)

## Como Usar

### Onde Executar

Execute **todos os comandos abaixo** a partir do **diretório raiz do projeto** (onde está o arquivo `docker-compose.monitoring.yml`):

```bash
cd /caminho/para/arandu  # Diretório raiz do projeto
```

### Iniciar o Stack

```bash
# Opção 1: Usar o script automatizado (recomendado)
./scripts/start-monitoring.sh

# Opção 2: Ou iniciar manualmente
docker-compose -f docker-compose.monitoring.yml up -d

# Verificar status dos containers
docker-compose -f docker-compose.monitoring.yml ps
```

### Acessar o Grafana

```bash
# Via túnel SSH (recomendado)
ssh -L 3000:localhost:3000 usuario@servidor-arandu

# Ou diretamente (se configurado corretamente)
http://localhost:3000
# Usuário: admin
# Senha: arandu2024
```

### Verificar Logs

```bash
# Ver logs do Loki
docker-compose -f docker-compose.monitoring.yml logs loki

# Ver logs do Promtail
docker-compose -f docker-compose.monitoring.yml logs promtail

# Ver logs do Grafana
docker-compose -f docker-compose.monitoring.yml logs grafana
```

### Parar o Stack

```bash
docker-compose -f docker-compose.monitoring.yml down

# Para remover também os volumes (apaga dados):
docker-compose -f docker-compose.monitoring.yml down -v
```

## Dashboards Disponíveis

### Arandu Ops - Technical Metrics

- **System Health**: Taxa de erros, latência P99, requests/seg, taxa de sucesso
- **Performance Metrics**: Latência por rota (P50/P99), distribuição de status HTTP
- **Tenant Activity**: Request rate por tenant (apenas UUIDs, sem identificação)
- **Log Explorer**: Visualização de logs com campos técnicos apenas

## Testes de Segurança

### Teste 1: Verificar ausência de PHI

```bash
# Fazer uma requisição de exemplo
curl http://localhost:8080/api/patients \
  -H "Authorization: Bearer <token>"

# No Grafana Explore, procurar por logs desta requisição
# Verificar: deve conter tenant_id e request_id
# NÃO deve conter: nomes, conteúdos clínicos, etc.
```

### Teste 2: Validar filtros do Promtail

```bash
# Verificar configuração
docker-compose -f docker-compose.monitoring.yml exec promtail \
  cat /etc/promtail/config.yml | grep -A5 "drop"

# Ver métricas de drop
curl http://localhost:9080/metrics | grep promtail_dropped_entries_total
```

### Teste 3: Simular erro 500

```bash
# Causar um erro proposital no endpoint
curl http://localhost:8080/api/patients/invalid-uuid

# Verificar no dashboard se o erro aparece em <10s
```

## Segurança de Rede

- Os containers rodam em uma rede Docker isolada (`arandu-monitoring`)
- Apenas o Grafana expõe porta para o host (3000)
- Loki e Promtail são acessíveis apenas internamente
- Recomendação: Usar VPN ou túnel SSH para acesso externo

## Manutenção

### Limpar logs antigos

O Loki está configurado com retenção automática de 7 dias. Para limpar manualmente:

```bash
# Acessar o container do Loki
docker-compose -f docker-compose.monitoring.yml exec loki sh

# Limpar dados antigos
rm -rf /loki/chunks/*
rm -rf /loki/index/*
```

### Backup

```bash
# Backup dos dados do Loki
docker-compose -f docker-compose.monitoring.yml exec loki tar -czf /tmp/loki-backup.tar.gz /loki

# Copiar para fora do container
docker cp loki:/tmp/loki-backup.tar.gz ./backup/
```

## Troubleshooting

### Docker Permission Denied

Se você receber erro `permission denied while trying to connect to the Docker daemon socket`:

```bash
# O usuário precisa estar no grupo docker
groups | grep docker

# Se não estiver, adicione (requer logout/login):
sudo usermod -aG docker $USER
newgrp docker

# Testar
docker ps
```

Veja [TROUBLESHOOTING.md](TROUBLESHOOTING.md) para mais detalhes.

### Promtail não está enviando logs

```bash
# Verificar configuração
docker-compose -f docker-compose.monitoring.yml logs promtail | tail -20

# Verificar conectividade com Loki
docker-compose -f docker-compose.monitoring.yml exec promtail \
  wget -qO- http://loki:3100/ready
```

### Grafana não conecta ao Loki

```bash
# Verificar datasource
curl -u admin:arandu2024 http://localhost:3000/api/datasources

# Testar conectividade
docker-compose -f docker-compose.monitoring.yml exec grafana \
  wget -qO- http://loki:3100/ready
```

### Não aparecem logs no dashboard

1. Verifique se a aplicação Arandu está rodando
2. Confira se os logs estão sendo enviados para stdout
3. Verifique se o Promtail está configurado para o container correto
4. Confira os filtros de drop no Promtail

## Contato

Para questões sobre segurança ou conformidade, entre em contato com o time de segurança da informação.
