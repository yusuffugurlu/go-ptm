package routes

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/cache"
	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/internal/services"
	"github.com/yusuffugurlu/go-project/pkg/middleware"
)

func RegisterTransactionRoutes(e *echo.Group, cacheService *cache.CacheService) {
	service := services.NewTransactionServiceWithCache(cacheService)
	controller := controllers.NewTransactionControllerWithService(service)
	route := e.Group("/transactions")

	cacheConfig := middleware.CacheConfig{
		Duration: 2 * time.Minute,
		KeyFunc:  middleware.CacheByUserID,
	}
	cacheMiddleware := middleware.NewCacheMiddleware(cacheService, cacheConfig)
	route.Use(cacheMiddleware.Cache())

	// route.POST("/deposit", controller.Deposit)
	route.POST("/withdraw", controller.Withdraw)

	route.POST("/debit", controller.Debit, middleware.RoleBasedAuth("user"))
	route.POST("/transfer", controller.Transfer, middleware.RoleBasedAuth("user"))
	route.GET("/history", controller.GetHistory, middleware.RoleBasedAuth("user"))
	route.GET("/:id", controller.GetByID, middleware.RoleBasedAuth("user"))

	route.GET("/all", controller.GetAllTransactions, middleware.RoleBasedAuth("admin"))
}
