package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ykhdr/mss-bot/internal/storage"
	"github.com/ykhdr/mss-bot/internal/storage/models"
)

// Storage implements ServerStorage interface using SQLite
type Storage struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

// New creates a new SQLite storage instance
func New(dbPath string) (*Storage, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Storage{
		db: db,
		sb: Builder(),
	}, nil
}

// GetByChatID returns server configuration for a specific chat
func (s *Storage) GetByChatID(ctx context.Context, chatID int64) (*models.Server, error) {
	query, args, err := s.sb.
		Select("id", "chat_id", "ip", "port", "name", "created_at", "updated_at").
		From("servers").
		Where(squirrel.Eq{"chat_id": chatID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var server models.Server
	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&server.ID,
		&server.ChatID,
		&server.IP,
		&server.Port,
		&server.Name,
		&server.CreatedAt,
		&server.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound{ChatID: chatID}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	return &server, nil
}

// Upsert creates or updates server configuration for a chat
func (s *Storage) Upsert(ctx context.Context, server *models.Server) error {
	now := time.Now()

	// Try to get existing record
	existing, err := s.GetByChatID(ctx, server.ChatID)
	if err != nil && !isNotFound(err) {
		return fmt.Errorf("failed to check existing server: %w", err)
	}

	if existing != nil {
		// Update existing record
		query, args, err := s.sb.
			Update("servers").
			Set("ip", server.IP).
			Set("port", server.Port).
			Set("name", server.Name).
			Set("updated_at", now).
			Where(squirrel.Eq{"chat_id": server.ChatID}).
			ToSql()
		if err != nil {
			return fmt.Errorf("failed to build update query: %w", err)
		}

		_, err = s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update server: %w", err)
		}

		server.ID = existing.ID
		server.CreatedAt = existing.CreatedAt
		server.UpdatedAt = now
	} else {
		// Insert new record
		query, args, err := s.sb.
			Insert("servers").
			Columns("chat_id", "ip", "port", "name", "created_at", "updated_at").
			Values(server.ChatID, server.IP, server.Port, server.Name, now, now).
			ToSql()
		if err != nil {
			return fmt.Errorf("failed to build insert query: %w", err)
		}

		result, err := s.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert server: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}

		server.ID = id
		server.CreatedAt = now
		server.UpdatedAt = now
	}

	return nil
}

// Delete removes server configuration for a chat
func (s *Storage) Delete(ctx context.Context, chatID int64) error {
	query, args, err := s.sb.
		Delete("servers").
		Where(squirrel.Eq{"chat_id": chatID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

func isNotFound(err error) bool {
	_, ok := err.(storage.ErrNotFound)
	return ok
}
