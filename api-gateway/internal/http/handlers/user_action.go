package handlers

import (
	"api-gateway/internal/events"
	"api-gateway/internal/events/schemas"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserActionHandler struct {
	Producer EventProducer
}

type EventProducer interface {
	Publish(ctx context.Context, topic, key string, event any) error
}

func (uah *UserActionHandler) Track(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID    string `json:"user_id"`
		ContentID string `json:"content_id"`
		Action    string `json:"action"`
		Rating    *int32 `json:"rating"`
		Duration  *int32 `json:"duration_sec"`
	}
	event := schemas.UserActionEvent{
		EventID:   uuid.NewString(),
		UserID:    req.UserID,
		ContentID: req.ContentID,
		Action:    req.Action,
		Rating:    req.Rating,
		Duration:  req.Duration,
		Timestamp: time.Now().UTC(),
		Source:    "api-gateway",
	}

	_ = uah.Producer.Publish(
		ctx,
		events.TopicUserActions,
		req.UserID,
		event,
	)

	w.WriteHeader(http.StatusAccepted)
}
