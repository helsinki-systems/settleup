package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"

	huma "github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

const (
	httpReadTimeout  = 10 * time.Second
	httpWriteTimeout = 5 * time.Second
)

type Server struct {
	conf Config

	apiClient *api.Client
}

type (
	humaContext   huma.Context
	customContext struct {
		humaContext
		custom context.Context //nolint:containedctx // Required due to missing huma SetContext
	}
)

func (c *customContext) Context() context.Context {
	return c.custom
}

// TODO: https://github.com/danielgtaylor/huma/pull/275#issuecomment-1975073339
func setContextValue(ctx huma.Context, key, val any) huma.Context {
	return &customContext{
		humaContext: ctx,
		custom:      context.WithValue(ctx.Context(), key, val),
	}
}

type ctxType string

const (
	loggerKey ctxType = "logger"
	ridKey    ctxType = "rid"
)

func loggerFromRequest(ctx context.Context) *slog.Logger {
	l := ctx.Value(loggerKey)
	if l, ok := l.(*slog.Logger); ok {
		return l
	}
	panic("Failed to type-assert to logger. This should never happen.") //nolint:forbidigo // Will never happen
}

func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	hapi := humago.New(mux, huma.DefaultConfig("API", "1.0.0"))

	registerMiddlewares(hapi, s.conf.Logger)
	registerRoutes(hapi, s.apiClient)

	srv := &http.Server{
		Addr: s.conf.HTTPTCPBind,

		Handler: mux,

		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,

		ErrorLog: slog.NewLogLogger(s.conf.Logger.Handler(), slog.LevelError),
	}

	s.conf.Logger.Info("HTTP server listening on TCP " + s.conf.HTTPTCPBind)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to serve using %+v: %w", s.conf, err)
	}

	return nil
}

func New(c Config, apiClient *api.Client) *Server {
	return &Server{
		conf: c,

		apiClient: apiClient,
	}
}
