-- Migration: 0004_observations_interventions_fts5
-- Description: Remove suporte FTS5 para observations e interventions
-- Date: 2025-03-17

-- Remove triggers de observations
DROP TRIGGER IF EXISTS observations_ai;
DROP TRIGGER IF EXISTS observations_ad;
DROP TRIGGER IF EXISTS observations_au;

-- Remove triggers de interventions
DROP TRIGGER IF EXISTS interventions_ai;
DROP TRIGGER IF EXISTS interventions_ad;
DROP TRIGGER IF EXISTS interventions_au;

-- Remove tabelas virtuais FTS5
DROP TABLE IF EXISTS observations_fts;
DROP TABLE IF EXISTS interventions_fts;
