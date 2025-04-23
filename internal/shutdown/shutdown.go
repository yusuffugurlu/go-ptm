package shutdown
import (
    "context"
    "os/signal"
    "syscall"
    "time"

   "github.com/yusuffugurlu/go-project/config/logger"
)

func Handle(onShutdown func()) {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    logger.Log.Info("Application is running. Press Ctrl+C to exit.")
    <-ctx.Done()

    logger.Log.Info("Shutting down gracefully...")
    onShutdown()
    time.Sleep(1 * time.Second)
    logger.Log.Info("Shutdown complete.")
}