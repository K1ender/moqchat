package api

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/K1ender/moqchat/internal/config"
	"github.com/K1ender/moqchat/internal/logger"
	httptransport "github.com/K1ender/moqchat/internal/transport/http"
)

const shutdownTimeout = 10 * time.Second

func Run(ctx context.Context) error {
	cfg := config.MustInit()
	ctx = logger.WithContext(ctx, logger.New(cfg.Env))
	log := logger.FromContext(ctx)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.InfoContext(ctx, "Starting API...")

	srv := httptransport.NewServer(cfg.HTTP)

	go func() {
		err := srv.Run()
		if err != nil {
			log.ErrorContext(ctx, "API server stopped")
		}
	}()

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer shutdownCancel()

	log.InfoContext(shutdownCtx, "Shutting down API...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.ErrorContext(shutdownCtx, "API server shutdown failed")
	}

	return nil
}
