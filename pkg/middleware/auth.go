package middleware

import (
	"github.com/labstack/echo/v4"
)

func RoleBasedAuth(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			return next(e)
		}
	}
}