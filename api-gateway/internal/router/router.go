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
	SSO            *handlers.SSOHandler
}

func NewRouter(h Handlers, mw *middleware.Manager) http.Handler {
	r := mux.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.RequestID)

	r.HandleFunc("/auth/register", h.SSO.Register).Methods(http.MethodPost)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(mw.Auth)

	api.HandleFunc(
		"/recommendations/explicit",
		h.Recommendation.GetExplicit,
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/events/user-action",
		h.UserAction.Track,
	).Methods(http.MethodPost)

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
