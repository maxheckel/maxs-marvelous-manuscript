package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

type DB struct {
	*sql.DB
}

type Config struct {
	DataDir string
	DBName  string
}

// New creates a new database connection and runs migrations
func New(cfg Config) (*DB, error) {
	if cfg.DBName == "" {
		cfg.DBName = "dnd_assistant.db"
	}

	dbPath := filepath.Join(cfg.DataDir, cfg.DBName)

	sqlDB, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{DB: sqlDB}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// Migrate runs all pending migrations
func (db *DB) Migrate() error {
	migrations := GetMigrations()

	n, err := migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if n > 0 {
		fmt.Printf("Applied %d migrations\n", n)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
