package app

import (
	grpcapp "gRPC_ContactManagement_Service/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage
	// TODO: init cm service
	grpcApp := grpcapp.New(log, port)
	return &App{GRPCSrv: grpcApp}
}
