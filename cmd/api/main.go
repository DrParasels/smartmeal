package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yourname/smartmeal/api"
	natsout "github.com/yourname/smartmeal/output/nats"
	"github.com/yourname/smartmeal/service/handler"
	"github.com/yourname/smartmeal/storages/postgresql"
)

func main() {
	cfg := CfgLoad()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error().Err(err).Msg("postgres connect failed")
		os.Exit(1)
	}
	defer pool.Close()
	log.Info().Msg("postgres connected")

	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Error().Err(err).Msg("nats connect failed")
		os.Exit(1)
	}
	defer func() {
		if err := nc.Drain(); err != nil {
			log.Error().Err(err).Msg("nats drain failed")
		}
	}()

	log.Info().Msg("nats connected")

	store := postgresql.New(pool)
	pub := natsout.NewPublisher(nc, cfg.NATSSubject)
	h := handler.NewMealHandler(store, pub)

	srv, err := api.NewServer(cfg.HTTPAddr, h, logger)
	if err != nil {
		log.Error().Err(err).Msg("api server create failed")
		os.Exit(1)
	}
	if err := srv.Run(ctx); err != nil {
		log.Error().Err(err).Msg("http server failed")
		os.Exit(1)
	}
	log.Info().Msg("shutdown complete")
}
