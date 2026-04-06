-- Migration: 0012_add_intervention_tags
-- Description: Adds classification tags for interventions (many-to-many relationship)

-- Tags reference table for interventions (pre-defined tags)
CREATE TABLE IF NOT EXISTS intervention_tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tag_type TEXT NOT NULL CHECK (tag_type IN ('cognitive', 'behavioral', 'emotional', 'psychoeducation', 'narrative', 'body')),
    color TEXT NOT NULL,
    icon TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL
);

-- Intervention classifications table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS intervention_classifications (
    id TEXT PRIMARY KEY,
    intervention_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    intensity INTEGER CHECK (intensity >= 1 AND intensity <= 5),
    created_at DATETIME NOT NULL,
    FOREIGN KEY (intervention_id) REFERENCES interventions(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES intervention_tags(id) ON DELETE CASCADE,
    UNIQUE(intervention_id, tag_id)
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_intervention_classifications_intervention ON intervention_classifications(intervention_id);
CREATE INDEX IF NOT EXISTS idx_intervention_classifications_tag ON intervention_classifications(tag_id);
CREATE INDEX IF NOT EXISTS idx_intervention_tags_type ON intervention_tags(tag_type);

-- Insert pre-defined intervention tags
-- Técnica Cognitiva (cognitive)
INSERT INTO intervention_tags (id, name, tag_type, color, icon, sort_order, created_at) VALUES
('int-cog-01', 'Reestruturação cognitiva', 'cognitive', '#7C3AED', 'brain', 1, datetime('now')),
('int-cog-02', 'Questionamento socrático', 'cognitive', '#7C3AED', 'question-circle', 2, datetime('now')),
('int-cog-03', 'Identificação de distorções', 'cognitive', '#7C3AED', 'search', 3, datetime('now')),
('int-cog-04', 'Registro de pensamentos', 'cognitive', '#7C3AED', 'clipboard', 4, datetime('now')),
('int-cog-05', 'Experimento comportamental', 'cognitive', '#7C3AED', 'flask', 5, datetime('now')),
('int-cog-06', 'Técnicas de decatastrofização', 'cognitive', '#7C3AED', 'cloud', 6, datetime('now')),
('int-cog-07', 'Reframe positivo', 'cognitive', '#7C3AED', 'sync-alt', 7, datetime('now'));

-- Técnica Comportamental (behavioral)
INSERT INTO intervention_tags (id, name, tag_type, color, icon, sort_order, created_at) VALUES
('int-beh-01', 'Exposição gradual', 'behavioral', '#1D9E75', 'stairs', 1, datetime('now')),
('int-beh-02', 'Ativação comportamental', 'behavioral', '#1D9E75', 'running', 2, datetime('now')),
('int-beh-03', 'Treino de habilidades', 'behavioral', '#1D9E75', 'dumbbell', 3, datetime('now')),
('int-beh-04', 'Reforço positivo', 'behavioral', '#1D9E75', 'plus-circle', 4, datetime('now')),
('int-beh-05', 'Modelagem', 'behavioral', '#1D9E75', 'copy', 5, datetime('now')),
('int-beh-06', 'Contrabalanço comportamental', 'behavioral', '#1D9E75', 'balance-scale', 6, datetime('now')),
('int-beh-07', 'Encadeamento de respostas', 'behavioral', '#1D9E75', 'link', 7, datetime('now'));

-- Técnica Emocional (emotional)
INSERT INTO intervention_tags (id, name, tag_type, color, icon, sort_order, created_at) VALUES
('int-emo-01', 'Validação emocional', 'emotional', '#0F6E56', 'heart', 1, datetime('now')),
('int-emo-02', 'Expressão de sentimentos', 'emotional', '#0F6E56', 'comment-dots', 2, datetime('now')),
('int-emo-03', 'Regulação emocional', 'emotional', '#0F6E56', 'sliders-h', 3, datetime('now')),
('int-emo-04', 'Mindfulness emocional', 'emotional', '#0F6E56', 'om', 4, datetime('now')),
('int-emo-05', 'Processamento emocional', 'emotional', '#0F6E56', 'stream', 5, datetime('now')),
('int-emo-06', 'Tolerância à afetividade', 'emotional', '#0F6E56', 'shield-alt', 6, datetime('now'));

-- Psicoeducação (psychoeducation)
INSERT INTO intervention_tags (id, name, tag_type, color, icon, sort_order, created_at) VALUES
('int-psy-01', 'Explicação sobre transtorno', 'psychoeducation', '#F59E0B', 'book-medical', 1, datetime('now')),
('int-psy-02', 'Informações sobre medicação', 'psychoeducation', '#F59E0B', 'pills', 2, datetime('now')),
('int-psy-03', 'Orientação familiar', 'psychoeducation', '#F59E0B', 'users', 3, datetime('now')),
('int-psy-04', 'Prevenção de recaída', 'psychoeducation', '#F59E0B', 'umbrella', 4, datetime('now')),
('int-psy-05', 'Estratégias de coping', 'psychoeducation', '#F59E0B', 'toolbox', 5, datetime('now')),
('int-psy-06', 'Psicoeducação sobre sono', 'psychoeducation', '#F59E0B', 'bed', 6, datetime('now'));

-- Exploração Narrativa (narrative)
INSERT INTO intervention_tags (id, name, tag_type, color, icon, sort_order, created_at) VALUES
('int-nar-01', 'Externalização', 'narrative', '#3B82F6', 'external-link-alt', 1, datetime('now')),
('int-nar-02', 'Reautorização', 'narrative', '#3B82F6', 'pen-fancy', 2, datetime('now')),
('int-nar-03', 'Identificação de exceções', 'narrative', '#3B82F6', 'search-plus', 3, datetime('now')),
('int-nar-04', 'Perguntas circulares', 'narrative', '#3B82F6', 'sync', 4, datetime('now')),
('int-nar-05', 'Genograma', 'narrative', '#3B82F6', 'sitemap', 5, datetime('now')),
('int-nar-06', 'Cartas para o problema', 'narrative', '#3B82F6', 'envelope', 6, datetime('now'));

-- Intervenção Corporal (body)
INSERT INTO intervention_tags (id, name, tag_type, color, icon, sort_order, created_at) VALUES
('int-bod-01', 'Respiração diafragmática', 'body', '#DC2626', 'wind', 1, datetime('now')),
('int-bod-02', 'Relaxamento muscular', 'body', '#DC2626', 'spa', 2, datetime('now')),
('int-bod-03', 'Grounding', 'body', '#DC2626', 'anchor', 3, datetime('now')),
('int-bod-04', 'Técnicas de ancoragem', 'body', '#DC2626', 'thumbtack', 4, datetime('now')),
('int-bod-05', 'Consciência corporal', 'body', '#DC2626', 'body-static', 5, datetime('now')),
('int-bod-06', 'Alongamento guiado', 'body', '#DC2626', 'child', 6, datetime('now'));
