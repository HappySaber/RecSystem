package handlers

import (
	"api-gateway/internal/dto"
	"api-gateway/internal/middleware"
	"context"
	"encoding/json"
	"net/http"
)

type UserActionHandler struct {
	Producer         EventProducer
	UserActionClient UserActionClient
}

type UserActionClient interface {
	TrackUserAction(ctx context.Context, userID, contentID, action string, rating, duration *int32) error
}

func (uah *UserActionHandler) TrackUserAction(w http.ResponseWriter, r *http.Request) {
	var req dto.TrackUserActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := middleware.GetUserID(r.Context())
	err := uah.UserActionClient.TrackUserAction(r.Context(), userID, req.ContentID, req.Action, req.Rating, req.Duration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(err)

}
