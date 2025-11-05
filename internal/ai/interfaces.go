package ai

import (
	"context"
	"io"
)

// TranscriptionSegment represents a segment of transcribed audio with speaker information
type TranscriptionSegment struct {
	Speaker    string  // Speaker identifier (e.g., "SPEAKER_00", "SPEAKER_01")
	Text       string  // Transcribed text
	Start      float64 // Start time in seconds
	End        float64 // End time in seconds
	Confidence float64 // Confidence score (0.0 to 1.0)
}

// TranscriptionResult contains the full transcription with speaker diarization
type TranscriptionResult struct {
	Segments   []TranscriptionSegment
	Language   string  // Detected language
	Duration   float64 // Total audio duration in seconds
	FullText   string  // Complete transcription without speaker labels
}

// SummarySection represents a section of the session summary
type SummarySection struct {
	Title   string
	Content string
}

// SessionSummary contains structured summary of a D&D session
type SessionSummary struct {
	Overview      string           // High-level overview
	KeyEvents     []string         // Important events that occurred
	NPCs          []string         // Non-player characters encountered
	Locations     []string         // Locations visited
	Items         []string         // Items obtained or discussed
	Combat        []string         // Combat encounters
	Decisions     []string         // Important decisions made by the party
	Cliffhangers  []string         // Unresolved plot points
	CustomSections []SummarySection // Additional sections
}

// Embedding represents a vector embedding
type Embedding struct {
	Vector []float64
	Model  string
}

// Transcriber handles audio transcription with speaker diarization
type Transcriber interface {
	// TranscribeFile transcribes an audio file and performs speaker diarization
	TranscribeFile(ctx context.Context, filePath string) (*TranscriptionResult, error)

	// TranscribeStream transcribes an audio stream with real-time speaker diarization
	TranscribeStream(ctx context.Context, audioStream io.Reader) (*TranscriptionResult, error)
}

// Summarizer generates summaries of D&D sessions
type Summarizer interface {
	// SummarizeSession generates a structured summary from a transcription
	SummarizeSession(ctx context.Context, transcription *TranscriptionResult) (*SessionSummary, error)

	// SummarizeText generates a summary from raw text
	SummarizeText(ctx context.Context, text string) (*SessionSummary, error)
}

// EmbeddingGenerator creates vector embeddings for semantic search
type EmbeddingGenerator interface {
	// GenerateEmbedding creates an embedding for the given text
	GenerateEmbedding(ctx context.Context, text string) (*Embedding, error)

	// GenerateBatchEmbeddings creates embeddings for multiple texts
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]*Embedding, error)
}

// AIService combines all AI capabilities
type AIService interface {
	Transcriber
	Summarizer
	EmbeddingGenerator
}
