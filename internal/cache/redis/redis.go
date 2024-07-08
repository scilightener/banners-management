package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Status int

const (
	// StatusExists status indicates that an item persists in redis cache.
	StatusExists Status = iota
	// StatusNotFound status indicates that an item does not persist in redis cache.
	StatusNotFound
	// StatusNotExists status indicates that an item does not persist in redis cache, nor anywhere else.
	// It means, there's no need to try to query it from persistent data storage (db).
	StatusNotExists
)

type CacheItem[T any] struct {
	Value  T      `json:"value"`
	Status Status `json:"status"`
}

func NewCacheItem[T any](value T, status Status) *CacheItem[T] {
	return &CacheItem[T]{
		Value:  value,
		Status: status,
	}
}

type Cache struct {
	client *redis.Client
}

func NewCache(connString string) (*Cache, error) {
	opt, err := redis.ParseURL(connString)
	if err != nil {
		return nil, fmt.Errorf("cache.redis.NewCache: %w", err)
	}

	client := redis.NewClient(opt)
	return &Cache{client: client}, nil
}

func Set[T any](c *Cache, ctx context.Context, key string, item *CacheItem[T], exp time.Duration) error {
	value, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("cache.redis.Set: %w", err)
	}
	err = c.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return fmt.Errorf("cache.redis.Set: %w", err)
	}

	return nil
}

func Get[T any](c *Cache, ctx context.Context, key string) (*CacheItem[T], error) {
	v, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var t T
			return NewCacheItem[T](t, StatusNotFound), nil
		}

		return nil, fmt.Errorf("cache.redis.Get: %w", err)
	}

	item := new(CacheItem[T])
	err = json.Unmarshal([]byte(v), item)
	if err != nil {
		return nil, fmt.Errorf("cache.redis.Get: %w", err)
	}

	return item, nil
}
