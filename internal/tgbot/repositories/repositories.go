package repositories

import (
	"context"
	"time"

	"github.com/nickzhog/crypto-bot/internal/tgbot/config"
	"github.com/nickzhog/crypto-bot/internal/tgbot/service/cryptocurrency"
	currencydb "github.com/nickzhog/crypto-bot/internal/tgbot/service/cryptocurrency/db"
	"github.com/nickzhog/crypto-bot/internal/tgbot/service/request"
	requestdb "github.com/nickzhog/crypto-bot/internal/tgbot/service/request/db"
	"github.com/nickzhog/crypto-bot/pkg/logging"
	"github.com/nickzhog/crypto-bot/pkg/postgres"
)

const (
	maxAttempts      = 3
	dbConnectTimeOut = time.Second * 5
)

type Repositories struct {
	Currency cryptocurrency.Repository
	Request  request.Repository
}

func GetRepositories(ctx context.Context, logger *logging.Logger, cfg *config.Config) Repositories {
	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeOut)
	defer cancel()

	pool, err := postgres.NewClient(ctx, maxAttempts, cfg.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}
	if err = pool.Ping(ctx); err != nil {
		logger.Fatal(err)
	}
	return Repositories{
		Currency: currencydb.NewRepository(pool, logger),
		Request:  requestdb.NewRepository(pool, logger),
	}
}
