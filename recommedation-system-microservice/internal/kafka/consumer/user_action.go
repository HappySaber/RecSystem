package consumer

import (
	"context"
	"encoding/json"
	"rec-system-microservice/internal/domain/models"
	useractions "rec-system-microservice/internal/services/user_actions"
)

func UserActionHandler(service *useractions.UserActionTracker) MessageHandler {
	return func(ctx context.Context, msg []byte) error {
		var event models.UserActionEvent

		if err := json.Unmarshal(msg, &event); err != nil {
			return err
		}

		return service.TrackUserAction(ctx, event)
	}
}
