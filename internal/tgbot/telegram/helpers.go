package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *telegramBot) sendMessage(msg tgbotapi.MessageConfig) error {
	_, err := t.bot.Send(msg)
	if err != nil {
		t.logger.Error(err)
	}
	return err
}
