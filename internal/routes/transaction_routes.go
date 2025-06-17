package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/controllers"
)


func RegisterTransactionRoutes(e *echo.Group) {
	controller := controllers.NewTransactionController()

	route := e.Group("/transaction")

	route.POST("/deposit", controller.Deposit)
}