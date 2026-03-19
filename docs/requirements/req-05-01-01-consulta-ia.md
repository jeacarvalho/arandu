REQ-05-01-01 — Síntese Reflexiva de Evolução (IA)

Identificação

ID: REQ-05-01-01

Capability: CAP-05-01 — Assistente Reflexivo com IA

Vision: VISION-05 — Assistência Reflexiva com IA

Status: implemented

História do Utilizador

Como psicólogo clínico, quero que o sistema gere uma síntese reflexiva sobre a evolução do paciente num determinado período, para identificar padrões latentes, mudanças de comportamento e correlações entre a biologia (vitals/medicação) e a narrativa psíquica.

Contexto

Com anos de dados, é difícil para o terapeuta manter a visão macro do caso. O Arandu utilizará modelos de linguagem avançados (LLMs) para atuar como um "supervisor clínico silencioso". A síntese não é um diagnóstico, mas uma provocação reflexiva baseada exclusivamente nos dados registados no prontuário.

Descrição Funcional

O sistema deve fornecer uma funcionalidade de "Gerar Reflexão" dentro do histórico do paciente.

Seleção de Contexto: O terapeuta pode selecionar o período (ex: últimos 3 meses, últimas 10 sessões).

Input para a IA: O sistema enviará para o modelo:

Snippets das observações clínicas.

Variações nos sinais vitais (peso, sono).

Alterações no histórico farmacológico.

Dados demográficos SOTA (ocupação, etnia, etc.).

Output Reflexivo: A IA deve devolver um texto estruturado contendo:

Temas Dominantes: O que mais ocupou o espaço psíquico.

Pontos de Inflexão: Mudanças notáveis na narrativa ou humor.

Correlações Sugeridas: Ex: "Observou-se maior resistência nas sessões após a troca da medicação X".

Provocação Clínica: Uma pergunta para o terapeuta considerar na próxima sessão.

Interface (Tecnologia Silenciosa)

Gatilho: Um botão "Reflexão IA" com seletor de período (3 meses, 6 meses, 1 ano, todo histórico) na seção "Ações Rápidas" da página do paciente.

Exibição: O conteúdo aparece num componente InsightCard com fundo âmbar (#FFFBEB), borda sutil (#FDE68A), utilizando a fonte clínica configurada (Source Serif 4 via CSS custom property).

Estado de Carregamento: Um indicador de carregamento com animação "pulse" e placeholders que não bloqueia a navegação.

Lógica Técnica (IA e Segurança)

Modelo: Gemini 2.5 Flash Lite via API Google Generative AI.

Prompt System: "Atue como um supervisor clínico experiente. Analise os dados abaixo e identifique padrões longitudinais. Seja breve, técnico e reflexivo. Nunca dê diagnósticos fechados. Foque em: 1. Temas Dominantes, 2. Pontos de Inflexão, 3. Correlações Sugeridas, 4. Provocação Clínica."

Privacidade: Nenhum dado identificável (Nome, CPF) é enviado para a API. Apenas conteúdo clínico anonimizado.

Retry Exponencial: Implementado com backoff exponencial (até 5 tentativas) para lidar com falhas temporárias da API.

Cache: Respostas cacheadas em memória por 24 horas usando chave SHA256(patientID:timeframe).

Critérios de Aceitação

CA-01: A síntese deve integrar dados de observações, intervenções e sinais vitais.

CA-02: O texto gerado deve ser renderizado em Source Serif 4, diferenciando-se visualmente da UI administrativa.

CA-03: O sistema deve armazenar a última reflexão gerada para evitar chamadas repetitivas à API (Cache).

CA-04: Deve haver um aviso legal: "Esta é uma análise gerada por IA para apoio à reflexão. A decisão clínica é exclusiva do profissional."

Fora do Escopo

Transcrição de áudio de sessões.

Chat em tempo real com a IA (será REQ-05-01-02).

Envio de relatórios por e-mail.