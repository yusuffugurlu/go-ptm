package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	RegisterLogRoutes(v1)
	RegisterUserRoutes(v1)
	RegisterAuthRoutes(v1)
	RegisterBalanceRoutes(v1)
	RegisterTransactionRoutes(v1)
}
