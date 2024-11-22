package app

import (
	grpcapp "gRPC_ContactManagement_Service/internal/app/grpc"
	"gRPC_ContactManagement_Service/internal/service/cm"
	"gRPC_ContactManagement_Service/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	// TODO: init cm service
	cmService := cm.New(log, storage, storage, storage)
	grpcApp := grpcapp.New(log, cmService, port)
	return &App{GRPCSrv: grpcApp}
}
