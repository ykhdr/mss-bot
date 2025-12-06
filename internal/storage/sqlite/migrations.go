package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Up      func(ctx context.Context, db *sql.DB) error
	Down    func(ctx context.Context, db *sql.DB) error
}

// migrations contains all database migrations in order
var migrations = []Migration{
	{
		Version: 1,
		Up:      upCreateServersTable,
		Down:    downCreateServersTable,
	},
}

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
	ctx := context.Background()

	// Create migrations table if not exists
	if err := createMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current version
	currentVersion, err := getCurrentVersion(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Run pending migrations
	for _, m := range migrations {
		if m.Version > currentVersion {
			if err := m.Up(ctx, db); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", m.Version, err)
			}

			if err := setVersion(ctx, db, m.Version); err != nil {
				return fmt.Errorf("failed to set version %d: %w", m.Version, err)
			}
		}
	}

	return nil
}

func createMigrationsTable(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

func getCurrentVersion(ctx context.Context, db *sql.DB) (int, error) {
	var version int
	err := db.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func setVersion(ctx context.Context, db *sql.DB, version int) error {
	_, err := db.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", version)
	return err
}

func upCreateServersTable(ctx context.Context, db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL UNIQUE,
			ip TEXT NOT NULL,
			port INTEGER NOT NULL DEFAULT 25565,
			name TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	// Create index for chat_id lookup
	indexQuery := `CREATE INDEX IF NOT EXISTS idx_servers_chat_id ON servers(chat_id)`
	_, err = db.ExecContext(ctx, indexQuery)
	return err
}

func downCreateServersTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS servers")
	return err
}

// Builder returns a squirrel statement builder configured for SQLite
func Builder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
}
