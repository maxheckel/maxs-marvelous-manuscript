package db

import (
	"embed"
	"fmt"
	"io/fs"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed ../../migrations/*.sql
var migrationsFS embed.FS

// GetMigrations returns the migration source
func GetMigrations() *migrate.MemoryMigrationSource {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{},
	}

	// Read all migration files
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		fmt.Printf("Error reading migrations: %v\n", err)
		return migrations
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := fs.ReadFile(migrationsFS, "migrations/"+file.Name())
		if err != nil {
			fmt.Printf("Error reading migration %s: %v\n", file.Name(), err)
			continue
		}

		migrations.Migrations = append(migrations.Migrations, &migrate.Migration{
			Id:   file.Name(),
			Up:   []string{string(content)},
			Down: []string{},
		})
	}

	return migrations
}
