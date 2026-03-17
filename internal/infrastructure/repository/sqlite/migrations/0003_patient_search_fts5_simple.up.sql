-- Migration: 0003_patient_search_fts5_simple
-- Description: Adiciona suporte básico a Full Text Search (FTS5) para busca de pacientes
-- Date: 2025-03-17

-- Tabela virtual FTS5 simples para busca de pacientes
-- Não usa content= para evitar complexidade com triggers
CREATE VIRTUAL TABLE IF NOT EXISTS patients_fts USING fts5(
    patient_id UNINDEXED,
    name,
    notes
);

-- Popula a tabela FTS5 com dados existentes
INSERT INTO patients_fts(patient_id, name, notes)
SELECT id, name, notes FROM patients;

-- Trigger simples para novos pacientes
CREATE TRIGGER IF NOT EXISTS patients_ai_simple AFTER INSERT ON patients BEGIN
    INSERT INTO patients_fts(patient_id, name, notes) 
    VALUES (new.id, new.name, new.notes);
END;