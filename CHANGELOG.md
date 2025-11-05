# Changelog

## [Unreleased]

### Added
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

### Technical Details
- All repositories use Jet's type-safe query builder
- Models auto-generated in `internal/db/gen/`
- Comprehensive relationships between campaigns, sessions, players, and recordings
- Foreign key constraints properly enforced
