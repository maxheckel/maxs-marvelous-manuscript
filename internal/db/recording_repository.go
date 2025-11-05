package db

import (
	"database/sql"
	"fmt"
	"time"

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
	query := `
		INSERT INTO recordings (file_id, filename, file_path, status, created_at)
		VALUES (?, ?, ?, 'recording', CURRENT_TIMESTAMP)
		RETURNING id, file_id, filename, file_path, duration_seconds, file_size_bytes,
		          status, created_at, completed_at, transcription_status, notes
	`

	var rec models.Recording
	var completedAt sql.NullTime
	var notes sql.NullString

	err := r.db.QueryRow(query, params.FileID, params.Filename, params.FilePath).Scan(
		&rec.ID,
		&rec.FileID,
		&rec.Filename,
		&rec.FilePath,
		&rec.DurationSeconds,
		&rec.FileSizeBytes,
		&rec.Status,
		&rec.CreatedAt,
		&completedAt,
		&rec.TranscriptionStatus,
		&notes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create recording: %w", err)
	}

	if completedAt.Valid {
		rec.CompletedAt = &completedAt.Time
	}
	if notes.Valid {
		rec.Notes = &notes.String
	}

	return &rec, nil
}

// GetByID retrieves a recording by ID
func (r *RecordingRepository) GetByID(id int64) (*models.Recording, error) {
	query := `
		SELECT id, file_id, filename, file_path, duration_seconds, file_size_bytes,
		       status, created_at, completed_at, transcription_status, notes
		FROM recordings
		WHERE id = ?
	`

	var rec models.Recording
	var completedAt sql.NullTime
	var notes sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&rec.ID,
		&rec.FileID,
		&rec.Filename,
		&rec.FilePath,
		&rec.DurationSeconds,
		&rec.FileSizeBytes,
		&rec.Status,
		&rec.CreatedAt,
		&completedAt,
		&rec.TranscriptionStatus,
		&notes,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recording not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get recording: %w", err)
	}

	if completedAt.Valid {
		rec.CompletedAt = &completedAt.Time
	}
	if notes.Valid {
		rec.Notes = &notes.String
	}

	return &rec, nil
}

// GetByFileID retrieves a recording by file ID
func (r *RecordingRepository) GetByFileID(fileID string) (*models.Recording, error) {
	query := `
		SELECT id, file_id, filename, file_path, duration_seconds, file_size_bytes,
		       status, created_at, completed_at, transcription_status, notes
		FROM recordings
		WHERE file_id = ?
	`

	var rec models.Recording
	var completedAt sql.NullTime
	var notes sql.NullString

	err := r.db.QueryRow(query, fileID).Scan(
		&rec.ID,
		&rec.FileID,
		&rec.Filename,
		&rec.FilePath,
		&rec.DurationSeconds,
		&rec.FileSizeBytes,
		&rec.Status,
		&rec.CreatedAt,
		&completedAt,
		&rec.TranscriptionStatus,
		&notes,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recording not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get recording: %w", err)
	}

	if completedAt.Valid {
		rec.CompletedAt = &completedAt.Time
	}
	if notes.Valid {
		rec.Notes = &notes.String
	}

	return &rec, nil
}

// List retrieves all recordings
func (r *RecordingRepository) List() ([]*models.Recording, error) {
	query := `
		SELECT id, file_id, filename, file_path, duration_seconds, file_size_bytes,
		       status, created_at, completed_at, transcription_status, notes
		FROM recordings
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list recordings: %w", err)
	}
	defer rows.Close()

	var recordings []*models.Recording
	for rows.Next() {
		var rec models.Recording
		var completedAt sql.NullTime
		var notes sql.NullString

		err := rows.Scan(
			&rec.ID,
			&rec.FileID,
			&rec.Filename,
			&rec.FilePath,
			&rec.DurationSeconds,
			&rec.FileSizeBytes,
			&rec.Status,
			&rec.CreatedAt,
			&completedAt,
			&rec.TranscriptionStatus,
			&notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recording: %w", err)
		}

		if completedAt.Valid {
			rec.CompletedAt = &completedAt.Time
		}
		if notes.Valid {
			rec.Notes = &notes.String
		}

		recordings = append(recordings, &rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recordings: %w", err)
	}

	return recordings, nil
}

// Update updates a recording
func (r *RecordingRepository) Update(id int64, params models.UpdateRecordingParams) error {
	// Build dynamic update query
	query := "UPDATE recordings SET "
	args := []interface{}{}
	updates := []string{}

	if params.DurationSeconds != nil {
		updates = append(updates, "duration_seconds = ?")
		args = append(args, *params.DurationSeconds)
	}
	if params.FileSizeBytes != nil {
		updates = append(updates, "file_size_bytes = ?")
		args = append(args, *params.FileSizeBytes)
	}
	if params.Status != nil {
		updates = append(updates, "status = ?")
		args = append(args, *params.Status)
	}
	if params.CompletedAt != nil {
		updates = append(updates, "completed_at = ?")
		args = append(args, *params.CompletedAt)
	}
	if params.TranscriptionStatus != nil {
		updates = append(updates, "transcription_status = ?")
		args = append(args, *params.TranscriptionStatus)
	}
	if params.Notes != nil {
		updates = append(updates, "notes = ?")
		args = append(args, *params.Notes)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update recording: %w", err)
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

// Delete deletes a recording
func (r *RecordingRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM recordings WHERE id = ?", id)
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
