package main

import (
	"context"
	"gRPC_ContactManagement_Service/internal/app"
	ssogrpc "gRPC_ContactManagement_Service/internal/clients/sso/grpc"
	"gRPC_ContactManagement_Service/internal/config"
	"gRPC_ContactManagement_Service/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: LOAD CONFIG
	cfg := config.MustLoad()

	// TODO: LOGGER
	log := setupLogger(cfg.Env)
	log.Info("logger setup")

	// TODO: init auth client
	authClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)

	if err != nil {
		panic(err)
	}
	authClientInterceptor := ssogrpc.SSOMiddleware(authClient, cfg.Clients.SSO.AppID)

	// TODO: INIT APP
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, authClientInterceptor)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
