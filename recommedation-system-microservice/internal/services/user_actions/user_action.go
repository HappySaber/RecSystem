package useractions

import (
	"context"
	"log/slog"
	"rec-system-microservice/internal/domain/models"
)

type UserActionTracker struct {
	log       *slog.Logger
	uaTracker UserActionTrackerSaver
}

type UserActionTrackerSaver interface {
	SaveTrackAction(
		ctx context.Context,
		event models.UserActionEvent,
	) error

	// SaveTrackActions(
	// 	ctx context.Context,
	// 	events []models.UserActionEvent,
	// ) error
}

func New(
	log *slog.Logger,
	uaTracker UserActionTrackerSaver,
) *UserActionTracker {
	return &UserActionTracker{
		log:       log,
		uaTracker: uaTracker,
	}
}

func (ua *UserActionTracker) TrackUserAction(
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

// func (ua *UserAction) TrackActions(
// 	ctx context.Context,
// 	events []models.UserActionEvent,
// ) error {
// 	const op = "useractions.TrackUserActions"

// 	log := ua.log.With(
// 		slog.String("op", op),
// 	)
// 	log.Info("tracking user actions",
// 		slog.Int("numberOfEvents", len(events)),
// 	)
// 	err := ua.uaTracker.SaveTrackActions(ctx, events)
// 	if err != nil {
// 		ua.log.Error("failed to track user actions", "error", err.Error())
// 		return err
// 	}
// 	return nil
// }
