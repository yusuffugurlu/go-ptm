package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func InitializeLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	//config.DisableStacktrace = true

	//config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Log = logger.Sugar()

	Log.Info("Logger initialized successfully")
}
