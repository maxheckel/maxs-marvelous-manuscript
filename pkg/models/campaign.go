package models

import "time"

type Campaign struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCampaignParams struct {
	Name        string
	Description *string
}

type UpdateCampaignParams struct {
	Name        *string
	Description *string
}

// CampaignWithPlayers includes player information
type CampaignWithPlayers struct {
	Campaign
	Players []Player `json:"players"`
}

// CampaignWithSessions includes session information
type CampaignWithSessions struct {
	Campaign
	Sessions []Session `json:"sessions"`
}
