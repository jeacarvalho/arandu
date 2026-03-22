#!/bin/bash

# Script de verificação do stack de monitoramento
# Execute em um ambiente com Docker instalado

set -e

echo "🚀 Iniciando stack de monitoramento Arandu..."
echo ""

# Verificar se Docker está instalado
if ! command -v docker &> /dev/null; then
    echo "❌ Docker não encontrado. Instale o Docker primeiro."
    exit 1
fi

# Verificar se docker-compose está instalado
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose não encontrado. Instale o docker-compose primeiro."
    exit 1
fi

# Subir o stack
echo "📦 Subindo containers..."
docker-compose -f docker-compose.monitoring.yml up -d

echo ""
echo "⏳ Aguardando serviços iniciarem..."
sleep 5

# Verificar status
echo ""
echo "📊 Status dos containers:"
docker-compose -f docker-compose.monitoring.yml ps

echo ""
echo "🔍 Verificando saúde dos serviços..."

# Verificar Loki
if curl -s http://localhost:3100/ready > /dev/null 2>&1; then
    echo "✅ Loki está respondendo"
else
    echo "⏳ Loki ainda está iniciando..."
fi

# Verificar Grafana
if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
    echo "✅ Grafana está respondendo"
else
    echo "⏳ Grafana ainda está iniciando..."
fi

echo ""
echo "🌐 Acesso ao Grafana:"
echo "   URL: http://localhost:3000"
echo "   Usuário: admin"
echo "   Senha: arandu2024"
echo ""
echo "📋 Comandos úteis:"
echo "   Ver logs:        docker-compose -f docker-compose.monitoring.yml logs -f"
echo "   Parar:           docker-compose -f docker-compose.monitoring.yml down"
echo "   Ver config:      cat monitoring/promtail-config.yml"
echo ""
echo "🔒 Nota de segurança:"
echo "   Este stack está configurado para NÃO capturar dados clínicos."
echo "   Apenas métricas técnicas (tenant_id, latência, status) são logadas."
echo ""
echo "⚠️  IMPORTANTE: Altere a senha do Grafana em ambiente de produção!"
echo "   docker-compose -f docker-compose.monitoring.yml exec grafana grafana-cli admin reset-admin-password 'nova-senha'"
