package app

import (
	"log/slog"
	grpcapp "recommedation-system-microservice/internal/app/grpc"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	storage, err := postgresql.New()
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
