package app

import (
	"log/slog"
	grpcapp "rec-system-microservice/internal/app/grpc"
	"rec-system-microservice/internal/kafka/consumer"
	recommendation "rec-system-microservice/internal/services/recommendation"
	"rec-system-microservice/internal/services/recommendation/airecommendation"
	useractions "rec-system-microservice/internal/services/user_actions"
	userpreferences "rec-system-microservice/internal/services/user_preferences"
	"rec-system-microservice/internal/storage/postgresql"

	"github.com/segmentio/kafka-go"
)

type App struct {
	GRPCSrv  *grpcapp.App
	Consumer *consumer.UserActionConsumer
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

	engine := recommendation.New(log, nil, 0)
	prefs := &userpreferences.UserPreferences{}
	actions := useractions.New(log, storage)

	grpcApp := grpcapp.New(log, grpcPort, engine, prefs, actions, aiRecService)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-action",
		GroupID: "recommendation",
	})

	uaConsumer := consumer.NewUserActionConsumer(reader, actions)
	return &App{
		GRPCSrv:  grpcApp,
		Consumer: uaConsumer,
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
