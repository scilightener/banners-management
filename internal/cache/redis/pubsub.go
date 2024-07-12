package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (c *Cache) Publish(ctx context.Context, channel string, message any) error {
	const comp = "cache.redis.pubsub.Publish"
	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("%s: %w", comp, err)
	}

	err = c.client.Publish(ctx, channel, bytes).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", comp, err)
	}

	return nil
}

func (c *Cache) Subscribe(ctx context.Context, channel string) <-chan *redis.Message {
	subscriber := c.client.Subscribe(ctx, channel)
	return subscriber.Channel()
}
