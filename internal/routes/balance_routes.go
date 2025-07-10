package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/internal/controllers"
	"github.com/yusuffugurlu/go-project/pkg/middleware"
)

func RegisterBalanceRoutes(e *echo.Group) {
	controller := controllers.NewBalanceController()

	route := e.Group("/balances")

	route.GET("/current", controller.GetCurrentBalance, middleware.RoleBasedAuth("user"))
	route.GET("/historical", controller.GetHistoricalBalances, middleware.RoleBasedAuth("user"))
}
