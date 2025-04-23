package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func InitializeLogger() {
    config := zap.NewDevelopmentConfig()
    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    logger, err := config.Build()
    if err != nil {
        panic(err)
    }

    Log = logger.Sugar()
    Log.Info("Logger initialized successfully")
}