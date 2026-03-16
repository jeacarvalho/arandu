-- Migration: 0001_initial_schema
-- Description: Drops all tables and indexes

-- Drop indexes
DROP INDEX IF EXISTS idx_insights_created_at;
DROP INDEX IF EXISTS idx_interventions_session_id;
DROP INDEX IF EXISTS idx_observations_created_at;
DROP INDEX IF EXISTS idx_observations_session_id;
DROP INDEX IF EXISTS idx_sessions_date;
DROP INDEX IF EXISTS idx_sessions_patient_id;
DROP INDEX IF EXISTS idx_patients_name;
DROP INDEX IF EXISTS idx_patients_created_at;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS insights;
DROP TABLE IF EXISTS interventions;
DROP TABLE IF EXISTS observations;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS patients;