package main

import (
	"flag"
	"stregy/internal/app"
	"stregy/internal/config"
	"stregy/pkg/logging"

	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()

	logging.Init(logrus.DebugLevel)

	cfg := config.GetConfig()

	app.Run(cfg)
}
