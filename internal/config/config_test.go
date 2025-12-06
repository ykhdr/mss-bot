package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_ValidConfig(t *testing.T) {
	content := `
bot {
    token "test-token-123"
}

database {
    path "./test.db"
}

minecraft {
    timeout "10s"
}
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.kdl")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)

	assert.Equal(t, "test-token-123", cfg.Bot.Token)
	assert.Equal(t, "./test.db", cfg.Database.Path)
	assert.Equal(t, 10*time.Second, cfg.Minecraft.Timeout)
}

func TestLoad_MissingToken(t *testing.T) {
	content := `
bot {
    token ""
}

database {
    path "./test.db"
}

minecraft {
    timeout "5s"
}
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.kdl")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	_, err = Load(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bot token is not configured")
}

func TestLoad_DefaultValues(t *testing.T) {
	content := `
bot {
    token "valid-token"
}

database {
}

minecraft {
}
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.kdl")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)

	assert.Equal(t, "./data/mss-bot.db", cfg.Database.Path)
	assert.Equal(t, 5*time.Second, cfg.Minecraft.Timeout)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.kdl")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

func TestLoad_InvalidKDL(t *testing.T) {
	content := `this is not valid kdl {{{{`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.kdl")
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	_, err = Load(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}
