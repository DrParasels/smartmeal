package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/yourname/smartmeal/internal/config"
	"github.com/yourname/smartmeal/internal/logger"
	eventsv1 "github.com/yourname/smartmeal/internal/pb/events/v1"
	"github.com/yourname/smartmeal/internal/stats"
	"github.com/yourname/smartmeal/internal/storages/sqlc"
	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("postgres connect failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Error("NATS connect failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := nc.Drain(); err != nil {
			log.Error("nats drain failed", "err", err)
		}
	}()

	q := sqlc.New(pool)
	p := stats.NewMealStats(q)

	sub, err := nc.Subscribe(cfg.NATSSubject, func(msg *nats.Msg) {
		var event eventsv1.MealCreated
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			log.Error("bad protobuf", "err", err)
			return
		}
		if _, err := p.HandleMealCreated(ctx, &event); err != nil {
			log.Error("handle meal created failed", "err", err)
			return
		}

		log.Debug("meal processed", "calories", event.Calories)
	})
	if err != nil {
		log.Error("NATS subscribe failed", "err", err, "subject", cfg.NATSSubject)
		os.Exit(1)
	}
	defer sub.Unsubscribe()

	log.Info("waiting for messages", "subject", cfg.NATSSubject)
	<-ctx.Done()
	log.Info("shutting down")

	if err := sub.Unsubscribe(); err != nil {
		log.Error("unsubscribe failed", "err", err)
	}
}
