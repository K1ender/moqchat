package logger

import (
	"context"
	"log/slog"

	"github.com/K1ender/moqchat/internal/config"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
)

func New(env config.Env) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.Development:
		zaplogger := zap.Must(zap.NewDevelopment())

		handler := slogzap.Option{
			Level:     slog.LevelDebug,
			AddSource: true,
			Logger:    zaplogger,
		}.NewZapHandler()

		logger = slog.New(handler)
	case config.Production:
		zaplogger := zap.Must(zap.NewProduction())

		handler := slogzap.Option{
			Level:     slog.LevelInfo,
			AddSource: true,
			Logger:    zaplogger,
		}.NewZapHandler()

		logger = slog.New(handler)
	}

	slog.SetDefault(logger)

	return logger
}

func L() *slog.Logger {
	return slog.Default()
}

type loggerkey struct{}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	ctx = context.WithValue(ctx, loggerkey{}, logger)
	return ctx
}

func FromContext(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerkey{}).(*slog.Logger)
}
