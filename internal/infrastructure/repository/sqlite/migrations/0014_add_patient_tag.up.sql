-- Migration: 0014_add_patient_tag
-- Description: Add clinical tag field to patients (e.g. "BURNOUT", "ANSIEDADE", "TRIAGEM")

ALTER TABLE patients ADD COLUMN tag TEXT;

CREATE INDEX IF NOT EXISTS idx_patients_tag ON patients(tag);
