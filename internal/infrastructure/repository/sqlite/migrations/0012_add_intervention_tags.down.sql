-- Migration: 0012_add_intervention_tags
-- Rollback: Remove intervention tags tables

DROP TABLE IF EXISTS intervention_classifications;
DROP TABLE IF EXISTS intervention_tags;
