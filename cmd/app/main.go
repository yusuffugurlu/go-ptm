package main

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config"
	"github.com/yusuffugurlu/go-project/config/logger"

	"github.com/yusuffugurlu/go-project/internal/cache"
	"github.com/yusuffugurlu/go-project/internal/database"
	"github.com/yusuffugurlu/go-project/internal/process"
	"github.com/yusuffugurlu/go-project/internal/routes"
	"github.com/yusuffugurlu/go-project/internal/server"
	"github.com/yusuffugurlu/go-project/pkg/validator"
)

func main() {
	e := echo.New()
	e.Validator = validator.New()

	logger.InitializeLogger()
	cfg := config.InitializeConfig()
	database.InitializeDb()

	redisClient := cache.NewRedisClient(cfg)
	if redisClient == nil {
		logger.Log.Fatal("Failed to initialize Redis client")
	}
	defer redisClient.Close()

	cacheService := cache.NewCacheService(redisClient)
	warmupService := cache.NewWarmupService(cacheService)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := warmupService.WarmupFrequentlyAccessedData(ctx); err != nil {
			logger.Log.Error("Cache warm-up failed", err)
		}

		warmupService.ScheduleWarmup(1 * time.Hour)
	}()

	process.InitWorkerPool(10)

	routes.InitRoutes(e, cacheService)
	server.StartServer(e)
}
