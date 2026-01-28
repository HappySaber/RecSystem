package router

import (
	"api-gateway/internal/http/handlers"
	"api-gateway/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	Recommendation *handlers.RecommendationHandler
	UserAction     *handlers.UserActionHandler
}

func NewRouter(h Handlers) http.Handler {
	r := mux.NewRouter()

	// global middleware
	r.Use(middleware.Recover)
	r.Use(middleware.RequestID)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc(
		"/recommendations/explicit",
		h.Recommendation.GetExplicit,
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/events/user-action",
		h.UserAction.Track,
	).Methods(http.MethodPost)

	// healthcheck
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
