package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/recorder"
)

const (
	defaultDataDir = "./data"
)

func main() {
	// Ensure data directory exists
	dataDir := defaultDataDir
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	database, err := db.New(db.Config{
		DataDir: dataDir,
		DBName:  "dnd_assistant.db",
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	recordingRepo := db.NewRecordingRepository(database)

	// Create recorder
	rec := recorder.New(recorder.Config{
		DataDir: dataDir,
		Format:  recorder.DefaultAudioFormat(),
		DB:      recordingRepo,
	})

	// Start recording
	fmt.Println("🎙️  D&D Session Recorder")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("Starting recording...")

	if err := rec.Start(); err != nil {
		log.Fatalf("Failed to start recording: %v", err)
	}

	fmt.Println("✅ Recording started!")
	fmt.Println()
	fmt.Println("Controls:")
	fmt.Println("  Press Ctrl+C to stop recording")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Display recording status
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := false
	for !done {
		select {
		case <-sigChan:
			fmt.Println("\n\n⏹️  Stopping recording...")
			if err := rec.Stop(); err != nil {
				log.Printf("Error stopping recorder: %v", err)
			}
			done = true

		case <-ticker.C:
			state := rec.GetState()
			duration := rec.GetDuration()
			fmt.Printf("\r⏺️  Recording: %s | Duration: %s",
				state.String(),
				formatDuration(duration))
		}
	}

	fmt.Println("\n✅ Recording saved!")
	fmt.Printf("📁 Location: %s\n", filepath.Join(dataDir, "recording_*.wav"))
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
