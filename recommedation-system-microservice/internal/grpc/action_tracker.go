package recommendation

import (
	"context"
	"fmt"
	"rec-system-microservice/internal/domain/models"
	recommendationpb "rec-system-microservice/internal/pb/recommendation"
	recs "rec-system-microservice/internal/pb/recommendation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserActionTracker interface {
	TrackUserAction(
		ctx context.Context,
		event models.UserActionEvent,
	) error
}

func (s serverAPI) TrackUserAction(
	ctx context.Context,
	req *recs.TrackUserActionRequest,
) (*recs.TrackUserActionResponse, error) {

	if req.GetUserId() == "" || req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	action, err := mapUserAction(req.GetAction())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong action type")
	}

	event := models.UserActionEvent{
		UserID:      req.GetUserId(),
		ContentID:   req.GetContentId(),
		Action:      action,
		Rating:      nil,
		DurationSec: nil,
	}

	err = s.actions.TrackUserAction(
		ctx,
		event,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.TrackUserActionResponse{}, nil
}

func mapUserAction(
	action recommendationpb.UserAction,
) (models.ActionType, error) {

	switch action {
	case recommendationpb.UserAction_VIEW:
		return models.ActionView, nil
	case recommendationpb.UserAction_LIKE:
		return models.ActionLike, nil
	case recommendationpb.UserAction_DISLIKE:
		return models.ActionDislike, nil
	case recommendationpb.UserAction_RATE:
		return models.ActionRate, nil
	case recommendationpb.UserAction_ADD_TO_FAVORITES:
		return models.ActionFavorite, nil
	default:
		return "", fmt.Errorf("unknown action: %v", action)
	}
}
