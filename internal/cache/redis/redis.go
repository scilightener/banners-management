package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Status is a status of cache item in redis cache.
type Status int

const (
	// StatusExists status indicates that an item persists in redis cache.
	StatusExists Status = iota
	// StatusNotFound status indicates that an item does not persist in redis cache.
	StatusNotFound
	// StatusNotExists status indicates that an item does not persist in redis cache, nor anywhere else.
	// It means, there's no need to try to query it from persistent data storage (db).
	// See: https://en.wikipedia.org/wiki/Negative_cache.
	StatusNotExists
)

// CacheItem is a struct containing the Value that needs to be stored or restored in/from redis cache
// and a Status. For more information, check descriptions: StatusExists, StatusNotFound, StatusNotExists.
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

// Cache is a struct containing redis client. It is passed to methods Set, Get.
type Cache struct {
	client *redis.Client
}

// NewCache parses provided connection string and returns a ready-to-use redis client.
func NewCache(ctx context.Context, connString string) (*Cache, error) {
	opt, err := redis.ParseURL(connString)
	if err != nil {
		return nil, fmt.Errorf("cache.redis.NewCache: %w", err)
	}

	client := redis.NewClient(opt)
	err = client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return &Cache{client: client}, nil
}

// Set serializes the item into a json struct and sets this string in redis cache by the provided key.
// Note: it is not a method of Cache, but a function that accepts it. It is because for now methods can't be generic.
// See: https://github.com/golang/go/issues/49085.
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

// Get retrieves the value by the given key from the redis cache and deserializes it into a go struct of type T.
// Note: it is not a method of Cache, but a function that accepts it. It is because for now methods can't be generic.
// See: https://github.com/golang/go/issues/49085.
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
