package request

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Request struct {
	UserID       uint64
	CreateAt     time.Time
	CurrencyName string
	Price        string
}

var ErrNoRows = errors.New("no result")

func PrepareStatistic(ctx context.Context, usrID int64, rep Repository) (string, error) {
	requests, err := rep.FindForUser(ctx, usrID)
	if err != nil {
		if errors.Is(err, ErrNoRows) {
			return "Похоже, вы еще не делали запросов", nil
		}
		return "", err
	}

	text := "В статистике отображены 5 последних запросов от вас:\n\n"
	for _, v := range requests {
		text += fmt.Sprintf("Криптовалюта: %s, Цена(на момент запроса): %s, Дата и время: %s\n",
			v.CurrencyName, v.Price, v.CreateAt.Format("02.01 15:04"))
	}
	return text, nil
}
