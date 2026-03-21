-- Migration: 0007_add_therapeutic_goals.up.sql
-- Description: Add therapeutic_goals table for clinical planning and goal tracking

CREATE TABLE IF NOT EXISTS therapeutic_goals (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'in_progress' CHECK(status IN ('in_progress', 'achieved', 'archived')),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_goals_patient_id ON therapeutic_goals(patient_id);
CREATE INDEX IF NOT EXISTS idx_goals_status ON therapeutic_goals(status);
CREATE INDEX IF NOT EXISTS idx_goals_patient_status ON therapeutic_goals(patient_id, status);
