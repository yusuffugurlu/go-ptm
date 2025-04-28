package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	RegisterUserRoutes(v1)
}