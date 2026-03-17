-- Migration: 0003_patient_search_fts5_simple
-- Description: Remove suporte básico a Full Text Search (FTS5) para busca de pacientes
-- Date: 2025-03-17

-- Remove trigger
DROP TRIGGER IF EXISTS patients_ai_simple;

-- Remove tabela virtual FTS5
DROP TABLE IF EXISTS patients_fts;