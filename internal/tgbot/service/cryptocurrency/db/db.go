package db

import (
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
