package grpcapp

import (
	"fmt"
	cmgrpc "gRPC_ContactManagement_Service/internal/grpc/cm"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// TODO: add new param CMService when you implement
func New(
	log *slog.Logger,
	port int,
) *App {
	gRPC := grpc.NewServer()
	cmgrpc.Register(gRPC)
	return &App{
		log:        log,
		gRPCServer: gRPC,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcapp.run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("gRPC server is running")

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
