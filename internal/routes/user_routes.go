package routes

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/yusuffugurlu/go-project/internal/cache"
	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	"github.com/yusuffugurlu/go-project/internal/services"
	"github.com/yusuffugurlu/go-project/pkg/middleware"
)

func RegisterUserRoutes(e *echo.Group, cacheService *cache.CacheService) {
	logService := services.NewAuditLogService(repositories.NewAuditLogRepository(database.Db))
	repo := repositories.NewUserRepository(database.Db)
	service := services.NewUserServiceWithCache(repo, logService, cacheService)
	controller := controllers.NewUserController(service)

	route := e.Group("/users")

	cacheConfig := middleware.CacheConfig{
		Duration: 5 * time.Minute,
		KeyFunc:  middleware.DefaultCacheKey,
	}
	cacheMiddleware := middleware.NewCacheMiddleware(cacheService, cacheConfig)
	route.Use(cacheMiddleware.Cache())

	//route.Use(middleware.RoleBasedAuth("admin"))

	route.GET("/", controller.GetAllUsers)
	route.GET("/:id", controller.GetUserById)
	route.POST("/create", controller.CreateUser)
	route.PUT("/:id", controller.UpdateUser)
	route.DELETE("/:id", controller.DeleteUser)
}
