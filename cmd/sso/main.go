package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/m/internal/app"
	"sso/m/internal/config"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Debug("config", slog.Any("config", cfg))

	application := app.New(logger, cfg.GRPC.Port, cfg.Postgres, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stop

	logger.Info("sso stopping", slog.String("signal", sig.String()))
	application.GRPCServer.Stop();

	logger.Info("sso stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	case "development":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	case "production":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	}
	return logger
}