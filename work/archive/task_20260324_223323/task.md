Task: Refinação da Anamnese Multidimensional (SOTA)

ID da Tarefa: task_20260324_anamnesis_sota_refinement
Requirement: REQ-01-06-01 (Anamnese Clínica Multidimensional)
Dependência: task_20260323_refactor_layout_slp e scripts/arandu_e2e_audit.sh
Stack: Go, templ, HTMX (hx-patch), Source Serif 4.

🎯 Objetivo

Transformar o componente web/components/patient/anamnesis.templ em uma interface de escrita imersiva, garantindo salvamento automático campo a campo e conformidade total com o Standardized Layout Protocol (SLP).

🏗️ 1. Arquitetura da Interface (Escrita Distraída)

A tela deve ser dividida em blocos narrativos claros, utilizando a tipografia Source Serif 4 para todo o conteúdo inserido pelo terapeuta:

A. Campos de Domínio (Narrativa)

Queixa Principal: O motivo imediato da busca por terapia.

Histórico do Problema: Evolução dos sintomas e tentativas anteriores de tratamento.

Contexto Familiar e Social: Dinâmicas relacionais e suporte social.

Marcos de Desenvolvimento: História de vida e eventos significativos.

B. Comportamento HTMX (Salvamento Silencioso)

Cada campo (textarea) deve possuir um trigger de salvamento individual:
hx-patch="/patients/{uuid}/anamnesis/{field_name}"
hx-trigger="keyup changed delay:2s, blur"

Feedback: Um indicador visual discreto (.save-indicator) que aparece apenas durante o tráfego e desaparece com um "check" sutil ao confirmar o 200 OK.

🎨 2. Design System & Ergonomia

Tipografia: Rótulos dos campos em Inter (Sans); Conteúdo dos campos em Source Serif 4 (Clinical).

Fundo: Utilizar o .main-canvas com --arandu-bg (#E1F5EE).

Sidebar: O componente deve ser renderizado dentro do layout.Layout utilizando a sidebar_patient.templ.

🛠️ 3. Implementação Técnica

Camada de Handlers (internal/web/handlers/patient_handler.go)

Implementar o método UpdateAnamnesisField que processa o PATCH.

Validar a propriedade do paciente (Tenant Isolation).

Retornar apenas o fragmento HTML do indicador de sucesso/erro (HTMX Out-of-band se necessário).

🛡️ 4. Protocolo de Auditoria "Ironclad" (Obrigatório)

Esta tarefa NÃO pode ser dada como concluída sem a execução e aprovação do script de integridade.

Checagens Manuais e Automatizadas:

Purificação: Execute ./scripts/arandu_e2e_audit.sh.

Falha se: Encontrar qualquer atributo style="..." no HTML da Anamnese.

Falha se: Encontrar vazamento de variáveis como { p.ID }.

Layout: Confirmar se o conteúdo principal mantém o margin-left: 280px (Desktop) e se a sidebar contextual está correta.

Sintaxe: Garantir que o templ generate não gerou ficheiros órfãos ou erros de placeholders.

📋 Checklist de Entrega

[ ] Os campos de texto usam Source Serif 4?

[ ] O salvamento é automático via hx-patch sem recarregar a página?

[ ] O script ./scripts/arandu_e2e_audit.sh retornou [SUCCESS]?

[ ] A interface é 100% responsiva (Mobile-First com Bottom Nav)?

Instrução de Persona: Você é um Engenheiro de UX com alma de terapeuta. A interface deve ser tão invisível quanto uma folha de papel de alta qualidade.