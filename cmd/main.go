package main

import (
	"garantex-monitor/config"
	"garantex-monitor/internal/app"
	"garantex-monitor/internal/infrastructure/logs"
)

func main() {
	conf := config.MustLoad()

	logger := logs.NewLogger(conf)

	app := app.NewApp(conf, logger)

	app.Bootstrap().Run()
}
