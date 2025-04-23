package main

import (
	"github.com/spf13/viper"
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config"
	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/shutdown"
)

func main() {
	logger.InitializeLogger()
	config.InitializeConfig()

	port := viper.GetString("APP_PORT")
    if port == "" {
        port = "8080"
    }

	e := echo.New()
	e.Logger.Fatal(e.Start(":" + port))


	shutdown.Handle(func() {
        logger.Log.Info("Custom shutdown logic executed.")
    })
}