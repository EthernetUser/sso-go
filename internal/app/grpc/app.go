package grpc

import (
	"fmt"
	"log/slog"
	"net"
	authGRPC "sso/m/internal/grpc/auth"
	"sso/m/internal/services/auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger, 
	authService *auth.Auth,
	port int,
	) *App {
	gRPCServer := grpc.NewServer()

	authGRPC.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	netListener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))

	if err != nil {
		return err
	}

	if err := app.gRPCServer.Serve(netListener); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop() {
	app.gRPCServer.GracefulStop()
}