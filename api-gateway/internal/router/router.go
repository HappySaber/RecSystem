package router

import (
	"api-gateway/internal/http/handlers"
	"api-gateway/internal/middleware"
	"net/http"

	pbRecs "api-gateway/internal/pbs/rec-system/recommendation"

	"github.com/gorilla/mux"
)

type Handlers struct {
	Recommendation *handlers.RecommendationHandler
	UserAction     *handlers.UserActionHandler
	SSO            *handlers.SSOHandler
	Catalog        *handlers.CatalogHandler
}

func NewRouter(h Handlers, mw *middleware.Manager) http.Handler {
	r := mux.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.RequestID)

	r.HandleFunc("/auth/register", h.SSO.Register).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", h.SSO.Login).Methods(http.MethodPost)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(mw.Auth)

	// Recommendation endpoints

	api.HandleFunc(
		"/recommendations",
		h.Recommendation.GetRecommendations,
	).Methods(http.MethodGet)

	api.HandleFunc(
		"/recommendations/explicit",
		h.Recommendation.GetExplicit,
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/recommendations/genres",
		h.Recommendation.GetRecommendationsByGenres,
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/recommendations/similar",
		h.Recommendation.GetSimilarContent,
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/recommendations/trending",
		h.Recommendation.GetTrendingContent,
	).Methods(http.MethodGet)

	api.HandleFunc(
		"/recommendations/popular",
		h.Recommendation.GetPopularContent,
	).Methods(http.MethodGet)

	// Catalog

	api.HandleFunc(
		"/catalog/content/{content_id}",
		h.Catalog.GetContentDetails,
	).Methods(http.MethodGet)

	api.HandleFunc(
		"/catalog/get_content",
		h.Catalog.GetContent,
	).Methods(http.MethodGet)

	api.HandleFunc(
		"/catalog/get-content-by-external",
		h.Catalog.FindContentByExternal,
	).Methods(http.MethodGet)

	// User actions

	api.HandleFunc(
		"/content/{content_id}/like",
		h.UserAction.HandleAction(pbRecs.UserAction_LIKE),
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/content/{content_id}/dislike",
		h.UserAction.HandleAction(pbRecs.UserAction_DISLIKE),
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/content/{content_id}/favorite",
		h.UserAction.HandleAction(pbRecs.UserAction_ADD_TO_FAVORITES),
	).Methods(http.MethodPost)

	api.HandleFunc(
		"/content/{content_id}/view",
		h.UserAction.HandleAction(pbRecs.UserAction_VIEW),
	).Methods(http.MethodPost)

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
