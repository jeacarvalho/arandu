-- Create patients table with notes column
CREATE TABLE IF NOT EXISTS patients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    notes TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    patient_id TEXT NOT NULL,
    date DATE NOT NULL,
    summary TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);

-- Create observations table
CREATE TABLE IF NOT EXISTS observations (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Insert test data
INSERT OR IGNORE INTO patients (id, name, notes, created_at, updated_at) 
VALUES ('test-patient-1', 'Paciente Teste', 'Notas do paciente teste', datetime('now'), datetime('now'));

INSERT OR IGNORE INTO sessions (id, patient_id, date, summary, created_at, updated_at)
VALUES ('test-session-1', 'test-patient-1', date('now'), 'Sessão inicial para teste', datetime('now'), datetime('now'));
