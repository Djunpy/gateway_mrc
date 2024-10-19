package main

import (
	"context"
	config2 "gateway_mrc/config"
	db "gateway_mrc/db/sqlc"
	"gateway_mrc/infrastructure/database"
	"gateway_mrc/infrastructure/server"
	logger2 "gateway_mrc/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := config2.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	logger, file := logger2.ConfigureLogger(config.Environment)
	defer file.Close()
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	connPool, err := pgxpool.New(ctx, config.PostgresSource)
	defer connPool.Close()
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = connPool.Ping(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = database.RunDBMigration(config.MigrationURL, config.PostgresSource)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot run migration")
	}
	store := db.NewStore(connPool)
	err = server.RunGinServer(config, store)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot run gin server")
	}
}