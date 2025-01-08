package main

import (
	"garantex-monitor/config"
	"garantex-monitor/internal/infrastructure/logs"

	"go.uber.org/zap"
)

func main() {

	conf := config.MustLoad()

	logger := logs.NewLogger(conf)

	logger.Info("logger", zap.String("key", "value"))
}
