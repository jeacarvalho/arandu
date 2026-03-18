-- Migration: 0005_add_biopsychosocial_tables.down.sql
-- Description: Remove patient_medications and patient_vitals tables

DROP INDEX IF EXISTS idx_medications_patient_status;
DROP INDEX IF EXISTS idx_medications_status;
DROP INDEX IF EXISTS idx_medications_patient_id;
DROP TABLE IF EXISTS patient_medications;

DROP INDEX IF EXISTS idx_vitals_date;
DROP INDEX IF EXISTS idx_vitals_patient_id;
DROP TABLE IF EXISTS patient_vitals;
