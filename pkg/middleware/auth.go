package middleware

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"github.com/yusuffugurlu/go-project/pkg/jwt"
)

func RoleBasedAuth(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return appErrors.NewUnauthorized(nil, "authorization header is required")
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.ValidateJWT(tokenString, os.Getenv("JWT_SECRET_KEY"))
			if err != nil {
				return appErrors.NewUnauthorized(err, "invalid token")
			}

			userClaims, err := jwt.GetUserClaims(token)
			if err != nil {
				return appErrors.NewUnauthorized(err, "failed to parse user claims")
			}

			if userClaims.Role != requiredRole {
				return appErrors.NewUnauthorized(nil, "insufficient permissions")
			}

			c.Set("user", userClaims)

			return next(c)
		}
	}
}
