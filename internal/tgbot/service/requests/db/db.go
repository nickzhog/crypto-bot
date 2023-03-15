package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
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

func (r *repository) Create(ctx context.Context, usrID int64, cryptoID string) error {
	q := `
		INSERT INTO public.requests 
		    (user_id, crypto_id) 
		VALUES 
		    ($1, $2)
	`
	_, err := r.client.Exec(ctx, q, usrID, cryptoID)
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
