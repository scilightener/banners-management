package banner

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"avito-test-task/internal/cache/redis"
	"avito-test-task/internal/lib/logger/sl"
	"avito-test-task/internal/models/entity"
	"avito-test-task/internal/storage/repo"
)

const (
	CacheTTL = 5 * time.Minute

	cacheSetOpTimeout = 20 * time.Second
)

type CacheKey struct {
	featureID, tagID int64
}

func (ck CacheKey) ToRedisKeyFormat() string {
	return strconv.FormatInt(ck.featureID, 10) + ":" + strconv.FormatInt(ck.tagID, 10)
}

type CacheReader struct {
	reader repo.BannerReader
	cache  *redis.Cache
	logger *slog.Logger
}

func NewCacheReader(reader repo.BannerReader, cache *redis.Cache, logger *slog.Logger) *CacheReader {
	return &CacheReader{
		reader: reader,
		cache:  cache,
		logger: logger,
	}
}

func (cbr *CacheReader) BannersByFeatureTag(
	ctx context.Context,
	featureID, tagID *int64,
	limit, offset *int,
	useLastRevision *bool,
) ([]*entity.Banner, error) {
	return cbr.reader.BannersByFeatureTag(ctx, featureID, tagID, limit, offset, useLastRevision)
}

func (cbr *CacheReader) BannerByFeatureTag(
	ctx context.Context,
	featureID, tagID int64,
	useLastRevision bool,
) (*entity.Banner, error) {
	const comp = "service.banner.cached_banner.BannerByFeatureTag"
	log := cbr.logger.With(slog.String("comp", comp))
	if useLastRevision {
		return cbr.reader.BannerByFeatureTag(ctx, featureID, tagID, useLastRevision)
	}

	key := CacheKey{featureID, tagID}.ToRedisKeyFormat()
	v, err := redis.Get[*entity.Banner](cbr.cache, ctx, key)
	if err != nil {
		log.Error("redis cache get error", sl.Err(err), slog.String("key", key))
		return cbr.getDataUpdateCache(ctx, featureID, tagID, useLastRevision)
	}

	switch v.Status {
	case redis.StatusExists:
		return v.Value, nil
	case redis.StatusNotFound:
		return cbr.getDataUpdateCache(ctx, featureID, tagID, useLastRevision)
	case redis.StatusNotExists:
		return nil, repo.ErrBannerNotFound
	}

	return v.Value, nil
}

func (cbr *CacheReader) getDataUpdateCache(
	ctx context.Context,
	featureID, tagID int64,
	useLastRevision bool,
) (*entity.Banner, error) {
	const comp = "service.banner.cached_banner.getDataUpdateCache"
	log := cbr.logger.With(slog.String("comp", comp))
	status := redis.StatusExists
	v, err := cbr.reader.BannerByFeatureTag(ctx, featureID, tagID, useLastRevision)
	if err != nil {
		if errors.Is(err, repo.ErrBannerNotFound) {
			status = redis.StatusNotExists
		} else {
			return v, err
		}
	}
	item := redis.NewCacheItem(v, status)
	key := CacheKey{featureID, tagID}.ToRedisKeyFormat()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheSetOpTimeout)
		defer cancel()
		err = redis.Set(cbr.cache, ctx, key, item, CacheTTL)
		if err != nil {
			log.Error("redis cache set error", sl.Err(err), slog.String("key", key))
		}
	}()

	return v, err
}
