Task: Implementação do REQ-05-01-01 — Síntese Reflexiva de Evolução (IA)

ID da Tarefa: task_20260319_ai_evolution_synthesis

Requirement: REQ-05-01-01

Stack Técnica: Go, templ, HTMX, Gemini API.

🎯 Objetivo

Implementar a primeira funcionalidade de Inteligência Assistida do Arandu. O sistema deve ser capaz de coletar os dados clínicos recentes de um paciente e gerar uma síntese reflexiva utilizando a API do Gemini, apresentando-a de forma elegante e imersiva.

🛠️ Escopo Técnico

1. Camada de Integração (internal/infrastructure/ai)

Criar o GeminiClient para comunicação com a API.

Implementar lógica de retry exponencial (até 5 vezes) conforme as regras do ecossistema Gemini.

Prompt Engineering: Definir o SystemInstruction que garanta o tom de "supervisor clínico" e o uso de terminologia técnica adequada.

2. Camada de Aplicação (internal/application/services)

Criar o AIService:

Método GeneratePatientSynthesis(ctx, patientID, timeframe).

Este método deve buscar dados das tabelas observations, interventions, patient_vitals e patient_medications para compor o prompt.

Implementar sanitização: garantir que nomes de pacientes não saiam do ambiente do servidor.

3. Camada Web (Componentes templ)

InsightCard: Componente de exibição da reflexão.

Fundo: bg-[#FFFBEB] (âmbar muito suave).

Fonte: Source Serif 4, text-xl.

Ícone: Sutil (ex: sparkles da Lucide).

GenerateInsightButton: Botão com hx-post, hx-indicator e hx-target.

4. Camada Web (Handlers)

POST /patients/{id}/analysis/synthesis:

Handler que dispara a geração.

Deve ser assíncrono do ponto de vista do utilizador (mostrar loading).

🎨 Design System (Silent UI)

O insight não deve "saltar" na cara do terapeuta. Ele deve aparecer como uma nota de rodapé sofisticada ou um card lateral que repousa sobre o fundo cinza papel.

Use a dualidade: Labels em Inter (ex: "REFLEXÃO DA IA"), Conteúdo em Source Serif 4.

🧪 Protocolo de Testes "Ironclad"

A. Teste de Integração (Mock)

Validar que o AIService monta o prompt corretamente incluindo os sinais vitais e as últimas observações.

B. Teste E2E (Playwright)

Abrir o prontuário de um paciente com dados de vitais e notas.

Clicar em "Gerar Reflexão".

Validar: O indicador de carregamento aparece.

Validar: Após alguns segundos, o InsightCard surge com o texto formatado em Serif.

Validar: O layout do histórico permanece intacto abaixo do card.

🛡️ Checklist de Integridade

[ ] A API Key do Gemini está configurada como const apiKey = "" no código (seguindo o contrato do ambiente)?

[ ] Os dados enviados para a IA são anónimos (apenas IDs e conteúdo clínico)?

[ ] O componente InsightCard é responsivo?

[ ] O scripts/arandu_guard.sh passou?

