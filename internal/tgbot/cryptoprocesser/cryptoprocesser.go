package cryptoprocesser

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nickzhog/crypto-bot/internal/tgbot/config"
	"github.com/nickzhog/crypto-bot/internal/tgbot/service/cryptocurrency"
	"github.com/nickzhog/crypto-bot/pkg/logging"
)

type CryptoProcesser interface {
	StartScan(ctx context.Context) error
}

type cryptoProcesser struct {
	logger      *logging.Logger
	cfg         *config.Config
	currencyRep cryptocurrency.Repository
}

func NewProcesser(logger *logging.Logger, cfg *config.Config, currencyRep cryptocurrency.Repository) *cryptoProcesser {
	return &cryptoProcesser{
		logger:      logger,
		cfg:         cfg,
		currencyRep: currencyRep,
	}
}

func (p *cryptoProcesser) StartScan(ctx context.Context) error {
	ticker := time.NewTicker(p.cfg.UpdateInterval)
	for {
		select {
		case <-ticker.C:
			p.scan(ctx)
		case <-ctx.Done():
			p.logger.Trace("orders processing exited properly")
			return nil
		}
	}
}

type Answer struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (p *cryptoProcesser) scan(ctx context.Context) {

	answer, err := requestToBinance(ctx)
	if err != nil {
		p.logger.Error(err)
	}

	currencies := make([]cryptocurrency.Currency, len(answer)/10)
	for _, v := range answer {
		re := regexp.MustCompile(`^.+USDT$`)
		match := re.MatchString(v.Symbol)
		if !match {
			continue
		}

		v.Symbol = strings.TrimSuffix(v.Symbol, "USDT")

		priceFloat, err := strconv.ParseFloat(v.Price, 64)
		if err != nil {
			p.logger.Error(err)
			continue
		}

		currencies = append(currencies, cryptocurrency.NewCurrency(v.Symbol, priceFloat))
	}

	if len(currencies) < 1 {
		return
	}

	err = p.currencyRep.UpsertMany(ctx, currencies)
	if err != nil {
		p.logger.Error(err)
	}
}
