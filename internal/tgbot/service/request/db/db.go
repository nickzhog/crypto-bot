package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/nickzhog/crypto-bot/internal/tgbot/service/request"
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

func (r *repository) Create(ctx context.Context, usrID int64, currencyName, price string) error {
	q := `
		INSERT INTO public.request 
		    (user_id, currency_name, price) 
		VALUES 
		    ($1, $2, $3)
	`
	_, err := r.client.Exec(ctx, q, usrID, currencyName, price)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error("err:", newErr.Error())
		}
		r.logger.Error("err:", err.Error())
	}
	return err
}

func (r *repository) FindForUser(ctx context.Context, usrID int64) ([]request.Request, error) {
	q := `
	SELECT 
		user_id, currency_name, price, create_at
	FROM 
		public.request
	WHERE 
		user_id = $1
	ORDER BY create_at DESC
	LIMIT 5;
`

	rows, err := r.client.Query(ctx, q, usrID)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	requests := make([]request.Request, 0, rows.CommandTag().RowsAffected())

	for rows.Next() {
		var req request.Request

		err = rows.Scan(&req.UserID, &req.CurrencyName, &req.Price, &req.CreateAt)

		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, request.ErrNoRows
		}
		r.logger.Error(err)
		return nil, err
	}

	if len(requests) < 1 {
		return nil, request.ErrNoRows
	}

	return requests, nil
}
