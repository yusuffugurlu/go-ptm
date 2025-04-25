package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config/logger"
)

func StartGracefully(e *echo.Echo, port string) {
	go func() {
		logger.Log.Infof("Starting server on port %s", port)
		if err := e.Start(":" + port); err != nil {
			if err.Error() != "http: Server closed" {
				logger.Log.Errorf("Error starting server: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Log.Info("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Log.Info("Attempting to shut down the server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		logger.Log.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Log.Info("Server gracefully stopped.")
	}
}
