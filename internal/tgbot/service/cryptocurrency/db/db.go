package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/nickzhog/crypto-bot/internal/tgbot/service/cryptocurrency"
	"github.com/nickzhog/crypto-bot/pkg/logging"
	"github.com/nickzhog/crypto-bot/pkg/postgres"
)

type repository struct {
	client postgres.Client
	logger *logging.Logger
}

func NewRepository(client postgres.Client, logger *logging.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) UpsertMany(ctx context.Context, currencies []cryptocurrency.Currency) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := `
	INSERT 
	INTO metrics
		(name, price, last_update) 
	VALUES 
		($1, $2, $3)
	ON CONFLICT (name) DO UPDATE 
	SET price=$2, last_update=$3
	`

	batch := &pgx.Batch{}
	for _, v := range currencies {
		batch.Queue(q, v.Name, v.PriceToDollarString, time.Now())
	}

	result := tx.SendBatch(ctx, batch)
	err = result.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
