package config

import (
	"os"

	"github.com/spf13/viper"
	"github.com/joho/godotenv"
	"github.com/yusuffugurlu/go-project/config/logger"
)

func InitializeConfig() {
	err := godotenv.Load()
    if err != nil {
		logger.Log.Error("Error loading .env file")
	}
	
    viper.SetDefault("APP_NAME", os.Getenv("APP_NAME"))
    viper.SetDefault("APP_PORT", os.Getenv("APP_PORT"))
	viper.SetDefault("DATABASE_CONNECTION_URL", os.Getenv("DATABASE_CONNECTION_URL"))

	logger.Log.Info("Config initialized successfully")
}