package api

import (
	"log/slog"
	"net/http"
)

type Config struct {
	Logger *slog.Logger

	HTTPClient *http.Client

	APIConf APIConfig
}

type APIConfig struct {
	BaseURL string
	Key     string

	SettleUpConf SettleUpConfig
}

type SettleUpConfig struct {
	Username string
	Password string

	GroupID string
}
