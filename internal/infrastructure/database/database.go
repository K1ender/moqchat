package database

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/K1ender/moqchat/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const timeout = 20 * time.Second

func New(
	ctx context.Context,
	cfg config.Database,
) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := pgxpool.New(ctx, dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		conn.Close()

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}

func dsn(cfg config.Database) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Pass,
		net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		cfg.Name,
	)
}
