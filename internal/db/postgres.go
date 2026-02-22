package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDb(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return pool, nil
}
