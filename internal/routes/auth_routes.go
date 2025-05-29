package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	"github.com/yusuffugurlu/go-project/internal/services"
)


func RegisterAuthRoutes(e *echo.Group) {
	userService := services.NewUserService(
		repositories.NewUserRepository(database.Db),
		services.NewAuditLogService(repositories.NewAuditLogRepository(database.Db)),
	)
	autService := services.NewAuthService(userService)
	authController := controllers.NewAuthController(autService)

	route := e.Group("/auth")

	route.POST("/login", authController.Login)
	route.POST("/register", authController.Register)
}