-- Migration: 0011_add_observation_tags
-- Rollback: Remove observation tags tables

DROP TABLE IF EXISTS observation_tags;
DROP TABLE IF EXISTS tags;
