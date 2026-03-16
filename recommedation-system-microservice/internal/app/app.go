package app

import (
	"context"
	"log"
	"log/slog"
	grpcapp "rec-system-microservice/internal/app/grpc"
	"rec-system-microservice/internal/kafka/consumer"
	contentgenre "rec-system-microservice/internal/services/content_genre"
	recommendation "rec-system-microservice/internal/services/recommendation"
	"rec-system-microservice/internal/services/recommendation/airecommendation"
	useractions "rec-system-microservice/internal/services/user_actions"
	userpreferences "rec-system-microservice/internal/services/user_preferences"
	"rec-system-microservice/internal/storage/postgresql"
	"time"

	"github.com/segmentio/kafka-go"
)

type App struct {
	GRPCSrv  *grpcapp.App
	Consumer []consumer.Consumer
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
	aiRecService := initAIRecommendationService(log, storage)

	engine := recommendation.New(log, storage, 0)
	prefs := &userpreferences.UserPreferences{}
	actions := useractions.New(log, storage)

	contentGenre := contentgenre.New(log, storage)

	grpcApp := grpcapp.New(log, grpcPort, engine, prefs, actions, aiRecService)

	userActionReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "user-action",
		GroupID: "recommendation",
	})

	genreReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "content-genre",
		GroupID: "recommendation",
	})

	userConsumer := consumer.NewConsumer(
		userActionReader,
		consumer.UserActionHandler(actions),
	)

	genreConsumer := consumer.NewConsumer(
		genreReader,
		consumer.ContentGenreHandler(contentGenre),
	)

	consumers := []consumer.Consumer{
		userConsumer,
		genreConsumer,
	}

	//uaConsumer := consumer.NewUserActionConsumer(reader, actions)
	return &App{
		GRPCSrv:  grpcApp,
		Consumer: consumers,
	}
}

func initAIRecommendationService(log *slog.Logger, storage *postgresql.Storage) *airecommendation.AIRecommendation {
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
	return airecommendation.NewAIRecommendation(log, aiClient, promptBuilder, aiParser, storage)
}

func (a *App) StartConsumers(ctx context.Context) {
	for _, c := range a.Consumer {
		go func(cons consumer.Consumer) {
			for {
				if err := cons.Start(ctx); err != nil {
					log.Println("consumer stopped:", err)
					time.Sleep(3 * time.Second)
					continue
				}
			}
		}(c)
	}
}
