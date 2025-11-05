package models

import "time"

type Session struct {
	ID            int64      `json:"id"`
	CampaignID    int64      `json:"campaign_id"`
	Name          string     `json:"name"`
	SessionNumber int        `json:"session_number"`
	SessionDate   *time.Time `json:"session_date,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CreateSessionParams struct {
	CampaignID    int64
	Name          string
	SessionNumber int
	SessionDate   *time.Time
	Notes         *string
}

type UpdateSessionParams struct {
	Name          *string
	SessionNumber *int
	SessionDate   *time.Time
	Notes         *string
}

// SessionWithRecordings includes recording information
type SessionWithRecordings struct {
	Session
	Recordings []Recording `json:"recordings"`
}

// SessionWithPlayers includes player attendance information
type SessionWithPlayers struct {
	Session
	Players []PlayerAttendance `json:"players"`
}

// SessionWithDetails includes everything
type SessionWithDetails struct {
	Session
	Campaign   *Campaign          `json:"campaign,omitempty"`
	Recordings []Recording        `json:"recordings"`
	Players    []PlayerAttendance `json:"players"`
}

// PlayerAttendance represents a player's attendance at a session
type PlayerAttendance struct {
	Player
	Attended bool `json:"attended"`
}
