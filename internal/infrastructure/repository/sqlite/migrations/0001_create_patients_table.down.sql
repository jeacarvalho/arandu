-- Migration: 0001_create_patients_table
-- Description: Drops the patients table and its indexes

DROP INDEX IF EXISTS idx_patients_name;
DROP INDEX IF EXISTS idx_patients_created_at;
DROP TABLE IF EXISTS patients;