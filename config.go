package main

import (
	"github.com/helsinki-systems/settleup/pkg/api"
	"github.com/helsinki-systems/settleup/pkg/server"
)

type Config struct {
	serverConf server.Config

	apiConf api.Config
}
