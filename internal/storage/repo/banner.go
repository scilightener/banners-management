package repo

import (
	"context"

	"avito-test-task/internal/models/entity"
)

type BannerSaver interface {
	SaveBanner(ctx context.Context, banner *entity.Banner) (int64, error)
}

type BannerReader interface {
	BannerByFeatureTag(
		ctx context.Context,
		featureID, tagID int64,
		useLastRevision bool,
	) (*entity.Banner, error)

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
