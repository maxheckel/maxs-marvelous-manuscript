# Max's Marvelous Manuscript

A D&D session recording and analysis assistant that runs on your local device (laptop or Raspberry Pi). Record your D&D sessions, transcribe them with AI, and get insights about your campaigns.

## Features

- **Audio Recording**: Record long D&D sessions with optimized audio quality (16kHz mono) to save space
- **Graphical Recorder UI**: Fullscreen GUI with pause/stop buttons, real-time duration display, and file size tracking
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
- Node.js 18+ and npm (only needed for web interface development)
- SQLite3 (usually comes pre-installed)
- For the GUI recorder: OpenGL libraries (usually pre-installed on modern systems)
  - **Raspberry Pi OS**: `sudo apt-get install libgl1-mesa-dev xorg-dev libasound2-dev`
  - **Linux**: `libgl1-mesa-dev` and `xorg-dev`
  - **macOS**: Comes pre-installed with Xcode Command Line Tools
  - **Windows**: Comes pre-installed

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

### Raspberry Pi Setup

For running the recorder on a Raspberry Pi with a small touchscreen:

1. Install system dependencies:
```bash
sudo apt-get update
sudo apt-get install -y libgl1-mesa-dev xorg-dev libasound2-dev
```

2. Install Go 1.25+ (if not already installed):
```bash
wget https://go.dev/dl/go1.25.1.linux-arm64.tar.gz  # or armv7l for 32-bit
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.1.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

3. Build the recorder (only the recorder is needed for Pi):
```bash
git clone <repository-url>
cd maxs-marvelous-manuscript
go mod download
go build -o bin/recorder ./cmd/recorder
```

4. Run the recorder (adjust scale for your screen size):
```bash
# For 3.5" screens, use larger scale
FYNE_SCALE=2.0 ./bin/recorder

# For 5-7" screens, use default
./bin/recorder
```

5. Optional - Auto-start on boot:
Create `/etc/systemd/system/dnd-recorder.service`:
```ini
[Unit]
Description=D&D Session Recorder
After=graphical.target

[Service]
Type=simple
User=pi
Environment="DISPLAY=:0"
Environment="FYNE_SCALE=1.5"
WorkingDirectory=/home/pi/maxs-marvelous-manuscript
ExecStart=/home/pi/maxs-marvelous-manuscript/bin/recorder
Restart=on-failure

[Install]
WantedBy=graphical.target
```

Then enable it:
```bash
sudo systemctl enable dnd-recorder
sudo systemctl start dnd-recorder
```

## Usage

### Recording a Session

Start the recorder GUI application:

```bash
./bin/recorder
```

The application will open in fullscreen mode with:
- **Status Display**: Shows current state (Recording/Paused/Stopped)
- **Duration Counter**: Real-time display of recording duration
- **File Info**: Current filename and file size
- **Pause Button**: Toggle between recording and paused states (highlighted in blue)
- **Stop Button**: Stop recording and save to database (highlighted in red)

The recorder starts automatically when launched. The audio file will be saved to the `data/` directory with metadata stored in the database.

#### Raspberry Pi & Small Screens

The UI is optimized for small touchscreens (3.5" - 7") commonly used with Raspberry Pi:
- Automatically scales UI elements (default 1.3x for better readability)
- Large, touch-friendly buttons
- Compact layout that fits small displays
- Reduced CPU usage (updates 4 times per second instead of 10)

To adjust UI scale for your screen size:
```bash
# Make UI larger for very small screens (e.g., 3.5" display)
FYNE_SCALE=2.0 ./bin/recorder

# Default scaling (recommended for 5-7" screens)
./bin/recorder

# Smaller UI for larger screens
FYNE_SCALE=1.0 ./bin/recorder
```

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
   - Fullscreen GUI application using Fyne framework
   - Real-time display of recording status, duration, and file size
   - Pause/Resume and Stop buttons for control
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
- [x] Fullscreen GUI with Fyne (pause/stop buttons)
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
- **GUI recorder** uses Fyne for a native, fullscreen interface with pause/stop buttons.
- Ensure your system has a working microphone/audio input device.
- AI features require API keys and may incur costs (OpenAI API).
- This runs entirely locally except for AI API calls.
- The recorder writes WAV files directly - no encoding overhead during recording.
