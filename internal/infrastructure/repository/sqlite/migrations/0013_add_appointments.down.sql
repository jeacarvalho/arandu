-- Drop indexes
DROP INDEX IF EXISTS idx_appointments_date;
DROP INDEX IF EXISTS idx_appointments_patient_id;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_date_range;
DROP INDEX IF EXISTS idx_appointments_session_id;

-- Drop tables
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS agenda_settings;
