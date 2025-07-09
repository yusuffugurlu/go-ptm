package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/pkg/middleware"
)

func RegisterTransactionRoutes(e *echo.Group) {
	controller := controllers.NewTransactionController()
	route := e.Group("/transactions")

	// route.POST("/deposit", controller.Deposit)
	route.POST("/withdraw", controller.Withdraw)

	route.POST("/debit", controller.Debit, middleware.RoleBasedAuth("user"))
	route.POST("/transfer", controller.Transfer, middleware.RoleBasedAuth("user"))
	route.GET("/history", controller.GetHistory, middleware.RoleBasedAuth("user"))
	route.GET("/:id", controller.GetByID, middleware.RoleBasedAuth("user"))

	route.GET("/all", controller.GetAllTransactions, middleware.RoleBasedAuth("admin"))
}
