package repo

import (
	"avito-test-task/internal/models/entity"
	"context"
)

type BannerSaver interface {
	SaveBanner(ctx context.Context, banner *entity.Banner) (int64, error)
}

type BannerReader interface {
	BannersByFeatureTag(
		ctx context.Context,
		featureID, tagID *int64,
		limit, offset *int,
		lastRevision *bool,
	) ([]*entity.Banner, error)
}

type BannerDeleter interface {
	DeleteBanner(ctx context.Context, bannerID int64) error
}

type BannerUpdater interface {
	UpdateBanner(ctx context.Context, banner *entity.UpdatableBanner) error
}
