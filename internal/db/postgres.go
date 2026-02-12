package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDb(ctx context.Context) (*pgxpool.Pool, error) {
	dbUrl := "postgres://obscurity:1010@localhost:5432/rating-app"

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping %v", err)
	}

	return pool, nil
}
