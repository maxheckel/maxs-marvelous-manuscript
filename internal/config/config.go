package config

import (
	"os"
	"path/filepath"
	"strconv"
)

// Config holds application configuration
type Config struct {
	// Database configuration
	DataDir string
	DBName  string

	// Web server configuration
	WebPort string
	APIHost string

	// AI services configuration
	OpenAIAPIKey string

	// Recorder configuration
	AudioSampleRate int
	AudioChannels   int
	AudioBitDepth   int
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	homeDir, _ := os.UserHomeDir()
	defaultDataDir := filepath.Join(homeDir, ".dnd-assistant", "data")

	cfg := &Config{
		// Database defaults
		DataDir: getEnvOrDefault("DATA_DIR", defaultDataDir),
		DBName:  getEnvOrDefault("DB_NAME", "dnd_assistant.db"),

		// Web server defaults
		WebPort: getEnvOrDefault("PORT", "8080"),
		APIHost: getEnvOrDefault("API_HOST", "http://localhost:8080"),

		// AI services
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),

		// Recorder defaults (optimized for long recordings)
		AudioSampleRate: getEnvIntOrDefault("AUDIO_SAMPLE_RATE", 16000),
		AudioChannels:   getEnvIntOrDefault("AUDIO_CHANNELS", 1),
		AudioBitDepth:   getEnvIntOrDefault("AUDIO_BIT_DEPTH", 16),
	}

	// Ensure data directory exists
	os.MkdirAll(cfg.DataDir, 0755)

	return cfg
}

// getEnvOrDefault gets an environment variable or returns the default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets an integer environment variable or returns the default
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
