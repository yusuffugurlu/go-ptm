package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yusuffugurlu/go-project/config/logger"
)

type CacheService struct {
	redisClient *RedisClient
}

func NewCacheService(redisClient *RedisClient) *CacheService {
	return &CacheService{
		redisClient: redisClient,
	}
}

func (c *CacheService) Get(ctx context.Context, key string) (string, error) {
	return c.redisClient.Get(ctx, key)
}

func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.redisClient.Set(ctx, key, value, expiration)
}

func (c *CacheService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.redisClient.Set(ctx, key, jsonData, expiration)
}

func (c *CacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	jsonData, err := c.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonData), dest)
}

func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.redisClient.Delete(ctx, key)
}

func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	return c.redisClient.Exists(ctx, key)
}

func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.redisClient.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.redisClient.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (c *CacheService) WarmUpCache(ctx context.Context) error {
	logger.Log.Info("Starting cache warm-up process")

	logger.Log.Info("Cache warm-up completed")
	return nil
}

func (c *CacheService) GenerateCacheKey(prefix, identifier string) string {
	return fmt.Sprintf("%s:%s", prefix, identifier)
}
