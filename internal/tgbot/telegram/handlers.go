package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nickzhog/crypto-bot/internal/tgbot/service/request"
)

const (
	MsgHello = `
Привет, %s %s
Данный бот позволяет получить актуальный курс для криптовалют, по отношению к доллару
Доступные команды:
/currencies - Выбор криптовалюты
/statistics - Ваша статистика`
)

func (t *telegramBot) sendHelloMsg(u tgbotapi.Update) {
	text := fmt.Sprintf(MsgHello, u.Message.From.FirstName, u.Message.From.LastName)
	msg := tgbotapi.NewMessage(u.Message.From.ID, text)
	t.sendMessage(msg)
}

func (t *telegramBot) sendError(usrID int64) {
	text := "Извините, произошла ошибка, попробуйте еще раз"
	msg := tgbotapi.NewMessage(usrID, text)
	t.sendMessage(msg)
}

func (t *telegramBot) handleCallbackData(ctx context.Context, c tgbotapi.CallbackQuery) {

	if len(c.Data) < 1 {
		t.sendError(c.From.ID)
		return
	}

	// Проверяем, является ли callbackData запросом на смену страницы
	re := regexp.MustCompile(`^page_(\d+)_(prev|next)$`)
	match := re.FindStringSubmatch(c.Data)
	if len(match) == 3 {
		// Это запрос на смену страницы, извлекаем индекс
		index, err := strconv.Atoi(match[1])
		if err != nil {
			t.sendError(c.From.ID)
			return
		}

		// Обрабатываем запрос на смену страницы
		if match[2] == "prev" {
			t.showCurrencies(ctx, c.From.ID, index-1)
		} else if match[2] == "next" {
			t.showCurrencies(ctx, c.From.ID, index+1)
		}

		return
	}

	currency, err := t.reps.Currency.FindOne(ctx, c.Data)
	if err != nil {
		t.sendError(c.From.ID)
		return
	}

	err = t.reps.Request.Create(ctx, c.From.ID,
		currency.Name, currency.PriceToDollarString)
	if err != nil {
		t.sendError(c.From.ID)
		return
	}

	text := fmt.Sprintf("Цена за 1 монету %s - %s долларов\nОбновлено: %s",
		currency.Name, currency.PriceToDollarString, currency.LastUpdate.Format("02.01 15:04"))
	msg := tgbotapi.NewMessage(c.From.ID, text)
	t.sendMessage(msg)
}

func (t *telegramBot) showCurrencies(ctx context.Context, chatID int64, index int) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	currencies, err := t.reps.Currency.FindCurrencies(ctx, index)
	if err != nil {
		t.logger.Error(err)
		t.sendError(chatID)
		return
	}

	for _, v := range currencies {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				v.Name, v.Name),
		))
	}

	prevBtn := tgbotapi.NewInlineKeyboardButtonData("<< Назад", fmt.Sprintf("page_%d_prev", index-1))
	nextBtn := tgbotapi.NewInlineKeyboardButtonData("Вперед >>", fmt.Sprintf("page_%d_next", index+1))

	row := tgbotapi.NewInlineKeyboardRow()
	if index > 0 {
		row = append(row, prevBtn)
	}

	if len(currencies) >= 10 {
		row = append(row, nextBtn)
	}

	if len(row) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg := tgbotapi.NewMessage(chatID, "Выберите криптовалюту:")
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func (t *telegramBot) showStatistics(ctx context.Context, u tgbotapi.Update) {
	answer, err := request.PrepareStatistic(ctx, u.Message.From.ID, t.reps.Request)
	if err != nil {
		t.sendError(u.Message.From.ID)
		return
	}

	msg := tgbotapi.NewMessage(u.Message.From.ID, answer)
	t.sendMessage(msg)
}
