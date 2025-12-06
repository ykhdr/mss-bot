package storage

import (
	"context"

	"github.com/ykhdr/mss-bot/internal/storage/models"
)

// ServerStorage defines the interface for server configuration storage
type ServerStorage interface {
	// GetByChatID returns server configuration for a specific chat
	GetByChatID(ctx context.Context, chatID int64) (*models.Server, error)

	// Upsert creates or updates server configuration for a chat
	Upsert(ctx context.Context, server *models.Server) error

	// Delete removes server configuration for a chat
	Delete(ctx context.Context, chatID int64) error

	// Close closes the storage connection
	Close() error
}

// ErrNotFound is returned when a server configuration is not found
type ErrNotFound struct {
	ChatID int64
}

func (e ErrNotFound) Error() string {
	return "server configuration not found"
}
