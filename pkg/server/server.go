package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"
)

const (
	httpReadTimeout  = 10 * time.Second
	httpWriteTimeout = 5 * time.Second
)

const (
	requestIDHeaderName = "X-Request-ID"
	requestIDLength     = 36 // UUIDv4
)

type Server struct {
	conf Config

	apiClient *api.Client
}

func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/members",
		makeAllMembersHandler(
			s.conf.Logger.With("hdl", "/members"),
			s.apiClient,
		),
	)
	mux.HandleFunc(
		"/transactions",
		makeTransactionsHandler(
			s.conf.Logger.With("hdl", "/transactions"),
			s.apiClient,
		),
	)

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
