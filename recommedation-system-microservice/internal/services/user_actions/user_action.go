package useractions

import (
	"context"
	"log/slog"
	"rec-system-microservice/internal/domain/models"
)

type UserAction struct {
	log       *slog.Logger
	uaTracker UserActionTrackerSaver
}

type UserActionTrackerSaver interface {
	SaveTrackAction(
		ctx context.Context,
		event models.UserActionEvent,
	) error

	SaveTrackActions(
		ctx context.Context,
		events []models.UserActionEvent,
	) error
}

func (ua *UserAction) TrackAction(
	ctx context.Context,
	event models.UserActionEvent,
) error {
	const op = "useractions.TrackUserAction"

	log := ua.log.With(
		slog.String("op", op),
	)
	log.Info("tracking user action",
		slog.String("userID", event.UserID),
		slog.String("actionType", string(event.Action)),
		slog.String("contentID", event.ContentID),
	)

	err := ua.uaTracker.SaveTrackAction(ctx, event)
	if err != nil {
		ua.log.Error("failed to track user action", "error", err.Error())
		return err
	}
	return nil
}
