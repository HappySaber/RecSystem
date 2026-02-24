package app

import (
	"api-gateway/internal/client"
	"api-gateway/internal/events/adapter"
	"api-gateway/internal/events/producer"
	"api-gateway/internal/http/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/router"
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
	userActionProd := adapter.NewUserActionProducer(kafkaProducer, "user-action")

	// gRPC clients
	recommendationClient, err := client.NewRecommendationClient(
		"localhost:44045",
	)
	if err != nil {
		return nil, err
	}

	catalogClient, err := client.NewCatalogClient(
		"localhost:44044",
	)

	if err != nil {
		return nil, err
	}
	ssoClient, err := client.NewSSOClient(
		"localhost:44043",
	)

	if err != nil {
		return nil, err
	}
	// Handlers
	recommendationHandler := &handlers.RecommendationHandler{
		Producer:             userActionProd,
		RecommendationClient: recommendationClient,
	}

	userActionHandler := &handlers.UserActionHandler{
		Producer: userActionProd,
	}

	ssoHandler := &handlers.SSOHandler{
		//Producer: kafkaProducer,
		SSOClient: ssoClient,
	}

	catalogHandler := &handlers.CatalogHandler{
		Producer:      userActionProd,
		CatalogClient: catalogClient,
	}
	// Middleware
	mw := middleware.NewManager(ssoClient)

	// Router
	httpRouter := router.NewRouter(
		router.Handlers{
			Recommendation: recommendationHandler,
			UserAction:     userActionHandler,
			SSO:            ssoHandler,
			Catalog:        catalogHandler,
		},
		mw,
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
