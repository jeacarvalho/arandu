# Checklist LGPD por Feature

## Ao criar qualquer nova feature que toca dados de pacientes

### 1. Mapeamento de dados
- [ ] Quais dados pessoais são coletados?
- [ ] Quais dados sensíveis (saúde) são coletados?
- [ ] Onde são armazenados? (tabela, campo)
- [ ] Por quanto tempo são retidos?
- [ ] Quem tem acesso?

### 2. Base legal
- [ ] Qual a base legal para este tratamento?
  - Consentimento (TCLE) — para dados do prontuário
  - Tutela da saúde — para tratamento clínico
  - Obrigação legal — para retenção mínima
- [ ] O consentimento é específico para esta finalidade?
- [ ] O titular pode revogar sem prejuízo?

### 3. Direitos do titular
- [ ] Como o paciente exerce direito de acesso?
- [ ] Como o paciente exerce direito de portabilidade?
- [ ] O que acontece ao revogar consentimento?
  - Novos atendimentos: encerrados
  - Prontuário existente: mantido (obrigação legal)
  - Dados de conta: anonimizados após período

### 4. Segurança
- [ ] Dados em trânsito são criptografados (TLS)?
- [ ] Dados em repouso são criptografados?
- [ ] Acesso é controlado por autenticação e autorização?
- [ ] Operações sensíveis têm audit trail?

### 5. Compartilhamento externo
- [ ] Dados saem do sistema? (APIs, IA, integrações)
- [ ] Se sim: estão anonimizados conforme Tier 3?
- [ ] O contrato com o fornecedor externo cobre tratamento de dados sensíveis?
- [ ] Há DPA (Data Processing Agreement) com o fornecedor de LLM?
