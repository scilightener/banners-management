package pgs

import (
	"avito-test-task/internal/storage/repo"
	"context"
	"fmt"

	"avito-test-task/internal/model/entity"
)

// BannerByFeatureTag finds a banner by provided featureID and tagID.
func (s *Storage) BannerByFeatureTag(
	ctx context.Context,
	featureID, tagID int64,
	_ bool,
) (*entity.Banner, error) {
	const comp = "storage.pgs.BannerByFeatureTag"

	rows, err := s.dbPool.Query(ctx,
		`WITH banners AS (
				SELECT id, title, text, url, is_active, feature_id, created_at, updated_at 
				FROM banner b JOIN banner_tag bt ON b.id = bt.banner_id 
				WHERE b.feature_id = $1 AND bt.tag_id = $2
			) SELECT id, title, text, url, is_active, feature_id, tag_id, created_at, updated_at
			FROM banners JOIN banner_tag bt ON banners.id = bt.banner_id
			ORDER BY id, tag_id;`,
		featureID, tagID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", comp, err)
	}

	banner := new(entity.Banner)
	tagIDs := make([]int64, 1, 16)
	if !rows.Next() {
		return nil, repo.ErrBannerNotFound
	}
	err = rows.Scan(
		&banner.ID,
		&banner.Title,
		&banner.Text,
		&banner.URL,
		&banner.IsActive,
		&banner.FeatureID,
		&tagIDs[0],
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", comp, err)
	}
	i := 1
	id := new(int64)
	for rows.Next() {
		tagIDs = append(tagIDs, 0)
		err = rows.Scan(
			id,
			nil,
			nil,
			nil,
			nil,
			nil,
			&tagIDs[i],
			nil,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", comp, err)
		}
		if *id != banner.ID {
			return nil, fmt.Errorf("%s: %w", comp, repo.ErrBannerNotUnique)
		}
		i++
	}
	banner.TagIDs = tagIDs

	return banner, nil
}
