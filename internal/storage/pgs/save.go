package pgs

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"banners-management/internal/model/entity"
	"banners-management/internal/storage/pgs/common/bannertag"
	"banners-management/internal/storage/repo"
)

// SaveBanner saves a banner to the database.
// It returns the ID of the common banner if successful, otherwise error.
func (s *Storage) SaveBanner(ctx context.Context, b *entity.Banner) (bannerID int64, err error) {
	const comp = "storage.pgs.SaveBanner"

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", comp, err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			bannerID, err = 0, fmt.Errorf("%s: %w", comp, err)
		}
	}()

	row := tx.QueryRow(
		ctx,
		`INSERT INTO Banner (title, text, url, is_active, feature_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
		b.Title,
		b.Text,
		b.URL,
		b.IsActive,
		b.FeatureID,
		b.CreatedAt,
		b.UpdatedAt,
	)

	err = row.Scan(&bannerID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", comp, err)
	}

	q := bannertag.InsertTagsQuery(bannerID, b.TagIDs)
	_, err = tx.Exec(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", comp, err)
	}

	err = tx.Commit(ctx)
	pgErr := new(pgconn.PgError)
	if errors.As(err, &pgErr) && pgErr.Code == "P0001" { // P0001 when trigger is fired
		return 0, fmt.Errorf("%s: %w", comp, repo.ErrBannerAlreadyExists)
	} else if err != nil {
		return 0, fmt.Errorf("%s: %w", comp, err)
	}

	return bannerID, nil
}
