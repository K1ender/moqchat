package database

import (
	"context"
	"fmt"
	"time"

	"github.com/K1ender/moqchat/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(
	ctx context.Context,
	cfg config.Database,
) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}

func dsn(cfg config.Database) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
}
