package main

import (
	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"
	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/server"
)

type Config struct {
	serverConf server.Config

	apiConf api.Config
}
