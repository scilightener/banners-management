package banner

import (
	"context"
	"encoding/json"
	"log/slog"

	goredis "github.com/redis/go-redis/v9"

	"avito-test-task/internal/cache/redis"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/storage/repo"
)

const (
	RedisBannerDeleterByFeatureTagChannelName = "redis_deleter_job_delayer"
)

type redisDeleteMessage struct {
	FeatureID int64 `json:"feature_id"`
	TagID     int64 `json:"tag_id"`
}

type RedisChannelDeleter struct {
	cache   *redis.Cache
	deleter repo.BannerDeleter
	logger  *slog.Logger
}

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

func (r *RedisChannelDeleter) DeleteBanner(ctx context.Context, bannerID int64) error {
	return r.deleter.DeleteBanner(ctx, bannerID)
}

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

func (r *RedisChannelDeleter) runDeleterDaemon(ctx context.Context, ch <-chan *goredis.Message) {
	for m := range ch {
		if m.Channel != RedisBannerDeleterByFeatureTagChannelName {
			continue
		}

		go func() {
			res := new(redisDeleteMessage)
			err := json.Unmarshal([]byte(m.Payload), res)
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
		}()
	}
}
