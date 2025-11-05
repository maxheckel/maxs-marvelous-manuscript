# Max's Marvelous Manuscript

A D&D session recording and analysis assistant that runs on your local device (laptop or Raspberry Pi). Record your D&D sessions, transcribe them with AI, and get insights about your campaigns.

## Features

- **Audio Recording**: Record long D&D sessions with optimized audio quality (16kHz mono) to save space
- **Session Management**: Store and manage multiple recording sessions in a local SQLite database
- **Web Interface**: Vue 3 TypeScript frontend for viewing and managing recordings
- **AI Integration**: Interfaces for transcription (Whisper), speaker diarization, session summarization, and embeddings
- **Local Storage**: Everything runs and stores data locally on your device

## Project Structure

```
.
├── cmd/
│   ├── recorder/          # CLI recording application
│   └── web/              # Web server application
├── internal/
│   ├── ai/               # AI services (transcription, summarization, embeddings)
│   ├── api/              # HTTP API handlers
│   ├── config/           # Configuration management
│   ├── db/               # Database connection and migrations
│   └── recorder/         # Audio recording logic
├── pkg/
│   └── models/           # Shared data models
├── migrations/           # SQL database migrations
├── web/
│   └── frontend/         # Vue 3 + TypeScript frontend
└── data/                # Local data storage (gitignored)
```

## Prerequisites

- Go 1.25.1 or later
- Node.js 18+ and npm
- SQLite3 (usually comes pre-installed)

## Installation

### Backend (Go)

1. Clone the repository:
```bash
git clone <repository-url>
cd maxs-marvelous-manuscript
```

2. Install Go dependencies:
```bash
go mod download
```

3. Build the applications:
```bash
# Build recorder
go build -o bin/recorder ./cmd/recorder

# Build web server
go build -o bin/web ./cmd/web
```

### Frontend (Vue)

1. Navigate to the frontend directory:
```bash
cd web/frontend
```

2. Install dependencies:
```bash
npm install
```

3. Build the frontend:
```bash
npm run build
```

For development with hot reload:
```bash
npm run dev
```

## Usage

### Recording a Session

Start the recording CLI application:

```bash
./bin/recorder
```

Controls:
- The recorder starts automatically
- Press `Ctrl+C` to stop recording

The audio file will be saved to the `data/` directory with metadata stored in the database.

### Web Interface

1. Start the web server:

```bash
./bin/web
```

2. Open your browser to `http://localhost:8080`

The web interface allows you to:
- View all recorded sessions
- Play back recordings
- See session metadata (duration, file size, etc.)
- Delete recordings

### Configuration

Configuration is done through environment variables:

```bash
# Data storage
export DATA_DIR="./data"              # Where audio files and DB are stored
export DB_NAME="dnd_assistant.db"     # Database filename

# Web server
export PORT="8080"                    # Web server port
export API_HOST="http://localhost:8080"

# AI services (for future use)
export OPENAI_API_KEY="your-key-here"

# Audio recording settings
export AUDIO_SAMPLE_RATE="16000"      # 16kHz for speech
export AUDIO_CHANNELS="1"             # Mono
export AUDIO_BIT_DEPTH="16"           # 16-bit
```

## Architecture

### Two Entry Points

1. **Recorder Application** (`cmd/recorder/main.go`)
   - CLI application for recording audio
   - Simple terminal UI showing recording status
   - Saves files and creates database records

2. **Web Application** (`cmd/web/main.go`)
   - HTTP API server
   - Serves Vue frontend
   - Provides REST API for managing recordings

### Database

- **SQLite** database stored locally
- **Migrations** managed with sql-migrate
- **Query Builder** using go-jet/jet for type-safe SQL queries
- **Repository pattern** for data access
- Tables:
  - `campaigns` - D&D campaigns
  - `sessions` - Individual game sessions within campaigns
  - `recordings` - Audio recordings for sessions
  - `players` - Player information
  - `campaign_players` - Many-to-many relationship between campaigns and players
  - `session_players` - Session attendance tracking

### Audio Recording

- **Real-time audio capture** using malgo (mini audio Go bindings)
- **Optimized for long sessions** (4+ hours)
- **Low quality audio** (16kHz, mono, 16-bit) to reduce file size
  - Approximately 115 MB per hour of recording
- **WAV format** for simplicity and compatibility
- **Pause/resume functionality** - pause recording without stopping
- **Cross-platform** - works on macOS, Linux, Windows, and Raspberry Pi

### AI Services (Interfaces)

Located in `internal/ai/`, these are interface definitions for future implementation:

- **Transcriber**: Convert audio to text with speaker diarization
  - Planned: OpenAI Whisper integration
  - Planned: Speaker diarization (may require additional services)

- **Summarizer**: Generate session summaries
  - Extract key events, NPCs, locations
  - Identify combat encounters
  - Note important decisions
  - Highlight cliffhangers

- **EmbeddingGenerator**: Create embeddings for semantic search
  - Search across sessions
  - Find similar moments
  - Query session content

## Development

### Running in Development Mode

Terminal 1 - Backend:
```bash
go run ./cmd/web
```

Terminal 2 - Frontend:
```bash
cd web/frontend
npm run dev
```

The frontend dev server (port 5173) will proxy API requests to the backend (port 8080).

### Database Migrations

Migrations are automatically run when the application starts. To create a new migration:

1. Create a new file in `migrations/` with the naming pattern `###_description.sql`
2. Use the sql-migrate format:

```sql
-- +migrate Up
CREATE TABLE ...;

-- +migrate Down
DROP TABLE ...;
```

3. After creating a migration, regenerate Jet models:

```bash
make generate-jet
```

This will:
- Create a temporary database with the new schema
- Generate type-safe Go models and query builders in `internal/db/gen/`

## Roadmap

### Completed ✅
- [x] Real-time audio capture with malgo
- [x] Pause/resume functionality
- [x] Campaign and session management
- [x] Player tracking and attendance
- [x] Type-safe database queries with Jet

### In Progress 🚧
- [ ] Implement OpenAI Whisper transcription
- [ ] Add speaker diarization
- [ ] Implement AI summarization

### Planned 📋
- [ ] Add embedding generation and semantic search
- [ ] Session comparison and analysis
- [ ] Export session notes and summaries
- [ ] Character and NPC tracking
- [ ] Timeline visualization
- [ ] Web UI for campaigns and sessions
- [ ] API endpoints for all resources

## Contributing

This is a personal project, but suggestions and contributions are welcome!

## License

[Add your license here]

## Notes

- Audio files are optimized for long sessions: approximately **115 MB per hour** at 16kHz, mono, 16-bit.
- **Audio recording is fully implemented** using malgo for cross-platform capture.
- Ensure your system has a working microphone/audio input device.
- AI features require API keys and may incur costs (OpenAI API).
- This runs entirely locally except for AI API calls.
- The recorder writes WAV files directly - no encoding overhead during recording.
