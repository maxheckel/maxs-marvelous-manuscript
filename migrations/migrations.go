package migrations

import (
	"embed"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed *.sql
var FS embed.FS

// GetMigrations returns the migration source for sql-migrate
func GetMigrations() migrate.MigrationSource {
	return &migrate.EmbedFileSystemMigrationSource{
		FileSystem: FS,
		Root:       ".",
	}
}
