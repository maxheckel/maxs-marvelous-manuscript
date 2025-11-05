package db

import (
	"fmt"
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/model"
	. "github.com/maxheckel/maxs-marvelous-manuscript/internal/db/gen/table"
	"github.com/maxheckel/maxs-marvelous-manuscript/pkg/models"
)

type RecordingRepository struct {
	db *DB
}

func NewRecordingRepository(db *DB) *RecordingRepository {
	return &RecordingRepository{db: db}
}

// Create creates a new recording
func (r *RecordingRepository) Create(params models.CreateRecordingParams) (*models.Recording, error) {
	jetModel := model.Recordings{
		FileID:   params.FileID,
		Filename: params.Filename,
		FilePath: params.FilePath,
		Status:   "recording",
	}

	if params.SessionID != nil {
		sessionID := int32(*params.SessionID)
		jetModel.SessionID = &sessionID
	}

	stmt := Recordings.
		INSERT(Recordings.MutableColumns).
		MODEL(jetModel).
		RETURNING(Recordings.AllColumns)

	var dest model.Recordings
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording: %w", err)
	}

	return jetModelToRecording(&dest), nil
}

// GetByID retrieves a recording by ID
func (r *RecordingRepository) GetByID(id int64) (*models.Recording, error) {
	stmt := SELECT(Recordings.AllColumns).
		FROM(Recordings).
		WHERE(Recordings.ID.EQ(Int32(int32(id))))

	var dest model.Recordings
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get recording: %w", err)
	}

	return jetModelToRecording(&dest), nil
}

// GetByFileID retrieves a recording by file ID
func (r *RecordingRepository) GetByFileID(fileID string) (*models.Recording, error) {
	stmt := SELECT(Recordings.AllColumns).
		FROM(Recordings).
		WHERE(Recordings.FileID.EQ(String(fileID)))

	var dest model.Recordings
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to get recording: %w", err)
	}

	return jetModelToRecording(&dest), nil
}

// List retrieves all recordings
func (r *RecordingRepository) List() ([]*models.Recording, error) {
	stmt := SELECT(Recordings.AllColumns).
		FROM(Recordings).
		ORDER_BY(Recordings.CreatedAt.DESC())

	var dest []model.Recordings
	err := stmt.Query(r.db.DB, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to list recordings: %w", err)
	}

	recordings := make([]*models.Recording, len(dest))
	for i, d := range dest {
		recordings[i] = jetModelToRecording(&d)
	}

	return recordings, nil
}

// Update updates a recording
func (r *RecordingRepository) Update(id int64, params models.UpdateRecordingParams) error {
	columns := ColumnList{}
	values := make(map[Column]interface{})

	if params.SessionID != nil {
		sessionID := int32(*params.SessionID)
		columns = append(columns, Recordings.SessionID)
		values[Recordings.SessionID] = &sessionID
	}
	if params.DurationSeconds != nil {
		duration := int32(*params.DurationSeconds)
		columns = append(columns, Recordings.DurationSeconds)
		values[Recordings.DurationSeconds] = &duration
	}
	if params.FileSizeBytes != nil {
		fileSize := int32(*params.FileSizeBytes)
		columns = append(columns, Recordings.FileSizeBytes)
		values[Recordings.FileSizeBytes] = &fileSize
	}
	if params.Status != nil {
		columns = append(columns, Recordings.Status)
		values[Recordings.Status] = *params.Status
	}
	if params.CompletedAt != nil {
		columns = append(columns, Recordings.CompletedAt)
		values[Recordings.CompletedAt] = params.CompletedAt
	}
	if params.TranscriptionStatus != nil {
		columns = append(columns, Recordings.TranscriptionStatus)
		values[Recordings.TranscriptionStatus] = params.TranscriptionStatus
	}
	if params.Notes != nil {
		columns = append(columns, Recordings.Notes)
		values[Recordings.Notes] = params.Notes
	}

	if len(columns) == 0 {
		return fmt.Errorf("no fields to update")
	}

	stmt := Recordings.
		UPDATE(columns).
		MODEL(values).
		WHERE(Recordings.ID.EQ(Int32(int32(id))))

	_, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to update recording: %w", err)
	}

	return nil
}

// Delete deletes a recording
func (r *RecordingRepository) Delete(id int64) error {
	stmt := Recordings.
		DELETE().
		WHERE(Recordings.ID.EQ(Int32(int32(id))))

	result, err := stmt.Exec(r.db.DB)
	if err != nil {
		return fmt.Errorf("failed to delete recording: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("recording not found")
	}

	return nil
}

// MarkCompleted marks a recording as completed
func (r *RecordingRepository) MarkCompleted(id int64, duration int, fileSize int64) error {
	now := time.Now()
	status := "completed"
	return r.Update(id, models.UpdateRecordingParams{
		DurationSeconds: &duration,
		FileSizeBytes:   &fileSize,
		Status:          &status,
		CompletedAt:     &now,
	})
}

// Helper function to convert Jet model to our domain model
func jetModelToRecording(m *model.Recordings) *models.Recording {
	rec := &models.Recording{
		ID:                  int64(*m.ID),
		FileID:              m.FileID,
		Filename:            m.Filename,
		FilePath:            m.FilePath,
		DurationSeconds:     0,
		FileSizeBytes:       0,
		Status:              m.Status,
		CreatedAt:           m.CreatedAt,
		TranscriptionStatus: "pending",
	}

	if m.SessionID != nil {
		sessionID := int64(*m.SessionID)
		rec.SessionID = &sessionID
	}
	if m.DurationSeconds != nil {
		rec.DurationSeconds = int(*m.DurationSeconds)
	}
	if m.FileSizeBytes != nil {
		rec.FileSizeBytes = int64(*m.FileSizeBytes)
	}
	if m.CompletedAt != nil {
		rec.CompletedAt = m.CompletedAt
	}
	if m.TranscriptionStatus != nil {
		rec.TranscriptionStatus = *m.TranscriptionStatus
	}
	if m.Notes != nil {
		rec.Notes = m.Notes
	}

	return rec
}
