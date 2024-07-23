package banner

import (
	"context"
	"encoding/json"
	"log/slog"

	goredis "github.com/redis/go-redis/v9"

	"banners-management/internal/cache/redis"
	"banners-management/internal/lib/logger/sl"
	"banners-management/internal/storage/repo"
)

const (
	RedisBannerDeleterByFeatureTagChannelName = "redis_deleter_job_delayer"
)

// redisDeleteMessage is a dto for a RedisChannelDeleter.
type redisDeleteMessage struct {
	FeatureID int64 `json:"feature_id"`
	TagID     int64 `json:"tag_id"`
}

// RedisChannelDeleter is a decorator for repo.BannerDeleter that allows an asynchronous operation executing.
type RedisChannelDeleter struct {
	cache   *redis.Cache
	deleter repo.BannerDeleter
	logger  *slog.Logger
}

// NewRedisChannelDeleter returns a new RedisChannelDeleter instance.
func NewRedisChannelDeleter(
	ctx context.Context,
	redis *redis.Cache,
	deleter repo.BannerDeleter,
	logger *slog.Logger,
) *RedisChannelDeleter {
	res := &RedisChannelDeleter{
		cache:   redis,
		deleter: deleter,
		logger:  logger,
	}

	go res.runDeleterDaemon(ctx, res.cache.Subscribe(ctx, RedisBannerDeleterByFeatureTagChannelName))

	return res
}

// DeleteBanner does nothing and just proxies the request to the decorated repo.BannerDeleter.
func (r *RedisChannelDeleter) DeleteBanner(ctx context.Context, bannerID int64) error {
	return r.deleter.DeleteBanner(ctx, bannerID)
}

// DeleteByFeatureTag asynchronously executes operation of banner deletion by featureID and tagID.
func (r *RedisChannelDeleter) DeleteByFeatureTag(ctx context.Context, featureID, tagID int64) error {
	const comp = "service.banner.job_delayer"

	message := redisDeleteMessage{featureID, tagID}
	err := r.cache.Publish(ctx, RedisBannerDeleterByFeatureTagChannelName, message)
	if err != nil {
		r.logger.Error("error publishing message to redis", slog.String("comp", comp), sl.Err(err))
		return ErrUnknown
	}

	return nil
}

// runDeleterDaemon reads channel ch and executes all the operations of banner deletion, received from this channel.
func (r *RedisChannelDeleter) runDeleterDaemon(ctx context.Context, ch <-chan *goredis.Message) {
	for m := range ch {
		if m.Channel != RedisBannerDeleterByFeatureTagChannelName {
			continue
		}

		go func() {
			r.handleReceivedPayload(ctx, []byte(m.Payload))
		}()
	}
}

// handleReceivedPayload handles the received payload from redis channel.
func (r *RedisChannelDeleter) handleReceivedPayload(ctx context.Context, payload []byte) {
	res := new(redisDeleteMessage)
	err := json.Unmarshal(payload, res)
	if err != nil {
		r.logger.Error("unable to parse payload", sl.Err(err))
		return
	}

	err = r.deleter.DeleteByFeatureTag(ctx, res.FeatureID, res.TagID)
	if err != nil {
		r.logger.Error("unable to delete banner by feature & tag",
			slog.Int64("featureID", res.FeatureID),
			slog.Int64("tagID", res.TagID),
			sl.Err(err),
		)
	}
}
