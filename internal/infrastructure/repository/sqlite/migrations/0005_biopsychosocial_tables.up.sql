-- Migration: 0005_add_biopsychosocial_tables.up.sql
-- Description: Add patient_medications and patient_vitals tables for biological context tracking

-- Patient Medications Table
CREATE TABLE IF NOT EXISTS patient_medications (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    name TEXT NOT NULL,
    dosage TEXT,
    frequency TEXT,
    prescriber TEXT,
    status TEXT DEFAULT 'active' CHECK(status IN ('active', 'suspended', 'finished')),
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_medications_patient_id ON patient_medications(patient_id);
CREATE INDEX IF NOT EXISTS idx_medications_status ON patient_medications(status);
CREATE INDEX IF NOT EXISTS idx_medications_patient_status ON patient_medications(patient_id, status);

-- Patient Vitals Table (Time-series)
CREATE TABLE IF NOT EXISTS patient_vitals (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATE NOT NULL,
    sleep_hours REAL CHECK(sleep_hours >= 0 AND sleep_hours <= 24),
    appetite_level INTEGER CHECK(appetite_level >= 1 AND appetite_level <= 10),
    weight REAL,
    physical_activity INTEGER DEFAULT 0,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_vitals_patient_id ON patient_vitals(patient_id);
CREATE INDEX IF NOT EXISTS idx_vitals_date ON patient_vitals(patient_id, date DESC);
