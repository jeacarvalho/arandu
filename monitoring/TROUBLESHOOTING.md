# Guia de Troubleshooting - Docker

## Problema de Permissão com Docker Socket

Se você estiver recebendo erros de `permission denied` ao tentar usar o Docker, é porque o usuário precisa ter permissão para acessar o socket do Docker.

### Diagnóstico

```bash
# Verificar se o usuário está no grupo docker
groups | grep docker

# Verificar permissões do socket Docker
ls -la /var/run/docker.sock
# Deve mostrar: srw-rw---- 1 root docker ...
```

### Soluções

#### Opção 1: Adicionar usuário ao grupo docker (Requer logout/login)

```bash
# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# Fazer logout e login novamente para a mudança fazer efeito
# Ou iniciar um novo shell com o grupo atualizado:
newgrp docker

# Verificar se funcionou
docker ps
```

#### Opção 2: Usar sudo (Temporário)

```bash
sudo docker-compose -f docker-compose.monitoring.yml up -d
```

#### Opção 3: Verificar serviço Docker

```bash
# Verificar se o Docker está rodando
sudo systemctl status docker

# Iniciar se estiver parado
sudo systemctl start docker

# Habilitar para iniciar automaticamente
sudo systemctl enable docker
```

### Erro: "URLSchemeUnknown: Not supported URL scheme http+docker"

Este erro pode ocorrer se houver incompatibilidade entre versões do docker-compose e a biblioteca Python Docker.

```bash
# Verificar versões
docker --version
docker-compose --version

# Se necessário, reinstalar docker-compose
pip3 install --upgrade docker-compose
# ou
sudo apt-get update && sudo apt-get install docker-compose
```

### Verificação Final

Após aplicar as correções:

```bash
# Deve mostrar containers em execução (ou vazio, sem erro)
docker ps

# Testar com hello-world
docker run hello-world
```

### Configuração Alternativa

Se o problema persistir, você pode usar uma configuração TCP (menos seguro, apenas para desenvolvimento local):

```bash
# Editar /etc/docker/daemon.json
{
  "hosts": ["unix:///var/run/docker.sock", "tcp://127.0.0.1:2375"]
}

# Reiniciar Docker
sudo systemctl restart docker

# Configurar variável de ambiente
export DOCKER_HOST=tcp://127.0.0.1:2375
```

**⚠️ Atenção**: Não use configuração TCP em produção sem TLS.
