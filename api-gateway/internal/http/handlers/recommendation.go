package handlers

import (
	"api-gateway/internal/dto"
	"api-gateway/internal/events"
	"api-gateway/internal/events/schemas"
	"api-gateway/internal/middleware"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	if req.UserID == "" || req.Query == "" || req.Limit <= 0 {
		http.Error(w, "missing or invalid fields", http.StatusBadRequest)
		return
	}

	reqID, _ := ctx.Value(middleware.RequestIDKey).(string)

	searchEvent := schemas.UserSearchEvent{
		EventID:   uuid.NewString(),
		RequestID: reqID,
		UserID:    req.UserID,
		Query:     req.Query,
		Timestamp: time.Now().UTC(),
		Source:    "api-gateway",
	}

	if err := h.Producer.Publish(
		ctx,
		events.TopicUserSearch,
		req.UserID,
		searchEvent,
	); err != nil {
	}

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
