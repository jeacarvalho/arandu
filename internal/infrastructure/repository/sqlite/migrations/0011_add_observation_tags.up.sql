-- Migration: 0011_add_observation_tags
-- Description: Adds classification tags for observations (many-to-many relationship)

-- Tags reference table (pre-defined tags)
CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tag_type TEXT NOT NULL CHECK (tag_type IN ('emotion', 'behavior', 'cognition', 'relationship', 'somatic', 'context')),
    color TEXT NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL
);

-- Observation tags table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS observation_tags (
    id TEXT PRIMARY KEY,
    observation_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    intensity INTEGER CHECK (intensity >= 1 AND intensity <= 5),
    created_at DATETIME NOT NULL,
    FOREIGN KEY (observation_id) REFERENCES observations(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE(observation_id, tag_id)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_observation_tags_observation_id ON observation_tags(observation_id);
CREATE INDEX IF NOT EXISTS idx_observation_tags_tag_id ON observation_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_tags_type ON tags(tag_type);

-- Insert pre-defined tags
INSERT INTO tags (id, name, tag_type, color, sort_order, created_at) VALUES
    ('tag-emotion-01', 'Ansiedade', 'emotion', '#0F6E56', 1, datetime('now')),
    ('tag-emotion-02', 'Tristeza', 'emotion', '#0F6E56', 2, datetime('now')),
    ('tag-emotion-03', 'Raiva', 'emotion', '#0F6E56', 3, datetime('now')),
    ('tag-emotion-04', 'Alegria', 'emotion', '#0F6E56', 4, datetime('now')),
    ('tag-emotion-05', 'Medo', 'emotion', '#0F6E56', 5, datetime('now')),
    ('tag-emotion-06', 'Frustração', 'emotion', '#0F6E56', 6, datetime('now')),
    ('tag-behavior-01', 'Evitação', 'behavior', '#1D9E75', 1, datetime('now')),
    ('tag-behavior-02', 'Confronto', 'behavior', '#1D9E75', 2, datetime('now')),
    ('tag-behavior-03', 'Isolamento', 'behavior', '#1D9E75', 3, datetime('now')),
    ('tag-behavior-04', 'Hiperatividade', 'behavior', '#1D9E75', 4, datetime('now')),
    ('tag-behavior-05', 'Impassividade', 'behavior', '#1D9E75', 5, datetime('now')),
    ('tag-cognition-01', 'Pensamento catastrófico', 'cognition', '#7C3AED', 1, datetime('now')),
    ('tag-cognition-02', 'Perfeccionismo', 'cognition', '#7C3AED', 2, datetime('now')),
    ('tag-cognition-03', 'Ruminar', 'cognition', '#7C3AED', 3, datetime('now')),
    ('tag-cognition-04', 'Distorção cognitiva', 'cognition', '#7C3AED', 4, datetime('now')),
    ('tag-cognition-05', 'Insight', 'cognition', '#7C3AED', 5, datetime('now')),
    ('tag-relationship-01', 'Conflito familiar', 'relationship', '#F59E0B', 1, datetime('now')),
    ('tag-relationship-02', 'Dificuldade social', 'relationship', '#F59E0B', 2, datetime('now')),
    ('tag-relationship-03', 'Limiares', 'relationship', '#F59E0B', 3, datetime('now')),
    ('tag-relationship-04', 'Vínculo terapêutico', 'relationship', '#F59E0B', 4, datetime('now')),
    ('tag-somatic-01', 'Tensão muscular', 'somatic', '#DC2626', 1, datetime('now')),
    ('tag-somatic-02', 'Insônia', 'somatic', '#DC2626', 2, datetime('now')),
    ('tag-somatic-03', 'Sintomas físicos', 'somatic', '#DC2626', 3, datetime('now')),
    ('tag-somatic-04', 'Mobilidade', 'somatic', '#DC2626', 4, datetime('now')),
    ('tag-context-01', 'Evento recente', 'context', '#6B7280', 1, datetime('now')),
    ('tag-context-02', 'Transição de vida', 'context', '#6B7280', 2, datetime('now')),
    ('tag-context-03', 'Estresse ocupacional', 'context', '#6B7280', 3, datetime('now')),
    ('tag-context-04', 'Crise', 'context', '#6B7280', 4, datetime('now'));
