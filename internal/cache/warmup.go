package cache

import (
	"context"
	"time"

	"github.com/yusuffugurlu/go-project/config/logger"
)

type WarmupService struct {
	cacheService *CacheService
}

func NewWarmupService(cacheService *CacheService) *WarmupService {
	return &WarmupService{
		cacheService: cacheService,
	}
}

type WarmupData struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
}

func (w *WarmupService) StartWarmup(ctx context.Context, data []WarmupData) error {
	logger.Log.Info("Starting cache warm-up process")

	for _, item := range data {
		select {
		case <-ctx.Done():
			logger.Log.Warn("Cache warm-up cancelled")
			return ctx.Err()
		default:
			if err := w.cacheService.SetJSON(ctx, item.Key, item.Value, item.Expiration); err != nil {
				logger.Log.Error("Failed to warm up cache for key", "key", item.Key, "error", err)
				continue
			}
			logger.Log.Debug("Cache warmed up", "key", item.Key)
		}
	}

	logger.Log.Info("Cache warm-up completed successfully")
	return nil
}

func (w *WarmupService) WarmupUserData(ctx context.Context, userID uint, userData interface{}) error {
	key := w.cacheService.GenerateCacheKey("user", string(rune(userID)))
	return w.cacheService.SetJSON(ctx, key, userData, 1*time.Hour)
}

func (w *WarmupService) WarmupTransactionData(ctx context.Context, transactionID uint, transactionData interface{}) error {
	key := w.cacheService.GenerateCacheKey("transaction", string(rune(transactionID)))
	return w.cacheService.SetJSON(ctx, key, transactionData, 30*time.Minute)
}

func (w *WarmupService) WarmupBalanceData(ctx context.Context, userID uint, balanceData interface{}) error {
	key := w.cacheService.GenerateCacheKey("balance", string(rune(userID)))
	return w.cacheService.SetJSON(ctx, key, balanceData, 15*time.Minute)
}

func (w *WarmupService) WarmupFrequentlyAccessedData(ctx context.Context) error {
	logger.Log.Info("Warming up frequently accessed data")
	warmupData := []WarmupData{
		{
			Key:        "system:config",
			Value:      map[string]interface{}{"version": "1.0.0", "environment": "production"},
			Expiration: 24 * time.Hour,
		},
		{
			Key:        "system:currencies",
			Value:      []string{"USD", "EUR", "TRY"},
			Expiration: 12 * time.Hour,
		},
	}

	return w.StartWarmup(ctx, warmupData)
}

func (w *WarmupService) ScheduleWarmup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				if err := w.WarmupFrequentlyAccessedData(ctx); err != nil {
					logger.Log.Error("Scheduled warm-up failed", err)
				}
				cancel()
			}
		}
	}()
}
