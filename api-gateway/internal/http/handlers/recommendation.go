package handlers

import (
	"api-gateway/internal/dto"
	"api-gateway/internal/middleware"
	"context"
	"encoding/json"
	"net/http"
)

type RecommendationHandler struct {
	Producer             EventProducer
	RecommendationClient RecommendationClient
}

type RecommendationClient interface {
	GetExplicit(
		ctx context.Context,
		query string,
		limit int,
	) ([]string, error)

	GetRecommendations(
		ctx context.Context,
		userID string,
		limit int,
	) ([]string, error)

	GetRecommendationsByGenres(
		ctx context.Context,
		genres []string,
		limit int,
	) ([]string, error)

	GetSimilarContent(
		ctx context.Context,
		contentID string,
		limit int,
	) ([]string, error)

	GetTrendingContent(
		ctx context.Context,
		limit int,
	) ([]string, error)

	GetPopularContent(
		ctx context.Context,
		limit int,
	) ([]string, error)
}

func (h *RecommendationHandler) GetExplicit(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	var req dto.GetExplicitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// reqID, _ := ctx.Value(middleware.RequestIDKey).(string)

	// searchEvent := schemas.UserSearchEvent{
	// 	EventID:   uuid.NewString(),
	// 	RequestID: reqID,
	// 	UserID:    req.UserID,
	// 	Query:     req.Query,
	// 	Timestamp: time.Now().UTC(),
	// 	Source:    "api-gateway",
	// }

	// if err := h.Producer.Publish(
	// 	ctx,
	// 	events.TopicUserSearch,
	// 	req.UserID,
	// 	searchEvent,
	// ); err != nil {
	// }

	contentIDs, err := h.RecommendationClient.GetExplicit(
		ctx,
		req.Query,
		req.Limit,
	)
	if err != nil {
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		return
	}

	resp := dto.GetExplicitResponse{
		ContentIDs: contentIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *RecommendationHandler) GetRecommendations(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	var req dto.GetRecommendationsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	contentIDs, err := h.RecommendationClient.GetRecommendations(
		ctx,
		userID,
		req.Limit,
	)
	if err != nil {
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		return
	}

	resp := dto.GetRecommendationsResponse{
		ContentIDs: contentIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *RecommendationHandler) GetRecommendationsByGenres(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	var req dto.GetRecommendationsByGenresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Genres) == 0 || req.Limit <= 0 {
		http.Error(w, "missing or invalid fields", http.StatusBadRequest)
		return
	}

	contentIDs, err := h.RecommendationClient.GetRecommendationsByGenres(
		ctx,
		req.Genres,
		req.Limit,
	)
	if err != nil {
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		return
	}

	resp := dto.GetRecommendationsByGenresResponse{
		ContentIDs: contentIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *RecommendationHandler) GetSimilarContent(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	var req dto.GetSimilarContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ContentID == "" || req.Limit <= 0 {
		http.Error(w, "missing or invalid fields", http.StatusBadRequest)
		return
	}

	contentIDs, err := h.RecommendationClient.GetSimilarContent(
		ctx,
		req.ContentID,
		req.Limit,
	)
	if err != nil {
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		return
	}

	resp := dto.GetSimilarContentResponse{
		ContentIDs: contentIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *RecommendationHandler) GetTrendingContent(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	var req dto.GetTrendingContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Limit <= 0 {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	contentIDs, err := h.RecommendationClient.GetTrendingContent(
		ctx,
		req.Limit,
	)
	if err != nil {
		http.Error(w, "failed to get trending content", http.StatusInternalServerError)
		return
	}

	resp := dto.GetTrendingContentResponse{
		ContentIDs: contentIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *RecommendationHandler) GetPopularContent(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	var req dto.GetPopularContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Limit <= 0 {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	contentIDs, err := h.RecommendationClient.GetPopularContent(
		ctx,
		req.Limit,
	)
	if err != nil {
		http.Error(w, "failed to get popular content", http.StatusInternalServerError)
		return
	}

	resp := dto.GetPopularContentResponse{
		ContentIDs: contentIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
