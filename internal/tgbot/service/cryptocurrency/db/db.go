package db

import (
	"context"
	"strconv"
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
	INTO public.cryptocurrency
		(currency_name, price, last_update) 
	VALUES 
		($1, $2, $3)
	ON CONFLICT (currency_name) DO UPDATE 
	SET price=$2, last_update=$3
	`

	batch := new(pgx.Batch)
	for _, v := range currencies {
		if len(v.PriceToDollarString) < 2 {
			continue
		}
		batch.Queue(q, v.Name, v.PriceToDollarString, time.Now())
	}

	result := tx.SendBatch(ctx, batch)
	err = result.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repository) FindCurrencies(ctx context.Context, index int) ([]cryptocurrency.Currency, error) {
	limit := 10 // количество записей на странице
	offset := index * limit

	q := `
    SELECT 
		currency_name, price
    FROM 
		public.cryptocurrency
    WHERE 
		last_update >= NOW() - INTERVAL '1 hour'
    ORDER BY 
		currency_name ASC
    LIMIT $1
    OFFSET $2
    `

	rows, err := r.client.Query(ctx, q, limit, offset)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	currencies := make([]cryptocurrency.Currency, 0, rows.CommandTag().RowsAffected())

	for rows.Next() {
		var c cryptocurrency.Currency

		err = rows.Scan(&c.Name, &c.PriceToDollarString)

		if err != nil {
			r.logger.Error(err)
			continue
		}

		if c.PriceToDollarFloat, err = strconv.ParseFloat(c.PriceToDollarString, 64); err != nil {
			r.logger.Error(err)
			continue
		}

		currencies = append(currencies, c)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return currencies, nil
}

func (r *repository) FindOne(ctx context.Context, name string) (cryptocurrency.Currency, error) {
	q := `
	SELECT
		currency_name, price, last_update
	FROM 
		public.cryptocurrency 
	WHERE 
		currency_name = $1
	`

	var c cryptocurrency.Currency
	err := r.client.QueryRow(ctx, q, name).
		Scan(&c.Name, &c.PriceToDollarString, &c.LastUpdate)
	if err != nil {
		return cryptocurrency.Currency{}, err
	}

	if c.PriceToDollarFloat, err = strconv.ParseFloat(c.PriceToDollarString, 64); err != nil {
		return cryptocurrency.Currency{}, err
	}

	return c, nil
}
