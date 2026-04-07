-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
    id TEXT PRIMARY KEY,
    patient_id TEXT,
    patient_name TEXT,
    date DATE NOT NULL,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    duration INTEGER NOT NULL DEFAULT 50,
    type TEXT NOT NULL DEFAULT 'session',
    status TEXT NOT NULL DEFAULT 'scheduled',
    notes TEXT,
    session_id TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    cancelled_at DATETIME,
    completed_at DATETIME,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL
);

-- Indexes for performance
CREATE INDEX idx_appointments_date ON appointments(date);
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_appointments_date_range ON appointments(date, start_time, end_time);
CREATE INDEX idx_appointments_session_id ON appointments(session_id);

-- Create agenda_settings table for therapist configuration
CREATE TABLE IF NOT EXISTS agenda_settings (
    user_id TEXT PRIMARY KEY,
    slot_duration INTEGER DEFAULT 50,
    work_start_time TEXT DEFAULT '08:00',
    work_end_time TEXT DEFAULT '18:00',
    work_days TEXT DEFAULT '1,2,3,4,5',
    break_between_slots INTEGER DEFAULT 10,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
