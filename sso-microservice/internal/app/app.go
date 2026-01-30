package app

import (
	"log/slog"
	grpcapp "sso-microservice/internal/app/grpc"
	"sso-microservice/internal/config"
	"sso-microservice/internal/kafka/producer"
	"sso-microservice/internal/services/auth"
	"sso-microservice/internal/storage/postgresql"
	"time"
)

type App struct {
	GRPCSrv  *grpcapp.App
	Producer *producer.KafkaProducer
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
	kafkaConfig config.KafkaProducerConfig,
) *App {
	storage, err := postgresql.New()
	if err != nil {
		panic(err)
	}

	producer := producer.NewKafkaProducer(kafkaConfig)

	authService := auth.New(log, storage, storage, storage, tokenTTL, producer)

	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCSrv:  grpcApp,
		Producer: producer,
	}
}

func (a *App) Stop() {
	_ = a.Producer.Close()
}
