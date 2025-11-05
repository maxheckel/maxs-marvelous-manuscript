package models

import "time"

type Recording struct {
	ID                   int64      `json:"id"`
	SessionID            *int64     `json:"session_id,omitempty"`
	FileID               string     `json:"file_id"`
	Filename             string     `json:"filename"`
	FilePath             string     `json:"file_path"`
	DurationSeconds      int        `json:"duration_seconds"`
	FileSizeBytes        int64      `json:"file_size_bytes"`
	Status               string     `json:"status"` // recording, completed, failed
	CreatedAt            time.Time  `json:"created_at"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	TranscriptionStatus  string     `json:"transcription_status"` // pending, processing, completed, failed
	Notes                *string    `json:"notes,omitempty"`
}

type CreateRecordingParams struct {
	SessionID *int64
	FileID    string
	Filename  string
	FilePath  string
}

type UpdateRecordingParams struct {
	SessionID           *int64
	DurationSeconds     *int
	FileSizeBytes       *int64
	Status              *string
	CompletedAt         *time.Time
	TranscriptionStatus *string
	Notes               *string
}
