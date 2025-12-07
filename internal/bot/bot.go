package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ykhdr/mss-bot/internal/service"
)

// Bot represents the Telegram bot
type Bot struct {
	api          *tgbotapi.BotAPI
	handlers     *Handlers
	stateManager *StateManager
}

// New creates a new bot instance
func New(token string, svc *service.ServerService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	sm := NewStateManager()
	handlers := NewHandlers(api, svc, sm)

	return &Bot{
		api:          api,
		handlers:     handlers,
		stateManager: sm,
	}, nil
}

// Start begins processing updates
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			b.processUpdate(ctx, &update)
		}
	}
}

// Stop gracefully stops the bot
func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}

func (b *Bot) processUpdate(ctx context.Context, update *tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.IsCommand() {
			b.handlers.HandleCommand(ctx, update.Message)
		}
	}

	if update.CallbackQuery != nil {
		b.handlers.HandleCallback(ctx, update.CallbackQuery)
	}
}
