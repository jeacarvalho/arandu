-- Migration: 0002_term_frequency_analysis
-- Description: Create infrastructure for term frequency analysis (without FTS5)
-- Created: 2026-03-17

-- Create table to store stop words (for filtering common words)
CREATE TABLE IF NOT EXISTS stop_words_portuguese (
    term TEXT PRIMARY KEY
);

-- Insert Portuguese stop words
INSERT OR IGNORE INTO stop_words_portuguese (term) VALUES 
    ('a'), ('à'), ('ao'), ('aos'), ('aquela'), ('aquelas'), ('aquele'), ('aqueles'), ('aquilo'), 
    ('as'), ('às'), ('até'), ('com'), ('como'), ('da'), ('das'), ('de'), ('dela'), ('delas'), 
    ('dele'), ('deles'), ('depois'), ('do'), ('dos'), ('e'), ('é'), ('ela'), ('elas'), ('ele'), 
    ('eles'), ('em'), ('entre'), ('era'), ('eram'), ('essa'), ('essas'), ('esse'), ('esses'), 
    ('esta'), ('está'), ('estão'), ('estas'), ('estava'), ('estavam'), ('este'), ('estes'), 
    ('esteve'), ('estive'), ('estivemos'), ('estiveram'), ('estivesse'), ('estivessem'), 
    ('estou'), ('eu'), ('foi'), ('fomos'), ('for'), ('fora'), ('foram'), ('fosse'), ('fossem'), 
    ('fui'), ('há'), ('havia'), ('haviam'), ('isso'), ('isto'), ('já'), ('lhe'), ('lhes'), 
    ('lo'), ('mais'), ('mas'), ('me'), ('mesmo'), ('meu'), ('meus'), ('minha'), ('minhas'), 
    ('muito'), ('muitos'), ('na'), ('não'), ('nas'), ('nem'), ('no'), ('nos'), ('nós'), 
    ('nossa'), ('nossas'), ('nosso'), ('nossos'), ('num'), ('numa'), ('o'), ('os'), ('ou'), 
    ('para'), ('pela'), ('pelas'), ('pelo'), ('pelos'), ('por'), ('qual'), ('quando'), 
    ('que'), ('quem'), ('são'), ('se'), ('seja'), ('sejam'), ('sem'), ('ser'), ('será'), 
    ('serão'), ('seria'), ('seriam'), ('seu'), ('seus'), ('só'), ('somos'), ('sou'), 
    ('sua'), ('suas'), ('também'), ('te'), ('tem'), ('tém'), ('tendo'), ('tenha'), ('tenham'), 
    ('temos'), ('tenho'), ('ter'), ('terá'), ('terão'), ('teria'), ('teriam'), ('teu'), 
    ('teus'), ('teve'), ('tinha'), ('tinham'), ('tive'), ('tivemos'), ('tiveram'), 
    ('tivesse'), ('tivessem'), ('tu'), ('tua'), ('tuas'), ('um'), ('uma'), ('umas'), 
    ('uns'), ('vai'), ('vamos'), ('vão'), ('você'), ('vocês'), ('vos'), ('vossa'), 
    ('vossas'), ('vosso'), ('vossos'), ('foi'), ('sido'), ('estou'), ('está'), ('estamos'), 
    ('estão'), ('estive'), ('esteve'), ('estivemos'), ('estiveram'), ('estava'), ('estávamos'),
    ('estavam'), ('estivera'), ('estivéramos'), ('seja'), ('sejamos'), ('sejam'), ('fosse'), 
    ('fôssemos'), ('fossem'), ('for'), ('formos'), ('forem'), ('for'), ('formos'), ('forem');

-- Create index for stop words performance
CREATE INDEX IF NOT EXISTS idx_stop_words_term ON stop_words_portuguese(term);

-- Create table to cache term frequencies for better performance
CREATE TABLE IF NOT EXISTS term_frequency_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    patient_id TEXT NOT NULL,
    term TEXT NOT NULL,
    source TEXT NOT NULL,
    frequency INTEGER NOT NULL,
    computed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(patient_id, term, source)
);

CREATE INDEX IF NOT EXISTS idx_term_cache_patient ON term_frequency_cache(patient_id);
CREATE INDEX IF NOT EXISTS idx_term_cache_term ON term_frequency_cache(term);
CREATE INDEX IF NOT EXISTS idx_term_cache_computed ON term_frequency_cache(computed_at);

-- Create view for term frequency (without FTS5)
-- This view counts full content entries per patient/source combination
-- For detailed word analysis, use the cache table
CREATE VIEW IF NOT EXISTS patient_content_count AS
SELECT 
    s.patient_id,
    'observation' AS source,
    COUNT(*) AS content_count,
    SUM(LENGTH(o.content)) AS total_chars
FROM observations o
JOIN sessions s ON s.id = o.session_id
GROUP BY s.patient_id
UNION ALL
SELECT 
    s.patient_id,
    'intervention' AS source,
    COUNT(*) AS content_count,
    SUM(LENGTH(i.content)) AS total_chars
FROM interventions i
JOIN sessions s ON s.id = i.session_id
GROUP BY s.patient_id;

-- Create function to compute term frequencies (simplified without FTS5)
-- This will be called from Go to populate the cache
CREATE VIEW IF NOT EXISTS term_analysis_ready AS
SELECT 
    s.patient_id,
    o.id AS observation_id,
    o.content
FROM observations o
JOIN sessions s ON s.id = o.session_id
UNION ALL
SELECT 
    s.patient_id,
    i.id AS observation_id,
    i.content
FROM interventions i
JOIN sessions s ON s.id = i.session_id;
