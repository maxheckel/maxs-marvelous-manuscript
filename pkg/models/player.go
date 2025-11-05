package models

import "time"

type Player struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Email         *string   `json:"email,omitempty"`
	CharacterName *string   `json:"character_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreatePlayerParams struct {
	Name          string
	Email         *string
	CharacterName *string
}

type UpdatePlayerParams struct {
	Name          *string
	Email         *string
	CharacterName *string
}

// PlayerWithCampaigns includes campaign information
type PlayerWithCampaigns struct {
	Player
	Campaigns []Campaign `json:"campaigns"`
}
