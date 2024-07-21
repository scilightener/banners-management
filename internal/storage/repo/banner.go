package repo

import (
	"context"

	"avito-test-task/internal/model/entity"
)

// BannerSaver is an interface that supports banner creating.
type BannerSaver interface {
	SaveBanner(ctx context.Context, banner *entity.Banner) (int64, error)
}

// BannerReader is an interface that supports retrieving banners by featureID and/or tagID.
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

// BannerDeleter is an interface that supports deleting banners by id and by featureID and tagID.
type BannerDeleter interface {
	DeleteBanner(ctx context.Context, bannerID int64) error
	DeleteByFeatureTag(ctx context.Context, featureID, tagID int64) error
}

// BannerUpdater is an interface that supports updating banners.
type BannerUpdater interface {
	UpdateBanner(ctx context.Context, banner *entity.UpdatableBanner) error
}
