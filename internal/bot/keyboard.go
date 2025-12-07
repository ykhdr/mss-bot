package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Callback data constants
const (
	CallbackStatus   = "status"
	CallbackSettings = "settings"
	CallbackPlayers  = "players"
	CallbackBack     = "back"
	CallbackRefresh  = "refresh"
)

// MainMenuKeyboard returns the main menu inline keyboard
func MainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статус", CallbackStatus),
			tgbotapi.NewInlineKeyboardButtonData("👥 Игроки", CallbackPlayers),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", CallbackSettings),
		),
	)
}

// StatusKeyboard returns the status view inline keyboard
func StatusKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", CallbackRefresh),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", CallbackBack),
		),
	)
}

// SettingsKeyboard returns the settings view inline keyboard
func SettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", CallbackBack),
		),
	)
}

// PlayersKeyboard returns the players view inline keyboard
func PlayersKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", CallbackRefresh),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", CallbackBack),
		),
	)
}

// BackKeyboard returns a simple back button keyboard
func BackKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", CallbackBack),
		),
	)
}
