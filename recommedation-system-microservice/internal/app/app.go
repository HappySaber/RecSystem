package app

import (
	"log/slog"
	grpcapp "rec-system-microservice/internal/app/grpc"
	"rec-system-microservice/internal/services/recommendation/airecommendation"
	"rec-system-microservice/internal/storage/postgresql"
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
	_ = storage // Placeholder to avoid unused variable error
	aiRecService := initAIRecommendationService(log)

	grpcApp := grpcapp.New(log, grpcPort, aiRecService)
	return &App{
		GRPCSrv: grpcApp,
	}
}

func initAIRecommendationService(log *slog.Logger) *airecommendation.AIRecommendation {
	aiClient, err := airecommendation.NewOpenAIClient()
	if err != nil {
		panic(err)
	}
	promptBuilder, err := airecommendation.NewPromptBuilder()
	if err != nil {
		panic(err)
	}
	aiParser, err := airecommendation.NewAIResponseParser()
	if err != nil {
		panic(err)
	}
	return airecommendation.NewAIRecommendation(log, aiClient, promptBuilder, aiParser)
}
