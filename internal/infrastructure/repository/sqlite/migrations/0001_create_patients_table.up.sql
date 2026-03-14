-- Migration: 0001_create_patients_table
-- Description: Creates the patients table for storing patient information

CREATE TABLE IF NOT EXISTS patients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Create index on created_at for faster sorting (if not exists)
CREATE INDEX IF NOT EXISTS idx_patients_created_at ON patients(created_at DESC);

-- Create index on name for faster searching (if not exists)
CREATE INDEX IF NOT EXISTS idx_patients_name ON patients(name);