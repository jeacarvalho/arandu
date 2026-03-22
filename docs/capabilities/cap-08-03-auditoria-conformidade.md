CAP-08-03 — Auditoria e Conformidade

Identificação

ID: CAP-08-03

Título: Auditoria e Conformidade

Vision: VISION-07 — Organização Operacional do Consultório

Status: draft

Descrição

Capacidade de garantir a integridade, a soberania e a segurança jurídica dos dados clínicos. Esta capacidade foca no registro inalterável de acessos e operações, garantindo que o sistema esteja em total conformidade com a LGPD e normas éticas profissionais.

Objetivos

Trilha de Auditoria: Registrar quem, quando e o que foi acessado no Control Plane.

Imutabilidade: Garantir que registros de auditoria não possam ser apagados ou editados.

Transparência de Acesso: Permitir ao administrador e ao profissional a verificação de acessos indevidos.

Conformidade Legal: Prover evidências técnicas para auditorias de conselhos de classe ou planos de saúde.

Contexto Técnico (SOTA)

Diferente dos logs de sistema (técnicos), a Auditoria de Conformidade é um log de negócio que reside no Banco Central (arandu_central.db). Ela protege o profissional ao provar a soberania e o sigilo sobre os arquivos SQLite individuais.

Requisitos Relacionados

REQ-08-03-01: Auditoria de Acessos Clínicos.

REQ-08-01-02: Serviço de Auditoria Centralizada.

Valor de Negócio / Clínico

Reduz riscos jurídicos para o profissional e para a plataforma. Aumenta a confiança do psicólogo de que seus dados estão protegidos contra acessos não autorizados por parte da equipe de infraestrutura ou terceiros.

Critérios de Sucesso

Toda visualização de prontuário gera uma entrada na tabela audit_logs.

O log de auditoria inclui IP e User-Agent.

A gravação é assíncrona e não interfere na performance da interface.