package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ykhdr/mss-bot/internal/storage"
	"github.com/ykhdr/mss-bot/internal/storage/models"
)

func setupTestDB(t *testing.T) *Storage {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	s, err := New(dbPath)
	require.NoError(t, err)

	t.Cleanup(func() {
		s.Close()
	})

	return s
}

func TestStorage_Upsert_Insert(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	server := &models.Server{
		ChatID: 12345,
		IP:     "mc.example.com",
		Port:   25565,
		Name:   "Test Server",
	}

	err := s.Upsert(ctx, server)
	require.NoError(t, err)

	assert.NotZero(t, server.ID)
	assert.NotZero(t, server.CreatedAt)
	assert.NotZero(t, server.UpdatedAt)
}

func TestStorage_Upsert_Update(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	server := &models.Server{
		ChatID: 12345,
		IP:     "mc.example.com",
		Port:   25565,
		Name:   "Test Server",
	}

	err := s.Upsert(ctx, server)
	require.NoError(t, err)

	originalID := server.ID
	originalCreatedAt := server.CreatedAt

	// Update the server
	server.IP = "new.example.com"
	server.Port = 25566
	server.Name = "Updated Server"

	err = s.Upsert(ctx, server)
	require.NoError(t, err)

	assert.Equal(t, originalID, server.ID)
	assert.WithinDuration(t, originalCreatedAt, server.CreatedAt, 0)
	assert.True(t, server.UpdatedAt.After(originalCreatedAt) || server.UpdatedAt.Equal(originalCreatedAt))
}

func TestStorage_GetByChatID_Found(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	server := &models.Server{
		ChatID: 12345,
		IP:     "mc.example.com",
		Port:   25565,
		Name:   "Test Server",
	}

	err := s.Upsert(ctx, server)
	require.NoError(t, err)

	found, err := s.GetByChatID(ctx, 12345)
	require.NoError(t, err)

	assert.Equal(t, server.ID, found.ID)
	assert.Equal(t, server.ChatID, found.ChatID)
	assert.Equal(t, server.IP, found.IP)
	assert.Equal(t, server.Port, found.Port)
	assert.Equal(t, server.Name, found.Name)
}

func TestStorage_GetByChatID_NotFound(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	_, err := s.GetByChatID(ctx, 99999)
	assert.Error(t, err)

	var notFound storage.ErrNotFound
	assert.ErrorAs(t, err, &notFound)
	assert.Equal(t, int64(99999), notFound.ChatID)
}

func TestStorage_Delete(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	server := &models.Server{
		ChatID: 12345,
		IP:     "mc.example.com",
		Port:   25565,
		Name:   "Test Server",
	}

	err := s.Upsert(ctx, server)
	require.NoError(t, err)

	err = s.Delete(ctx, 12345)
	require.NoError(t, err)

	_, err = s.GetByChatID(ctx, 12345)
	assert.Error(t, err)

	var notFound storage.ErrNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestStorage_Delete_NonExistent(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	// Should not error when deleting non-existent record
	err := s.Delete(ctx, 99999)
	assert.NoError(t, err)
}

func TestStorage_MultipleChatIDs(t *testing.T) {
	s := setupTestDB(t)
	ctx := context.Background()

	servers := []*models.Server{
		{ChatID: 111, IP: "server1.com", Port: 25565, Name: "Server 1"},
		{ChatID: 222, IP: "server2.com", Port: 25566, Name: "Server 2"},
		{ChatID: 333, IP: "server3.com", Port: 25567, Name: "Server 3"},
	}

	for _, server := range servers {
		err := s.Upsert(ctx, server)
		require.NoError(t, err)
	}

	for _, server := range servers {
		found, err := s.GetByChatID(ctx, server.ChatID)
		require.NoError(t, err)
		assert.Equal(t, server.IP, found.IP)
		assert.Equal(t, server.Name, found.Name)
	}
}
