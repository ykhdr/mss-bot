package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ykhdr/mss-bot/internal/minecraft"
	"github.com/ykhdr/mss-bot/internal/storage"
	"github.com/ykhdr/mss-bot/internal/storage/models"
)

// MockStorage is a mock implementation of storage.ServerStorage
type MockStorage struct {
	servers map[int64]*models.Server
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		servers: make(map[int64]*models.Server),
	}
}

func (m *MockStorage) GetByChatID(ctx context.Context, chatID int64) (*models.Server, error) {
	server, ok := m.servers[chatID]
	if !ok {
		return nil, storage.ErrNotFound{ChatID: chatID}
	}
	return server, nil
}

func (m *MockStorage) Upsert(ctx context.Context, server *models.Server) error {
	now := time.Now()
	if existing, ok := m.servers[server.ChatID]; ok {
		server.ID = existing.ID
		server.CreatedAt = existing.CreatedAt
	} else {
		server.ID = int64(len(m.servers) + 1)
		server.CreatedAt = now
	}
	server.UpdatedAt = now
	m.servers[server.ChatID] = server
	return nil
}

func (m *MockStorage) Delete(ctx context.Context, chatID int64) error {
	delete(m.servers, chatID)
	return nil
}

func (m *MockStorage) Close() error {
	return nil
}

func TestServerService_SetServerConfig(t *testing.T) {
	mockStorage := NewMockStorage()
	mcClient := minecraft.NewClient(5 * time.Second)
	service := NewServerService(mockStorage, mcClient)

	ctx := context.Background()
	err := service.SetServerConfig(ctx, 12345, "mc.example.com", 25565, "Test Server")

	require.NoError(t, err)

	server, err := mockStorage.GetByChatID(ctx, 12345)
	require.NoError(t, err)

	assert.Equal(t, "mc.example.com", server.IP)
	assert.Equal(t, 25565, server.Port)
	assert.Equal(t, "Test Server", server.Name)
}

func TestServerService_GetServerConfig(t *testing.T) {
	mockStorage := NewMockStorage()
	mcClient := minecraft.NewClient(5 * time.Second)
	service := NewServerService(mockStorage, mcClient)

	ctx := context.Background()

	// Set a server config first
	err := service.SetServerConfig(ctx, 12345, "mc.example.com", 25565, "Test Server")
	require.NoError(t, err)

	// Get it back
	server, err := service.GetServerConfig(ctx, 12345)
	require.NoError(t, err)

	assert.Equal(t, "mc.example.com", server.IP)
	assert.Equal(t, 25565, server.Port)
	assert.Equal(t, "Test Server", server.Name)
}

func TestServerService_GetServerConfig_NotFound(t *testing.T) {
	mockStorage := NewMockStorage()
	mcClient := minecraft.NewClient(5 * time.Second)
	service := NewServerService(mockStorage, mcClient)

	ctx := context.Background()

	_, err := service.GetServerConfig(ctx, 99999)
	assert.Error(t, err)

	var notFound storage.ErrNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestFormatConfig_NoServer(t *testing.T) {
	result := FormatConfig(nil)

	assert.Contains(t, result, "Сервер не настроен")
	assert.Contains(t, result, "/mss-set")
}

func TestFormatConfig_WithServer(t *testing.T) {
	server := &models.Server{
		IP:   "mc.example.com",
		Port: 25565,
		Name: "Test Server",
	}

	result := FormatConfig(server)

	assert.Contains(t, result, "mc.example.com")
	assert.Contains(t, result, "25565")
	assert.Contains(t, result, "Test Server")
}

func TestEscapeMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "hello\\_world"},
		{"test*bold*", "test\\*bold\\*"},
		{"normal text", "normal text"},
		{"[link](url)", "\\[link\\]\\(url\\)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := escapeMarkdown(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestServerStatusResult_FormatStatus_NoServer(t *testing.T) {
	result := &ServerStatusResult{
		Server: nil,
	}

	formatted := result.FormatStatus()
	assert.Contains(t, formatted, "Сервер не настроен")
}

func TestServerStatusResult_FormatStatus_Offline(t *testing.T) {
	result := &ServerStatusResult{
		Server: &models.Server{
			IP:   "mc.example.com",
			Port: 25565,
			Name: "Test Server",
		},
		Status: &minecraft.ServerStatus{
			Online: false,
		},
	}

	formatted := result.FormatStatus()
	assert.Contains(t, formatted, "🔴")
	assert.Contains(t, formatted, "Недоступен")
}

func TestServerStatusResult_FormatStatus_Online(t *testing.T) {
	result := &ServerStatusResult{
		Server: &models.Server{
			IP:   "mc.example.com",
			Port: 25565,
			Name: "Test Server",
		},
		Status: &minecraft.ServerStatus{
			Online:  true,
			Version: "1.20.4",
			Players: minecraft.PlayersInfo{
				Online: 5,
				Max:    20,
				Sample: []minecraft.Player{
					{Name: "Player1"},
					{Name: "Player2"},
				},
			},
		},
	}

	formatted := result.FormatStatus()
	assert.Contains(t, formatted, "🟢")
	assert.Contains(t, formatted, "1\\.20\\.4") // escaped for MarkdownV2
	assert.Contains(t, formatted, "5/20")
	assert.Contains(t, formatted, "Player1")
	assert.Contains(t, formatted, "Player2")
}
