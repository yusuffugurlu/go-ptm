package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusuffugurlu/go-project/internal/cache"
)

func InitRoutes(e *echo.Echo, cacheService *cache.CacheService) {
	v1 := e.Group("/api/v1")

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	RegisterLogRoutes(v1)
	RegisterUserRoutes(v1, cacheService)
	RegisterAuthRoutes(v1)
	RegisterBalanceRoutes(v1)
	RegisterTransactionRoutes(v1, cacheService)
}
