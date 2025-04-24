package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config/logger" // Logger'ı kullanmak için import
)

// StartGracefully starts the Echo server and handles graceful shutdown.
// It blocks until a shutdown signal is received and the server is stopped.
func StartGracefully(e *echo.Echo, port string) {
	// Sunucuyu bir goroutine içinde başlat
	go func() {
		logger.Log.Infof("Starting server on port %s", port)
		if err := e.Start(":" + port); err != nil {
			// Hata "http: Server closed" değilse logla
			// Bu hata Shutdown çağrıldığında normal olarak döner
			if err.Error() != "http: Server closed" {
				// e.Logger yerine kendi logger'ımızı kullanalım
				logger.Log.Errorf("Error starting server: %v", err)
				// Burada os.Exit(1) çağırmak yerine hatayı yukarı iletmek daha iyi olabilir,
				// ancak şimdilik loglamak yeterli. Uygulamanın ana akışı zaten sinyal bekliyor olacak.
			}
		}
	}()

	// Graceful shutdown için sinyal bekleyelim
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) ve SIGTERM sinyallerini dinle
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Sinyal gelene kadar burada bekle
	<-quit
	logger.Log.Info("Shutdown signal received, initiating graceful shutdown...")

	// Graceful shutdown için context ve timeout belirle (örneğin 10 saniye)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Echo sunucusunu kapatmayı dene
	logger.Log.Info("Attempting to shut down the server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		// e.Logger yerine kendi logger'ımızı kullanalım
		logger.Log.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Log.Info("Server gracefully stopped.")
	}
}
