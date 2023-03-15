package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nickzhog/crypto-bot/internal/tgbot/repositories"
	"github.com/nickzhog/crypto-bot/pkg/logging"
)

type TelegramBot interface {
	StartUpdates(ctx context.Context) error
}

type telegramBot struct {
	bot    *tgbotapi.BotAPI
	logger *logging.Logger
	reps   repositories.Repositories
}

func PrepareBot(token string, logger *logging.Logger, reps repositories.Repositories) *telegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Tracef("Authorized on account %s", bot.Self.UserName)

	return &telegramBot{
		bot:    bot,
		logger: logger,
		reps:   reps,
	}
}

func (t *telegramBot) StartUpdates(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			t.handleUpdate(update)

		case <-ctx.Done():
			return nil
		}
	}
}
