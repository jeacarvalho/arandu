-- 0001_initial_central.down.sql
-- Rollback initial schema for Control Plane (Central DB)

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_tenants_status;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;