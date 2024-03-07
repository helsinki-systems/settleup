package main

import (
	"context"
	"errors"
	"flag" //nolint:depguard // We only allow to import the flag package in here
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"
	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/server"

	"github.com/joho/godotenv"
)

//nolint:gochecknoglobals // Nice to use as a global
var logTarget = os.Stderr

func run(ctx context.Context, c Config) error {
	ac := api.New(c.apiConf)

	if _, err := ac.Login(ctx); err != nil {
		return fmt.Errorf("failed to login at API: %w", err)
	}

	srv := server.New(c.serverConf, ac)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to listen and serve: %w", err)
	}

	return nil
}

func main() {
	httpTCPBind := flag.String("http.tcp.bind", ":8080", "the TCP socket to bind to")

	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	ctx := context.Background()

	ll := new(slog.LevelVar)
	ll.Set(slog.LevelInfo)
	l := slog.New(slog.NewJSONHandler(logTarget, &slog.HandlerOptions{
		Level: ll,
	}))
	slog.SetDefault(l)

	if err := godotenv.Load(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			l.Info("no .env file found, service may fail to start")
		} else {
			l.Error(
				"failed to load env",
				"err", err,
			)
		}
	}

	// We have a debug env var as well as a debug CLI flag
	if getenv("DEBUG", "false") == "true" {
		*debug = true
	}

	if *debug {
		ll.Set(slog.LevelDebug)
	}

	c := Config{
		serverConf: server.Config{
			Logger: l.With("svc", "server"),

			HTTPTCPBind: *httpTCPBind,
		},

		apiConf: api.Config{
			Logger: l.With("svc", "api"),

			HTTPClient: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
				},
			},

			APIConf: api.APIConfig{
				BaseURL: getenv("API_BASE_URL", ""),
				Key:     getenv("API_KEY", ""),

				SettleUpConf: api.SettleUpConfig{
					Username: getenv("SETTLEUP_USERNAME", ""),
					Password: getenv("SETTLEUP_PASSWORD", ""),

					GroupID: getenv("SETTLEUP_GROUP_ID", ""),
				},
			},
		},
	}

	if err := run(ctx, c); err != nil {
		l.Error(err.Error())
	}
}
