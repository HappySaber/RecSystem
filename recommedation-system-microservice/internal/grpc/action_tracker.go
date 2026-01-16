package recommendation

import (
	"context"
	recs "rec-system-microservice/internal/pb/recommendation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserActionTracker interface {
	TrackUserAction(
		ctx context.Context,
		userID string,
		contentID string,
		action recs.UserAction,
		rating *int32,
		duration *int32,
	) error
}

//
// user actions
//

func (s serverAPI) TrackUserAction(
	ctx context.Context,
	req *recs.TrackUserActionRequest,
) (*recs.TrackUserActionResponse, error) {

	if req.GetUserId() == "" || req.GetContentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "wrong arguments")
	}

	err := s.actions.TrackUserAction(
		ctx,
		req.GetUserId(),
		req.GetContentId(),
		req.GetAction(),
		req.Rating,
		req.DurationSec,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &recs.TrackUserActionResponse{}, nil
}
