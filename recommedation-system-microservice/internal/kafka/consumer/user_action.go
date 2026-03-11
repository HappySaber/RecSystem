package consumer

import (
	"context"
	"encoding/json"
	"log"
	"rec-system-microservice/internal/domain/models"
	useractions "rec-system-microservice/internal/services/user_actions"
)

func UserActionHandler(service *useractions.UserActionTracker) MessageHandler {
	return func(ctx context.Context, msg []byte) error {
		var event models.UserActionEvent
		log.Printf("RAW EVENT: %s", string(msg))
		if err := json.Unmarshal(msg, &event); err != nil {
			return err
		}
		log.Printf("PARSED EVENT: %+v", event)

		return service.TrackUserAction(ctx, event)
	}
}
