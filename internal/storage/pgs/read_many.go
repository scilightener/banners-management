package pgs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"avito-test-task/internal/model/entity"
)

// BannersByFeatureTag returns slice of banners associated with given feature and tag.
// It respects the limit and offset parameters, if provided,
// where limit is the maximum number of banners to return and
// offset is the number of banners to skip.
// All the parameters are optional. If they're set to nil, they're ignored.
func (s *Storage) BannersByFeatureTag(
	ctx context.Context,
	featureID, tagID *int64,
	limit, offset *int,
	_ *bool,
) ([]*entity.Banner, error) {
	const comp = "storage.pgs.BannersByFeatureTag"

	q, args := buildReadManyQuery(featureID, tagID, limit, offset)

	rows, err := s.dbPool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", comp, err)
	}

	defer rows.Close()

	banners := make([]*entity.Banner, 0)
	buf := new(entity.Banner)
	for rows.Next() {
		var tagID int64
		err := rows.Scan(
			&buf.ID,
			&buf.Title,
			&buf.Text,
			&buf.URL,
			&buf.IsActive,
			&buf.FeatureID,
			&tagID,
			&buf.CreatedAt,
			&buf.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", comp, err)
		}
		if len(banners) > 0 && buf.ID == banners[len(banners)-1].ID {
			banners[len(banners)-1].TagIDs = append(banners[len(banners)-1].TagIDs, tagID)
		} else {
			buf.TagIDs = append(buf.TagIDs, tagID)
			banners = append(banners, buf)
			buf = new(entity.Banner)
		}
	}

	return banners, nil
}

// buildReadManyQuery builds a sql query based on the provided parameters.
// It returns the query string and the arguments to be passed to the query.
// If a parameter is nil, it's ignored.
func buildReadManyQuery(featureID, tagID *int64, limit, offset *int) (string, []any) {
	var sb strings.Builder

	sb.WriteString(`WITH banners AS (`)
	q, args := getBannersQuery(featureID, tagID, limit, offset)
	sb.WriteString(q)
	sb.WriteString(`) SELECT id, title, text, url, is_active, feature_id, tag_id, created_at, updated_at
			FROM banners JOIN banner_tag bt ON banners.id = bt.banner_id
			ORDER BY id, tag_id;`)

	return sb.String(), args
}

func getBannersQuery(featureID, tagID *int64, limit, offset *int) (string, []any) {
	var (
		args = make([]any, 0, 4)
		sb   strings.Builder
	)

	sb.WriteString(`SELECT id, title, text, url, is_active, feature_id, created_at, updated_at FROM banner b`)

	if tagID != nil {
		sb.WriteString(` JOIN banner_tag bt ON b.id = bt.banner_id`)
	}

	if featureID != nil {
		sb.WriteString(" WHERE b.feature_id = $")
		sb.WriteString(strconv.Itoa(len(args)+1) + " ")
		args = append(args, *featureID)
	}

	if tagID != nil {
		if featureID != nil {
			sb.WriteString(" AND ")
		} else {
			sb.WriteString(" WHERE ")
		}
		sb.WriteString(" bt.tag_id = $")
		sb.WriteString(strconv.Itoa(len(args)+1) + " ")
		args = append(args, *tagID)
	}

	if limit != nil {
		sb.WriteString(" LIMIT $")
		sb.WriteString(strconv.Itoa(len(args)+1) + " ")
		args = append(args, *limit)
	}

	if offset != nil {
		sb.WriteString(" OFFSET $")
		sb.WriteString(strconv.Itoa(len(args)+1) + " ")
		args = append(args, *offset)
	}

	return sb.String(), args
}
