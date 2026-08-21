package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	api "github.com/yourname/smartmeal/api/openapi/ogen"
)

type Server struct {
	http *http.Server
}

func NewServer(addr string, h api.Handler, logger zerolog.Logger) (*Server, error) {
	ogenSrv, err := api.NewServer(h, api.WithMiddleware(ZerologMiddleware(logger)))
	if err != nil {
		return nil, err
	}

	return &Server{
		http: &http.Server{
			Addr:    addr,
			Handler: ogenSrv,
		},
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("addr", s.http.Addr).Msg("http server starting")
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		log.Info().Msg("shutdown signal received")
	case err := <-errCh:
		log.Error().Err(err).Msg("http server failed")
		return err
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.http.Shutdown(shutdownCtx)
}
