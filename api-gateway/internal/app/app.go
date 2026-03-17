package app

import (
	"api-gateway/internal/client"
	"api-gateway/internal/events/adapter"
	"api-gateway/internal/events/producer"
	"api-gateway/internal/http/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/router"
	"context"
	"os"
)

type App struct {
	server *Server
}

func New(addr string) (*App, error) {
	// читаем адреса из env — в Docker это имена сервисов из compose
	// локально — localhost с портами
	kafkaBroker := os.Getenv("KAFKA_BROKERS")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	recAddr := os.Getenv("REC_GRPC_ADDR")
	if recAddr == "" {
		recAddr = "localhost:44045"
	}

	catalogAddr := os.Getenv("CATALOG_GRPC_ADDR")
	if catalogAddr == "" {
		catalogAddr = "localhost:44044"
	}

	ssoAddr := os.Getenv("SSO_GRPC_ADDR")
	if ssoAddr == "" {
		ssoAddr = "localhost:44043"
	}

	kafkaProducer := producer.NewProducer([]string{kafkaBroker})
	userActionProd := adapter.NewUserActionProducer(kafkaProducer, "user-action")

	recommendationClient, err := client.NewRecommendationClient(recAddr)
	if err != nil {
		return nil, err
	}

	catalogClient, err := client.NewCatalogClient(catalogAddr)
	if err != nil {
		return nil, err
	}

	ssoClient, err := client.NewSSOClient(ssoAddr)
	if err != nil {
		return nil, err
	}

	// Handlers
	recommendationHandler := &handlers.RecommendationHandler{
		Producer:             userActionProd,
		RecommendationClient: recommendationClient,
	}

	userActionHandler := &handlers.UserActionHandler{
		Producer:         userActionProd,
		UserActionClient: recommendationClient,
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

	mw := middleware.NewManager(os.Getenv("JWT_SECRET_KEY"))

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
