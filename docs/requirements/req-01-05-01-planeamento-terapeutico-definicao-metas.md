REQ-01-05-01 — Planeamento Terapêutico e Definição de Metas

Identificação

ID: REQ-01-05-01

Capability: CAP-01-05 — Gestão de Plano Terapêutico

Vision: VISION-01 — Registro da Prática Clínica

Status: draft

História do utilizador

Como psicólogo clínico, quero definir um plano terapêutico com metas claras para cada paciente, para orientar a minha prática clínica e monitorizar o progresso de forma objetiva ao longo do tempo.

Contexto

O Arandu diferencia-se de um simples bloco de notas por ser um instrumento de reflexão estratégica. O Planeamento Terapêutico é o "mapa" desenhado após a avaliação inicial. Ele permite que o terapeuta saiba exatamente para onde está a conduzir o processo, evitando que as sessões se tornem apenas conversas informais.

Descrição funcional

O sistema deve permitir a gestão de um plano de metas vinculado ao paciente:

Metas Atómicas: Capacidade de criar objetivos específicos (ex: "Reduzir esquiva social", "Processar luto materno").

Estados da Meta: Em Progresso, Alcançada ou Arquivada.

Racional Clínico: Um campo de texto longo para a fundamentação teórica do plano (Ex: Formulação de Caso).

Visibilidade na Sessão: Durante o registro de uma nova sessão, o terapeuta deve poder visualizar as metas ativas num painel lateral sutil.

Interface (Padrão Arandu SOTA)

Estética: Lista de metas com estilo de "check-list" elegante.

Tipografia: - Título da Meta: Inter (Sans).

Descrição/Racional: Source Serif 4 (Serif) para imersão técnica.

Cores: Fundo branco (--arandu-paper) sobre o papel de seda (--arandu-bg).

Lógica Técnica (SQL Multi-tenant)

-- Tabela de Metas Terapêuticas
CREATE TABLE IF NOT EXISTS therapeutic_goals (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'in_progress', 
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);


Critérios de Aceitação

CA-01: O utilizador pode criar, editar e mudar o status de metas via HTMX.

CA-02: A mudança de status (ex: marcar como Alcançada) deve ter uma transição visual suave.

CA-03: O plano deve ser isolado no SQLite individual do utilizador.

CA-04: O campo de descrição deve suportar texto longo e ser renderizado em Source Serif 4.

Fora do escopo

Gráficos de progresso quantitativo.

Metas compartilhadas com o paciente (o Arandu é exclusivo para o terapeuta).