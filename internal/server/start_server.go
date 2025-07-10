package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yusuffugurlu/go-project/config/logger"
	customMiddleware "github.com/yusuffugurlu/go-project/pkg/middleware"
	"golang.org/x/time/rate"
)


func StartServer(e *echo.Echo) {
	e.HTTPErrorHandler = customMiddleware.GlobalErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoPrometheus.MetricsMiddleware())
	e.Use(customMiddleware.PerformanceMetrics)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(2))))

	go func() {
		logger.Log.Infof("Starting server on port %s", os.Getenv("APP_PORT"))
		if err := e.Start(os.Getenv("APP_PORT")); err != nil {
			if err.Error() != "http: Server closed" {
				logger.Log.Info("Error starting server: ", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Log.Info("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	logger.Log.Info("Attempting to shut down the server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Log.Info("Server gracefully stopped.")
	}
}