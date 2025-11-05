package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/recorder"
)

const (
	defaultDataDir = "./data"
)

type RecorderUI struct {
	app          fyne.App
	window       fyne.Window
	rec          *recorder.Recorder
	statusText   binding.String
	durationText binding.String
	filenameText binding.String
	filesizeText binding.String
	pauseButton  *widget.Button
	stopButton   *widget.Button
	dataDir      string
}

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

	// Create and run UI
	ui := NewRecorderUI(rec, dataDir)
	ui.Run()
}

func NewRecorderUI(rec *recorder.Recorder, dataDir string) *RecorderUI {
	a := app.New()
	w := a.NewWindow("D&D Session Recorder")

	// Set UI scale from environment variable if provided
	// Useful for small Raspberry Pi screens - try FYNE_SCALE=2.0 for larger UI
	// Default is 1.0, use 1.5-2.0 for small screens
	if scale := os.Getenv("FYNE_SCALE"); scale == "" {
		// Set a reasonable default for small screens
		os.Setenv("FYNE_SCALE", "1.3")
	}

	ui := &RecorderUI{
		app:          a,
		window:       w,
		rec:          rec,
		dataDir:      dataDir,
		statusText:   binding.NewString(),
		durationText: binding.NewString(),
		filenameText: binding.NewString(),
		filesizeText: binding.NewString(),
	}

	// Initialize bindings
	ui.statusText.Set("⏺️  Recording")
	ui.durationText.Set("00:00:00")
	ui.filenameText.Set("Initializing...")
	ui.filesizeText.Set("")

	ui.setupUI()
	return ui
}

func (ui *RecorderUI) setupUI() {
	// Create labels with data binding (thread-safe updates)
	statusLabel := widget.NewLabelWithData(ui.statusText)
	statusLabel.Alignment = fyne.TextAlignCenter
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	durationLabel := widget.NewLabelWithData(ui.durationText)
	durationLabel.Alignment = fyne.TextAlignCenter
	durationLabel.TextStyle = fyne.TextStyle{Bold: true}

	filenameLabel := widget.NewLabelWithData(ui.filenameText)
	filenameLabel.Alignment = fyne.TextAlignCenter
	filenameLabel.Wrapping = fyne.TextTruncate

	filesizeLabel := widget.NewLabelWithData(ui.filesizeText)
	filesizeLabel.Alignment = fyne.TextAlignCenter

	// Pause button - text updated in togglePause()
	ui.pauseButton = widget.NewButton("⏸️  Pause", func() {
		ui.togglePause()
	})
	ui.pauseButton.Importance = widget.HighImportance

	ui.stopButton = widget.NewButton("⏹️  Stop", func() {
		ui.stop()
	})
	ui.stopButton.Importance = widget.DangerImportance

	// Compact layout optimized for small screens
	content := container.NewBorder(
		// Top: Status and time (most important info)
		container.NewVBox(
			statusLabel,
			durationLabel,
		),
		// Bottom: Large buttons
		container.NewGridWithColumns(2,
			ui.pauseButton,
			ui.stopButton,
		),
		// Left/Right: nil
		nil, nil,
		// Center: File info (use full width available)
		container.NewVBox(
			filenameLabel,
			filesizeLabel,
		),
	)

	ui.window.SetContent(content)

	// Don't set a fixed size - let it adapt to the screen
	// Fullscreen will use the entire available display
	ui.window.SetFullScreen(true)
	ui.window.SetMaster()

	// Handle window close
	ui.window.SetCloseIntercept(func() {
		ui.stop()
	})
}

func (ui *RecorderUI) Run() {
	// Start recording
	if err := ui.rec.Start(); err != nil {
		log.Fatalf("Failed to start recording: %v", err)
	}

	// Start update loop using AfterFunc (runs on main thread)
	ui.scheduleUpdate()

	ui.window.ShowAndRun()
}

func (ui *RecorderUI) scheduleUpdate() {
	// This goroutine updates the bindings which are thread-safe
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			state := ui.rec.GetState()
			duration := ui.rec.GetDuration()
			currentFile := ui.rec.GetCurrentFile()
			fileSize := ui.rec.GetFileSize()

			// Update status
			statusText := ""
			switch state {
			case recorder.StateRecording:
				statusText = "⏺️  Recording"
			case recorder.StatePaused:
				statusText = "⏸️  Paused"
			case recorder.StateStopped:
				statusText = "⏹️  Stopped"
				// Don't schedule next update if stopped
				return
			}

			durationText := formatDuration(duration)

			var filename, sizeText string
			if currentFile != "" {
				filename = filepath.Base(currentFile)
				sizeText = formatBytes(fileSize)
			} else {
				filename = "Initializing..."
				sizeText = ""
			}

			// Update bindings (thread-safe)
			ui.statusText.Set(statusText)
			ui.durationText.Set(durationText)
			ui.filenameText.Set(filename)
			ui.filesizeText.Set(sizeText)
		}
	}()
}

func (ui *RecorderUI) togglePause() {
	state := ui.rec.GetState()
	if state == recorder.StatePaused {
		if err := ui.rec.Resume(); err != nil {
			log.Printf("Error resuming: %v", err)
		} else {
			ui.pauseButton.SetText("⏸️  Pause")
		}
	} else if state == recorder.StateRecording {
		if err := ui.rec.Pause(); err != nil {
			log.Printf("Error pausing: %v", err)
		} else {
			ui.pauseButton.SetText("▶️  Resume")
		}
	}
}

func (ui *RecorderUI) stop() {
	state := ui.rec.GetState()
	if state == recorder.StateStopped {
		ui.app.Quit()
		return
	}

	ui.statusText.Set("⏹️  Stopping...")
	ui.pauseButton.Disable()
	ui.stopButton.Disable()

	// Do the stop operation in a background goroutine
	// to keep UI responsive, but don't call UI methods from it
	go func() {
		if err := ui.rec.Stop(); err != nil {
			log.Printf("Error stopping recorder: %v", err)
		}

		// Update status text via binding (thread-safe)
		ui.statusText.Set("✅ Recording Saved!")

		// Wait a bit for user to see the message
		time.Sleep(2 * time.Second)

		// Quit the app - this is safe to call from goroutine
		ui.app.Quit()
	}()
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
