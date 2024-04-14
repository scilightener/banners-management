package pgs

import (
	"avito-test-task/internal/storage"
	"context"
	"fmt"
)

func (s *Storage) DeleteBanner(ctx context.Context, id int64) error {
	const comp = "storage.pgs.DeleteBanner"

	r, err := s.dbPool.Exec(ctx, `DELETE FROM banner WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", comp, err)
	}

	if r.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", comp, storage.ErrBannerNotFound)
	}

	return nil
}
