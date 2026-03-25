-- Migration: 0010_add_patient_identity_fields
-- Description: Remove identity SOTA fields from patients table

-- Drop indexes first
DROP INDEX IF EXISTS idx_patients_gender;
DROP INDEX IF EXISTS idx_patients_ethnicity;
DROP INDEX IF EXISTS idx_patients_occupation;
DROP INDEX IF EXISTS idx_patients_education;

-- Remove columns
ALTER TABLE patients DROP COLUMN gender;
ALTER TABLE patients DROP COLUMN ethnicity;
ALTER TABLE patients DROP COLUMN occupation;
ALTER TABLE patients DROP COLUMN education;