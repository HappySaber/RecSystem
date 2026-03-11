package handlers

import (
	"api-gateway/internal/middleware"
	"context"
	"net/http"

	pbRecs "api-gateway/internal/pbs/rec-system/recommendation"

	"github.com/gorilla/mux"
)

type UserActionHandler struct {
	Producer         EventProducer
	UserActionClient UserActionClient
}

type UserActionClient interface {
	TrackUserAction(
		ctx context.Context,
		userID string,
		contentID string,
		action pbRecs.UserAction,
		rating *int32,
		duration *int32,
	) error
}

// func (uah *UserActionHandler) TrackUserAction(w http.ResponseWriter, r *http.Request) {
// 	var req dto.TrackUserActionRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	userID, _ := middleware.GetUserID(r.Context())
// 	err := uah.UserActionClient.TrackUserAction(r.Context(), userID, req.ContentID, req.Action, req.Rating, req.Duration)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(err)

// }
func (uah *UserActionHandler) HandleAction(action pbRecs.UserAction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		contentID := vars["content_id"]

		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		err := uah.UserActionClient.TrackUserAction(
			r.Context(),
			userID,
			contentID,
			action,
			nil,
			nil,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
