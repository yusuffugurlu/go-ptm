package middleware

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/cache"
)

type CacheConfig struct {
	Duration time.Duration
	KeyFunc  func(c echo.Context) string
}

type CacheMiddleware struct {
	cacheService *cache.CacheService
	config       CacheConfig
}

func NewCacheMiddleware(cacheService *cache.CacheService, config CacheConfig) *CacheMiddleware {
	return &CacheMiddleware{
		cacheService: cacheService,
		config:       config,
	}
}

func DefaultCacheKey(c echo.Context) string {
	path := c.Request().URL.Path
	query := c.Request().URL.RawQuery
	method := c.Request().Method

	key := fmt.Sprintf("%s:%s?%s", method, path, query)
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("cache:%x", hash)
}

func (cm *CacheMiddleware) Cache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			key := cm.config.KeyFunc(c)
			if key == "" {
				key = DefaultCacheKey(c)
			}

			ctx := context.Background()
			cachedResponse, err := cm.cacheService.Get(ctx, key)
			if err == nil && cachedResponse != "" {
				logger.Log.Debug("Cache hit", "key", key)
				return c.JSONBlob(http.StatusOK, []byte(cachedResponse))
			}

			logger.Log.Debug("Cache miss", "key", key)

			rec := &responseRecorder{
				ResponseWriter: c.Response().Writer,
				body:           make([]byte, 0),
			}
			c.Response().Writer = rec

			err = next(c)
			if err != nil {
				return err
			}

			if c.Response().Status == http.StatusOK {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					if err := cm.cacheService.Set(ctx, key, string(rec.body), cm.config.Duration); err != nil {
						logger.Log.Error("Failed to cache response", err)
					} else {
						logger.Log.Debug("Response cached", "key", key)
					}
				}()
			}

			return nil
		}
	}
}

type responseRecorder struct {
	http.ResponseWriter
	body []byte
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

func CacheByUserID(c echo.Context) string {
	userID := c.Get("user_id")
	if userID == nil {
		return DefaultCacheKey(c)
	}

	path := c.Request().URL.Path
	query := c.Request().URL.RawQuery
	key := fmt.Sprintf("user:%v:%s?%s", userID, path, query)
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("cache:%x", hash)
}

func CacheByQueryParams(params ...string) func(echo.Context) string {
	return func(c echo.Context) string {
		path := c.Request().URL.Path
		query := c.Request().URL.Query()

		var keyParts []string
		keyParts = append(keyParts, path)

		for _, param := range params {
			if value := query.Get(param); value != "" {
				keyParts = append(keyParts, fmt.Sprintf("%s=%s", param, value))
			}
		}

		key := fmt.Sprintf("%s", keyParts)
		hash := md5.Sum([]byte(key))
		return fmt.Sprintf("cache:%x", hash)
	}
}
