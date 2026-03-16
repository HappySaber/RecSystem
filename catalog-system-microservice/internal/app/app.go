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

// New — для продакшна, DSN читается из env
func New(log *slog.Logger, grpcPort int) *App {
	storage, err := postgresql.New()
	if err != nil {
		panic(err)
	}
	return NewApp(log, grpcPort, storage)
}

// NewWithDSN — для тестов, DSN передаётся снаружи от testcontainers
func NewWithDSN(log *slog.Logger, grpcPort int, dsn string) *App {
	storage, err := postgresql.NewWithDSN(dsn)
	if err != nil {
		panic(err)
	}
	return NewApp(log, grpcPort, storage)
}

// внутренний конструктор, общий для New и NewWithDSN
func NewApp(log *slog.Logger, grpcPort int, storage *postgresql.Storage) *App {
	catalogService := catalog.New(log, storage)
	grpcApp := grpcapp.New(log, grpcPort, catalogService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
