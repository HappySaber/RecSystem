package app

import (
	grpcapp "catalog-microservice/internal/app/grpc"
	catalog "catalog-microservice/internal/services"
	"catalog-microservice/internal/storage/postgresql"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
) *App {
	storage, err := postgresql.New()
	if err != nil {
		panic(err)
	}
	catalogService := catalog.New(log, storage)
	grpcApp := grpcapp.New(log, grpcPort, catalogService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
