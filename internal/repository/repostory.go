package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func safeRollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}
