-- +migrate Up
ALTER TABLE recordings ADD COLUMN session_id INTEGER REFERENCES sessions(id) ON DELETE SET NULL;

CREATE INDEX idx_recordings_session_id ON recordings(session_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_recordings_session_id;
ALTER TABLE recordings DROP COLUMN session_id;
