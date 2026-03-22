CAP-08-02 — Observabilidade e Diagnóstico

Identificação

ID: CAP-08-02

Título: Observabilidade e Diagnóstico

Vision: VISION-01 — Registro da Prática Clínica (Suporte à Operação)

Status: draft

Descrição

Prover visibilidade total sobre o comportamento técnico e operacional do Arandu. Esta capacidade permite que a equipe de engenharia identifique falhas, gargalos de performance e comportamentos anômalos em tempo real, garantindo a alta disponibilidade do sistema para os terapeutas.

Objetivos

Rastreabilidade: Identificar o percurso de uma requisição desde o clique do utilizador até à persistência no banco de dados.

Diagnóstico Rápido: Reduzir o MTTR (Mean Time To Repair) através de logs ricos e contextualizados.

Monitorização de Performance: Medir latências e tempos de resposta de forma segmentada por rota e tenant.

Paridade de Ambientes: Garantir que o comportamento de diagnóstico seja idêntico em desenvolvimento e produção.

Contexto Técnico (SOTA)

A observabilidade no Arandu baseia-se em Logs Estruturados (JSON) e rastreio via Request ID.

Engine: log/slog (Standard Library Go).

Coleta: As saídas de sistema (stdout) são capturadas por agentes de log (como Loki) sem necessidade de escrita em ficheiros locais.

Versionamento: Metadados de Git Tags são injetados no log para correlacionar erros com versões específicas do código.

Requisitos Relacionados

REQ-08-02-01: Padronização de Logs Estruturados.

REQ-08-02-02: Rastreabilidade de Requisições (Request ID).

Valor de Negócio

Para o administrador do sistema, esta capability garante que o Arandu seja uma plataforma "transparente". Em caso de erro de um médico, a causa raiz pode ser isolada sem violar a privacidade do conteúdo clínico, protegendo a reputação da plataforma.

Critérios de Sucesso

Logs emitidos em JSON válido contendo tenant_id e version.

Tempo de resposta de cada rota logado automaticamente.

Zero menção a conteúdos sensíveis (PII) nos logs de sistema.