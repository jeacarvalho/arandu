-- Rollback: 0014_add_patient_tag
-- SQLite does not support DROP COLUMN before 3.35.0; recreate table without tag

CREATE TABLE patients_backup AS SELECT id, name, gender, ethnicity, occupation, education, notes, created_at, updated_at FROM patients;
DROP TABLE patients;
ALTER TABLE patients_backup RENAME TO patients;
