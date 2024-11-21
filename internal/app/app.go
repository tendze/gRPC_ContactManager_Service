package app

import (
	grpcapp "gRPC_ContactManagement_Service/internal/app/grpc"
	"gRPC_ContactManagement_Service/internal/service/cm"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	storagePath string,
) *App {
	// TODO: init storage
	// TODO: init cm service
	cmService := cm.New()
	grpcApp := grpcapp.New(log, port)
	return &App{GRPCSrv: grpcApp}
}
