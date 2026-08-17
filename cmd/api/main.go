package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	api "github.com/yourname/smartmeal/api/ogen"
	"github.com/yourname/smartmeal/internal/config"
	"github.com/yourname/smartmeal/internal/handler"
	"github.com/yourname/smartmeal/internal/logger"
	"github.com/yourname/smartmeal/internal/storages/sqlc"
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
	defer nc.Drain() // поменял с close

	query := sqlc.New(pool)
	h := handler.NewMealHandler(query, nc, cfg.NATSSubject)

	ogenSrv, err := api.NewServer(h)
	if err != nil {
		log.Error("api server create failed", "err", err)
		os.Exit(1)
	}

	httpSrv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: ogenSrv,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Info("http server starting", "addr", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-errCh:
		log.Error("http server failed", "err", err)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown failed", "err", err)
	}

	log.Info("shutdown complete")
}
