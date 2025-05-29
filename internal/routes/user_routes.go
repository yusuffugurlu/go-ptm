package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	"github.com/yusuffugurlu/go-project/internal/services"
)

func RegisterUserRoutes(e *echo.Group) {
	logService := services.NewAuditLogService(repositories.NewAuditLogRepository(database.Db))
	repo := repositories.NewUserRepository(database.Db)
	service := services.NewUserService(repo, logService)
	controller := controllers.NewUserController(service)

	route := e.Group("/users")

	//route.Use(middleware.RoleBasedAuth("admin"))

	route.GET("/", controller.GetAllUsers)
	route.GET("/:id", controller.GetUserById)
	route.POST("/create", controller.CreateUser)
	route.PUT("/:id", controller.UpdateUser)
	route.DELETE("/:id", controller.DeleteUser)
}
