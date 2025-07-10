package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config/logger"
)

func PerformanceMetrics(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        err := next(c)
        duration := time.Since(start)

        logger.Log.Info("%s %s took %v", c.Request().Method, c.Request().URL.Path, duration)

        return err
    }
}
