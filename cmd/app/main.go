package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusuffugurlu/go-project/config"
	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/db"
	"github.com/yusuffugurlu/go-project/internal/server"

	echoPrometheus "github.com/globocom/echo-prometheus"
)

func main() {
	logger.InitializeLogger()

	cfg := config.InitializeConfig()

	db.Connect(cfg.DatabaseConnectionURL)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(echoPrometheus.MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	server.StartGracefully(e, cfg.AppPort)
}
