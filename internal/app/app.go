package app

import (
	"log/slog"
	gRPCapp "sso/m/internal/app/grpc"
	"sso/m/internal/config"
	"sso/m/internal/services/auth"
	"sso/m/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCServer *gRPCapp.App
}

func New(log *slog.Logger, gRPCPort int, postgresCfg config.PostgresConfig, tokenTTL time.Duration) *App {
	postgres := postgres.New(postgresCfg)
	authService := auth.NewAuth(log, postgres, postgres, postgres, tokenTTL)
	gRPCServer := gRPCapp.New(log, authService, gRPCPort)


	return &App{
		GRPCServer: gRPCServer,
	}
}
