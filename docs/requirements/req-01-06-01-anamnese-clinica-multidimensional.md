REQ-01-06-01 — Anamnese Clínica Multidimensional

Identificação

ID: REQ-01-06-01

Capability: CAP-01-06 — Avaliação Inicial e Anamnese

Vision: VISION-01 — Registro da Prática Clínica

Status: draft

História do utilizador

Como psicólogo clínico, quero registrar a anamnese completa do paciente em campos estruturados (história familiar, desenvolvimento, queixa principal), para ter um mapa claro da subjetividade do sujeito desde o início do tratamento e facilitar consultas rápidas durante o processo terapêutico.

Contexto

Diferente das notas de sessão (que são temporais), a anamnese é a fundação estrutural do caso. No Arandu, ela não deve ser um formulário estático chato, mas uma "Árvore de História" que pode ser preenchida gradualmente (Progressive Disclosure).

Descrição funcional

O sistema deve prover um módulo de Anamnese com as seguintes dimensões:

Queixa Principal: Texto livre em Source Serif 4.

História Pessoal e Social: Relacionamentos, escolaridade, lazer.

História Familiar: Composição familiar e dinâmicas (Base para o Genograma futuro).

Exame das Funções Mentais: Estado atual (atenção, humor, sono, apetite - integrado com REQ-01-04-01).

Interface (Padrão Arandu SOTA)

Navegação Lateral Interna: Sub-menu no perfil do paciente: "Dados Cadastrais" | "Anamnese" | "Prontuário".

Silent Input: Blocos de escrita fluida. Ao terminar um bloco, o sistema salva automaticamente via HTMX.

Visual: Seções expansíveis (Accordion) para não sobrecarregar o olhar do terapeuta.

Lógica Técnica (SQL)

CREATE TABLE IF NOT EXISTS patient_anamnesis (
    patient_id TEXT PRIMARY KEY,
    chief_complaint TEXT,           -- Queixa Principal
    personal_history TEXT,          -- História Pessoal
    family_history TEXT,            -- História Familiar
    mental_state_exam TEXT,         -- Exame das funções mentais
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);


Critérios de Aceitação

CA-01: O utilizador pode navegar para a anamnese a partir do perfil do paciente.

CA-02: Todo o texto clínico deve obrigatoriamente usar Source Serif 4.

CA-03: O salvamento deve ser atómico por seção via HTMX (Silent Save).

CA-04: Os dados devem estar isolados no SQLite do tenant.

Nota de Soberania: Este dado reside no SQLite individual do profissional (clinical_{user_uuid}.db)