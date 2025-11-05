package db

import (
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/model"
	. "github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/table"
	"github.com/maxheckel/maxs-marvelous-manuscript/pkg/models"
)

type SessionRepository struct {
	db *DB
}

func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(params models.CreateSessionParams) (*models.Session, error) {
	jetModel := model.Sessions{
		CampaignID:    int32(params.CampaignID),
		Name:          params.Name,
		SessionNumber: int32(params.SessionNumber),
		SessionDate:   params.SessionDate,
		Notes:         params.Notes,
	}

	stmt := Sessions.
		INSERT(Sessions.MutableColumns).
		MODEL(jetModel).
		RETURNING(Sessions.AllColumns)

	var dest model.Sessions
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return jetModelToSession(&dest), nil
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(id int64) (*models.Session, error) {
	stmt := SELECT(Sessions.AllColumns).
		FROM(Sessions).
		WHERE(Sessions.ID.EQ(Int32(int32(id))))

	var dest model.Sessions
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return jetModelToSession(&dest), nil
}

// ListByCampaign retrieves all sessions for a campaign
func (r *SessionRepository) ListByCampaign(campaignID int64) ([]*models.Session, error) {
	stmt := SELECT(Sessions.AllColumns).
		FROM(Sessions).
		WHERE(Sessions.CampaignID.EQ(Int32(int32(campaignID)))).
		ORDER_BY(Sessions.SessionNumber.ASC())

	var dest []model.Sessions
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	sessions := make([]*models.Session, len(dest))
	for i, d := range dest {
		sessions[i] = jetModelToSession(&d)
	}

	return sessions, nil
}

// List retrieves all sessions
func (r *SessionRepository) List() ([]*models.Session, error) {
	stmt := SELECT(Sessions.AllColumns).
		FROM(Sessions).
		ORDER_BY(Sessions.SessionDate.DESC(), Sessions.CreatedAt.DESC())

	var dest []model.Sessions
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	sessions := make([]*models.Session, len(dest))
	for i, d := range dest {
		sessions[i] = jetModelToSession(&d)
	}

	return sessions, nil
}

// Update updates a session
func (r *SessionRepository) Update(id int64, params models.UpdateSessionParams) error {
	columns := ColumnList{}
	values := make(map[Column]interface{})

	if params.Name != nil {
		columns = append(columns, Sessions.Name)
		values[Sessions.Name] = *params.Name
	}
	if params.SessionNumber != nil {
		sessionNumber := int32(*params.SessionNumber)
		columns = append(columns, Sessions.SessionNumber)
		values[Sessions.SessionNumber] = sessionNumber
	}
	if params.SessionDate != nil {
		columns = append(columns, Sessions.SessionDate)
		values[Sessions.SessionDate] = params.SessionDate
	}
	if params.Notes != nil {
		columns = append(columns, Sessions.Notes)
		values[Sessions.Notes] = params.Notes
	}

	if len(columns) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Always update the updated_at timestamp
	columns = append(columns, Sessions.UpdatedAt)
	values[Sessions.UpdatedAt] = CURRENT_TIMESTAMP()

	stmt := Sessions.
		UPDATE(columns).
		MODEL(values).
		WHERE(Sessions.ID.EQ(Int32(int32(id))))

	_, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// Delete deletes a session
func (r *SessionRepository) Delete(id int64) error {
	stmt := Sessions.
		DELETE().
		WHERE(Sessions.ID.EQ(Int32(int32(id))))

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// AddPlayer adds a player to a session with attendance tracking
func (r *SessionRepository) AddPlayer(sessionID, playerID int64, attended bool) error {
	jetModel := model.SessionPlayers{
		SessionID: int32(sessionID),
		PlayerID:  int32(playerID),
		Attended:  attended,
	}

	stmt := SessionPlayers.
		INSERT(SessionPlayers.AllColumns).
		MODEL(jetModel).
		ON_CONFLICT(SessionPlayers.SessionID, SessionPlayers.PlayerID).
		DO_UPDATE(
			SET(SessionPlayers.Attended.SET(Bool(attended))),
		)

	_, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to add player to session: %w", err)
	}

	return nil
}

// RemovePlayer removes a player from a session
func (r *SessionRepository) RemovePlayer(sessionID, playerID int64) error {
	stmt := SessionPlayers.
		DELETE().
		WHERE(
			SessionPlayers.SessionID.EQ(Int32(int32(sessionID))).
				AND(SessionPlayers.PlayerID.EQ(Int32(int32(playerID)))),
		)

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to remove player from session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player not found in session")
	}

	return nil
}

// GetPlayers retrieves all players in a session with attendance info
func (r *SessionRepository) GetPlayers(sessionID int64) ([]models.PlayerAttendance, error) {
	stmt := SELECT(
		Players.AllColumns,
		SessionPlayers.Attended,
	).FROM(
		Players.
			INNER_JOIN(SessionPlayers, SessionPlayers.PlayerID.EQ(Players.ID)),
	).WHERE(SessionPlayers.SessionID.EQ(Int32(int32(sessionID)))).
		ORDER_BY(Players.Name.ASC())

	type Result struct {
		model.Players
		Attended bool
	}

	var dest []Result
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get session players: %w", err)
	}

	players := make([]models.PlayerAttendance, len(dest))
	for i, d := range dest {
		player := jetModelToPlayer(&d.Players)
		players[i] = models.PlayerAttendance{
			Player:   *player,
			Attended: d.Attended,
		}
	}

	return players, nil
}

// GetWithDetails retrieves a session with all related data
func (r *SessionRepository) GetWithDetails(sessionID int64) (*models.SessionWithDetails, error) {
	session, err := r.GetByID(sessionID)
	if err != nil {
		return nil, err
	}

	// Get campaign
	var campaign model.Campaigns
	campaignStmt := SELECT(Campaigns.AllColumns).
		FROM(Campaigns).
		WHERE(Campaigns.ID.EQ(Int32(int32(session.CampaignID))))

	err = campaignStmt.Query(r.db.DB, &campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}

	// Get recordings
	var recordings []model.Recordings
	recordingsStmt := SELECT(Recordings.AllColumns).
		FROM(Recordings).
		WHERE(Recordings.SessionID.EQ(Int32(int32(sessionID)))).
		ORDER_BY(Recordings.CreatedAt.ASC())

	err = recordingsStmt.Query(r.db.DB, &recordings)
	if err != nil {
		return nil, fmt.Errorf("failed to get recordings: %w", err)
	}

	// Get players
	players, err := r.GetPlayers(sessionID)
	if err != nil {
		return nil, err
	}

	result := &models.SessionWithDetails{
		Session:    *session,
		Campaign:   jetModelToCampaign(&campaign),
		Recordings: make([]models.Recording, len(recordings)),
		Players:    players,
	}

	for i, rec := range recordings {
		result.Recordings[i] = *jetModelToRecording(&rec)
	}

	return result, nil
}

// Helper function to convert Jet model to our domain model
func jetModelToSession(m *model.Sessions) *models.Session {
	session := &models.Session{
		ID:            int64(*m.ID),
		CampaignID:    int64(m.CampaignID),
		Name:          m.Name,
		SessionNumber: int(m.SessionNumber),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}

	if m.SessionDate != nil {
		session.SessionDate = m.SessionDate
	}
	if m.Notes != nil {
		session.Notes = m.Notes
	}

	return session
}
