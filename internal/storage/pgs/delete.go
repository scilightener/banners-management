package pgs

import (
	"context"
	"fmt"

	"banners-management/internal/storage/repo"
)

// DeleteBanner deletes banner by id.
func (s *Storage) DeleteBanner(ctx context.Context, id int64) error {
	const comp = "storage.pgs.DeleteBanner"

	r, err := s.dbPool.Exec(ctx, `DELETE FROM banner WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", comp, err)
	}

	if r.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", comp, repo.ErrBannerNotFound)
	}

	return nil
}

// DeleteByFeatureTag deletes banner by featureID and tagID.
func (s *Storage) DeleteByFeatureTag(ctx context.Context, featureID, tagID int64) error {
	const comp = "storage.pgs.DeleteByFeatureTag"

	banner, err := s.BannerByFeatureTag(ctx, featureID, tagID, true)
	if err != nil {
		return fmt.Errorf("%s: %w", comp, err)
	}

	err = s.DeleteBanner(ctx, banner.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", comp, err)
	}

	return nil
}
