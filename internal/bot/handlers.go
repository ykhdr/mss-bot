package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	"github.com/ykhdr/mss-bot/internal/minecraft"
	"github.com/ykhdr/mss-bot/internal/service"
	"github.com/ykhdr/mss-bot/internal/storage"
	"github.com/ykhdr/mss-bot/internal/storage/models"
)

// Handlers contains all bot command and callback handlers
type Handlers struct {
	bot          *tgbotapi.BotAPI
	service      *service.ServerService
	stateManager *StateManager
}

// NewHandlers creates a new handlers instance
func NewHandlers(bot *tgbotapi.BotAPI, svc *service.ServerService, sm *StateManager) *Handlers {
	return &Handlers{
		bot:          bot,
		service:      svc,
		stateManager: sm,
	}
}

// HandleCommand processes incoming commands
func (h *Handlers) HandleCommand(ctx context.Context, message *tgbotapi.Message) {
	switch message.Command() {
	case "mss":
		h.handleMSS(ctx, message)
	case "set":
		h.handleSet(ctx, message)
	case "start":
		h.handleStart(ctx, message)
	case "help":
		h.handleHelp(ctx, message)
	}
}

// HandleCallback processes inline keyboard callbacks
func (h *Handlers) HandleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	// Answer callback to remove loading state
	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	if _, err := h.bot.Request(callbackResponse); err != nil {
		log.Printf("Failed to answer callback: %v", err)
	}

	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	switch callback.Data {
	case CallbackStatus:
		h.showStatus(ctx, chatID, messageID)
	case CallbackSettings:
		h.showSettings(ctx, chatID, messageID)
	case CallbackBack:
		h.showMainMenu(ctx, chatID, messageID)
	case CallbackRefresh:
		h.showStatus(ctx, chatID, messageID)
	}
}

func (h *Handlers) handleStart(ctx context.Context, message *tgbotapi.Message) {
	text := "👋 Привет! Я бот для проверки статуса Minecraft серверов.\n\n" +
		"Используйте /mss для открытия меню."

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	if _, err := h.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("Failed to send start message")
	}
}

func (h *Handlers) handleHelp(ctx context.Context, message *tgbotapi.Message) {
	text := "📖 *Справка*\n\n" +
		"*Команды:*\n" +
		"/mss \\- Открыть главное меню\n" +
		"/mss\\-set \\<ip:port\\> \\<name\\> \\- Настроить сервер\n\n" +
		"*Пример:*\n" +
		"`/mss-set mc.example.com:25565 My Server`"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	if _, err := h.bot.Send(msg); err != nil {
		log.Error().Err(err).Msg("Failed to send help message")
	}
}

func (h *Handlers) handleMSS(ctx context.Context, message *tgbotapi.Message) {
	text := "🎮 *Minecraft Server Status*\n\nВыберите действие:"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.ReplyMarkup = MainMenuKeyboard()

	sent, err := h.bot.Send(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send main menu")
		return
	}

	h.stateManager.SetState(message.Chat.ID, StateMainMenu, sent.MessageID)
}

func (h *Handlers) handleSet(ctx context.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// Check if we're in settings state
	if !h.stateManager.IsInState(chatID, StateSettings) {
		msg := tgbotapi.NewMessage(chatID,
			"⚠️ Эта команда доступна только из меню настроек.\n"+
				"Используйте /mss и нажмите кнопку Настройки.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	// Parse arguments
	args := message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(chatID,
			"❌ Неверный формат.\n\n"+
				"Использование: `/mss-set <ip:port> <name>`\n"+
				"Пример: `/mss-set mc.example.com:25565 My Server`")
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		if _, err := h.bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	// Parse address and name
	parts := strings.SplitN(args, " ", 2)
	address := parts[0]
	name := ""
	if len(parts) > 1 {
		name = parts[1]
	}

	host, port, err := minecraft.ParseAddress(address)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Неверный адрес: %v", err))
		if _, err := h.bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	// Save server config
	if err := h.service.SetServerConfig(ctx, chatID, host, port, name); err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка сохранения: %v", err))
		if _, err := h.bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("Failed to send error message")
		}
		return
	}

	// Update settings message
	messageID := h.stateManager.GetMessageID(chatID)
	h.showSettings(ctx, chatID, messageID)

	// Send confirmation
	confirmMsg := tgbotapi.NewMessage(chatID, "✅ Сервер успешно настроен!")
	if _, err := h.bot.Send(confirmMsg); err != nil {
		log.Error().Err(err).Msg("Failed to send confirmation")
	}
}

func (h *Handlers) showMainMenu(ctx context.Context, chatID int64, messageID int) {
	text := "🎮 *Minecraft Server Status*\n\nВыберите действие:"

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = tgbotapi.ModeMarkdownV2
	edit.ReplyMarkup = pointerTo(MainMenuKeyboard())

	if _, err := h.bot.Send(edit); err != nil {
		log.Error().Err(err).Msg("Failed to edit message to main menu")
	}

	h.stateManager.SetState(chatID, StateMainMenu, messageID)
}

func (h *Handlers) showStatus(ctx context.Context, chatID int64, messageID int) {
	result, err := h.service.GetServerStatus(ctx, chatID)

	var text string
	if err != nil {
		if isNotFound(err) {
			text = "⚠️ Сервер не настроен\\.\n\nИспользуйте настройки для добавления сервера\\."
		} else {
			text = fmt.Sprintf("❌ Ошибка: %v", escapeMarkdownV2(err.Error()))
		}
	} else {
		text = result.FormatStatus()
	}

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = tgbotapi.ModeMarkdownV2
	edit.ReplyMarkup = pointerTo(StatusKeyboard())

	if _, err := h.bot.Send(edit); err != nil {
		if strings.Contains(err.Error(), "message is not modified") {
			// No need to log this as an error
			return
		}
		log.Error().Err(err).Msg("Failed to edit message to status")
	}

	h.stateManager.SetState(chatID, StateStatus, messageID)
}

func (h *Handlers) showSettings(ctx context.Context, chatID int64, messageID int) {
	server, err := h.service.GetServerConfig(ctx, chatID)

	var text string
	if err != nil && !isNotFound(err) {
		text = fmt.Sprintf("❌ Ошибка: %v", escapeMarkdownV2(err.Error()))
	} else {
		text = service.FormatConfig(server)
	}

	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = tgbotapi.ModeMarkdownV2
	edit.ReplyMarkup = pointerTo(SettingsKeyboard())

	if _, err := h.bot.Send(edit); err != nil {
		log.Error().Err(err).Msg("Failed to edit message to settings")
	}

	h.stateManager.SetState(chatID, StateSettings, messageID)
}

func isNotFound(err error) bool {
	_, ok := err.(storage.ErrNotFound)
	return ok
}

func pointerTo[T any](v T) *T {
	return &v
}

func escapeMarkdownV2(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(s)
}

// Ensure models is used
var _ = models.Server{}
