package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	eventsv1 "github.com/yourname/smartmeal/events"
	"github.com/yourname/smartmeal/service/stats"
	"github.com/yourname/smartmeal/storages/postgresql"
	"google.golang.org/protobuf/proto"
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

	store := postgresql.New(pool)
	p := stats.NewMealStats(store)

	sub, err := nc.Subscribe(cfg.NATSSubject, func(msg *nats.Msg) {
		var event eventsv1.MealCreated
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			log.Error().Err(err).Msg("bad protobuf")
			return
		}

		msgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		msgCtx = logger.With().
			Int64("meal_id", event.MealId).
			Int32("calories", event.Calories).
			Logger().
			WithContext(msgCtx)

		if err := p.HandleMealCreated(msgCtx, &event); err != nil {
			log.Error().Err(err).Msg("handle meal created failed")
			return
		}
		log.Debug().Int32("calories", event.Calories).Msg("meal processed")
	})
	if err != nil {
		log.Error().Str("subject", cfg.NATSSubject).Err(err).Msg("NATS subscribe failed")
		os.Exit(1)
	}
	defer sub.Unsubscribe()

	log.Info().Str("subject", cfg.NATSSubject).Msg("waiting for messages")
	<-ctx.Done()
	log.Info().Msg("shutting down")

	if err := sub.Unsubscribe(); err != nil {
		log.Error().Err(err).Msg("unsubscribe failed")
	}
}
