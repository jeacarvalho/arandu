CAP-07-04 — Recuperação de Informação e Performance

Identificação

ID: CAP-07-04

Título: Recuperação de Informação e Performance

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

Descrição

Capacidade de navegar, filtrar e localizar eventos clínicos específicos dentro de grandes volumes de dados (Big Data Clínico), garantindo respostas sub-segundo e baixa carga cognitiva para o terapeuta.

Objetivos

Performance SOTA: Garantir que buscas em bases com mais de 50.000 registros retornem resultados em milissegundos.

Navegação Inteligente: Implementar mecanismos de busca que transcendam a listagem simples, como Autocomplete e Full-Text Search (FTS).

Escalabilidade de UI: Utilizar técnicas de carregamento progressivo (Infinite Scroll) para evitar o travamento do navegador.

Precisão Clínica: Permitir a localização de termos exatos dentro da narrativa clínica para fins de revisão e supervisão.

Contexto Técnico

Esta capability utiliza o motor SQLite FTS5 para indexação e a biblioteca HTMX para atualizações parciais de interface (fragments), permitindo que o sistema escale sem perder a característica de "Tecnologia Silenciosa".

Requisitos Relacionados

REQ-01-00-03: Busca e Localização de Pacientes.

REQ-07-04-01: Paginação e Infinite Scroll (Planejado).

REQ-07-04-02: Busca Contextual no Prontuário (FTS5).

Valor de Negócio

Para o terapeuta, esta capacidade significa que o Arandu nunca se tornará "pesado" ou "lento" com o passar dos anos. A ferramenta mantém a agilidade de um caderno de papel, mesmo contendo o histórico de uma década de prática clínica.

Critérios de Sucesso

Busca por nome de paciente retorna resultados em < 50ms.

Busca por termos dentro do prontuário (MATCH) retorna snippets em < 300ms.

O uso de memória no cliente permanece estável independentemente do tamanho do histórico carregado.