package client

import (
	pbRecs "api-gateway/internal/pbs/rec-system/recommendation"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type UserActionClient struct {
	client pbRecs.RecommendationServiceClient
}

func NewUserActionClient(addr string) (*UserActionClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user action service: %w", err)
	}

	c := pbRecs.NewRecommendationServiceClient(conn)
	return &UserActionClient{
		client: c,
	}, nil
}

func (c *UserActionClient) TrackUserAction(
	ctx context.Context,
	userID string,
	contentID string,
	action string,
	rating *int32,
	duration *int32,
) error {

	act, err := parseUserAction(action)
	if err != nil {
		return err
	}

	_, err = c.client.TrackUserAction(ctx, &pbRecs.TrackUserActionRequest{
		UserId:      userID,
		ContentId:   contentID,
		Action:      act,
		Rating:      rating,
		DurationSec: duration,
	})
	if err != nil {
		return fmt.Errorf("failed to track user action: %w", err)
	}
	return nil
}

func parseUserAction(a string) (pbRecs.UserAction, error) {
	v, ok := pbRecs.UserAction_value[a]
	if !ok {
		return pbRecs.UserAction_VIEW, fmt.Errorf("invalid action: %s", a)
	}
	return pbRecs.UserAction(v), nil
}
