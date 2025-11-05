-- +migrate Up
CREATE TABLE IF NOT EXISTS recordings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_id TEXT NOT NULL UNIQUE,
    filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    duration_seconds INTEGER DEFAULT 0,
    file_size_bytes INTEGER DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'recording',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    transcription_status TEXT DEFAULT 'pending',
    notes TEXT
);

CREATE INDEX idx_recordings_status ON recordings(status);
CREATE INDEX idx_recordings_created_at ON recordings(created_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_recordings_created_at;
DROP INDEX IF EXISTS idx_recordings_status;
DROP TABLE IF EXISTS recordings;
