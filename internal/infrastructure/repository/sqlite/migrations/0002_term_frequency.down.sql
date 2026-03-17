-- Migration: 0002_term_frequency_analysis
-- Description: Rollback term frequency analysis
-- Created: 2026-03-17

DROP VIEW IF EXISTS term_analysis_ready;
DROP VIEW IF EXISTS patient_content_count;
DROP TABLE IF EXISTS term_frequency_cache;
DROP TABLE IF EXISTS stop_words_portuguese;
