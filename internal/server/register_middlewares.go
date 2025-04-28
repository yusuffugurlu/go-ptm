package server

import (
	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	customMiddleware "github.com/yusuffugurlu/go-project/pkg/middleware"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.HTTPErrorHandler = customMiddleware.GlobalErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoPrometheus.MetricsMiddleware())
}
