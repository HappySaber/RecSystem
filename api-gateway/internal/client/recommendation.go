package client

import (
	pbRecs "api-gateway/internal/pbs/rec-system/recommendation"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type RecommendationClient struct {
	client pbRecs.RecommendationServiceClient
}

func NewRecommendationClient(addr string) (*RecommendationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to recommendation service: %w", err)
	}

	c := pbRecs.NewRecommendationServiceClient(conn)
	return &RecommendationClient{
		client: c,
	}, nil
}

func (c *RecommendationClient) GetExplicit(ctx context.Context, query string, limit int) ([]string, error) {
	resp, err := c.client.GetAIRecommendations(
		ctx,
		&pbRecs.GetAIRecommendationsRequest{
			UserId: "",
			Mode:   pbRecs.RecommendationMode_EXPLICIT,
			Query:  &query,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.ContentIds, nil
}
