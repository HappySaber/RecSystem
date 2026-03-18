// internal/app/app.go
package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	grpcapp "rec-system-microservice/internal/app/grpc"
	"rec-system-microservice/internal/kafka/consumer"
	contentgenre "rec-system-microservice/internal/services/content_genre"
	recommendation "rec-system-microservice/internal/services/recommendation"
	"rec-system-microservice/internal/services/recommendation/airecommendation"
	useractions "rec-system-microservice/internal/services/user_actions"
	userpreferences "rec-system-microservice/internal/services/user_preferences"
	"rec-system-microservice/internal/storage/postgresql"

	"github.com/segmentio/kafka-go"
)

type App struct {
	GRPCSrv  *grpcapp.App
	Consumer []consumer.Consumer
}

// New — для продакшна
func New(log *slog.Logger, grpcPort int) *App {
	storage, err := postgresql.New()
	if err != nil {
		panic(err)
	}

	aiRecService := initAIRecommendationService(log, storage)
	return newApp(log, grpcPort, storage, aiRecService)
}

// NewWithDSN — для тестов, принимает DSN от testcontainers
// Kafka консьюмеры не запускаются — в тестах не нужны
func NewWithDSN(log *slog.Logger, grpcPort int, dsn string) *App {
	storage, err := postgresql.NewWithDSN(dsn)
	if err != nil {
		panic(err)
	}

	// используем заглушку вместо реального AI клиента
	// чтобы не требовать HF_API_KEY в тестах
	aiRecService := initAIRecommendationServiceStub(log, storage)
	return newApp(log, grpcPort, storage, aiRecService)
}

// newApp — общий конструктор
func newApp(
	log *slog.Logger,
	grpcPort int,
	storage *postgresql.Storage,
	aiRecService *airecommendation.AIRecommendation,
) *App {
	engine := recommendation.New(log, storage, 0)
	prefs := &userpreferences.UserPreferences{}
	actions := useractions.New(log, storage)
	contentGenre := contentgenre.New(log, storage)

	grpcApp := grpcapp.New(log, grpcPort, engine, prefs, actions, aiRecService)

	// читаем брокеры из env — в Docker это kafka:9092, локально localhost:9092
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	brokers := []string{kafkaBrokers}

	userActionReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "user-action",
		GroupID:     "recommendation",
		StartOffset: kafka.FirstOffset,
	})

	genreReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "content-genre",
		GroupID:     "recommendation",
		StartOffset: kafka.FirstOffset,
	})

	userConsumer := consumer.NewConsumer(
		userActionReader,
		consumer.UserActionHandler(actions),
	)

	genreConsumer := consumer.NewConsumer(
		genreReader,
		consumer.ContentGenreHandler(contentGenre),
	)

	return &App{
		GRPCSrv:  grpcApp,
		Consumer: []consumer.Consumer{userConsumer, genreConsumer},
	}
}

// initAIRecommendationService — для продакшна, требует HF_API_KEY
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

// initAIRecommendationServiceStub — для тестов, без реального API ключа
func initAIRecommendationServiceStub(log *slog.Logger, storage *postgresql.Storage) *airecommendation.AIRecommendation {
	promptBuilder, _ := airecommendation.NewPromptBuilder()
	aiParser, _ := airecommendation.NewAIResponseParser()

	// заглушка клиента — всегда возвращает пустой список
	stubClient := airecommendation.NewStubAIClient()

	return airecommendation.NewAIRecommendation(log, stubClient, promptBuilder, aiParser, storage)
}

func (a *App) StartConsumers(ctx context.Context) {
	log.Printf("starting %d consumers", len(a.Consumer))
	for i, c := range a.Consumer {
		go func(idx int, cons consumer.Consumer) {
			log.Printf("consumer %d starting", idx)
			for {
				if err := cons.Start(ctx); err != nil {
					log.Printf("consumer %d stopped: %v", idx, err)
					time.Sleep(3 * time.Second)
					continue
				}
			}
		}(i, c)
	}
}
