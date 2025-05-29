package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	"github.com/yusuffugurlu/go-project/internal/services"
	"github.com/yusuffugurlu/go-project/pkg/middleware"
)

func RegisterLogRoutes(e *echo.Group) {
	repository := repositories.NewAuditLogRepository(database.Db)
	service := services.NewAuditLogService(repository)
	controllers := controllers.NewLogController(service)
	
	route := e.Group("/logs")

	route.Use(middleware.RoleBasedAuth("admin"))

	route.GET("/", controllers.GetAllLogs)
}
