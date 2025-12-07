package config

import (
	"fmt"
	"time"

	kdlconfig "github.com/ykhdr/kdl-config"
)

// Config represents the application configuration
type Config struct {
	Bot       BotConfig
	Database  DatabaseConfig
	Minecraft MinecraftConfig
	Logging   LoggingConfig
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level string
}

// BotConfig contains Telegram bot settings
type BotConfig struct {
	Token string
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Path string
}

// MinecraftConfig contains Minecraft query settings
type MinecraftConfig struct {
	Timeout time.Duration
}

// kdlConfig is the internal KDL structure for parsing
type kdlConfig struct {
	Bot       kdlBotConfig       `kdl:"bot"`
	Database  kdlDatabaseConfig  `kdl:"database"`
	Minecraft kdlMinecraftConfig `kdl:"minecraft"`
	Logging   kdlLoggingConfig   `kdl:"logging"`
}

type kdlLoggingConfig struct {
	Level string `kdl:"level"`
}

type kdlBotConfig struct {
	Token string `kdl:"token" required:"true"`
}

type kdlDatabaseConfig struct {
	Path string `kdl:"path"`
}

type kdlMinecraftConfig struct {
	Timeout string `kdl:"timeout"`
}

// Load reads and parses the KDL configuration file
func Load(path string) (*Config, error) {
	var kdlCfg kdlConfig

	loader := kdlconfig.NewLoader()
	if err := loader.Load(&kdlCfg, path); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	timeout, err := time.ParseDuration(kdlCfg.Minecraft.Timeout)
	if err != nil && kdlCfg.Minecraft.Timeout != "" {
		return nil, fmt.Errorf("invalid timeout format: %w", err)
	}

	cfg := &Config{
		Bot: BotConfig{
			Token: kdlCfg.Bot.Token,
		},
		Database: DatabaseConfig{
			Path: kdlCfg.Database.Path,
		},
		Minecraft: MinecraftConfig{
			Timeout: timeout,
		},
		Logging: LoggingConfig{
			Level: kdlCfg.Logging.Level,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that all required configuration values are set
func (c *Config) validate() error {
	if c.Bot.Token == "" || c.Bot.Token == "YOUR_TELEGRAM_BOT_TOKEN" {
		return fmt.Errorf("bot token is not configured")
	}

	if c.Database.Path == "" {
		c.Database.Path = "./data/mss-bot.db"
	}

	if c.Minecraft.Timeout == 0 {
		c.Minecraft.Timeout = 5 * time.Second
	}

	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}

	return nil
}

// String returns a string representation of the configuration (for logging)
func (c *Config) String() string {
	return fmt.Sprintf(
		"Bot.Token: [REDACTED], Database.Path: %s, Minecraft.Timeout: %s, Logging.Level: %s",
		c.Database.Path,
		c.Minecraft.Timeout,
		c.Logging.Level,
	)
}
