package server

import (
	"log/slog"
)

type Config struct {
	Logger *slog.Logger

	HTTPTCPBind string
}
