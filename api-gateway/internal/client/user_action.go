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

func (c *RecommendationClient) TrackUserAction(
	ctx context.Context,
	userID string,
	contentID string,
	action pbRecs.UserAction,
	rating *int32,
	duration *int32,
) error {

	_, err := c.client.TrackUserAction(ctx, &pbRecs.TrackUserActionRequest{
		UserId:      userID,
		ContentId:   contentID,
		Action:      action,
		Rating:      rating,
		DurationSec: duration,
	})
	fmt.Println("grpc response:")

	if err != nil {
		fmt.Println("grpc error:", err)
		return fmt.Errorf("failed to track user action: %w", err)
	}

	return nil
}
