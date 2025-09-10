package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/yusuffugurlu/go-project/config/logger"
)

type Config struct {
	AppPort               string
	AppName               string
	DatabaseConnectionURL string
	RedisURL              string
	RedisPassword         string
}

func InitializeConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Warn("No .env file found or error loading .env, relying on actual environment variables")
	}

	config := &Config{
		AppPort:               os.Getenv("APP_PORT"),
		AppName:               os.Getenv("APP_NAME"),
		DatabaseConnectionURL: os.Getenv("DATABASE_CONNECTION_URL"),
		RedisURL:              os.Getenv("REDIS_URL"),
		RedisPassword:         os.Getenv("REDIS_PASSWORD"),
	}

	if config.AppPort == "" {
		logger.Log.Warn("APP_PORT not set, defaulting to 8080")
		config.AppPort = "8080"
	}

	if config.RedisURL == "" {
		logger.Log.Warn("REDIS_URL not set, defaulting to localhost:6379")
		config.RedisURL = "localhost:6379"
	}

	logger.Log.Info("Config initialized using os.Getenv and godotenv")

	return config
}
