package recorder

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/google/uuid"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db"
	"github.com/maxheckel/maxs-marvelous-manuscript/pkg/models"
)

// AudioFormat defines the audio recording parameters
type AudioFormat struct {
	SampleRate  int // 16000 Hz for lower quality, longer recordings
	Channels    int // 1 for mono
	BitDepth    int // 16 bits
}

// DefaultAudioFormat returns the default format optimized for long D&D sessions
func DefaultAudioFormat() AudioFormat {
	return AudioFormat{
		SampleRate: 16000, // 16 kHz is good enough for speech
		Channels:   1,     // Mono
		BitDepth:   16,    // 16-bit
	}
}

// RecorderState represents the current state of the recorder
type RecorderState int

const (
	StateIdle RecorderState = iota
	StateRecording
	StatePaused
	StateStopped
)

func (s RecorderState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRecording:
		return "recording"
	case StatePaused:
		return "paused"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// Recorder handles audio recording
type Recorder struct {
	format       AudioFormat
	dataDir      string
	db           *db.RecordingRepository
	state        RecorderState
	currentFile  *os.File
	currentID    int64
	fileID       string
	startTime    time.Time
	pauseTime    time.Time
	pausedTotal  time.Duration
	mu           sync.RWMutex
	stopChan     chan struct{}
	audioBuffer  []byte
	malgoCtx     *malgo.AllocatedContext
	malgoDevice  *malgo.Device
	captureWg    sync.WaitGroup
}

// Config holds recorder configuration
type Config struct {
	DataDir string
	Format  AudioFormat
	DB      *db.RecordingRepository
}

// New creates a new recorder
func New(cfg Config) *Recorder {
	if cfg.Format.SampleRate == 0 {
		cfg.Format = DefaultAudioFormat()
	}

	return &Recorder{
		format:  cfg.Format,
		dataDir: cfg.DataDir,
		db:      cfg.DB,
		state:   StateIdle,
	}
}

// Start begins recording
func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StateIdle && r.state != StateStopped {
		return fmt.Errorf("recorder is already active")
	}

	// Generate unique file ID
	r.fileID = uuid.New().String()
	filename := fmt.Sprintf("recording_%s.wav", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(r.dataDir, filename)

	// Create the audio file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create audio file: %w", err)
	}

	// Write WAV header (will be updated when recording stops)
	if err := r.writeWAVHeader(file, 0); err != nil {
		file.Close()
		return fmt.Errorf("failed to write WAV header: %w", err)
	}

	r.currentFile = file
	r.startTime = time.Now()
	r.pausedTotal = 0
	r.audioBuffer = make([]byte, 0)
	r.stopChan = make(chan struct{})

	// Create database record
	rec, err := r.db.Create(models.CreateRecordingParams{
		FileID:   r.fileID,
		Filename: filename,
		FilePath: filePath,
	})
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to create recording record: %w", err)
	}

	r.currentID = rec.ID
	r.state = StateRecording

	// Start audio capture in a goroutine
	r.captureWg.Add(1)
	go func() {
		defer r.captureWg.Done()
		r.captureAudio()
	}()

	return nil
}

// Pause pauses the recording
func (r *Recorder) Pause() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StateRecording {
		return fmt.Errorf("recorder is not recording")
	}

	r.pauseTime = time.Now()
	r.state = StatePaused
	return nil
}

// Resume resumes the recording
func (r *Recorder) Resume() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StatePaused {
		return fmt.Errorf("recorder is not paused")
	}

	r.pausedTotal += time.Since(r.pauseTime)
	r.state = StateRecording
	return nil
}

// Stop stops the recording
func (r *Recorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != StateRecording && r.state != StatePaused {
		return fmt.Errorf("recorder is not active")
	}

	// Signal the audio capture to stop
	close(r.stopChan)

	// Release lock while waiting for capture to finish
	r.mu.Unlock()
	r.captureWg.Wait()
	r.mu.Lock()

	// Calculate final duration and file size
	duration := time.Since(r.startTime) - r.pausedTotal
	durationSeconds := int(duration.Seconds())

	// Update WAV header with final size
	if err := r.finalizeWAVFile(); err != nil {
		return fmt.Errorf("failed to finalize WAV file: %w", err)
	}

	// Get file size
	fileInfo, err := r.currentFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// Close the file
	if err := r.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	// Update database record
	if err := r.db.MarkCompleted(r.currentID, durationSeconds, fileSize); err != nil {
		return fmt.Errorf("failed to update recording: %w", err)
	}

	r.state = StateStopped
	r.currentFile = nil

	return nil
}

// GetState returns the current recorder state
func (r *Recorder) GetState() RecorderState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

// GetDuration returns the current recording duration
func (r *Recorder) GetDuration() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.state == StateIdle || r.state == StateStopped {
		return 0
	}

	elapsed := time.Since(r.startTime) - r.pausedTotal
	if r.state == StatePaused {
		elapsed -= time.Since(r.pauseTime)
	}

	return elapsed
}

// captureAudio captures audio from the microphone using malgo
func (r *Recorder) captureAudio() {
	// Initialize malgo context
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		fmt.Printf("Failed to initialize malgo context: %v\n", err)
		return
	}
	r.malgoCtx = ctx
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	// Configure capture device
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = uint32(r.format.Channels)
	deviceConfig.SampleRate = uint32(r.format.SampleRate)
	deviceConfig.Alsa.NoMMap = 1

	// Data callback - called when audio data is available
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		r.mu.RLock()
		state := r.state
		file := r.currentFile
		r.mu.RUnlock()

		// Only write if we're actively recording (not paused)
		if state == StateRecording && file != nil {
			_, err := file.Write(pSample)
			if err != nil {
				fmt.Printf("Failed to write audio data: %v\n", err)
			}
		}
	}

	// Initialize the device
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		fmt.Printf("Failed to initialize capture device: %v\n", err)
		return
	}
	r.malgoDevice = device
	defer device.Uninit()

	// Start the device
	err = device.Start()
	if err != nil {
		fmt.Printf("Failed to start capture device: %v\n", err)
		return
	}

	// Wait for stop signal
	<-r.stopChan

	// Stop the device
	err = device.Stop()
	if err != nil {
		fmt.Printf("Failed to stop capture device: %v\n", err)
	}
}

// writeWAVHeader writes a WAV file header
func (r *Recorder) writeWAVHeader(file *os.File, dataSize uint32) error {
	// WAV header structure
	sampleRate := uint32(r.format.SampleRate)
	numChannels := uint16(r.format.Channels)
	bitsPerSample := uint16(r.format.BitDepth)
	byteRate := sampleRate * uint32(numChannels) * uint32(bitsPerSample) / 8
	blockAlign := numChannels * bitsPerSample / 8

	header := make([]byte, 44)

	// RIFF chunk
	copy(header[0:4], "RIFF")
	writeUint32(header[4:8], dataSize+36)
	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	writeUint32(header[16:20], 16)                // fmt chunk size
	writeUint16(header[20:22], 1)                 // PCM format
	writeUint16(header[22:24], numChannels)
	writeUint32(header[24:28], sampleRate)
	writeUint32(header[28:32], byteRate)
	writeUint16(header[32:34], blockAlign)
	writeUint16(header[34:36], bitsPerSample)

	// data chunk
	copy(header[36:40], "data")
	writeUint32(header[40:44], dataSize)

	_, err := file.WriteAt(header, 0)
	return err
}

// finalizeWAVFile updates the WAV header with the final file size
func (r *Recorder) finalizeWAVFile() error {
	fileInfo, err := r.currentFile.Stat()
	if err != nil {
		return err
	}

	dataSize := uint32(fileInfo.Size() - 44)
	return r.writeWAVHeader(r.currentFile, dataSize)
}

// Helper functions for writing WAV header
func writeUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func writeUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}
