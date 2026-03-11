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

func (c *RecommendationClient) GetRecommendations(
	ctx context.Context,
	userID string,
	limit int,
) ([]string, error) {

	resp, err := c.client.GetRecommendations(
		ctx,
		&pbRecs.GetRecommendationsRequest{
			UserId: userID,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.ContentIds, nil
}

func (c *RecommendationClient) GetRecommendationsByGenres(
	ctx context.Context,
	genres []string,
	limit int,
) ([]string, error) {

	resp, err := c.client.GetRecommendationsByGenres(
		ctx,
		&pbRecs.GetRecommendationsByGenresRequest{
			Genres: genres,
			Limit:  int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.ContentIds, nil
}

func (c *RecommendationClient) GetSimilarContent(
	ctx context.Context,
	contentID string,
	limit int,
) ([]string, error) {

	resp, err := c.client.GetSimilarContent(
		ctx,
		&pbRecs.GetSimilarContentRequest{
			ContentId: contentID,
			Limit:     int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.ContentIds, nil
}

func (c *RecommendationClient) GetTrendingContent(
	ctx context.Context,
	limit int,
) ([]string, error) {

	resp, err := c.client.GetTrendingContent(
		ctx,
		&pbRecs.GetTrendingContentRequest{
			Limit: int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.ContentIds, nil
}

func (c *RecommendationClient) GetPopularContent(
	ctx context.Context,
	limit int,
) ([]string, error) {

	resp, err := c.client.GetPopularContent(
		ctx,
		&pbRecs.GetPopularContentRequest{
			Limit: int32(limit),
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.ContentIds, nil
}
