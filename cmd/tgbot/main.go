package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/nickzhog/crypto-bot/internal/tgbot/config"
	"github.com/nickzhog/crypto-bot/internal/tgbot/cryptoprocesser"
	"github.com/nickzhog/crypto-bot/internal/tgbot/repositories"
	"github.com/nickzhog/crypto-bot/internal/tgbot/telegram"
	"github.com/nickzhog/crypto-bot/migration"
	"github.com/nickzhog/crypto-bot/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()
	logger.Tracef("config: %+v", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Tracef("system call:%+v", oscall)
		cancel()
	}()

	err := migration.Migrate(cfg.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}

	reps := repositories.GetRepositories(ctx, logger, cfg)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		err := cryptoprocesser.NewProcesser(logger, cfg, reps.Currency).StartScan(ctx)
		if err != nil {
			logger.Errorf("crypto processer error: %s", err.Error())
		}
		wg.Done()
	}()

	go func() {
		bot := telegram.PrepareBot(cfg.TgToken, logger, reps)
		if err := bot.StartUpdates(ctx); err != nil {
			logger.Error(err)
		}
		wg.Done()
	}()

	wg.Wait()

	logger.Trace("graceful shutdown")
}
