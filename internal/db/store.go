package db

import (
	"WBL0/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

// New init storage with connections pool
func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode,
	)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("init db storage failed: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db storage failed: %w", err)
	}
	return &Storage{pool: pool}, nil
}

// Close closes all connections in the connections pool
func (s *Storage) Close() {
	s.pool.Close()
}
