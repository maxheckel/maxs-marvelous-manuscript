package db

import (
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/model"
	. "github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/table"
	"github.com/maxheckel/maxs-marvelous-manuscript/pkg/models"
)

type CampaignRepository struct {
	db *DB
}

func NewCampaignRepository(db *DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

// Create creates a new campaign
func (r *CampaignRepository) Create(params models.CreateCampaignParams) (*models.Campaign, error) {
	jetModel := model.Campaigns{
		Name:        params.Name,
		Description: params.Description,
	}

	stmt := Campaigns.
		INSERT(Campaigns.Name, Campaigns.Description).
		MODEL(jetModel).
		RETURNING(Campaigns.AllColumns)

	var dest model.Campaigns
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return jetModelToCampaign(&dest), nil
}

// GetByID retrieves a campaign by ID
func (r *CampaignRepository) GetByID(id int64) (*models.Campaign, error) {
	stmt := SELECT(Campaigns.AllColumns).
		FROM(Campaigns).
		WHERE(Campaigns.ID.EQ(Int32(int32(id))))

	var dest model.Campaigns
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	return jetModelToCampaign(&dest), nil
}

// List retrieves all campaigns
func (r *CampaignRepository) List() ([]*models.Campaign, error) {
	stmt := SELECT(Campaigns.AllColumns).
		FROM(Campaigns).
		ORDER_BY(Campaigns.CreatedAt.DESC())

	var dest []model.Campaigns
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %w", err)
	}

	campaigns := make([]*models.Campaign, len(dest))
	for i, d := range dest {
		campaigns[i] = jetModelToCampaign(&d)
	}

	return campaigns, nil
}

// Update updates a campaign
func (r *CampaignRepository) Update(id int64, params models.UpdateCampaignParams) error {
	stmt := Campaigns.UPDATE().
		SET(Campaigns.UpdatedAt.SET(CURRENT_TIMESTAMP()))

	if params.Name != nil {
		stmt = stmt.SET(Campaigns.Name.SET(String(*params.Name)))
	}
	if params.Description != nil {
		stmt = stmt.SET(Campaigns.Description.SET(String(*params.Description)))
	}

	stmt = stmt.WHERE(Campaigns.ID.EQ(Int32(int32(id))))

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to update campaign: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("campaign not found")
	}

	return nil
}

// Delete deletes a campaign
func (r *CampaignRepository) Delete(id int64) error {
	stmt := Campaigns.
		DELETE().
		WHERE(Campaigns.ID.EQ(Int32(int32(id))))

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to delete campaign: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("campaign not found")
	}

	return nil
}

// AddPlayer adds a player to a campaign
func (r *CampaignRepository) AddPlayer(campaignID, playerID int64) error {
	jetModel := model.CampaignPlayers{
		CampaignID: int32(campaignID),
		PlayerID:   int32(playerID),
	}

	stmt := CampaignPlayers.
		INSERT(CampaignPlayers.CampaignID, CampaignPlayers.PlayerID).
		MODEL(jetModel)

	_, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to add player to campaign: %w", err)
	}

	return nil
}

// RemovePlayer removes a player from a campaign
func (r *CampaignRepository) RemovePlayer(campaignID, playerID int64) error {
	stmt := CampaignPlayers.
		DELETE().
		WHERE(
			CampaignPlayers.CampaignID.EQ(Int32(int32(campaignID))).
				AND(CampaignPlayers.PlayerID.EQ(Int32(int32(playerID)))),
		)

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to remove player from campaign: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player not found in campaign")
	}

	return nil
}

// GetPlayers retrieves all players in a campaign
func (r *CampaignRepository) GetPlayers(campaignID int64) ([]*models.Player, error) {
	stmt := SELECT(Players.AllColumns).
		FROM(
			Players.
				INNER_JOIN(CampaignPlayers, CampaignPlayers.PlayerID.EQ(Players.ID)),
		).
		WHERE(CampaignPlayers.CampaignID.EQ(Int32(int32(campaignID)))).
		ORDER_BY(CampaignPlayers.JoinedAt.ASC())

	var dest []model.Players
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign players: %w", err)
	}

	players := make([]*models.Player, len(dest))
	for i, d := range dest {
		players[i] = jetModelToPlayer(&d)
	}

	return players, nil
}

// GetWithPlayers retrieves a campaign with all its players
func (r *CampaignRepository) GetWithPlayers(campaignID int64) (*models.CampaignWithPlayers, error) {
	campaign, err := r.GetByID(campaignID)
	if err != nil {
		return nil, err
	}

	players, err := r.GetPlayers(campaignID)
	if err != nil {
		return nil, err
	}

	result := &models.CampaignWithPlayers{
		Campaign: *campaign,
		Players:  make([]models.Player, len(players)),
	}

	for i, p := range players {
		result.Players[i] = *p
	}

	return result, nil
}

// Helper function to convert Jet model to our domain model
func jetModelToCampaign(m *model.Campaigns) *models.Campaign {
	campaign := &models.Campaign{
		ID:        int64(*m.ID),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.Description != nil {
		campaign.Description = m.Description
	}

	return campaign
}

// Helper function to convert Jet model to our domain model
func jetModelToPlayer(m *model.Players) *models.Player {
	player := &models.Player{
		ID:        int64(*m.ID),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}

	if m.Email != nil {
		player.Email = m.Email
	}
	if m.CharacterName != nil {
		player.CharacterName = m.CharacterName
	}

	return player
}
