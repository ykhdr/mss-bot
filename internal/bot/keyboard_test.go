package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainMenuKeyboard(t *testing.T) {
	kb := MainMenuKeyboard()

	assert.Len(t, kb.InlineKeyboard, 1)
	assert.Len(t, kb.InlineKeyboard[0], 2)

	assert.Equal(t, "📊 Статус", kb.InlineKeyboard[0][0].Text)
	assert.Equal(t, CallbackStatus, *kb.InlineKeyboard[0][0].CallbackData)

	assert.Equal(t, "⚙️ Настройки", kb.InlineKeyboard[0][1].Text)
	assert.Equal(t, CallbackSettings, *kb.InlineKeyboard[0][1].CallbackData)
}

func TestStatusKeyboard(t *testing.T) {
	kb := StatusKeyboard()

	assert.Len(t, kb.InlineKeyboard, 2)

	assert.Equal(t, "🔄 Обновить", kb.InlineKeyboard[0][0].Text)
	assert.Equal(t, CallbackRefresh, *kb.InlineKeyboard[0][0].CallbackData)

	assert.Equal(t, "◀️ Назад", kb.InlineKeyboard[1][0].Text)
	assert.Equal(t, CallbackBack, *kb.InlineKeyboard[1][0].CallbackData)
}

func TestSettingsKeyboard(t *testing.T) {
	kb := SettingsKeyboard()

	assert.Len(t, kb.InlineKeyboard, 1)
	assert.Len(t, kb.InlineKeyboard[0], 1)

	assert.Equal(t, "◀️ Назад", kb.InlineKeyboard[0][0].Text)
	assert.Equal(t, CallbackBack, *kb.InlineKeyboard[0][0].CallbackData)
}

func TestBackKeyboard(t *testing.T) {
	kb := BackKeyboard()

	assert.Len(t, kb.InlineKeyboard, 1)
	assert.Len(t, kb.InlineKeyboard[0], 1)

	assert.Equal(t, "◀️ Назад", kb.InlineKeyboard[0][0].Text)
	assert.Equal(t, CallbackBack, *kb.InlineKeyboard[0][0].CallbackData)
}

func TestCallbackConstants(t *testing.T) {
	assert.Equal(t, "status", CallbackStatus)
	assert.Equal(t, "settings", CallbackSettings)
	assert.Equal(t, "back", CallbackBack)
	assert.Equal(t, "refresh", CallbackRefresh)
}
