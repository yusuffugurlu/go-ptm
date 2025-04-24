package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/yusuffugurlu/go-project/config/logger"
)

type Config struct {
	AppPort               string `mapstructure:"APP_PORT"`
	AppName               string `mapstructure:"APP_NAME"`
	DatabaseConnectionURL string `mapstructure:"DATABASE_CONNECTION_URL"`
}

func InitializeConfig() *Config {
	projectRoot, err := os.Getwd()
	if err != nil {
		logger.Log.Warnf("Could not get working directory: %v. Assuming './'", err)
		projectRoot = "."
	}

	envFilePath := filepath.Join(projectRoot, ".env")

	viper.SetConfigFile(envFilePath)

	err = viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Log.Warnf("Config file not found at %s; relying on environment variables or defaults", envFilePath)
		} else {
			logger.Log.Errorf("Error reading config file '%s': %v", envFilePath, err)
		}
	}

	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Log.Errorf("Unable to decode config into struct: %v", err)
	}

	if config.AppPort == "" {
		logger.Log.Warn("APP_PORT not found in config or environment, defaulting to 8080")
		config.AppPort = "8080"
	}

	logger.Log.Info("Config initialized successfully")

	return &config
}
