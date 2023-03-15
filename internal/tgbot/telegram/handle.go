package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *telegramBot) handleUpdate(u tgbotapi.Update) {
	if u.Message == nil && u.CallbackQuery == nil {
		return
	}

}
