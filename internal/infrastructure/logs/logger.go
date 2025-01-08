package logs

import (
	"errors"
	"garantex-monitor/config"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"
)

var (
	ErrorLoggerConf = errors.New("failed configure logger")
)

func NewLogger(conf *config.Config) *zap.Logger {
	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	switch conf.Env {
	case Local:
		// Level: DEBUG; Output: stdout; fmt: console
		zapConf.OutputPaths = []string{"stdout"}
		zapConf.ErrorOutputPaths = []string{"stderr"}
		zapConf.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapConf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		zapConf.EncoderConfig.ConsoleSeparator = "   "
		zapConf.Encoding = "console"
	case Dev:
		// Level: DEBUG; Output: stdout, file; fmt: json
		zapConf.OutputPaths = []string{"stdout", "logs.log"}
		zapConf.ErrorOutputPaths = []string{"stderr", "logs.log"}
		zapConf.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case Prod:
		// Level: INFO; Output: stdout, file; fmt: json
		zapConf.OutputPaths = []string{"stdout", "logs.log"}
		zapConf.ErrorOutputPaths = []string{"stderr", "logs.log"}
		zapConf.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	default:
		log.Fatalf("error logger configuration for env=%s", conf.Env)
	}

	logger, err := zapConf.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger.Named(conf.Name)
}
