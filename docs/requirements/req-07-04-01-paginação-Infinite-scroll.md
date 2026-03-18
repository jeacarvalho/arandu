REQ-07-04-01 — Paginação e Infinite Scroll

Identificação

ID: REQ-07-04-01

Capability: CAP-07-04 — Recuperação de Informação e Performance

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

História do utilizador

Como psicólogo clínico, quero que as listas de pacientes e sessões sejam carregadas automaticamente à medida que faço scroll, para manter a agilidade do sistema e evitar a sobrecarga de informação de uma lista interminável.

Contexto

Com 63.000 sessões e centenas de pacientes, carregar o HTML de todos os registos de uma vez destruiria a performance do navegador e a experiência do utilizador. Este requisito implementa o padrão de "Carregamento Progressivo" (Infinite Scroll), substituindo a paginação tradicional por um fluxo contínuo e silencioso.

Descrição funcional

O sistema deve carregar dados em "lotes" (batches) de tamanho fixo.

Tamanho do Lote: 20 itens por carregamento.

Mecanismo de Gatilho: Utilizar o evento revealed do HTMX no último item da lista atual.

Indicador de Carregamento: Exibir um sinal visual sutil (ex: um pequeno spinner ou mensagem "A carregar mais...") enquanto os dados são recuperados.

Estado da UI: A adição de novos itens não deve fazer a página "saltar" ou perder a posição do scroll.

Interface (Padrão Arandu SOTA)

Seguindo a filosofia de Tecnologia Silenciosa:

Sem Botões: Eliminar os botões de "Próxima Página" ou "Anterior". A navegação deve ser orgânica através do scroll.

Transição: Os novos itens devem aparecer com uma transição de opacidade sutil (fade-in) para evitar quebras visuais bruscas.

Final da Lista: Quando não houver mais dados, o sistema deve exibir uma mensagem discreta: "Fim dos registos históricos".

Lógica Técnica (HTMX + SQL)

1. Parâmetros de Query

O Handler deve aceitar os parâmetros limit (padrão 20) e offset (calculado pelo cliente ou incrementado).

2. Implementação HTMX

O último elemento de cada lote deve conter os atributos de carregamento:

<div class="patient-item" 
     hx-get="/patients?offset=20" 
     hx-trigger="revealed" 
     hx-swap="afterend">
    <!-- Conteúdo do 20º paciente -->
</div>


3. Persistência (SQL)

SELECT id, name, created_at 
FROM patients 
ORDER BY created_at DESC 
LIMIT 20 OFFSET ?;


Fluxo

O utilizador abre a lista de pacientes.

O sistema carrega os primeiros 20 pacientes.

O utilizador faz scroll até ao final da lista.

O HTMX deteta que o último item está visível (revealed).

O sistema solicita os próximos 20 itens via GET.

O backend processa a query com o novo offset.

O fragmento HTML é injetado após o último item atual.

Critérios de Aceitação

CA-01: O carregamento do próximo lote deve ocorrer de forma automática sem necessidade de clique.

CA-02: O sistema não deve carregar dados duplicados se o utilizador fizer scroll rápido (uso de ID único para verificação se necessário).

CA-03: A performance de renderização de cada lote deve ser inferior a 200ms.

CA-04: O estado da sidebar e da top bar deve permanecer inalterado durante o carregamento.

CA-05: Se houver um erro de rede, o sistema deve exibir uma opção de "Tentar Novamente" apenas para aquele lote.

Fora do escopo

"Salto" para páginas específicas (ex: ir direto para a página 50).

Ordenação dinâmica pelo utilizador durante o scroll (a ordenação é fixa por contexto).