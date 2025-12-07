package app

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ykhdr/mss-bot/internal/bot"
	"github.com/ykhdr/mss-bot/internal/config"
	"github.com/ykhdr/mss-bot/internal/logging"
	"github.com/ykhdr/mss-bot/internal/minecraft"
	"github.com/ykhdr/mss-bot/internal/service"
	"github.com/ykhdr/mss-bot/internal/storage"
	"github.com/ykhdr/mss-bot/internal/storage/sqlite"
)

// App represents the application
type App struct {
	cfg     *config.Config
	storage storage.ServerStorage
	bot     *bot.Bot
	cancel  context.CancelFunc
}

// New creates a new application instance
func New(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logging
	logging.Setup(cfg.Logging)
	log.Info().Msg("loaded config\n" + cfg.String())

	// Initialize storage
	store, err := sqlite.New(cfg.Database.Path)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize storage")
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize Minecraft client
	mcClient := minecraft.NewClient(cfg.Minecraft.Timeout)

	// Initialize service
	svc := service.NewServerService(store, mcClient)

	// Initialize bot
	b, err := bot.New(cfg.Bot.Token, svc)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize bot")
		store.Close()
		return nil, fmt.Errorf("failed to initialize bot: %w", err)
	}

	log.Info().Msg("successful initialization")

	return &App{
		cfg:     cfg,
		storage: store,
		bot:     b,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	log.Info().Msg("starting bot")
	return a.bot.Start(ctx)
}

// Shutdown gracefully stops the application
func (a *App) Shutdown() error {
	log.Info().Msg("shutting down application")

	if a.cancel != nil {
		a.cancel()
	}

	if a.bot != nil {
		a.bot.Stop()
		log.Debug().Msg("bot stopped")
	}

	if a.storage != nil {
		if err := a.storage.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close storage")
			return err
		}
		log.Debug().Msg("storage closed")
	}

	log.Info().Msg("application shutdown complete")
	return nil
}
