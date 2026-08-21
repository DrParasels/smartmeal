package api

import (
	"time"

	"github.com/ogen-go/ogen/middleware"
	"github.com/rs/zerolog"
	ogen "github.com/yourname/smartmeal/api/openapi/ogen"
)

func ZerologMiddleware(base zerolog.Logger) ogen.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		l := base.With().
			Str("operation_id", req.OperationID).
			Str("method", req.Raw.Method).
			Str("path", req.Raw.URL.Path).
			Logger()
		req.SetContext(l.WithContext(req.Context))

		start := time.Now()
		resp, err := next(req)

		ev := l.Info().Dur("duration", time.Since(start))
		if err != nil {
			ev = l.Error().Err(err).Dur("duration", time.Since(start))
		}
		ev.Msg("request completed")
		return resp, err
	}
}
