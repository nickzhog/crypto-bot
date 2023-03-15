package cryptoprocesser

import (
	"context"
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

func (p *cryptoProcesser) scan(ctx context.Context) {

}
