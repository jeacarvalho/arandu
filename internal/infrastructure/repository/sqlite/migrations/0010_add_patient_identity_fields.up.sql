-- Migration: 0010_add_patient_identity_fields
-- Description: Add identity SOTA fields to patients table (gender, ethnicity, occupation, education)

-- Add new identity fields to patients table
ALTER TABLE patients ADD COLUMN gender TEXT;
ALTER TABLE patients ADD COLUMN ethnicity TEXT;
ALTER TABLE patients ADD COLUMN occupation TEXT;
ALTER TABLE patients ADD COLUMN education TEXT;

-- Create indexes for the new fields
CREATE INDEX IF NOT EXISTS idx_patients_gender ON patients(gender);
CREATE INDEX IF NOT EXISTS idx_patients_ethnicity ON patients(ethnicity);
CREATE INDEX IF NOT EXISTS idx_patients_occupation ON patients(occupation);
CREATE INDEX IF NOT EXISTS idx_patients_education ON patients(education);