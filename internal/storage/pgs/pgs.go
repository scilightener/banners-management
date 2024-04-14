package pgs

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func New(ctx context.Context, connectionString string) (*Storage, error) {
	const comp = "storage.pgs.New"

	dbPool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", comp, err)
	}

	return &Storage{dbPool: dbPool}, nil
}

func (s *Storage) Close(_ context.Context) error {
	s.dbPool.Close()
	return nil
}
