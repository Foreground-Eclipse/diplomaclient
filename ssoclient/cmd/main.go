package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	ssogrpc "ssoclient/inter/clients/sso/grpc"
	"ssoclient/inter/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	ssoClient, err := ssogrpc.New(
		context.Background(),
		cfg.Clients.SSO.Adress,
		log,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {

	}
	isAdminResponse, err := ssoClient.IsAdmin(context.Background(), 1)

	fmt.Println(isAdminResponse)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
