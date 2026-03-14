#!/bin/bash

PORT=$1
shift
CMD="$@"

if [ -z "$PORT" ]; then
  echo "Uso: $0 <porta> <comando>"
  echo "Exemplo: $0 8080 go run main.go"
  exit 1
fi

if [ -z "$CMD" ]; then
  echo "Erro: informe o comando a executar após liberar a porta"
  exit 1
fi

echo "Verificando processos na porta $PORT..."

PIDS=$(lsof -ti tcp:$PORT)

if [ -n "$PIDS" ]; then
  echo "Processos encontrados:"
  ps -fp $PIDS

  echo "Matando processos..."
  kill $PIDS

  sleep 1

  for PID in $PIDS
  do
    if ps -p $PID > /dev/null
    then
      echo "Forçando kill -9 no PID $PID"
      kill -9 $PID
    fi
  done

  echo "Porta $PORT liberada."
else
  echo "Nenhum processo usando a porta $PORT."
fi

echo "Executando comando:"
echo "$CMD"
echo "---------------------------"

exec $CMD
