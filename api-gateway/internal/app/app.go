package app

import (
	router "api-gateway/internal"
	"api-gateway/internal/client"
	"api-gateway/internal/events/producer"
	"api-gateway/internal/http/handlers"
	"context"
)

type App struct {
	server *Server
}

func New(addr string) (*App, error) {
	// Kafka producer
	kafkaProducer := producer.NewProducer([]string{
		"localhost:9092",
	})

	// gRPC clients
	recommendationClient, err := client.NewRecommendationClient(
		"localhost:50051",
	)
	if err != nil {
		return nil, err
	}

	// Handlers
	recommendationHandler := &handlers.RecommendationHandler{
		Producer:             kafkaProducer,
		RecommendationClient: recommendationClient,
	}

	userActionHandler := &handlers.UserActionHandler{
		Producer: kafkaProducer,
	}

	// Router
	httpRouter := router.NewRouter(
		router.Handlers{
			Recommendation: recommendationHandler,
			UserAction:     userActionHandler,
		},
	)

	// HTTP server
	server := NewServer(addr, httpRouter)

	return &App{
		server: server,
	}, nil
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
