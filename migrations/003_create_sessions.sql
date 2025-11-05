-- +migrate Up
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    campaign_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    session_number INTEGER NOT NULL,
    session_date DATE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_campaign_id ON sessions(campaign_id);
CREATE INDEX idx_sessions_session_date ON sessions(session_date);
CREATE INDEX idx_sessions_session_number ON sessions(campaign_id, session_number);

-- +migrate Down
DROP INDEX IF EXISTS idx_sessions_session_number;
DROP INDEX IF EXISTS idx_sessions_session_date;
DROP INDEX IF EXISTS idx_sessions_campaign_id;
DROP TABLE IF EXISTS sessions;
