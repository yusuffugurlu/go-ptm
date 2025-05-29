package main

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config"
	"github.com/yusuffugurlu/go-project/config/logger"

	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/routes"
	"github.com/yusuffugurlu/go-project/internal/server"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

func main() {
	e := echo.New()
	e.Validator = validator.New()

	logger.InitializeLogger()
	config.InitializeConfig()
	database.InitializeDb()

	routes.InitRoutes(e)
	server.StartServer(e)
}