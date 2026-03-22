-- Migration: 0009_add_anamnesis.up.sql
-- Description: Add patient_anamnesis table for clinical history foundation

CREATE TABLE IF NOT EXISTS patient_anamnesis (
    patient_id TEXT PRIMARY KEY,
    chief_complaint TEXT,
    personal_history TEXT,
    family_history TEXT,
    mental_state_exam TEXT,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_anamnesis_patient_id ON patient_anamnesis(patient_id);
