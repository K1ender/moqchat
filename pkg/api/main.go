package api

import (
	"context"

	"github.com/K1ender/moqchat/internal/config"
	"github.com/K1ender/moqchat/internal/logger"
)

func Run(ctx context.Context) error {
	cfg := config.MustInit()
	ctx = logger.WithContext(ctx, logger.New(cfg.Env))
	log := logger.FromContext(ctx)

	log.InfoContext(ctx, "Starting API...")

	return nil
}
