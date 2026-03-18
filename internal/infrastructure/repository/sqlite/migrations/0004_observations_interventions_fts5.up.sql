-- Migration: 0004_observations_interventions_fts5
-- Description: Adiciona suporte FTS5 para busca em observations e interventions
-- Date: 2025-03-17

-- ============================================
-- 1. TABELAS VIRTUAIS FTS5
-- ============================================

-- Tabela virtual FTS5 para observations
-- Usa uma coluna separada source_id para armazenar o ID TEXT
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(
    source_id UNINDEXED,
    content
);

-- Tabela virtual FTS5 para interventions
CREATE VIRTUAL TABLE IF NOT EXISTS interventions_fts USING fts5(
    source_id UNINDEXED,
    content
);

-- ============================================
-- 2. TRIGGERS DE SINCRONIZAÇÃO: OBSERVATIONS
-- ============================================

-- Sincroniza INSERT em observations
CREATE TRIGGER IF NOT EXISTS observations_ai AFTER INSERT ON observations BEGIN
  INSERT INTO observations_fts(source_id, content) VALUES (new.id, new.content);
END;

-- Sincroniza DELETE em observations
CREATE TRIGGER IF NOT EXISTS observations_ad AFTER DELETE ON observations BEGIN
  DELETE FROM observations_fts WHERE source_id = old.id;
END;

-- Sincroniza UPDATE em observations
CREATE TRIGGER IF NOT EXISTS observations_au AFTER UPDATE ON observations BEGIN
  UPDATE observations_fts SET content = new.content WHERE source_id = old.id;
END;

-- ============================================
-- 3. TRIGGERS DE SINCRONIZAÇÃO: INTERVENTIONS
-- ============================================

-- Sincroniza INSERT em interventions
CREATE TRIGGER IF NOT EXISTS interventions_ai AFTER INSERT ON interventions BEGIN
  INSERT INTO interventions_fts(source_id, content) VALUES (new.id, new.content);
END;

-- Sincroniza DELETE em interventions
CREATE TRIGGER IF NOT EXISTS interventions_ad AFTER DELETE ON interventions BEGIN
  DELETE FROM interventions_fts WHERE source_id = old.id;
END;

-- Sincroniza UPDATE em interventions
CREATE TRIGGER IF NOT EXISTS interventions_au AFTER UPDATE ON interventions BEGIN
  UPDATE interventions_fts SET content = new.content WHERE source_id = old.id;
END;

-- ============================================
-- 4. POPULAÇÃO INICIAL
-- ============================================
-- Para dados já existentes, insere no índice FTS5

INSERT INTO observations_fts(source_id, content)
SELECT id, content FROM observations;

INSERT INTO interventions_fts(source_id, content)
SELECT id, content FROM interventions;
