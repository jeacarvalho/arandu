-- Migration: 0008_add_goal_closure_fields.down.sql
-- Description: Rollback closure_note and closed_at columns

ALTER TABLE therapeutic_goals DROP COLUMN closed_at;
ALTER TABLE therapeutic_goals DROP COLUMN closure_note;
