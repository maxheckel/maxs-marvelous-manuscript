# Changelog

## [Unreleased]

### Added
- **Real-time Audio Recording**: Fully functional audio capture using malgo
  - Cross-platform support (macOS, Linux, Windows, Raspberry Pi)
  - 16kHz, mono, 16-bit WAV recording optimized for speech
  - Real-time writing to WAV files (no post-processing needed)
  - Pause/resume functionality
  - Approximately 115 MB per hour of recording
  - Thread-safe state management
- **Campaign Management**: Track D&D campaigns with name, description, and timestamps
- **Session Management**: Record individual game sessions within campaigns
  - Session number tracking
  - Session date and notes
  - Link recordings to sessions
- **Player Management**: Track players across campaigns and sessions
  - Player name, email, and character name
  - Many-to-many relationships with campaigns
  - Session attendance tracking
- **Jet Query Builder**: Refactored all repositories to use go-jet for type-safe SQL queries
  - Auto-generated models and query builders from schema
  - Type-safe query construction
  - Reduced boilerplate code
- **Database Schema**:
  - `campaigns` table
  - `sessions` table with foreign key to campaigns
  - `players` table
  - `campaign_players` junction table
  - `session_players` junction table with attendance tracking
  - `recordings` table now links to sessions
- **Make Target**: Added `make generate-jet` to regenerate Jet models from schema

### Changed
- **Repository Pattern**: All repositories now use Jet instead of raw SQL
- **Recordings**: Can now be associated with specific sessions
- **Migrations**: Moved migration embed to `migrations/` package for better organization

### Fixed
- **Jet Update Methods**: Fixed all repository Update methods to use SET() chaining instead of MODEL() with maps
  - RecordingRepository.Update() - proper handling of timestamps and nullable fields
  - CampaignRepository.Update() - consistent SET() pattern
  - SessionRepository.Update() - correct DATE vs TIMESTAMP handling
  - PlayerRepository.Update() - added rows affected check for consistency

### Technical Details
- All repositories use Jet's type-safe query builder
- Models auto-generated in `internal/db/gen/`
- Comprehensive relationships between campaigns, sessions, players, and recordings
- Foreign key constraints properly enforced
