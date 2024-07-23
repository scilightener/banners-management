package banner

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"banners-management/internal/cache/redis"
	"banners-management/internal/lib/logger/sl"
	"banners-management/internal/model/entity"
	"banners-management/internal/storage/repo"
)

const (
	// CacheTTL is time for a single redis.CacheItem to be stored in redis cache.
	CacheTTL = 5 * time.Minute

	cacheSetOpTimeout = 20 * time.Second
)

// CacheKey is a composite redis key.
type CacheKey struct {
	featureID, tagID int64
}

// ToRedisKeyFormat returns a string that can be used as redis key.
func (ck CacheKey) ToRedisKeyFormat() string {
	return strconv.FormatInt(ck.featureID, 10) + ":" + strconv.FormatInt(ck.tagID, 10)
}

// CacheReader is a decorator for repo.BannerReader that caches all recent read results in redis cache.
type CacheReader struct {
	reader repo.BannerReader
	cache  *redis.Cache
	logger *slog.Logger
}

// NewCacheReader returns a new CacheReader instance.
func NewCacheReader(reader repo.BannerReader, cache *redis.Cache, logger *slog.Logger) *CacheReader {
	return &CacheReader{
		reader: reader,
		cache:  cache,
		logger: logger,
	}
}

// BannersByFeatureTag does nothing and just proxies the request to the decorated repo.BannerReader.
func (cbr *CacheReader) BannersByFeatureTag(
	ctx context.Context,
	featureID, tagID *int64,
	limit, offset *int,
	useLastRevision *bool,
) ([]*entity.Banner, error) {
	return cbr.reader.BannersByFeatureTag(ctx, featureID, tagID, limit, offset, useLastRevision)
}

// BannerByFeatureTag checks if requested data is stored in redis, and if not,
// returns a request result from decorated repo.BannerReader and asynchronously updates cache.
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

// getDataUpdateCache retrieves data from the original repo.BannerReader and
// asynchronously updates cache with this data.
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
