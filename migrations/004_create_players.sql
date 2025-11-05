-- +migrate Up
CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT,
    character_name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_players_name ON players(name);
CREATE INDEX idx_players_email ON players(email);

-- Campaign players junction table
CREATE TABLE IF NOT EXISTS campaign_players (
    campaign_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (campaign_id, player_id),
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

CREATE INDEX idx_campaign_players_player_id ON campaign_players(player_id);

-- Session players junction table (for tracking attendance)
CREATE TABLE IF NOT EXISTS session_players (
    session_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    attended BOOLEAN NOT NULL DEFAULT 1,
    PRIMARY KEY (session_id, player_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

CREATE INDEX idx_session_players_player_id ON session_players(player_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_session_players_player_id;
DROP TABLE IF EXISTS session_players;

DROP INDEX IF EXISTS idx_campaign_players_player_id;
DROP TABLE IF EXISTS campaign_players;

DROP INDEX IF EXISTS idx_players_email;
DROP INDEX IF EXISTS idx_players_name;
DROP TABLE IF EXISTS players;
