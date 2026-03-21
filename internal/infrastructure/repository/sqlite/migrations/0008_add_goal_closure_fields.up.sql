-- Migration: 0008_add_goal_closure_fields.up.sql
-- Description: Add closure_note and closed_at fields to therapeutic_goals

ALTER TABLE therapeutic_goals ADD COLUMN closure_note TEXT;
ALTER TABLE therapeutic_goals ADD COLUMN closed_at DATETIME;
