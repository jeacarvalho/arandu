-- Migration: 0007_add_therapeutic_goals.down.sql
-- Description: Rollback therapeutic_goals table

DROP INDEX IF EXISTS idx_goals_patient_status;
DROP INDEX IF EXISTS idx_goals_status;
DROP INDEX IF EXISTS idx_goals_patient_id;
DROP TABLE IF EXISTS therapeutic_goals;
