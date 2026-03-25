# Missão: Construção do Auditor E2E Supremo (SOTA Guard)

## Status
**CONCLUIDA** ✅

Contexto: O sistema cresceu e falhas de "ghosting" (vazamento de código como { p.ID }) e estilos inline persistem. Sua tarefa é criar o ./scripts/arandu_e2e_audit.sh, um script que simula uma jornada clínica completa, validando cada saída HTML contra o Design System Botânico.

🏗️ 1. Configuração do Ambiente de Teste (In-Memory)

O script deve instruir a execução do app em modo de auditoria para garantir um estado limpo:

Banco de Dados: Utilize a infraestrutura definida em internal/platform/storage/factory.go. Ao definir APP_ENV=test, o sistema deve usar SQLite in-memory (cache=shared).

Migrations: O script deve garantir que o internal/infrastructure/repository/sqlite/migrator.go seja executado tanto para o Control Plane (storage/arandu_central.db) quanto para o novo Tenant no início do teste.

Limpeza de Lixo: Antes de iniciar, o script deve remover tmp/cookies.txt e qualquer banco residual em storage/tenants/test_*.

🏃 2. A Jornada do Guerreiro (Passos do Teste)

O script deve utilizar curl e gerir o estado em tmp/cookies.txt, testando os handlers localizados em internal/web/handlers/. Atenção: Verifique as rotas exatas definidas no main.go e auth_handler.go antes de prosseguir.

Passo 1: Identidade (Auth)

Ação: POST /login (ou a rota definida em auth_handler.go).

Validação: Verificar se o servidor responde com sucesso e se o cookie arandu_session (ou similar) foi gerado. O teste deve validar se o arquivo .db do tenant foi provisionado no caminho hashed (ex: storage/tenants/xx/yy/...).

Passo 2: Criação de Paciente SOTA

Ação: POST /patients enviando dados de REQ-01-00-01 (Nome, Etnia, Gênero).

Validação: Capturar o HTML de /patients/{uuid} e validar a presença das classes do SLP.

Passo 3: Anamnese Multidimensional

Ação: PATCH /patients/{uuid}/anamnesis/chief_complaint (Handler: patient_handler.go).

Validação: O fragmento retornado deve conter a classe .save-indicator e o texto em .font-clinical.

Passo 4: Registro de Sessão e Observação

Ação: POST /sessions e POST /sessions/{id}/observations (Handler: observation_handler.go).

Validação: Simular requisição HTMX (Header HX-Request: true). Validar se o servidor retorna apenas o componente web/components/session/observation_item.templ e não o layout completo.

🔍 3. Níveis de Auditoria "Pente-Fino"

Para cada captura de HTML, o script deve rodar estas checagens:

A. Detecção de Ghosting (Vazamento de Código)

Regex: grep -P '\{ \.?[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+ \}'.

Alvo: Se encontrar { p.Name }, { p.ID } ou qualquer placeholder do templ não processado, o teste falha.

B. Purificação Estética (Zero-Style)

Regex: grep -o 'style="' | wc -l.

Alvo: O contador deve ser 0. Toda estilização deve vir do web/static/css/style.css.

C. Protocolo SLP e Sidebar de Contexto

Check: Validar se no perfil do paciente a sidebar contém os links de web/components/layout/sidebar_patient.templ (Anamnese, Prontuário, Metas).

D. Dualidade Tipográfica

Check: Validar se o conteúdo de observation_item está envolto por tags com a classe .font-clinical (Source Serif 4).

📊 4. Saída e Relatório Técnico

O script deve gerar um log visual no terminal e salvar os artefatos em tmp/audit_logs/:

[STEP 1] Auth & Provisioning ........ [PASS]
[STEP 2] SOTA Patient Profile ....... [PASS]
[STEP 3] Anamnesis Auto-save ........ [FAIL] -> Ghosting detected: { p.ID }
[STEP 4] HTMX Observation Fragment .. [PASS]


Código de Saída: 0 (Sucesso) ou 1 (Falha).

🛡️ Instrução de Conclusão para o Agente

"A tarefa só será considerada concluída se ./scripts/arandu_e2e_audit.sh retornar SUCESSO. Se houver falha, analise o HTML em tmp/audit_logs/, corrija o componente .templ ou o CSS em web/static/css/style.css, rode templ generate e repita o teste."

Persona: Você é um QA Engineer Sênior. Sua obsessão é o HTML purificado. O código Go é o motor, mas o HTML é a experiência sagrada do terapeuta.