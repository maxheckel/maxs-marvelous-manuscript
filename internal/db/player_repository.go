package db

import (
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/model"
	. "github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/table"
	"github.com/maxheckel/maxs-marvelous-manuscript/pkg/models"
)

type PlayerRepository struct {
	db *DB
}

func NewPlayerRepository(db *DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

// Create creates a new player
func (r *PlayerRepository) Create(params models.CreatePlayerParams) (*models.Player, error) {
	jetModel := model.Players{
		Name:          params.Name,
		Email:         params.Email,
		CharacterName: params.CharacterName,
	}

	stmt := Players.
		INSERT(Players.MutableColumns).
		MODEL(jetModel).
		RETURNING(Players.AllColumns)

	var dest model.Players
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return jetModelToPlayer(&dest), nil
}

// GetByID retrieves a player by ID
func (r *PlayerRepository) GetByID(id int64) (*models.Player, error) {
	stmt := SELECT(Players.AllColumns).
		FROM(Players).
		WHERE(Players.ID.EQ(Int32(int32(id))))

	var dest model.Players
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return jetModelToPlayer(&dest), nil
}

// GetByEmail retrieves a player by email
func (r *PlayerRepository) GetByEmail(email string) (*models.Player, error) {
	stmt := SELECT(Players.AllColumns).
		FROM(Players).
		WHERE(Players.Email.EQ(String(email)))

	var dest model.Players
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return jetModelToPlayer(&dest), nil
}

// List retrieves all players
func (r *PlayerRepository) List() ([]*models.Player, error) {
	stmt := SELECT(Players.AllColumns).
		FROM(Players).
		ORDER_BY(Players.Name.ASC())

	var dest []model.Players
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to list players: %w", err)
	}

	players := make([]*models.Player, len(dest))
	for i, d := range dest {
		players[i] = jetModelToPlayer(&d)
	}

	return players, nil
}

// Update updates a player
func (r *PlayerRepository) Update(id int64, params models.UpdatePlayerParams) error {
	columns := ColumnList{}
	values := make(map[Column]interface{})

	if params.Name != nil {
		columns = append(columns, Players.Name)
		values[Players.Name] = *params.Name
	}
	if params.Email != nil {
		columns = append(columns, Players.Email)
		values[Players.Email] = params.Email
	}
	if params.CharacterName != nil {
		columns = append(columns, Players.CharacterName)
		values[Players.CharacterName] = params.CharacterName
	}

	if len(columns) == 0 {
		return fmt.Errorf("no fields to update")
	}

	stmt := Players.
		UPDATE(columns).
		MODEL(values).
		WHERE(Players.ID.EQ(Int32(int32(id))))

	_, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	return nil
}

// Delete deletes a player
func (r *PlayerRepository) Delete(id int64) error {
	stmt := Players.
		DELETE().
		WHERE(Players.ID.EQ(Int32(int32(id))))

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player not found")
	}

	return nil
}

// GetCampaigns retrieves all campaigns a player is part of
func (r *PlayerRepository) GetCampaigns(playerID int64) ([]*models.Campaign, error) {
	stmt := SELECT(Campaigns.AllColumns).
		FROM(
			Campaigns.
				INNER_JOIN(CampaignPlayers, CampaignPlayers.CampaignID.EQ(Campaigns.ID)),
		).
		WHERE(CampaignPlayers.PlayerID.EQ(Int32(int32(playerID)))).
		ORDER_BY(CampaignPlayers.JoinedAt.DESC())

	var dest []model.Campaigns
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get player campaigns: %w", err)
	}

	campaigns := make([]*models.Campaign, len(dest))
	for i, d := range dest {
		campaigns[i] = jetModelToCampaign(&d)
	}

	return campaigns, nil
}

// GetSessions retrieves all sessions a player attended
func (r *PlayerRepository) GetSessions(playerID int64) ([]*models.Session, error) {
	stmt := SELECT(Sessions.AllColumns).
		FROM(
			Sessions.
				INNER_JOIN(SessionPlayers, SessionPlayers.SessionID.EQ(Sessions.ID)),
		).
		WHERE(
			SessionPlayers.PlayerID.EQ(Int32(int32(playerID))).
				AND(SessionPlayers.Attended.EQ(Bool(true))),
		).
		ORDER_BY(Sessions.SessionDate.DESC())

	var dest []model.Sessions
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get player sessions: %w", err)
	}

	sessions := make([]*models.Session, len(dest))
	for i, d := range dest {
		sessions[i] = jetModelToSession(&d)
	}

	return sessions, nil
}

// GetWithCampaigns retrieves a player with all their campaigns
func (r *PlayerRepository) GetWithCampaigns(playerID int64) (*models.PlayerWithCampaigns, error) {
	player, err := r.GetByID(playerID)
	if err != nil {
		return nil, err
	}

	campaigns, err := r.GetCampaigns(playerID)
	if err != nil {
		return nil, err
	}

	result := &models.PlayerWithCampaigns{
		Player:    *player,
		Campaigns: make([]models.Campaign, len(campaigns)),
	}

	for i, c := range campaigns {
		result.Campaigns[i] = *c
	}

	return result, nil
}
